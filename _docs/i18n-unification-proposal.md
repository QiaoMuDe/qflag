# qflag 多语言统一管理方案

## 背景

目前 qflag 库的 `UseChinese` 字段只能控制帮助信息的中文显示，但其他提示信息（如错误信息、验证提示等）仍然分散在各个模块中，无法统一控制。需要设计一个集中的多语言管理方案。

## 设计目标

1. **统一管理**：所有用户可见的文本信息集中管理
2. **一键切换**：通过 `UseChinese` 字段控制所有文本的语言
3. **易于扩展**：支持未来添加更多语言
4. **向后兼容**：不影响现有 API 和使用方式

## 技术方案

### 1. 核心数据结构

```go
// internal/i18n/messages.go

// MessageKey 消息键类型
type MessageKey string

// 预定义消息键常量
const (
    // 错误信息
    ErrUnknownSubcommand    MessageKey = "err.unknown_subcommand"
    ErrUnknownFlag          MessageKey = "err.unknown_flag"
    ErrMutexFlags           MessageKey = "err.mutex_flags"
    ErrRequiredFlags        MessageKey = "err.required_flags"
    ErrInvalidFlagValue     MessageKey = "err.invalid_flag_value"
    ErrNoRunFunc            MessageKey = "err.no_run_func"
    ErrParseFailed          MessageKey = "err.parse_failed"
    
    // 提示信息
    HintSeeHelp             MessageKey = "hint.see_help"
    HintSimilarCommands     MessageKey = "hint.similar_commands"
    HintSimilarFlags        MessageKey = "hint.similar_flags"
    
    // 验证信息
    ValidateRequired        MessageKey = "validate.required"
    ValidateMutexGroup      MessageKey = "validate.mutex_group"
    
    // 内置标志帮助
    HelpFlagHelp            MessageKey = "help.flag_help"
    HelpFlagVersion         MessageKey = "help.flag_version"
    HelpFlagCompletion      MessageKey = "help.flag_completion"
)

// MessageManager 消息管理器
type MessageManager struct {
    useChinese bool
    messages   map[MessageKey]map[bool]string // key -> (isChinese -> message)
}

// 全局消息管理器实例
var globalManager = &MessageManager{
    useChinese: false,
    messages:   make(map[MessageKey]map[bool]string),
}

// Init 初始化消息管理器
func Init(useChinese bool) {
    globalManager.useChinese = useChinese
    globalManager.registerDefaultMessages()
}

// SetUseChinese 设置是否使用中文
func SetUseChinese(useChinese bool) {
    globalManager.useChinese = useChinese
}

// Get 获取消息
func Get(key MessageKey, args ...interface{}) string {
    return globalManager.Get(key, args...)
}

// Get 获取消息（管理器方法）
func (m *MessageManager) Get(key MessageKey, args ...interface{}) string {
    msgMap, exists := m.messages[key]
    if !exists {
        return string(key) // 返回键名作为回退
    }
    
    msg, exists := msgMap[m.useChinese]
    if !exists {
        // 如果当前语言不存在，返回英文
        msg = msgMap[false]
    }
    
    if len(args) > 0 {
        return fmt.Sprintf(msg, args...)
    }
    return msg
}

// Register 注册消息
func Register(key MessageKey, english, chinese string) {
    globalManager.messages[key] = map[bool]string{
        false: english,
        true:  chinese,
    }
}
```

### 2. 默认消息注册

```go
// internal/i18n/defaults.go

func (m *MessageManager) registerDefaultMessages() {
    // 错误信息
    Register(ErrUnknownSubcommand, 
        "'%s' is not a valid command. See '%s --help'.",
        "'%s' 不是有效的命令。参见 '%s --help'。")
    
    Register(ErrUnknownFlag,
        "unknown flag: '%s'",
        "未知标志: '%s'")
    
    Register(ErrMutexFlags,
        "mutually exclusive flags %v in group '%s' cannot be used together",
        "互斥组 '%s' 中的标志 %v 不能同时使用")
    
    Register(ErrRequiredFlags,
        "required flags %v in group '%s' must be set",
        "必需组 '%s' 中的标志 %v 必须设置")
    
    Register(ErrInvalidFlagValue,
        "invalid value for flag '%s': %v",
        "标志 '%s' 的值无效: %v")
    
    Register(ErrNoRunFunc,
        "command '%s' has no run function set",
        "命令 '%s' 没有设置运行函数")
    
    Register(ErrParseFailed,
        "parse failed: %v",
        "解析失败: %v")
    
    // 提示信息
    Register(HintSeeHelp,
        "See '%s --help' for more information.",
        "参见 '%s --help' 获取更多信息。")
    
    Register(HintSimilarCommands,
        "The most similar commands are",
        "最相似的命令是")
    
    Register(HintSimilarFlags,
        "The most similar flags are",
        "最相似的标志是")
    
    // 验证信息
    Register(ValidateRequired,
        "flag '%s' is required",
        "标志 '%s' 是必需的")
    
    Register(ValidateMutexGroup,
        "one of flags %v in mutex group '%s' must be set",
        "互斥组 '%s' 中的标志 %v 必须设置一个")
    
    // 内置标志帮助
    Register(HelpFlagHelp,
        "Show help information",
        "显示帮助信息")
    
    Register(HelpFlagVersion,
        "Show version information",
        "显示版本信息")
    
    Register(HelpFlagCompletion,
        "Generate completion script (bash|pwsh|fish)",
        "生成补全脚本 (bash|pwsh|fish)")
}
```

### 3. 使用方式

#### 3.1 在错误类型中使用

```go
// internal/types/error.go

func (e *UnknownSubcommandError) Error() string {
    var sb strings.Builder
    // 使用 i18n 获取消息
    sb.WriteString(i18n.Get(i18n.ErrUnknownSubcommand, e.Input, e.Command))
    
    if len(e.Suggestions) > 0 {
        sb.WriteString("\n\n")
        sb.WriteString(i18n.Get(i18n.HintSimilarCommands))
        sb.WriteString("\n")
        for _, sug := range e.Suggestions {
            sb.WriteString(fmt.Sprintf("        %s\n", sug))
        }
    }
    
    return sb.String()
}
```

#### 3.2 在解析器中使用

```go
// internal/parser/parser.go

// 修改前
return fmt.Errorf("cmd %q has no run function set", cmd.Name())

// 修改后
return fmt.Errorf(i18n.Get(i18n.ErrNoRunFunc, cmd.Name()))
```

#### 3.3 在验证器中使用

```go
// internal/parser/parser_validation.go

// 修改前
return fmt.Errorf("mutually exclusive flags %v in group '%s' cannot be used together", 
    setFlagsList, group.Name)

// 修改后
return fmt.Errorf(i18n.Get(i18n.ErrMutexFlags, setFlagsList, group.Name))
```

#### 3.4 在内置标志中使用

```go
// internal/builtin/builtin.go

// 注册帮助标志时
helpFlag := &BoolFlag{
    LongName:  "help",
    ShortName: "h",
    Desc:      i18n.Get(i18n.HelpFlagHelp),  // 自动根据语言返回对应描述
    Default:   false,
}
```

### 4. 与 Cmd 集成

```go
// internal/cmd/cmd_config.go

// ApplyOpts 应用配置时初始化 i18n
func (c *Cmd) ApplyOpts(opts *CmdOpts) error {
    // ... 其他代码 ...
    
    // 初始化多语言（根命令设置全局语言）
    if c.IsRoot() && opts.UseChinese {
        i18n.SetUseChinese(true)
    }
    
    // ... 其他代码 ...
}

// 或者提供更明确的控制方法
func (c *Cmd) SetUseChinese(useChinese bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.config.UseChinese = useChinese
    
    // 设置全局语言
    i18n.SetUseChinese(useChinese)
}
```

### 5. 扩展性设计

支持未来添加更多语言：

```go
// Language 语言类型
type Language string

const (
    LangEnglish Language = "en"
    LangChinese Language = "zh"
    LangJapanese Language = "ja"
    // ... 更多语言
)

// MessageManager 扩展
type MessageManager struct {
    currentLang Language
    messages    map[MessageKey]map[Language]string
}

// Register 注册多语言消息
func Register(key MessageKey, translations map[Language]string) {
    globalManager.messages[key] = translations
}

// 使用示例
Register(ErrUnknownSubcommand, map[Language]string{
    LangEnglish:  "'%s' is not a valid command.",
    LangChinese:  "'%s' 不是有效的命令。",
    LangJapanese: "'%s' は有効なコマンドではありません。",
})
```

## 实施步骤

### 第一阶段：基础架构
1. 创建 `internal/i18n` 包
2. 实现 `MessageManager` 和基础 API
3. 注册默认的中英文消息

### 第二阶段：错误信息迁移
1. 修改 `internal/types/error.go` 中的错误类型
2. 修改 `internal/parser/parser.go` 中的错误返回
3. 修改 `internal/parser/parser_validation.go` 中的验证错误

### 第三阶段：提示信息迁移
1. 修改 `internal/builtin/builtin.go` 中的内置标志描述
2. 修改帮助信息生成器
3. 修改补全脚本生成器

### 第四阶段：测试和文档
1. 添加单元测试验证中英文切换
2. 更新 API 文档
3. 添加使用示例

## 注意事项

1. **性能考虑**：消息获取使用 map 查找，性能开销极小
2. **线程安全**：`MessageManager` 使用读写锁保证并发安全
3. **回退机制**：如果当前语言的消息不存在，自动回退到英文
4. **动态切换**：支持运行时动态切换语言（需要重新生成帮助信息等）
5. **参数格式化**：支持 `fmt.Sprintf` 风格的参数替换

## 相关文件

- `internal/i18n/messages.go` - 消息管理器核心实现
- `internal/i18n/defaults.go` - 默认消息定义
- `internal/types/error.go` - 错误类型（需要修改）
- `internal/parser/*.go` - 解析器（需要修改）
- `internal/builtin/builtin.go` - 内置标志（需要修改）
- `internal/help/*.go` - 帮助生成器（需要修改）
