# Makefile for NETCONF Search Tool

.PHONY: build clean test install run-example help

# 默认目标
all: build

# 构建程序
build:
	@echo "Building netconf-search..."
	@go build -o netconf-search .

# 安装程序到$GOPATH/bin
install:
	@echo "Installing netconf-search..."
	@go install .

# 运行测试
test:
	@echo "Running tests..."
	@go test ./...

# 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -f netconf-search

# 运行示例
run-example: build
	@echo "Running example searches..."
	@echo "\n=== Searching for 'interface' ==="
	@./netconf-search interface ./testdata
	@echo "\n=== Searching for 'hostname' ==="
	@./netconf-search hostname ./testdata

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 检查代码
lint:
	@echo "Running go vet..."
	@go vet ./...

# 显示帮助
help:
	@echo "Available targets:"
	@echo "  build       - Build the netconf-search binary"
	@echo "  install     - Install netconf-search to \$$GOPATH/bin"
	@echo "  test        - Run all tests"
	@echo "  clean       - Remove build artifacts"
	@echo "  run-example - Build and run example searches"
	@echo "  fmt         - Format Go code"
	@echo "  lint        - Run go vet"
	@echo "  help        - Show this help message"
