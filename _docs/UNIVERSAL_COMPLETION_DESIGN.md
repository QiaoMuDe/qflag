# 通用Shell补全生成器设计文档

## 1. 问题分析

### 当前实现的问题
1. **硬编码Shell类型**: 当前实现直接针对Bash和PowerShell, 添加新Shell需要修改核心代码
2. **紧耦合设计**: Shell特定逻辑与核心补全逻辑混合在一起
3. **重复代码**: 不同Shell的实现存在相似的逻辑, 但没有抽象
4. **难以测试**: Shell特定逻辑与核心逻辑耦合, 导致单元测试困难
5. **配置不灵活**: 没有统一的配置系统, 无法根据Shell特性调整行为

## 2. 设计目标

1. **可扩展性**: 轻松添加新Shell支持, 无需修改现有代码
2. **可维护性**: 清晰的职责分离, 便于维护和调试
3. **可测试性**: 各组件可独立测试, 提高代码质量
4. **性能优化**: 保持或提升当前实现的性能优势
5. **配置灵活性**: 支持Shell特定的配置选项

## 3. 架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    CompletionGenerator                    │
│                    (核心协调器)                          │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                 CompletionContext                         │
│              (补全上下文与数据模型)                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                ShellRenderer接口                         │
│              (Shell渲染器抽象接口)                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│BashRenderer │ │PwshRenderer │ │ZshRenderer  │
│             │ │             │ │             │
└─────────────┘ └─────────────┘ └─────────────┘
```

### 3.2 核心组件

#### 3.2.1 CompletionGenerator (核心协调器)

```go
type CompletionGenerator struct {
    renderers map[ShellType]ShellRenderer
    config    *GeneratorConfig
}

// 主要职责: 
// 1. 管理所有Shell渲染器
// 2. 协调数据收集和渲染过程
// 3. 提供统一的公共API
```

#### 3.2.2 CompletionContext (补全上下文)

```go
type CompletionContext struct {
    Command       Command
    FlagParams   []FlagParam
    CommandTree  CommandTreeNode
    ShellConfig  ShellSpecificConfig
}

// 主要职责: 
// 1. 封装补全所需的所有数据
// 2. 提供统一的数据访问接口
// 3. 支持Shell特定的配置
```

#### 3.2.3 ShellRenderer接口 (Shell渲染器抽象)

```go
type ShellRenderer interface {
    // 渲染完整的补全脚本
    Render(ctx *CompletionContext) (string, error)
    
    // 获取Shell类型
    GetShellType() ShellType
    
    // 验证Shell特定配置
    ValidateConfig(config ShellSpecificConfig) error
    
    // 获取默认配置
    GetDefaultConfig() ShellSpecificConfig
}

// 主要职责: 
// 1. 定义Shell渲染器的统一接口
// 2. 强制实现必要的渲染方法
// 3. 支持Shell特定的配置验证
```

## 4. 接口设计

### 4.1 核心接口

```go
// ShellRenderer Shell渲染器接口
type ShellRenderer interface {
    // Render 渲染完整的补全脚本
    Render(ctx *CompletionContext) (string, error)
    
    // GetShellType 获取Shell类型
    GetShellType() ShellType
    
    // ValidateConfig 验证Shell特定配置
    ValidateConfig(config ShellSpecificConfig) error
    
    // GetDefaultConfig 获取默认配置
    GetDefaultConfig() ShellSpecificConfig
    
    // GetSupportedFeatures 获取支持的特性列表
    GetSupportedFeatures() []ShellFeature
}

// DataCollector 数据收集器接口
type DataCollector interface {
    // CollectFlagParams 收集标志参数
    CollectFlagParams(cmd Command) []FlagParam
    
    // BuildCommandTree 构建命令树
    BuildCommandTree(cmd Command) CommandTreeNode
    
    // CollectCompletionOptions 收集补全选项
    CollectCompletionOptions(cmd Command) []string
}

// TemplateEngine 模板引擎接口
type TemplateEngine interface {
    // ExecuteTemplate 执行模板
    ExecuteTemplate(name string, data interface{}) (string, error)
    
    // LoadTemplates 加载模板
    LoadTemplates() error
}
```

### 4.2 数据模型

```go
// ShellType Shell类型枚举
type ShellType string

const (
    Bash       ShellType = "bash"
    PowerShell ShellType = "powershell"
    Zsh        ShellType = "zsh"
    Fish       ShellType = "fish"
)

// ShellFeature Shell特性枚举
type ShellFeature string

const (
    FeatureFuzzyMatch     ShellFeature = "fuzzy_match"
    FeatureEnumCompletion ShellFeature = "enum_completion"
    FeatureNestedCommands ShellFeature = "nested_commands"
)

// FlagParam 标志参数模型
type FlagParam struct {
    Name         string            // 标志名称
    CommandPath  string            // 命令路径
    Type         FlagType          // 标志类型
    ValueType    string            // 值类型
    EnumOptions  []string          // 枚举选项
    Description  string            // 描述信息
    Required     bool              // 是否必需
    ShellSpecific map[string]interface{} // Shell特定属性
}

// CommandTreeNode 命令树节点
type CommandTreeNode struct {
    Name         string                // 命令名称
    Path         string                // 命令路径
    Description  string                // 描述信息
    Flags        []FlagParam           // 标志列表
    SubCommands  []CommandTreeNode      // 子命令列表
    ShellSpecific map[string]interface{} // Shell特定属性
}

// ShellSpecificConfig Shell特定配置
type ShellSpecificConfig map[string]interface{}

// GeneratorConfig 生成器配置
type GeneratorConfig struct {
    EnableCache        bool                        // 是否启用缓存
    CacheSize         int                         // 缓存大小
    EnableFuzzyMatch  bool                        // 是否启用模糊匹配
    MaxCompletionItems int                         // 最大补全项数
    ShellConfigs      map[ShellType]ShellSpecificConfig // Shell特定配置
}
```

## 5. 实现策略

### 5.1 渲染器实现模式

```go
// BaseShellRenderer 基础渲染器, 提供通用功能
type BaseShellRenderer struct {
    templateEngine TemplateEngine
    config       ShellSpecificConfig
}

// BashRenderer Bash渲染器
type BashRenderer struct {
    BaseShellRenderer
}

// 实现ShellRenderer接口
func (r *BashRenderer) Render(ctx *CompletionContext) (string, error) {
    // 使用模板引擎渲染Bash脚本
    return r.templateEngine.ExecuteTemplate("bash_completion", ctx)
}

func (r *BashRenderer) GetShellType() ShellType {
    return Bash
}

func (r *BashRenderer) GetSupportedFeatures() []ShellFeature {
    return []ShellFeature{FeatureFuzzyMatch, FeatureEnumCompletion, FeatureNestedCommands}
}
```

### 5.2 插件式注册机制

```go
// RendererRegistry 渲染器注册表
type RendererRegistry struct {
    renderers map[ShellType]ShellRenderer
    mu        sync.RWMutex
}

// RegisterRenderer 注册渲染器
func (r *RendererRegistry) RegisterRenderer(renderer ShellRenderer) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    shellType := renderer.GetShellType()
    if _, exists := r.renderers[shellType]; exists {
        return fmt.Errorf("renderer for %s already registered", shellType)
    }
    
    r.renderers[shellType] = renderer
    return nil
}

// GetRenderer 获取渲染器
func (r *RendererRegistry) GetRenderer(shellType ShellType) (ShellRenderer, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    renderer, exists := r.renderers[shellType]
    if !exists {
        return nil, fmt.Errorf("no renderer found for %s", shellType)
    }
    
    return renderer, nil
}
```

### 5.3 模板系统设计

```
templates/
├── base/
│   ├── header.tmpl          # 通用脚本头部
│   ├── footer.tmpl          # 通用脚本尾部
│   └── functions.tmpl       # 通用函数定义
├── bash/
│   ├── completion.tmpl       # Bash主补全脚本
│   ├── command-tree.tmpl     # Bash命令树
│   └── flag-handler.tmpl     # Bash标志处理
├── powershell/
│   ├── completion.tmpl       # PowerShell主补全脚本
│   ├── command-tree.tmpl     # PowerShell命令树
│   └── flag-handler.tmpl     # PowerShell标志处理
└── zsh/
    ├── completion.tmpl       # Zsh主补全脚本
    ├── command-tree.tmpl     # Zsh命令树
    └── flag-handler.tmpl     # Zsh标志处理
```

## 6. 扩展新Shell的步骤

### 6.1 添加新Shell类型

```go
const (
    // 现有Shell类型...
    Fish ShellType = "fish"  // 新增Fish Shell
)
```

### 6.2 实现ShellRenderer接口

```go
type FishRenderer struct {
    BaseShellRenderer
}

func (r *FishRenderer) Render(ctx *CompletionContext) (string, error) {
    return r.templateEngine.ExecuteTemplate("fish_completion", ctx)
}

func (r *FishRenderer) GetShellType() ShellType {
    return Fish
}

func (r *FishRenderer) GetSupportedFeatures() []ShellFeature {
    return []ShellFeature{FeatureFuzzyMatch, FeatureEnumCompletion}
}
```

### 6.3 创建模板文件

在`templates/fish/`目录下创建Fish Shell特定的模板文件。

### 6.4 注册新渲染器

```go
func init() {
    renderer := &FishRenderer{
        BaseShellRenderer: BaseShellRenderer{
            templateEngine: NewTemplateEngine("fish"),
            config:        make(ShellSpecificConfig),
        },
    }
    
    CompletionGenerator.RegisterRenderer(renderer)
}
```

## 7. 配置系统设计

### 7.1 分层配置

```go
// 全局默认配置
var DefaultConfig = &GeneratorConfig{
    EnableCache:        true,
    CacheSize:         100,
    EnableFuzzyMatch:  true,
    MaxCompletionItems: 1000,
    ShellConfigs: make(map[ShellType]ShellSpecificConfig),
}

// Shell特定默认配置
var ShellDefaultConfigs = map[ShellType]ShellSpecificConfig{
    Bash: {
        "enable_fuzzy_match": true,
        "max_items": 1000,
        "cache_enabled": true,
    },
    PowerShell: {
        "enable_fuzzy_match": false,
        "max_items": 500,
        "cache_enabled": true,
    },
}
```

### 7.2 配置验证

```go
func (g *CompletionGenerator) ValidateConfig(config *GeneratorConfig) error {
    // 验证全局配置
    if config.MaxCompletionItems <= 0 {
        return errors.New("MaxCompletionItems must be positive")
    }
    
    // 验证Shell特定配置
    for shellType, shellConfig := range config.ShellConfigs {
        renderer, err := g.GetRenderer(shellType)
        if err != nil {
            return err
        }
        
        if err := renderer.ValidateConfig(shellConfig); err != nil {
            return fmt.Errorf("invalid config for %s: %w", shellType, err)
        }
    }
    
    return nil
}
```

## 8. 性能优化策略

### 8.1 对象池优化

```go
// 对象池管理
var (
    contextPool = sync.Pool{
        New: func() interface{} {
            return &CompletionContext{}
        },
    }
    
    flagParamPool = sync.Pool{
        New: func() interface{} {
            return make([]FlagParam, 0, 32)
        },
    }
)

// 获取上下文对象
func GetCompletionContext() *CompletionContext {
    return contextPool.Get().(*CompletionContext)
}

// 归还上下文对象
func PutCompletionContext(ctx *CompletionContext) {
    ctx.Reset()
    contextPool.Put(ctx)
}
```

### 8.2 缓存机制

```go
// CompletionCache 补全缓存
type CompletionCache struct {
    cache map[string]*CacheEntry
    mu    sync.RWMutex
    size  int
    maxSize int
}

type CacheEntry struct {
    Value      string
    Expiration time.Time
}

// Get 获取缓存项
func (c *CompletionCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    entry, exists := c.cache[key]
    if !exists || time.Now().After(entry.Expiration) {
        return "", false
    }
    
    return entry.Value, true
}
```

## 9. 测试策略

### 9.1 单元测试

```go
// 测试渲染器接口实现
func TestShellRenderer_Render(t *testing.T) {
    tests := []struct {
        name     string
        renderer ShellRenderer
        context  *CompletionContext
        expected string
    }{
        {
            name:     "Bash basic completion",
            renderer: &BashRenderer{},
            context:  createTestContext(Bash),
            expected: "bash_completion_script",
        },
        // 更多测试用例...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := tt.renderer.Render(tt.context)
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 9.2 集成测试

```go
// 测试完整的补全生成流程
func TestCompletionGenerator_Generate(t *testing.T) {
    generator := NewCompletionGenerator()
    
    // 注册测试渲染器
    generator.RegisterRenderer(&TestRenderer{})
    
    // 创建测试命令
    cmd := createTestCommand()
    
    // 生成补全脚本
    result, err := generator.Generate(cmd, "test")
    
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

## 10. 迁移策略

### 10.1 渐进式迁移

1. **第一阶段**: 实现新的接口和核心架构
2. **第二阶段**: 将现有Bash和PowerShell实现迁移到新架构
3. **第三阶段**: 添加新Shell支持 (如Zsh、Fish) 
4. **第四阶段**: 优化性能和添加高级特性

### 10.2 兼容性保证

```go
// 保持向后兼容的API
func GenerateCompletion(cmd Command, shellType string) (string, error) {
    // 转换为新的API调用
    return NewCompletionGenerator().Generate(cmd, ShellType(shellType))
}
```

## 11. 总结

这个通用设计通过以下方式解决了当前实现的问题: 

1. **插件式架构**: 通过ShellRenderer接口实现插件式扩展, 添加新Shell无需修改核心代码
2. **职责分离**: 核心逻辑与Shell特定逻辑完全分离, 提高可维护性
3. **统一配置**: 提供分层的配置系统, 支持全局和Shell特定配置
4. **模板系统**: 使用模板引擎实现Shell特定的脚本生成, 提高代码复用
5. **性能优化**: 保持现有的性能优化策略, 如对象池和缓存机制

这个设计不仅解决了当前的问题, 还为未来的扩展提供了坚实的基础。