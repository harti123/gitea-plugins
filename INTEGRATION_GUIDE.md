# Gitea 插件系统集成指南

本文档说明如何将插件系统集成到 Gitea 源码中。

## ✅ 集成状态

- ✅ 插件系统初始化代码（cmd/web.go）
- ✅ 插件管理路由（routers/web/web.go）
- ✅ 动态插件路由注册（routers/web/web.go）
- ✅ 插件配置支持（modules/setting/plugin.go）
- ✅ 数据库迁移（models/migrations/v1_26/v326.go, v327.go）
- ✅ 管理员菜单项（templates/admin/navbar.tmpl）
- ✅ 国际化支持（options/locale/locale_zh-CN_plugins.ini）

## 已完成的集成工作

所有必要的代码已经集成完成！现在可以直接编译和使用。

## 1. ✅ 初始化插件系统

已在 `cmd/web.go` 中添加初始化代码：

```go
package cmd

import (
    // ... 其他导入
    plugin_service "code.gitea.io/gitea/services/plugin"
)

func runWeb(ctx *cli.Context) error {
    // ... 现有代码

    // 初始化插件系统
    pluginsDir := setting.PluginsDir
    if pluginsDir == "" {
        pluginsDir = "./plugins"
    }
    
    marketplaceURL := setting.PluginMarketplaceURL
    if marketplaceURL == "" {
        marketplaceURL = "https://plugins.gitea.io"
    }

    if err := plugin_service.InitManager(pluginsDir, marketplaceURL); err != nil {
        log.Fatal("Failed to initialize plugin manager: %v", err)
    }

    // 加载所有已安装的插件
    loader := plugin_service.GetLoader()
    if err := loader.LoadAll(graceful.GetManager().HammerContext()); err != nil {
        log.Error("Failed to load plugins: %v", err)
    }

    // ... 继续现有代码
}
```

## 2. ✅ 注册插件管理路由

已在 `routers/web/web.go` 中添加路由：

```go
package web

import (
    // ... 其他导入
    "code.gitea.io/gitea/routers/web/admin"
    plugin_service "code.gitea.io/gitea/services/plugin"
)

func Routes(ctx gocontext.Context) *web.Route {
    // ... 现有代码

    // 管理员路由
    m.Group("/admin", func() {
        // ... 现有管理员路由

        // 插件管理路由
        m.Group("/plugins", func() {
            m.Get("", admin.PluginsList)
            m.Get("/market", admin.PluginsMarket)
            m.Post("/install", admin.PluginInstall)
            m.Post("/{id}/uninstall", admin.PluginUninstall)
            m.Post("/{id}/toggle", admin.PluginToggle)
        })
    }, reqAdmin)

    // 动态注册插件路由
    registerPluginRoutes(m)

    return m
}

// registerPluginRoutes 注册所有插件的路由
func registerPluginRoutes(m *web.Route) {
    loader := plugin_service.GetLoader()
    if loader == nil {
        return
    }

    plugins := loader.GetAllPlugins()
    for pluginID, plugin := range plugins {
        log.Info("Registering routes for plugin: %s", pluginID)
        
        // 注册 Web 路由
        plugin.RegisterRoutes(m)
        
        // 注册 API 路由
        plugin.RegisterAPIRoutes(m)
    }
}
```

## 3. ✅ 添加配置支持

已创建 `modules/setting/plugin.go` 并在 `setting.go` 中加载：

```go
package setting

var (
    // Plugin settings
    PluginsDir           string
    PluginMarketplaceURL string
    PluginEnabled        bool
    AllowMarketInstall   bool
)

func loadPluginFrom(rootCfg ConfigProvider) {
    sec := rootCfg.Section("plugin")
    PluginsDir = sec.Key("PLUGINS_DIR").MustString("./plugins")
    PluginMarketplaceURL = sec.Key("MARKETPLACE_URL").MustString("https://plugins.gitea.io")
    PluginEnabled = sec.Key("ENABLED").MustBool(true)
    AllowMarketInstall = sec.Key("ALLOW_MARKETPLACE_INSTALL").MustBool(true)
}

func LoadSettings() {
    // ... 现有代码
    
    loadPluginFrom(CfgProvider)
    
    // ... 继续现有代码
}
```

### 在 `custom/conf/app.ini.sample` 中添加配置示例

```ini
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; Plugin Settings
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

[plugin]
;; 是否启用插件系统
ENABLED = true
;; 插件目录
PLUGINS_DIR = ./plugins
;; 插件市场地址
MARKETPLACE_URL = https://plugins.gitea.io
;; 是否允许从市场安装插件
ALLOW_MARKETPLACE_INSTALL = true
```

## 4. ✅ 数据库迁移

已创建迁移文件：
- `models/migrations/v1_26/v326.go` - 插件表
- `models/migrations/v1_26/v327.go` - 授权设备表

已在 `models/migrations/migrations.go` 中注册：

```go
package migrations

import (
    "xorm.io/xorm"
)

func addPluginTable(x *xorm.Engine) error {
    type Plugin struct {
        ID          int64  `xorm:"pk autoincr"`
        PluginID    string `xorm:"VARCHAR(100) UNIQUE NOT NULL INDEX"`
        Name        string `xorm:"VARCHAR(200) NOT NULL"`
        Version     string `xorm:"VARCHAR(50) NOT NULL"`
        Description string `xorm:"TEXT"`
        Author      string `xorm:"VARCHAR(200)"`
        Homepage    string `xorm:"VARCHAR(500)"`
        License     string `xorm:"VARCHAR(50)"`
        IsEnabled   bool   `xorm:"NOT NULL DEFAULT false"`
        IsInstalled bool   `xorm:"NOT NULL DEFAULT false"`
        InstallPath string `xorm:"VARCHAR(500)"`
        Config      string `xorm:"TEXT"`
        CreatedUnix int64  `xorm:"created"`
        UpdatedUnix int64  `xorm:"updated"`
    }

    return x.Sync2(new(Plugin))
}
```

### 在 `models/migrations/migrations.go` 中注册迁移

```go
var migrations = []Migration{
    // ... 现有迁移
    
    // v1.XX.0 -> v1.XX.1
    NewMigration("Add plugin table", addPluginTable),
}
```

## 5. ✅ 添加管理员菜单项

已在 `templates/admin/navbar.tmpl` 中添加菜单：

```html
<div class="item{{if .PageIsAdminPlugins}} active{{end}}">
    <a href="{{AppSubUrl}}/admin/plugins">
        {{svg "octicon-plug"}} {{ctx.Locale.Tr "admin.plugins.title"}}
    </a>
</div>
```

## 6. 编译和测试

### 编译 Gitea

```bash
cd gitea-license
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### 运行数据库迁移

```bash
./gitea migrate
```

### 启动 Gitea

```bash
./gitea web
```

### 测试插件系统

1. 访问 `http://localhost:3000/admin/plugins`
2. 应该能看到插件管理页面
3. 尝试安装插件

## 7. 安装授权管理插件

### 方式 1：手动安装

```bash
# 解压插件到 plugins 目录
unzip license-manager-1.0.0.zip -d ./plugins/installed/

# 重启 Gitea
./gitea web
```

### 方式 2：通过后台安装

1. 访问插件市场
2. 找到授权管理插件
3. 点击安装

## 8. 验证插件功能

### 检查插件是否加载

查看 Gitea 日志，应该看到：

```
Plugin loaded: 授权管理 v1.0.0
Registering routes for plugin: license-manager
```

### 访问授权管理页面

```
http://localhost:3000/user/settings/license
```

### 测试 API

```bash
curl -X POST http://localhost:3000/api/v1/license/verify \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_code": "TEST123",
    "license_key": "TEST456"
  }'
```

## 9. 故障排除

### 插件无法加载

1. 检查 Go 编译环境
2. 查看 Gitea 日志
3. 确认插件目录权限
4. 检查 plugin.json 格式

### 路由无法访问

1. 确认插件已启用
2. 检查路由注册代码
3. 清除浏览器缓存

### 编译失败

1. 确认 Go 版本 >= 1.19
2. 检查依赖包
3. 查看编译错误日志

## 10. 开发新插件

参考 `license-manager-plugin/` 的结构：

1. 创建 `plugin.json`
2. 实现 `IPlugin` 接口
3. 编写业务逻辑
4. 创建模板和资源
5. 打包测试

## 总结

完成以上步骤后，Gitea 插件系统就完全集成好了。你可以：

- ✅ 在后台管理插件
- ✅ 从市场安装插件
- ✅ 开发自定义插件
- ✅ 动态加载/卸载插件

---

**注意**：这是一个二次开发方案，需要修改 Gitea 源码并重新编译。
