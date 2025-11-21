# Atlas 配置文件
# 文档: https://atlasgo.io/atlas-schema/hcl
# 
# 重要提示：在 Go workspace 模式下使用 Atlas 时，需要设置环境变量 GOWORK=off
# 例如: GOWORK=off atlas migrate diff initial --env dev

# 定义环境变量
variable "db_host" {
  type    = string
  default = "localhost"
}

variable "db_port" {
  type    = number
  default = 3306
}

variable "db_name" {
  type    = string
  default = "go-micro-scaffold"
}

variable "db_user" {
  type    = string
  default = "root"
}

variable "db_password" {
  type    = string
  default = ""  # MySQL 容器使用空密码
}

# 定义开发环境
env "dev" {
  # 数据库连接 URL
  url = "mysql://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?parseTime=true"

  # 迁移文件存储目录
  migration {
    dir = "file://internal/infrastructure/persistence/ent/migrations"
  }

  # 从 Ent Schema 生成迁移
  # 注意：在 Go workspace 模式下运行时，需要设置 GOWORK=off 环境变量
  src = "ent://internal/infrastructure/persistence/ent/schema"

  # 开发数据库 URL（用于验证迁移）
  dev = "docker://mysql/8/dev"

  # 格式化配置
  format {
    migrate {
      diff = "{{ sql . \"  \" }}" # 缩进 SQL
    }
  }
}

# 定义生产环境
env "prod" {
  url = "mysql://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?parseTime=true"

  migration {
    dir = "file://internal/infrastructure/persistence/ent/migrations"
  }

  # 生产环境的保护措施
  diff {
    # 跳过破坏性更改的自动应用
    skip {
      drop_schema = true
      drop_table  = true
    }
  }

  # 迁移前备份
  backup = true
}

# Docker 环境（用于 docker-compose）
env "docker" {
  url = "mysql://root:root@mysql:3306/go-micro-scaffold?parseTime=true"

  migration {
    dir = "file://internal/infrastructure/persistence/ent/migrations"
  }

  src = "ent://internal/infrastructure/persistence/ent/schema"
  dev = "docker://mysql/8/dev"
}
