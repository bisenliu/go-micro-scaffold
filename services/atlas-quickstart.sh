#!/bin/bash

# Atlas 迁移快速启动脚本
# 用法: ./atlas-quickstart.sh

set -e

echo "🚀 Atlas 数据库迁移快速启动"
echo "================================"

# 检查是否安装了 Atlas
if ! command -v atlas &> /dev/null; then
    echo "❌ Atlas 未安装"
    echo ""
    echo "请选择安装方式:"
    echo "1) macOS (使用 Homebrew)"
    echo "2) Linux/macOS (使用官方安装脚本)"
    echo "3) 使用 Go 安装"
    echo "4) 跳过安装，手动安装"
    read -p "请选择 (1-4): " choice
    
    case $choice in
        1)
            echo "📦 使用 Homebrew 安装 Atlas..."
            brew install ariga/tap/atlas
            ;;
        2)
            echo "📦 使用官方脚本安装 Atlas..."
            curl -sSf https://atlasgo.sh | sh
            ;;
        3)
            echo "📦 使用 Go 安装 Atlas..."
            go install ariga.io/atlas/cmd/atlas@latest
            ;;
        4)
            echo "⏭️  跳过安装"
            echo "请访问 https://atlasgo.io/getting-started#installation 手动安装"
            exit 1
            ;;
        *)
            echo "❌ 无效选择"
            exit 1
            ;;
    esac
    
    echo "✅ Atlas 安装完成"
    atlas version
fi

echo ""
echo "✅ Atlas 已安装: $(atlas version)"
echo ""

# 检查 MySQL 是否运行
echo "🔍 检查 MySQL 连接..."
if nc -z localhost 3306 2>/dev/null; then
    echo "✅ MySQL 正在运行 (localhost:3306)"
else
    echo "⚠️  MySQL 未运行"
    echo ""
    read -p "是否启动 Docker Compose? (y/n): " start_docker
    if [[ $start_docker == "y" ]]; then
        echo "🐳 启动 Docker Compose..."
        cd ..
        docker-compose up -d mysql redis
        cd services
        echo "⏳ 等待 MySQL 启动..."
        sleep 10
    else
        echo "❌ 无法继续，请先启动 MySQL"
        exit 1
    fi
fi

echo ""
echo "📝 选择操作:"
echo "1) 生成初始迁移文件"
echo "2) 应用迁移到数据库"
echo "3) 查看迁移状态"
echo "4) 生成新的迁移文件"
echo "5) 回滚迁移"
echo "6) 验证迁移文件"
read -p "请选择 (1-6): " action

case $action in
    1)
        echo "📝 生成初始迁移文件..."
        GOWORK=off atlas migrate diff initial \
          --env dev \
          --to ent://internal/infrastructure/persistence/ent/schema
        
        echo ""
        echo "✅ 迁移文件已生成"
        echo "📂 位置: internal/infrastructure/persistence/ent/migrations/"
        ls -lah internal/infrastructure/persistence/ent/migrations/
        ;;
    
    2)
        echo "🔄 应用迁移到数据库..."
        GOWORK=off atlas migrate apply --env dev
        
        echo ""
        echo "✅ 迁移已应用"
        ;;
    
    3)
        echo "📊 查看迁移状态..."
        GOWORK=off atlas migrate status --env dev
        ;;
    
    4)
        read -p "输入迁移名称 (例如: add_user_email): " migration_name
        
        echo "📝 生成迁移文件: $migration_name"
        GOWORK=off atlas migrate diff "$migration_name" --env dev
        
        echo ""
        echo "✅ 迁移文件已生成"
        ls -lah internal/infrastructure/persistence/ent/migrations/
        
        echo ""
        read -p "查看生成的 SQL? (y/n): " view_sql
        if [[ $view_sql == "y" ]]; then
            latest_file=$(ls -t internal/infrastructure/persistence/ent/migrations/*.sql | head -1)
            echo "📄 文件内容: $latest_file"
            cat "$latest_file"
        fi
        
        echo ""
        read -p "立即应用此迁移? (y/n): " apply_now
        if [[ $apply_now == "y" ]]; then
            GOWORK=off atlas migrate apply --env dev
            echo "✅ 迁移已应用"
        fi
        ;;
    
    5)
        echo "⚠️  回滚迁移"
        GOWORK=off atlas migrate status --env dev
        echo ""
        read -p "确认回滚最后一次迁移? (y/n): " confirm
        if [[ $confirm == "y" ]]; then
            GOWORK=off atlas migrate down --env dev
            echo "✅ 迁移已回滚"
        fi
        ;;
    
    6)
        echo "🔍 验证迁移文件..."
        GOWORK=off atlas migrate validate --env dev
        echo "✅ 迁移文件验证通过"
        ;;
    
    *)
        echo "❌ 无效选择"
        exit 1
        ;;
esac

echo ""
echo "🎉 操作完成!"
echo ""
echo "💡 提示:"
echo "  - 查看完整指南: cat ../docs/guides/ATLAS_GUIDE.md"
echo "  - 迁移文件位置: internal/infrastructure/persistence/ent/migrations/"
echo "  - 配置文件: atlas.hcl"
