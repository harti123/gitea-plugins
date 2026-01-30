# 宝塔面板部署 Gitea 插件系统指南

## 📋 前提条件

- 已安装宝塔面板
- 服务器系统：Ubuntu/Debian/CentOS
- 至少 2GB 内存
- 至少 10GB 磁盘空间

## 🚀 部署步骤

### 步骤 1：通过宝塔面板安装依赖

#### 1.1 安装软件

登录宝塔面板，在"软件商店"中安装：

1. **Nginx** 或 **Apache**（用于反向代理）
2. **MySQL 8.0** 或 **PostgreSQL**（数据库）

#### 1.2 通过 SSH 终端安装编译依赖

在宝塔面板中：
1. 点击左侧"终端"
2. 或者点击"文件" → "终端"
3. 执行以下命令：

```bash
# Ubuntu/Debian 系统
apt-get update
apt-get install -y git golang-go nodejs npm build-essential sqlite3

# CentOS 系统
yum install -y git golang nodejs npm gcc gcc-c++ make sqlite
```

### 步骤 2：创建部署目录

在宝塔面板的"文件"管理中：

1. 进入 `/www/wwwroot/`
2. 创建新文件夹：`gitea`
3. 或者在终端执行：

```bash
mkdir -p /www/wwwroot/gitea
cd /www/wwwroot/gitea
```

### 步骤 3：克隆代码

在宝塔终端中执行：

```bash
cd /www/wwwroot/gitea
git clone https://gitee.com/harti/gitea-plugins.git
cd gitea-plugins
```

### 步骤 4：编译 Gitea

**重要**：编译需要 10-20 分钟，请耐心等待。

```bash
# 设置 Go 环境变量
export GOPROXY=https://goproxy.cn,direct

# 开始编译
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

编译完成后，会在当前目录生成 `gitea` 可执行文件。

### 步骤 5：创建数据库

#### 方式 1：使用 MySQL（推荐）

在宝塔面板中：
1. 点击"数据库"
2. 点击"添加数据库"
3. 填写信息：
   - 数据库名：`gitea`
   - 用户名：`gitea`
   - 密码：自动生成或自定义
4. 记录数据库信息

#### 方式 2：使用 SQLite（简单）

不需要创建数据库，Gitea 会自动创建 SQLite 文件。

### 步骤 6：创建配置文件

在宝塔面板的"文件"管理中：

1. 进入 `/www/wwwroot/gitea/gitea-plugins/`
2. 创建文件夹：`custom/conf`
3. 在 `custom/conf/` 中创建文件：`app.ini`
4. 编辑 `app.ini`，填入以下内容：

```ini
APP_NAME = Gitea: Git with a cup of tea
RUN_MODE = prod

[server]
PROTOCOL         = http
DOMAIN           = 你的域名或IP
ROOT_URL         = http://你的域名或IP:3000/
HTTP_PORT        = 3000
DISABLE_SSH      = false
SSH_PORT         = 22

[database]
# 使用 MySQL
DB_TYPE  = mysql
HOST     = 127.0.0.1:3306
NAME     = gitea
USER     = gitea
PASSWD   = 你的数据库密码

# 或使用 SQLite（注释掉上面的 MySQL 配置）
# DB_TYPE = sqlite3
# PATH    = /www/wwwroot/gitea/gitea-plugins/data/gitea.db

[repository]
ROOT = /www/wwwroot/gitea/gitea-plugins/data/gitea-repositories

[log]
MODE      = file
LEVEL     = info
ROOT_PATH = /www/wwwroot/gitea/gitea-plugins/log

[security]
INSTALL_LOCK   = false

[plugin]
ENABLED     = true
PLUGINS_DIR = /www/wwwroot/gitea/gitea-plugins/data/plugins
```

### 步骤 7：创建必要的目录

在宝塔终端执行：

```bash
cd /www/wwwroot/gitea/gitea-plugins
mkdir -p custom/conf
mkdir -p data/gitea-repositories
mkdir -p data/plugins/installed
mkdir -p log
```

### 步骤 8：使用宝塔的"进程守护管理器"

#### 8.1 安装进程守护管理器

1. 在宝塔面板"软件商店"中搜索"进程守护管理器"
2. 点击"安装"

#### 8.2 添加 Gitea 进程

1. 打开"进程守护管理器"
2. 点击"添加守护进程"
3. 填写信息：
   - **进程名称**：`gitea`
   - **启动命令**：`/www/wwwroot/gitea/gitea-plugins/gitea web`
   - **运行目录**：`/www/wwwroot/gitea/gitea-plugins`
   - **进程数量**：`1`
4. 点击"提交"
5. 点击"启动"

### 步骤 9：配置防火墙

在宝塔面板中：
1. 点击"安全"
2. 添加端口规则：
   - 端口：`3000`
   - 协议：`TCP`
   - 备注：`Gitea Web`
3. 如果需要 SSH 克隆，也添加：
   - 端口：`22`
   - 协议：`TCP`
   - 备注：`Gitea SSH`

### 步骤 10：配置反向代理（可选但推荐）

#### 10.1 在宝塔面板添加网站

1. 点击"网站"
2. 点击"添加站点"
3. 填写信息：
   - 域名：`git.你的域名.com`
   - 根目录：`/www/wwwroot/gitea`
   - PHP 版本：纯静态
4. 点击"提交"

#### 10.2 配置反向代理

1. 点击刚创建的网站的"设置"
2. 点击"反向代理"
3. 点击"添加反向代理"
4. 填写信息：
   - 代理名称：`gitea`
   - 目标 URL：`http://127.0.0.1:3000`
   - 发送域名：`$host`
5. 点击"提交"

#### 10.3 配置 SSL（推荐）

1. 在网站设置中点击"SSL"
2. 选择"Let's Encrypt"
3. 点击"申请"
4. 等待证书申请完成

### 步骤 11：访问 Gitea

1. **直接访问**：`http://你的服务器IP:3000`
2. **通过域名访问**（如果配置了反向代理）：`https://git.你的域名.com`

### 步骤 12：完成初始化

首次访问会进入安装向导：

1. **数据库设置**：
   - 如果使用 MySQL，填写之前创建的数据库信息
   - 如果使用 SQLite，保持默认

2. **一般设置**：
   - 站点名称：自定义
   - 仓库根目录：保持默认
   - Git LFS 根目录：保持默认

3. **管理员账号**：
   - 用户名：自定义
   - 密码：设置强密码
   - 邮箱：你的邮箱

4. 点击"立即安装"

## 🔧 管理命令

### 启动/停止/重启 Gitea

在宝塔的"进程守护管理器"中：
- 点击"启动"按钮
- 点击"停止"按钮
- 点击"重启"按钮

或者在终端执行：

```bash
# 查看进程
ps aux | grep gitea

# 停止进程
kill -9 进程ID

# 启动进程
cd /www/wwwroot/gitea/gitea-plugins
./gitea web &
```

### 查看日志

在宝塔面板的"文件"管理中：
- 进入 `/www/wwwroot/gitea/gitea-plugins/log/`
- 查看 `gitea.log` 文件

或在终端执行：

```bash
tail -f /www/wwwroot/gitea/gitea-plugins/log/gitea.log
```

### 备份数据

在宝塔面板中：
1. 点击"计划任务"
2. 添加任务：
   - 任务类型：备份目录
   - 目录：`/www/wwwroot/gitea/gitea-plugins/data`
   - 执行周期：每天
3. 点击"添加任务"

## 🐛 常见问题

### 问题 1：编译失败

**原因**：Go 版本太低或依赖下载失败

**解决**：
```bash
# 设置 Go 代理
export GOPROXY=https://goproxy.cn,direct

# 重新编译
make clean
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### 问题 2：无法访问 3000 端口

**原因**：防火墙未开放端口

**解决**：
1. 在宝塔面板"安全"中添加端口 3000
2. 检查云服务器安全组是否开放 3000 端口

### 问题 3：数据库连接失败

**原因**：数据库配置错误

**解决**：
1. 检查 `custom/conf/app.ini` 中的数据库配置
2. 确认数据库用户名和密码正确
3. 测试数据库连接：
   ```bash
   mysql -u gitea -p -h 127.0.0.1
   ```

### 问题 4：进程守护失败

**原因**：可执行文件权限不足

**解决**：
```bash
chmod +x /www/wwwroot/gitea/gitea-plugins/gitea
```

## 📊 性能优化

### 1. 使用 Redis 缓存

在宝塔面板安装 Redis，然后在 `app.ini` 中添加：

```ini
[cache]
ENABLED = true
ADAPTER = redis
HOST    = 127.0.0.1:6379
```

### 2. 配置 Nginx 缓存

在反向代理配置中添加缓存规则。

### 3. 定期清理日志

在宝塔"计划任务"中添加：
```bash
find /www/wwwroot/gitea/gitea-plugins/log -name "*.log" -mtime +30 -delete
```

## 🎉 完成！

现在你的 Gitea 插件系统已经在宝塔面板上运行了！

**下一步**：
1. ✅ 创建第一个仓库
2. ✅ 配置 SSH 密钥
3. ✅ 安装授权管理插件
4. ✅ 配置邮件服务

---

**需要帮助？**
- Gitee 仓库：https://gitee.com/harti/gitea-plugins
- 宝塔论坛：https://www.bt.cn/bbs/
