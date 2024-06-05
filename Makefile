.PHONY: init build run

all: init build run

init:
	go mod tidy
	go mod vendor

build:
	go build -o build/wallet-system.exe cmd/main.go 

run:
	docker compose up --build -d

stop:
	docker compose down --volumes