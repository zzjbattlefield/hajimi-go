# Makefile for Hajimi GO - GitHub Secret Scanner

# Binary name
BINARY=hajimi-go

# Main package path
MAIN_PACKAGE=./scanner

# Build the binary
build:
	go build -o ${BINARY} ${MAIN_PACKAGE}

# Install dependencies
deps:
	go mod tidy

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run the application
run:
	go run ${MAIN_PACKAGE}

# Clean build artifacts
clean:
	rm -f ${BINARY}

# Install the binary
install:
	go install ${MAIN_PACKAGE}

# Build for different platforms
build-all: build-linux build-windows build-mac

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}-linux ${MAIN_PACKAGE}

build-windows:
	GOOS=windows GOARCH=amd64 go build -o ${BINARY}-windows.exe ${MAIN_PACKAGE}

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY}-mac ${MAIN_PACKAGE}

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install the binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-mac     - Build for Mac"
	@echo "  help          - Show this help"

.PHONY: build deps test test-coverage run clean install build-all build-linux build-windows build-mac help