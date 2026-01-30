# Gitea 插件系统配置

## 配置文件

在 `custom/conf/app.ini` 中添加以下配置：

```ini
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; 插件设置
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

## 配置说明

### ENABLED
- 类型：布尔值
- 默认值：`true`
- 说明：是否启用插件系统

### PLUGINS_DIR
- 类型：字符串
- 默认值：`./plugins`
- 说明：插件安装目录，相对于 Gitea 工作目录

### MARKETPLACE_URL
- 类型：字符串
- 默认值：`https://plugins.gitea.io`
- 说明：插件市场的 URL 地址

### ALLOW_MARKETPLACE_INSTALL
- 类型：布尔值
- 默认值：`true`
- 说明：是否允许从插件市场安装插件

## 目录结构

```
plugins/
├── installed/          # 已安装的插件
│   ├── license-manager/
│   │   ├── plugin.so
│   │   ├── plugin.json
│   │   ├── models/
│   │   ├── services/
│   │   ├── routers/
│   │   ├── templates/
│   │   └── locales/
│   └── ...
├── temp/              # 临时文件
└── config/            # 插件配置
```

## 使用方法

1. 修改配置文件 `custom/conf/app.ini`
2. 重启 Gitea
3. 访问管理后台 -> 插件管理
4. 安装或管理插件
