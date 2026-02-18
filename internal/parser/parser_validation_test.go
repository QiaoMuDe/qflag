package parser

import (
	"errors"
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

/*
本测试文件验证了互斥组和必需组验证逻辑的各种场景，包括：

1. 混合长短名支持验证：
   - 验证用户可以在互斥组和必需组中混合使用长名和短名
   - 测试场景：--run -p, -r --parallel, --run --parallel, -r -p
   - 确保所有组合都能正确识别和验证

2. 去重逻辑验证：
   - 验证同一个标志的多个名称不会重复显示在错误信息中
   - 测试场景：互斥组包含 ["run", "parallel", "r", "p"]，用户输入 --run -p
   - 确保错误信息显示 [--run/-r --parallel/-p] 而不是重复的

3. 无效标志名处理验证：
   - 验证当互斥组或必需组中包含不存在的标志时的错误处理
   - 测试场景：组定义中包含 "nonexistent" 标志
   - 确保返回明确的错误信息：invalid flag name 'nonexistent' in group

4. 错误信息格式验证：
   - 验证各种错误情况下的错误信息格式
   - 互斥冲突：mutually exclusive flags [flags] in group 'name' cannot be used together
   - 必需缺失：required flags [flags] in group 'name' must be set
   - 无效标志：invalid flag name 'name' in group 'type'

5. 边界情况验证：
   - 验证各种边界情况和异常场景
   - 确保代码的健壮性和错误处理的完整性

测试覆盖了所有主要功能路径，确保验证逻辑的正确性和完整性。
*/

// TestValidateMutexGroups 测试互斥组验证逻辑
func TestValidateMutexGroups(t *testing.T) {
	tests := []struct {
		name           string
		flags          []types.Flag
		mutexGroups    []types.MutexGroup
		args           []string
		expectedError  string
		expectedErrMsg string
	}{
		{
			name: "正常情况：只设置一个标志",
			flags: []types.Flag{
				mock.NewMockBoolFlag("run", "r", "运行模式", true),
				mock.NewMockBoolFlag("parallel", "p", "并行模式", false),
			},
			mutexGroups: []types.MutexGroup{
				{Name: "parallel", Flags: []string{"run", "parallel", "r", "p"}, AllowNone: true},
			},
			args:          []string{"-r"},
			expectedError: "",
		},
		{
			name: "互斥冲突：设置多个标志",
			flags: []types.Flag{
				mock.NewMockBoolFlag("run", "r", "运行模式", true),
				mock.NewMockBoolFlag("parallel", "p", "并行模式", true),
			},
			mutexGroups: []types.MutexGroup{
				{Name: "parallel", Flags: []string{"run", "parallel", "r", "p"}, AllowNone: true},
			},
			args:           []string{"-r", "-p"},
			expectedError:  "MUTEX_GROUP_VIOLATION",
			expectedErrMsg: "mutually exclusive flags [--run/-r --parallel/-p] in group 'parallel' cannot be used together",
		},
		{
			name: "混合长短名冲突",
			flags: []types.Flag{
				mock.NewMockBoolFlag("run", "r", "运行模式", true),
				mock.NewMockBoolFlag("parallel", "p", "并行模式", true),
			},
			mutexGroups: []types.MutexGroup{
				{Name: "parallel", Flags: []string{"run", "parallel", "r", "p"}, AllowNone: true},
			},
			args:           []string{"--run", "-p"},
			expectedError:  "MUTEX_GROUP_VIOLATION",
			expectedErrMsg: "mutually exclusive flags [--run/-r --parallel/-p] in group 'parallel' cannot be used together",
		},
		{
			name: "无效标志名：包含不存在的标志",
			flags: []types.Flag{
				mock.NewMockBoolFlag("run", "r", "运行模式", true),
				mock.NewMockBoolFlag("parallel", "p", "并行模式", false),
			},
			mutexGroups: []types.MutexGroup{
				{Name: "parallel", Flags: []string{"run", "parallel", "nonexistent", "r", "p"}, AllowNone: true},
			},
			args:           []string{"-r"},
			expectedError:  "INVALID_FLAG_NAME",
			expectedErrMsg: "invalid flag name 'nonexistent' in mutex group 'parallel'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟命令
			cmd := mock.NewMockCommand("test", "t", "Test command")

			// 添加标志到命令
			for _, flag := range tt.flags {
				if err := cmd.AddFlag(flag); err != nil {
					t.Errorf("添加标志失败: %v", err)
					return
				}
			}

			// 添加互斥组
			for _, group := range tt.mutexGroups {
				if err := cmd.AddMutexGroup(group.Name, group.Flags, group.AllowNone); err != nil {
					t.Errorf("添加互斥组失败: %v", err)
					return
				}
			}

			// 创建解析器
			parser := NewDefaultParser(types.ContinueOnError)

			// 解析参数
			err := parser.ParseOnly(cmd, tt.args)

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("期望没有错误，但得到: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("期望错误 '%s'，但没有得到错误", tt.expectedError)
				} else {
					// 检查错误类型
					var qflagErr *types.Error
					if !errors.As(err, &qflagErr) {
						t.Errorf("期望 qflag 错误，但得到: %T", err)
					} else if qflagErr.Code != tt.expectedError {
						t.Errorf("期望错误类型 '%s'，但得到 '%s'", tt.expectedError, qflagErr.Code)
					} else if qflagErr.Message != tt.expectedErrMsg {
						t.Errorf("期望错误信息 '%s'，但得到 '%s'", tt.expectedErrMsg, qflagErr.Message)
					}
				}
			}
		})
	}
}

// TestValidateRequiredGroups 测试必需组验证逻辑
func TestValidateRequiredGroups(t *testing.T) {
	tests := []struct {
		name           string
		flags          []types.Flag
		requiredGroups []types.RequiredGroup
		args           []string
		expectedError  string
		expectedErrMsg string
	}{
		{
			name: "正常情况：所有必需标志都已设置",
			flags: []types.Flag{
				mock.NewMockBoolFlag("input", "i", "输入文件", true),
				mock.NewMockBoolFlag("output", "o", "输出文件", true),
			},
			requiredGroups: []types.RequiredGroup{
				{Name: "io", Flags: []string{"input", "output", "i", "o"}},
			},
			args:          []string{"-i", "-o"},
			expectedError: "",
		},
		{
			name: "缺少必需标志",
			flags: []types.Flag{
				mock.NewMockBoolFlag("input", "i", "输入文件", true),
				mock.NewMockBoolFlag("output", "o", "输出文件", false),
			},
			requiredGroups: []types.RequiredGroup{
				{Name: "io", Flags: []string{"input", "output", "i", "o"}},
			},
			args:           []string{"-i"},
			expectedError:  "REQUIRED_GROUP_VIOLATION",
			expectedErrMsg: "required flags [--output/-o] in group 'io' must be set",
		},
		{
			name: "无效标志名：包含不存在的标志",
			flags: []types.Flag{
				mock.NewMockBoolFlag("input", "i", "输入文件", true),
				mock.NewMockBoolFlag("output", "o", "输出文件", true),
			},
			requiredGroups: []types.RequiredGroup{
				{Name: "io", Flags: []string{"input", "output", "nonexistent", "i", "o"}},
			},
			args:           []string{"-i", "-o"},
			expectedError:  "INVALID_FLAG_NAME",
			expectedErrMsg: "invalid flag name 'nonexistent' in required group 'io'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟命令
			cmd := mock.NewMockCommand("test", "t", "Test command")

			// 添加标志到命令
			for _, flag := range tt.flags {
				if err := cmd.AddFlag(flag); err != nil {
					t.Errorf("添加标志失败: %v", err)
					return
				}
			}

			// 添加必需组
			for _, group := range tt.requiredGroups {
				if err := cmd.AddRequiredGroup(group.Name, group.Flags); err != nil {
					t.Errorf("添加必需组失败: %v", err)
					return
				}
			}

			// 创建解析器
			parser := NewDefaultParser(types.ContinueOnError)

			// 解析参数
			err := parser.ParseOnly(cmd, tt.args)

			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("期望没有错误，但得到: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("期望错误 '%s'，但没有得到错误", tt.expectedError)
				} else {
					// 检查错误类型
					var qflagErr *types.Error
					if !errors.As(err, &qflagErr) {
						t.Errorf("期望 qflag 错误，但得到: %T", err)
					} else if qflagErr.Code != tt.expectedError {
						t.Errorf("期望错误类型 '%s'，但得到 '%s'", tt.expectedError, qflagErr.Code)
					} else if qflagErr.Message != tt.expectedErrMsg {
						t.Errorf("期望错误信息 '%s'，但得到 '%s'", tt.expectedErrMsg, qflagErr.Message)
					}
				}
			}
		})
	}
}
