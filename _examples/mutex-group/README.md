# 互斥组示例

本示例演示了 qflag 库中互斥组 (mutex group) 功能的使用方法。

## 功能说明

互斥组是一组标志, 其中最多只能有一个被设置。当用户设置了互斥组中的多个标志时, 解析器会返回错误。

## 两种模式

1. **允许为空模式** (Allow=true): 可以不设置任何标志
2. **必须设置一个模式** (Allow=false): 必须设置其中一个标志

## 示例中的互斥组

### 输出格式互斥组
```go
app.AddMutexGroup("output_format", []string{"json", "xml", "yaml"}, true)
```
- 包含三个标志: `--json`、`--xml`、`--yaml`
- 允许为空模式: 可以不设置任何输出格式 (使用默认格式) 
- 最多只能设置其中一个输出格式

### 输入源互斥组
```go
app.AddMutexGroup("input_source", []string{"file", "url", "stdin"}, false)
```
- 包含三个标志: `--file`、`--url`、`--stdin`
- 必须设置一个模式: 必须指定其中一个输入源
- 最多只能设置其中一个输入源

## 运行示例

### 有效的用法

```bash
# 只指定输入源, 使用默认输出格式
go run . --file input.txt

# 指定输入源和输出格式
go run . --file input.txt --json
go run . --url http://example.com --xml
go run . --stdin --yaml
```

### 无效的用法

```bash
# 不指定输入源 (必须指定一个) 
go run .
# 错误: one of the mutually exclusive flags [file url stdin] in group 'input_source' must be set

# 指定多个输出格式
go run . --file input.txt --json --xml
# 错误: mutually exclusive flags [json xml] in group 'output_format' cannot be used together

# 指定多个输入源
go run . --file input.txt --url http://example.com
# 错误: mutually exclusive flags [file url] in group 'input_source' cannot be used together
```

## 代码结构

- `main.go`: 示例主程序, 演示互斥组的使用
- 添加了两个互斥组, 分别控制输出格式和输入源
- 解析参数后检查设置了哪些标志, 并显示相应的信息

## 适用场景

互斥组功能适用于以下场景: 

1. **输出格式选择**: JSON、XML、YAML 等格式互斥
2. **操作模式选择**: 创建、更新、删除等操作互斥
3. **输入源选择**: 文件、URL、标准输入等源互斥
4. **认证方式选择**: 用户名密码、API密钥、证书等认证方式互斥

互斥组功能可以有效防止用户使用冲突的参数, 提高命令行工具的健壮性和用户体验。