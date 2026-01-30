# GitHub Actions 自动构建指南

## 🎯 概述

使用 GitHub Actions 自动构建 Docker 镜像并推送到阿里云容器镜像服务（华南2-河源）。

## 📦 镜像信息

- **镜像地址**：`registry.cn-heyuan.aliyuncs.com/harti123/gitea-plugins`
- **区域**：华南2（河源）
- **命名空间**：`harti123`

## 🚀 配置步骤

### 步骤 1：在 GitHub 创建仓库

1. 访问：https://github.com/new
2. 仓库名称：`gitea-plugins`
3. 设置为私有或公开
4. 不要初始化 README、.gitignore 或 LICENSE
5. 点击"Create repository"

### 步骤 2：添加 GitHub 远程仓库

在本地执行：

```bash
cd gitea-license

# 添加 GitHub 远程仓库（替换为你的 GitHub 用户名）
git remote add github https://github.com/你的用户名/gitea-plugins.git

# 推送代码到 GitHub
git push github main
```

### 步骤 3：配置 GitHub Secrets

在 GitHub 仓库中配置密钥：

1. **进入仓库设置**
   - 访问：`https://github.com/你的用户名/gitea-plugins/settings/secrets/actions`
   - 或者：仓库页面 → Settings → Secrets and variables → Actions

2. **添加以下 Secrets**

   点击"New repository secret"，添加：

   | Secret 名称 | 值 | 说明 |
   |------------|-----|------|
   | `DOCKER_USERNAME` | 你的阿里云账号用户名 | 用于登录阿里云镜像服务 |
   | `DOCKER_PASSWORD` | 你的阿里云镜像服务密码 | 固定密码，不是登录密码 |

3. **获取阿里云镜像服务密码**

   如果还没有设置：
   - 访问：https://cr.console.aliyun.com
   - 点击右上角头像 → 访问凭证
   - 设置固定密码
   - 记录用户名和密码

### 步骤 4：触发自动构建

配置完成后，每次推送代码到 GitHub 都会自动触发构建：

```bash
# 方式 1：推送代码触发
git add .
git commit -m "触发自动构建"
git push github main

# 方式 2：创建标签触发（推荐用于版本发布）
git tag -a v1.0.0 -m "Release version 1.0.0"
git push github v1.0.0
```

### 步骤 5：查看构建状态

1. **访问 Actions 页面**
   - `https://github.com/你的用户名/gitea-plugins/actions`

2. **查看构建日志**
   - 点击最新的 workflow run
   - 查看详细日志
   - 等待构建完成（约 15-30 分钟）

3. **构建成功标志**
   - ✅ 绿色对勾表示成功
   - ❌ 红色叉号表示失败

## 📊 构建完成后

构建成功后，镜像会自动推送到：

```
registry.cn-heyuan.aliyuncs.com/harti123/gitea-plugins:latest
```

### 在服务器上使用

```bash
# 登录阿里云镜像服务
docker login --username=你的用户名 registry.cn-heyuan.aliyuncs.com

# 拉取镜像
docker pull registry.cn-heyuan.aliyuncs.com/harti123/gitea-plugins:latest

# 运行容器
docker run -d \
  --name gitea \
  -p 3000:3000 \
  -p 2222:22 \
  -v gitea-data:/data \
  registry.cn-heyuan.aliyuncs.com/harti123/gitea-plugins:latest
```

## 🔄 同时推送到 Gitee 和 GitHub

你可以同时维护两个仓库：

```bash
# 查看当前远程仓库
git remote -v

# 同时推送到两个仓库
git push origin main    # 推送到 Gitee
git push github main    # 推送到 GitHub

# 或者一次性推送到所有远程仓库
git remote set-url --add --push origin https://gitee.com/harti/gitea-plugins.git
git remote set-url --add --push origin https://github.com/你的用户名/gitea-plugins.git
git push origin main    # 同时推送到两个仓库
```

## 🎯 自动化工作流

配置完成后，你的工作流程是：

1. **本地开发**
   ```bash
   # 修改代码
   git add .
   git commit -m "更新功能"
   ```

2. **推送到 GitHub**
   ```bash
   git push github main
   ```

3. **自动构建**
   - GitHub Actions 自动触发
   - 构建 Docker 镜像
   - 推送到阿里云

4. **在服务器更新**
   ```bash
   docker pull registry.cn-heyuan.aliyuncs.com/harti123/gitea-plugins:latest
   docker-compose restart gitea
   ```

## 🐛 故障排除

### 问题 1：构建失败 - Authentication failed

**原因**：GitHub Secrets 配置错误

**解决**：
1. 检查 `DOCKER_USERNAME` 是否正确
2. 检查 `DOCKER_PASSWORD` 是否是镜像服务密码（不是登录密码）
3. 重新设置 Secrets

### 问题 2：构建超时

**原因**：Gitea 源码较大，构建时间长

**解决**：
- 等待更长时间（首次构建可能需要 30 分钟）
- GitHub Actions 免费版有 6 小时限制，足够使用

### 问题 3：推送失败 - Permission denied

**原因**：没有推送权限

**解决**：
```bash
# 检查远程仓库 URL
git remote -v

# 使用 HTTPS 并输入 GitHub Token
git push https://你的用户名:你的Token@github.com/你的用户名/gitea-plugins.git main
```

### 问题 4：Dockerfile not found

**原因**：Dockerfile 路径错误

**解决**：
- 确认 `Dockerfile.custom` 在仓库根目录
- 检查 `.github/workflows/docker-build.yml` 中的路径配置

## 📝 GitHub Actions 配置说明

当前配置文件：`.github/workflows/docker-build.yml`

**触发条件**：
- 推送到 `main` 或 `master` 分支
- 创建以 `v` 开头的标签（如 `v1.0.0`）
- Pull Request 到 `main` 或 `master` 分支

**构建平台**：
- `linux/amd64`（x86_64）
- `linux/arm64`（ARM64）

**镜像标签**：
- `latest`：最新的 main 分支构建
- `v1.0.0`：版本标签
- `main`：分支名称

## 🎉 优势

使用 GitHub Actions 的优势：

1. ✅ **完全免费**：公开仓库无限制，私有仓库每月 2000 分钟
2. ✅ **自动化**：推送代码自动构建
3. ✅ **多平台**：支持 AMD64 和 ARM64
4. ✅ **快速**：GitHub 服务器速度快
5. ✅ **可靠**：构建失败会收到邮件通知
6. ✅ **缓存**：使用 GitHub Actions Cache 加速构建

## 📞 需要帮助？

- GitHub 仓库：`https://github.com/你的用户名/gitea-plugins`
- Gitee 仓库：https://gitee.com/harti/gitea-plugins
- GitHub Actions 文档：https://docs.github.com/actions
- 阿里云镜像服务：https://cr.console.aliyun.com/cn-heyuan/instances/repositories

---

**下一步**：创建 GitHub 仓库并配置 Secrets，然后推送代码触发自动构建！🚀
