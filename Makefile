# Go Micro Scaffold - 项目 Makefile
# 提供统一的开发、构建、测试和部署命令

.DEFAULT_GOAL := help
.PHONY: help

# ==================== 通用配置 ====================
PROJECT_NAME := go-micro-scaffold
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := 1.24.1

# 颜色输出
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

# ==================== 帮助信息 ====================
help: ## 📖 显示帮助信息
	@echo "$(COLOR_BOLD)$(PROJECT_NAME) - Makefile 命令列表$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)使用方法:$(COLOR_RESET) make [target]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(COLOR_YELLOW)💡 提示:$(COLOR_RESET) 查看所有可用命令: make help"

# ==================== 开发环境设置 ====================
.PHONY: setup install deps

setup: ## 🛠️  初始化开发环境
	@echo "$(COLOR_BLUE)初始化开发环境...$(COLOR_RESET)"
	@if [ ! -f "go.work" ]; then \
		go work init; \
		go work use ./services; \
		go work use ./common; \
	fi
	@$(MAKE) deps
	@echo "$(COLOR_GREEN)✅ 开发环境初始化完成$(COLOR_RESET)"

install: setup ## 📦 安装项目依赖（同 setup）

deps: ## 📥 下载并整理依赖
	@echo "$(COLOR_BLUE)下载依赖...$(COLOR_RESET)"
	@cd services && go mod download && go mod tidy
	@cd common && go mod download && go mod tidy
	@echo "$(COLOR_GREEN)✅ 依赖下载完成$(COLOR_RESET)"

# ==================== 代码生成 ====================
.PHONY: generate ent-generate

generate: ent-generate ## 🔄 生成所有代码

ent-generate: ## 🔄 生成 Ent 代码
	@echo "$(COLOR_BLUE)生成 Ent 代码...$(COLOR_RESET)"
	@cd services/internal/infrastructure/persistence/ent && \
		go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
	@echo "$(COLOR_GREEN)✅ Ent 代码生成完成$(COLOR_RESET)"

# ==================== 构建 ====================
.PHONY: build build-server build-cli clean

build: build-server build-cli ## 🔨 构建所有二进制文件

build-server: ## 🔨 构建服务端
	@echo "$(COLOR_BLUE)构建服务端...$(COLOR_RESET)"
	@cd services && go build -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" \
		-o bin/server cmd/server/main.go
	@echo "$(COLOR_GREEN)✅ 服务端构建完成: services/bin/server$(COLOR_RESET)"

build-cli: ## 🔨 构建 CLI 工具
	@echo "$(COLOR_BLUE)构建 CLI 工具...$(COLOR_RESET)"
	@cd services && go build -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" \
		-o bin/cli cmd/cli/main.go
	@echo "$(COLOR_GREEN)✅ CLI 工具构建完成: services/bin/cli$(COLOR_RESET)"

clean: ## 🧹 清理构建文件
	@echo "$(COLOR_BLUE)清理构建文件...$(COLOR_RESET)"
	@rm -rf services/bin/
	@rm -rf services/logs/*
	@echo "$(COLOR_GREEN)✅ 清理完成$(COLOR_RESET)"

# ==================== 运行 ====================
.PHONY: run run-server dev

run: run-server ## 🚀 运行服务（同 run-server）

run-server: ## 🚀 运行服务端
	@echo "$(COLOR_BLUE)启动服务端...$(COLOR_RESET)"
	@cd services && go run cmd/server/main.go

dev: run-server ## 💻 开发模式运行（同 run-server）

# ==================== 数据库迁移 (Atlas) ====================
.PHONY: migrate-install migrate-create migrate-apply migrate-status migrate-down migrate-validate migrate-docker

migrate-install: ## 📦 安装 Atlas CLI
	@echo "$(COLOR_BLUE)安装 Atlas...$(COLOR_RESET)"
	@if command -v brew >/dev/null 2>&1; then \
		brew install ariga/tap/atlas; \
	elif command -v go >/dev/null 2>&1; then \
		go install ariga.io/atlas/cmd/atlas@latest; \
	else \
		echo "$(COLOR_YELLOW)⚠️  请手动安装 Atlas: curl -sSf https://atlasgo.sh | sh$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✅ Atlas 安装完成$(COLOR_RESET)"

migrate-create: ## 📝 创建新的迁移文件
	@echo "$(COLOR_BLUE)创建迁移文件...$(COLOR_RESET)"
	@read -p "输入迁移名称 (例如: add_user_email): " name; \
	cd services && GOWORK=off atlas migrate diff $$name --env dev
	@echo "$(COLOR_GREEN)✅ 迁移文件已创建$(COLOR_RESET)"

migrate-apply: ## ✅ 应用迁移到数据库
	@echo "$(COLOR_BLUE)应用迁移...$(COLOR_RESET)"
	@cd services && GOWORK=off atlas migrate apply --env dev
	@echo "$(COLOR_GREEN)✅ 迁移已应用$(COLOR_RESET)"

migrate-status: ## 📊 查看迁移状态
	@cd services && GOWORK=off atlas migrate status --env dev

migrate-down: ## ⬇️  回滚最后一次迁移
	@echo "$(COLOR_YELLOW)⚠️  即将回滚最后一次迁移$(COLOR_RESET)"
	@read -p "确认继续? (y/n): " confirm; \
	if [ "$$confirm" = "y" ]; then \
		cd services && GOWORK=off atlas migrate down --env dev; \
		echo "$(COLOR_GREEN)✅ 迁移已回滚$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)已取消$(COLOR_RESET)"; \
	fi

migrate-validate: ## 🔍 验证迁移文件
	@cd services && GOWORK=off atlas migrate validate --env dev

migrate-docker: ## 🐳 在 Docker 环境应用迁移
	@echo "$(COLOR_BLUE)在 Docker 环境应用迁移...$(COLOR_RESET)"
	@cd services && GOWORK=off atlas migrate apply --env docker
	@echo "$(COLOR_GREEN)✅ Docker 迁移已应用$(COLOR_RESET)"

migrate-quick: ## 🚀 快速启动迁移（交互式）
	@cd services && ./atlas-quickstart.sh

# ==================== 测试 ====================
.PHONY: test test-unit test-integration test-coverage

test: test-unit ## 🧪 运行所有测试

test-unit: ## 🧪 运行单元测试
	@echo "$(COLOR_BLUE)运行单元测试...$(COLOR_RESET)"
	@go test -v -race -short ./services/... ./common/...

test-integration: ## 🧪 运行集成测试
	@echo "$(COLOR_BLUE)运行集成测试...$(COLOR_RESET)"
	@go test -v -race ./services/... ./common/...

test-coverage: ## 📊 生成测试覆盖率报告
	@echo "$(COLOR_BLUE)生成测试覆盖率报告...$(COLOR_RESET)"
	@go test -coverprofile=coverage.out ./services/... ./common/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✅ 覆盖率报告: coverage.html$(COLOR_RESET)"

# ==================== 代码质量 ====================
.PHONY: lint fmt vet

lint: ## 🔍 代码检查 (golangci-lint)
	@echo "$(COLOR_BLUE)运行代码检查...$(COLOR_RESET)"
	@golangci-lint run ./services/... ./common/...

fmt: ## 🎨 格式化代码
	@echo "$(COLOR_BLUE)格式化代码...$(COLOR_RESET)"
	@go fmt ./services/... ./common/...
	@echo "$(COLOR_GREEN)✅ 代码格式化完成$(COLOR_RESET)"

vet: ## 🔍 Go vet 检查
	@echo "$(COLOR_BLUE)运行 go vet...$(COLOR_RESET)"
	@go vet ./services/... ./common/...

# ==================== Docker ====================
.PHONY: docker-build docker-up docker-down docker-restart docker-logs docker-clean

docker-build: ## 🐳 构建 Docker 镜像
	@echo "$(COLOR_BLUE)构建 Docker 镜像...$(COLOR_RESET)"
	@docker-compose build
	@echo "$(COLOR_GREEN)✅ Docker 镜像构建完成$(COLOR_RESET)"

docker-up: ## 🐳 启动 Docker 服务
	@echo "$(COLOR_BLUE)启动 Docker 服务...$(COLOR_RESET)"
	@docker-compose up -d
	@echo "$(COLOR_GREEN)✅ Docker 服务已启动$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)💡 查看日志: make docker-logs$(COLOR_RESET)"

docker-down: ## 🐳 停止 Docker 服务
	@echo "$(COLOR_BLUE)停止 Docker 服务...$(COLOR_RESET)"
	@docker-compose down
	@echo "$(COLOR_GREEN)✅ Docker 服务已停止$(COLOR_RESET)"

docker-restart: docker-down docker-up ## 🐳 重启 Docker 服务

docker-logs: ## 📋 查看 Docker 日志
	@docker-compose logs -f

docker-clean: ## 🧹 清理 Docker 资源
	@echo "$(COLOR_BLUE)清理 Docker 资源...$(COLOR_RESET)"
	@docker-compose down -v
	@docker system prune -f
	@echo "$(COLOR_GREEN)✅ Docker 资源清理完成$(COLOR_RESET)"

# ==================== 开发工具 ====================
.PHONY: swagger-install swagger-generate

swagger-install: ## 📦 安装 Swagger 工具
	@echo "$(COLOR_BLUE)安装 Swag...$(COLOR_RESET)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(COLOR_GREEN)✅ Swag 安装完成$(COLOR_RESET)"

swagger-generate: ## 📝 生成 Swagger 文档
	@echo "$(COLOR_BLUE)生成 Swagger 文档...$(COLOR_RESET)"
	@cd services && swag init -g cmd/server/main.go -o docs
	@echo "$(COLOR_GREEN)✅ Swagger 文档生成完成$(COLOR_RESET)"

# ==================== 数据库管理 ====================
.PHONY: db-reset db-seed

db-reset: ## 🗄️  重置数据库（危险操作）
	@echo "$(COLOR_YELLOW)⚠️  即将删除并重建数据库$(COLOR_RESET)"
	@read -p "确认继续? (y/n): " confirm; \
	if [ "$$confirm" = "y" ]; then \
		mysql -h localhost -u root -p -e "DROP DATABASE IF EXISTS \`go-micro-scaffold\`; CREATE DATABASE \`go-micro-scaffold\`;"; \
		$(MAKE) migrate-apply; \
		echo "$(COLOR_GREEN)✅ 数据库已重置$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)已取消$(COLOR_RESET)"; \
	fi

db-seed: ## 🌱 填充测试数据
	@echo "$(COLOR_BLUE)填充测试数据...$(COLOR_RESET)"
	@cd services && go run cmd/cli/main.go seed
	@echo "$(COLOR_GREEN)✅ 测试数据填充完成$(COLOR_RESET)"

# ==================== 依赖图 ====================
.PHONY: deps-graph

deps-graph: ## 📊 生成依赖关系图
	@echo "$(COLOR_BLUE)生成依赖关系图...$(COLOR_RESET)"
	@cd services && go run cmd/server/main.go -graph=true -graph-output=../assets/dependency-graph.dot
	@dot -Tpng assets/dependency-graph.dot -o assets/dependency-graph.png
	@echo "$(COLOR_GREEN)✅ 依赖图生成完成: assets/dependency-graph.png$(COLOR_RESET)"

# ==================== 一键启动 ====================
.PHONY: quickstart

quickstart: setup migrate-apply run-server ## 🚀 一键启动（安装依赖 → 迁移 → 运行）

# ==================== 版本信息 ====================
.PHONY: version

version: ## 📌 显示版本信息
	@echo "$(COLOR_BOLD)$(PROJECT_NAME)$(COLOR_RESET)"
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
