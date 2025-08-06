package cmd

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// =============================================================================
// 测试辅助函数
// =============================================================================

// createExtendedTestCmd 创建扩展测试用的命令实例
func createExtendedTestCmd() *Cmd {
	return NewCmd("extended-test", "et", flag.ContinueOnError)
}

// createExtendedTestCmdWithBuiltins 创建包含内置标志的扩展测试命令
func createExtendedTestCmdWithBuiltins() *Cmd {
	cmd := NewCmd("extended-test", "et", flag.ContinueOnError)
	// 添加内置标志到映射中
	cmd.ctx.BuiltinFlags.NameMap.Store("help", true)
	cmd.ctx.BuiltinFlags.NameMap.Store("h", true)
	cmd.ctx.BuiltinFlags.NameMap.Store("version", true)
	cmd.ctx.BuiltinFlags.NameMap.Store("v", true)
	return cmd
}

// =============================================================================
// 枚举类型标志测试
// =============================================================================

func TestCmd_Enum(t *testing.T) {
	tests := []struct {
		name        string
		longName    string
		shortName   string
		defValue    string
		usage       string
		options     []string
		expectPanic bool
		panicMsg    string
		setupCmd    func() *Cmd
	}{
		{
			name:      "正常创建枚举标志",
			longName:  "mode",
			shortName: "m",
			defValue:  "debug",
			usage:     "运行模式",
			options:   []string{"debug", "release", "test"},
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "仅长名称枚举标志",
			longName:  "level",
			shortName: "",
			defValue:  "info",
			usage:     "日志级别",
			options:   []string{"debug", "info", "warn", "error"},
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "仅短名称枚举标志",
			longName:  "",
			shortName: "l",
			defValue:  "low",
			usage:     "优先级",
			options:   []string{"low", "medium", "high"},
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "空选项列表",
			longName:  "empty",
			shortName: "e",
			defValue:  "",
			usage:     "空选项测试",
			options:   []string{},
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "nil选项列表",
			longName:  "nil",
			shortName: "n",
			defValue:  "",
			usage:     "nil选项测试",
			options:   nil,
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:        "使用保留的长名称",
			longName:    "help",
			shortName:   "x",
			defValue:    "test",
			usage:       "测试保留名称",
			options:     []string{"test"},
			expectPanic: true,
			panicMsg:    "flag long name help is reserved",
			setupCmd:    createExtendedTestCmdWithBuiltins,
		},
		{
			name:        "使用保留的短名称",
			longName:    "test",
			shortName:   "h",
			defValue:    "test",
			usage:       "测试保留名称",
			options:     []string{"test"},
			expectPanic: true,
			panicMsg:    "flag short name h is reserved",
			setupCmd:    createExtendedTestCmdWithBuiltins,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Enum() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Enum() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("Enum() 期望panic但未发生")
				}
			}()

			flag := cmd.Enum(tt.longName, tt.shortName, tt.defValue, tt.usage, tt.options)

			if !tt.expectPanic {
				if flag == nil {
					t.Error("Enum() 返回nil")
					return
				}

				// 验证标志属性
				if flag.LongName() != tt.longName {
					t.Errorf("LongName() = %v, 期望 %v", flag.LongName(), tt.longName)
				}
				if flag.ShortName() != tt.shortName {
					t.Errorf("ShortName() = %v, 期望 %v", flag.ShortName(), tt.shortName)
				}
				if flag.Usage() != tt.usage {
					t.Errorf("Usage() = %v, 期望 %v", flag.Usage(), tt.usage)
				}
				if flag.Get() != tt.defValue {
					t.Errorf("Get() = %v, 期望 %v", flag.Get(), tt.defValue)
				}
				if flag.Type() != flags.FlagTypeEnum {
					t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeEnum)
				}
			}
		})
	}
}

func TestCmd_EnumVar(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.EnumFlag
		longName    string
		shortName   string
		defValue    string
		usage       string
		options     []string
		expectPanic bool
		panicMsg    string
		setupCmd    func() *Cmd
	}{
		{
			name:      "正常绑定枚举标志",
			flag:      &flags.EnumFlag{},
			longName:  "format",
			shortName: "f",
			defValue:  "json",
			usage:     "输出格式",
			options:   []string{"json", "xml", "yaml"},
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			longName:    "test",
			shortName:   "t",
			defValue:    "test",
			usage:       "测试",
			options:     []string{"test"},
			expectPanic: true,
			panicMsg:    "EnumFlag pointer cannot be nil",
			setupCmd:    createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("EnumVar() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("EnumVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("EnumVar() 期望panic但未发生")
				}
			}()

			cmd.EnumVar(tt.flag, tt.longName, tt.shortName, tt.defValue, tt.usage, tt.options)

			if !tt.expectPanic && tt.flag != nil {
				// 验证标志属性
				if tt.flag.LongName() != tt.longName {
					t.Errorf("LongName() = %v, 期望 %v", tt.flag.LongName(), tt.longName)
				}
				if tt.flag.ShortName() != tt.shortName {
					t.Errorf("ShortName() = %v, 期望 %v", tt.flag.ShortName(), tt.shortName)
				}
			}
		})
	}
}

// =============================================================================
// 无符号整数类型标志测试
// =============================================================================

func TestCmd_Uint16(t *testing.T) {
	tests := []struct {
		name        string
		longName    string
		shortName   string
		defValue    uint16
		usage       string
		expectPanic bool
		panicMsg    string
		setupCmd    func() *Cmd
	}{
		{
			name:      "正常创建uint16标志",
			longName:  "port",
			shortName: "p",
			defValue:  8080,
			usage:     "端口号",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "最小值uint16标志",
			longName:  "min",
			shortName: "m",
			defValue:  0,
			usage:     "最小值",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "最大值uint16标志",
			longName:  "max",
			shortName: "x",
			defValue:  65535,
			usage:     "最大值",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:        "使用保留的长名称",
			longName:    "version",
			shortName:   "p",
			defValue:    1,
			usage:       "版本号",
			expectPanic: true,
			panicMsg:    "flag long name version is reserved",
			setupCmd:    createExtendedTestCmdWithBuiltins,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Uint16() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Uint16() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("Uint16() 期望panic但未发生")
				}
			}()

			flag := cmd.Uint16(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if !tt.expectPanic {
				if flag == nil {
					t.Error("Uint16() 返回nil")
					return
				}

				// 验证标志属性
				if flag.Get() != tt.defValue {
					t.Errorf("Get() = %v, 期望 %v", flag.Get(), tt.defValue)
				}
				if flag.Type() != flags.FlagTypeUint16 {
					t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeUint16)
				}
			}
		})
	}
}

func TestCmd_Uint16Var(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.Uint16Flag
		longName    string
		shortName   string
		defValue    uint16
		usage       string
		expectPanic bool
		panicMsg    string
		setupCmd    func() *Cmd
	}{
		{
			name:      "正常绑定uint16标志",
			flag:      &flags.Uint16Flag{},
			longName:  "timeout",
			shortName: "t",
			defValue:  300,
			usage:     "超时时间",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			longName:    "test",
			shortName:   "t",
			defValue:    1,
			usage:       "测试",
			expectPanic: true,
			panicMsg:    "Uint16Flag pointer cannot be nil",
			setupCmd:    createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Uint16Var() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Uint16Var() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("Uint16Var() 期望panic但未发生")
				}
			}()

			cmd.Uint16Var(tt.flag, tt.longName, tt.shortName, tt.defValue, tt.usage)

			if !tt.expectPanic && tt.flag != nil {
				if tt.flag.Get() != tt.defValue {
					t.Errorf("Get() = %v, 期望 %v", tt.flag.Get(), tt.defValue)
				}
			}
		})
	}
}

func TestCmd_Uint32(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  uint32
		usage     string
		setupCmd  func() *Cmd
	}{
		{
			name:      "正常创建uint32标志",
			longName:  "size",
			shortName: "s",
			defValue:  1024,
			usage:     "大小",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "最大值uint32标志",
			longName:  "max",
			shortName: "m",
			defValue:  4294967295,
			usage:     "最大值",
			setupCmd:  createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			flag := cmd.Uint32(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if flag == nil {
				t.Error("Uint32() 返回nil")
				return
			}

			if flag.Get() != tt.defValue {
				t.Errorf("Get() = %v, 期望 %v", flag.Get(), tt.defValue)
			}
			if flag.Type() != flags.FlagTypeUint32 {
				t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeUint32)
			}
		})
	}
}

func TestCmd_Uint32Var(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.Uint32Flag
		expectPanic bool
		panicMsg    string
	}{
		{
			name: "正常绑定uint32标志",
			flag: &flags.Uint32Flag{},
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			expectPanic: true,
			panicMsg:    "Uint32Flag pointer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createExtendedTestCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Uint32Var() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Uint32Var() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("Uint32Var() 期望panic但未发生")
				}
			}()

			cmd.Uint32Var(tt.flag, "test", "t", 100, "测试")
		})
	}
}

func TestCmd_Uint64(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  uint64
		usage     string
		setupCmd  func() *Cmd
	}{
		{
			name:      "正常创建uint64标志",
			longName:  "memory",
			shortName: "m",
			defValue:  1073741824,
			usage:     "内存大小",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "最大值uint64标志",
			longName:  "max",
			shortName: "x",
			defValue:  18446744073709551615,
			usage:     "最大值",
			setupCmd:  createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			flag := cmd.Uint64(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if flag == nil {
				t.Error("Uint64() 返回nil")
				return
			}

			if flag.Get() != tt.defValue {
				t.Errorf("Get() = %v, 期望 %v", flag.Get(), tt.defValue)
			}
			if flag.Type() != flags.FlagTypeUint64 {
				t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeUint64)
			}
		})
	}
}

func TestCmd_Uint64Var(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.Uint64Flag
		expectPanic bool
		panicMsg    string
	}{
		{
			name: "正常绑定uint64标志",
			flag: &flags.Uint64Flag{},
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			expectPanic: true,
			panicMsg:    "Uint64Flag pointer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createExtendedTestCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Uint64Var() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("Uint64Var() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("Uint64Var() 期望panic但未发生")
				}
			}()

			cmd.Uint64Var(tt.flag, "test", "t", 100, "测试")
		})
	}
}

// =============================================================================
// 时间类型标志测试
// =============================================================================

func TestCmd_Time(t *testing.T) {
	defaultTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  time.Time
		usage     string
		setupCmd  func() *Cmd
	}{
		{
			name:      "正常创建时间标志",
			longName:  "start",
			shortName: "s",
			defValue:  defaultTime,
			usage:     "开始时间",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "零值时间标志",
			longName:  "zero",
			shortName: "z",
			defValue:  time.Time{},
			usage:     "零值时间",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "当前时间标志",
			longName:  "now",
			shortName: "n",
			defValue:  time.Now(),
			usage:     "当前时间",
			setupCmd:  createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			flag := cmd.Time(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if flag == nil {
				t.Error("Time() 返回nil")
				return
			}

			// 验证标志属性
			if !flag.Get().Equal(tt.defValue) {
				t.Errorf("Get() = %v, 期望 %v", flag.Get(), tt.defValue)
			}
			if flag.Type() != flags.FlagTypeTime {
				t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeTime)
			}
		})
	}
}

func TestCmd_TimeVar(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.TimeFlag
		expectPanic bool
		panicMsg    string
	}{
		{
			name: "正常绑定时间标志",
			flag: &flags.TimeFlag{},
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			expectPanic: true,
			panicMsg:    "TimeFlag pointer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createExtendedTestCmd()
			defaultTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("TimeVar() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("TimeVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("TimeVar() 期望panic但未发生")
				}
			}()

			cmd.TimeVar(tt.flag, "test", "t", defaultTime, "测试时间")

			if !tt.expectPanic && tt.flag != nil {
				if !tt.flag.Get().Equal(defaultTime) {
					t.Errorf("Get() = %v, 期望 %v", tt.flag.Get(), defaultTime)
				}
			}
		})
	}
}

// =============================================================================
// 时间间隔类型标志测试
// =============================================================================

func TestCmd_Duration(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  time.Duration
		usage     string
		setupCmd  func() *Cmd
	}{
		{
			name:      "正常创建时间间隔标志",
			longName:  "timeout",
			shortName: "t",
			defValue:  30 * time.Second,
			usage:     "超时时间",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "零值时间间隔标志",
			longName:  "zero",
			shortName: "z",
			defValue:  0,
			usage:     "零值时间间隔",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "长时间间隔标志",
			longName:  "long",
			shortName: "l",
			defValue:  24 * time.Hour,
			usage:     "长时间间隔",
			setupCmd:  createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			flag := cmd.Duration(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if flag == nil {
				t.Error("Duration() 返回nil")
				return
			}

			if flag.Get() != tt.defValue {
				t.Errorf("Get() = %v, 期望 %v", flag.Get(), tt.defValue)
			}
			if flag.Type() != flags.FlagTypeDuration {
				t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeDuration)
			}
		})
	}
}

func TestCmd_DurationVar(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.DurationFlag
		expectPanic bool
		panicMsg    string
	}{
		{
			name: "正常绑定时间间隔标志",
			flag: &flags.DurationFlag{},
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			expectPanic: true,
			panicMsg:    "DurationFlag pointer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createExtendedTestCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("DurationVar() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("DurationVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("DurationVar() 期望panic但未发生")
				}
			}()

			cmd.DurationVar(tt.flag, "test", "t", 5*time.Second, "测试时间间隔")
		})
	}
}

// =============================================================================
// 键值对类型标志测试
// =============================================================================

func TestCmd_Map(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  map[string]string
		usage     string
		setupCmd  func() *Cmd
	}{
		{
			name:      "正常创建键值对标志",
			longName:  "config",
			shortName: "c",
			defValue:  map[string]string{"key1": "value1", "key2": "value2"},
			usage:     "配置项",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "空键值对标志",
			longName:  "empty",
			shortName: "e",
			defValue:  map[string]string{},
			usage:     "空配置",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "nil键值对标志",
			longName:  "nil",
			shortName: "n",
			defValue:  nil,
			usage:     "nil配置",
			setupCmd:  createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			flag := cmd.Map(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if flag == nil {
				t.Error("Map() 返回nil")
				return
			}

			// 验证标志属性
			got := flag.Get()
			if tt.defValue == nil {
				// nil应该被转换为空map
				if len(got) != 0 {
					t.Errorf("Get() = %v, 期望空map", got)
				}
			} else if !reflect.DeepEqual(got, tt.defValue) {
				t.Errorf("Get() = %v, 期望 %v", got, tt.defValue)
			}
			if flag.Type() != flags.FlagTypeMap {
				t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeMap)
			}
		})
	}
}

func TestCmd_MapVar(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.MapFlag
		expectPanic bool
		panicMsg    string
	}{
		{
			name: "正常绑定键值对标志",
			flag: &flags.MapFlag{},
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			expectPanic: true,
			panicMsg:    "MapFlag pointer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createExtendedTestCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("MapVar() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("MapVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("MapVar() 期望panic但未发生")
				}
			}()

			cmd.MapVar(tt.flag, "test", "t", map[string]string{"key": "value"}, "测试键值对")
		})
	}
}

// =============================================================================
// 切片类型标志测试
// =============================================================================

func TestCmd_Slice(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		defValue  []string
		usage     string
		setupCmd  func() *Cmd
	}{
		{
			name:      "正常创建切片标志",
			longName:  "files",
			shortName: "f",
			defValue:  []string{"file1.txt", "file2.txt"},
			usage:     "文件列表",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "空切片标志",
			longName:  "empty",
			shortName: "e",
			defValue:  []string{},
			usage:     "空文件列表",
			setupCmd:  createExtendedTestCmd,
		},
		{
			name:      "nil切片标志",
			longName:  "nil",
			shortName: "n",
			defValue:  nil,
			usage:     "nil文件列表",
			setupCmd:  createExtendedTestCmd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			flag := cmd.Slice(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if flag == nil {
				t.Error("Slice() 返回nil")
				return
			}

			// 验证标志属性
			got := flag.Get()
			if tt.defValue == nil {
				// nil应该被转换为空切片
				if len(got) != 0 {
					t.Errorf("Get() = %v, 期望空切片", got)
				}
			} else if !reflect.DeepEqual(got, tt.defValue) {
				t.Errorf("Get() = %v, 期望 %v", got, tt.defValue)
			}
			if flag.Type() != flags.FlagTypeSlice {
				t.Errorf("Type() = %v, 期望 %v", flag.Type(), flags.FlagTypeSlice)
			}
		})
	}
}

func TestCmd_SliceVar(t *testing.T) {
	tests := []struct {
		name        string
		flag        *flags.SliceFlag
		expectPanic bool
		panicMsg    string
	}{
		{
			name: "正常绑定切片标志",
			flag: &flags.SliceFlag{},
		},
		{
			name:        "nil标志指针",
			flag:        nil,
			expectPanic: true,
			panicMsg:    "SliceFlag pointer cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createExtendedTestCmd()

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("SliceVar() 意外的panic: %v", r)
					} else if tt.panicMsg != "" && !strings.Contains(fmt.Sprintf("%v", r), tt.panicMsg) {
						t.Errorf("SliceVar() panic信息 = %v, 期望包含 %v", r, tt.panicMsg)
					}
				} else if tt.expectPanic {
					t.Error("SliceVar() 期望panic但未发生")
				}
			}()

			cmd.SliceVar(tt.flag, "test", "t", []string{"item1", "item2"}, "测试切片")
		})
	}
}

// =============================================================================
// 并发安全测试
// =============================================================================

func TestCmd_ConcurrentFlagCreation(t *testing.T) {
	cmd := createExtendedTestCmd()

	var wg sync.WaitGroup
	numGoroutines := 10

	// 测试并发创建不同类型的标志
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// 创建不同类型的标志
			cmd.Enum(fmt.Sprintf("enum-%d", id), fmt.Sprintf("e%d", id), "default", "测试枚举", []string{"default", "option1", "option2"})
			cmd.Uint16(fmt.Sprintf("uint16-%d", id), fmt.Sprintf("u%d", id), uint16(id), "测试uint16")
			cmd.Uint32(fmt.Sprintf("uint32-%d", id), fmt.Sprintf("v%d", id), uint32(id), "测试uint32")
			cmd.Uint64(fmt.Sprintf("uint64-%d", id), fmt.Sprintf("w%d", id), uint64(id), "测试uint64")
			cmd.Time(fmt.Sprintf("time-%d", id), fmt.Sprintf("t%d", id), time.Now(), "测试时间")
			cmd.Duration(fmt.Sprintf("duration-%d", id), fmt.Sprintf("d%d", id), time.Duration(id)*time.Second, "测试时间间隔")
			cmd.Map(fmt.Sprintf("map-%d", id), fmt.Sprintf("m%d", id), map[string]string{"key": fmt.Sprintf("value-%d", id)}, "测试键值对")
			cmd.Slice(fmt.Sprintf("slice-%d", id), fmt.Sprintf("s%d", id), []string{fmt.Sprintf("item-%d", id)}, "测试切片")
		}(i)
	}

	wg.Wait()

	// 验证所有标志都已正确创建
	expectedFlags := numGoroutines * 8 // 每个goroutine创建8个标志
	actualFlags := cmd.NFlag()
	if actualFlags != expectedFlags+1 { // +1 for built-in help flag
		t.Errorf("期望创建 %d 个标志，实际创建 %d 个", expectedFlags+1, actualFlags)
	}
}

// =============================================================================
// 边界条件测试
// =============================================================================

func TestCmd_ExtendedBoundaryConditions(t *testing.T) {
	t.Run("极长标志名称", func(t *testing.T) {
		cmd := createExtendedTestCmd()
		longName := strings.Repeat("a", 1000)
		shortName := "x"

		flag := cmd.Enum(longName, shortName, "default", "极长标志名称测试", []string{"default", "option1"})
		if flag == nil {
			t.Errorf("标志 %s 不应为 nil", longName)
			return
		}
		if flag.LongName() != longName {
			t.Error("极长标志名称不匹配")
		}
	})

	t.Run("空字符串标志名称组合", func(t *testing.T) {
		cmd := createExtendedTestCmd()

		// 测试仅长名称
		flag1 := cmd.Uint16("only-long", "", 100, "仅长名称")
		if flag1 == nil {
			t.Error("仅长名称标志创建失败")
		}

		// 测试仅短名称
		flag2 := cmd.Uint32("", "s", 200, "仅短名称")
		if flag2 == nil {
			t.Error("仅短名称标志创建失败")
		}
	})

	t.Run("极值测试", func(t *testing.T) {
		cmd := createExtendedTestCmd()

		// uint16最大值
		flag1 := cmd.Uint16("max-uint16", "m1", 65535, "uint16最大值")
		if flag1.Get() != 65535 {
			t.Errorf("uint16最大值 = %v, 期望 65535", flag1.Get())
		}

		// uint32最大值
		flag2 := cmd.Uint32("max-uint32", "m2", 4294967295, "uint32最大值")
		if flag2.Get() != 4294967295 {
			t.Errorf("uint32最大值 = %v, 期望 4294967295", flag2.Get())
		}

		// uint64最大值
		flag3 := cmd.Uint64("max-uint64", "m3", 18446744073709551615, "uint64最大值")
		if flag3.Get() != 18446744073709551615 {
			t.Errorf("uint64最大值 = %v, 期望 18446744073709551615", flag3.Get())
		}
	})

	t.Run("特殊时间值测试", func(t *testing.T) {
		cmd := createExtendedTestCmd()

		// Unix纪元时间
		epochTime := time.Unix(0, 0)
		flag1 := cmd.Time("epoch", "e", epochTime, "Unix纪元时间")
		if !flag1.Get().Equal(epochTime) {
			t.Errorf("Unix纪元时间 = %v, 期望 %v", flag1.Get(), epochTime)
		}

		// 最大时间间隔
		maxDuration := time.Duration(1<<63 - 1)
		flag2 := cmd.Duration("max-duration", "md", maxDuration, "最大时间间隔")
		if flag2.Get() != maxDuration {
			t.Errorf("最大时间间隔 = %v, 期望 %v", flag2.Get(), maxDuration)
		}
	})
}

// =============================================================================
// 错误处理测试
// =============================================================================

func TestCmd_ExtendedErrorHandling(t *testing.T) {
	t.Run("重复标志名称", func(t *testing.T) {
		cmd := createExtendedTestCmd()

		// 创建第一个标志
		flag1 := cmd.Enum("duplicate", "d", "value1", "第一个标志", []string{"value1", "value2"})
		if flag1 == nil {
			t.Error("第一个标志创建失败")
		}

		// 尝试创建重复的标志应该panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望重复标志名称时panic，但没有panic")
			}
		}()

		cmd.Uint16("duplicate", "x", 100, "重复的标志")
	})

	t.Run("内置标志冲突", func(t *testing.T) {
		cmd := createExtendedTestCmdWithBuiltins()

		// 测试与help标志冲突
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望与内置标志冲突时panic，但没有panic")
			} else if !strings.Contains(fmt.Sprintf("%v", r), "reserved") {
				t.Errorf("panic信息应包含'reserved'，实际: %v", r)
			}
		}()

		cmd.Enum("help", "x", "default", "与help冲突", []string{"default"})
	})
}

// =============================================================================
// 性能测试
// =============================================================================

func BenchmarkCmd_EnumCreation(b *testing.B) {
	cmd := createExtendedTestCmd()
	options := []string{"option1", "option2", "option3", "option4", "option5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Enum(fmt.Sprintf("enum-%d", i), fmt.Sprintf("e%d", i), "option1", "基准测试枚举", options)
	}
}

func BenchmarkCmd_UintCreation(b *testing.B) {
	cmd := createExtendedTestCmd()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Uint64(fmt.Sprintf("uint-%d", i), fmt.Sprintf("u%d", i), uint64(i), "基准测试uint64")
	}
}

func BenchmarkCmd_SliceCreation(b *testing.B) {
	cmd := createExtendedTestCmd()
	defaultSlice := []string{"item1", "item2", "item3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Slice(fmt.Sprintf("slice-%d", i), fmt.Sprintf("s%d", i), defaultSlice, "基准测试切片")
	}
}

// =============================================================================
// 集成测试
// =============================================================================

func TestCmd_ExtendedIntegration(t *testing.T) {
	t.Run("混合标志类型创建", func(t *testing.T) {
		cmd := createExtendedTestCmd()

		// 创建各种类型的标志
		enumFlag := cmd.Enum("format", "f", "json", "输出格式", []string{"json", "xml", "yaml"})
		uint16Flag := cmd.Uint16("port", "p", 8080, "端口号")
		uint32Flag := cmd.Uint32("size", "s", 1024, "大小")
		uint64Flag := cmd.Uint64("memory", "m", 1073741824, "内存")
		timeFlag := cmd.Time("start", "st", time.Now(), "开始时间")
		durationFlag := cmd.Duration("timeout", "t", 30*time.Second, "超时")
		mapFlag := cmd.Map("config", "c", map[string]string{"key": "value"}, "配置")
		sliceFlag := cmd.Slice("files", "fl", []string{"file1", "file2"}, "文件列表")

		// 验证所有标志都已创建
		flags := [...]interface{}{enumFlag, uint16Flag, uint32Flag, uint64Flag, timeFlag, durationFlag, mapFlag, sliceFlag}
		for i, flag := range flags {
			if flag == nil {
				t.Errorf("标志 %d 创建失败", i)
			}
		}

		// 验证标志数量
		expectedCount := len(flags) + 1 // +1 for built-in help flag
		if cmd.NFlag() != expectedCount {
			t.Errorf("期望 %d 个标志，实际 %d 个", expectedCount, cmd.NFlag())
		}
	})

	t.Run("标志注册验证", func(t *testing.T) {
		cmd := createExtendedTestCmd()

		// 创建标志
		cmd.Enum("test-enum", "te", "default", "测试枚举", []string{"default", "other"})
		cmd.Uint16("test-uint16", "tu", 100, "测试uint16")

		// 验证标志是否已注册
		if !cmd.FlagExists("test-enum") {
			t.Error("枚举标志未正确注册")
		}
		if !cmd.FlagExists("te") {
			t.Error("枚举标志短名称未正确注册")
		}
		if !cmd.FlagExists("test-uint16") {
			t.Error("uint16标志未正确注册")
		}
		if !cmd.FlagExists("tu") {
			t.Error("uint16标志短名称未正确注册")
		}
	})
}
