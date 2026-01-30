# Gitea 授权管理功能

## 功能概述

为 Gitea 添加了设备授权管理功能，允许每个用户管理自己的授权设备。

### 核心特性

✅ **用户级别功能** - 每个用户都可以使用，不是管理员专属  
✅ **数据隔离** - 用户只能看到和管理自己的授权设备  
✅ **授权码生成** - 自动生成唯一的授权码  
✅ **设备管理** - 添加、编辑、删除、启用/禁用设备  
✅ **到期管理** - 支持永久授权或设置有效期  
✅ **API 支持** - 提供 RESTful API 供客户端验证授权  

## 数据库表结构

```sql
CREATE TABLE `authorized_device` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT(20) NOT NULL COMMENT '所属用户ID',
  `device_id` VARCHAR(64) NOT NULL COMMENT '设备唯一ID',
  `machine_code` VARCHAR(64) NOT NULL COMMENT '机器码',
  `machine_name` VARCHAR(200) DEFAULT NULL COMMENT '机器名称',
  `license_key` VARCHAR(128) NOT NULL COMMENT '授权码',
  `is_enabled` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
  `expiry_date` BIGINT(20) DEFAULT NULL COMMENT '到期时间',
  `created_unix` BIGINT(20) NOT NULL COMMENT '创建时间',
  `updated_unix` BIGINT(20) NOT NULL COMMENT '更新时间',
  `last_verified_at` BIGINT(20) DEFAULT NULL COMMENT '最后验证时间',
  `remarks` TEXT COMMENT '备注',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_machine_code` (`machine_code`),
  KEY `idx_license_key` (`license_key`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='授权设备表';
```

## 使用流程

### 1. 用户端操作

#### 步骤 1：登录 Gitea
用户需要先登录 Gitea 账号。

#### 步骤 2：进入授权管理
访问：`用户设置` → `授权管理` (http://your-gitea.com/user/settings/license)

#### 步骤 3：添加设备授权
1. 点击"新建授权"按钮
2. 输入机器码（由客户端程序提供）
3. 填写机器名称（可选，便于识别）
4. 设置有效期（0 表示永久授权）
5. 添加备注（可选）
6. 点击"创建授权"

#### 步骤 4：获取授权码
创建成功后，系统会显示生成的授权码，格式如：
```
XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX
```

#### 步骤 5：将授权码提供给客户端
将授权码复制并输入到客户端程序中。

### 2. 客户端集成

#### 获取机器码
```csharp
using tsServerManager;

string machineCode = LicenseVerifier.GetMachineCode();
Console.WriteLine($"机器码: {machineCode}");
```

#### 验证授权
```csharp
var verifier = new LicenseVerifier("http://your-gitea.com");
var result = await verifier.VerifyLicenseAsync(machineCode, licenseKey);

if (result.IsAuthorized)
{
    Console.WriteLine("授权验证成功");
    if (result.ExpiryDate.HasValue)
    {
        Console.WriteLine($"到期时间: {result.ExpiryDate.Value}");
    }
}
else
{
    Console.WriteLine($"授权验证失败: {result.Message}");
}
```

## API 接口

### 1. 验证授权

**请求：**
```http
POST /api/v1/license/verify
Content-Type: application/json
Authorization: token YOUR_GITEA_TOKEN

{
  "machine_code": "ABCD1234EFGH5678IJKL9012MNOP3456",
  "license_key": "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
}
```

**响应：**
```json
{
  "is_authorized": true,
  "expiry_date": "2027-01-30T00:00:00Z",
  "message": "授权验证成功",
  "server_time": "2026-01-30T10:30:00Z"
}
```

### 2. 列出设备

**请求：**
```http
GET /api/v1/user/license/devices?page=1&limit=20
Authorization: token YOUR_GITEA_TOKEN
```

**响应：**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "device_id": "abc123def456",
    "machine_code": "ABCD1234EFGH5678IJKL9012MNOP3456",
    "machine_name": "办公室电脑",
    "license_key": "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX",
    "is_enabled": true,
    "expiry_date": 1738281600,
    "created_unix": 1738195200,
    "updated_unix": 1738195200,
    "last_verified_at": 1738281600,
    "remarks": "主要开发机器"
  }
]
```

### 3. 创建设备授权

**请求：**
```http
POST /api/v1/user/license/devices
Content-Type: application/json
Authorization: token YOUR_GITEA_TOKEN

{
  "machine_code": "ABCD1234EFGH5678IJKL9012MNOP3456",
  "machine_name": "办公室电脑",
  "expiry_days": 365,
  "remarks": "主要开发机器"
}
```

**响应：**
```json
{
  "success": true,
  "message": "授权创建成功",
  "device_id": "abc123def456",
  "license_key": "XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
}
```

### 4. 删除设备授权

**请求：**
```http
DELETE /api/v1/user/license/devices/1
Authorization: token YOUR_GITEA_TOKEN
```

**响应：**
```http
HTTP/1.1 204 No Content
```

### 5. 切换设备状态

**请求：**
```http
POST /api/v1/user/license/devices/toggle
Content-Type: application/json
Authorization: token YOUR_GITEA_TOKEN

{
  "id": 1
}
```

**响应：**
```json
{
  "success": true,
  "message": "设备状态已更新"
}
```

## 路由配置

需要在 Gitea 的路由配置中添加以下路由：

### Web 路由 (routers/web/web.go)
```go
// 用户设置 - 授权管理
m.Group("/settings", func() {
    m.Group("/license", func() {
        m.Get("", user.LicenseList)
        m.Get("/new", user.LicenseNew)
        m.Post("/new", user.LicenseNewPost)
        m.Get("/{id}/edit", user.LicenseEdit)
        m.Post("/{id}/edit", user.LicenseEditPost)
        m.Post("/{id}/delete", user.LicenseDelete)
        m.Post("/{id}/toggle", user.LicenseToggle)
    })
}, reqSignIn)
```

### API 路由 (routers/api/v1/api.go)
```go
// 授权管理 API
m.Group("/license", func() {
    m.Post("/verify", license.Verify)
    m.Post("/register", license.Register)
}, reqToken())

m.Group("/user/license", func() {
    m.Get("/devices", license.ListDevices)
    m.Post("/devices", license.CreateDevice)
    m.Delete("/devices/{id}", license.DeleteDevice)
    m.Post("/devices/toggle", license.ToggleDevice)
}, reqToken())
```

## 编译和部署

### 1. 编译 Gitea

```bash
cd gitea-license
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### 2. 运行数据库迁移

```bash
./gitea migrate
```

### 3. 启动 Gitea

```bash
./gitea web
```

## 安全注意事项

1. **API Token 保护**：所有 API 请求都需要有效的 Gitea Token
2. **数据隔离**：用户只能访问自己的授权设备
3. **HTTPS 传输**：生产环境建议使用 HTTPS 保护授权码传输
4. **授权码保密**：授权码应妥善保管，不要泄露给他人

## 常见问题

### Q1: 如何获取机器码？
A: 使用客户端程序的 `LicenseVerifier.GetMachineCode()` 方法获取。

### Q2: 授权码丢失了怎么办？
A: 在 Gitea 的授权管理页面可以查看已创建的授权码。

### Q3: 可以为同一台机器创建多个授权吗？
A: 不可以，每个机器码在同一用户下只能创建一个授权。

### Q4: 授权到期后会怎样？
A: 客户端验证时会返回"授权已过期"，需要在 Gitea 中延长有效期。

### Q5: 如何临时禁用某个设备？
A: 在授权管理页面点击"禁用"按钮，不需要删除授权。

## 许可证

本功能遵循 Gitea 的 MIT 许可证。
