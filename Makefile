.PHONY: build run fmt lint vet test clean install setup-hooks

# Build the binary
build:
	go build -o peachy .

# Run the application
run: build
	./peachy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Run linter with auto-fix
lint-fix:
	golangci-lint run --fix

# Run go vet
vet:
	go vet ./...

# Run all checks (format, vet, lint, build)
check: fmt vet lint build

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f peachy
	go clean

# Install the binary to GOPATH/bin
install:
	go install .

# Setup git hooks
setup-hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/*

# Development: run with live reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Run: go install github.com/air-verse/air@latest"; \
	fi
