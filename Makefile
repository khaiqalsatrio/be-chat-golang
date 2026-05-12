.PHONY: run build test clean

APP_NAME=chat-golang
MAIN_PATH=./src/cmd/server/main.go

run:
	go run $(MAIN_PATH)

build:
	go build -o $(APP_NAME).exe $(MAIN_PATH)

test:
	go test ./...

clean:
	if exist $(APP_NAME).exe del $(APP_NAME).exe
