# Simple Banking System

## API

### Create Account

- Endpoint: `POST /accounts`
- Description: Creates a new account.
- Response: JSON object with the created account ID.

```json
{
  "AccountID": 1
}
```

### Deposit to Account

Endpoint: POST /accounts/deposit
Request Body: JSON object with account_id (int64) and amount (int).

```json
{
  "account_id": 1,
  "amount": 100
}
```

### Withdraw from Account

Endpoint: POST /accounts/withdraw
Request Body: JSON object with account_id (int64) and amount (int).

```json
{
  "account_id": 1,
  "amount": 100
}
```

### Transfer between Accounts

Endpoint: POST /accounts/transfer
Request Body: JSON object with from_account_id (int64), to_account_id (int64), and amount (int).

```json
    {
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 100
    }
```

## Docker

```bash
docker build -t simple-banking-system .
docker run -p 8080:8080 simple-banking-system
```
