SERVER=server
CLIENT=client

.PHONY: run clean all

build:
	go build -o $(SERVER) ./$(SERVER)
	go build -o $(CLIENT) ./$(CLIENT)

run: build
	@echo "Running server..."
	server/$(SERVER) &

	@echo "Running client..."
	client/$(CLIENT)

clean:
	rm -f server/$(SERVER) client/$(CLIENT)

all: build run clean
