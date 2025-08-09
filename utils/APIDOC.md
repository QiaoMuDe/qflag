# Package utils

**导入路径:** `gitee.com/MM-Q/qflag/utils`

utils 工具包提供了qflag包使用的各种通用工具函数，包括字符串处理、类型转换、文件操作等辅助功能，为其他模块提供基础支持。

## 函数

### GetExecutablePath

```go
func GetExecutablePath() string
```

GetExecutablePath 获取程序的绝对安装路径。

**功能描述:**
如果无法通过 `os.Executable` 获取路径，则使用 `os.Args[0]` 作为替代方案。

**返回值:**
- `string`: 程序的绝对路径字符串

**使用场景:**
- 获取当前执行程序的完整路径
- 用于配置文件路径计算
- 日志文件路径确定
- 资源文件定位

**示例:**
```go
// 获取程序路径
execPath := GetExecutablePath()
fmt.Println("程序路径:", execPath)
// 输出: 程序路径: /usr/local/bin/myapp

// 基于程序路径构建配置文件路径
configPath := filepath.Join(filepath.Dir(execPath), "config.json")
```

**实现说明:**
1. 首先尝试使用 `os.Executable()` 获取程序路径
2. 如果失败，则回退到使用 `os.Args[0]`
3. 确保返回的是绝对路径