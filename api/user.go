package api

import (
	"go-practice/db/tutorial"
	"go-practice/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// user go 를 만드는 작업부터 해야합니다. 이거부터 다시 작업합시다.
type createUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	UserName          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user tutorial.User) userResponse {
	return userResponse{
		UserName:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.CreateUser(ctx, tutorial.CreateUserParams{
		Username:       req.UserName,
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		if tutorial.ErrorCode(err) == tutorial.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type loginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	ReturnUser         userResponse `json:"user"`
	AccessToken        string       `json:"token"`
	AccessTokenExpired time.Time    `json:"access_token_expired"`
	SessionID          string       `json:"session_id"`
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var loginRequest loginRequest
	if err := ctx.BindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.SelectUser(ctx, loginRequest.UserName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	err2 := util.CheckPassword(loginRequest.Password, user.HashedPassword)
	if err2 != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err2))
		return
	}

	token, payload, err3 := server.tokenMaker.CreateToken(loginRequest.UserName, server.config.AccessTokenDuration)
	if err3 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err3))
		return
	}

	loginSession := tutorial.CreateSessionParams{
		ID:           payload.ID.String(),
		ExpiresAt:    payload.ExpiredAt,
		RefreshToken: token,
		Username:     loginRequest.UserName,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
	}
	session, err4 := server.store.CreateSession(ctx, loginSession)
	if err4 != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err4))
		return
	}

	loginResponse := loginResponse{
		ReturnUser: userResponse{
			UserName:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
		AccessToken:        token,
		AccessTokenExpired: payload.ExpiredAt,
		SessionID:          session.ID,
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

func (server *Server) ValidUser(ctx *gin.Context, userName string) bool {
	_, err := server.store.SelectUser(ctx, userName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return true
	}

	return false
}
