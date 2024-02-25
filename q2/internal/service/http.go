package service

import (
	"context"
	"net/http"

	"github.com/Yougigun/meepshop_q2/internal/handler"
	"github.com/Yougigun/meepshop_q2/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpService struct {
	Engine *gin.Engine
}

func Build(ctx context.Context, log *zap.Logger, repo *repository.Repository) *http.Server {
	r := gin.Default()
	h := handler.NewAccountHandler(ctx, log, repo)

	r.POST("/accounts", h.CreateAccount)

	r.POST("/accounts/deposit", h.DepositAccount)

	r.POST("/accounts/withdraw", h.WithdrawAccount)

	r.POST("/accounts/transfer", h.TransferAccount)
	{
		// internal api for admin. todo: add auth middleware
		r.GET("/accounts/:id", h.GetAccount)
		r.GET("/transactions", h.GetTransactionLog)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	return srv
}
