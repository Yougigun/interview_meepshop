package repository

import (
	"context"
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type RepositoryI interface {
	CreateAccount(account *account) error
	GetAccount(id accountID) (*account, error)
	DepositAccount(account *account, amount int) error
	WithdrawAccount(id accountID) (*account, error)
	TransferAccount(id accountID) (*account, error)
}

type accountID int64

var idCounter int64

type account struct {
	ID      accountID
	Balance int
	rw      sync.RWMutex
}

type TransactionLog struct {
	From   accountID
	To     accountID
	Amount int64
	When   time.Time
}

type transactions struct {
	transactions []TransactionLog
	rw           sync.RWMutex
}

type Repository struct {
	Accounts     map[accountID]*account
	Transactions transactions
}

func NewRepository() *Repository {
	return &Repository{
		Accounts: make(map[accountID]*account),
		Transactions: transactions{
			transactions: make([]TransactionLog, 0),
		},
	}
}

func (r *Repository) CreateAccount(ctx context.Context) (accountID, error) {
	// use uuid to generate account id
	id := atomic.AddInt64(&idCounter, 1)
	r.Accounts[accountID(id)] = &account{
		ID:      accountID(id),
		Balance: 0,
	}
	return accountID(id), nil
}

func (r *Repository) GetAccount(ctx context.Context, id int64) (*account, error) {
	// check if account exists
	if r.Accounts[accountID(id)] == nil {
		return nil, errors.New("account not found")
	}
	r.Accounts[accountID(id)].rw.RLock()
	defer r.Accounts[accountID(id)].rw.RUnlock()
	readAccount := &account{
		ID:      r.Accounts[accountID(id)].ID,
		Balance: r.Accounts[accountID(id)].Balance,
	}
	return readAccount, nil
}

func (r *Repository) DepositAccount(ctx context.Context, aid int64, amount int) error {
	if account := r.Accounts[accountID(aid)]; account == nil {
		return errors.New("account not found")
	} else {
		account.rw.Lock()
		defer account.rw.Unlock()
		account.Balance += amount
		return nil
	}
}

func (r *Repository) WithdrawAccount(ctx context.Context, id int64, amount int) error {

	// check if account exists
	if r.Accounts[accountID(id)] == nil {
		return errors.New("account not found")
	}
	r.Accounts[accountID(id)].rw.Lock()
	defer r.Accounts[accountID(id)].rw.Unlock()
	r.Accounts[accountID(id)].Balance -= amount
	if r.Accounts[accountID(id)].Balance < 0 {
		r.Accounts[accountID(id)].Balance += amount
		return errors.New("insufficient funds")
	}
	return nil
}

func (r *Repository) TransferAccount(ctx context.Context, from int64, to int64, amount int) error {
	fromID := accountID(from)
	toID := accountID(to)
	// check if account exists
	if r.Accounts[fromID] == nil || r.Accounts[toID] == nil {
		return errors.New("account not found")
	}

	if fromID == toID {
		// Handle the case where from and to are the same, which could be a no-op or an error
		return errors.New("cannot transfer to the same account")
	}

	// Ensure consistent locking order
	ids := []int{int(from), int(to)}
	sort.Ints(ids)
	first, second := accountID(ids[0]), accountID(ids[1])

	r.Accounts[first].rw.Lock()
	defer r.Accounts[first].rw.Unlock()

	r.Accounts[second].rw.Lock()
	defer r.Accounts[second].rw.Unlock()

	// Perform the transfer
	r.Accounts[fromID].Balance -= amount
	if r.Accounts[fromID].Balance < 0 {
		r.Accounts[fromID].Balance += amount // Roll back the change
		return errors.New("insufficient funds")
	}
	r.Accounts[toID].Balance += amount

	return nil
}

// GetTransactions returns a copy of the transaction log
func (r *Repository) GetTransactions(ctx context.Context) []TransactionLog {
	r.Transactions.rw.RLock()
	defer r.Transactions.rw.RUnlock()
	return append([]TransactionLog(nil), r.Transactions.transactions...)
}

type BatchTransaction []struct {
	From   int64
	To     int64
	Amount int
	When   time.Time
}

// AddTransaction adds a new transaction to the log
func (r *Repository) AddTransaction(ctx context.Context, batch BatchTransaction) {
	r.Transactions.rw.Lock()
	defer r.Transactions.rw.Unlock()
	for _, transaction := range batch {
		r.Transactions.transactions = append(r.Transactions.transactions, TransactionLog{
			From:   accountID(transaction.From),
			To:     accountID(transaction.To),
			Amount: int64(transaction.Amount),
			When:   transaction.When,
		})
	}
}
