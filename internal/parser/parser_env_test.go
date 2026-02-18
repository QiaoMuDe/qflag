package parser

import (
	"os"
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

/*
本测试文件验证了环境变量加载功能的各种场景，包括：

1. 基本环境变量加载验证：
   - 验证环境变量能正确加载到标志中
   - 测试场景：设置环境变量，验证标志值正确加载
   - 确保环境变量优先级低于命令行参数

2. 环境变量前缀处理验证：
   - 验证带前缀的环境变量名处理
   - 测试场景：使用不同前缀，验证环境变量名组合正确
   - 确保前缀处理逻辑正确

3. 环境变量不存在处理验证：
   - 验证环境变量不存在时的处理
   - 测试场景：未设置环境变量，验证标志保持默认值
   - 确保不会因为环境变量不存在而出错

4. 命令行参数优先级验证：
   - 验证命令行参数优先级高于环境变量
   - 测试场景：同时设置命令行参数和环境变量，验证使用命令行参数值
   - 确保优先级逻辑正确

5. 无环境变量绑定标志处理验证：
   - 验证未绑定环境变量的标志不会被处理
   - 测试场景：标志未设置环境变量名，验证跳过处理
   - 确保只处理绑定了环境变量的标志

测试覆盖了所有主要功能路径，确保环境变量加载逻辑的正确性和完整性。
*/

// 创建带有环境变量的字符串标志
func createMockStringFlag(name, short, desc, defaultValue, envVar string) *mock.MockFlag {
	flag := mock.NewMockFlag(name, short, desc, types.FlagTypeString, defaultValue)
	flag.BindEnv(envVar)
	return flag
}

// 创建带有环境变量的整数标志
func createMockIntFlag(name, short, desc string, defaultValue int64, envVar string) *mock.MockFlag {
	flag := mock.NewMockFlag(name, short, desc, types.FlagTypeInt, defaultValue)
	flag.BindEnv(envVar)
	return flag
}

// TestLoadEnvVars 测试环境变量加载功能
func TestLoadEnvVars(t *testing.T) {
	tests := []struct {
		name           string
		flags          []types.Flag
		envVars        map[string]string
		envPrefix      string
		args           []string
		expectedValues map[string]interface{}
		expectError    bool
	}{
		{
			name: "基本环境变量加载",
			flags: []types.Flag{
				createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
				createMockIntFlag("count", "c", "计数", 0, "COUNT"),
			},
			envVars: map[string]string{
				"INPUT_FILE": "test.txt",
				"COUNT":      "42",
			},
			envPrefix: "",
			args:      []string{},
			expectedValues: map[string]interface{}{
				"input": "test.txt",
				"count": "42",
			},
			expectError: false,
		},
		{
			name: "带前缀的环境变量加载",
			flags: []types.Flag{
				createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
				createMockIntFlag("count", "c", "计数", 0, "COUNT"),
			},
			envVars: map[string]string{
				"APP_INPUT_FILE": "test.txt",
				"APP_COUNT":      "42",
			},
			envPrefix: "APP_",
			args:      []string{},
			expectedValues: map[string]interface{}{
				"input": "test.txt",
				"count": "42",
			},
			expectError: false,
		},
		{
			name: "命令行参数优先级高于环境变量",
			flags: []types.Flag{
				createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
				createMockIntFlag("count", "c", "计数", 0, "COUNT"),
			},
			envVars: map[string]string{
				"INPUT_FILE": "env.txt",
				"COUNT":      "100",
			},
			envPrefix: "",
			args:      []string{"--input", "cmd.txt", "--count", "50"},
			expectedValues: map[string]interface{}{
				"input": "cmd.txt",
				"count": "50",
			},
			expectError: false,
		},
		{
			name: "环境变量不存在",
			flags: []types.Flag{
				createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
				createMockIntFlag("count", "c", "计数", 0, "COUNT"),
			},
			envVars:   map[string]string{},
			envPrefix: "",
			args:      []string{},
			expectedValues: map[string]interface{}{
				"input": "",
				"count": "",
			},
			expectError: false,
		},
		{
			name: "未绑定环境变量的标志",
			flags: []types.Flag{
				createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
				mock.NewMockFlag("output", "o", "输出文件", types.FlagTypeString, ""), // 未绑定环境变量
			},
			envVars: map[string]string{
				"INPUT_FILE": "test.txt",
			},
			envPrefix: "",
			args:      []string{},
			expectedValues: map[string]interface{}{
				"input":  "test.txt",
				"output": "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// 清理环境变量
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// 创建模拟命令
			cmd := mock.NewMockCommand("test", "t", "Test command")

			// 添加标志到命令
			for _, flag := range tt.flags {
				if err := cmd.AddFlag(flag); err != nil {
					t.Errorf("添加标志失败: %v", err)
					return
				}
			}

			// 设置环境变量前缀
			if tt.envPrefix != "" {
				cmd.SetEnvPrefix(tt.envPrefix)
			}

			// 创建解析器
			parser := NewDefaultParser(types.ContinueOnError)

			// 使用 Parse 方法来测试环境变量加载
			err := parser.Parse(cmd, tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误，但没有得到错误")
				}
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("期望没有错误，但得到: %v", err)
				return
			}

			// 验证标志值
			for _, flag := range tt.flags {
				flagName := flag.Name()
				expectedValue, exists := tt.expectedValues[flagName]
				if !exists {
					continue
				}

				// 获取标志的实际值
				actualValue := flag.GetStr()

				if actualValue != expectedValue {
					t.Errorf("标志 '%s' 的值不正确，期望: %v，实际: %v", flagName, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestSingleFlagEnv 测试单个标志的环境变量加载
func TestSingleFlagEnv(t *testing.T) {
	tests := []struct {
		name          string
		flag          types.Flag
		envVar        string
		envValue      string
		envPrefix     string
		expectedValue string
		expectError   bool
	}{
		{
			name:          "无环境变量绑定",
			flag:          mock.NewMockFlag("input", "i", "输入文件", types.FlagTypeString, ""),
			envVar:        "",
			envValue:      "",
			envPrefix:     "",
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "环境变量存在",
			flag:          createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
			envVar:        "INPUT_FILE",
			envValue:      "test.txt",
			envPrefix:     "",
			expectedValue: "test.txt",
			expectError:   false,
		},
		{
			name:          "环境变量不存在",
			flag:          createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
			envVar:        "INPUT_FILE",
			envValue:      "",
			envPrefix:     "",
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "带前缀的环境变量",
			flag:          createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
			envVar:        "APP_INPUT_FILE",
			envValue:      "test.txt",
			envPrefix:     "APP_",
			expectedValue: "test.txt",
			expectError:   false,
		},
		{
			name:          "标志已被命令行参数设置",
			flag:          createMockStringFlag("input", "i", "输入文件", "", "INPUT_FILE"),
			envVar:        "INPUT_FILE",
			envValue:      "env_value",
			envPrefix:     "",
			expectedValue: "cmd_value", // 应该保持命令行参数的值
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.envValue != "" {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			}

			// 创建模拟命令
			cmd := mock.NewMockCommand("test", "t", "Test command")
			if err := cmd.AddFlag(tt.flag); err != nil {
				t.Errorf("添加标志失败: %v", err)
				return
			}

			// 设置环境变量前缀
			if tt.envPrefix != "" {
				cmd.SetEnvPrefix(tt.envPrefix)
			}

			// 创建解析器
			parser := NewDefaultParser(types.ContinueOnError)

			// 如果测试命令行参数优先级，先设置标志值
			var args []string
			if tt.name == "标志已被命令行参数设置" {
				args = []string{"--input", "cmd_value"}
			}

			// 使用 Parse 方法来测试环境变量加载
			err := parser.Parse(cmd, args)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误，但没有得到错误")
				}
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("期望没有错误，但得到: %v", err)
				return
			}

			// 验证标志值
			actualValue := tt.flag.GetStr()

			if actualValue != tt.expectedValue {
				t.Errorf("标志值不正确，期望: %v，实际: %v", tt.expectedValue, actualValue)
			}
		})
	}
}
