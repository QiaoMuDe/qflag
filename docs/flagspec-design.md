# FlagSpec 设计方案

## 概述

本文档详细描述了 FlagSpec 结构体的设计方案, 该设计与 CmdSpec 形成配套的规范模式, 用于通过结构化配置创建标志。

## 设计目标

1. **一致性** - 与 CmdSpec 形成完整的设计体系
2. **结构化** - 集中定义所有标志属性
3. **类型安全** - 编译时类型检查和验证
4. **扩展性** - 易于添加新类型和验证器
5. **配置友好** - 支持从配置文件加载

## 核心设计

### 1. FlagSpec 结构体

```go
// FlagSpec 标志规格结构体
//
// FlagSpec 提供了通过规格创建标志的方式, 包含标志的所有属性。
// 这种方式比工厂函数更加直观和集中。
type FlagSpec struct {
    // 基本属性
    LongName    string        // 长名称
    ShortName   string        // 短名称
    Desc        string        // 描述
    Default     any           // 默认值
    
    // 类型定义
    Type        FlagType      // 标志类型枚举
    
    // 验证相关
    Validator   Validator     // 验证器
    
    // 枚举值 (仅用于枚举类型) 
    EnumValues  []any         // 枚举值
}
```

### 2. 标志类型枚举

```go
// FlagType 标志类型枚举
//
// 直接使用 types.FlagType, 避免重复定义
type FlagType = types.FlagType
```

### 3. 创建方法

```go
// NewFlagFromSpec 从规格创建标志
//
// 参数:
//   - spec: 标志规格结构体
//
// 返回值:
//   - types.Flag: 创建的标志实例
//   - error: 创建失败时返回错误
//
// 功能说明: 
//   - 根据规格结构体创建标志并添加到当前命令
//   - 自动设置所有属性和配置
//   - 应用验证器和枚举值
//   - 使用defer捕获panic, 转换为错误返回
func (c *Cmd) NewFlagFromSpec(spec *FlagSpec) (types.Flag, error) {
    // 使用defer捕获panic, 转换为错误返回
    defer func() {
        if r := recover(); r != nil {
            // 将panic转换为错误
            switch x := r.(type) {
            case string:
                err = types.NewError("PANIC", x, nil)
            case error:
                err = types.NewError("PANIC", x.Error(), x)
            default:
                err = types.NewError("PANIC", fmt.Sprintf("%v", x), nil)
            }
            flag = nil
        }
    }()

    // 验证参数
    if spec == nil {
        return nil, types.NewError("INVALID_FLAG_SPEC", "flag spec cannot be nil", nil)
    }

    // 根据类型创建标志
    switch spec.Type {
    case types.FlagTypeBool:
        defaultValue := spec.Default.(bool)
        flag = c.Bool(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeString:
        defaultValue := spec.Default.(string)
        flag = c.String(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeInt:
        defaultValue := spec.Default.(int)
        flag = c.Int(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeInt64:
        defaultValue := spec.Default.(int64)
        flag = c.Int64(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeUint:
        defaultValue := spec.Default.(uint)
        flag = c.Uint(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeUint8:
        defaultValue := spec.Default.(uint8)
        flag = c.Uint8(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeUint16:
        defaultValue := spec.Default.(uint16)
        flag = c.Uint16(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeUint32:
        defaultValue := spec.Default.(uint32)
        flag = c.Uint32(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeUint64:
        defaultValue := spec.Default.(uint64)
        flag = c.Uint64(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeFloat64:
        defaultValue := spec.Default.(float64)
        flag = c.Float64(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeDuration:
        defaultValue := spec.Default.(time.Duration)
        flag = c.Duration(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeTime:
        defaultValue := spec.Default.(time.Time)
        flag = c.Time(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeSize:
        defaultValue := spec.Default.(int64)
        flag = c.Size(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeStringSlice:
        defaultValue := spec.Default.([]string)
        flag = c.StringSlice(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeIntSlice:
        defaultValue := spec.Default.([]int)
        flag = c.IntSlice(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeInt64Slice:
        defaultValue := spec.Default.([]int64)
        flag = c.Int64Slice(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeMap:
        defaultValue := spec.Default.(map[string]string)
        flag = c.Map(spec.LongName, spec.ShortName, spec.Desc, defaultValue)
    case types.FlagTypeEnum:
        defaultValue := spec.Default.(string)
        // 将枚举值转换为字符串切片
        enumStrValues := make([]string, len(spec.EnumValues))
        for i, value := range spec.EnumValues {
            enumStrValues[i] = fmt.Sprintf("%v", value)
        }
        flag = c.Enum(spec.LongName, spec.ShortName, spec.Desc, defaultValue, enumStrValues)
    default:
        return nil, types.NewError("UNKNOWN_FLAG_TYPE", "unknown flag type", nil)
    }
    
    // 应用验证器
    if spec.Validator != nil {
        flag.SetValidator(spec.Validator)
    }
    
    return flag, nil
}
```

### 4. 便捷构造函数

```go
// NewFlagSpec 创建新的标志规格
//
// 参数:
//   - longName: 长名称
//   - shortName: 短名称
//   - flagType: 标志类型
//
// 返回值:
//   - *FlagSpec: 创建的标志规格实例
//
// 功能说明: 
//   - 创建带有默认值的标志规格
//   - 初始化所有字段为合理默认值
//   - 便于链式调用配置
func NewFlagSpec(longName, shortName string, flagType FlagType) *FlagSpec {
    return &FlagSpec{
        LongName:   longName,
        ShortName:  shortName,
        Desc:       "",
        Default:    nil,
        Type:       flagType,
        Validator:  nil,
        EnumValues: []any{},
    }
}
```

## 使用示例

### 1. 基本用法

```go
// 创建布尔标志规格
verboseSpec := NewFlagSpec("verbose", "v", types.FlagTypeBool)
verboseSpec.Desc = "详细输出"
verboseSpec.Default = false

// 创建字符串标志规格
inputSpec := NewFlagSpec("input", "i", types.FlagTypeString)
inputSpec.Desc = "输入文件"
inputSpec.Default = ""
inputSpec.Validator = validator.Required()

// 从规格创建标志
verboseFlag, err := cmd.NewFlagFromSpec(verboseSpec)
if err != nil {
    log.Fatal(err)
}

inputFlag, err := cmd.NewFlagFromSpec(inputSpec)
if err != nil {
    log.Fatal(err)
}
```

### 2. 枚举标志

```go
// 创建枚举标志规格
modeSpec := NewFlagSpec("mode", "m", types.FlagTypeEnum)
modeSpec.Desc = "运行模式"
modeSpec.Default = "auto"
modeSpec.EnumValues = []any{"auto", "manual", "debug"}

// 从规格创建标志
modeFlag, err := cmd.NewFlagFromSpec(modeSpec)
if err != nil {
    log.Fatal(err)
}
```

### 3. 高级类型示例

```go
// 创建时间标志规格
timeSpec := NewFlagSpec("start-time", "t", types.FlagTypeTime)
timeSpec.Desc = "开始时间"
timeSpec.Default = time.Now()

// 创建字符串切片标志规格
filesSpec := NewFlagSpec("files", "f", types.FlagTypeStringSlice)
filesSpec.Desc = "文件列表"
filesSpec.Default = []string{}

// 创建映射标志规格
configSpec := NewFlagSpec("config", "c", types.FlagTypeMap)
configSpec.Desc = "配置项"
configSpec.Default = map[string]string{}

// 从规格创建标志
timeFlag, err := cmd.NewFlagFromSpec(timeSpec)
if err != nil {
    log.Fatal(err)
}

filesFlag, err := cmd.NewFlagFromSpec(filesSpec)
if err != nil {
    log.Fatal(err)
}

configFlag, err := cmd.NewFlagFromSpec(configSpec)
if err != nil {
    log.Fatal(err)
}
```

```go
// 创建数值标志规格
countSpec := NewFlagSpec("count", "c", types.FlagTypeInt)
countSpec.Desc = "计数器"
countSpec.Default = 10
countSpec.Validator = validator.Range(0, 100)

// 创建字符串标志规格
nameSpec := NewFlagSpec("name", "n", types.FlagTypeString)
nameSpec.Desc = "名称"
nameSpec.Default = ""
nameSpec.Validator = validator.Regex("^[a-zA-Z]+$")

// 从规格创建标志
countFlag, err := cmd.NewFlagFromSpec(countSpec)
if err != nil {
    log.Fatal(err)
}

nameFlag, err := cmd.NewFlagFromSpec(nameSpec)
if err != nil {
    log.Fatal(err)
}
```

### 4. 基本验证

```go
// 创建数值标志规格
countSpec := NewFlagSpec("count", "c", types.FlagTypeInt)
countSpec.Desc = "计数器"
countSpec.Default = 10
countSpec.Validator = validator.Range(0, 100)

// 创建字符串标志规格
nameSpec := NewFlagSpec("name", "n", types.FlagTypeString)
nameSpec.Desc = "名称"
nameSpec.Default = ""
nameSpec.Validator = validator.Regex("^[a-zA-Z]+$")

// 从规格创建标志
countFlag, err := NewFlagFromSpec(cmd, countSpec)
if err != nil {
    log.Fatal(err)
}

nameFlag, err := NewFlagFromSpec(cmd, nameSpec)
if err != nil {
    log.Fatal(err)
}
```

### 5. 与 CmdSpec 结合使用

```go
// 创建命令规格
appSpec := cmd.NewCmdSpec("myapp", "app")
appSpec.Desc = "我的应用程序"
appSpec.RunFunc = func(cmd types.Command) error {
    // 命令执行逻辑
    return nil
}

// 从命令规格创建命令
app, err := cmd.NewCmdFromSpec(appSpec)
if err != nil {
    log.Fatal(err)
}

// 创建标志规格
verboseSpec := NewFlagSpec("verbose", "v", types.FlagTypeBool)
verboseSpec.Desc = "详细输出"
verboseSpec.Default = false

inputSpec := NewFlagSpec("input", "i", types.FlagTypeString)
inputSpec.Desc = "输入文件"
inputSpec.Default = ""

// 从规格创建标志并添加到命令
verboseFlag, err := app.NewFlagFromSpec(verboseSpec)
if err != nil {
    log.Fatal(err)
}

inputFlag, err := app.NewFlagFromSpec(inputSpec)
if err != nil {
    log.Fatal(err)
}
```

## 配置文件支持

### 1. JSON 配置示例

```json
{
  "flags": [
    {
      "longName": "verbose",
      "shortName": "v",
      "desc": "详细输出",
      "default": false,
      "type": 0
    },
    {
      "longName": "input",
      "shortName": "i",
      "desc": "输入文件",
      "default": "",
      "type": 1,
      "validators": ["required", "file_exists"]
    },
    {
      "longName": "mode",
      "shortName": "m",
      "desc": "运行模式",
      "default": "auto",
      "type": 7,
      "enumValues": ["auto", "manual", "debug"]
    }
  ]
}
```

### 2. 从配置文件加载

```go
// 从JSON配置加载标志规格
func LoadFlagSpecsFromJSON(jsonData []byte) ([]*FlagSpec, error) {
    var config struct {
        Flags []*FlagSpec `json:"flags"`
    }
    
    if err := json.Unmarshal(jsonData, &config); err != nil {
        return nil, err
    }
    
    return config.Flags, nil
}
```

## 优势分析

### 1. 一致性优势
- 与 CmdSpec 形成完整的设计模式
- 统一的配置方式, 无论是命令还是标志
- 一致的错误处理和验证机制

### 2. 结构化优势
- 所有标志属性集中在一个结构体中
- 便于从配置文件加载 (JSON/YAML) 
- 便于代码生成和文档生成

### 3. 扩展性优势
- 可以轻松添加新的标志类型
- 验证器系统灵活可扩展
- 枚举类型支持提供更好的类型安全

### 4. 类型安全优势
- 编译时检查标志类型
- 避免运行时类型错误
- 明确的类型定义和验证

## 实现考虑

### 1. 默认值处理
不同类型的标志需要不同类型的默认值, 使用 `any` 类型, 然后在创建时进行类型断言。

### 2. 验证器集成
通过标志的 `SetValidator` 方法设置验证器, 与现有实现保持一致。

### 3. 与现有代码兼容
- 保持现有的标志工厂函数
- 新增 `NewFlagFromSpec` 函数
- 两种方式可以并存

### 4. 标志创建方式
将 `NewFlagFromSpec` 作为 `Cmd` 的方法, 确保标志正确注册到当前命令, 避免归属混乱问题。

## 总结

FlagSpec 设计方案完全可行, 而且有很多优势: 

1. **设计一致性** - 与 CmdSpec 形成完整的设计体系
2. **配置简洁性** - 结构化的配置方式, 易于理解和使用
3. **验证集成** - 内置验证器支持
4. **扩展性强** - 易于添加新类型
5. **类型安全** - 编译时检查和类型断言
6. **配置文件支持** - 可以从 JSON/YAML 加载标志定义

这样的设计可以让整个框架更加一致和强大, 同时保持向后兼容性。通过 FlagSpec, 开发者可以更直观地定义标志, 同时获得更好的类型安全和验证支持。