# qflag API文档

## API文档
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