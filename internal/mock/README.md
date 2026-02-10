# Mock 包使用指南

`internal/mock` 包提供了用于测试的模拟实现, 包括模拟命令、标志、注册表和解析器等。这些模拟实现可以帮助您编写单元测试, 而无需依赖真实的实现。

## 主要组件

### 1. MockFlag

模拟标志实现, 支持多种标志类型: 

```go
// 创建基础模拟标志
flag := mock.NewMockFlagBasic("name", "description")

// 创建布尔标志
boolFlag := mock.NewMockBoolFlag("debug", "d", "Enable debug mode", false)

// 创建枚举标志
enumFlag := mock.NewMockEnumFlag("mode", "m", "Operation mode", "normal", []string{"normal", "debug", "release"})
```

### 2. MockCommand / MockCommandBasic

模拟命令实现: 

```go
// 创建基础模拟命令
cmd := mock.NewMockCommandBasic("test", "t", "Test command")

// 创建扩展的模拟命令
extendedCmd := mock.NewMockCommand("test", "t", "Test command")

// 创建子命令
subCmd := mock.NewMockSubCommand("sub", "s", "Sub command", parentCmd)
```

### 3. MockFlagRegistry / MockCmdRegistry

模拟注册表实现: 

```go
// 创建标志注册表
flagReg := mock.NewMockFlagRegistry()

// 创建命令注册表
cmdReg := mock.NewMockCmdRegistry()

// 注册标志和命令
flagReg.Register(flag)
cmdReg.Register(cmd)
```

### 4. MockParser

模拟解析器实现: 

```go
// 创建模拟解析器
parser := mock.NewMockParser()

// 创建带有错误的模拟解析器
errorParser := mock.NewMockParserWithError(parseError, routeError)
```

### 5. TestHelper

测试辅助工具, 提供便捷的方法创建测试对象: 

```go
helper := mock.NewTestHelper()

// 创建带有标志的命令
cmd := helper.CreateMockCommandWithFlags("test", "t", "Test command", flag1, flag2)

// 创建命令树
tree := helper.CreateMockCommandTree()

// 创建带有验证器的标志
validatedFlag := helper.CreateMockFlagWithValidator("token", "t", "Auth token", validator)
```

### 6. CustomHandler

自定义处理器, 用于测试内置标志处理器注册: 

```go
// 创建自定义处理器
customHandler := mock.NewCustomHandler(types.BuiltinFlagType(999))

// 注册到内置标志管理器
manager.RegisterHandler(customHandler)
```

## 使用示例

### 基本使用

```go
func TestMyFunction(t *testing.T) {
    helper := mock.NewTestHelper()
    
    // 创建模拟命令
    cmd := helper.CreateMockCommandWithFlags(
        "test",
        "t",
        "Test command",
        helper.CreateMockBoolFlag("verbose", "v", "Verbose output", false),
        helper.CreateMockEnumFlag("mode", "m", "Operation mode", "normal", []string{"normal", "debug"}),
    )
    
    // 测试代码
    result := MyFunction(cmd)
    
    // 验证结果
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### 创建命令树

```go
func TestCommandTree(t *testing.T) {
    helper := mock.NewTestHelper()
    
    // 创建命令树
    root := helper.CreateMockCommandTree()
    
    // 获取子命令
    sub1, exists := root.GetSubCmd("sub1")
    if !exists {
        t.Error("Sub command sub1 should exist")
    }
    
    // 获取孙子命令
    sub1_1, exists := sub1.GetSubCmd("sub1-1")
    if !exists {
        t.Error("Sub command sub1-1 should exist")
    }
}
```

### 测试解析器

```go
func TestParser(t *testing.T) {
    helper := mock.NewTestHelper()
    
    // 创建模拟解析器
    parser := mock.NewMockParser()
    
    // 创建命令
    cmd := helper.CreateMockCommandWithRunFunc(
        "test",
        "t",
        "Test command",
        func(c types.Command) error {
            // 验证解析结果
            return nil
        },
    )
    
    // 测试解析
    err := parser.ParseAndRoute(cmd, []string{"-v", "--mode", "debug"})
    if err != nil {
        t.Errorf("ParseAndRoute error = %v", err)
    }
}
```

### 测试自定义处理器

```go
func TestCustomHandler(t *testing.T) {
    manager := builtin.NewBuiltinFlagManager()
    
    // 创建自定义处理器
    customHandler := mock.NewCustomHandler(types.BuiltinFlagType(999))
    
    // 注册处理器
    manager.RegisterHandler(customHandler)
    
    // 验证处理器已注册
    if _, exists := manager.GetHandlers()[types.BuiltinFlagType(999)]; !exists {
        t.Error("Custom handler should be registered")
    }
}
```

## 设计原则

1. **简单易用**: 提供便捷的构造函数和辅助方法
2. **灵活可扩展**: 支持自定义行为和状态
3. **完整实现**: 实现了所有必要的接口方法
4. **测试友好**: 提供验证和检查方法

## 注意事项

1. 模拟实现主要用于测试, 不建议在生产代码中使用
2. 某些复杂行为可能需要自定义实现
3. 测试时注意验证模拟对象的状态和方法调用

## 更多示例

更多使用示例请参考 `example_test.go` 文件。