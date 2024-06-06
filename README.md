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
The wallet system will running on http://localhost:8000

## List Available API
1. Register new user (http://localhost:8000/create_user)
```
curl --location --request POST 'http://localhost:8000/create_user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "exampleusername"
}'
```
2. Read balance (http://localhost:8000/balance_read)
```
curl --location 'http://localhost:8000/balance_read' \
--header 'Authorization: Bearer ••••••'
```
3. Balance top-up (http://localhost:8000/balance_topup)
```
curl --location --request POST 'http://localhost:8000/balance_topup' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ••••••' \
--data-raw '{
    "amount": 1000000
}'
```
4. Money transfer between wallets (http://localhost:8000/transfer)
```
curl --location --request POST 'http://localhost:8000/transfer' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ••••••' \
--data-raw '{
    "to_username": "targetusername",
    "amount": 50000
}'
```
5. List top N transactions by value per user  (http://localhost:8000/top_transaction_per_user)
```
curl --location 'http://localhost:8000/top_transaction_per_user' \
--header 'Authorization: Bearer ••••••'
```
6. List overall top N transacting users by value  (http://localhost:8000/top_users)
```
curl --location 'http://localhost:8000/top_users' \
--header 'Authorization: Bearer ••••••'
```

## Testing
To run test, run the following command:
```
make test
```

## API Testing
To run api test, run the following command:
```
make test_api
```
