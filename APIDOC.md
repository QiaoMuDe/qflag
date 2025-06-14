# qflag API文档

## API文档

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

### 标志注册表

Cmd结构体中的flagRegistry字段类型已更新为：

```go
map[interface{}]Flag
```

可以存储任何实现了Flag接口的标志类型实例。
### Cmd结构体
#### NewCmd

创建新的命令实例

```go
func NewCmd(name string, shortName string, errorHandling flag.ErrorHandling) *Cmd
```

参数:
- name: 命令名称
- shortName: 命令短名称
- errorHandling: 错误处理方式（flag.ContinueOnError、flag.ExitOnError、flag.PanicOnError）

#### String

添加字符串类型标志

```go
func (c *Cmd) String(name, shortName, defValue, help string) *StringFlag
```

#### Int

添加整数类型标志

```go
func (c *Cmd) Int(name, shortName string, defValue int, help string) *IntFlag
```

#### Bool

添加布尔类型标志

```go
func (c *Cmd) Bool(name, shortName string, defValue bool, help string) *BoolFlag
```

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

#### Float

添加浮点数类型标志

```go
func (c *Cmd) Float(name, shortName string, defValue float64, help string) *FloatFlag
```

#### StringVar

添加字符串类型标志变量

```go
func (c *Cmd) StringVar(p *string, name, shortName, defValue, help string)
```

#### IntVar

添加整数类型标志变量

```go
func (c *Cmd) IntVar(p *int, name, shortName string, defValue int, help string)
```

#### BoolVar

添加布尔类型标志变量

```go
func (c *Cmd) BoolVar(p *bool, name, shortName string, defValue bool, help string)
```

#### FloatVar

添加浮点数类型标志变量

```go
func (c *Cmd) FloatVar(p *float64, name, shortName string, defValue float64, help string)
```