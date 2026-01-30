# 宝塔面板图形化部署 Gitea 指南

## 📌 像建网站一样部署 Gitea

本指南将教你如何在宝塔面板中像添加网站一样部署 Gitea。

---

## 第一步：登录宝塔面板

1. 打开浏览器
2. 访问：`http://你的服务器IP:8888`
3. 输入宝塔账号和密码
4. 点击"登录"

---

## 第二步：创建数据库（像建网站一样）

### 2.1 点击左侧菜单"数据库"

### 2.2 点击"添加数据库"按钮

### 2.3 填写数据库信息

在弹出的窗口中填写：

```
数据库名：gitea
用户名：gitea
密码：（点击"随机"按钮生成，或自己设置）
访问权限：本地服务器
备注：Gitea 数据库
```

### 2.4 点击"提交"

### 2.5 记录数据库信息

**重要**：复制并保存以下信息（后面会用到）：
- 数据库名：`gitea`
- 用户名：`gitea`
- 密码：`刚才生成的密码`

---

## 第三步：添加网站（为 Gitea 创建站点）

### 3.1 点击左侧菜单"网站"

### 3.2 点击"添加站点"按钮

### 3.3 填写网站信息

在弹出的窗口中填写：

```
域名：git.你的域名.com（或者直接用 IP）
备注：Gitea 代码托管平台
根目录：/www/wwwroot/gitea
FTP：不创建
数据库：不创建（已经创建过了）
PHP 版本：纯静态
```

### 3.4 点击"提交"

---

## 第四步：配置反向代理（让域名指向 Gitea）

### 4.1 找到刚创建的网站

在"网站"列表中找到 `git.你的域名.com`

### 4.2 点击"设置"按钮

### 4.3 点击左侧"反向代理"

### 4.4 点击"添加反向代理"

### 4.5 填写代理信息

```
代理名称：gitea
目标 URL：http://127.0.0.1:3000
发送域名：$host
```

### 4.6 点击"保存"

---

## 第五步：使用文件管理器上传代码

### 5.1 点击左侧菜单"文件"

### 5.2 进入目录

依次点击进入：`www` → `wwwroot` → `gitea`

### 5.3 使用终端克隆代码

在文件管理器右上角，点击"终端"按钮

在终端中输入以下命令：

```bash
cd /www/wwwroot/gitea
git clone https://gitee.com/harti/gitea-plugins.git
cd gitea-plugins
```

---

## 第六步：安装编译依赖（在终端中）

在刚才打开的终端中，继续输入：

```bash
# Ubuntu/Debian 系统
apt-get update
apt-get install -y git golang-go nodejs npm build-essential

# CentOS 系统（如果是 CentOS）
# yum install -y git golang nodejs npm gcc gcc-c++ make
```

等待安装完成（约 2-5 分钟）

---

## 第七步：编译 Gitea（在终端中）

在终端中继续输入：

```bash
# 设置 Go 代理（加速下载）
export GOPROXY=https://goproxy.cn,direct

# 开始编译（需要 10-20 分钟，请耐心等待）
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

**提示**：编译过程中会显示很多信息，这是正常的。等待出现 "Build succeeded" 或类似提示。

---

## 第八步：创建配置文件（使用文件管理器）

### 8.1 在文件管理器中操作

1. 确保你在 `/www/wwwroot/gitea/gitea-plugins/` 目录
2. 点击"新建文件夹"，创建文件夹：`custom`
3. 进入 `custom` 文件夹
4. 再创建文件夹：`conf`
5. 进入 `conf` 文件夹

### 8.2 创建配置文件

1. 点击"新建文件"
2. 文件名：`app.ini`
3. 点击"确定"

### 8.3 编辑配置文件

1. 点击 `app.ini` 文件
2. 点击"编辑"
3. 粘贴以下内容（**记得修改数据库密码和域名**）：

```ini
APP_NAME = Gitea: Git with a cup of tea
RUN_MODE = prod

[server]
PROTOCOL         = http
DOMAIN           = git.你的域名.com
ROOT_URL         = http://git.你的域名.com/
HTTP_PORT        = 3000
DISABLE_SSH      = false
SSH_PORT         = 22

[database]
DB_TYPE  = mysql
HOST     = 127.0.0.1:3306
NAME     = gitea
USER     = gitea
PASSWD   = 你的数据库密码（第二步记录的）
CHARSET  = utf8mb4

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

4. 点击"保存"

---

## 第九步：创建数据目录（使用文件管理器）

### 9.1 回到 `/www/wwwroot/gitea/gitea-plugins/` 目录

### 9.2 创建以下文件夹

依次点击"新建文件夹"，创建：
- `data`
- `log`

### 9.3 进入 `data` 文件夹，再创建

- `gitea-repositories`
- `plugins`

### 9.4 进入 `plugins` 文件夹，再创建

- `installed`

---

## 第十步：安装"进程守护管理器"（像安装插件一样）

### 10.1 点击左侧菜单"软件商店"

### 10.2 在搜索框输入

```
进程守护管理器
```

### 10.3 找到"进程守护管理器"

### 10.4 点击"安装"按钮

等待安装完成（约 1 分钟）

---

## 第十一步：添加 Gitea 进程（像添加网站一样）

### 11.1 点击左侧菜单"软件商店"

### 11.2 找到"进程守护管理器"

### 11.3 点击"设置"按钮

### 11.4 点击"添加守护进程"

### 11.5 填写进程信息

```
进程名称：gitea
启动用户：www
启动命令：/www/wwwroot/gitea/gitea-plugins/gitea web
运行目录：/www/wwwroot/gitea/gitea-plugins
进程数量：1
```

### 11.6 点击"提交"

### 11.7 点击"启动"按钮

---

## 第十二步：开放端口（像添加安全规则一样）

### 12.1 点击左侧菜单"安全"

### 12.2 在"端口列表"中点击"添加"

### 12.3 填写端口信息

```
端口：3000
协议：TCP
备注：Gitea Web 端口
```

### 12.4 点击"确定"

---

## 第十三步：配置 SSL 证书（可选，让网站支持 HTTPS）

### 13.1 点击左侧菜单"网站"

### 13.2 找到 `git.你的域名.com`

### 13.3 点击"设置"

### 13.4 点击左侧"SSL"

### 13.5 选择"Let's Encrypt"

### 13.6 点击"申请"

等待证书申请完成（约 1 分钟）

---

## 第十四步：访问 Gitea

### 14.1 打开浏览器

### 14.2 访问你的网站

```
http://git.你的域名.com
或
http://你的服务器IP:3000
```

### 14.3 完成初始化向导

首次访问会看到安装页面：

1. **数据库设置**（已自动填写）：
   - 数据库类型：MySQL
   - 主机：127.0.0.1:3306
   - 用户名：gitea
   - 密码：（已自动填写）
   - 数据库名：gitea

2. **一般设置**：
   - 站点名称：自定义（如：我的代码仓库）
   - 仓库根目录：保持默认
   - Git LFS 根目录：保持默认
   - 运行系统用户：www

3. **管理员账号设置**：
   - 管理员用户名：admin（或自定义）
   - 密码：设置一个强密码
   - 邮箱：你的邮箱

4. 点击"立即安装"

等待安装完成（约 10 秒）

---

## 🎉 完成！

现在你可以：

1. ✅ 使用管理员账号登录
2. ✅ 创建第一个代码仓库
3. ✅ 邀请团队成员
4. ✅ 开始使用 Git 推送代码

---

## 📊 日常管理（像管理网站一样）

### 查看 Gitea 运行状态

1. 点击"软件商店"
2. 找到"进程守护管理器"
3. 点击"设置"
4. 查看 gitea 进程状态

### 重启 Gitea

在进程守护管理器中：
1. 找到 gitea 进程
2. 点击"重启"按钮

### 查看日志

1. 点击"文件"
2. 进入 `/www/wwwroot/gitea/gitea-plugins/log/`
3. 点击 `gitea.log` 文件
4. 点击"查看"

### 备份数据

1. 点击"计划任务"
2. 点击"添加任务"
3. 任务类型：备份目录
4. 目录：`/www/wwwroot/gitea/gitea-plugins/data`
5. 执行周期：每天
6. 点击"添加"

---

## ❓ 常见问题

### 问题 1：编译失败

**解决**：在终端重新执行：
```bash
cd /www/wwwroot/gitea/gitea-plugins
export GOPROXY=https://goproxy.cn,direct
make clean
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### 问题 2：无法访问 3000 端口

**解决**：
1. 检查宝塔"安全"中是否开放了 3000 端口
2. 检查云服务器安全组是否开放了 3000 端口

### 问题 3：进程启动失败

**解决**：
1. 在终端执行：`chmod +x /www/wwwroot/gitea/gitea-plugins/gitea`
2. 重新启动进程

### 问题 4：数据库连接失败

**解决**：
1. 检查 `app.ini` 中的数据库密码是否正确
2. 在"数据库"中查看 gitea 数据库是否存在

---

## 📞 需要帮助？

如果遇到问题：
1. 查看日志文件：`/www/wwwroot/gitea/gitea-plugins/log/gitea.log`
2. 在宝塔论坛提问：https://www.bt.cn/bbs/
3. 查看 Gitea 文档：https://docs.gitea.io/

---

**恭喜你完成部署！** 🎉
