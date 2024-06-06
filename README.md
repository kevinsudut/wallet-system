# Wallet System

## Requirements
To run this project you need to have the following installed:
1. [Go](https://golang.org/doc/install) version 1.21
2. [mock](https://github.com/uber-go/mock)
    Install the latest version with:
    ```
    go install go.uber.org/mock/mockgen@latest
    ```
3. [Docker](https://docs.docker.com/get-docker/) version 20   
4. [Docker Compose](https://docs.docker.com/compose/install/) version 1.29

## Initiate The Project
To start working, execute

```
make
```
The backend system will running on http://localhost:8000

## List Available API
1. Register new user (http://localhost:8000/create_user)
2. Read balance (http://localhost:8000/balance_read)
3. Balance top-up (http://localhost:8000/balance_topup)
4. Money transfer between wallets (http://localhost:8000/transfer)
5. List top N transactions by value per user  (http://localhost:8000/top_users)
6. List overall top N transacting users by value  (http://localhost:8000/top_transaction_per_user)