package registry

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockFlag 模拟标志实现
type MockFlag struct {
	longName  string
	shortName string
	value     interface{}
	isSet     bool
}

func NewMockFlag(longName, shortName string, value interface{}) *MockFlag {
	return &MockFlag{
		longName:  longName,
		shortName: shortName,
		value:     value,
		isSet:     false,
	}
}

func (m *MockFlag) LongName() string {
	return m.longName
}

func (m *MockFlag) ShortName() string {
	return m.shortName
}

func (m *MockFlag) Usage() string {
	return "mock flag usage"
}

func (m *MockFlag) Type() flags.FlagType {
	return flags.FlagTypeString
}

func (m *MockFlag) GetDefaultAny() interface{} {
	return m.value
}

func (m *MockFlag) String() string {
	if m.value == nil {
		return ""
	}
	return fmt.Sprintf("%v", m.value)
}

func (m *MockFlag) IsSet() bool {
	return m.isSet
}

func (m *MockFlag) Reset() {
	m.isSet = false
}

func (m *MockFlag) GetEnvVar() string {
	return ""
}

// createTestContext 创建测试上下文，预先标记内置标志
func createTestContext() *types.CmdContext {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	// 标记内置标志名称
	ctx.BuiltinFlags.MarkAsBuiltin("help", "h", "version", "v", "generate-shell-completion", "gsc")

	return ctx
}

// TestRegisterFlag 测试标志注册功能
func TestRegisterFlag(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		flag      flags.Flag
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "注册有效的长标志",
			longName:  "test-flag",
			shortName: "",
			flag:      NewMockFlag("test-flag", "", "default"),
			wantErr:   false,
		},
		{
			name:      "注册有效的短标志",
			longName:  "",
			shortName: "t",
			flag:      NewMockFlag("", "t", "default"),
			wantErr:   false,
		},
		{
			name:      "注册长短标志都有效",
			longName:  "test-flag",
			shortName: "t",
			flag:      NewMockFlag("test-flag", "t", "default"),
			wantErr:   false,
		},
		{
			name:      "长短标志名都为空",
			longName:  "",
			shortName: "",
			flag:      NewMockFlag("", "", "default"),
			wantErr:   true,
			errMsg:    "flag long name and short name cannot both be empty",
		},
		{
			name:      "长标志名包含感叹号",
			longName:  "test!flag",
			shortName: "",
			flag:      NewMockFlag("test!flag", "", "default"),
			wantErr:   true,
		},
		{
			name:      "短标志名包含感叹号",
			longName:  "",
			shortName: "t!",
			flag:      NewMockFlag("", "t!", "default"),
			wantErr:   true,
		},
		{
			name:      "使用内置标志名help",
			longName:  "help",
			shortName: "",
			flag:      NewMockFlag("help", "", "default"),
			wantErr:   true,
			errMsg:    "flag long name 'help' is reserved",
		},
		{
			name:      "使用内置标志名h",
			longName:  "",
			shortName: "h",
			flag:      NewMockFlag("", "h", "default"),
			wantErr:   true,
			errMsg:    "flag short name 'h' is reserved",
		},
		{
			name:      "使用内置标志名version",
			longName:  "version",
			shortName: "",
			flag:      NewMockFlag("version", "", "default"),
			wantErr:   true,
			errMsg:    "flag long name 'version' is reserved",
		},
		{
			name:      "使用内置标志名v",
			longName:  "",
			shortName: "v",
			flag:      NewMockFlag("", "v", "default"),
			wantErr:   true,
			errMsg:    "flag short name 'v' is reserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext()

			err := RegisterFlag(ctx, tt.flag, tt.longName, tt.shortName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("RegisterFlag() 期望错误但未返回错误")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("RegisterFlag() 错误信息 = %v, 期望 %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("RegisterFlag() 意外错误 = %v", err)
				}
			}
		})
	}
}

// TestValidateFlagNames 测试标志名称验证
func TestValidateFlagNames(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "有效的长标志名",
			longName:  "valid-flag",
			shortName: "",
			wantErr:   false,
		},
		{
			name:      "有效的短标志名",
			longName:  "",
			shortName: "f",
			wantErr:   false,
		},
		{
			name:      "长短标志名都有效",
			longName:  "valid-flag",
			shortName: "f",
			wantErr:   false,
		},
		{
			name:      "长短标志名都为空",
			longName:  "",
			shortName: "",
			wantErr:   true,
			errMsg:    "flag long name and short name cannot both be empty",
		},
		{
			name:      "长标志名包含空格",
			longName:  "flag name",
			shortName: "",
			wantErr:   true,
		},
		{
			name:      "长标志名包含空格",
			longName:  "flag name",
			shortName: "",
			wantErr:   true,
		},
		{
			name:      "短标志名包含感叹号",
			longName:  "",
			shortName: "f!",
			wantErr:   true,
		},
		{
			name:      "使用内置长标志名",
			longName:  "help",
			shortName: "",
			wantErr:   true,
			errMsg:    "flag long name 'help' is reserved",
		},
		{
			name:      "使用内置短标志名",
			longName:  "",
			shortName: "h",
			wantErr:   true,
			errMsg:    "flag short name 'h' is reserved",
		},
		{
			name:      "使用内置version标志",
			longName:  "version",
			shortName: "",
			wantErr:   true,
			errMsg:    "flag long name 'version' is reserved",
		},
		{
			name:      "使用内置v标志",
			longName:  "",
			shortName: "v",
			wantErr:   true,
			errMsg:    "flag short name 'v' is reserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext()

			err := ValidateFlagNames(ctx, tt.longName, tt.shortName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateFlagNames() 期望错误但未返回错误")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidateFlagNames() 错误信息 = %v, 期望 %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFlagNames() 意外错误 = %v", err)
				}
			}
		})
	}
}

// TestValidateSingleFlagName 测试单个标志名称验证
func TestValidateSingleFlagName(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		nameType string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "有效标志名",
			flagName: "valid-flag",
			nameType: "long name",
			wantErr:  false,
		},
		{
			name:     "包含非法字符感叹号",
			flagName: "flag!value",
			nameType: "long name",
			wantErr:  true,
		},
		{
			name:     "包含非法字符空格",
			flagName: "flag name",
			nameType: "long name",
			wantErr:  true,
		},
		{
			name:     "内置标志help",
			flagName: "help",
			nameType: "long name",
			wantErr:  true,
			errMsg:   "flag long name 'help' is reserved",
		},
		{
			name:     "内置标志h",
			flagName: "h",
			nameType: "short name",
			wantErr:  true,
			errMsg:   "flag short name 'h' is reserved",
		},
		{
			name:     "内置标志version",
			flagName: "version",
			nameType: "long name",
			wantErr:  true,
			errMsg:   "flag long name 'version' is reserved",
		},
		{
			name:     "内置标志v",
			flagName: "v",
			nameType: "short name",
			wantErr:  true,
			errMsg:   "flag short name 'v' is reserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext()

			err := validateSingleFlagName(ctx, tt.flagName, tt.nameType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateSingleFlagName() 期望错误但未返回错误")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("validateSingleFlagName() 错误信息 = %v, 期望 %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateSingleFlagName() 意外错误 = %v", err)
				}
			}
		})
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	t.Run("极长标志名", func(t *testing.T) {
		longName := strings.Repeat("a", 1000)

		err := ValidateFlagNames(ctx, longName, "")
		if err != nil {
			t.Errorf("极长标志名验证失败: %v", err)
		}
	})

	t.Run("特殊字符测试", func(t *testing.T) {
		specialChars := []string{
			"flag-with-dash",
			"flag_with_underscore",
			"flag123",
			"123flag",
		}

		for _, name := range specialChars {
			err := ValidateFlagNames(ctx, name, "")
			if err != nil && !containsInvalidChars(name) {
				t.Errorf("特殊字符标志名 %q 验证失败: %v", name, err)
			}
		}
	})

	t.Run("非法字符测试", func(t *testing.T) {
		invalidNames := []string{
			"flag!value",
			"flag value",
			"flag@value",
			"flag#value",
		}

		for _, name := range invalidNames {
			err := ValidateFlagNames(ctx, name, "")
			if err == nil {
				t.Errorf("包含非法字符的标志名 %q 应该验证失败", name)
			}
		}
	})

	t.Run("空值测试", func(t *testing.T) {
		err := ValidateFlagNames(ctx, "", "")
		if err == nil {
			t.Error("长短标志名都为空应该返回错误")
		}
	})

	t.Run("nil上下文测试", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("nil上下文应该导致panic或错误")
			}
		}()
		err := ValidateFlagNames(nil, "test", "")
		if err == nil {
			t.Error("期望验证标志名称时返回错误")
		}
	})
}

// containsInvalidChars 检查字符串是否包含非法字符
func containsInvalidChars(s string) bool {
	// 使用flags包中定义的InvalidFlagChars
	return strings.ContainsAny(s, flags.InvalidFlagChars)
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	t.Run("并发验证标志名", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100
		errors := make(chan error, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				flagName := fmt.Sprintf("flag-%d", id)
				err := ValidateFlagNames(ctx, flagName, "")
				if err != nil {
					errors <- err
				}
			}(i)
		}
		wg.Wait()
		close(errors)

		for err := range errors {
			t.Errorf("并发验证标志名失败: %v", err)
		}
	})

	t.Run("并发注册标志", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 50
		errors := make(chan error, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				flagName := fmt.Sprintf("concurrent-flag-%d", id)
				flag := NewMockFlag(flagName, "", "default")
				err := RegisterFlag(ctx, flag, flagName, "")
				if err != nil {
					errors <- err
				}
			}(i)
		}
		wg.Wait()
		close(errors)

		for err := range errors {
			t.Errorf("并发注册标志失败: %v", err)
		}
	})
}

// TestPerformance 性能测试
func TestPerformance(t *testing.T) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	t.Run("大量标志验证性能", func(t *testing.T) {
		numFlags := 10000

		for i := 0; i < numFlags; i++ {
			flagName := fmt.Sprintf("perf-flag-%d", i)
			err := ValidateFlagNames(ctx, flagName, "")
			if err != nil {
				t.Errorf("性能测试中标志验证失败: %v", err)
				break
			}
		}
	})
}

// TestNilInputs 测试nil输入处理
func TestNilInputs(t *testing.T) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	t.Run("nil标志注册", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("注册nil标志应该导致panic或返回错误")
			}
		}()
		err := RegisterFlag(ctx, nil, "test", "")
		if err == nil {
			t.Error("期望注册标志时返回错误")
		}
	})

	t.Run("nil上下文标志注册", func(t *testing.T) {
		flag := NewMockFlag("test", "", "default")
		defer func() {
			if r := recover(); r == nil {
				t.Error("nil上下文应该导致panic")
			}
		}()
		err := RegisterFlag(nil, flag, "test", "")
		if err == nil {
			t.Error("期望注册标志时返回错误")
		}
	})
}

// TestTypeConversion 测试类型转换边界
func TestTypeConversion(t *testing.T) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	t.Run("不同类型标志", func(t *testing.T) {
		// 测试不同类型的标志
		testCases := []struct {
			name  string
			value interface{}
		}{
			{"string-flag", "test"},
			{"int-flag", 42},
			{"bool-flag", true},
		}

		for _, tc := range testCases {
			flag := NewMockFlag(tc.name, "", tc.value)
			err := RegisterFlag(ctx, flag, tc.name, "")
			if err != nil {
				t.Errorf("注册%s类型标志失败: %v", tc.name, err)
			}
		}
	})
}

// TestInvalidFlagChars 测试非法字符常量
func TestInvalidFlagChars(t *testing.T) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	t.Run("测试所有非法字符", func(t *testing.T) {
		// 测试flags.InvalidFlagChars中的每个字符
		for _, char := range flags.InvalidFlagChars {
			flagName := "flag" + string(char) + "name"
			err := ValidateFlagNames(ctx, flagName, "")
			if err == nil {
				t.Errorf("包含非法字符 %q 的标志名应该验证失败", char)
			}
		}
	})
}

// BenchmarkValidateFlagNames 基准测试标志名验证
func BenchmarkValidateFlagNames(b *testing.B) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flagName := fmt.Sprintf("benchmark-flag-%d", i%1000)
		_ = ValidateFlagNames(ctx, flagName, "")
	}
}

// BenchmarkRegisterFlag 基准测试标志注册
func BenchmarkRegisterFlag(b *testing.B) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flagName := fmt.Sprintf("benchmark-flag-%d", i%1000)
		flag := NewMockFlag(flagName, "", "default")
		_ = RegisterFlag(ctx, flag, flagName, "")
	}
}

// BenchmarkConcurrentValidation 基准测试并发验证
func BenchmarkConcurrentValidation(b *testing.B) {
	ctx := &types.CmdContext{
		FlagRegistry: flags.NewFlagRegistry(),
		BuiltinFlags: types.NewBuiltinFlags(),
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			flagName := fmt.Sprintf("concurrent-flag-%d", i%1000)
			_ = ValidateFlagNames(ctx, flagName, "")
			i++
		}
	})
}
