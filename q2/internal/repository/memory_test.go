package repository

import (
	"context"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	accID, err := repo.CreateAccount(ctx)
	if err != nil {
		t.Errorf("CreateAccount() error = %v, wantErr %v", err, false)
	}
	if accID == 0 {
		t.Errorf("CreateAccount() got = %v, want %v", accID, "non-zero ID")
	}
}

func TestGetAccount(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	// Pre-create an account to test retrieval
	expectedID, _ := repo.CreateAccount(ctx)

	acc, err := repo.GetAccount(ctx, int64(expectedID))
	if err != nil {
		t.Errorf("GetAccount() error = %v, wantErr %v", err, false)
	}
	if acc.ID != expectedID {
		t.Errorf("GetAccount() got = %v, want %v", acc.ID, expectedID)
	}
}

func TestDepositAccount(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	// Create an account for deposit testing
	accID, _ := repo.CreateAccount(ctx)

	err := repo.DepositAccount(ctx, int64(accID), 100)
	if err != nil {
		t.Errorf("DepositAccount() error = %v, wantErr %v", err, false)
	}

	acc, _ := repo.GetAccount(ctx, int64(accID))
	if acc.Balance != 100 {
		t.Errorf("DepositAccount() got = %v, want %v", acc.Balance, 100)
	}
}

func TestWithdrawAccount(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	// Create an account and deposit an initial amount
	accID, _ := repo.CreateAccount(ctx)
	_ = repo.DepositAccount(ctx, int64(accID), 200)

	// Withdraw a valid amount
	if err := repo.WithdrawAccount(ctx, int64(accID), 100); err != nil {
		t.Errorf("WithdrawAccount() error = %v, wantErr %v", err, false)
	}

	// Assert the balance is as expected
	acc, _ := repo.GetAccount(ctx, int64(accID))
	if acc.Balance != 100 {
		t.Errorf("WithdrawAccount() got = %v, want %v", acc.Balance, 100)
	}

	// Attempt to withdraw more than the balance
	if err := repo.WithdrawAccount(ctx, int64(accID), 200); err == nil {
		t.Errorf("WithdrawAccount() expected error for insufficient funds, got nil")
	}
}

func TestTransferAccount(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	// Create two accounts
	fromAccID, _ := repo.CreateAccount(ctx)
	toAccID, _ := repo.CreateAccount(ctx)

	// Deposit into the first account
	_ = repo.DepositAccount(ctx, int64(fromAccID), 300)

	// Transfer funds
	if err := repo.TransferAccount(ctx, int64(fromAccID), int64(toAccID), 150); err != nil {
		t.Errorf("TransferAccount() error = %v, wantErr %v", err, false)
	}

	// Assert balances are as expected
	fromAcc, _ := repo.GetAccount(ctx, int64(fromAccID))
	toAcc, _ := repo.GetAccount(ctx, int64(toAccID))
	if fromAcc.Balance != 150 || toAcc.Balance != 150 {
		t.Errorf("TransferAccount() gotFrom = %v, want %v; gotTo = %v, want %v", fromAcc.Balance, 150, toAcc.Balance, 150)
	}

	// Test transferring with insufficient funds
	if err := repo.TransferAccount(ctx, int64(fromAccID), int64(toAccID), 300); err == nil {
		t.Errorf("TransferAccount() expected error for insufficient funds, got nil")
	}
}

func TestAddTransaction(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()

	// Create a batch of transactions
	batch := BatchTransaction{
		{From: 1, To: 2, Amount: 100, When: time.Now()},
		{From: 2, To: 1, Amount: 50, When: time.Now()},
	}

	// Add transactions to the log
	repo.AddTransaction(ctx, batch)

	// Retrieve the transactions log
	trans := repo.GetTransactions(ctx)
	if len(trans) != 2 {
		t.Errorf("AddTransaction() got = %v transactions, want %v", len(trans), 2)
	}

	// Verify the first transaction details
	if trans[0].From != 1 || trans[0].To != 2 || trans[0].Amount != 100 {
		t.Errorf("AddTransaction() gotFirstTransaction = %+v, want From=1, To=2, Amount=100", trans[0])
	}
}



