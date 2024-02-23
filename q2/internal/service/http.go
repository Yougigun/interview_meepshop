package service

import (
	"github.com/Yougigun/meepshop_q2/internal/handler"
	"github.com/Yougigun/meepshop_q2/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpService struct {
	Engine *gin.Engine
}

func Build(log *zap.Logger, ar *repository.Repository) *HttpService {
	r := gin.Default()
	h := handler.NewAccountHandler(log, ar)
	r.POST("/accounts", h.CreateAccount)

	r.POST("/accounts/deposit", h.DepositAccount)

	r.POST("/accounts/withdraw", h.WithdrawAccount)

	r.POST("/accounts/transfer", h.TransferAccount)

	return &HttpService{
		Engine: r,
	}
}
