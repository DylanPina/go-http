SERVER_BIN=server
CLIENT_BIN=client
SERVER_SRC=cmd/server/server.go
CLIENT_SRC=cmd/client/client.go

.PHONY: run clean all

build:
	go build -o $(SERVER_BIN) $(SERVER_SRC)
	go build -o $(CLIENT_BIN) $(CLIENT_SRC)

run: build
	@echo "Running server..."
	./$(SERVER_BIN) &

	@echo "Running client..."
	./$(CLIENT_BIN)

run-server:
	./$(SERVER_BIN)

run-client:
	./$(CLIENT_BIN)

clean:
	rm -f $(SERVER_BIN) $(CLIENT_BIN)

all: build run clean
