// Package completion 内部补全测试
// 本文件包含了补全系统内部功能的单元测试，测试补全算法、
// 匹配策略等核心功能的底层支持的正确性。
package completion

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestGetValueTypeByFlagType 测试根据标志类型获取值类型
func TestGetValueTypeByFlagType(t *testing.T) {
	tests := []struct {
		name     string
		flagType flags.FlagType
		expected string
	}{
		{
			name:     "布尔类型",
			flagType: flags.FlagTypeBool,
			expected: "bool",
		},
		{
			name:     "枚举类型",
			flagType: flags.FlagTypeEnum,
			expected: "enum",
		},
		{
			name:     "字符串类型",
			flagType: flags.FlagTypeString,
			expected: "string",
		},
		{
			name:     "整数类型",
			flagType: flags.FlagTypeInt,
			expected: "string",
		},
		{
			name:     "浮点数类型",
			flagType: flags.FlagTypeFloat64,
			expected: "string",
		},
		{
			name:     "时间间隔类型",
			flagType: flags.FlagTypeDuration,
			expected: "string",
		},
		{
			name:     "切片类型",
			flagType: flags.FlagTypeSlice,
			expected: "string",
		},
		{
			name:     "时间类型",
			flagType: flags.FlagTypeTime,
			expected: "string",
		},
		{
			name:     "映射类型",
			flagType: flags.FlagTypeMap,
			expected: "string",
		},
		{
			name:     "未知类型",
			flagType: flags.FlagTypeUnknown,
			expected: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getValueTypeByFlagType(tt.flagType)
			if result != tt.expected {
				t.Errorf("getValueTypeByFlagType(%v) = %v, 期望 %v", tt.flagType, result, tt.expected)
			}
		})
	}
}

// TestGetParamTypeByFlagType 测试根据标志类型获取参数需求类型
func TestGetParamTypeByFlagType(t *testing.T) {
	tests := []struct {
		name     string
		flagType flags.FlagType
		expected string
	}{
		{
			name:     "布尔类型",
			flagType: flags.FlagTypeBool,
			expected: "none",
		},
		{
			name:     "字符串类型",
			flagType: flags.FlagTypeString,
			expected: "required",
		},
		{
			name:     "整数类型",
			flagType: flags.FlagTypeInt,
			expected: "required",
		},
		{
			name:     "枚举类型",
			flagType: flags.FlagTypeEnum,
			expected: "required",
		},
		{
			name:     "浮点数类型",
			flagType: flags.FlagTypeFloat64,
			expected: "required",
		},
		{
			name:     "时间间隔类型",
			flagType: flags.FlagTypeDuration,
			expected: "required",
		},
		{
			name:     "切片类型",
			flagType: flags.FlagTypeSlice,
			expected: "required",
		},
		{
			name:     "时间类型",
			flagType: flags.FlagTypeTime,
			expected: "required",
		},
		{
			name:     "映射类型",
			flagType: flags.FlagTypeMap,
			expected: "required",
		},
		{
			name:     "未知类型",
			flagType: flags.FlagTypeUnknown,
			expected: "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getParamTypeByFlagType(tt.flagType)
			if result != tt.expected {
				t.Errorf("getParamTypeByFlagType(%v) = %v, 期望 %v", tt.flagType, result, tt.expected)
			}
		})
	}
}

// TestInternalValidateCompletionGeneration 测试补全生成验证（内部版本）
func TestInternalValidateCompletionGeneration(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *types.CmdContext
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil上下文",
			ctx:     nil,
			wantErr: true,
			errMsg:  "command instance is nil",
		},
		{
			name: "非根命令",
			ctx: func() *types.CmdContext {
				ctx := types.NewCmdContext("sub", "", flag.ContinueOnError)
				parent := types.NewCmdContext("parent", "", flag.ContinueOnError)
				ctx.Parent = parent
				return ctx
			}(),
			wantErr: true,
			errMsg:  "invalid command state: not a root command",
		},
		{
			name: "标志注册表为nil",
			ctx: func() *types.CmdContext {
				ctx := types.NewCmdContext("test", "", flag.ContinueOnError)
				ctx.FlagRegistry = nil
				return ctx
			}(),
			wantErr: true,
			errMsg:  "invalid command state: flag registry is nil",
		},
		{
			name:    "有效的根命令",
			ctx:     types.NewCmdContext("valid", "", flag.ContinueOnError),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCompletionGeneration(tt.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateCompletionGeneration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("validateCompletionGeneration() error = %v, 期望错误信息包含 %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestInternalValidateCompletionGenerationEdgeCases 测试补全生成验证的边界情况（内部版本）
func TestInternalValidateCompletionGenerationEdgeCases(t *testing.T) {
	// 测试有效的根命令但没有标志
	t.Run("有效根命令无标志", func(t *testing.T) {
		ctx := types.NewCmdContext("empty", "", flag.ContinueOnError)
		err := validateCompletionGeneration(ctx)
		if err != nil {
			t.Errorf("validateCompletionGeneration() 对于空标志的根命令应该成功, 但得到错误: %v", err)
		}
	})

	// 测试有子命令的根命令
	t.Run("有子命令的根命令", func(t *testing.T) {
		rootCtx := types.NewCmdContext("root", "", flag.ContinueOnError)
		subCtx := types.NewCmdContext("sub", "", flag.ContinueOnError)
		subCtx.Parent = rootCtx
		rootCtx.SubCmds = []*types.CmdContext{subCtx}

		err := validateCompletionGeneration(rootCtx)
		if err != nil {
			t.Errorf("validateCompletionGeneration() 对于有子命令的根命令应该成功, 但得到错误: %v", err)
		}
	})
}

// BenchmarkGetValueTypeByFlagType 基准测试获取值类型
func BenchmarkGetValueTypeByFlagType(b *testing.B) {
	flagTypes := []flags.FlagType{
		flags.FlagTypeBool,
		flags.FlagTypeString,
		flags.FlagTypeInt,
		flags.FlagTypeEnum,
		flags.FlagTypeFloat64,
		flags.FlagTypeDuration,
		flags.FlagTypeSlice,
		flags.FlagTypeTime,
		flags.FlagTypeMap,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ft := range flagTypes {
			_ = getValueTypeByFlagType(ft)
		}
	}
}

// BenchmarkGetParamTypeByFlagType 基准测试获取参数类型
func BenchmarkGetParamTypeByFlagType(b *testing.B) {
	flagTypes := []flags.FlagType{
		flags.FlagTypeBool,
		flags.FlagTypeString,
		flags.FlagTypeInt,
		flags.FlagTypeEnum,
		flags.FlagTypeFloat64,
		flags.FlagTypeDuration,
		flags.FlagTypeSlice,
		flags.FlagTypeTime,
		flags.FlagTypeMap,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ft := range flagTypes {
			_ = getParamTypeByFlagType(ft)
		}
	}
}

// BenchmarkInternalValidateCompletionGeneration 基准测试补全生成验证（内部版本）
func BenchmarkInternalValidateCompletionGeneration(b *testing.B) {
	ctx := types.NewCmdContext("benchmark", "", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateCompletionGeneration(ctx)
	}
}
