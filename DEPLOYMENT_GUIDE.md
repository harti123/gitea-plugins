# Gitea 插件系统部署指南

本指南将帮助你完成 Gitea 插件系统的完整部署。

## 前置要求

- Go 1.21 或更高版本
- Git
- Make
- 数据库（SQLite/MySQL/PostgreSQL）

## 步骤 1：编译 Gitea

```bash
cd gitea-license

# 编译（包含 SQLite 支持）
TAGS="bindata sqlite sqlite_unlock_notify" make build

# 或者使用其他数据库
# MySQL: TAGS="bindata" make build
# PostgreSQL: TAGS="bindata" make build
```

编译完成后，会在当前目录生成 `gitea` 可执行文件。

## 步骤 2：配置 Gitea

### 2.1 创建配置目录

```bash
mkdir -p custom/conf
```

### 2.2 创建配置文件

创建 `custom/conf/app.ini` 文件，添加以下内容：

```ini
[database]
DB_TYPE = sqlite3
PATH = data/gitea.db

[server]
HTTP_PORT = 3000
DOMAIN = localhost
ROOT_URL = http://localhost:3000/

[plugin]
ENABLED = true
PLUGINS_DIR = ./plugins
MARKETPLACE_URL = https://plugins.gitea.io
ALLOW_MARKETPLACE_INSTALL = true
```

## 步骤 3：初始化数据库

```bash
# 运行数据库迁移
./gitea migrate
```

这将创建所有必要的数据库表，包括：
- `plugin` 表（插件信息）
- `authorized_device` 表（授权设备信息）

## 步骤 4：启动 Gitea

```bash
./gitea web
```

首次启动时，访问 `http://localhost:3000` 完成安装向导。

## 步骤 5：安装授权管理插件

### 方式 1：手动安装

```bash
# 创建插件目录
mkdir -p plugins/installed

# 编译授权管理插件
cd ../license-manager-plugin
./build.sh  # Linux/Mac
# 或
build.bat   # Windows

# 复制到插件目录
cp -r . ../gitea-license/plugins/installed/license-manager/
```

### 方式 2：通过后台安装（需要插件市场）

1. 以管理员身份登录 Gitea
2. 访问 `管理后台` -> `插件管理`
3. 点击 `插件市场`
4. 找到 `授权管理` 插件
5. 点击 `安装`

## 步骤 6：启用插件

1. 访问 `管理后台` -> `插件管理`
2. 找到 `授权管理` 插件
3. 点击 `启用`
4. 重启 Gitea（可选，某些插件可能需要）

## 步骤 7：验证安装

### 7.1 检查日志

查看 Gitea 日志，应该看到：

```
Plugin loaded: 授权管理 v1.0.0
Registering routes for plugin: license-manager
```

### 7.2 访问授权管理页面

1. 以普通用户身份登录
2. 访问 `个人设置` -> `授权管理`
3. 应该能看到授权设备列表页面

### 7.3 测试 API

```bash
# 创建授权
curl -X POST http://localhost:3000/api/v1/user/license/devices \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_code": "TEST123",
    "machine_name": "测试机器",
    "expiry_days": 365
  }'

# 验证授权
curl -X POST http://localhost:3000/api/v1/license/verify \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_code": "TEST123",
    "license_key": "YOUR_LICENSE_KEY"
  }'
```

## 故障排除

### 插件无法加载

**问题**：日志中没有看到 "Plugin loaded" 消息

**解决方案**：
1. 检查插件目录权限
2. 确认 `plugin.so` 文件存在
3. 查看详细错误日志
4. 确认 Go 版本兼容性

### 路由无法访问

**问题**：访问授权管理页面返回 404

**解决方案**：
1. 确认插件已启用
2. 重启 Gitea
3. 检查路由注册日志
4. 清除浏览器缓存

### 编译失败

**问题**：编译插件时出错

**解决方案**：
1. 确认 Go 版本 >= 1.21
2. 检查依赖包是否完整
3. 查看编译错误详情
4. 尝试清理并重新编译：
   ```bash
   go clean -cache
   go mod tidy
   ./build.sh
   ```

### 数据库迁移失败

**问题**：运行 `./gitea migrate` 时出错

**解决方案**：
1. 检查数据库连接配置
2. 确认数据库用户权限
3. 查看迁移错误日志
4. 手动创建表（不推荐）

## 生产环境部署建议

### 1. 使用系统服务

创建 systemd 服务文件 `/etc/systemd/system/gitea.service`：

```ini
[Unit]
Description=Gitea
After=network.target

[Service]
Type=simple
User=git
WorkingDirectory=/home/git/gitea
ExecStart=/home/git/gitea/gitea web
Restart=always
Environment=USER=git HOME=/home/git

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable gitea
sudo systemctl start gitea
```

### 2. 使用反向代理

Nginx 配置示例：

```nginx
server {
    listen 80;
    server_name git.example.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. 配置 HTTPS

使用 Let's Encrypt：

```bash
sudo certbot --nginx -d git.example.com
```

### 4. 数据备份

定期备份：
- 数据库
- `data/` 目录
- `custom/` 目录
- `plugins/` 目录

### 5. 监控和日志

- 配置日志轮转
- 设置监控告警
- 定期检查插件状态

## 更新和维护

### 更新 Gitea

```bash
# 停止服务
sudo systemctl stop gitea

# 备份
cp gitea gitea.backup
cp -r data data.backup

# 编译新版本
git pull
TAGS="bindata sqlite sqlite_unlock_notify" make build

# 运行迁移
./gitea migrate

# 启动服务
sudo systemctl start gitea
```

### 更新插件

```bash
# 停止服务
sudo systemctl stop gitea

# 更新插件
cd license-manager-plugin
git pull
./build.sh
cp -r . ../gitea-license/plugins/installed/license-manager/

# 启动服务
sudo systemctl start gitea
```

## 安全建议

1. **限制插件安装权限**：只允许管理员安装插件
2. **审查插件代码**：安装前检查插件源码
3. **使用 HTTPS**：保护数据传输安全
4. **定期更新**：及时更新 Gitea 和插件
5. **备份数据**：定期备份重要数据
6. **监控日志**：关注异常日志

## 获取帮助

- 查看日志：`logs/gitea.log`
- 查看文档：`INTEGRATION_GUIDE.md`
- 插件文档：`license-manager-plugin/README.md`
- 提交问题：GitHub Issues

## 总结

完成以上步骤后，你应该拥有一个完整的 Gitea 插件系统，包括：

✅ 插件管理后台
✅ 插件加载器
✅ 授权管理插件
✅ 用户级别的授权功能
✅ 完整的 API 接口

现在你可以：
- 管理已安装的插件
- 开发新的插件
- 使用授权管理功能
- 通过 API 集成到其他系统
