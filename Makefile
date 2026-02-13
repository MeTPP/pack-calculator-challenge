GO = go
BINARY_NAME = pack-calculator
MAIN_PATH = ./cmd/main.go

.PHONY: test build run stop clean lint

test:
	$(GO) test -v -race ./...

build:
	$(GO) build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_PATH)

run:
	docker compose up --build

stop:
	docker compose down

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME)
	docker compose down --rmi all --volumes --remove-orphans
