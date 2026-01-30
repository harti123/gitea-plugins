# Gitea 插件系统 + 授权管理

为 Gitea 开发的完整插件系统，支持动态加载、热插拔和后台管理。包含授权管理插件作为示例。

## 🎯 项目目标

1. 为 Gitea 添加插件系统，支持后台安装和管理
2. 实现授权管理插件，提供用户级别的设备授权功能
3. 提供完整的 API 接口，方便客户端集成

## ✨ 核心特性

### 插件系统
- ✅ 动态加载插件（无需重启）
- ✅ 后台管理界面
- ✅ 插件市场集成
- ✅ 版本管理
- ✅ 权限控制
- ✅ 配置持久化

### 授权管理插件
- ✅ 用户级别授权（数据隔离）
- ✅ 设备授权码生成
- ✅ 授权验证 API
- ✅ 过期时间管理
- ✅ 启用/禁用控制
- ✅ Web 界面管理

## 📦 项目结构

```
gitea-license/
├── models/
│   ├── plugin/              # 插件数据模型
│   ├── license/             # 授权数据模型
│   └── migrations/          # 数据库迁移
├── services/
│   ├── plugin/              # 插件服务
│   └── license/             # 授权服务
├── routers/
│   ├── web/admin/           # 管理后台
│   └── api/v1/license/      # API 接口
├── templates/
│   ├── admin/plugins/       # 插件管理页面
│   └── user/settings/       # 用户授权页面
└── modules/setting/         # 配置管理

license-manager-plugin/      # 授权管理插件
├── main.go                  # 插件入口
├── plugin.json              # 插件元数据
├── models/                  # 数据模型
├── services/                # 业务逻辑
├── templates/               # 页面模板
└── locales/                 # 国际化

jssevrer/tsServerManager/    # C# 客户端示例
├── LicenseVerifier.cs       # 授权验证器
├── LicenseDialog.xaml       # 授权对话框
└── LicenseDialog.xaml.cs    # 对话框逻辑
```

## 🚀 快速开始

### 1. 编译 Gitea

```bash
cd gitea-license
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### 2. 配置

创建 `custom/conf/app.ini`：

```ini
[database]
DB_TYPE = sqlite3
PATH = data/gitea.db

[server]
HTTP_PORT = 3000

[plugin]
ENABLED = true
PLUGINS_DIR = ./plugins
```

### 3. 初始化

```bash
./gitea migrate
```

### 4. 安装插件

```bash
cd ../license-manager-plugin
./build.sh
cp -r . ../gitea-license/plugins/installed/license-manager/
```

### 5. 启动

```bash
cd ../gitea-license
./gitea web
```

访问 `http://localhost:3000`

## 📖 文档

- **[快速开始](QUICK_START.md)** - 5 分钟快速体验
- **[部署指南](DEPLOYMENT_GUIDE.md)** - 完整的部署流程
- **[集成指南](INTEGRATION_GUIDE.md)** - 代码集成说明
- **[配置说明](PLUGIN_CONFIG.md)** - 配置参数详解
- **[数据库指南](PLUGIN_DATABASE_GUIDE.md)** - 插件数据库支持
- **[项目总结](PLUGIN_SYSTEM_SUMMARY.md)** - 完整的项目总结
- **[设计方案](Gitea插件系统设计方案.md)** - 系统设计文档
- **[使用指南](Gitea插件系统使用指南.md)** - 使用说明

### 插件文档
- **[插件 README](../license-manager-plugin/README.md)** - 插件功能说明
- **[插件安装](../license-manager-plugin/INSTALL.md)** - 插件安装指南

## 🔧 使用示例

### Web 界面

1. **管理员**：访问 `管理后台` -> `插件管理`
2. **用户**：访问 `个人设置` -> `授权管理`

### API 调用

```bash
# 创建授权
curl -X POST http://localhost:3000/api/v1/user/license/devices \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_code": "MACHINE-001",
    "machine_name": "我的电脑",
    "expiry_days": 365
  }'

# 验证授权
curl -X POST http://localhost:3000/api/v1/license/verify \
  -H "Content-Type: application/json" \
  -d '{
    "machine_code": "MACHINE-001",
    "license_key": "YOUR-LICENSE-KEY"
  }'
```

### C# 客户端

```csharp
var verifier = new LicenseVerifier("http://localhost:3000");
bool isValid = await verifier.VerifyLicenseAsync(machineCode, licenseKey);

if (!isValid)
{
    var dialog = new LicenseDialog();
    dialog.ShowDialog();
}
```

## 🏗️ 架构设计

### 插件系统架构

```
Gitea
  ├── Plugin Manager
  │   ├── Loader (动态加载)
  │   ├── Registry (插件注册)
  │   └── Lifecycle (生命周期)
  │
  └── Plugins
      ├── License Manager
      ├── Plugin 2
      └── Plugin N
```

### 授权管理流程

```
Client App
    ↓ (HTTP Request)
Gitea API
    ↓
License Plugin
    ↓
Database (authorized_device)
    ↓
Response (is_authorized: true/false)
```

## 🔐 安全特性

- ✅ 用户数据隔离
- ✅ API 身份验证
- ✅ 授权码加密（SHA256）
- ✅ 过期时间验证
- ✅ 权限控制

## 📊 完成度

| 模块 | 完成度 | 说明 |
|------|--------|------|
| 插件系统核心 | 100% | 完整实现 |
| 插件管理后台 | 100% | 完整实现 |
| 授权管理插件 | 100% | 完整实现 |
| 系统集成 | 100% | 完整集成 |
| 数据库迁移 | 100% | 完整实现 |
| API 接口 | 100% | 完整实现 |
| Web 界面 | 100% | 完整实现 |
| 文档 | 100% | 完整文档 |
| 客户端示例 | 100% | C# 示例 |

## 🎓 开发新插件

### 1. 创建插件结构

```
my-plugin/
├── plugin.json
├── main.go
├── models/
├── services/
└── templates/
```

### 2. 实现 IPlugin 接口

```go
type MyPlugin struct{}

func (p *MyPlugin) Info() *plugin.PluginInfo { ... }
func (p *MyPlugin) Init() error { ... }
func (p *MyPlugin) RegisterRoutes(r chi.Router) { ... }
func (p *MyPlugin) RegisterAPIRoutes(r chi.Router) { ... }
// ...
```

### 3. 编译和安装

```bash
go build -buildmode=plugin -o plugin.so
cp -r . /path/to/gitea/plugins/installed/my-plugin/
```

## 🐛 故障排除

### 插件无法加载
- 检查 Go 版本是否匹配
- 确认 `plugin.so` 文件存在
- 查看 Gitea 日志

### 路由 404
- 确认插件已启用
- 重启 Gitea
- 清除浏览器缓存

### 编译失败
- 确认 Go >= 1.21
- 运行 `go mod tidy`
- 查看编译错误详情

## 📝 更新日志

### v1.0.0 (2026-01-30)
- ✅ 完成插件系统核心功能
- ✅ 完成授权管理插件
- ✅ 完成系统集成
- ✅ 完成文档编写

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 👥 作者

Kiro Team

## 🔗 相关链接

- [Gitea 官网](https://gitea.io/)
- [Go 插件文档](https://pkg.go.dev/plugin)
- [项目文档](QUICK_START.md)

---

**现在就开始使用 Gitea 插件系统吧！** 🚀
