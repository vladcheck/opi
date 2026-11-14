build:
	@go build -o bin/opi ./cmd

run: build
	@./bin/opi
