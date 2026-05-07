# 标志依赖关系示例

本示例演示了 qflag 库中标志依赖关系 (flag dependency) 功能的使用方法。

## 功能说明

标志依赖关系允许定义标志之间的条件约束。当触发标志被设置时，目标标志必须满足特定的约束条件。

## 依赖类型

### 1. 必需依赖 (DepRequired)
当触发标志设置时，目标标志必须同时被设置。

```go
app.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert", "key"}, qflag.DepRequired)
```
- 当 `--ssl` 被设置时，`--cert` 和 `--key` 必须同时被设置
- 适用于：启用某个功能时，必须提供相关的配置项

### 2. 互斥依赖 (DepMutex)
当触发标志设置时，目标标志不能被设置。

```go
app.AddFlagDependency("debug_mutex_ssl", "debug", []string{"ssl"}, qflag.DepMutex)
```
- 当 `--debug` 被设置时，`--ssl` 不能被设置
- 适用于：某些功能不能同时使用

## 示例场景

本示例模拟一个 SSL 服务器配置工具：

- `--ssl`: 启用 SSL/TLS 加密
- `--cert`: SSL 证书文件路径
- `--key`: SSL 私钥文件路径
- `--ca-cert`: CA 证书文件路径（可选）
- `--port`: 服务器端口（默认 8080）
- `--debug`: 启用调试模式

### 依赖关系

1. **SSL 必需证书和私钥**: 启用 SSL 时，必须提供证书和私钥文件
2. **调试模式与 SSL 互斥**: 调试模式下不能使用 SSL（简化调试）

## 运行示例

### 有效的用法

```bash
# 普通HTTP服务器
go run .

# 指定端口
go run . --port 3000

# 启用调试模式
go run . --debug

# 启用SSL（必须提供证书和私钥）
go run . --ssl --cert server.crt --key server.key

# 启用SSL并指定端口
go run . --ssl --cert server.crt --key server.key --port 443
```

### 无效的用法

```bash
# 启用SSL但不提供证书（缺少必需依赖）
go run . --ssl
# 错误: 标志 'ssl' 设置了, 但以下标志未设置: [cert key]

# 启用SSL但只提供证书（缺少部分必需依赖）
go run . --ssl --cert server.crt
# 错误: 标志 'ssl' 设置了, 但以下标志未设置: [key]

# 同时启用调试模式和SSL（互斥冲突）
go run . --debug --ssl --cert server.crt --key server.key
# 错误: 标志 'debug' 设置了, 但以下标志不能设置: [ssl]
```

## 代码要点

### 添加标志依赖关系

```go
// 必需依赖
app.AddFlagDependency(
    "ssl_requires_cert",           // 依赖关系名称
    "ssl",                         // 触发标志
    []string{"cert", "key"},       // 目标标志列表
    qflag.DepRequired,             // 依赖类型：必需
)

// 互斥依赖
app.AddFlagDependency(
    "debug_mutex_ssl",             // 依赖关系名称
    "debug",                       // 触发标志
    []string{"ssl"},               // 目标标志列表
    qflag.DepMutex,                // 依赖类型：互斥
)
```

### 批量配置（通过 CmdOpts）

```go
opts := &qflag.CmdOpts{
    FlagDependencies: []qflag.FlagDependency{
        {
            Name:    "ssl_requires_cert",
            Trigger: "ssl",
            Targets: []string{"cert", "key"},
            Type:    qflag.DepRequired,
        },
        {
            Name:    "debug_mutex_ssl",
            Trigger: "debug",
            Targets: []string{"ssl"},
            Type:    qflag.DepMutex,
        },
    },
}
app.ApplyOpts(opts)
```

## 与互斥组和必需组的区别

| 特性 | 互斥组 | 必需组 | 标志依赖 |
|------|--------|--------|----------|
| 作用对象 | 一组标志之间 | 一组标志之间 | 触发标志 → 目标标志 |
| 关系方向 | 双向互斥 | 无方向 | 单向依赖 |
| 条件性 | 无条件 | 可选条件 | 基于触发标志 |
| 典型场景 | 多选一 | 全部必需 | 条件约束 |

标志依赖关系提供了更灵活的条件约束能力，适用于复杂的命令行参数验证场景。
