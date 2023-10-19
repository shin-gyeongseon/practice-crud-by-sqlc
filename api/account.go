package api

import (
	"database/sql"
	"fmt"
	"go-practice/db/tutorial"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Queries.CreateAccount(ctx, tutorial.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Queries.SelectAccount(ctx, req.ID) // 여기서 ctx 를 넣어줘야하는건지 context.Context 를 넣어줘야하는건지 모르겠네 ;;
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("err address : %v , sql.ErrNoRows address : %v \n", &err, &sql.ErrNoRows)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, account)
		return
	}

	ctx.JSON(http.StatusOK, account)
}
