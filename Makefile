# Go文件分类系统 Makefile

# 变量定义
BINARY_NAME=file-classifier
MAIN_PATH=./cmd/server
BUILD_DIR=./bin

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

.PHONY: build
build: ## 构建应用程序
	@echo "构建应用程序..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: run
run: ## 运行应用程序
	@echo "启动文件分类服务器..."
	go run $(MAIN_PATH)/main.go

.PHONY: dev
dev: ## 开发模式运行（自动重启）
	@echo "开发模式启动..."
	@which air > /dev/null || (echo "请先安装 air: go install github.com/cosmtrek/air@latest" && exit 1)
	air

.PHONY: test
test: ## 运行测试
	go test -v ./...

.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -rf uploads/
	@echo "清理完成"

.PHONY: tidy
tidy: ## 整理依赖
	go mod tidy

.PHONY: fmt
fmt: ## 格式化代码
	go fmt ./...

.PHONY: vet
vet: ## 静态检查
	go vet ./...

.PHONY: lint
lint: ## 代码检查
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint" && exit 1)
	golangci-lint run

.PHONY: install
install: build ## 安装到系统
	@echo "安装应用程序..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "安装完成，可以直接使用 $(BINARY_NAME) 命令"

.PHONY: docker-build
docker-build: ## 构建Docker镜像
	docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run: ## 运行Docker容器
	docker run -p 3000:3000 $(BINARY_NAME):latest