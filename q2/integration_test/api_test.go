package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yougigun/meepshop_q2/internal/repository"
	"github.com/Yougigun/meepshop_q2/internal/service"
	"go.uber.org/zap"
)

func TestCreateAccountAPI(t *testing.T) {
	// Setup
	logger := zap.NewNop()             
	repo := repository.NewRepository() 

	router := service.Build(context.Background(), logger, repo)

	reqBody := bytes.NewBufferString(`{}`)
	req, err := http.NewRequest("POST", "/accounts", reqBody) // Adjust the HTTP method and endpoint as necessary
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method directly and pass in our Request and ResponseRecorder.
	router.Handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"AccountID":1}` // Adjust expected response as necessary
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDepositAccountAPI(t *testing.T) {
	// Setup
	logger := zap.NewNop()             
	repo := repository.NewRepository() 

	router := service.Build(context.Background(), logger, repo)
	
	reqBody := bytes.NewBufferString(`{}`)
	req, err := http.NewRequest("POST", "/accounts", reqBody) 
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"AccountID":1}` // Adjust expected response as necessary
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	depositBody := bytes.NewBufferString(`{"account_id":1,"amount":100}`)
	req, err = http.NewRequest("POST", "/accounts/deposit", depositBody)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.Handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected = `"success"`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}


func TestWithdrawAccountAPI(t *testing.T) {
    // Setup
    logger := zap.NewNop()             
    repo := repository.NewRepository() 
    router := service.Build(context.Background(), logger, repo)

    // Create an account
    createAccBody := bytes.NewBufferString(`{}`)
    createReq, err := http.NewRequest("POST", "/accounts", createAccBody)
    if err != nil {
        t.Fatal(err)
    }

    createRR := httptest.NewRecorder()
    router.Handler.ServeHTTP(createRR, createReq)

    if status := createRR.Code; status != http.StatusOK {
        t.Errorf("Create account handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Assuming the account creation returns JSON like {"AccountID": 1}
    // Extract the account ID from the response for further operations
    // This part needs adjustment based on your actual response structure
    var account struct {
        AccountID int64 `json:"AccountID"`
    }
    err = json.Unmarshal(createRR.Body.Bytes(), &account)
    if err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }

    // Deposit to the account before withdrawal to ensure sufficient balance
    depositBody := bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"amount":100}`, account.AccountID))
    depositReq, err := http.NewRequest("POST", "/accounts/deposit", depositBody)
    if err != nil {
        t.Fatal(err)
    }

    depositRR := httptest.NewRecorder()
    router.Handler.ServeHTTP(depositRR, depositReq)

    if status := depositRR.Code; status != http.StatusOK {
        t.Errorf("Deposit handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Withdraw from the account
    withdrawBody := bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"amount":50}`, account.AccountID))
    withdrawReq, err := http.NewRequest("POST", "/accounts/withdraw", withdrawBody)
    if err != nil {
        t.Fatal(err)
    }

    withdrawRR := httptest.NewRecorder()
    router.Handler.ServeHTTP(withdrawRR, withdrawReq)

    if status := withdrawRR.Code; status != http.StatusOK {
        t.Errorf("Withdraw handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Here you should adjust the expected response as per your API's specification
    expectedWithdrawResponse := `"success"` // Adjust this based on your actual success response
    if withdrawRR.Body.String() != expectedWithdrawResponse {
        t.Errorf("Withdraw handler returned unexpected body: got %v want %v", withdrawRR.Body.String(), expectedWithdrawResponse)
    }
}



func TestTransferAccountAPI(t *testing.T) {
	logger := zap.NewNop()
	repo := repository.NewRepository() // Initialize your repository here
	router := service.Build(context.Background(), logger, repo)

	// Helper function to create an account and return its ID
	createAccount := func() int64 {
		reqBody := bytes.NewBufferString(`{}`)
		req, err := http.NewRequest("POST", "/accounts", reqBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.Handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code during account creation: got %v want %v", status, http.StatusOK)
		}

		var account struct {
			AccountID int64 `json:"AccountID"`
		}
		if err := json.Unmarshal(rr.Body.Bytes(), &account); err != nil {
			t.Fatalf("Failed to unmarshal response during account creation: %v", err)
		}

		return account.AccountID
	}

	// Create two accounts
	fromAccountID := createAccount()
	toAccountID := createAccount()

	// Deposit into the first account to ensure sufficient balance
	depositBody := bytes.NewBufferString(fmt.Sprintf(`{"account_id":%d,"amount":100}`, fromAccountID))
	depositReq, err := http.NewRequest("POST", "/accounts/deposit", depositBody)
	if err != nil {
		t.Fatal(err)
	}

	depositRR := httptest.NewRecorder()
	router.Handler.ServeHTTP(depositRR, depositReq)

	if status := depositRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code during deposit: got %v want %v", status, http.StatusOK)
	}

	// Perform the transfer
	transferBody := bytes.NewBufferString(fmt.Sprintf(`{"from_account_id":%d,"to_account_id":%d,"amount":50}`, fromAccountID, toAccountID))
	transferReq, err := http.NewRequest("POST", "/accounts/transfer", transferBody)
	if err != nil {
		t.Fatal(err)
	}

	transferRR := httptest.NewRecorder()
	router.Handler.ServeHTTP(transferRR, transferReq)

	if status := transferRR.Code; status != http.StatusOK {
		t.Errorf("Transfer handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedTransferResponse := `"success"` // Adjust based on your actual success response
	if transferRR.Body.String() != expectedTransferResponse {
		t.Errorf("Transfer handler returned unexpected body: got %v want %v", transferRR.Body.String(), expectedTransferResponse)
	}
}
