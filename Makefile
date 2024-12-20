.PHONY: build run

build:
	@echo "Building..."
	@go build -o bin/hard cmd/hard/main.go
run: build
	@./bin/hard
