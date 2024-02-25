package handler

import (
	"context"
	"time"

	"github.com/Yougigun/meepshop_q2/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AccountHandler struct {
	logger              *zap.Logger
	repository          *repository.Repository
	transactionLogQueue chan *TransactionLog
}

func NewAccountHandler(ctx context.Context, logger *zap.Logger, repo *repository.Repository) *AccountHandler {
	queue := make(chan *TransactionLog, 10000)
	batchLogs := make([]struct {
		From   int64
		To     int64
		Amount int
		When   time.Time
	}, 0, 300)
	// deal with transaction log, this may lose some logs if the server is down. todo: use kafka or other message queue
	go func() {
		for {
			select {
			case tl := <-queue:
				batchLogs = append(batchLogs, struct {
					From   int64
					To     int64
					Amount int
					When   time.Time
				}{
					From:   tl.From,
					To:     tl.To,
					Amount: tl.Amount,
					When:   tl.When,
				},
				)
				// if more than 3000 logs, then add to repository
				if len(batchLogs) > 300 {
					repo.AddTransaction(ctx, repository.BatchTransaction(batchLogs))
					batchLogs = make([]struct {
						From   int64
						To     int64
						Amount int
						When   time.Time
					}, 0, 300)
				}

			case <-time.After(5 * time.Second):
				repo.AddTransaction(ctx, repository.BatchTransaction(batchLogs))
			}
		}
	}()

	return &AccountHandler{
		logger:              logger,
		repository:          repo,
		transactionLogQueue: queue,
	}
}

func (h *AccountHandler) CreateAccount(gCtx *gin.Context) {
	ctx := gCtx.Request.Context()
	if account, err := h.repository.CreateAccount(ctx); err != nil {
		gCtx.JSON(500, err.Error())
	} else {
		h.logger.Info("create account", zap.Any("account", account))
		gCtx.JSON(200, account)
	}
}

type DepositAccountRequest struct {
	AccountID int64 `json:"account_id"`
	Amount    int   `json:"amount"`
}

func (h *AccountHandler) DepositAccount(ctx *gin.Context) {
	reqBody := &DepositAccountRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	// deposit account
	if err := h.repository.DepositAccount(ctx, reqBody.AccountID, reqBody.Amount); err != nil {
		ctx.JSON(500, err.Error())
	} else {
		h.logger.Info("deposit account", zap.Any("account_id", reqBody.AccountID), zap.Any("amount", reqBody.Amount))
		ctx.JSON(200, "success")
	}
}

type WithdrawAccountRequest struct {
	AccountID int64 `json:"account_id"`
	Amount    int   `json:"amount"`
}

func (h *AccountHandler) WithdrawAccount(ctx *gin.Context) {
	reqBody := &WithdrawAccountRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	// withdraw account
	if err := h.repository.WithdrawAccount(ctx, reqBody.AccountID, reqBody.Amount); err != nil {
		ctx.JSON(500, err.Error())
	} else {
		h.logger.Info("withdraw account", zap.Any("account_id", reqBody.AccountID), zap.Any("amount", reqBody.Amount))
		ctx.JSON(200, "success")
	}

}

type TransferAccountRequest struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int   `json:"amount"`
}

type TransactionLog struct {
	From   int64     `json:"from"`
	To     int64     `json:"to"`
	Amount int       `json:"amount"`
	When   time.Time `json:"when"`
}

func (h *AccountHandler) TransferAccount(ctx *gin.Context) {
	reqBody := &TransferAccountRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	// transfer account
	if err := h.repository.TransferAccount(ctx, reqBody.FromAccountID, reqBody.ToAccountID, reqBody.Amount); err != nil {
		ctx.JSON(500, err.Error())
	} else {
		ctx.JSON(200, "success")
	}
	// log transaction
	tl := &TransactionLog{
		From:   reqBody.FromAccountID,
		To:     reqBody.ToAccountID,
		Amount: reqBody.Amount,
		When:   time.Now(),
	}
	h.transactionLogQueue <- tl
	h.logger.Info("transaction log", zap.Any("log", tl))
}

type GetAccountRequest struct {
	ID      int64 `json:"id"`
	Balance int   `json:"balance"`
}

func (h *AccountHandler) GetAccount(ctx *gin.Context) {
	reqBody := &GetAccountRequest{}
	if err := ctx.ShouldBindJSON(reqBody); err != nil {
		ctx.JSON(400, err.Error())
		return
	}
	// get account
	if account, err := h.repository.GetAccount(ctx, reqBody.ID); err != nil {
		ctx.JSON(500, err.Error())
	} else {
		h.logger.Info("get account", zap.Any("account_id", reqBody.ID))
		ctx.JSON(200, account)
	}
}

func (h *AccountHandler) GetTransactionLog(ctx *gin.Context) {
	// get transaction log
	tl := h.repository.GetTransactions(ctx)
	// h.logger.Info("get transaction log", zap.Any("log", tl))
	ctx.JSON(200, tl)
}
