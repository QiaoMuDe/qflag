# qflag API 文档

## Package qflag

`qflag` 包提供对标准库 `flag` 的封装，自动实现长短标志，并默认绑定 `-h/--help` 标志打印帮助信息。用户可通过 `Cmd.Help` 字段自定义帮助内容，支持直接赋值字符串或从文件加载。

## 变量

- **NewCmd**：导出 `cmd` 包中的 `NewCommand` 函数。
  ```go
  var NewCmd = cmd.NewCommand
  ```
- **QCommandLine**：导出 `cmd` 包的全局默认 `Command` 实例。
  ```go
  var QCommandLine = cmd.QCommandLine
  ```

## 函数

### AddExample

添加示例，这些示例将在命令行帮助信息中显示。

```go
func AddExample(e cmd.ExampleInfo)
```

- **参数**：
  - `e`：示例信息，`ExampleInfo` 类型。

### AddNote

添加注意事项，这些注意事项将在命令行帮助信息中显示。

```go
func AddNote(note string)
```

- **参数**：
  - `note`：注意事项内容，字符串类型。

### AddSubCmd

向全局默认命令实例 `QCommandLine` 添加一个或多个子命令。

```go
func AddSubCmd(subCmds ...*cmd.Cmd) error
```

- **参数**：
  - `subCmds`：可变参数，接收一个或多个 `*Cmd` 类型的子命令实例。
- **返回值**：
  - 若添加子命令过程中出现错误（如子命令为 `nil` 或存在循环引用），则返回错误信息；否则返回 `nil`。

### Arg

获取全局默认命令实例 `QCommandLine` 解析后的指定索引位置的非标志参数。

```go
func Arg(i int) string
```

- **参数**：
  - `i`：非标志参数的索引位置，从 0 开始计数。
- **返回值**：
  - 指定索引位置的非标志参数；若索引越界，则返回空字符串。

### Args

获取全局默认命令实例 `QCommandLine` 解析后的非标志参数切片。

```go
func Args() []string
```

- **返回值**：
  - 包含所有非标志参数的字符串切片。

### Bool

为全局默认命令创建一个布尔类型的命令行标志。

```go
func Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
- **返回值**：
  - 指向新创建的布尔标志对象的指针。

### BoolVar

将布尔类型的命令行标志绑定到全局默认命令实例 `QCommandLine` 中。

```go
func BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
```

- **参数**：
  - `f`：指向 `BoolFlag` 类型的指针，用于存储和管理布尔类型命令行标志的相关信息。
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。

### Duration

为全局默认命令定义一个时间间隔类型的命令行标志。

```go
func Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
```

- **参数**：
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。
- **返回值**：
  - 指向新创建的时间间隔类型标志对象的指针。

### DurationVar

为全局默认命令将一个时间间隔类型的命令行标志绑定到指定的 `DurationFlag` 指针。

```go
func DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
```

- **参数**：
  - `f`：指向 `DurationFlag` 类型的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### Enum

为全局默认命令定义一个枚举类型的命令行标志。

```go
func Enum(longName, shortName string, defValue string, usage string, enumValues []string) *flags.EnumFlag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
  - `enumValues`：枚举值的集合，用于指定标志可接受的取值范围。
- **返回值**：
  - 指向新创建的枚举类型标志对象的指针。

### EnumVar

为全局默认命令将一个枚举类型的命令行标志绑定到指定的 `EnumFlag` 指针。

```go
func EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string)
```

- **参数**：
  - `f`：指向 `EnumFlag` 类型的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。
  - `enumValues`：枚举值的集合，用于指定标志可接受的取值范围。

### FlagExists

检查全局默认命令实例 `QCommandLine` 中是否存在指定名称的标志。

```go
func FlagExists(name string) bool
```

- **参数**：
  - `name`：要检查的标志名称，可以是长名称或短名称。
- **返回值**：
  - 若存在指定名称的标志，则返回 `true`；否则返回 `false`。

### Float64

为全局默认命令创建一个浮点数类型的命令行标志。

```go
func Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
- **返回值**：
  - 指向新创建的浮点数标志对象的指针。

### Float64Var

为全局默认命令绑定一个浮点数类型的命令行标志到指定的 `FloatFlag` 指针。

```go
func Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
```

- **参数**：
  - `f`：指向 `FloatFlag` 的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### GetDescription

获取命令描述信息。

```go
func GetDescription() string
```

- **返回值**：
  - 命令描述信息。

### GetExamples

获取命令行标志的示例信息列表。

```go
func GetExamples() []cmd.ExampleInfo
```

- **返回值**：
  - 示例信息列表，每个元素为 `ExampleInfo` 类型。

### GetHelp

返回全局默认命令实例 `QCommandLine` 的帮助信息。

```go
func GetHelp() string
```

- **返回值**：
  - 命令行帮助信息。

### GetLogoText

获取全局默认命令实例 `QCommandLine` 的 logo 文本。

```go
func GetLogoText() string
```

- **返回值**：
  - 配置的 logo 文本。

### GetModuleHelps

获取模块帮助信息。

```go
func GetModuleHelps() string
```

- **返回值**：
  - 模块帮助信息。

### GetNotes

获取所有备注信息。

```go
func GetNotes() []string
```

- **返回值**：
  - 备注信息列表。

### GetUsageSyntax

获取全局默认命令实例 `QCommandLine` 的用法信息。

```go
func GetUsageSyntax() string
```

- **返回值**：
  - 命令行用法信息。

### GetUseChinese

获取当前命令行标志是否使用中文。

```go
func GetUseChinese() bool
```

- **返回值**：
  - 如果使用中文，则返回 `true`；否则返回 `false`。

### GetVersion

获取全局默认命令的版本信息。

```go
func GetVersion() string
```

- **返回值**：
  - 版本信息字符串。

### Int

为全局默认命令创建一个整数类型的命令行标志。

```go
func Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
- **返回值**：
  - 指向新创建的整数标志对象的指针。

### Int64

为全局默认命令定义一个 64 位整数类型的命令行标志。

```go
func Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
```

- **参数**：
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。
- **返回值**：
  - 指向新创建的 64 位整数类型标志对象的指针。

### Int64Var

为全局默认命令绑定一个 64 位整数类型的命令行标志到指定的 `Int64Flag` 指针。

```go
func Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
```

- **参数**：
  - `f`：指向 `Int64Flag` 对象的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### IntVar

将整数类型的命令行标志绑定到全局默认命令的 `IntFlag` 指针上。

```go
func IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
```

- **参数**：
  - `f`：指向 `IntFlag` 类型的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### LongName

获取命令长名称。

```go
func LongName() string
```

- **返回值**：
  - 命令长名称。

### Map

为全局默认命令创建一个键值对类型的命令行标志。

```go
func Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
- **返回值**：
  - 指向新创建的键值对标志对象的指针。

### MapVar

为全局默认命令将一个键值对类型的命令行标志绑定到指定的 `MapFlag` 指针。

```go
func MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
```

- **参数**：
  - `f`：指向 `MapFlag` 的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### NArg

获取全局默认命令实例 `QCommandLine` 解析后的非标志参数的数量。

```go
func NArg() int
```

- **返回值**：
  - 非标志参数的数量。

### NFlag

获取全局默认命令实例 `QCommandLine` 解析后已定义和使用的标志的数量。

```go
func NFlag() int
```

- **返回值**：
  - 标志的数量。

### Parse

完整解析命令行参数（含子命令处理）。

```go
func Parse() error
```

- **主要功能**：
  1. 解析当前命令的长短标志及内置标志。
  2. 自动检测并解析子命令及其参数（若存在）。
  3. 验证枚举类型标志的有效性。
- **参数**：
  - `args`：原始命令行参数切片（包含可能的子命令及参数）。
- **返回值**：
  - 解析过程中遇到的错误（如标志格式错误、子命令解析失败等）。
- **注意事项**：
  - 每个 `Cmd` 实例仅会被解析一次（线程安全）。
  - 若检测到子命令，会将剩余参数传递给子命令的 `Parse` 方法。
  - 处理内置标志执行逻辑。

### ParseFlagsOnly

仅解析当前命令的标志参数（忽略子命令）。

```go
func ParseFlagsOnly() error
```

- **主要功能**：
  1. 解析当前命令的长短标志及内置标志。
  2. 验证枚举类型标志的有效性。
  3. 明确忽略所有子命令及后续参数。
- **参数**：
  - `args`：原始命令行参数切片（子命令及后续参数会被忽略）。
- **返回值**：
  - 解析过程中遇到的错误（如标志格式错误等）。
- **注意事项**：
  - 每个 `Cmd` 实例仅会被解析一次（线程安全）。
  - 不会处理任何子命令，所有参数均视为当前命令的标志或位置参数。
  - 处理内置标志逻辑。

### Path

为全局默认命令创建一个路径类型的命令行标志。

```go
func Path(longName, shortName string, defValue string, usage string) *flags.PathFlag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
- **返回值**：
  - 指向新创建的路径标志对象的指针。

### PathVar

为全局默认命令将一个路径类型的命令行标志绑定到指定的 `PathFlag` 指针。

```go
func PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
```

- **参数**：
  - `f`：指向 `PathFlag` 的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### PrintHelp

输出全局默认命令实例 `QCommandLine` 的帮助信息。

```go
func PrintHelp()
```

### SetDescription

设置命令描述信息。

```go
func SetDescription(desc string)
```

- **参数**：
  - `desc`：命令描述信息。

### SetHelp

配置全局默认命令实例 `QCommandLine` 的帮助信息。

```go
func SetHelp(help string)
```

- **参数**：
  - `help`：新的帮助信息。

### LoadHelp

从文件中加载帮助信息。

```go
func LoadHelp(filePath string)
```

- **参数**：
  - `filePath`：帮助信息文件路径。

### SetLogoText

配置全局默认命令实例 `QCommandLine` 的 logo 文本。

```go
func SetLogoText(logoText string)
```

- **参数**：
  - `logoText`：配置的 logo 文本。

### SetModuleHelps

配置模块帮助信息。

```go
func SetModuleHelps(moduleHelps string)
```

- **参数**：
  - `moduleHelps`：模块帮助信息。

### SetUsageSyntax

配置全局默认命令实例 `QCommandLine` 的用法信息。

```go
func SetUsageSyntax(usageSyntax string)
```

- **参数**：
  - `usageSyntax`：新的用法信息。

### SetUseChinese

设置当前命令行标志是否使用中文。

```go
func SetUseChinese(useChinese bool)
```

- **参数**：
  - `useChinese`：如果使用中文，则传入 `true`；否则传入 `false`。

### SetVersion

为全局默认命令设置版本信息。

```go
func SetVersion(version string)
```

- **参数**：
  - `version`：版本信息字符串。

### ShortName

获取命令短名称。

```go
func ShortName() string
```

- **返回值**：
  - 命令短名称。

### Slice

为全局默认命令定义一个字符串切片类型的命令行标志。

```go
func Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
```

- **参数**：
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。
- **返回值**：
  - 指向新创建的字符串切片类型标志对象的指针。

### SliceVar

为全局默认命令将一个字符串切片类型的命令行标志绑定到指定的 `SliceFlag` 指针。

```go
func SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
```

- **参数**：
  - `f`：指向 `SliceFlag` 对象的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### String

为全局默认命令创建一个字符串类型的命令行标志。

```go
func String(longName, shortName, defValue, usage string) *flags.StringFlag
```

- **参数**：
  - `longName`：标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：标志的默认值，当命令行未指定该标志时使用。
  - `usage`：标志的帮助说明信息，用于在显示帮助信息时展示。
- **返回值**：
  - 指向新创建的字符串标志对象的指针。

### StringVar

将字符串类型的命令行标志绑定到全局默认命令的 `StringFlag` 指针上。

```go
func StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
```

- **参数**：
  - `f`：指向 `StringFlag` 的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### SubCmds

获取所有已注册的子命令列表。

```go
func SubCmds() []*cmd.Cmd
```

- **返回值**：
  - 子命令列表。

### Time

为全局默认命令定义一个时间类型的命令行标志。

```go
func Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
```

- **参数**：
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。
- **返回值**：
  - 指向新创建的时间类型标志对象的指针。

### TimeVar

为全局默认命令定义一个时间类型的命令行标志，并将其绑定到指定的 `TimeFlag` 指针。

```go
func TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
```

- **参数**：
  - `f`：指向 `TimeFlag` 对象的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

### Uint16

为全局默认命令定义一个无符号 16 位整数类型的命令行标志。

```go
func Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
```

- **参数**：
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。
- **返回值**：
  - 指向新创建的无符号 16 位整数类型标志对象的指针。

### Uint16Var

为全局默认命令创建一个无符号 16 位整数类型标志，并将其绑定到指定的 `Uint16Flag` 指针。

```go
func Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
```

- **参数**：
  - `f`：指向 `Uint16Flag` 对象的指针。
  - `longName`：命令行标志的长名称，在命令行中以 `--longName` 的形式使用。
  - `shortName`：命令行标志的短名称，在命令行中以 `-shortName` 的形式使用。
  - `defValue`：命令行标志的默认值。
  - `usage`：命令行标志的帮助说明信息。

## 类型

### ExampleInfo

导出 `cmd` 包中的 `ExampleInfo` 结构体。

```go
type ExampleInfo = cmd.ExampleInfo
```

### QCommandLineInterface

定义了全局默认命令行接口，提供统一的命令行参数管理功能。

```go
type QCommandLineInterface interface {
	// 元数据操作方法
	LongName() string
	ShortName() string
	GetDescription() string
	SetDescription(desc string)
	GetHelp() string
	SetHelp(help string)
	SetUsageSyntax(usageSyntax string)
	GetUsageSyntax() string
	GetUseChinese() bool
	SetUseChinese(useChinese bool)
	AddSubCmd(subCmd *cmd.Cmd)
	SubCmds() []*cmd.Cmd
	Parse() error
	ParseFlagsOnly() error
	Args() []string
	Arg(i int) string
	NArg() int
	NFlag() int
	PrintHelp()
	FlagExists(name string) bool
	AddNote(note string)
	GetNotes() []string
	AddExample(e cmd.ExampleInfo)
	GetExamples() []cmd.ExampleInfo
	SetVersion(version string)
	GetVersion() string
	SetLogoText(logoText string)
	GetLogoText() string
	SetModuleHelps(moduleHelps string)
	GetModuleHelps() string

	// 添加标志方法
	String(longName, shortName, defValue, usage string) *flags.StringFlag
	Int(longName, shortName string, defValue int, usage string) *flags.IntFlag
	Bool(longName, shortName string, defValue bool, usage string) *flags.BoolFlag
	Float64(longName, shortName string, defValue float64, usage string) *flags.Float64Flag
	Duration(longName, shortName string, defValue time.Duration, usage string) *flags.DurationFlag
	Enum(longName, shortName string, defValue string, usage string, enumValues []string) *flags.EnumFlag
	Slice(longName, shortName string, defValue []string, usage string) *flags.SliceFlag
	Int64(longName, shortName string, defValue int64, usage string) *flags.Int64Flag
	Uint16(longName, shortName string, defValue uint16, usage string) *flags.Uint16Flag
	Time(longName, shortName string, defValue time.Time, usage string) *flags.TimeFlag
	Map(longName, shortName string, defValue map[string]string, usage string) *flags.MapFlag
	Path(longName, shortName string, defValue string, usage string) *flags.PathFlag

	// 绑定变量方法
	StringVar(f *flags.StringFlag, longName, shortName, defValue, usage string)
	IntVar(f *flags.IntFlag, longName, shortName string, defValue int, usage string)
	BoolVar(f *flags.BoolFlag, longName, shortName string, defValue bool, usage string)
	Float64Var(f *flags.Float64Flag, longName, shortName string, defValue float64, usage string)
	DurationVar(f *flags.DurationFlag, longName, shortName string, defValue time.Duration, usage string)
	EnumVar(f *flags.EnumFlag, longName, shortName string, defValue string, usage string, enumValues []string)
	SliceVar(f *flags.SliceFlag, longName, shortName string, defValue []string, usage string)
	Int64Var(f *flags.Int64Flag, longName, shortName string, defValue int64, usage string)
	Uint16Var(f *flags.Uint16Flag, longName, shortName string, defValue uint16, usage string)
	TimeVar(f *flags.TimeFlag, longName, shortName string, defValue time.Time, usage string)
	MapVar(f *flags.MapFlag, longName, shortName string, defValue map[string]string, usage string)
	PathVar(f *flags.PathFlag, longName, shortName string, defValue string, usage string)
}
```