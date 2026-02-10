# qflag 开发进度计划

## 概述

本文档基于 `REFACTOR_PLAN.md` 的架构设计, 制定开发进度计划。开发分为 **4个阶段**, 优先级从高到低: 

- **P0 (必须) **: 核心类型定义、基础命令、基础标志、基础解析器
- **P1 (重要) **: 验证器、配置和错误处理
- **P2 (有用) **: 帮助信息生成、公共 API 优化
- **P3 (可选) **: 命令行补全、复杂标志类型

---

## 阶段一: 核心类型定义 (P0) 

### 目标
建立整个项目的基础类型和接口定义, 确保后续实现有坚实的根基。

### 任务列表

#### 1.1 创建项目结构
- [ ] 初始化 Go module (`gitee.com/MM-Q/qflag1`) 
- [ ] 创建目录结构
  ```
  qflag/
  ├── internal/
  │   ├── types/
  │   ├── cmd/
  │   ├── flag/
  │   ├── parser/
  │   ├── validator/
  │   ├── help/
  │   ├── registry/
  │   ├── completion/
  │   └── error/
  ├── exports.go
  └── qflag.go
  ```

#### 1.2 实现 types 包
- [ ] **cmd.go** - 定义 Cmd 接口
  ```go
  type Cmd interface {
      // 基本属性
      Name() string
      LongName() string
      ShortName() string
      Desc() string

      // 标志管理
      AddFlag(flag Flag) error
      GetFlag(name string) (Flag, bool)
      Flags() []Flag

      // 子命令管理
      AddSubCmd(cmd Cmd) error
      GetSubCmd(name string) (Cmd, bool)
      SubCmds() []Cmd
      HasSubCmd(name string) bool

      // 参数解析
      Parse(args []string) error
    ParseAndRoute(args []string) error
    ParseOnly(args []string) ([]string, error)  // 解析当前命令的参数, 不递归
    IsParsed() bool

      // 参数访问
      Args() []string
      Arg(index int) string
      NArg() int

      // 执行
      Run() error
      SetRun(fn func(Cmd) error)
      HasRunFunc() bool

      // 帮助信息
      Help() string
      PrintHelp()

      // 配置
      SetDesc(desc string)
      SetVersion(version string)
      SetChinese(useChinese bool)
      SetEnvPrefix(prefix string)
      Config() *CmdConfig
  }
  ```

- [ ] **flag.go** - 定义 Flag 接口和 FlagType 枚举
  ```go
  type Flag interface {
      Name() string
      LongName() string
      ShortName() string
      Desc() string
      Type() FlagType
      Get() any
      Set(value string) error
      GetDefault() any
      IsSet() bool
      Reset()
      String() string
      Validate() error
      SetValidator(validator Validator)
      BindEnv(name string)
      GetEnvVar() string
  }

  type FlagType int

  const (
      FlagTypeUnknown FlagType = iota
      // 标准库已实现的类型 (P0/P1) 
      FlagTypeString
      FlagTypeInt
      FlagTypeInt64
      FlagTypeUint
      FlagTypeUint64
      FlagTypeFloat64
      FlagTypeBool
      FlagTypeDuration
      // 需要自定义实现的类型 (P2/P3) 
      FlagTypeTime
      FlagTypeMap
      FlagTypeStringSlice
      FlagTypeIntSlice
      FlagTypeEnum
      FlagTypeSize
  )
  ```

- [ ] **validator.go** - 定义 Validator 接口
  ```go
  type Validator interface {
      Validate(value any) error
  }
  ```

- [ ] **registry.go** - 定义注册表接口
  ```go
  type FlagRegistry interface {
      Register(flag Flag) error
      Unregister(name string) error
      Get(name string) (Flag, bool)
      List() []Flag
      Has(name string) bool
      Count() int
      Clear()
  }

  type CmdRegistry interface {
      Register(cmd Cmd) error
      Unregister(name string) error
      Get(name string) (Cmd, bool)
      List() []Cmd
      Has(name string) bool
      Count() int
      Clear()
  }
  ```

- [ ] **config.go** - 定义配置类型
  ```go
  type CmdConfig struct {
      Version    string
      UseChinese bool
  }

  type FlagConfig struct {
      Required   bool
      EnvPrefix  string
  }
  ```

- [ ] **error.go** - 定义错误类型
  ```go
  type Error struct {
      Code    string
      Message string
      Cause   error
  }
  ```

### 验收标准
- [ ] 所有接口和类型定义完整
- [ ] 代码可以编译 (即使没有实现) 
- [ ] 编写接口单元测试

---

## 阶段二: 注册表和基础命令 (P0) 

### 目标
实现注册表管理命令和标志的存储, 基础命令结构提供命令的基本功能。

### 任务列表

#### 2.1 实现 registry 包
- [ ] **impl.go** - 通用注册表实现 (被 FlagRegistry 和 CmdRegistry 共用) 
  ```go
  type registry[T any] struct {
      mu   sync.RWMutex
      items map[string]T
  }

  func NewRegistry[T any]() *registry[T] {
      return &registry[T]{
          items: make(map[string]T),
      }
  }
  ```

- [ ] **flag_registry.go** - FlagRegistry 实现
  ```go
  type FlagRegistryImpl struct {
      *registry[Flag]
  }

  func NewFlagRegistry() FlagRegistry {
      return &FlagRegistryImpl{
          registry: NewRegistry[Flag](),
      }
  }
  ```

- [ ] **cmd_registry.go** - CmdRegistry 实现
  ```go
  type CmdRegistryImpl struct {
      *registry[Cmd]
  }

  func NewCmdRegistry() CmdRegistry {
      return &CmdRegistryImpl{
          registry: NewRegistry[Cmd](),
      }
  }
  ```

#### 2.2 实现 cmd 包
- [ ] **base_cmd.go** - Cmd 基础实现
  ```go
  type Cmd struct {
      mu sync.RWMutex

      longName         string
      shortName        string
      description      string
      config           *config.CmdConfig
      flagRegistry     types.FlagRegistry
      cmdRegistry  types.CmdRegistry
      args             []string
      parsed           bool
      parseOnce        sync.Once
      runFunc          func(Cmd) error
      parser           parser.Parser
      helpGen          help.HelpGenerator
  }

  func NewCmd(longName, shortName string) *Cmd {
      return &Cmd{
          longName:        longName,
          shortName:       shortName,
          config:          config.NewCmdConfig(),
          flagRegistry:    registry.NewFlagRegistry(),
          cmdRegistry: registry.NewCmdRegistry(),
          args:            []string{},
          parser:          parser.NewDefaultParser(),
          helpGen:         help.NewTextGenerator(),
      }
  }
  ```

- [ ] 实现 Cmd 接口的所有方法
  - [ ] 基本属性方法 (Name、LongName、ShortName、Description) 
  - [ ] 标志管理方法 (AddFlag、GetFlag、Flags) 
  - [ ] 子命令管理方法 (AddSubCmd、GetSubCmd、SubCmds、HasSubCmd) 
  - [ ] 参数解析方法 (Parse、ParseAndRoute、ParseOnly) 
  - [ ] 参数访问方法 (Args、Arg、NArg) 
  - [ ] 执行方法 (Run、SetRun、HasRunFunc) 
  - [ ] 帮助信息方法 (Help、PrintHelp) 
  - [ ] 配置方法 (SetDesc、SetVersion、SetChinese) 

- [ ] **root_cmd.go** - RootCmd 根命令 (可选, 继承 Cmd) 
- [ ] **sub_cmd.go** - SubCmd 子命令 (可选, 继承 Cmd) 

### 验收标准
- [ ] 注册表可以正确存储和查询命令/标志
- [ ] 基础命令结构可以正确管理标志和子命令
- [ ] 编写集成测试

---

## 阶段三: 基础标志类型 (P0) 

### 目标
实现基础标志类型, **优先使用标准库已实现的标志类型**, 对于标准库未实现的类型使用 `flag.Value` 接口实现。

### 标志实现策略

| 优先级 | 类型 | 实现方式 | 说明 |
|--------|------|----------|------|
| P0 | `String`, `Int`, `Bool`, `Float64`, `Duration` | 直接使用标准库方法 | `flag.StringVar`、`flag.IntVar` 等 |
| P1 | `Int64`, `Uint`, `Uint64` | 直接使用标准库方法 | `flag.Int64Var`、`flag.UintVar` 等 |
| P2 | `Time`, `StringSlice` | 自定义 `flag.Value` 接口 | 实现 `String()` 和 `Set(string)` |
| P3 | `Map`, `IntSlice`, `Enum`, `Size` | 自定义 `flag.Value` 接口 | 实现 `String()` 和 `Set(string)` |

### 任务列表

#### 3.1 实现 base_flag.go
- [ ] **BaseFlag[T any]** 泛型结构体
  ```go
  type BaseFlag[T any] struct {
      mu          sync.RWMutex
      longName    string
      shortName   string
      description string
      flagType    FlagType
      value       *T
      default_    T
      isSet       bool
      validator   validator.Validator
      envVar      string
  }
  ```

- [ ] 实现 Flag 接口的所有方法
  - [ ] 基本属性方法 (Name、LongName、ShortName、Description、Type) 
  - [ ] 值访问方法 (Get、Set、GetDefault、IsSet、Reset) 
  - [ ] 字符串表示 (String) 
  - [ ] 验证方法 (Validate、SetValidator) 
  - [ ] 环境变量方法 (BindEnv、GetEnvVar) 

- [ ] 添加泛型方法
  ```go
  func (f *BaseFlag[T]) GetValue() T
  func (f *BaseFlag[T]) SetValue(value T) error
  func (f *BaseFlag[T]) GetDefaultTyped() T
  func (f *BaseFlag[T]) GetValuePtr() *T  // 用于标准库注册
  ```

- [ ] 添加 Value 接口实现 (用于自定义类型) 
  ```go
  func (f *BaseFlag[T]) String() string {
      return fmt.Sprintf("%v", f.Get())
  }

  func (f *BaseFlag[T]) Set(value string) error {
      // 具体类型需要实现此方法
      return fmt.Errorf("not implemented")
  }
  ```

#### 3.2 实现标准库已实现的标志类型 (P0) 
- [ ] **string_flag.go** - 使用 `flag.StringVar`
  ```go
  type StringFlag struct {
      *BaseFlag[string]
  }

  func NewStringFlag(longName, shortName, description string, default_ string) *StringFlag {
      return &StringFlag{
          BaseFlag: NewBaseFlag(string)("string", longName, shortName, description, FlagTypeString, default_),
      }
  }
  ```

- [ ] **int_flag.go** - 使用 `flag.IntVar`
  ```go
  type IntFlag struct {
      *BaseFlag[int]
  }

  func NewIntFlag(longName, shortName, description string, default_ int) *IntFlag {
      return &IntFlag{
          BaseFlag: NewBaseFlag[int]("int", longName, shortName, description, FlagTypeInt, default_),
      }
  }
  ```

- [ ] **bool_flag.go** - 使用 `flag.BoolVar`
  ```go
  type BoolFlag struct {
      *BaseFlag[bool]
  }

  func NewBoolFlag(longName, shortName, description string, default_ bool) *BoolFlag {
      return &BoolFlag{
          BaseFlag: NewBaseFlag[bool]("bool", longName, shortName, description, FlagTypeBool, default_),
      }
  }
  ```

- [ ] **float64_flag.go** - 使用 `flag.Float64Var`
  ```go
  type Float64Flag struct {
      *BaseFlag[float64]
  }

  func NewFloat64Flag(longName, shortName, description string, default_ float64) *Float64Flag {
      return &Float64Flag{
          BaseFlag: NewBaseFlag[float64]("float64", longName, shortName, description, FlagTypeFloat64, default_),
      }
  }
  ```

- [ ] **duration_flag.go** - 使用 `flag.DurationVar`
  ```go
  type DurationFlag struct {
      *BaseFlag[time.Duration]
  }

  func NewDurationFlag(longName, shortName, description string, default_ time.Duration) *DurationFlag {
      return &DurationFlag{
          BaseFlag: NewBaseFlag[time.Duration]("duration", longName, shortName, description, FlagTypeDuration, default_),
      }
  }
  ```

#### 3.3 实现标准库已实现的扩展类型 (P1) 
- [ ] **int64_flag.go** - 使用 `flag.Int64Var`
- [ ] **uint_flag.go** - 使用 `flag.UintVar`
- [ ] **uint64_flag.go** - 使用 `flag.Uint64Var`

#### 3.4 实现需要自定义的类型 (P2/P3) 
- [ ] **time_flag.go** - 自定义 `flag.Value` 接口实现
  ```go
  type TimeFlag struct {
      *BaseFlag[time.Time]
  }

  func NewTimeFlag(longName, shortName, description string, default_ time.Time) *TimeFlag {
      return &TimeFlag{
          BaseFlag: NewBaseFlag[time.Time]("time", longName, shortName, description, FlagTypeTime, default_),
      }
  }

  // 实现 flag.Value 接口
  func (f *TimeFlag) Set(value string) error {
      t, err := time.Parse("2006-01-02 15:04:05", value)
      if err != nil {
          return err
      }
      *f.value = t
      f.isSet = true
      return nil
  }
  ```

- [ ] **string_slice_flag.go** - 自定义 `flag.Value` 接口实现
  ```go
  type StringSliceFlag struct {
      *BaseFlag[[]string]
  }

  func NewStringSliceFlag(longName, shortName, description string, default_ []string) *StringSliceFlag {
      return &StringSliceFlag{
          BaseFlag: NewBaseFlag[[]string]("stringSlice", longName, shortName, description, FlagTypeStringSlice, default_),
      }
  }

  // 实现 flag.Value 接口
  func (f *StringSliceFlag) Set(value string) error {
      *f.value = append(*f.value, value)
      f.isSet = true
      return nil
  }
  ```

- [ ] **map_flag.go** - 自定义 `flag.Value` 接口实现 (P3) 
- [ ] **int_slice_flag.go** - 自定义 `flag.Value` 接口实现 (P3) 
- [ ] **enum_flag.go** - 自定义 `flag.Value` 接口实现 (P3) 
- [ ] **size_flag.go** - 自定义 `flag.Value` 接口实现 (P3) 

### 验收标准
- [ ] P0 类型 (String、Int、Bool、Float64、Duration) 使用标准库方法实现
- [ ] P1 类型 (Int64、Uint、Uint64) 使用标准库方法实现
- [ ] P2/P3 类型使用 `flag.Value` 接口实现
- [ ] 长名字和短名字都可以使用
- [ ] 默认值和 IsSet 状态正确
- [ ] 编写单元测试

---

## 阶段四: 基础解析器 (P0) 

### 目标
实现参数解析器, 支持使用标准库 `flag` 包解析标志, 递归解析子命令。

### 任务列表

#### 4.1 实现 parser.go
- [x] **DefaultParser** 结构体
  ```go
  type DefaultParser struct {
      flagSet *flag.FlagSet
  }
  ```

- [x] **NewDefaultParser** 工厂函数
  ```go
  func NewDefaultParser() *DefaultParser
  ```

- [x] **ParseOnly** 方法 - 解析当前命令的参数, 不递归
  ```go
  func (p *DefaultParser) ParseOnly(cmd types.Cmd, args []string) ([]string, error)
  ```

- [x] **Parse** 方法 - 解析并递归子命令
  ```go
  func (p *DefaultParser) Parse(cmd types.Cmd, args []string) error
  ```

- [x] **ParseAndRoute** 方法 - 解析并执行
  ```go
  func (p *DefaultParser) ParseAndRoute(cmd types.Cmd, args []string) error
  ```

- [x] **registerFlag** 方法 - 将标志注册到标准库
  ```go
  func (p *DefaultParser) registerFlag(f types.Flag)
  ```

### 验收标准
- [ ] 可以解析长名字标志 (`--flag`) 
- [ ] 可以解析短名字标志 (`-f`) 
- [ ] 可以递归解析子命令
- [ ] ParseAndRoute 可以自动路由并执行
- [ ] 编写集成测试

---

## 阶段五: 基础验证器 (P1) 

### 目标
实现基础的验证器, 支持必填检查、范围检查等。

### 任务列表

#### 5.1 实现 validator 包
- [ ] **required_validator.go** - 必填验证
  ```go
  type RequiredValidator struct{}

  func (v *RequiredValidator) Validate(value any) error
  ```

- [ ] **range_validator.go** - 范围验证
  ```go
  type RangeValidator struct {
      Min any
      Max any
  }

  func (v *RangeValidator) Validate(value any) error
  ```

- [ ] **length_validator.go** - 长度验证 (适用于字符串、切片) 
  ```go
  type LengthValidator struct {
      Min int
      Max int
  }

  func (v *LengthValidator) Validate(value any) error
  ```

- [ ] **enum_validator.go** - 枚举验证
  ```go
  type EnumValidator struct {
      Allowed []any
  }

  func (v *EnumValidator) Validate(value any) error
  ```

### 验收标准
- [ ] 验证器可以正确验证标志值
- [ ] 验证失败时返回有意义的错误信息
- [ ] 编写单元测试

---

## 阶段六: 帮助信息生成 (P2) 

### 目标
实现帮助信息生成器, 支持文本格式的中英文帮助。

### 任务列表

#### 6.1 实现 help 包
- [ ] **text_generator.go** - 文本帮助生成器
  ```go
  type TextGenerator struct{}

  func (g *TextGenerator) Generate(cmd types.Cmd) string
  ```

- [ ] **template.go** - 帮助模板 (可选) 

- [ ] 支持中英文切换
- [ ] 显示标志的默认值、描述等信息

### 验收标准
- [ ] 可以生成格式良好的帮助信息
- [ ] 支持中文和英文
- [ ] 显示所有标志和子命令的信息

---

## 阶段七: 公共 API 和导出 (P2) 

### 目标
提供便捷的公共 API, 简化用户使用。

### 任务列表

#### 7.1 实现 exports.go
- [ ] 导出 types 包中的所有类型
  ```go
  type Cmd = types.Cmd
  type Flag = types.Flag
  type Validator = types.Validator
  type FlagRegistry = types.FlagRegistry
  type CmdRegistry = types.CmdRegistry
  type CmdConfig = types.CmdConfig
  type FlagConfig = types.FlagConfig
  type Error = types.Error
  ```

#### 7.2 实现 qflag.go
- [ ] 提供便捷函数
  ```go
  func NewCmd(longName, shortName string) *cmd.Cmd
  func NewStringFlag(...) *flag.StringFlag
  func NewIntFlag(...) *flag.IntFlag
  // ... 其他便捷函数
  ```

### 验收标准
- [ ] 用户可以简单导入并使用
- [ ] 编写使用示例

---

## 阶段八: 高级功能 (P3) 

### 目标
实现高级功能, 包括复杂标志类型、命令行补全等。

### 任务列表

#### 8.1 复杂标志类型
- [ ] **time_flag.go** - 时间类型
  ```go
  type TimeFlag struct {
      *BaseFlag[time.Time]
  }
  ```

- [ ] **map_flag.go** - Map 类型
  ```go
  type MapFlag struct {
      *BaseFlag[map[string]string]
  }
  ```

- [ ] **slice_flag.go** - Slice 类型
  ```go
  type StringSliceFlag struct {
      *BaseFlag[[]string]
  }
  ```

- [ ] **size_flag.go** - 文件大小类型
  ```go
  type SizeFlag struct {
      *BaseFlag[int64]
  }
  ```

#### 8.2 命令行补全
- [ ] **bash.go** - Bash 补全
  ```go
  type BashCompleter struct{}

  func (c *BashCompleter) Generate(cmd types.Cmd) string
  ```

- [ ] **powershell.go** - PowerShell 补全
  ```go
  type PowerShellCompleter struct{}

  func (c *PowerShellCompleter) Generate(cmd types.Cmd) string
  ```

#### 8.3 高级验证器
- [ ] **regex_validator.go** - 正则表达式验证
- [ ] **ip_validator.go** - IP 地址验证
- [ ] **email_validator.go** - 邮箱验证

### 验收标准
- [ ] 高级标志类型可以正常使用
- [ ] 命令行补全可以生成补全脚本
- [ ] 高级验证器可以正确验证

---

## 开发顺序建议

```
优先级    阶段              任务数    预计时间
─────────────────────────────────────────────
P0        阶段一: 核心类型    6         3天
P0        阶段二: 注册表      10        3天
P0        阶段三: 基础标志    12        5天
P0        阶段四: 解析器      6         3天
─────────────────────────────────────────────
P1        阶段五: 验证器      4         2天
─────────────────────────────────────────────
P2        阶段六: 帮助生成    3         2天
P2        阶段七: 公共 API    2         1天
─────────────────────────────────────────────
P3        阶段八: 高级功能    10        5天
─────────────────────────────────────────────
总计                        53        24天
```

---

## 里程碑

### M1: 核心类型可用 (第3天) 
- [ ] Cmd 接口定义完成
- [ ] Flag 接口定义完成
- [ ] 注册表接口定义完成

### M2: 命令和标志可用 (第11天) 
- [ ] 可以创建命令
- [ ] 可以添加和查询标志
- [ ] 可以添加和查询子命令

### M3: 完整解析能力 (第14天) 
- [ ] 可以解析长名字和短名字标志
- [ ] 可以递归解析子命令
- [ ] ParseAndRoute 可以自动执行

### M4: 生产就绪 (第17天) 
- [ ] 基础验证器可用
- [ ] 帮助信息生成可用
- [ ] 公共 API 简洁易用

### M5: 功能完整 (可选, 第24天) 
- [ ] 复杂标志类型可用
- [ ] 命令行补全可用
- [ ] 高级验证器可用

---

## 注意事项

1. **每个阶段结束后进行代码审查**
2. **确保所有功能都有单元测试**
3. **保持接口稳定, 避免频繁变更**
4. **优先保证核心功能稳定, 再开发高级功能**
5. **参考现有 qflag 代码, 确保兼容性**

---

## 相关文档

- [REFACTOR_PLAN.md](REFACTOR_PLAN.md) - 详细架构设计
- [README.md](README.md) - 使用说明 (待编写) 
- [EXAMPLES.md](EXAMPLES.md) - 使用示例 (待编写) 
