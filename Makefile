build:
	@go build -o bin/notify cmd/main.go

run: build
	@./bin/notify