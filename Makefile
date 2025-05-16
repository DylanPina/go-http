SERVER_BIN=go-http-server.exe
SERVER_SRC=cmd/server/main.go
SERVER_TEST_SRC=cmd/server/main_test.go

.PHONY: run clean all test

build:
	go build -o $(SERVER_BIN) $(SERVER_SRC)

run: build
	@echo "Running server..."
	./$(SERVER_BIN) &

test: build
	@echo "Starting server..."
	./$(SERVER_BIN) &

	go test -v $(SERVER_TEST_SRC)

	@echo "Cleaning up..."
	rm -f $(SERVER_BIN)

clean:
	rm -f $(SERVER_BIN)

all: build run clean
