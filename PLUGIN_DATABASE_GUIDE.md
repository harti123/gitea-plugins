# Gitea 插件数据库支持指南

## 概述

Gitea 插件系统支持插件自动创建和管理数据库表。当插件安装或加载时，系统会自动同步插件定义的数据库模型到数据库中。

## 工作原理

### 1. 插件定义数据模型

插件在 `RegisterModels()` 方法中返回需要创建的数据库模型：

```go
func (p *MyPlugin) RegisterModels() []interface{} {
    return []interface{}{
        &MyModel{},
        &AnotherModel{},
    }
}
```

### 2. 自动创建表

当插件加载时，系统会：
1. 调用插件的 `RegisterModels()` 方法获取模型列表
2. 使用 XORM 的 `Sync2()` 方法自动创建或更新表结构
3. 如果表已存在，会自动添加缺失的列（不会删除现有列）

### 3. 数据持久化

- 插件卸载时，**不会自动删除数据表**，以防止数据丢失
- 如需清理数据，需要在插件的 `Uninstall()` 方法中手动处理

## 使用示例

### 完整示例：授权管理插件

#### 1. 定义数据模型

```go
// models/device.go
package models

import (
    "code.gitea.io/gitea/modules/timeutil"
)

type AuthorizedDevice struct {
    ID          int64              `xorm:"pk autoincr"`
    DeviceID    string             `xorm:"VARCHAR(100) UNIQUE NOT NULL INDEX"`
    UserID      int64              `xorm:"INDEX NOT NULL"`
    MachineCode string             `xorm:"VARCHAR(200) NOT NULL INDEX"`
    MachineName string             `xorm:"VARCHAR(200)"`
    LicenseKey  string             `xorm:"VARCHAR(500) NOT NULL"`
    IsEnabled   bool               `xorm:"NOT NULL DEFAULT true"`
    ExpiryDate  timeutil.TimeStamp `xorm:"INDEX"`
    Remarks     string             `xorm:"TEXT"`
    CreatedUnix timeutil.TimeStamp `xorm:"created"`
    UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
}
```

#### 2. 在插件中注册模型

```go
// main.go
package main

import (
    "code.gitea.io/gitea/models/db"
    license_model "license-manager-plugin/models"
)

type LicenseManagerPlugin struct{}

func (p *LicenseManagerPlugin) Init() error {
    // 注册模型到 db 包（可选，用于类型识别）
    db.RegisterModel(new(license_model.AuthorizedDevice))
    return nil
}

func (p *LicenseManagerPlugin) RegisterModels() []interface{} {
    return []interface{}{
        &license_model.AuthorizedDevice{},
    }
}
```

#### 3. 在 plugin.json 中声明

```json
{
    "id": "license-manager",
    "name": "授权管理",
    "version": "1.0.0",
    "hooks": {
        "models": true
    }
}
```

## 数据库操作

### 创建记录

```go
func CreateDevice(ctx context.Context, device *AuthorizedDevice) error {
    _, err := db.GetEngine(ctx).Insert(device)
    return err
}
```

### 查询记录

```go
func GetDeviceByID(ctx context.Context, id int64) (*AuthorizedDevice, error) {
    device := &AuthorizedDevice{}
    has, err := db.GetEngine(ctx).ID(id).Get(device)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, ErrDeviceNotExist
    }
    return device, nil
}
```

### 更新记录

```go
func UpdateDevice(ctx context.Context, device *AuthorizedDevice) error {
    _, err := db.GetEngine(ctx).ID(device.ID).AllCols().Update(device)
    return err
}
```

### 删除记录

```go
func DeleteDevice(ctx context.Context, id int64) error {
    _, err := db.GetEngine(ctx).ID(id).Delete(&AuthorizedDevice{})
    return err
}
```

### 复杂查询

```go
func SearchDevices(ctx context.Context, opts *SearchOptions) ([]*AuthorizedDevice, int64, error) {
    sess := db.GetEngine(ctx).
        Where("user_id = ?", opts.UserID)
    
    if opts.IsEnabled != nil {
        sess = sess.And("is_enabled = ?", *opts.IsEnabled)
    }
    
    count, err := sess.Count(&AuthorizedDevice{})
    if err != nil {
        return nil, 0, err
    }
    
    devices := make([]*AuthorizedDevice, 0, opts.PageSize)
    err = sess.
        Limit(opts.PageSize, (opts.Page-1)*opts.PageSize).
        Desc("created_unix").
        Find(&devices)
    
    return devices, count, err
}
```

## 数据库标签说明

### XORM 标签

```go
type MyModel struct {
    ID          int64  `xorm:"pk autoincr"`           // 主键，自增
    Name        string `xorm:"VARCHAR(200) NOT NULL"` // 字符串，非空
    Email       string `xorm:"VARCHAR(200) UNIQUE"`   // 唯一索引
    UserID      int64  `xorm:"INDEX"`                 // 普通索引
    IsActive    bool   `xorm:"NOT NULL DEFAULT true"` // 布尔值，默认 true
    Content     string `xorm:"TEXT"`                  // 长文本
    CreatedUnix int64  `xorm:"created"`               // 创建时间（自动）
    UpdatedUnix int64  `xorm:"updated"`               // 更新时间（自动）
}
```

### 常用标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `pk` | 主键 | `xorm:"pk autoincr"` |
| `autoincr` | 自增 | `xorm:"pk autoincr"` |
| `NOT NULL` | 非空 | `xorm:"NOT NULL"` |
| `UNIQUE` | 唯一 | `xorm:"UNIQUE"` |
| `INDEX` | 索引 | `xorm:"INDEX"` |
| `DEFAULT` | 默认值 | `xorm:"DEFAULT true"` |
| `created` | 创建时间 | `xorm:"created"` |
| `updated` | 更新时间 | `xorm:"updated"` |
| `TEXT` | 长文本 | `xorm:"TEXT"` |
| `VARCHAR(n)` | 字符串长度 | `xorm:"VARCHAR(200)"` |

## 数据迁移

### 添加新字段

如果需要在已有表中添加新字段：

1. 在模型中添加新字段
2. 重新加载插件
3. XORM 会自动添加新列

```go
type AuthorizedDevice struct {
    // ... 现有字段
    
    // 新增字段
    LastVerifiedUnix timeutil.TimeStamp `xorm:"INDEX"`
}
```

### 修改字段类型

**注意**：XORM 的 `Sync2()` 不会自动修改字段类型，需要手动迁移：

```go
func (p *MyPlugin) Init() error {
    // 手动执行 SQL 迁移
    engine := db.GetEngine(context.Background())
    
    // 检查是否需要迁移
    _, err := engine.Exec("ALTER TABLE authorized_device MODIFY COLUMN machine_code VARCHAR(300)")
    if err != nil {
        log.Warn("Migration failed: %v", err)
    }
    
    return nil
}
```

## 数据清理

### 卸载时清理数据

如果需要在卸载插件时删除数据：

```go
func (p *MyPlugin) Uninstall() error {
    // 警告：这会删除所有数据！
    engine := db.GetEngine(context.Background())
    
    // 删除表
    if err := engine.DropTables(&AuthorizedDevice{}); err != nil {
        return fmt.Errorf("drop tables: %w", err)
    }
    
    log.Info("Plugin data cleaned up")
    return nil
}
```

### 保留数据（推荐）

默认情况下，卸载插件时不删除数据：

```go
func (p *MyPlugin) Uninstall() error {
    // 不删除数据表，只做清理工作
    log.Info("Plugin uninstalled, data preserved")
    return nil
}
```

## 多数据库支持

插件的数据库模型会自动适配 Gitea 配置的数据库类型：

- SQLite
- MySQL
- PostgreSQL
- MSSQL

XORM 会自动处理不同数据库的语法差异。

## 事务支持

### 使用事务

```go
func CreateDeviceWithTransaction(ctx context.Context, device *AuthorizedDevice) error {
    return db.WithTx(ctx, func(ctx context.Context) error {
        // 在事务中执行多个操作
        if _, err := db.GetEngine(ctx).Insert(device); err != nil {
            return err
        }
        
        // 其他操作...
        
        return nil
    })
}
```

## 性能优化

### 1. 添加索引

```go
type AuthorizedDevice struct {
    UserID      int64  `xorm:"INDEX"`           // 单列索引
    MachineCode string `xorm:"INDEX"`           // 单列索引
    IsEnabled   bool   `xorm:"INDEX"`           // 单列索引
}
```

### 2. 复合索引

```go
// 在 Init() 中创建复合索引
func (p *MyPlugin) Init() error {
    engine := db.GetEngine(context.Background())
    
    // 创建复合索引
    _, err := engine.Exec("CREATE INDEX idx_user_enabled ON authorized_device(user_id, is_enabled)")
    if err != nil {
        log.Warn("Create index failed: %v", err)
    }
    
    return nil
}
```

### 3. 查询优化

```go
// 只查询需要的字段
func GetDeviceNames(ctx context.Context, userID int64) ([]string, error) {
    var names []string
    err := db.GetEngine(ctx).
        Table("authorized_device").
        Where("user_id = ?", userID).
        Cols("machine_name").
        Find(&names)
    return names, err
}
```

## 最佳实践

### 1. 使用时间戳

```go
type MyModel struct {
    CreatedUnix timeutil.TimeStamp `xorm:"created"`
    UpdatedUnix timeutil.TimeStamp `xorm:"updated"`
}
```

### 2. 软删除

```go
type MyModel struct {
    DeletedUnix timeutil.TimeStamp `xorm:"INDEX"`
}

func DeleteDevice(ctx context.Context, id int64) error {
    device := &AuthorizedDevice{
        ID:          id,
        DeletedUnix: timeutil.TimeStampNow(),
    }
    _, err := db.GetEngine(ctx).ID(id).Cols("deleted_unix").Update(device)
    return err
}
```

### 3. 数据验证

```go
func (d *AuthorizedDevice) Validate() error {
    if d.MachineCode == "" {
        return errors.New("machine code is required")
    }
    if d.UserID <= 0 {
        return errors.New("invalid user id")
    }
    return nil
}
```

### 4. 错误处理

```go
var (
    ErrDeviceNotExist     = errors.New("device does not exist")
    ErrDeviceAlreadyExist = errors.New("device already exists")
)

func GetDevice(ctx context.Context, id int64) (*AuthorizedDevice, error) {
    device := &AuthorizedDevice{}
    has, err := db.GetEngine(ctx).ID(id).Get(device)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, ErrDeviceNotExist
    }
    return device, nil
}
```

## 调试技巧

### 1. 启用 SQL 日志

```go
func (p *MyPlugin) Init() error {
    // 开发环境启用 SQL 日志
    if !setting.IsProd {
        db.GetEngine(context.Background()).ShowSQL(true)
    }
    return nil
}
```

### 2. 检查表结构

```go
func (p *MyPlugin) Init() error {
    engine := db.GetEngine(context.Background())
    
    // 获取表信息
    tables, err := engine.DBMetas()
    if err != nil {
        return err
    }
    
    for _, table := range tables {
        log.Info("Table: %s", table.Name)
        for _, col := range table.Columns() {
            log.Info("  Column: %s, Type: %s", col.Name, col.SQLType.Name)
        }
    }
    
    return nil
}
```

## 总结

Gitea 插件系统提供了完整的数据库支持：

✅ 自动创建表结构
✅ 自动添加新字段
✅ 支持多种数据库
✅ 事务支持
✅ 灵活的数据清理策略
✅ 完整的 CRUD 操作

通过合理使用这些功能，插件可以轻松管理自己的数据，无需手动编写数据库迁移脚本。
