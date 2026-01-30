# Gitea 插件系统 Docker 部署指南

本指南介绍如何将 Gitea 插件系统编译为 Docker 镜像并部署到云端。

## 📋 目录

- [前置要求](#前置要求)
- [本地构建](#本地构建)
- [推送到云端](#推送到云端)
- [云端部署](#云端部署)
- [配置说明](#配置说明)
- [故障排除](#故障排除)

## 前置要求

### 1. 安装 Docker

**Windows:**
```bash
# 下载并安装 Docker Desktop
# https://www.docker.com/products/docker-desktop
```

**Linux:**
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 添加当前用户到 docker 组
sudo usermod -aG docker $USER
```

**macOS:**
```bash
# 下载并安装 Docker Desktop
# https://www.docker.com/products/docker-desktop
```

### 2. 注册 Docker 仓库

选择一个 Docker 仓库服务：

- **Docker Hub**: https://hub.docker.com
- **阿里云容器镜像服务**: https://cr.console.aliyun.com
- **腾讯云容器镜像服务**: https://console.cloud.tencent.com/tcr
- **私有仓库**: Harbor, GitLab Container Registry 等

## 本地构建

### 方式 1：使用构建脚本（推荐）

#### Linux/macOS:

```bash
cd gitea-license

# 修改配置
vim build-and-push.sh
# 修改 DOCKER_REGISTRY 为你的仓库地址
# 例如: DOCKER_REGISTRY="registry.cn-hangzhou.aliyuncs.com/your-namespace"

# 添加执行权限
chmod +x build-and-push.sh

# 登录 Docker 仓库
docker login your-registry.com

# 执行构建和推送
./build-and-push.sh
```

#### Windows:

```cmd
cd gitea-license

# 修改配置
notepad build-and-push.bat
# 修改 DOCKER_REGISTRY 为你的仓库地址

# 登录 Docker 仓库
docker login your-registry.com

# 执行构建和推送
build-and-push.bat
```

### 方式 2：手动构建

```bash
cd gitea-license

# 1. 构建镜像
docker build -f Dockerfile.custom -t gitea-with-plugins:1.0.0 .

# 2. 标记镜像
docker tag gitea-with-plugins:1.0.0 your-registry.com/gitea-with-plugins:1.0.0
docker tag gitea-with-plugins:1.0.0 your-registry.com/gitea-with-plugins:latest

# 3. 登录仓库
docker login your-registry.com

# 4. 推送镜像
docker push your-registry.com/gitea-with-plugins:1.0.0
docker push your-registry.com/gitea-with-plugins:latest
```

## 推送到云端

### Docker Hub

```bash
# 登录
docker login

# 标记镜像
docker tag gitea-with-plugins:1.0.0 your-username/gitea-with-plugins:1.0.0

# 推送
docker push your-username/gitea-with-plugins:1.0.0
```

### 阿里云容器镜像服务

```bash
# 登录
docker login --username=your-username registry.cn-hangzhou.aliyuncs.com

# 标记镜像
docker tag gitea-with-plugins:1.0.0 registry.cn-hangzhou.aliyuncs.com/your-namespace/gitea-with-plugins:1.0.0

# 推送
docker push registry.cn-hangzhou.aliyuncs.com/your-namespace/gitea-with-plugins:1.0.0
```

### 腾讯云容器镜像服务

```bash
# 登录
docker login --username=your-username ccr.ccs.tencentyun.com

# 标记镜像
docker tag gitea-with-plugins:1.0.0 ccr.ccs.tencentyun.com/your-namespace/gitea-with-plugins:1.0.0

# 推送
docker push ccr.ccs.tencentyun.com/your-namespace/gitea-with-plugins:1.0.0
```

## 云端部署

### 方式 1：使用 Docker Compose（推荐）

#### 1. 创建部署目录

```bash
mkdir -p ~/gitea-deploy
cd ~/gitea-deploy
```

#### 2. 创建 docker-compose.yml

```yaml
version: "3.8"

services:
  gitea:
    image: your-registry.com/gitea-with-plugins:latest
    container_name: gitea
    environment:
      - USER_UID=1000
      - USER_GID=1000
      - GITEA__database__DB_TYPE=postgres
      - GITEA__database__HOST=db:5432
      - GITEA__database__NAME=gitea
      - GITEA__database__USER=gitea
      - GITEA__database__PASSWD=gitea_password_here
      - GITEA__plugin__ENABLED=true
      - GITEA__plugin__PLUGINS_DIR=/data/gitea/plugins
      - GITEA__server__DOMAIN=your-domain.com
      - GITEA__server__ROOT_URL=https://your-domain.com
    restart: always
    networks:
      - gitea
    volumes:
      - gitea-data:/data
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    ports:
      - "3000:3000"
      - "2222:22"
    depends_on:
      - db

  db:
    image: postgres:14-alpine
    container_name: gitea-db
    restart: always
    environment:
      - POSTGRES_USER=gitea
      - POSTGRES_PASSWORD=gitea_password_here
      - POSTGRES_DB=gitea
    networks:
      - gitea
    volumes:
      - postgres-data:/var/lib/postgresql/data

networks:
  gitea:
    driver: bridge

volumes:
  gitea-data:
    driver: local
  postgres-data:
    driver: local
```

#### 3. 启动服务

```bash
# 启动
docker-compose up -d

# 查看日志
docker-compose logs -f gitea

# 停止
docker-compose down

# 重启
docker-compose restart
```

### 方式 2：使用 Docker 命令

```bash
# 创建网络
docker network create gitea

# 启动数据库
docker run -d \
  --name gitea-db \
  --network gitea \
  -e POSTGRES_USER=gitea \
  -e POSTGRES_PASSWORD=gitea \
  -e POSTGRES_DB=gitea \
  -v postgres-data:/var/lib/postgresql/data \
  postgres:14-alpine

# 启动 Gitea
docker run -d \
  --name gitea \
  --network gitea \
  -p 3000:3000 \
  -p 2222:22 \
  -e USER_UID=1000 \
  -e USER_GID=1000 \
  -e GITEA__database__DB_TYPE=postgres \
  -e GITEA__database__HOST=gitea-db:5432 \
  -e GITEA__database__NAME=gitea \
  -e GITEA__database__USER=gitea \
  -e GITEA__database__PASSWD=gitea \
  -e GITEA__plugin__ENABLED=true \
  -v gitea-data:/data \
  your-registry.com/gitea-with-plugins:latest
```

### 方式 3：Kubernetes 部署

#### 1. 创建 Namespace

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: gitea
```

#### 2. 创建 ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gitea-config
  namespace: gitea
data:
  app.ini: |
    [database]
    DB_TYPE = postgres
    HOST = gitea-db:5432
    NAME = gitea
    USER = gitea
    PASSWD = gitea
    
    [plugin]
    ENABLED = true
    PLUGINS_DIR = /data/gitea/plugins
    
    [server]
    DOMAIN = your-domain.com
    ROOT_URL = https://your-domain.com
```

#### 3. 创建 Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gitea
  namespace: gitea
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitea
  template:
    metadata:
      labels:
        app: gitea
    spec:
      containers:
      - name: gitea
        image: your-registry.com/gitea-with-plugins:latest
        ports:
        - containerPort: 3000
        - containerPort: 22
        env:
        - name: USER_UID
          value: "1000"
        - name: USER_GID
          value: "1000"
        volumeMounts:
        - name: gitea-data
          mountPath: /data
        - name: config
          mountPath: /data/gitea/conf
      volumes:
      - name: gitea-data
        persistentVolumeClaim:
          claimName: gitea-pvc
      - name: config
        configMap:
          name: gitea-config
```

#### 4. 创建 Service

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: gitea
  namespace: gitea
spec:
  type: LoadBalancer
  ports:
  - name: http
    port: 3000
    targetPort: 3000
  - name: ssh
    port: 22
    targetPort: 22
  selector:
    app: gitea
```

#### 5. 部署

```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# 查看状态
kubectl get pods -n gitea
kubectl logs -f -n gitea deployment/gitea
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `USER_UID` | 运行用户 UID | 1000 |
| `USER_GID` | 运行用户 GID | 1000 |
| `GITEA__database__DB_TYPE` | 数据库类型 | postgres |
| `GITEA__database__HOST` | 数据库地址 | db:5432 |
| `GITEA__database__NAME` | 数据库名称 | gitea |
| `GITEA__database__USER` | 数据库用户 | gitea |
| `GITEA__database__PASSWD` | 数据库密码 | - |
| `GITEA__plugin__ENABLED` | 启用插件系统 | true |
| `GITEA__plugin__PLUGINS_DIR` | 插件目录 | /data/gitea/plugins |
| `GITEA__server__DOMAIN` | 域名 | localhost |
| `GITEA__server__ROOT_URL` | 根 URL | http://localhost:3000 |

### 数据持久化

重要的数据目录：

- `/data/gitea` - Gitea 数据目录
- `/data/git` - Git 仓库目录
- `/data/gitea/plugins` - 插件目录

确保这些目录使用 Docker Volume 或持久化存储。

### 反向代理配置

#### Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### Traefik

```yaml
# docker-compose.yml
services:
  gitea:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.gitea.rule=Host(`your-domain.com`)"
      - "traefik.http.routers.gitea.entrypoints=websecure"
      - "traefik.http.routers.gitea.tls.certresolver=letsencrypt"
      - "traefik.http.services.gitea.loadbalancer.server.port=3000"
```

## 故障排除

### 1. 镜像构建失败

**问题**: 编译时内存不足

**解决方案**:
```bash
# 增加 Docker 内存限制
# Docker Desktop -> Settings -> Resources -> Memory
# 建议至少 4GB
```

### 2. 插件无法加载

**问题**: 插件目录权限问题

**解决方案**:
```bash
# 进入容器
docker exec -it gitea bash

# 检查权限
ls -la /data/gitea/plugins

# 修复权限
chown -R git:git /data/gitea/plugins
```

### 3. 数据库连接失败

**问题**: 数据库未就绪

**解决方案**:
```yaml
# 在 docker-compose.yml 中添加健康检查
services:
  db:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gitea"]
      interval: 10s
      timeout: 5s
      retries: 5
  
  gitea:
    depends_on:
      db:
        condition: service_healthy
```

### 4. 端口冲突

**问题**: 端口已被占用

**解决方案**:
```yaml
# 修改端口映射
ports:
  - "8080:3000"  # 使用 8080 代替 3000
  - "2222:22"
```

### 5. 查看日志

```bash
# Docker Compose
docker-compose logs -f gitea

# Docker
docker logs -f gitea

# Kubernetes
kubectl logs -f -n gitea deployment/gitea
```

## 更新和维护

### 更新镜像

```bash
# 1. 构建新版本
docker build -f Dockerfile.custom -t gitea-with-plugins:1.1.0 .

# 2. 推送到仓库
docker push your-registry.com/gitea-with-plugins:1.1.0

# 3. 更新部署
docker-compose pull
docker-compose up -d
```

### 备份数据

```bash
# 备份数据卷
docker run --rm \
  -v gitea-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/gitea-backup-$(date +%Y%m%d).tar.gz /data

# 备份数据库
docker exec gitea-db pg_dump -U gitea gitea > gitea-db-backup-$(date +%Y%m%d).sql
```

### 恢复数据

```bash
# 恢复数据卷
docker run --rm \
  -v gitea-data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/gitea-backup-20260130.tar.gz -C /

# 恢复数据库
docker exec -i gitea-db psql -U gitea gitea < gitea-db-backup-20260130.sql
```

## 生产环境建议

### 1. 使用 HTTPS

- 配置 SSL 证书
- 使用 Let's Encrypt 自动续期
- 强制 HTTPS 重定向

### 2. 配置备份

- 定期备份数据
- 备份到远程存储
- 测试恢复流程

### 3. 监控和日志

- 配置日志收集
- 设置监控告警
- 定期检查健康状态

### 4. 安全加固

- 限制容器权限
- 使用非 root 用户
- 定期更新镜像
- 扫描安全漏洞

### 5. 性能优化

- 使用 SSD 存储
- 配置缓存
- 优化数据库
- 使用 CDN

## 云服务商部署

### 阿里云 ECS

```bash
# 1. 购买 ECS 实例
# 2. 安装 Docker
# 3. 拉取镜像
docker pull registry.cn-hangzhou.aliyuncs.com/your-namespace/gitea-with-plugins:latest

# 4. 启动服务
docker-compose up -d
```

### 腾讯云 CVM

```bash
# 1. 购买 CVM 实例
# 2. 安装 Docker
# 3. 拉取镜像
docker pull ccr.ccs.tencentyun.com/your-namespace/gitea-with-plugins:latest

# 4. 启动服务
docker-compose up -d
```

### AWS EC2

```bash
# 1. 启动 EC2 实例
# 2. 安装 Docker
# 3. 拉取镜像
docker pull your-ecr-registry/gitea-with-plugins:latest

# 4. 启动服务
docker-compose up -d
```

## 总结

完成以上步骤后，你的 Gitea 插件系统就成功部署到云端了！

访问 `http://your-domain.com:3000` 开始使用。

## 获取帮助

- 查看日志：`docker-compose logs -f`
- 进入容器：`docker exec -it gitea bash`
- 检查状态：`docker-compose ps`
- 重启服务：`docker-compose restart`

---

**部署完成！** 🎉
