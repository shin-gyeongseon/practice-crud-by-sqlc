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
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user tutorial.User) userResponse {
	return userResponse{
		Username:          user.Username,
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
