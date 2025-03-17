.PHONY: build run clean swagger

APP_NAME=accounts

build:
	go build -o $(APP_NAME) main.go

run:
	go run main.go

clean:
	rm -f $(APP_NAME)
	rm -f accounts.db

swagger:
	swag init

test:
	go test ./...

.DEFAULT_GOAL := build