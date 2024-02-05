run-dev:
	@APP_ENV=development APP_LOCAL_ENV_FILE_PATH=.env.example go run ./...

mocks:
	@echo "Generating mocks..."
	@mockery

lint:
	@echo "Run golangci-lint..."
	@golangci-lint run -v

tests:
	@echo "Run tests..."
	@go test ./... -race -v -coverprofile cover.out

cover:
	@echo "Run coverage..."
	@go tool cover -html=cover.out

.PHONY: run-dev mocks lint tests cover
