# Package utils

utils 工具包

## FUNCTIONS

### GetExecutablePath

```go
func GetExecutablePath() string
```

GetExecutablePath 获取程序的绝对安装路径 如果无法通过 os.Executable 获取路径，则使用 os.Args[0] 作为替代

**返回：**

  * 程序的绝对路径字符串