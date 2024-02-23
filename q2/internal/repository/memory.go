package repository

import (
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

type transactionLog struct {
	From   accountID
	To     accountID
	Amount int64
	When   time.Time
}

type transactions struct {
	transactions []transactionLog
	rw           sync.RWMutex
}

type Repository struct {
	Accounts     map[accountID]*account
	Transactions transactions
}

func NewRepository() *Repository {
	return &Repository{
		Accounts: make(map[accountID]*account),
	}
}

func (r *Repository) CreateAccount() (accountID, error) {
	// use uuid to generate account id
	id := atomic.AddInt64(&idCounter, 1)
	r.Accounts[accountID(id)] = &account{
		ID:      accountID(id),
		Balance: 0,
	}
	return accountID(id), nil
}

func (r *Repository) GetAccount(id accountID) (*account, error) {
	r.Accounts[id].rw.RLock()
	defer r.Accounts[id].rw.RUnlock()
	readAccount := &account{
		ID:      r.Accounts[id].ID,
		Balance: r.Accounts[id].Balance,
	}
	return readAccount, nil
}

func (r *Repository) DepositAccount(account *account, amount int) error {
	account.rw.Lock()
	defer account.rw.Unlock()
	account.Balance += amount
	return nil
}

func (r *Repository) WithdrawAccount(id accountID, amount int) (*account, error) {
	r.Accounts[id].rw.Lock()
	defer r.Accounts[id].rw.Unlock()
	r.Accounts[id].Balance -= amount
	if r.Accounts[id].Balance < 0 {
		r.Accounts[id].Balance += amount
		return nil, errors.New("insufficient funds")
	}
	return r.Accounts[id], nil
}

func (r *Repository) TransferAccount(from accountID, to accountID, amount int) (*account, error) {
	// Ensure consistent locking order
	ids := []int{int(from), int(to)}
	sort.Ints(ids)
	first, second := accountID(ids[0]), accountID(ids[1])

	r.Accounts[first].rw.Lock()
	defer r.Accounts[first].rw.Unlock()

	r.Accounts[second].rw.Lock()
	defer r.Accounts[second].rw.Unlock()

	if from == to {
		// Handle the case where from and to are the same, which could be a no-op or an error
		return nil, errors.New("cannot transfer to the same account")
	}

	// Perform the transfer
	r.Accounts[from].Balance -= amount
	if r.Accounts[from].Balance < 0 {
		r.Accounts[from].Balance += amount // Roll back the change
		return nil, errors.New("insufficient funds")
	}
	r.Accounts[to].Balance += amount

	return r.Accounts[from], nil
}
