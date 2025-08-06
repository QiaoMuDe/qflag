package parser

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockFlag 模拟标志实现，用于测试
type MockFlag struct {
	longName  string
	shortName string
	value     string
	isSet     bool
	envVar    string
}

func NewMockFlag(longName, shortName, defaultValue string) *MockFlag {
	return &MockFlag{
		longName:  longName,
		shortName: shortName,
		value:     defaultValue,
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
	return m.value
}

func (m *MockFlag) IsSet() bool {
	return m.isSet
}

func (m *MockFlag) Reset() {
	m.isSet = false
}

func (m *MockFlag) GetEnvVar() string {
	return m.envVar
}

func (m *MockFlag) Set(value string) error {
	m.value = value
	m.isSet = true
	return nil
}

func (m *MockFlag) BindEnv(envName string) {
	m.envVar = envName
}

// ErrorMockFlag 用于测试错误情况的Mock标志
type ErrorMockFlag struct {
	*MockFlag
	shouldError bool
	errorMsg    string
}

func NewErrorMockFlag(longName, shortName, defaultValue string, shouldError bool, errorMsg string) *ErrorMockFlag {
	return &ErrorMockFlag{
		MockFlag:    NewMockFlag(longName, shortName, defaultValue),
		shouldError: shouldError,
		errorMsg:    errorMsg,
	}
}

func (e *ErrorMockFlag) Set(value string) error {
	if e.shouldError {
		return fmt.Errorf("%s", e.errorMsg)
	}
	return e.MockFlag.Set(value)
}

// createTestContext 创建用于测试的命令上下文
func createTestContext() *types.CmdContext {
	ctx := types.NewCmdContext("test-cmd", "tc", flag.ContinueOnError)
	return ctx
}

// createTestContextWithFlags 创建带有标志的测试上下文
func createTestContextWithFlags() *types.CmdContext {
	ctx := createTestContext()

	// 添加一些测试标志
	mockFlag1 := NewMockFlag("verbose", "v", "false")
	mockFlag2 := NewMockFlag("output", "o", "stdout")
	mockFlag3 := NewMockFlag("config", "c", "")

	// 绑定环境变量
	mockFlag2.BindEnv("TEST_OUTPUT")
	mockFlag3.BindEnv("TEST_CONFIG")

	// 注册标志到 FlagSet
	ctx.FlagSet.Var(mockFlag1, "verbose", "verbose output")
	ctx.FlagSet.Var(mockFlag1, "v", "verbose output")
	ctx.FlagSet.Var(mockFlag2, "output", "output destination")
	ctx.FlagSet.Var(mockFlag2, "o", "output destination")
	ctx.FlagSet.Var(mockFlag3, "config", "config file path")
	ctx.FlagSet.Var(mockFlag3, "c", "config file path")

	return ctx
}

// TestParseArgs 测试参数解析功能
func TestParseArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		parseSubcmds bool
		setupEnv     map[string]string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "解析空参数",
			args:         []string{},
			parseSubcmds: false,
			wantErr:      false,
		},
		{
			name:         "解析有效标志",
			args:         []string{"-v", "true", "--output", "file.txt"},
			parseSubcmds: false,
			wantErr:      false,
		},
		{
			name:         "解析带环境变量的标志",
			args:         []string{"-v", "true"},
			parseSubcmds: false,
			setupEnv: map[string]string{
				"TEST_OUTPUT": "env_output.txt",
				"TEST_CONFIG": "env_config.yaml",
			},
			wantErr: false,
		},
		{
			name:         "解析非标志参数",
			args:         []string{"-v", "true", "arg1", "arg2"},
			parseSubcmds: false,
			wantErr:      false,
		},
		{
			name:         "解析子命令",
			args:         []string{"subcmd", "-v", "true"},
			parseSubcmds: true,
			wantErr:      false,
		},
		{
			name:         "解析无效标志",
			args:         []string{"--invalid-flag", "value"},
			parseSubcmds: false,
			wantErr:      true,
			errContains:  "flag provided but not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for key, value := range tt.setupEnv {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			ctx := createTestContextWithFlags()

			// 如果测试子命令，添加子命令
			if tt.parseSubcmds && len(tt.args) > 0 && tt.args[0] == "subcmd" {
				subCtx := createTestContextWithFlags()
				ctx.SubCmds = append(ctx.SubCmds, subCtx)
				ctx.SubCmdMap["subcmd"] = subCtx
			}

			err := ParseArgs(ctx, tt.args, tt.parseSubcmds)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseArgs() 期望错误但未返回错误")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseArgs() 错误信息 = %v, 期望包含 %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ParseArgs() 意外错误 = %v", err)
				}
			}
		})
	}
}

// TestParseSubCommandSafe 测试子命令解析功能
func TestParseSubCommandSafe(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setupSubCmd bool
		wantErr     bool
		errContains string
	}{
		{
			name:        "解析空参数",
			args:        []string{},
			setupSubCmd: false,
			wantErr:     false,
		},
		{
			name:        "解析存在的子命令",
			args:        []string{"subcmd", "-v", "true"},
			setupSubCmd: true,
			wantErr:     false,
		},
		{
			name:        "解析不存在的子命令",
			args:        []string{"nonexistent", "-v", "true"},
			setupSubCmd: false,
			wantErr:     false, // 不存在的子命令不会报错，只是不处理
		},
		{
			name:        "解析子命令但无剩余参数",
			args:        []string{"subcmd"},
			setupSubCmd: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContextWithFlags()

			// 设置子命令
			if tt.setupSubCmd {
				subCtx := createTestContextWithFlags()
				ctx.SubCmds = append(ctx.SubCmds, subCtx)
				ctx.SubCmdMap["subcmd"] = subCtx
			}

			err := ParseSubCommandSafe(ctx, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseSubCommandSafe() 期望错误但未返回错误")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseSubCommandSafe() 错误信息 = %v, 期望包含 %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ParseSubCommandSafe() 意外错误 = %v", err)
				}
			}
		})
	}
}

// TestLoadEnvVars 测试环境变量加载功能
func TestLoadEnvVars(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "加载有效环境变量",
			envVars: map[string]string{
				"TEST_OUTPUT": "env_output.txt",
				"TEST_CONFIG": "env_config.yaml",
			},
			wantErr: false,
		},
		{
			name:    "加载空环境变量",
			envVars: map[string]string{},
			wantErr: false,
		},
		{
			name: "加载部分环境变量",
			envVars: map[string]string{
				"TEST_OUTPUT": "env_output.txt",
			},
			wantErr: false,
		},
		{
			name: "环境变量值为空",
			envVars: map[string]string{
				"TEST_OUTPUT": "",
				"TEST_CONFIG": "config.yaml",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for key, value := range tt.envVars {
				if value != "" {
					os.Setenv(key, value)
				}
				defer os.Unsetenv(key)
			}

			ctx := createTestContextWithFlags()
			err := LoadEnvVars(ctx)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadEnvVars() 期望错误但未返回错误")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("LoadEnvVars() 错误信息 = %v, 期望包含 %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("LoadEnvVars() 意外错误 = %v", err)
				}
			}
		})
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("nil上下文测试", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("nil上下文应该导致panic")
			}
		}()
		ParseArgs(nil, []string{}, false)
	})

	t.Run("极长参数列表", func(t *testing.T) {
		ctx := createTestContextWithFlags()
		args := make([]string, 10000)
		for i := range args {
			args[i] = fmt.Sprintf("arg%d", i)
		}

		err := ParseArgs(ctx, args, false)
		if err != nil {
			t.Errorf("极长参数列表解析失败: %v", err)
		}

		if len(ctx.Args) != 10000 {
			t.Errorf("参数数量不匹配，期望 10000，实际 %d", len(ctx.Args))
		}
	})

	t.Run("特殊字符参数", func(t *testing.T) {
		ctx := createTestContextWithFlags()
		specialArgs := []string{
			"arg with spaces",
			"arg-with-dashes",
			"arg_with_underscores",
			"arg123",
			"中文参数",
			"🚀emoji",
		}

		err := ParseArgs(ctx, specialArgs, false)
		if err != nil {
			t.Errorf("特殊字符参数解析失败: %v", err)
		}
	})

	t.Run("重复环境变量处理", func(t *testing.T) {
		os.Setenv("TEST_OUTPUT", "duplicate_test")
		defer os.Unsetenv("TEST_OUTPUT")

		ctx := createTestContextWithFlags()

		// 多次调用LoadEnvVars
		err1 := LoadEnvVars(ctx)
		err2 := LoadEnvVars(ctx)

		if err1 != nil {
			t.Errorf("第一次LoadEnvVars失败: %v", err1)
		}
		if err2 != nil {
			t.Errorf("第二次LoadEnvVars失败: %v", err2)
		}
	})
}

// TestConcurrency 测试并发安全性
func TestConcurrency(t *testing.T) {
	t.Run("并发解析参数", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100
		errors := make(chan error, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				ctx := createTestContextWithFlags()
				args := []string{"-v", "true", fmt.Sprintf("arg%d", id)}
				err := ParseArgs(ctx, args, false)
				if err != nil {
					errors <- err
				}
			}(i)
		}
		wg.Wait()
		close(errors)

		for err := range errors {
			t.Errorf("并发解析参数失败: %v", err)
		}
	})

	t.Run("并发加载环境变量", func(t *testing.T) {
		os.Setenv("TEST_CONCURRENT", "concurrent_value")
		defer os.Unsetenv("TEST_CONCURRENT")

		var wg sync.WaitGroup
		numGoroutines := 50
		errors := make(chan error, numGoroutines)

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				ctx := createTestContextWithFlags()
				err := LoadEnvVars(ctx)
				if err != nil {
					errors <- err
				}
			}()
		}
		wg.Wait()
		close(errors)

		for err := range errors {
			t.Errorf("并发加载环境变量失败: %v", err)
		}
	})
}

// TestComplexScenarios 测试复杂场景
func TestComplexScenarios(t *testing.T) {
	t.Run("嵌套子命令解析", func(t *testing.T) {
		// 创建主命令
		mainCtx := createTestContextWithFlags()

		// 创建一级子命令
		subCtx1 := createTestContextWithFlags()
		mainCtx.SubCmds = append(mainCtx.SubCmds, subCtx1)
		mainCtx.SubCmdMap["sub1"] = subCtx1

		// 创建二级子命令
		subCtx2 := createTestContextWithFlags()
		subCtx1.SubCmds = append(subCtx1.SubCmds, subCtx2)
		subCtx1.SubCmdMap["sub2"] = subCtx2

		args := []string{"sub1", "sub2", "-v", "true", "final_arg"}
		err := ParseArgs(mainCtx, args, true)

		if err != nil {
			t.Errorf("嵌套子命令解析失败: %v", err)
		}
	})

	t.Run("混合标志和环境变量", func(t *testing.T) {
		// 设置环境变量
		os.Setenv("TEST_OUTPUT", "env_value")
		os.Setenv("TEST_CONFIG", "env_config")
		defer func() {
			os.Unsetenv("TEST_OUTPUT")
			os.Unsetenv("TEST_CONFIG")
		}()

		ctx := createTestContextWithFlags()
		args := []string{"-v", "true", "--config", "flag_config", "remaining_arg"}

		err := ParseArgs(ctx, args, false)
		if err != nil {
			t.Errorf("混合标志和环境变量解析失败: %v", err)
		}

		// 验证参数被正确解析
		if len(ctx.Args) == 0 {
			t.Error("期望有剩余参数")
		}
	})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	t.Run("标志解析错误", func(t *testing.T) {
		ctx := createTestContext()
		// 不注册任何标志，然后尝试解析未定义的标志
		args := []string{"--undefined-flag", "value"}

		err := ParseArgs(ctx, args, false)
		if err == nil {
			t.Error("期望解析未定义标志时返回错误")
		}
	})

	t.Run("环境变量解析错误", func(t *testing.T) {
		// 创建一个会在Set时返回错误的ErrorMockFlag
		ctx := createTestContext()

		errorFlag := NewErrorMockFlag("error-flag", "", "default", true, "模拟设置错误")
		errorFlag.BindEnv("ERROR_ENV")

		ctx.FlagSet.Var(errorFlag, "error-flag", "error flag")

		os.Setenv("ERROR_ENV", "some_value")
		defer os.Unsetenv("ERROR_ENV")

		err := LoadEnvVars(ctx)
		if err == nil {
			t.Error("期望环境变量解析错误时返回错误")
		}
	})
}

// TestPerformance 性能测试
func TestPerformance(t *testing.T) {
	t.Run("大量标志解析性能", func(t *testing.T) {
		ctx := createTestContext()

		// 创建大量标志
		numFlags := 1000
		args := make([]string, 0, numFlags*2)

		for i := 0; i < numFlags; i++ {
			flagName := fmt.Sprintf("flag%d", i)
			mockFlag := NewMockFlag(flagName, "", "default")
			ctx.FlagSet.Var(mockFlag, flagName, "test flag")

			args = append(args, fmt.Sprintf("--%s", flagName), fmt.Sprintf("value%d", i))
		}

		err := ParseArgs(ctx, args, false)
		if err != nil {
			t.Errorf("大量标志解析失败: %v", err)
		}
	})

	t.Run("大量环境变量加载性能", func(t *testing.T) {
		ctx := createTestContext()

		// 创建大量带环境变量的标志
		numFlags := 500

		for i := 0; i < numFlags; i++ {
			flagName := fmt.Sprintf("envflag%d", i)
			envName := fmt.Sprintf("TEST_ENV_%d", i)

			mockFlag := NewMockFlag(flagName, "", "default")
			mockFlag.BindEnv(envName)
			ctx.FlagSet.Var(mockFlag, flagName, "test env flag")

			os.Setenv(envName, fmt.Sprintf("env_value_%d", i))
			defer os.Unsetenv(envName)
		}

		err := LoadEnvVars(ctx)
		if err != nil {
			t.Errorf("大量环境变量加载失败: %v", err)
		}
	})
}

// BenchmarkParseArgs 基准测试参数解析
func BenchmarkParseArgs(b *testing.B) {
	ctx := createTestContextWithFlags()
	args := []string{"-v", "true", "--output", "file.txt", "arg1", "arg2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 重置上下文状态
		ctx.Args = []string{}
		ctx.Parsed.Store(false)

		ParseArgs(ctx, args, false)
	}
}

// BenchmarkLoadEnvVars 基准测试环境变量加载
func BenchmarkLoadEnvVars(b *testing.B) {
	os.Setenv("TEST_OUTPUT", "benchmark_output")
	os.Setenv("TEST_CONFIG", "benchmark_config")
	defer func() {
		os.Unsetenv("TEST_OUTPUT")
		os.Unsetenv("TEST_CONFIG")
	}()

	ctx := createTestContextWithFlags()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadEnvVars(ctx)
	}
}

// BenchmarkParseSubCommand 基准测试子命令解析
func BenchmarkParseSubCommand(b *testing.B) {
	ctx := createTestContextWithFlags()
	subCtx := createTestContextWithFlags()
	ctx.SubCmds = append(ctx.SubCmds, subCtx)
	ctx.SubCmdMap["subcmd"] = subCtx

	args := []string{"subcmd", "-v", "true"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseSubCommandSafe(ctx, args)
	}
}

// BenchmarkConcurrentParsing 基准测试并发解析
func BenchmarkConcurrentParsing(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := createTestContextWithFlags()
			args := []string{"-v", "true", "--output", "file.txt"}
			ParseArgs(ctx, args, false)
		}
	})
}
