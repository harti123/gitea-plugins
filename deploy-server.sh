#!/bin/bash

# Gitea 插件系统服务器部署脚本
# 直接在服务器上编译和部署

set -e

echo "🚀 开始部署 Gitea 插件系统..."

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then 
    echo "❌ 请使用 root 用户运行此脚本"
    exit 1
fi

# 1. 安装依赖
echo "📦 安装依赖..."
apt-get update
apt-get install -y git golang-go nodejs npm build-essential sqlite3

# 2. 创建 git 用户
echo "👤 创建 git 用户..."
if ! id -u git > /dev/null 2>&1; then
    useradd -m -s /bin/bash git
fi

# 3. 克隆代码
echo "📥 克隆代码..."
cd /home/git
if [ ! -d "gitea-plugins" ]; then
    sudo -u git git clone https://gitee.com/harti/gitea-plugins.git
fi
cd gitea-plugins

# 4. 编译 Gitea
echo "🔨 编译 Gitea（这可能需要 10-20 分钟）..."
sudo -u git TAGS="bindata sqlite sqlite_unlock_notify" make build

# 5. 创建必要的目录
echo "📁 创建目录..."
mkdir -p /var/lib/gitea/{custom,data,log}
mkdir -p /var/lib/gitea/data/gitea/plugins/installed
mkdir -p /etc/gitea
chown -R git:git /var/lib/gitea
chown -R git:git /etc/gitea

# 6. 复制二进制文件
echo "📋 安装 Gitea..."
cp gitea /usr/local/bin/gitea
chmod +x /usr/local/bin/gitea

# 7. 创建 systemd 服务
echo "⚙️  创建系统服务..."
cat > /etc/systemd/system/gitea.service << 'EOF'
[Unit]
Description=Gitea (Git with a cup of tea)
After=syslog.target
After=network.target

[Service]
RestartSec=2s
Type=simple
User=git
Group=git
WorkingDirectory=/var/lib/gitea/
ExecStart=/usr/local/bin/gitea web --config /etc/gitea/app.ini
Restart=always
Environment=USER=git HOME=/home/git GITEA_WORK_DIR=/var/lib/gitea

[Install]
WantedBy=multi-user.target
EOF

# 8. 创建配置文件
echo "📝 创建配置文件..."
cat > /etc/gitea/app.ini << 'EOF'
APP_NAME = Gitea: Git with a cup of tea
RUN_MODE = prod
RUN_USER = git

[server]
PROTOCOL         = http
DOMAIN           = localhost
ROOT_URL         = http://localhost:3000/
HTTP_PORT        = 3000
DISABLE_SSH      = false
SSH_PORT         = 22
LFS_START_SERVER = true

[database]
DB_TYPE  = sqlite3
PATH     = /var/lib/gitea/data/gitea.db

[repository]
ROOT = /var/lib/gitea/data/gitea-repositories

[log]
MODE      = file
LEVEL     = info
ROOT_PATH = /var/lib/gitea/log

[security]
INSTALL_LOCK   = false
SECRET_KEY     = 
INTERNAL_TOKEN = 

[plugin]
ENABLED     = true
PLUGINS_DIR = /var/lib/gitea/data/gitea/plugins
EOF

chown git:git /etc/gitea/app.ini
chmod 640 /etc/gitea/app.ini

# 9. 启动服务
echo "🎯 启动 Gitea 服务..."
systemctl daemon-reload
systemctl enable gitea
systemctl start gitea

# 10. 检查状态
echo "✅ 检查服务状态..."
sleep 3
systemctl status gitea --no-pager

echo ""
echo "🎉 部署完成！"
echo ""
echo "📝 下一步："
echo "1. 访问 http://你的服务器IP:3000 完成初始化"
echo "2. 查看日志: journalctl -u gitea -f"
echo "3. 重启服务: systemctl restart gitea"
echo ""
