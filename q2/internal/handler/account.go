package handler

import (
	"github.com/Yougigun/meepshop_q2/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccountHandler struct {
	logger     *zap.Logger
	repository *repository.Repository
	queue      chan any
}

func NewAccountHandler(logger *zap.Logger, repository *repository.Repository) *AccountHandler {
	return &AccountHandler{
		logger:     logger,
		repository: repository,
		queue:      make(chan any, 10000),
	}
}

func (h *AccountHandler) CreateAccount(ctx *gin.Context) {
	// generate account id
}

func (h *AccountHandler) DepositAccount(ctx *gin.Context) {
}

func (h *AccountHandler) WithdrawAccount(ctx *gin.Context) {
}

func (h *AccountHandler) TransferAccount(ctx *gin.Context) {

}

func (h *AccountHandler) GetAccount(ctx *gin.Context) {

}
