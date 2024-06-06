.PHONY: init build run

all: init build test run

init:
	go mod tidy
	go mod vendor
	make generate_mocks

build:
	go build -o build/wallet-system.exe cmd/main.go 

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

run:
	docker compose up --build -d

stop:
	docker compose down --volumes

generate_mocks:
	mockgen -source=app/domain/auth/interfaces.go -destination=app/domain/auth/mock.go -package=domainauth
	mockgen -source=app/domain/balance/interfaces.go -destination=app/domain/balance/mock.go -package=domainbalance
	mockgen -source=app/handler/template/template.go -destination=app/handler/template/mock.go -package=handlertemplate
	mockgen -source=app/usecase/auth/interfaces.go -destination=app/usecase/auth/mock.go -package=usecaseauth
	mockgen -source=app/usecase/balance/interfaces.go -destination=app/usecase/balance/mock.go -package=usecasebalance
	mockgen -source=app/usecase/transaction/interfaces.go -destination=app/usecase/transaction/mock.go -package=usecasetransaction
	mockgen -source=pkg/helper/singleflight/interfaces.go -destination=pkg/helper/singleflight/mock.go -package=singleflight
	mockgen -source=pkg/lib/database/interfaces.go -destination=pkg/lib/database/mock.go -package=database
	mockgen -source=pkg/lib/lru-cache/interfaces.go -destination=pkg/lib/lru-cache/mock.go -package=lrucache
	mockgen -source=pkg/lib/token/interfaces.go -destination=pkg/lib/token/mock.go -package=token
