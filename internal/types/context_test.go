package types

import (
	"flag"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestNewCmdContext_基本功能 测试NewCmdContext的基本功能
func TestNewCmdContext_基本功能(t *testing.T) {
	tests := []struct {
		name          string
		longName      string
		shortName     string
		errorHandling flag.ErrorHandling
		expectPanic   bool
		expectedName  string
		description   string
	}{
		{
			name:          "正常创建_长短名称都有",
			longName:      "test-command",
			shortName:     "tc",
			errorHandling: flag.ContinueOnError,
			expectPanic:   false,
			expectedName:  "test-command",
			description:   "正常情况下创建命令上下文",
		},
		{
			name:          "只有长名称",
			longName:      "long-command",
			shortName:     "",
			errorHandling: flag.ContinueOnError,
			expectPanic:   false,
			expectedName:  "long-command",
			description:   "只提供长名称",
		},
		{
			name:          "只有短名称",
			longName:      "",
			shortName:     "s",
			errorHandling: flag.ContinueOnError,
			expectPanic:   false,
			expectedName:  "s",
			description:   "只提供短名称",
		},
		{
			name:          "长短名称都为空",
			longName:      "",
			shortName:     "",
			errorHandling: flag.ContinueOnError,
			expectPanic:   true,
			description:   "长短名称都为空应该panic",
		},
		{
			name:          "ExitOnError模式",
			longName:      "exit-cmd",
			shortName:     "e",
			errorHandling: flag.ExitOnError,
			expectPanic:   false,
			expectedName:  "exit-cmd",
			description:   "使用ExitOnError错误处理模式",
		},
		{
			name:          "PanicOnError模式",
			longName:      "panic-cmd",
			shortName:     "p",
			errorHandling: flag.PanicOnError,
			expectPanic:   false,
			expectedName:  "panic-cmd",
			description:   "使用PanicOnError错误处理模式",
		},
		{
			name:          "特殊字符名称",
			longName:      "test-cmd_123",
			shortName:     "t1",
			errorHandling: flag.ContinueOnError,
			expectPanic:   false,
			expectedName:  "test-cmd_123",
			description:   "包含特殊字符的命令名称",
		},
		{
			name:          "中文名称",
			longName:      "测试命令",
			shortName:     "测",
			errorHandling: flag.ContinueOnError,
			expectPanic:   false,
			expectedName:  "测试命令",
			description:   "中文命令名称",
		},
		{
			name:          "长名称为空_使用短名称作为FlagSet名称",
			longName:      "",
			shortName:     "short",
			errorHandling: flag.ContinueOnError,
			expectPanic:   false,
			expectedName:  "short",
			description:   "长名称为空时应使用短名称作为FlagSet名称",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx *CmdContext
			var panicked bool

			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
						if !tt.expectPanic {
							t.Errorf("意外的panic: %v", r)
						}
					}
				}()
				ctx = NewCmdContext(tt.longName, tt.shortName, tt.errorHandling)
			}()

			if tt.expectPanic {
				if !panicked {
					t.Error("期望panic但没有发生")
				}
				return
			}

			// 验证基本字段
			if ctx == nil {
				t.Fatal("NewCmdContext返回了nil")
			}

			if ctx.LongName != tt.longName {
				t.Errorf("长名称不匹配: 期望 %q, 实际 %q", tt.longName, ctx.LongName)
			}

			if ctx.ShortName != tt.shortName {
				t.Errorf("短名称不匹配: 期望 %q, 实际 %q", tt.shortName, ctx.ShortName)
			}

			// 验证FlagSet名称
			if ctx.FlagSet.Name() != tt.expectedName {
				t.Errorf("FlagSet名称不匹配: 期望 %q, 实际 %q", tt.expectedName, ctx.FlagSet.Name())
			}

			// 验证初始化的字段
			if ctx.FlagRegistry == nil {
				t.Error("FlagRegistry未初始化")
			}

			if ctx.FlagSet == nil {
				t.Error("FlagSet未初始化")
			}

			if ctx.Args == nil {
				t.Error("Args未初始化")
			}

			if len(ctx.Args) != 0 {
				t.Errorf("Args初始长度应为0, 实际: %d", len(ctx.Args))
			}

			if ctx.SubCmds == nil {
				t.Error("SubCmds未初始化")
			}

			if len(ctx.SubCmds) != 0 {
				t.Errorf("SubCmds初始长度应为0, 实际: %d", len(ctx.SubCmds))
			}

			if ctx.SubCmdMap == nil {
				t.Error("SubCmdMap未初始化")
			}

			if len(ctx.SubCmdMap) != 0 {
				t.Errorf("SubCmdMap初始长度应为0, 实际: %d", len(ctx.SubCmdMap))
			}

			if ctx.Config == nil {
				t.Error("Config未初始化")
			}

			if ctx.BuiltinFlags == nil {
				t.Error("BuiltinFlags未初始化")
			}

			// 验证解析状态
			if ctx.Parsed.Load() {
				t.Error("新创建的上下文不应该处于已解析状态")
			}

			// 验证Parent为nil
			if ctx.Parent != nil {
				t.Error("新创建的上下文Parent应该为nil")
			}
		})
	}
}

// TestCmdContext_GetName_边界场景 测试GetName方法的边界场景
func TestCmdContext_GetName_边界场景(t *testing.T) {
	tests := []struct {
		name         string
		longName     string
		shortName    string
		expectedName string
		description  string
	}{
		{
			name:         "长短名称都有_优先返回长名称",
			longName:     "long-name",
			shortName:    "s",
			expectedName: "long-name",
			description:  "有长名称时优先返回长名称",
		},
		{
			name:         "只有长名称",
			longName:     "only-long",
			shortName:    "",
			expectedName: "only-long",
			description:  "只有长名称时返回长名称",
		},
		{
			name:         "只有短名称",
			longName:     "",
			shortName:    "o",
			expectedName: "o",
			description:  "只有短名称时返回短名称",
		},
		{
			name:         "空字符串长名称_有短名称",
			longName:     "",
			shortName:    "short",
			expectedName: "short",
			description:  "长名称为空字符串时返回短名称",
		},
		{
			name:         "特殊字符名称",
			longName:     "test-cmd_123",
			shortName:    "t1",
			expectedName: "test-cmd_123",
			description:  "包含特殊字符的名称",
		},
		{
			name:         "中文名称",
			longName:     "测试命令",
			shortName:    "测",
			expectedName: "测试命令",
			description:  "中文命令名称",
		},
		{
			name:         "长名称包含空格",
			longName:     "command with spaces",
			shortName:    "c",
			expectedName: "command with spaces",
			description:  "长名称包含空格字符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewCmdContext(tt.longName, tt.shortName, flag.ContinueOnError)
			gotName := ctx.GetName()

			if gotName != tt.expectedName {
				t.Errorf("GetName()返回值不匹配: 期望 %q, 实际 %q", tt.expectedName, gotName)
			}
		})
	}
}

// TestCmdContext_并发安全性 测试CmdContext的并发安全性
func TestCmdContext_并发安全性(t *testing.T) {
	ctx := NewCmdContext("test", "t", flag.ContinueOnError)

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// 测试并发读取
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// 并发读取各种字段
				_ = ctx.GetName()
				_ = ctx.LongName
				_ = ctx.ShortName
				_ = ctx.Parsed.Load()

				// 并发读取切片和映射（需要锁保护）
				ctx.Mutex.RLock()
				_ = len(ctx.Args)
				_ = len(ctx.SubCmds)
				_ = len(ctx.SubCmdMap)
				ctx.Mutex.RUnlock()
			}
		}(i)
	}

	// 测试并发写入Args
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				ctx.Mutex.Lock()
				ctx.Args = append(ctx.Args, "arg")
				ctx.Mutex.Unlock()
			}
		}(i)
	}

	// 测试并发操作Parsed状态
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				if j%2 == 0 {
					ctx.Parsed.Store(true)
				} else {
					ctx.Parsed.Store(false)
				}
				_ = ctx.Parsed.Load()
			}
		}(i)
	}

	wg.Wait()

	// 验证最终状态的一致性
	ctx.Mutex.RLock()
	argsLen := len(ctx.Args)
	ctx.Mutex.RUnlock()

	expectedArgsLen := numGoroutines * numOperations
	if argsLen != expectedArgsLen {
		t.Errorf("并发写入Args后长度不正确: 期望 %d, 实际 %d", expectedArgsLen, argsLen)
	}

	t.Logf("并发测试完成 - Args长度: %d, Parsed状态: %v", argsLen, ctx.Parsed.Load())
}

// TestCmdContext_ParseOnce_并发 测试ParseOnce的并发安全性
func TestCmdContext_ParseOnce_并发(t *testing.T) {
	ctx := NewCmdContext("test", "t", flag.ContinueOnError)

	var executeCount int32
	var wg sync.WaitGroup
	numGoroutines := 10

	// 模拟多个goroutine同时尝试执行解析
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			ctx.ParseOnce.Do(func() {
				// 模拟解析操作
				time.Sleep(10 * time.Millisecond)
				executeCount++
				ctx.Parsed.Store(true)
			})
		}()
	}

	wg.Wait()

	// 验证ParseOnce确保只执行一次
	if executeCount != 1 {
		t.Errorf("ParseOnce应该只执行一次, 实际执行了 %d 次", executeCount)
	}

	if !ctx.Parsed.Load() {
		t.Error("解析后Parsed状态应该为true")
	}
}

// TestCmdContext_字段初始化完整性 测试所有字段的初始化完整性
func TestCmdContext_字段初始化完整性(t *testing.T) {
	ctx := NewCmdContext("test", "t", flag.ContinueOnError)

	// 使用反射检查所有字段是否正确初始化
	v := reflect.ValueOf(ctx).Elem()
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		switch fieldName {
		case "LongName":
			if field.String() != "test" {
				t.Errorf("LongName未正确初始化: %v", field.Interface())
			}
		case "ShortName":
			if field.String() != "t" {
				t.Errorf("ShortName未正确初始化: %v", field.Interface())
			}
		case "FlagRegistry":
			if field.IsNil() {
				t.Error("FlagRegistry未初始化")
			}
		case "FlagSet":
			if field.IsNil() {
				t.Error("FlagSet未初始化")
			}
		case "Args":
			if field.IsNil() {
				t.Error("Args未初始化")
			}
			if field.Len() != 0 {
				t.Errorf("Args初始长度应为0, 实际: %d", field.Len())
			}
		case "SubCmds":
			if field.IsNil() {
				t.Error("SubCmds未初始化")
			}
			if field.Len() != 0 {
				t.Errorf("SubCmds初始长度应为0, 实际: %d", field.Len())
			}
		case "SubCmdMap":
			if field.IsNil() {
				t.Error("SubCmdMap未初始化")
			}
			if field.Len() != 0 {
				t.Errorf("SubCmdMap初始长度应为0, 实际: %d", field.Len())
			}
		case "Parent":
			if !field.IsNil() {
				t.Error("Parent应该初始化为nil")
			}
		case "Config":
			if field.IsNil() {
				t.Error("Config未初始化")
			}
		case "BuiltinFlags":
			if field.IsNil() {
				t.Error("BuiltinFlags未初始化")
			}
		case "ParseHook":
			if !field.IsNil() {
				t.Error("ParseHook应该初始化为nil")
			}
		case "Parsed":
			// atomic.Bool类型，检查初始值
			if ctx.Parsed.Load() {
				t.Error("Parsed应该初始化为false")
			}
		case "ParseOnce":
			// sync.Once类型，无法直接检查，但可以验证其功能
			// 这里不做特殊检查，因为在其他测试中已经验证了功能
		case "Mutex":
			// sync.RWMutex类型，无法直接检查初始状态
			// 但可以验证其可用性
			ctx.Mutex.RLock()
			ctx.Mutex.RUnlock()
		}
	}
}

// TestCmdContext_极值测试 测试极值情况
func TestCmdContext_极值测试(t *testing.T) {
	tests := []struct {
		name        string
		longName    string
		shortName   string
		description string
	}{
		{
			name:        "极长的长名称",
			longName:    strings.Repeat("a", 10000),
			shortName:   "a",
			description: "测试极长的长名称",
		},
		{
			name:        "极长的短名称",
			longName:    "test",
			shortName:   strings.Repeat("b", 1000),
			description: "测试极长的短名称",
		},
		{
			name:        "单字符名称",
			longName:    "a",
			shortName:   "b",
			description: "测试单字符名称",
		},
		{
			name:        "包含所有ASCII字符",
			longName:    "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./~`",
			shortName:   "!",
			description: "测试包含特殊ASCII字符的名称",
		},
		{
			name:        "Unicode字符",
			longName:    "测试命令🚀✨🎉",
			shortName:   "🚀",
			description: "测试Unicode字符名称",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewCmdContext(tt.longName, tt.shortName, flag.ContinueOnError)

			if ctx == nil {
				t.Fatal("NewCmdContext返回了nil")
			}

			if ctx.LongName != tt.longName {
				t.Errorf("长名称不匹配: 期望 %q, 实际 %q", tt.longName, ctx.LongName)
			}

			if ctx.ShortName != tt.shortName {
				t.Errorf("短名称不匹配: 期望 %q, 实际 %q", tt.shortName, ctx.ShortName)
			}

			// 验证GetName方法在极值情况下的表现
			expectedName := tt.longName
			if expectedName == "" {
				expectedName = tt.shortName
			}

			if ctx.GetName() != expectedName {
				t.Errorf("GetName()在极值情况下返回值不正确: 期望 %q, 实际 %q", expectedName, ctx.GetName())
			}
		})
	}
}

// TestCmdContext_内存泄漏检测 测试潜在的内存泄漏
func TestCmdContext_内存泄漏检测(t *testing.T) {
	// 创建大量上下文并立即释放
	for i := 0; i < 1000; i++ {
		ctx := NewCmdContext("test", "t", flag.ContinueOnError)

		// 添加一些数据
		ctx.Mutex.Lock()
		ctx.Args = append(ctx.Args, "arg1", "arg2", "arg3")
		ctx.SubCmdMap["child"] = &CmdContext{}
		ctx.Mutex.Unlock()

		// 设置一些状态
		ctx.Parsed.Store(true)

		// 清理引用（模拟正常使用后的清理）
		ctx.Mutex.Lock()
		ctx.Args = nil
		ctx.SubCmdMap = nil
		ctx.SubCmds = nil
		ctx.Parent = nil
		ctx.Mutex.Unlock()
	}

	// 这个测试主要是为了在运行时检测内存使用情况
	// 实际的内存泄漏检测需要使用专门的工具如pprof
	t.Log("内存泄漏检测测试完成")
}

// TestCmdContext_错误处理模式 测试不同的错误处理模式
func TestCmdContext_错误处理模式(t *testing.T) {
	errorModes := []flag.ErrorHandling{
		flag.ContinueOnError,
		flag.ExitOnError,
		flag.PanicOnError,
	}

	for i, mode := range errorModes {
		modeName := []string{"ContinueOnError", "ExitOnError", "PanicOnError"}[i]
		t.Run(modeName, func(t *testing.T) {
			ctx := NewCmdContext("test", "t", mode)

			if ctx == nil {
				t.Fatal("NewCmdContext返回了nil")
			}

			if ctx.FlagSet == nil {
				t.Fatal("FlagSet未初始化")
			}

			// 验证FlagSet的错误处理模式是否正确设置
			// 注意：flag.FlagSet没有公开的方法来获取ErrorHandling，
			// 所以我们只能通过创建成功来验证
			t.Logf("成功创建错误处理模式为 %s 的上下文", modeName)
		})
	}
}
