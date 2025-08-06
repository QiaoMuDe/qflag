package cmd

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// createTestCmd 创建用于测试的命令实例
func createTestCmd() *Cmd {
	return NewCmd("test-cmd", "tc", flag.ContinueOnError)
}

// createTestCmdWithBuiltins 创建带有内置标志的测试命令
func createTestCmdWithBuiltins() *Cmd {
	cmd := createTestCmd()
	// 标记内置标志
	cmd.ctx.BuiltinFlags.MarkAsBuiltin("help", "h", "version", "v")
	return cmd
}

// =============================================================================
// 布尔类型标志测试
// =============================================================================

// TestBoolVar 测试布尔类型标志变量绑定
func TestBoolVar(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  bool
		usage     string
		wantPanic bool
		panicMsg  string
		setupCmd  func() *Cmd
	}{
		{
			name:      "有效的长标志名",
			longName:  "verbose",
			shortName: "",
			defValue:  false,
			usage:     "启用详细输出",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "有效的短标志名",
			longName:  "",
			shortName: "v",
			defValue:  true,
			usage:     "启用详细输出",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "长短标志名都有效",
			longName:  "debug",
			shortName: "d",
			defValue:  false,
			usage:     "启用调试模式",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "nil标志指针",
			longName:  "test",
			shortName: "",
			defValue:  false,
			usage:     "测试标志",
			wantPanic: true,
			panicMsg:  "BoolFlag pointer cannot be nil",
			setupCmd:  createTestCmd,
		},
		{
			name:      "内置长标志名冲突",
			longName:  "help",
			shortName: "",
			defValue:  false,
			usage:     "帮助信息",
			wantPanic: true,
			panicMsg:  "flag long name help is reserved",
			setupCmd:  createTestCmdWithBuiltins,
		},
		{
			name:      "内置短标志名冲突",
			longName:  "",
			shortName: "h",
			defValue:  false,
			usage:     "帮助信息",
			wantPanic: true,
			panicMsg:  "flag short name h is reserved",
			setupCmd:  createTestCmdWithBuiltins,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("BoolVar() 期望panic但未发生")
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("BoolVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				}()

				var f *flags.BoolFlag
				if tt.name != "nil标志指针" {
					f = &flags.BoolFlag{}
				}
				cmd.BoolVar(f, tt.longName, tt.shortName, tt.defValue, tt.usage)
			} else {
				f := &flags.BoolFlag{}
				cmd.BoolVar(f, tt.longName, tt.shortName, tt.defValue, tt.usage)

				// 验证标志是否正确初始化
				if f.LongName() != tt.longName {
					t.Errorf("BoolVar() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
				}
				if f.ShortName() != tt.shortName {
					t.Errorf("BoolVar() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
				}
				if f.Get() != tt.defValue {
					t.Errorf("BoolVar() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
				}
				if f.Usage() != tt.usage {
					t.Errorf("BoolVar() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
				}
			}
		})
	}
}

// TestBool 测试布尔类型标志创建
func TestBool(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  bool
		usage     string
	}{
		{
			name:      "创建布尔标志",
			longName:  "verbose",
			shortName: "v",
			defValue:  false,
			usage:     "启用详细输出",
		},
		{
			name:      "创建默认为true的布尔标志",
			longName:  "quiet",
			shortName: "q",
			defValue:  true,
			usage:     "静默模式",
		},
		{
			name:      "只有长标志名",
			longName:  "debug",
			shortName: "",
			defValue:  false,
			usage:     "调试模式",
		},
		{
			name:      "只有短标志名",
			longName:  "",
			shortName: "f",
			defValue:  true,
			usage:     "强制模式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCmd()
			f := cmd.Bool(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if f == nil {
				t.Errorf("Bool() 返回nil指针")
				return
			}

			// 验证标志属性
			if f.LongName() != tt.longName {
				t.Errorf("Bool() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
			}
			if f.ShortName() != tt.shortName {
				t.Errorf("Bool() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
			}
			if f.Get() != tt.defValue {
				t.Errorf("Bool() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
			}
			if f.Usage() != tt.usage {
				t.Errorf("Bool() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
			}
		})
	}
}

// =============================================================================
// 字符串类型标志测试
// =============================================================================

// TestStringVar 测试字符串类型标志变量绑定
func TestStringVar(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  string
		usage     string
		wantPanic bool
		panicMsg  string
		setupCmd  func() *Cmd
	}{
		{
			name:      "有效的字符串标志",
			longName:  "output",
			shortName: "o",
			defValue:  "stdout",
			usage:     "输出目标",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "空默认值",
			longName:  "config",
			shortName: "c",
			defValue:  "",
			usage:     "配置文件路径",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "中文默认值",
			longName:  "message",
			shortName: "m",
			defValue:  "你好世界",
			usage:     "消息内容",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "nil标志指针",
			longName:  "test",
			shortName: "",
			defValue:  "default",
			usage:     "测试标志",
			wantPanic: true,
			panicMsg:  "StringFlag pointer cannot be nil",
			setupCmd:  createTestCmd,
		},
		{
			name:      "内置标志冲突",
			longName:  "version",
			shortName: "",
			defValue:  "1.0.0",
			usage:     "版本信息",
			wantPanic: true,
			panicMsg:  "flag long name version is reserved",
			setupCmd:  createTestCmdWithBuiltins,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("StringVar() 期望panic但未发生")
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("StringVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				}()

				var f *flags.StringFlag
				if tt.name != "nil标志指针" {
					f = &flags.StringFlag{}
				}
				cmd.StringVar(f, tt.longName, tt.shortName, tt.defValue, tt.usage)
			} else {
				f := &flags.StringFlag{}
				cmd.StringVar(f, tt.longName, tt.shortName, tt.defValue, tt.usage)

				// 验证标志属性
				if f.LongName() != tt.longName {
					t.Errorf("StringVar() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
				}
				if f.ShortName() != tt.shortName {
					t.Errorf("StringVar() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
				}
				if f.Get() != tt.defValue {
					t.Errorf("StringVar() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
				}
				if f.Usage() != tt.usage {
					t.Errorf("StringVar() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
				}
			}
		})
	}
}

// TestString 测试字符串类型标志创建
func TestString(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  string
		usage     string
	}{
		{
			name:      "创建字符串标志",
			longName:  "output",
			shortName: "o",
			defValue:  "stdout",
			usage:     "输出目标",
		},
		{
			name:      "特殊字符默认值",
			longName:  "pattern",
			shortName: "p",
			defValue:  "*.go",
			usage:     "文件模式",
		},
		{
			name:      "长字符串默认值",
			longName:  "description",
			shortName: "",
			defValue:  "这是一个很长的描述信息，用于测试长字符串的处理能力",
			usage:     "描述信息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCmd()
			f := cmd.String(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if f == nil {
				t.Errorf("String() 返回nil指针")
				return
			}

			// 验证标志属性
			if f.LongName() != tt.longName {
				t.Errorf("String() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
			}
			if f.ShortName() != tt.shortName {
				t.Errorf("String() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
			}
			if f.Get() != tt.defValue {
				t.Errorf("String() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
			}
			if f.Usage() != tt.usage {
				t.Errorf("String() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
			}
		})
	}
}

// =============================================================================
// 浮点数类型标志测试
// =============================================================================

// TestFloat64Var 测试浮点数类型标志变量绑定
func TestFloat64Var(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  float64
		usage     string
		wantPanic bool
		panicMsg  string
		setupCmd  func() *Cmd
	}{
		{
			name:      "有效的浮点数标志",
			longName:  "threshold",
			shortName: "t",
			defValue:  0.5,
			usage:     "阈值设置",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "零值",
			longName:  "zero",
			shortName: "",
			defValue:  0.0,
			usage:     "零值测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "负数值",
			longName:  "negative",
			shortName: "n",
			defValue:  -3.14,
			usage:     "负数测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "极大值",
			longName:  "max",
			shortName: "",
			defValue:  1.7976931348623157e+308,
			usage:     "最大值测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "nil标志指针",
			longName:  "test",
			shortName: "",
			defValue:  1.0,
			usage:     "测试标志",
			wantPanic: true,
			panicMsg:  "FloatFlag pointer cannot be nil",
			setupCmd:  createTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Float64Var() 期望panic但未发生")
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Float64Var() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				}()

				var f *flags.Float64Flag
				if tt.name != "nil标志指针" {
					f = &flags.Float64Flag{}
				}
				cmd.Float64Var(f, tt.longName, tt.shortName, tt.defValue, tt.usage)
			} else {
				f := &flags.Float64Flag{}
				cmd.Float64Var(f, tt.longName, tt.shortName, tt.defValue, tt.usage)

				// 验证标志属性
				if f.LongName() != tt.longName {
					t.Errorf("Float64Var() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
				}
				if f.ShortName() != tt.shortName {
					t.Errorf("Float64Var() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
				}
				if f.Get() != tt.defValue {
					t.Errorf("Float64Var() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
				}
				if f.Usage() != tt.usage {
					t.Errorf("Float64Var() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
				}
			}
		})
	}
}

// TestFloat64 测试浮点数类型标志创建
func TestFloat64(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  float64
		usage     string
	}{
		{
			name:      "创建浮点数标志",
			longName:  "rate",
			shortName: "r",
			defValue:  1.5,
			usage:     "速率设置",
		},
		{
			name:      "科学计数法",
			longName:  "scientific",
			shortName: "",
			defValue:  1.23e-4,
			usage:     "科学计数法测试",
		},
		{
			name:      "π值",
			longName:  "pi",
			shortName: "",
			defValue:  3.141592653589793,
			usage:     "圆周率",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCmd()
			f := cmd.Float64(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if f == nil {
				t.Errorf("Float64() 返回nil指针")
				return
			}

			// 验证标志属性
			if f.LongName() != tt.longName {
				t.Errorf("Float64() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
			}
			if f.ShortName() != tt.shortName {
				t.Errorf("Float64() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
			}
			if f.Get() != tt.defValue {
				t.Errorf("Float64() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
			}
			if f.Usage() != tt.usage {
				t.Errorf("Float64() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
			}
		})
	}
}

// =============================================================================
// 整数类型标志测试
// =============================================================================

// TestIntVar 测试整数类型标志变量绑定
func TestIntVar(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  int
		usage     string
		wantPanic bool
		panicMsg  string
		setupCmd  func() *Cmd
	}{
		{
			name:      "有效的整数标志",
			longName:  "count",
			shortName: "c",
			defValue:  10,
			usage:     "计数器",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "零值",
			longName:  "zero",
			shortName: "",
			defValue:  0,
			usage:     "零值测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "负数值",
			longName:  "negative",
			shortName: "n",
			defValue:  -100,
			usage:     "负数测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "nil标志指针",
			longName:  "test",
			shortName: "",
			defValue:  1,
			usage:     "测试标志",
			wantPanic: true,
			panicMsg:  "IntFlag pointer cannot be nil",
			setupCmd:  createTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("IntVar() 期望panic但未发生")
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("IntVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				}()

				var f *flags.IntFlag
				if tt.name != "nil标志指针" {
					f = &flags.IntFlag{}
				}
				cmd.IntVar(f, tt.longName, tt.shortName, tt.defValue, tt.usage)
			} else {
				f := &flags.IntFlag{}
				cmd.IntVar(f, tt.longName, tt.shortName, tt.defValue, tt.usage)

				// 验证标志属性
				if f.LongName() != tt.longName {
					t.Errorf("IntVar() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
				}
				if f.ShortName() != tt.shortName {
					t.Errorf("IntVar() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
				}
				if f.Get() != tt.defValue {
					t.Errorf("IntVar() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
				}
				if f.Usage() != tt.usage {
					t.Errorf("IntVar() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
				}
			}
		})
	}
}

// TestInt 测试整数类型标志创建
func TestInt(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  int
		usage     string
	}{
		{
			name:      "创建整数标志",
			longName:  "port",
			shortName: "p",
			defValue:  8080,
			usage:     "端口号",
		},
		{
			name:      "最大整数值",
			longName:  "max-int",
			shortName: "",
			defValue:  2147483647,
			usage:     "最大整数值",
		},
		{
			name:      "最小整数值",
			longName:  "min-int",
			shortName: "",
			defValue:  -2147483648,
			usage:     "最小整数值",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCmd()
			f := cmd.Int(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if f == nil {
				t.Errorf("Int() 返回nil指针")
				return
			}

			// 验证标志属性
			if f.LongName() != tt.longName {
				t.Errorf("Int() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
			}
			if f.ShortName() != tt.shortName {
				t.Errorf("Int() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
			}
			if f.Get() != tt.defValue {
				t.Errorf("Int() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
			}
			if f.Usage() != tt.usage {
				t.Errorf("Int() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
			}
		})
	}
}

// =============================================================================
// 64位整数类型标志测试
// =============================================================================

// TestInt64Var 测试64位整数类型标志变量绑定
func TestInt64Var(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  int64
		usage     string
		wantPanic bool
		panicMsg  string
		setupCmd  func() *Cmd
	}{
		{
			name:      "有效的64位整数标志",
			longName:  "size",
			shortName: "s",
			defValue:  1024,
			usage:     "文件大小",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "大数值",
			longName:  "big-number",
			shortName: "",
			defValue:  9223372036854775807, // int64最大值
			usage:     "大数值测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "负大数值",
			longName:  "negative-big",
			shortName: "",
			defValue:  -9223372036854775808, // int64最小值
			usage:     "负大数值测试",
			wantPanic: false,
			setupCmd:  createTestCmd,
		},
		{
			name:      "nil标志指针",
			longName:  "test",
			shortName: "",
			defValue:  1,
			usage:     "测试标志",
			wantPanic: true,
			panicMsg:  "Int64Flag pointer cannot be nil",
			setupCmd:  createTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Int64Var() 期望panic但未发生")
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Int64Var() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				}()

				var f *flags.Int64Flag
				if tt.name != "nil标志指针" {
					f = &flags.Int64Flag{}
				}
				cmd.Int64Var(f, tt.longName, tt.shortName, tt.defValue, tt.usage)
			} else {
				f := &flags.Int64Flag{}
				cmd.Int64Var(f, tt.longName, tt.shortName, tt.defValue, tt.usage)

				// 验证标志属性
				if f.LongName() != tt.longName {
					t.Errorf("Int64Var() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
				}
				if f.ShortName() != tt.shortName {
					t.Errorf("Int64Var() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
				}
				if f.Get() != tt.defValue {
					t.Errorf("Int64Var() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
				}
				if f.Usage() != tt.usage {
					t.Errorf("Int64Var() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
				}
			}
		})
	}
}

// TestInt64 测试64位整数类型标志创建
func TestInt64(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  int64
		usage     string
	}{
		{
			name:      "创建64位整数标志",
			longName:  "timestamp",
			shortName: "t",
			defValue:  1640995200, // 2022-01-01 00:00:00 UTC
			usage:     "时间戳",
		},
		{
			name:      "字节大小",
			longName:  "bytes",
			shortName: "b",
			defValue:  1073741824, // 1GB
			usage:     "字节大小",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createTestCmd()
			f := cmd.Int64(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if f == nil {
				t.Errorf("Int64() 返回nil指针")
				return
			}

			// 验证标志属性
			if f.LongName() != tt.longName {
				t.Errorf("Int64() 长标志名 = %v, 期望 %v", f.LongName(), tt.longName)
			}
			if f.ShortName() != tt.shortName {
				t.Errorf("Int64() 短标志名 = %v, 期望 %v", f.ShortName(), tt.shortName)
			}
			if f.Get() != tt.defValue {
				t.Errorf("Int64() 默认值 = %v, 期望 %v", f.Get(), tt.defValue)
			}
			if f.Usage() != tt.usage {
				t.Errorf("Int64() 使用说明 = %v, 期望 %v", f.Usage(), tt.usage)
			}
		})
	}
}

// =============================================================================
// 边界条件和错误处理测试
// =============================================================================

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("同时注册长短标志", func(t *testing.T) {
		cmd := createTestCmd()
		f := cmd.Bool("verbose", "v", false, "详细输出")

		// 验证两个标志都被注册
		if !cmd.FlagExists("verbose") {
			t.Error("长标志名未被注册")
		}
		if !cmd.FlagExists("v") {
			t.Error("短标志名未被注册")
		}

		// 验证标志对象相同
		if f.LongName() != "verbose" || f.ShortName() != "v" {
			t.Error("标志对象属性不正确")
		}
	})

	t.Run("只注册长标志", func(t *testing.T) {
		cmd := createTestCmd()
		f := cmd.String("output", "", "stdout", "输出目标")

		if !cmd.FlagExists("output") {
			t.Error("长标志名未被注册")
		}
		if f.ShortName() != "" {
			t.Error("短标志名应该为空")
		}
	})

	t.Run("只注册短标志", func(t *testing.T) {
		cmd := createTestCmd()
		f := cmd.Int("", "p", 8080, "端口号")

		if !cmd.FlagExists("p") {
			t.Error("短标志名未被注册")
		}
		if f.LongName() != "" {
			t.Error("长标志名应该为空")
		}
	})

	t.Run("极值测试", func(t *testing.T) {
		cmd := createTestCmd()

		// 测试极大浮点数
		f1 := cmd.Float64("max-float", "", 1.7976931348623157e+308, "最大浮点数")
		if f1.Get() != 1.7976931348623157e+308 {
			t.Error("极大浮点数设置失败")
		}

		// 测试极小浮点数
		f2 := cmd.Float64("min-float", "", 4.9406564584124654e-324, "最小浮点数")
		if f2.Get() != 4.9406564584124654e-324 {
			t.Error("极小浮点数设置失败")
		}

		// 测试最大int64
		f3 := cmd.Int64("max-int64", "", 9223372036854775807, "最大64位整数")
		if f3.Get() != 9223372036854775807 {
			t.Error("最大64位整数设置失败")
		}

		// 测试最小int64
		f4 := cmd.Int64("min-int64", "", -9223372036854775808, "最小64位整数")
		if f4.Get() != -9223372036854775808 {
			t.Error("最小64位整数设置失败")
		}
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		cmd := createTestCmd()

		// 测试包含特殊字符的使用说明
		f1 := cmd.String("test1", "", "default", "包含特殊字符: !@#$%^&*()")
		if f1.Usage() != "包含特殊字符: !@#$%^&*()" {
			t.Error("特殊字符使用说明处理失败")
		}

		// 测试中文使用说明
		f2 := cmd.Bool("test2", "", false, "这是中文使用说明")
		if f2.Usage() != "这是中文使用说明" {
			t.Error("中文使用说明处理失败")
		}

		// 测试emoji使用说明
		f3 := cmd.Int("test3", "", 0, "包含emoji: 🚀🎉✨")
		if f3.Usage() != "包含emoji: 🚀🎉✨" {
			t.Error("emoji使用说明处理失败")
		}
	})
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	t.Run("并发创建标志", func(t *testing.T) {
		cmd := createTestCmd()
		var wg sync.WaitGroup
		numGoroutines := 100

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				flagName := fmt.Sprintf("flag%d", id)
				f := cmd.Bool(flagName, "", false, fmt.Sprintf("标志%d", id))
				if f == nil {
					t.Errorf("并发创建标志%d失败", id)
				}
			}(i)
		}
		wg.Wait()

		// 验证所有标志都被创建
		for i := 0; i < numGoroutines; i++ {
			flagName := fmt.Sprintf("flag%d", i)
			if !cmd.FlagExists(flagName) {
				t.Errorf("标志%s未被正确创建", flagName)
			}
		}
	})

	t.Run("并发创建不同类型标志", func(t *testing.T) {
		cmd := createTestCmd()
		var wg sync.WaitGroup
		numGoroutines := 50

		wg.Add(numGoroutines * 5) // 5种类型的标志

		// 并发创建布尔标志
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				cmd.Bool(fmt.Sprintf("bool%d", id), "", false, "布尔标志")
			}(i)
		}

		// 并发创建字符串标志
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				cmd.String(fmt.Sprintf("string%d", id), "", "default", "字符串标志")
			}(i)
		}

		// 并发创建整数标志
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				cmd.Int(fmt.Sprintf("int%d", id), "", 0, "整数标志")
			}(i)
		}

		// 并发创建64位整数标志
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				cmd.Int64(fmt.Sprintf("int64_%d", id), "", 0, "64位整数标志")
			}(i)
		}

		// 并发创建浮点数标志
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				cmd.Float64(fmt.Sprintf("float64_%d", id), "", 0.0, "浮点数标志")
			}(i)
		}

		wg.Wait()

		// 验证所有标志都被创建 - 注意cmd已经有一个内置的help标志
		actualFlags := cmd.NFlag()
		if actualFlags < numGoroutines*5 {
			t.Errorf("期望创建至少%d个标志，实际创建%d个", numGoroutines*5, actualFlags)
		}
	})
}

// TestComplexScenarios 测试复杂场景
func TestComplexScenarios(t *testing.T) {
	t.Run("大量标志注册", func(t *testing.T) {
		cmd := createTestCmd()
		numFlags := 1000

		// 注册大量标志
		for i := 0; i < numFlags; i++ {
			switch i % 5 {
			case 0:
				cmd.Bool(fmt.Sprintf("bool%d", i), "", false, fmt.Sprintf("布尔标志%d", i))
			case 1:
				cmd.String(fmt.Sprintf("string%d", i), "", "default", fmt.Sprintf("字符串标志%d", i))
			case 2:
				cmd.Int(fmt.Sprintf("int%d", i), "", i, fmt.Sprintf("整数标志%d", i))
			case 3:
				cmd.Int64(fmt.Sprintf("int64_%d", i), "", int64(i), fmt.Sprintf("64位整数标志%d", i))
			case 4:
				cmd.Float64(fmt.Sprintf("float64_%d", i), "", float64(i), fmt.Sprintf("浮点数标志%d", i))
			}
		}

		// 验证标志数量
		if cmd.NFlag() < numFlags {
			t.Errorf("期望注册%d个标志，实际注册%d个", numFlags, cmd.NFlag())
		}

		// 随机验证一些标志
		testIndices := []int{0, 100, 500, 999}
		for _, i := range testIndices {
			flagName := fmt.Sprintf("bool%d", i)
			if i%5 == 0 && !cmd.FlagExists(flagName) {
				t.Errorf("标志%s未被正确注册", flagName)
			}
		}
	})

	t.Run("混合长短标志名", func(t *testing.T) {
		cmd := createTestCmd()

		// 创建各种组合的标志
		f1 := cmd.Bool("verbose", "v", false, "详细输出")
		f2 := cmd.String("output", "", "stdout", "输出目标")
		f3 := cmd.Int("", "p", 8080, "端口号")
		f4 := cmd.Float64("threshold", "t", 0.5, "阈值")
		f5 := cmd.Int64("size", "", 1024, "大小")

		// 验证标志属性
		if f1.LongName() != "verbose" || f1.ShortName() != "v" {
			t.Error("布尔标志长短名设置错误")
		}
		if f2.LongName() != "output" || f2.ShortName() != "" {
			t.Error("字符串标志长短名设置错误")
		}
		if f3.LongName() != "" || f3.ShortName() != "p" {
			t.Error("整数标志长短名设置错误")
		}
		if f4.LongName() != "threshold" || f4.ShortName() != "t" {
			t.Error("浮点数标志长短名设置错误")
		}
		if f5.LongName() != "size" || f5.ShortName() != "" {
			t.Error("64位整数标志长短名设置错误")
		}
	})
}

// TestPerformance 性能测试
func TestPerformance(t *testing.T) {
	t.Run("标志创建性能", func(t *testing.T) {
		cmd := createTestCmd()
		numFlags := 10000

		// 测试大量标志创建的性能
		for i := 0; i < numFlags; i++ {
			cmd.Bool(fmt.Sprintf("perf%d", i), "", false, "性能测试标志")
		}

		if cmd.NFlag() < numFlags {
			t.Errorf("性能测试失败，期望%d个标志，实际%d个", numFlags, cmd.NFlag())
		}
	})
}

// =============================================================================
// 基准测试
// =============================================================================

// BenchmarkBoolVar 基准测试布尔标志变量绑定
func BenchmarkBoolVar(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f := &flags.BoolFlag{}
		cmd.BoolVar(f, fmt.Sprintf("bool%d", i), "", false, "基准测试")
	}
}

// BenchmarkBool 基准测试布尔标志创建
func BenchmarkBool(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd.Bool(fmt.Sprintf("bool%d", i), "", false, "基准测试")
	}
}

// BenchmarkStringVar 基准测试字符串标志变量绑定
func BenchmarkStringVar(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f := &flags.StringFlag{}
		cmd.StringVar(f, fmt.Sprintf("string%d", i), "", "default", "基准测试")
	}
}

// BenchmarkString 基准测试字符串标志创建
func BenchmarkString(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd.String(fmt.Sprintf("string%d", i), "", "default", "基准测试")
	}
}

// BenchmarkIntVar 基准测试整数标志变量绑定
func BenchmarkIntVar(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f := &flags.IntFlag{}
		cmd.IntVar(f, fmt.Sprintf("int%d", i), "", 0, "基准测试")
	}
}

// BenchmarkInt 基准测试整数标志创建
func BenchmarkInt(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd.Int(fmt.Sprintf("int%d", i), "", 0, "基准测试")
	}
}

// BenchmarkFloat64Var 基准测试浮点数标志变量绑定
func BenchmarkFloat64Var(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f := &flags.Float64Flag{}
		cmd.Float64Var(f, fmt.Sprintf("float%d", i), "", 0.0, "基准测试")
	}
}

// BenchmarkFloat64 基准测试浮点数标志创建
func BenchmarkFloat64(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd.Float64(fmt.Sprintf("float%d", i), "", 0.0, "基准测试")
	}
}

// BenchmarkInt64Var 基准测试64位整数标志变量绑定
func BenchmarkInt64Var(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f := &flags.Int64Flag{}
		cmd.Int64Var(f, fmt.Sprintf("int64_%d", i), "", 0, "基准测试")
	}
}

// BenchmarkInt64 基准测试64位整数标志创建
func BenchmarkInt64(b *testing.B) {
	cmd := createTestCmd()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd.Int64(fmt.Sprintf("int64_%d", i), "", 0, "基准测试")
	}
}

// BenchmarkConcurrentFlagCreation 基准测试并发标志创建
func BenchmarkConcurrentFlagCreation(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cmd := createTestCmd()
			cmd.Bool(fmt.Sprintf("concurrent%d", i), "", false, "并发测试")
			i++
		}
	})
}
