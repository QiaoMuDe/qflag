// Package completion Bash补全测试
// 本文件包含了Bash自动补全功能的单元测试，测试Bash补全脚本
// 生成和标志、子命令的智能补全功能的正确性。
package completion

import (
	"bytes"
	"strings"
	"testing"
)

// TestGenerateBashCommandTreeEntry 测试生成Bash命令树条目
func TestGenerateBashCommandTreeEntry(t *testing.T) {
	tests := []struct {
		name     string
		cmdPath  string
		cmdOpts  []string
		expected string
	}{
		{
			name:     "根命令条目",
			cmdPath:  "/",
			cmdOpts:  []string{"--help", "-h", "start", "stop"},
			expected: `{{.ProgramName}}_cmd_tree[/]="--help|-h|start|stop"`,
		},
		{
			name:     "子命令条目",
			cmdPath:  "/start/",
			cmdOpts:  []string{"--verbose", "-v", "--config", "-c"},
			expected: `{{.ProgramName}}_cmd_tree[/start/]="--verbose|-v|--config|-c"`,
		},
		{
			name:     "空选项",
			cmdPath:  "/empty/",
			cmdOpts:  []string{},
			expected: `{{.ProgramName}}_cmd_tree[/empty/]=""`,
		},
		{
			name:     "单个选项",
			cmdPath:  "/single/",
			cmdOpts:  []string{"--only"},
			expected: `{{.ProgramName}}_cmd_tree[/single/]="--only"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			generateBashCommandTreeEntry(&buf, tt.cmdPath, tt.cmdOpts, "testprog")

			result := buf.String()
			// 更新期望值以匹配实际的程序名
			expectedWithProgName := strings.ReplaceAll(tt.expected, "{{.ProgramName}}", "testprog")
			if !strings.Contains(result, expectedWithProgName) {
				t.Errorf("generateBashCommandTreeEntry() = %q, 期望包含 %q", result, expectedWithProgName)
			}
		})
	}
}

// TestGenerateBashCompletion 测试生成Bash补全脚本
func TestGenerateBashCompletion(t *testing.T) {
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--verbose",
			Type:        "none",
			ValueType:   "bool",
		},
		{
			CommandPath: "/",
			Name:        "--mode",
			Type:        "required",
			ValueType:   "enum",
			EnumOptions: []string{"debug", "release", "test"},
		},
		{
			CommandPath: "/start/",
			Name:        "--config",
			Type:        "required",
			ValueType:   "string",
		},
	}

	rootCmdOpts := []string{"--help", "-h", "start", "stop"}
	cmdTreeEntries := `myapp_cmd_tree[/start/]="--config|-c"`
	programName := "myapp"

	var buf bytes.Buffer
	generateBashCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)

	result := buf.String()

	// 调试：输出实际生成的脚本内容
	t.Logf("实际生成的脚本内容:\n%s", result)

	// 检查必要的组件
	expectedComponents := []string{
		"#!/usr/bin/env bash",
		"declare -A myapp_cmd_tree",
		"declare -A myapp_flag_params",
		"declare -A myapp_enum_options",
		"_myapp()",
		"complete -F _myapp myapp",
		`myapp_cmd_tree[/]="--help|-h|start|stop"`,
		"myapp_cmd_tree[/start/]=\"--config|-c\"",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(result, component) {
			t.Errorf("Bash补全脚本不包含必要组件: %s", component)
		}
	}

	// 检查标志参数
	expectedFlags := []string{
		`myapp_flag_params["/|--verbose"]="none|bool"`,
		`myapp_flag_params["/|--mode"]="required|enum"`,
		`myapp_flag_params["/start/|--config"]="required|string"`,
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(result, flag) {
			t.Errorf("Bash补全脚本不包含标志参数: %s", flag)
		}
	}

	// 检查枚举选项 - 适配新的格式（使用|分隔符）
	expectedEnum := `myapp_enum_options["/|--mode"]="debug|release|test"`
	if !strings.Contains(result, expectedEnum) {
		t.Errorf("Bash补全脚本不包含枚举选项: %s", expectedEnum)
	}
}

// TestEscapeSpecialChars 测试特殊字符转义
func TestEscapeSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "普通字符串",
			input:    "normal",
			expected: "normal",
		},
		{
			name:     "包含反斜杠",
			input:    "path\\to\\file",
			expected: "path\\\\to\\\\file",
		},
		{
			name:     "包含双引号",
			input:    `say "hello"`,
			expected: `say\ \"hello\"`,
		},
		{
			name:     "包含空格",
			input:    "hello world",
			expected: "hello\\ world",
		},
		{
			name:     "混合特殊字符",
			input:    `path\to "my file"`,
			expected: `path\\to\ \"my\ file\"`,
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有特殊字符",
			input:    `\" `,
			expected: `\\\"\ `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("escapeSpecialChars(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestBashTemplateConstants 测试Bash模板常量
func TestBashTemplateConstants(t *testing.T) {
	// 测试标志参数项格式
	if BashFlagParamItem != "{{.ProgramName}}_flag_params[%q]=%q\n" {
		t.Errorf("BashFlagParamItem 格式错误: %s", BashFlagParamItem)
	}

	// 测试枚举选项格式
	if BashEnumOptions != "{{.ProgramName}}_enum_options[%q]=%q\n" {
		t.Errorf("BashEnumOptions 格式错误: %s", BashEnumOptions)
	}

	// 测试函数头部模板
	if !strings.Contains(BashFunctionHeader, "#!/usr/bin/env bash") {
		t.Error("BashFunctionHeader 不包含shebang")
	}

	if !strings.Contains(BashFunctionHeader, "_{{.ProgramName}}()") {
		t.Error("BashFunctionHeader 不包含函数定义")
	}

	if !strings.Contains(BashFunctionHeader, "complete -F _{{.ProgramName}} {{.ProgramName}}") {
		t.Error("BashFunctionHeader 不包含complete命令")
	}
}

// TestBashCompletionWithComplexScenario 测试复杂场景的Bash补全
func TestBashCompletionWithComplexScenario(t *testing.T) {
	// 创建复杂的参数场景
	params := []FlagParam{
		// 根命令标志
		{CommandPath: "/", Name: "--help", Type: "none", ValueType: "bool"},
		{CommandPath: "/", Name: "-h", Type: "none", ValueType: "bool"},
		{CommandPath: "/", Name: "--verbose", Type: "none", ValueType: "bool"},
		{CommandPath: "/", Name: "-v", Type: "none", ValueType: "bool"},
		{CommandPath: "/", Name: "--config", Type: "required", ValueType: "string"},
		{CommandPath: "/", Name: "-c", Type: "required", ValueType: "string"},
		{CommandPath: "/", Name: "--mode", Type: "required", ValueType: "enum", EnumOptions: []string{"dev", "prod", "test"}},
		{CommandPath: "/", Name: "-m", Type: "required", ValueType: "enum", EnumOptions: []string{"dev", "prod", "test"}},

		// 子命令标志
		{CommandPath: "/start/", Name: "--port", Type: "required", ValueType: "string"},
		{CommandPath: "/start/", Name: "-p", Type: "required", ValueType: "string"},
		{CommandPath: "/start/", Name: "--daemon", Type: "none", ValueType: "bool"},
		{CommandPath: "/start/", Name: "-d", Type: "none", ValueType: "bool"},

		// 深层子命令标志
		{CommandPath: "/config/set/", Name: "--key", Type: "required", ValueType: "string"},
		{CommandPath: "/config/set/", Name: "--value", Type: "required", ValueType: "string"},
	}

	rootCmdOpts := []string{"--help", "-h", "--verbose", "-v", "--config", "-c", "--mode", "-m", "start", "stop", "config"}

	cmdTreeEntries := `complexapp_cmd_tree[/start/]="--port|-p|--daemon|-d"
complexapp_cmd_tree[/stop/]="--force|-f"
complexapp_cmd_tree[/config/]="set|get|list"
complexapp_cmd_tree[/config/set/]="--key|--value"`

	programName := "complexapp"

	var buf bytes.Buffer
	generateBashCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)

	result := buf.String()

	// 验证复杂场景的各个部分
	tests := []struct {
		name     string
		contains string
	}{
		{"程序名称", "_complexapp()"},
		{"完成命令", "complete -F _complexapp complexapp"},
		{"根命令选项", `complexapp_cmd_tree[/]="--help|-h|--verbose|-v|--config|-c|--mode|-m|start|stop|config"`},
		{"子命令树", `complexapp_cmd_tree[/start/]="--port|-p|--daemon|-d"`},
		{"深层命令树", `complexapp_cmd_tree[/config/set/]="--key|--value"`},
		{"根命令枚举", `complexapp_enum_options["/|--mode"]="dev|prod|test"`},
		{"根命令标志", `complexapp_flag_params["/|--config"]="required|string"`},
		{"子命令标志", `complexapp_flag_params["/start/|--port"]="required|string"`},
		{"深层标志", `complexapp_flag_params["/config/set/|--key"]="required|string"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(result, tt.contains) {
				t.Errorf("复杂Bash补全脚本不包含: %s", tt.contains)
			}
		})
	}
}

// BenchmarkGenerateBashCompletion 基准测试Bash补全生成
func BenchmarkGenerateBashCompletion(b *testing.B) {
	// 准备测试数据
	params := make([]FlagParam, 50)
	for i := 0; i < 50; i++ {
		params[i] = FlagParam{
			CommandPath: "/",
			Name:        "--flag" + string(rune('a'+i%26)),
			Type:        "required",
			ValueType:   "string",
		}
	}

	rootCmdOpts := make([]string, 20)
	for i := 0; i < 20; i++ {
		rootCmdOpts[i] = "--opt" + string(rune('a'+i%26))
	}

	cmdTreeEntries := "cmd_tree[/sub1/]=\"--flag1|-f1\"\ncmd_tree[/sub2/]=\"--flag2|-f2\""
	programName := "benchapp"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		generateBashCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)
	}
}

// BenchmarkEscapeSpecialChars 基准测试特殊字符转义
func BenchmarkEscapeSpecialChars(b *testing.B) {
	testStrings := []string{
		"normal string",
		`path\to\file with "quotes" and spaces`,
		"simple",
		`complex\"string\\with\many\"special\chars`,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_ = escapeSpecialChars(s)
		}
	}
}
