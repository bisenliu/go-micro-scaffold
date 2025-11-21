# Atlas 数据库迁移完整指南

> 本文档是 Atlas 数据库迁移工具在 `go-micro-scaffold` 项目中的完整使用指南。

---

## 📖 目录

- [快速开始](#快速开始)
- [什么是 Atlas](#什么是-atlas)
- [安装 Atlas](#安装-atlas)
- [基础用法](#基础用法)
- [常用命令](#常用命令)
- [Go Workspace 兼容性](#go-workspace-兼容性)
- [故障排查](#故障排查)
- [最佳实践](#最佳实践)
- [完整示例](#完整示例)

---

## 快速开始

### 🚀 三种使用方式

#### 方式 1：交互式脚本（推荐新手）

```bash
cd services
./atlas-quickstart.sh
```

脚本会自动：

- ✅ 检查并引导安装 Atlas
- ✅ 验证 MySQL 连接
- ✅ 提供交互式菜单操作

#### 方式 2：Makefile（推荐日常使用）

```bash
# 查看所有命令
make help

# 创建迁移
make migrate-create

# 应用迁移
make migrate-apply

# 查看状态
make migrate-status
```

#### 方式 3：直接使用 Atlas CLI

```bash
cd services

# 创建迁移
GOWORK=off atlas migrate diff <迁移名> --env dev

# 应用迁移
GOWORK=off atlas migrate apply --env dev

# 查看状态
GOWORK=off atlas migrate status --env dev
```

> **⚠️ 重要**: 由于项目使用 Go workspace 模式，所有 Atlas 命令都需要添加 `GOWORK=off` 前缀。

---

## 什么是 Atlas

[Atlas](https://atlasgo.io/) 是一个现代化的数据库迁移工具，由 Ent 团队开发，特别适合与 Ent ORM 集成使用。

### 核心特性

- ✅ **声明式迁移**: 从 Ent Schema 自动生成迁移脚本
- ✅ **版本控制**: 所有迁移文件都可以纳入 Git 版本管理
- ✅ **安全检查**: 自动检测破坏性变更（如删除表、删除列）
- ✅ **多环境支持**: 开发、测试、生产环境配置隔离
- ✅ **可视化**: 生成 ER 图和数据库文档

### 与 Ent 自动迁移对比

| 特性           | Atlas (推荐) | Ent 自动迁移 |
| -------------- | ------------ | ------------ |
| 版本化迁移     | ✅ 支持      | ❌ 不支持    |
| 迁移历史       | ✅ 保留文件  | ❌ 无历史    |
| 生产环境       | ✅ 适用      | ❌ 不推荐    |
| 回滚能力       | ✅ 支持      | ❌ 困难      |
| 团队协作       | ✅ Git 管理  | ❌ 难以协作  |
| 破坏性变更检测 | ✅ 自动检测  | ❌ 无检测    |

---

## 安装 Atlas

### macOS

```bash
brew install ariga/tap/atlas
```

### Linux

```bash
curl -sSf https://atlasgo.sh | sh
```

### 使用 Go 安装

```bash
go install ariga.io/atlas/cmd/atlas@latest
```

### 验证安装

```bash
atlas version
```

---

## 基础用法

### 1. 生成初始迁移

```bash
cd services
GOWORK=off atlas migrate diff initial --env dev
```

这会基于你的 Ent Schema 生成迁移文件：

```
services/internal/infrastructure/persistence/ent/migrations/
├── 20251121100000_initial.sql
└── atlas.sum
```

### 2. 应用迁移

```bash
GOWORK=off atlas migrate apply --env dev
```

### 3. 查看迁移状态

```bash
GOWORK=off atlas migrate status --env dev
```

输出示例：

```
Migration Status: OK
  -- Current Version: 20251121100000
  -- Next Version:    Already at latest version
  -- Executed Files:  1
  -- Pending Files:   0
```

### 4. 验证数据库

```bash
# 使用 Docker 中的 MySQL
docker exec -i mysql mysql -uroot go-micro-scaffold -e "SHOW TABLES;"
```

---

## 常用命令

### 创建新迁移

当你修改了 Ent Schema（如添加新字段、新表），需要生成迁移：

```bash
# 方式 1: 使用 Makefile
make migrate-create
# 提示输入迁移名称，如: add_user_email

# 方式 2: 直接使用 Atlas
GOWORK=off atlas migrate diff add_user_email --env dev
```

### 查看迁移内容

```bash
# 查看最新的迁移文件
cat services/internal/infrastructure/persistence/ent/migrations/*.sql | tail -n 50
```

### 回滚迁移

```bash
# 方式 1: 使用 Makefile（带确认）
make migrate-down

# 方式 2: 直接使用 Atlas
GOWORK=off atlas migrate down --env dev
```

### 验证迁移文件

```bash
# 验证迁移文件完整性
GOWORK=off atlas migrate validate --env dev
```

### 预览迁移（Dry Run）

```bash
# 预览将执行的 SQL，但不实际执行
GOWORK=off atlas migrate apply --env dev --dry-run
```

---

## Go Workspace 兼容性

### 问题说明

本项目使用 Go workspace 模式管理 `common` 和 `services` 两个模块。Atlas 在加载 Ent schema 时会遇到兼容性问题。

**错误示例**:

```
Error: loading ent schema: go: -mod may only be set to readonly or vendor when in workspace mode
```

### 解决方案

项目已自动配置好兼容性方案：

#### 1. services/go.mod 中添加了 replace 指令

```go
// services/go.mod
module services

go 1.24.1

// replace 指令使得在 GOWORK=off 模式下也能找到 common 模块
replace common => ../common

require (
    common v0.0.0
    // ... 其他依赖
)
```

#### 2. 所有命令使用 GOWORK=off

```bash
# ✅ 正确
GOWORK=off atlas migrate diff initial --env dev

# ❌ 错误（会报错）
atlas migrate diff initial --env dev
```

#### 3. 自动化工具已配置

- ✅ `atlas-quickstart.sh` - 所有命令已添加 `GOWORK=off`
- ✅ `Makefile` - 所有 migrate-* 命令已添加 `GOWORK=off`

**详细说明**: 参见 `services/docs/ATLAS_GOWORKSPACE_FIX.md`

---

## 故障排查

### 问题 1: Atlas 未安装

**错误**: `zsh: command not found: atlas`

**解决方案**:

```bash
# 使用快速启动脚本，会自动引导安装
cd services
./atlas-quickstart.sh

# 或使用 Makefile
make migrate-install
```

---

### 问题 2: MySQL 连接失败

**错误**: `Error 1045 (28000): Access denied for user 'root'`

**原因**: 数据库密码不匹配

**解决方案**:

1. 检查 MySQL 容器的实际密码：

```bash
docker inspect mysql | grep -A 5 "Env"
```

2. 修改 `services/atlas.hcl` 中的密码配置：

```hcl
variable "db_password" {
  type    = string
  default = ""  # 根据实际情况修改
}
```

---

### 问题 3: 数据库不干净

**错误**: `found table "xxx" in schema, baseline version or allow-dirty is required`

**原因**: 首次迁移时数据库已有表存在

**解决方案**:

```bash
# ⚠️ 仅开发环境！会删除所有数据
docker exec -i mysql mysql -uroot -e "
DROP DATABASE IF EXISTS \`go-micro-scaffold\`;
CREATE DATABASE \`go-micro-scaffold\`;
"

# 然后重新应用迁移
make migrate-apply
```

---

### 问题 4: 迁移文件哈希校验失败

**错误**: `checksum mismatch for file "xxx.sql"`

**解决方案**:

```bash
# 重新计算哈希
GOWORK=off atlas migrate hash --force --env dev
```

---

## 最佳实践

### ✅ DO（推荐做法）

#### 1. 总是先生成迁移文件再应用

```bash
# ✅ 正确流程
make migrate-create  # 生成迁移
# 查看 SQL 内容
cat services/internal/infrastructure/persistence/ent/migrations/*_<迁移名>.sql
# 确认无误后应用
make migrate-apply
```

#### 2. 使用描述性的迁移名称

```bash
# ✅ 好的命名
GOWORK=off atlas migrate diff add_user_email_and_phone_verified_fields --env dev

# ❌ 不好的命名
GOWORK=off atlas migrate diff update --env dev
```

#### 3. 将迁移文件纳入版本控制

```bash
git add services/internal/infrastructure/persistence/ent/migrations/
git commit -m "feat: add user email field migration"
git push
```

#### 4. 在开发环境测试迁移

```bash
# 先在开发环境验证
GOWORK=off atlas migrate apply --env dev --dry-run

# 确认无误后实际应用
GOWORK=off atlas migrate apply --env dev
```

#### 5. 破坏性迁移前备份数据

```bash
# 生产环境迁移前备份
docker exec mysql mysqldump -uroot go-micro-scaffold > backup_$(date +%Y%m%d).sql
```

### ❌ DON'T（避免做法）

#### 1. 不要直接编辑已应用的迁移文件

迁移文件有哈希校验，修改会导致验证失败。应该创建新的迁移文件。

#### 2. 不要跳过迁移文件的审查

特别是自动生成的迁移，可能不符合预期。始终检查 SQL 内容。

#### 3. 不要在生产环境使用 atlas schema apply

这会跳过迁移历史记录。应该使用 `atlas migrate apply`。

---

## 完整示例

### 场景 1：添加新字段

#### 步骤 1: 修改 Ent Schema

```go
// services/internal/infrastructure/persistence/ent/schema/user.go
func (User) Fields() []ent.Field {
    return []ent.Field{
        // ... 现有字段
        field.String("email").Optional(),  // 新增字段
        field.String("avatar").Optional(), // 新增字段
    }
}
```

#### 步骤 2: 生成 Ent 代码

```bash
cd services/internal/infrastructure/persistence/ent
go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
cd -
```

#### 步骤 3: 生成迁移

```bash
make migrate-create
# 输入迁移名称: add_user_email_and_avatar
```

#### 步骤 4: 查看生成的 SQL

```bash
cat services/internal/infrastructure/persistence/ent/migrations/*_add_user_email_and_avatar.sql
```

示例输出：

```sql
-- Add column "email" to table "users"
ALTER TABLE `users` ADD COLUMN `email` varchar(255) NULL;

-- Add column "avatar" to table "users"
ALTER TABLE `users` ADD COLUMN `avatar` varchar(255) NULL;
```

#### 步骤 5: 应用迁移

```bash
make migrate-apply
```

#### 步骤 6: 验证

```bash
docker exec -i mysql mysql -uroot go-micro-scaffold -e "DESC users;"
```

#### 步骤 7: 提交到 Git

```bash
git add services/internal/infrastructure/persistence/ent/
git commit -m "feat: add email and avatar fields to user"
git push
```

---

### 场景 2：创建新表

#### 步骤 1: 创建新的 Ent Schema

```go
// services/internal/infrastructure/persistence/ent/schema/order.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "github.com/google/uuid"
    "time"
)

type Order struct {
    ent.Schema
}

func (Order) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New),
        field.String("user_id").NotEmpty(),
        field.Int64("total_amount"),
        field.Int("status").Default(1),
        field.Time("created_at").Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}
```

#### 步骤 2-7: 同场景 1

迁移名称建议: `create_orders_table`

---

## 配置文件说明

### atlas.hcl

项目的 Atlas 配置文件位于 `services/atlas.hcl`，定义了三个环境：

```hcl
# 开发环境
env "dev" {
  url = "mysql://root:@localhost:3306/go-micro-scaffold?parseTime=true"
  migration {
    dir = "file://internal/infrastructure/persistence/ent/migrations"
  }
  src = "ent://internal/infrastructure/persistence/ent/schema"
  dev = "docker://mysql/8/dev"
}

# Docker 环境
env "docker" {
  url = "mysql://root:root@mysql:3306/go-micro-scaffold?parseTime=true"
  # ... 其他配置
}

# 生产环境
env "prod" {
  url = "mysql://root:password@prod-db:3306/go-micro-scaffold?parseTime=true"
  # 生产环境的保护措施
  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
  backup = true
}
```

### 环境变量覆盖

可以通过环境变量覆盖配置：

```bash
DB_HOST=192.168.1.100 \
DB_PORT=3307 \
DB_NAME=my_database \
DB_PASSWORD=secret \
GOWORK=off atlas migrate apply --env dev
```

---

## 项目文件结构

```
go-micro-scaffold/
├── services/
│   ├── atlas.hcl                          # Atlas 配置文件
│   ├── atlas-quickstart.sh                # 快速启动脚本
│   ├── go.mod                             # 包含 replace 指令
│   ├── docs/
│   │   ├── ATLAS_MIGRATION_GUIDE.md       # 详细指南
│   │   ├── ATLAS_COMMANDS_CHEATSHEET.md   # 命令速查表
│   │   └── ATLAS_GOWORKSPACE_FIX.md       # Go Workspace 兼容性说明
│   └── internal/infrastructure/persistence/ent/
│       ├── schema/                        # Ent Schema 定义
│       └── migrations/                    # 迁移文件目录
│           ├── 20251121100000_initial.sql
│           └── atlas.sum
├── Makefile                               # 项目命令工具
└── ATLAS_GUIDE.md                         # 本文档
```

---

## 相关资源

### 项目文档

- **本文档**: `ATLAS_GUIDE.md` - Atlas 使用完整指南
- **详细教程**: `services/docs/ATLAS_MIGRATION_GUIDE.md` - 35KB+ 详细教程
- **命令速查**: `services/docs/ATLAS_COMMANDS_CHEATSHEET.md` - 常用命令参考
- **兼容性说明**: `services/docs/ATLAS_GOWORKSPACE_FIX.md` - Go Workspace 问题解决

### 官方文档

- [Atlas 官网](https://atlasgo.io/)
- [Ent 迁移指南](https://entgo.io/docs/versioned-migrations/)
- [Atlas CLI 参考](https://atlasgo.io/cli-reference)
- [Atlas Schema HCL 语法](https://atlasgo.io/atlas-schema/hcl)

### 工具链接

- [Atlas GitHub](https://github.com/ariga/atlas)
- [Ent GitHub](https://github.com/ent/ent)

---

## 快速命令参考

```bash
# === 安装 ===
brew install ariga/tap/atlas         # macOS
curl -sSf https://atlasgo.sh | sh    # Linux

# === 基础操作 ===
make migrate-create                   # 创建迁移
make migrate-apply                    # 应用迁移
make migrate-status                   # 查看状态
make migrate-down                     # 回滚迁移
make migrate-validate                 # 验证迁移

# === 数据库操作 ===
docker exec -i mysql mysql -uroot go-micro-scaffold -e "SHOW TABLES;"
docker exec -i mysql mysql -uroot go-micro-scaffold -e "DESC users;"

# === 其他 ===
make help                            # 查看所有命令
./services/atlas-quickstart.sh      # 交互式脚本
```

---
