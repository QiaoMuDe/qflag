# qflag API文档

## 接口与结构体

### Flag接口

所有标志类型实现的通用接口

```go
type Flag interface {
    Name() string       // 获取标志名称
    ShortName() string  // 获取短标志名称
    Usage() string      // 获取使用说明
    Type() FlagType     // 获取标志类型
    getDefaultAny() any // 获取默认值(any类型)
}
```

### TypedFlag接口

泛型标志接口，继承自Flag接口

```go
type TypedFlag[T any] interface {
    Flag               // 继承Flag接口
    GetDefault() T     // 获取类型化默认值
    GetValue() T       // 获取类型化当前值
    SetValue(T)        // 设置类型化值
}
```

### Cmd结构体

命令结构体，用于定义和解析命令行标志

#### NewCmd

创建新的命令实例

```go
func NewCmd(name string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

参数:
- name: 命令名称
- shortName: 命令短名称
- errorHandling: 错误处理方式（flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError）

#### GetExecutablePath

获取程序的绝对安装路径

```go
func GetExecutablePath() string
```

返回值:
- 程序的绝对路径字符串

#### FlagExists

检查指定名称的标志是否存在

```go
func (c *Cmd) FlagExists(name string) bool
```

参数:
- name: 标志名称

返回值:
- bool: 标志是否存在

#### AddSubCmd

添加子命令

```go
func (c *Cmd) AddSubCmd(subCmds ...*Cmd)
```

#### Parse

解析命令行参数

```go
func (c *Cmd) Parse(args []string) error
```

## 标志添加方法

#### String

添加字符串类型标志

```go
func (c *Cmd) String(name, shortName, defValue, usage string) *StringFlag
```

参数:
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

#### Int

添加整数类型标志

```go
func (c *Cmd) Int(name, shortName string, defValue int, usage string) *IntFlag
```

参数:
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

#### Bool

添加布尔类型标志

```go
func (c *Cmd) Bool(name, shortName string, defValue bool, usage string) *BoolFlag
```

参数:
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

#### Float

添加浮点数类型标志

```go
func (c *Cmd) Float(name, shortName string, defValue float64, usage string) *FloatFlag
```

参数:
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

## 变量绑定方法

#### StringVar

添加字符串类型标志变量

```go
func (c *Cmd) StringVar(p *string, name, shortName, defValue, usage string)
```

参数:
- p: 存储标志值的指针
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

#### IntVar

添加整数类型标志变量

```go
func (c *Cmd) IntVar(p *int, name, shortName string, defValue int, usage string)
```

参数:
- p: 存储标志值的指针
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

#### BoolVar

添加布尔类型标志变量

```go
func (c *Cmd) BoolVar(p *bool, name, shortName string, defValue bool, usage string)
```

参数:
- p: 存储标志值的指针
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明

#### FloatVar

添加浮点数类型标志变量

```go
func (c *Cmd) FloatVar(p *float64, name, shortName string, defValue float64, usage string)
```

参数:
- p: 存储标志值的指针
- name: 长标志名
- shortName: 短标志名
- defValue: 默认值
- usage: 标志说明