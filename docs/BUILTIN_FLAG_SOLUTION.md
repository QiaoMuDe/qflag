# 内置标志名称冲突解决方案 - 方案三：延迟注册

## 问题背景

### 当前问题

在当前的实现中，内置标志管理器在初始化时会预注册所有内置标志的名称到全局映射表中：

```go
// builtin/manager.go
func (m *BuiltinFlagManager) RegisterHandler(handler types.BuiltinFlagHandler) {
    flagType := handler.Type()
    m.handlers[flagType] = handler

    // 注册标志名映射 - 问题所在！
    switch flagType {
    case types.HelpFlag:
        m.flags[types.HelpFlagName] = types.HelpFlag
        m.flags[types.HelpFlagShortName] = types.HelpFlag
    case types.VersionFlag:
        m.flags[types.VersionFlagName] = types.VersionFlag      // ❌ 全局占用 "version"
        m.flags[types.VersionFlagShortName] = types.VersionFlag  // ❌ 全局占用 "v"
    case types.CompletionFlag:
        m.flags[types.CompletionFlagName] = types.CompletionFlag  // ❌ 全局占用 "completion"
    }
}
```

### 问题场景

```go
// 主命令
mainCmd := qflag.NewCmd("main", "m", qflag.ExitOnError)
mainCmd.SetVersion("1.0.0")  // 会注册版本标志

// 子命令
subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
subCmd.Bool("verbose", "v", "详细输出", false)  // ❌ 冲突！无法使用 -v

// 原因：即使子命令不会注册版本标志，"v" 也被全局映射认为是内置标志
```

---

## 方案三：延迟注册

### 核心思路

**不在初始化时预注册所有内置标志名称，而是在实际检查时动态判断是否为内置标志**

### 设计原则

1. **动态判断**：每次检查时都基于当前命令的配置
2. **延迟注册**：不在初始化时预注册名称
3. **上下文相关**：检查结果依赖于具体的命令实例
4. **向后兼容**：不改变现有的 `ShouldRegister` 逻辑

---

## 实现方案

### 1. 修改 `BuiltinFlagManager` 结构体

```go
type BuiltinFlagManager struct {
    handlers map[types.BuiltinFlagType]types.BuiltinFlagHandler
    // 移除 flags 字段，不再预注册所有标志名称
}
```

### 2. 修改 `RegisterHandler` 方法

```go
// RegisterHandler 注册内置标志处理器
//
// 参数:
//   - handler: 要注册的处理器
//
// 功能说明:
//   - 将处理器添加到处理器映射表
//   - 不再预注册标志名映射
func (m *BuiltinFlagManager) RegisterHandler(handler types.BuiltinFlagHandler) {
    flagType := handler.Type()
    m.handlers[flagType] = handler
    
    // 移除标志名映射的注册
    // 不再预注册所有标志名称到全局映射表
}
```

### 3. 修改 `isBuiltinFlag` 方法

```go
// isBuiltinFlag 检查是否是内置标志
//
// 参数:
//   - f: 要检查的标志
//   - cmd: 当前命令实例
//
// 返回值:
//   - types.BuiltinFlagType: 标志类型
//   - bool: 是否是内置标志
//
// 功能说明:
//   - 动态检查标志是否为内置标志
//   - 基于当前命令的配置和 ShouldRegister 判断
//   - 解决名称冲突问题
func (m *BuiltinFlagManager) isBuiltinFlag(f types.Flag, cmd types.Command) (types.BuiltinFlagType, bool) {
    // 遍历所有处理器，检查是否应该注册且名称匹配
    for _, handler := range m.handlers {
        // 检查该标志类型是否真的会在当前命令中注册
        if !handler.ShouldRegister(cmd) {
            continue
        }

        // 检查名称是否匹配
        switch handler.Type() {
        case types.HelpFlag:
            if f.LongName() == types.HelpFlagName || f.ShortName() == types.HelpFlagShortName {
                return types.HelpFlag, true
            }
            
        case types.VersionFlag:
            if f.LongName() == types.VersionFlagName || f.ShortName() == types.VersionFlagShortName {
                return types.VersionFlag, true
            }
            
        case types.CompletionFlag:
            if f.LongName() == types.CompletionFlagName {
                return types.CompletionFlag, true
            }
        }
    }

    return 0, false
}
```

### 4. 修改 `HandleBuiltinFlags` 方法

```go
// HandleBuiltinFlags 处理内置标志
//
// 参数:
//   - cmd: 要处理标志的命令
//
// 返回值:
//   - error: 处理失败时返回错误
//
// 功能说明:
//   - 遍历命令的所有标志, 检查是否是内置标志
//   - 如果是内置标志且被设置, 则执行对应的处理器
//   - 传入当前命令进行动态检查
func (m *BuiltinFlagManager) HandleBuiltinFlags(cmd types.Command) error {
    flags := cmd.Flags()

    for _, f := range flags {
        // 传入当前命令进行动态检查
        if flagType, isBuiltin := m.isBuiltinFlag(f, cmd); isBuiltin {
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
```

---

## 方案优势

| 优势 | 说明 |
|------|------|
| ✅ **解决名称冲突** | 子命令可以自由使用 `-v` 等标志名称 |
| ✅ **逻辑清晰** | 基于实际的 `ShouldRegister` 判断 |
| ✅ **动态判断** | 每次检查时都基于当前命令的配置 |
| ✅ **简单实现** | 不需要复杂的映射管理 |
| ✅ **向后兼容** | 不改变现有的 `ShouldRegister` 逻辑 |
| ✅ **性能良好** | 只在实际需要时进行检查 |

---

## 实际使用场景

### 场景1：子命令使用 `-v` 标志

```go
// 主命令
mainCmd := qflag.NewCmd("main", "m", qflag.ExitOnError)
mainCmd.SetVersion("1.0.0")  // 设置版本信息

// 子命令
subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
subCmd.Bool("verbose", "v", "详细输出", false)  // ✅ 可以使用 -v

// 注册
mainCmd.AddSubCmds(subCmd)
```

### 场景2：检查逻辑

```go
// 主命令的 -v 会被认为是内置版本标志
isBuiltin, flagType := builtinMgr.isBuiltinFlag(versionFlag, mainCmd)
// isBuiltin = true, flagType = VersionFlag
// 原因：mainCmd.IsRootCmd() = true 且 mainCmd.Config().Version != ""

// 子命令的 -v 不会被认为是内置版本标志
isBuiltin, flagType := builtinMgr.isBuiltinFlag(verboseFlag, subCmd)
// isBuiltin = false, flagType = 0
// 原因：subCmd.IsRootCmd() = false，VersionHandler.ShouldRegister() 返回 false
```

### 场景3：多层子命令

```go
// 主命令
rootCmd := qflag.NewCmd("app", "a", qflag.ExitOnError)
rootCmd.SetVersion("1.0.0")

// 子命令
subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
subCmd.Bool("verbose", "v", "详细输出", false)  // ✅ 可以使用 -v

// 子子命令
subSubCmd := qflag.NewCmd("subsub", "ss", qflag.ExitOnError)
subSubCmd.Bool("value", "v", "值", false)  // ✅ 可以使用 -v

rootCmd.AddSubCmds(subCmd)
subCmd.AddSubCmds(subSubCmd)
```

---

## 实现步骤

### 第一步：修改结构体

1. 移除 `BuiltinFlagManager` 中的 `flags` 字段
2. 更新 `NewBuiltinFlagManager` 初始化方法

### 第二步：修改注册方法

1. 移除 `RegisterHandler` 中的标志名映射代码
2. 保留处理器映射逻辑

### 第三步：修改检查方法

1. 重写 `isBuiltinFlag` 方法，添加 `cmd` 参数
2. 实现动态检查逻辑
3. 基于 `ShouldRegister` 进行判断

### 第四步：修改处理方法

1. 更新 `HandleBuiltinFlags` 调用
2. 传入当前命令参数

### 第五步：更新调用点

1. 找到所有调用 `isBuiltinFlag` 的地方
2. 添加当前命令参数

---

## 测试验证

### 测试用例1：名称冲突

```go
func TestBuiltinFlagNameConflict(t *testing.T) {
    // 创建主命令（有版本）
    mainCmd := qflag.NewCmd("main", "m", qflag.ExitOnError)
    mainCmd.SetVersion("1.0.0")
    
    // 创建子命令
    subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
    verboseFlag := subCmd.Bool("verbose", "v", "详细输出", false)
    
    // 验证：子命令的 -v 不被认为是内置标志
    builtinMgr := builtin.NewBuiltinFlagManager()
    flagType, isBuiltin := builtinMgr.isBuiltinFlag(verboseFlag, subCmd)
    
    if isBuiltin {
        t.Error("子命令的 -v 不应该被认为是内置标志")
    }
    
    if flagType != 0 {
        t.Errorf("标志类型应该为 0，得到 %d", flagType)
    }
}
```

### 测试用例2：多层嵌套

```go
func TestMultiLevelBuiltinFlags(t *testing.T) {
    // 创建三层命令
    rootCmd := qflag.NewCmd("root", "r", qflag.ExitOnError)
    rootCmd.SetVersion("1.0.0")
    
    subCmd := qflag.NewCmd("sub", "s", qflag.ExitOnError)
    subCmd.Bool("verbose", "v", "详细输出", false)
    
    subSubCmd := qflag.NewCmd("subsub", "ss", qflag.ExitOnError)
    subSubCmd.Bool("value", "v", "值", false)
    
    // 验证：只有根命令的 -v 是内置标志
    builtinMgr := builtin.NewBuiltinFlagManager()
    
    // 根命令
    rootVersionFlag := rootCmd.Bool("version", "v", "版本", false)
    flagType, isBuiltin := builtinMgr.isBuiltinFlag(rootVersionFlag, rootCmd)
    if !isBuiltin || flagType != types.VersionFlag {
        t.Error("根命令的 -v 应该被认为是内置版本标志")
    }
    
    // 子命令
    flagType, isBuiltin = builtinMgr.isBuiltinFlag(verboseFlag, subCmd)
    if isBuiltin {
        t.Error("子命令的 -v 不应该被认为是内置标志")
    }
    
    // 子子命令
    valueFlag := subSubCmd.Bool("value", "v", "值", false)
    flagType, isBuiltin = builtinMgr.isBuiltinFlag(valueFlag, subSubCmd)
    if isBuiltin {
        t.Error("子子命令的 -v 不应该被认为是内置标志")
    }
}
```

---

## 总结

方案三通过**延迟注册和动态检查**，完美解决了内置标志名称占用的问题：

1. ✅ **解决核心问题**：子命令可以自由使用 `-v` 等标志名称
2. ✅ **保持设计一致性**：基于现有的 `ShouldRegister` 逻辑
3. ✅ **实现简单清晰**：不需要复杂的映射管理
4. ✅ **性能良好**：只在需要时进行检查
5. ✅ **完全向后兼容**：不改变现有接口

这个方案是最优雅的解决方案，既解决了实际问题，又保持了代码的简洁性。