# Gitea 插件系统快速开始

5 分钟快速体验 Gitea 插件系统和授权管理功能。

## 快速安装（开发环境）

### 1. 编译 Gitea（2 分钟）

```bash
cd gitea-license
TAGS="bindata sqlite sqlite_unlock_notify" make build
```

### 2. 初始化（1 分钟）

```bash
# 创建配置
mkdir -p custom/conf
cat > custom/conf/app.ini << EOF
[database]
DB_TYPE = sqlite3
PATH = data/gitea.db

[server]
HTTP_PORT = 3000

[plugin]
ENABLED = true
PLUGINS_DIR = ./plugins
EOF

# 初始化数据库
./gitea migrate
```

### 3. 安装插件（1 分钟）

```bash
# 创建插件目录
mkdir -p plugins/installed

# 编译并安装授权管理插件
cd ../license-manager-plugin
chmod +x build.sh
./build.sh
cp -r . ../gitea-license/plugins/installed/license-manager/
cd ../gitea-license
```

### 4. 启动（1 分钟）

```bash
./gitea web
```

访问 `http://localhost:3000`，完成安装向导（创建管理员账号）。

## 快速体验

### 1. 启用插件

1. 登录管理员账号
2. 访问 `管理后台` -> `插件管理`
3. 找到 `授权管理` 插件，点击 `启用`

### 2. 创建授权

1. 访问 `个人设置` -> `授权管理`
2. 点击 `新建授权`
3. 填写：
   - 机器码：`TEST-MACHINE-001`
   - 机器名称：`测试机器`
   - 有效期：`365` 天
4. 点击 `创建`
5. 复制生成的授权码

### 3. 验证授权（API）

```bash
# 获取 API Token
# 访问：个人设置 -> 应用 -> 生成新令牌

# 验证授权
curl -X POST http://localhost:3000/api/v1/license/verify \
  -H "Authorization: token YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_code": "TEST-MACHINE-001",
    "license_key": "YOUR_LICENSE_KEY"
  }'
```

## 客户端集成示例

### C# 客户端

```csharp
using System;
using System.Net.Http;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

class Program
{
    static async Task Main()
    {
        var client = new HttpClient();
        var request = new
        {
            machine_code = "TEST-MACHINE-001",
            license_key = "YOUR_LICENSE_KEY"
        };

        var json = JsonSerializer.Serialize(request);
        var content = new StringContent(json, Encoding.UTF8, "application/json");
        
        var response = await client.PostAsync(
            "http://localhost:3000/api/v1/license/verify",
            content
        );

        var result = await response.Content.ReadAsStringAsync();
        Console.WriteLine(result);
    }
}
```

### Python 客户端

```python
import requests

def verify_license(machine_code, license_key):
    url = "http://localhost:3000/api/v1/license/verify"
    data = {
        "machine_code": machine_code,
        "license_key": license_key
    }
    
    response = requests.post(url, json=data)
    return response.json()

# 使用
result = verify_license("TEST-MACHINE-001", "YOUR_LICENSE_KEY")
print(result)
```

### JavaScript 客户端

```javascript
async function verifyLicense(machineCode, licenseKey) {
    const response = await fetch('http://localhost:3000/api/v1/license/verify', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            machine_code: machineCode,
            license_key: licenseKey
        })
    });
    
    return await response.json();
}

// 使用
verifyLicense('TEST-MACHINE-001', 'YOUR_LICENSE_KEY')
    .then(result => console.log(result));
```

## 下一步

- 📖 阅读完整文档：`DEPLOYMENT_GUIDE.md`
- 🔧 开发新插件：`INTEGRATION_GUIDE.md`
- 💡 查看示例：`license-manager-plugin/`
- 🚀 生产部署：参考 `DEPLOYMENT_GUIDE.md` 的生产环境部署建议

## 常见问题

**Q: 插件没有加载？**
A: 检查日志，确认 `plugin.so` 文件存在且有执行权限。

**Q: 路由 404？**
A: 确认插件已启用，尝试重启 Gitea。

**Q: 编译失败？**
A: 确认 Go 版本 >= 1.21，运行 `go mod tidy`。

**Q: 如何卸载插件？**
A: 在插件管理页面点击 `卸载`，或手动删除插件目录。

## 获取帮助

遇到问题？查看：
- 日志文件：`logs/gitea.log`
- 详细文档：`INTEGRATION_GUIDE.md`
- 插件文档：`license-manager-plugin/README.md`
