# 内置标志系统设计方案

## 1. 设计目标

实现一个内置标志系统, 用于处理常见的内置标志 (如 `-h/--help`, `-v/--version`) , 在解析完成后自动检查这些标志是否被调用, 并执行相应的操作。

## 2. 设计思路

1. **条件注册**: 根据命令配置决定是否注册某些内置标志 (如版本标志) 
2. **解析前注册**: 在解析前注册必要的内置标志
3. **解析后检查**: 在解析完成后检查内置标志是否被设置
4. **自动执行**: 如果内置标志被设置, 自动执行相应操作并退出

## 3. 实现方案

### 3.1 内置标志类型定义

在 `internal/types` 中添加内置标志类型: 

```go
// BuiltinFlagType 内置标志类型
type BuiltinFlagType int

const (
    // HelpFlag 帮助标志 - 总是注册
    HelpFlag BuiltinFlagType = iota
    // VersionFlag 版本标志 - 只有设置了版本信息时才注册
    VersionFlag
    // 可以继续添加其他内置标志
)
```

### 3.2 内置标志处理器接口

```go
// BuiltinFlagHandler 内置标志处理器接口
type BuiltinFlagHandler interface {
    // Handle 处理内置标志
    Handle(cmd Command) error
    // Type 返回标志类型
    Type() BuiltinFlagType
    // ShouldRegister 判断是否应该注册此标志
    ShouldRegister(cmd Command) bool
}
```

### 3.3 具体处理器实现

```go
// HelpHandler 帮助标志处理器
type HelpHandler struct{}

func (h *HelpHandler) Handle(cmd Command) error {
    cmd.PrintHelp()
    os.Exit(0)
    return nil
}

func (h *HelpHandler) Type() BuiltinFlagType {
    return HelpFlag
}

func (h *HelpHandler) ShouldRegister(cmd Command) bool {
    // 帮助标志总是注册
    return true
}

// VersionHandler 版本标志处理器
type VersionHandler struct{}

func (h *VersionHandler) Handle(cmd Command) error {
    fmt.Println(cmd.Version())
    os.Exit(0)
    return nil
}

func (h *VersionHandler) Type() BuiltinFlagType {
    return VersionFlag
}

func (h *VersionHandler) ShouldRegister(cmd Command) bool {
    // 只有设置了版本信息时才注册版本标志
    return cmd.Version() != ""
}
```

### 3.4 内置标志管理器

```go
// BuiltinFlagManager 内置标志管理器
type BuiltinFlagManager struct {
    handlers map[BuiltinFlagType]BuiltinFlagHandler
    flags    map[string]BuiltinFlagType // 标志名到类型的映射
}

// NewBuiltinFlagManager 创建内置标志管理器
func NewBuiltinFlagManager() *BuiltinFlagManager {
    m := &BuiltinFlagManager{
        handlers: make(map[BuiltinFlagType]BuiltinFlagHandler),
        flags:    make(map[string]BuiltinFlagType),
    }
    
    // 注册默认处理器
    m.RegisterHandler(&HelpHandler{})
    m.RegisterHandler(&VersionHandler{})
    
    return m
}

// RegisterHandler 注册内置标志处理器
func (m *BuiltinFlagManager) RegisterHandler(handler BuiltinFlagHandler) {
    flagType := handler.Type()
    m.handlers[flagType] = handler
    
    // 注册标志名映射
    switch flagType {
    case HelpFlag:
        m.flags["help"] = HelpFlag
        m.flags["h"] = HelpFlag
    case VersionFlag:
        m.flags["version"] = VersionFlag
        m.flags["v"] = VersionFlag
    }
}

// RegisterBuiltinFlags 注册内置标志
func (m *BuiltinFlagManager) RegisterBuiltinFlags(cmd Command) error {
    for _, handler := range m.handlers {
        // 检查是否应该注册此标志
        if handler.ShouldRegister(cmd) {
            switch handler.Type() {
            case HelpFlag:
                // 根据命令的语言设置使用相应的描述信息
                var desc string
                if cmd.Config().UseChinese {
                    desc = "显示帮助信息"
                } else {
                    desc = "Show help information"
                }
                helpFlag := flag.NewBoolFlag("help", "h", desc, false)
                if err := cmd.AddFlag(helpFlag); err != nil {
                    return err
                }
            case VersionFlag:
                // 根据命令的语言设置使用相应的描述信息
                var desc string
                if cmd.Config().UseChinese {
                    desc = "显示版本信息"
                } else {
                    desc = "Show version information"
                }
                versionFlag := flag.NewBoolFlag("version", "v", desc, false)
                if err := cmd.AddFlag(versionFlag); err != nil {
                    return err
                }
            }
        }
    }
    
    return nil
}

// HandleBuiltinFlags 处理内置标志
func (m *BuiltinFlagManager) HandleBuiltinFlags(cmd Command) error {
    flags := cmd.Flags()
    
    for _, f := range flags {
        // 检查是否是内置标志
        if flagType, isBuiltin := m.isBuiltinFlag(f); isBuiltin {
            // 检查是否被设置
            if f.IsSet() {
                // 执行处理器
                if handler, exists := m.handlers[flagType]; exists {
                    return handler.Handle(cmd)
                }
            }
        }
    }
    
    return nil
}

// isBuiltinFlag 检查是否是内置标志
func (m *BuiltinFlagManager) isBuiltinFlag(f Flag) (BuiltinFlagType, bool) {
    // 检查长名称
    if flagType, exists := m.flags[f.Name()]; exists {
        return flagType, true
    }
    
    // 检查短名称
    if flagType, exists := m.flags[f.ShortName()]; exists {
        return flagType, true
    }
    
    return 0, false
}
```

### 3.5 修改解析器

在 `DefaultParser` 中添加内置标志管理器, 并在解析前后进行相应操作: 

```go
// DefaultParser 默认解析器实现
type DefaultParser struct {
    flagSet       *flag.FlagSet
    errorHandling types.ErrorHandling
    builtinMgr    *BuiltinFlagManager // 添加内置标志管理器
}

// NewDefaultParser 创建默认解析器实例
func NewDefaultParser(errorHandling types.ErrorHandling) types.Parser {
    return &DefaultParser{
        errorHandling: errorHandling,
        builtinMgr:    NewBuiltinFlagManager(), // 初始化内置标志管理器
    }
}

// ParseOnly 仅解析命令行参数, 不执行子命令路由
func (p *DefaultParser) ParseOnly(cmd types.Command, args []string) error {
    // 注册内置标志
    if err := p.builtinMgr.RegisterBuiltinFlags(cmd); err != nil {
        return err
    }
    
    // ... 现有解析代码 ...
    
    // 在解析完成后检查内置标志
    if err := p.builtinMgr.HandleBuiltinFlags(cmd); err != nil {
        return err
    }
    
    return nil
}
```

## 4. 优点

1. **条件注册**: 根据命令配置决定是否注册某些内置标志, 更加灵活
2. **自动化**: 内置标志自动注册和处理, 无需用户手动操作
3. **可扩展**: 通过添加新的处理器可以轻松支持更多内置标志
4. **非侵入性**: 不影响现有代码结构, 只是添加了额外的功能
5. **用户体验**: 只显示相关的内置标志, 避免混淆
6. **国际化支持**: 根据命令的语言设置自动使用相应的描述信息, 支持中英文切换

## 5. 使用示例

```go
// 示例1: 有版本信息的命令 (中文) 
cmd := cmd.NewCmd("myapp", "m", types.ContinueOnError)
cmd.SetDesc("我的应用程序")
cmd.SetVersion("1.0.0") // 设置版本信息
cmd.SetChinese(true)     // 设置为中文

// 解析参数
if err := cmd.Parse(os.Args[1:]); err != nil {
    log.Fatal(err)
}

// 用户可以使用 -h/--help 查看帮助 (显示中文描述) 
// 用户可以使用 -v/--version 查看版本 (显示中文描述) 

// 示例2: 没有版本信息的命令 (英文) 
cmd2 := cmd.NewCmd("myapp2", "m2", types.ContinueOnError)
cmd2.SetDesc("My application 2")
// 不设置版本信息
cmd2.SetChinese(false)    // 设置为英文

// 解析参数
if err := cmd2.Parse(os.Args[1:]); err != nil {
    log.Fatal(err)
}

// 用户只能使用 -h/--help 查看帮助 (显示英文描述) 
// -v/--version 不可用, 因为没有注册版本标志

// 示例3: 混合语言设置
cmd3 := cmd.NewCmd("myapp3", "m3", types.ContinueOnError)
cmd3.SetDesc("我的应用程序3")
cmd3.SetVersion("1.0.0") // 设置版本信息
cmd3.SetChinese(true)     // 设置为中文

// 动态切换语言
if someCondition {
    cmd3.SetChinese(false) // 切换到英文
}

// 解析参数
if err := cmd3.Parse(os.Args[1:]); err != nil {
    log.Fatal(err)
}

// 根据最终的语言设置显示相应的描述信息
```

## 6. 实现步骤

1. 在 `internal/types` 中添加内置标志相关接口和类型
2. 创建 `internal/builtin` 包, 实现内置标志管理器和处理器
3. 修改 `internal/parser/parser.go`, 集成内置标志管理器
4. 编写测试用例验证功能
5. 更新文档和示例

## 7. 扩展性

这个设计具有良好的扩展性, 可以轻松添加新的内置标志: 

1. 在 `BuiltinFlagType` 中添加新的标志类型
2. 实现新的 `BuiltinFlagHandler`
3. 在 `RegisterHandler` 中添加标志名映射
4. 在 `RegisterBuiltinFlags` 中添加标志创建逻辑, 并根据语言设置使用相应的描述信息

例如, 可以添加一个 `--config` 标志来显示配置文件路径, 或者添加一个 `--verbose` 标志来切换详细输出模式。

对于国际化支持, 新增的内置标志也应该遵循相同的原则: 
- 根据命令的语言设置 (`cmd.Config().UseChinese`) 决定使用中文还是英文描述
- 在 `RegisterBuiltinFlags` 中实现条件逻辑, 提供不同语言的描述信息
- 确保所有内置标志都支持国际化, 提供一致的用户体验