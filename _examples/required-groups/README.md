# 必需组示例

本示例演示了QFlag库中两种必需组的使用方法：
1. **普通必需组** - 组中的所有标志都必须被设置
2. **条件性必需组** - 如果组中任何一个标志被设置，则所有标志都必须被设置

## 文件说明

- `required-groups-demo.go` - 统一的示例程序，通过子命令展示不同的使用场景
- `simple.go` - 简单示例，演示基本用法（已合并到统一示例中）
- `main.go` - 完整示例，包含更多功能（已合并到统一示例中）
- `advanced.go` - 高级示例，模拟实际应用场景（已合并到统一示例中）
- `test.bat` - Windows批处理脚本，用于测试各种场景
- `test.sh` - Shell脚本，用于测试各种场景
- `test.ps1` - PowerShell脚本，用于测试各种场景

## 运行示例

### 统一示例程序

```bash
cd examples/required-groups
go run required-groups-demo.go --help
```

### 子命令

#### 简单示例

```bash
go run required-groups-demo.go simple --server-host localhost --server-port 8080
```

#### 完整示例

```bash
go run required-groups-demo.go full --host localhost --port 8080
```

#### 高级示例

```bash
go run required-groups-demo.go advanced --access-key AK123456 --secret-key SK789012 --region us-west-1 --service myapp
```

## 使用示例

### 1. 仅使用普通必需组

```bash
go run required-groups-demo.go simple --server-host localhost --server-port 8080
```

输出：
```
服务器: localhost:8080
数据库: 未配置
```

### 2. 使用普通必需组和条件性必需组

```bash
go run required-groups-demo.go simple --server-host localhost --server-port 8080 --db-host localhost --db-port 3306 --db-name mydb
```

输出：
```
服务器: localhost:8080
数据库: localhost:3306/mydb
```

### 3. 高级示例 - 完整配置

```bash
go run required-groups-demo.go advanced --access-key AK123456 --secret-key SK789012 --region us-west-1 --service myapp --db-host localhost --db-port 3306 --db-name mydb --db-user admin --db-pass password --verbose
```

输出：
```
=== 云服务部署配置 ===
访问密钥: AK****56
秘密密钥: SK****12
部署区域: us-west-1
服务名称: myapp

--- 数据库配置 ---
主机: localhost:3306
数据库: mydb
用户: admin
密码: pa****rd

--- Redis配置: 未启用 ---

--- 监控配置: 未启用 ---

--- 其他选项 ---
详细模式: 已启用

=== 部署准备完成 ===
```

### 4. 普通必需组未完全设置（会报错）

```bash
go run required-groups-demo.go simple --server-host localhost
```

输出：
```
参数解析错误: required flags [--server-port/-sp] in group 'server' must be set

使用示例:
  简单示例: required-groups-demo simple --server-host localhost --server-port 8080
  完整示例: required-groups-demo full --host localhost --port 8080
  高级示例: required-groups-demo advanced --access-key AK123 --secret-key SK456 --region us-west-1 --service myapp

使用 'required-groups-demo <subcommand> --help' 查看子命令的详细帮助
```

### 5. 条件性必需组部分设置（会报错）

```bash
go run required-groups-demo.go simple --server-host localhost --server-port 8080 --db-host localhost
```

输出：
```
参数解析错误: since one of the flags in group 'database' is used, all flags [--db-port/-dp --db-name/-dn] must be set

使用示例:
  简单示例: required-groups-demo simple --server-host localhost --server-port 8080
  完整示例: required-groups-demo full --host localhost --port 8080
  高级示例: required-groups-demo advanced --access-key AK123 --secret-key SK456 --region us-west-1 --service myapp

使用 'required-groups-demo <subcommand> --help' 查看子命令的详细帮助
```

## 代码说明

### 普通必需组

```go
// 添加普通必需组 - 所有标志都必须设置
err := cmd.AddRequiredGroup("server", []string{"server-host", "server-port"}, false)
```

- 第三个参数 `false` 表示这是一个普通必需组
- 无论是否使用了组中的任何一个标志，所有标志都必须被设置

### 条件性必需组

```go
// 添加条件性必需组 - 如果使用其中一个则必须同时使用
err = cmd.AddRequiredGroup("database", []string{"db-host", "db-port", "db-name"}, true)
```

- 第三个参数 `true` 表示这是一个条件性必需组
- 只有当组中任何一个标志被设置时，所有标志才必须被设置
- 如果组中的所有标志都未被设置，则不会触发验证

### 子命令

```go
// 创建子命令
simpleCmd := qflag.NewCmd("simple", "s", qflag.ContinueOnError)
simpleCmd.SetDesc("简单必需组示例")

// 设置执行函数
simpleCmd.SetRun(func(cmd types.Command) error {
    // 实现逻辑
    return nil
})

// 添加到全局根命令
qflag.AddSubCmds(simpleCmd)
```

## 使用场景

### 普通必需组适用于：
- 必须同时提供的核心功能参数
- 相互依赖的基础配置
- 认证相关的参数（如用户名和密码）

### 条件性必需组适用于：
- 可选功能模块的配置参数
- 可选连接的参数
- 可选输出格式的参数

## 注意事项

1. 普通必需组和条件性必需组可以同时使用
2. 条件性必需组提供了更灵活的标志验证方式
3. 错误消息会明确指出哪些标志需要设置
4. 可以使用 `IsSet()` 方法检查标志是否被设置
5. 在示例中，我们使用了独特的短标志名以避免冲突
6. 统一的示例程序通过子命令组织不同的使用场景，避免了多个main函数的问题