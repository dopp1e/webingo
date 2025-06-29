#include .envrc

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: gen-docs
gen-docs:
	@echo "Generating API documentation..."
	swag init -g ./api/main.go -d cmd,internal && swag fmt