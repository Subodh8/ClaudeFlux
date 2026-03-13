BINARY_NAME := claudeflux
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS     := -ldflags "-X main.version=$(VERSION) -s -w"
BUILD_DIR   := ./bin
GOFLAGS     := CGO_ENABLED=1
 
## Build the claudeflux binary
build:
	@echo "→ Building $(BINARY_NAME) $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	$(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/claudeflux
	@echo "✓ Binary: $(BUILD_DIR)/$(BINARY_NAME)"
 
## Build for all platforms
build-all:
	GOOS=linux   GOARCH=amd64 $(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64   ./cmd/claudeflux
	GOOS=linux   GOARCH=arm64 $(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64   ./cmd/claudeflux
	GOOS=darwin  GOARCH=amd64 $(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64  ./cmd/claudeflux
	GOOS=darwin  GOARCH=arm64 $(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64  ./cmd/claudeflux
	GOOS=windows GOARCH=amd64 $(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/claudeflux
 
## Run tests
test:
	go test ./... -race -count=1 -timeout 120s
 
## Run tests with coverage
test-coverage:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
 
## Run linter
lint:
	golangci-lint run ./...
 
## Build dashboard for production
dashboard-build:
	cd dashboard && npm ci && npm run build
 
## Run claudeflux in dev mode (hot reload)
dev:
	air
 
## Install to $$GOPATH/bin
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
 
## Run the example research workflow
example:
	./$(BUILD_DIR)/$(BINARY_NAME) run examples/research-write-review/workflow.yaml --dashboard
 
## Start full stack with Docker
docker:
	docker compose up --build
 
## Clean build artifacts
clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html
