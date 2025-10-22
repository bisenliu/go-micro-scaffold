# 部署指南

## 概述

本指南详细说明如何在不同环境中部署 Go Micro Scaffold 应用，包括 Swagger 文档的环境特定配置和最佳实践。

## 🚀 快速部署

### 使用 Make 命令

```bash
# 开发环境部署
make deploy-dev

# 预发布环境部署
make deploy-staging

# 生产环境部署
make deploy-prod
```

### 使用部署脚本

```bash
# 开发环境
./scripts/deploy.sh dev

# 预发布环境
./scripts/deploy.sh staging --docker

# 生产环境（需要确认）
./scripts/deploy.sh prod --force
```

## 🏗️ 环境配置

### 开发环境 (Development)

**特点：**
- ✅ Swagger UI 完全启用
- ✅ 详细日志输出
- ✅ 调试功能开启
- ✅ 热重载支持

**配置：**
```yaml
swagger:
  enabled: true
  title: "Go Micro Scaffold API (Development)"
  host: "localhost:8080"

logging:
  level: "debug"
  format: "console"

system:
  debug: true
```

**部署命令：**
```bash
# 使用 Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# 使用部署脚本
./scripts/deploy.sh dev --docker

# 直接运行
make dev && make run
```

**访问地址：**
- 应用: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- 数据库管理: http://localhost:8081 (Adminer)
- Redis 管理: http://localhost:8082 (Redis Commander)

### 预发布环境 (Staging)

**特点：**
- ✅ Swagger UI 启用（可能需要认证）
- ⚠️ 生产级配置但允许调试
- ✅ 完整监控和日志
- ✅ 安全头部配置

**配置：**
```yaml
swagger:
  enabled: true
  title: "Go Micro Scaffold API (Staging)"
  host: "staging-api.example.com"

logging:
  level: "info"
  format: "json"

security:
  tls:
    enabled: true
```

**部署命令：**
```bash
# 使用 Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up -d

# 使用部署脚本
./scripts/deploy.sh staging --docker

# 使用 Make
make deploy-staging
```

**环境变量：**
```bash
export DB_HOST=staging-db.example.com
export DB_PASSWORD=staging_secure_password
export REDIS_HOST=staging-redis.example.com
export JWT_SECRET=staging-jwt-secret-key
```

### 生产环境 (Production)

**特点：**
- ❌ Swagger UI 禁用（安全考虑）
- ✅ 最高安全级别
- ✅ 性能优化配置
- ✅ 完整监控和告警

**配置：**
```yaml
swagger:
  enabled: false  # 生产环境禁用

logging:
  level: "warn"
  format: "json"
  output: ["file"]

security:
  tls:
    enabled: true
    min_version: "1.2"
```

**部署命令：**
```bash
# 使用 Docker Compose
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 使用部署脚本（需要确认）
./scripts/deploy.sh prod --force

# 使用 Make
make deploy-prod
```

**环境变量：**
```bash
export DB_HOST=prod-db.example.com
export DB_PASSWORD=production_secure_password
export REDIS_HOST=prod-redis.example.com
export JWT_SECRET=production-jwt-secret-key-32-chars
export SYSTEM_SECRET_KEY=production-system-secret-32-chars
```

## 🐳 Docker 部署

### 单容器部署

```bash
# 构建镜像
docker build -t go-micro-scaffold:latest .

# 运行容器
docker run -d \
  --name go-micro-scaffold \
  -p 8080:8080 \
  -e GO_ENV=production \
  -e DB_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  go-micro-scaffold:latest
```

### Docker Compose 部署

```bash
# 开发环境
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# 预发布环境
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up -d

# 生产环境
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 多阶段构建

Dockerfile 支持多阶段构建：

```bash
# 开发镜像（包含源码和工具）
docker build --target development -t go-micro-scaffold:dev .

# 生产镜像（最小化）
docker build --target production -t go-micro-scaffold:prod .
```

## ☸️ Kubernetes 部署

### 基础部署

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-micro-scaffold
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-micro-scaffold
  template:
    metadata:
      labels:
        app: go-micro-scaffold
    spec:
      containers:
      - name: app
        image: go-micro-scaffold:latest
        ports:
        - containerPort: 8080
        env:
        - name: GO_ENV
          value: "production"
        - name: SWAGGER_ENABLED
          value: "false"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

### 服务配置

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: go-micro-scaffold-service
spec:
  selector:
    app: go-micro-scaffold
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Ingress 配置

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-micro-scaffold-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - api.example.com
    secretName: api-tls
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-micro-scaffold-service
            port:
              number: 80
```

## 🔧 CI/CD 集成

### GitHub Actions

项目包含完整的 GitHub Actions 工作流：

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    # 代码质量检查和测试
  
  build:
    # 构建应用和 Docker 镜像
  
  deploy-dev:
    # 自动部署到开发环境
  
  deploy-staging:
    # 自动部署到预发布环境
  
  deploy-prod:
    # 手动部署到生产环境
```

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

variables:
  DOCKER_IMAGE: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

test:
  stage: test
  script:
    - make test
    - make swagger-validate

build:
  stage: build
  script:
    - docker build -t $DOCKER_IMAGE .
    - docker push $DOCKER_IMAGE

deploy_staging:
  stage: deploy
  script:
    - ./scripts/deploy.sh staging --docker
  environment:
    name: staging
    url: https://staging-api.example.com
  only:
    - main

deploy_production:
  stage: deploy
  script:
    - ./scripts/deploy.sh prod --force
  environment:
    name: production
    url: https://api.example.com
  when: manual
  only:
    - tags
```

## 📊 监控和日志

### Prometheus 监控

```yaml
# 监控配置
monitoring:
  enabled: true
  prometheus:
    enabled: true
    path: "/metrics"
```

**访问地址：**
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

### 日志聚合

```yaml
# 日志配置
logging:
  level: "info"
  format: "json"
  output: ["stdout", "file"]
```

**日志查看：**
```bash
# Docker 日志
docker logs go-micro-scaffold

# 文件日志
tail -f logs/app.log

# 结构化查询
cat logs/app.log | jq '.level == "error"'
```

## 🔒 安全配置

### SSL/TLS 配置

```bash
# 生成自签名证书（开发环境）
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# 使用 Let's Encrypt（生产环境）
certbot certonly --webroot -w /var/www/html -d api.example.com
```

### 环境变量管理

```bash
# 使用 .env 文件（开发环境）
cp .env.example .env

# 使用 Kubernetes Secrets（生产环境）
kubectl create secret generic app-secrets \
  --from-literal=db-password=secure-password \
  --from-literal=jwt-secret=jwt-secret-key
```

### 网络安全

```yaml
# 防火墙规则
iptables -A INPUT -p tcp --dport 8080 -s 10.0.0.0/8 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP

# Docker 网络隔离
networks:
  app-network:
    driver: bridge
    internal: true
```

## 🚨 故障排除

### 常见问题

**1. Swagger 文档不显示**
```bash
# 检查 Swagger 配置
grep -r "swagger.enabled" configs/

# 重新生成文档
make swagger-gen

# 验证文档
make swagger-validate
```

**2. 数据库连接失败**
```bash
# 检查数据库连接
docker exec -it go-micro-scaffold-mysql mysql -u root -p

# 查看连接配置
docker logs go-micro-scaffold | grep -i database
```

**3. 容器启动失败**
```bash
# 查看容器日志
docker logs go-micro-scaffold

# 检查健康状态
docker inspect go-micro-scaffold | jq '.[0].State.Health'

# 进入容器调试
docker exec -it go-micro-scaffold sh
```

### 性能调优

**1. 数据库优化**
```sql
-- 查看慢查询
SHOW VARIABLES LIKE 'slow_query_log';
SHOW VARIABLES LIKE 'long_query_time';

-- 优化配置
SET GLOBAL innodb_buffer_pool_size = 256M;
SET GLOBAL max_connections = 200;
```

**2. Redis 优化**
```bash
# 内存使用情况
redis-cli INFO memory

# 性能监控
redis-cli MONITOR
```

**3. 应用优化**
```bash
# 启用 pprof（开发环境）
go tool pprof http://localhost:8080/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap
```

## 📋 部署检查清单

### 部署前检查

- [ ] 代码已通过所有测试
- [ ] Swagger 文档已生成并验证
- [ ] 环境变量已正确配置
- [ ] 数据库迁移已完成
- [ ] SSL 证书已配置（生产环境）
- [ ] 监控和告警已设置
- [ ] 备份策略已实施

### 部署后验证

- [ ] 应用健康检查通过
- [ ] API 端点正常响应
- [ ] 数据库连接正常
- [ ] Redis 缓存工作正常
- [ ] 日志输出正常
- [ ] 监控指标正常
- [ ] Swagger UI 访问正常（开发/预发布环境）

### 生产环境特殊检查

- [ ] Swagger UI 已禁用
- [ ] 调试功能已关闭
- [ ] 敏感信息已隐藏
- [ ] 安全头部已配置
- [ ] 访问控制已启用
- [ ] 备份和恢复已测试

## 🔄 回滚策略

### 快速回滚

```bash
# Docker 回滚
docker stop go-micro-scaffold
docker run -d --name go-micro-scaffold go-micro-scaffold:previous-version

# Kubernetes 回滚
kubectl rollout undo deployment/go-micro-scaffold

# 使用部署脚本回滚
./scripts/deploy.sh prod --version=v1.0.0 --force
```

### 数据库回滚

```bash
# 使用 CLI 工具回滚迁移
./bin/cli migrate down

# 从备份恢复
mysql -u root -p go_micro_scaffold < backup_20231201.sql
```

---

通过遵循本部署指南，可以确保应用在不同环境中的稳定运行，同时保证 Swagger 文档的适当配置和安全性。