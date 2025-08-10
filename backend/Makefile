# Go Backend Makefile

APP_NAME=clashcoder
MAIN_FILE=main.go

.PHONY: run build clean fmt lint test

## Run the backend server
run:
	go run $(MAIN_FILE)

## Build the binary
build:
	go build -o $(APP_NAME) $(MAIN_FILE)

## Run tests
test:
	go test ./...

## Clean up binaries and temp files
clean:
	rm -f $(APP_NAME)

