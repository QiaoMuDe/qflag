package completion

import (
	"bytes"
	"strings"
	"testing"
)

// TestGeneratePwshCommandTreeEntry 测试生成PowerShell命令树条目
func TestGeneratePwshCommandTreeEntry(t *testing.T) {
	tests := []struct {
		name     string
		cmdPath  string
		cmdOpts  []string
		expected []string // 期望包含的字符串
	}{
		{
			name:    "根命令条目",
			cmdPath: "/",
			cmdOpts: []string{"--help", "-h", "start", "stop"},
			expected: []string{
				`Context = "/"`,
				`'--help'`,
				`'-h'`,
				`'start'`,
				`'stop'`,
			},
		},
		{
			name:    "子命令条目",
			cmdPath: "/start/",
			cmdOpts: []string{"--verbose", "-v", "--config", "-c"},
			expected: []string{
				`Context = "/start/"`,
				`'--verbose'`,
				`'-v'`,
				`'--config'`,
				`'-c'`,
			},
		},
		{
			name:    "单个选项",
			cmdPath: "/single/",
			cmdOpts: []string{"--only"},
			expected: []string{
				`Context = "/single/"`,
				`'--only'`,
			},
		},
		{
			name:    "包含特殊字符的选项",
			cmdPath: "/special/",
			cmdOpts: []string{"--path='C:\\Program Files'", "--quote=\"test\""},
			expected: []string{
				"Context = \"/special/\"",
				"'--path=''C:\\\\Program Files'''",
				"'--quote=\"test\"'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			generatePwshCommandTreeEntry(&buf, tt.cmdPath, tt.cmdOpts)

			result := buf.String()
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("generatePwshCommandTreeEntry() = %q, 期望包含 %q", result, expected)
				}
			}
		})
	}
}

// TestGeneratePwshCompletion 测试生成PowerShell补全脚本
func TestGeneratePwshCompletion(t *testing.T) {
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
	cmdTreeEntries := `	@{ Context = "/start/"; Options = @('--config', '-c') }`
	programName := "myapp"

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)

	result := buf.String()

	// 检查必要的组件
	expectedComponents := []string{
		"$myapp_commandName = \"myapp\"",
		"$myapp_cmdTree = @(",
		"$myapp_flagParams = @(",
		"Register-ArgumentCompleter",
		"-CommandName ${myapp_commandName}",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(result, component) {
			t.Errorf("PowerShell补全脚本不包含必要组件: %s", component)
		}
	}

	// 检查根命令条目
	expectedRoot := `Context = "/"; Options = @('--help', '-h', 'start', 'stop')`
	if !strings.Contains(result, expectedRoot) {
		t.Errorf("PowerShell补全脚本不包含根命令条目: %s", expectedRoot)
	}

	// 检查标志参数
	expectedFlags := []string{
		`Context = "/"; Parameter = "--verbose"; ParamType = "none"; ValueType = "bool"`,
		`Context = "/"; Parameter = "--mode"; ParamType = "required"; ValueType = "enum"`,
		`Context = "/start/"; Parameter = "--config"; ParamType = "required"; ValueType = "string"`,
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(result, flag) {
			t.Errorf("PowerShell补全脚本不包含标志参数: %s", flag)
		}
	}

	// 检查枚举选项
	expectedEnum := `Options = @('debug', 'release', 'test')`
	if !strings.Contains(result, expectedEnum) {
		t.Errorf("PowerShell补全脚本不包含枚举选项: %s", expectedEnum)
	}
}

// TestEscapePwshString 测试PowerShell字符串转义
func TestEscapePwshString(t *testing.T) {
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
			name:     "包含单引号",
			input:    "don't",
			expected: "don''t",
		},
		{
			name:     "包含反斜杠",
			input:    "path\\to\\file",
			expected: "path\\\\to\\\\file",
		},
		{
			name:     "混合特殊字符",
			input:    "path\\to 'my file'",
			expected: "path\\\\to ''my file''",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有特殊字符",
			input:    "'\\",
			expected: "''\\\\",
		},
		{
			name:     "多个连续单引号",
			input:    "'''",
			expected: "''''''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapePwshString(tt.input)
			if result != tt.expected {
				t.Errorf("escapePwshString(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFormatOptions 测试选项格式化
func TestFormatOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []string
		expected string
	}{
		{
			name:     "普通选项",
			options:  []string{"option1", "option2", "option3"},
			expected: "'option1', 'option2', 'option3'",
		},
		{
			name:     "包含特殊字符的选项",
			options:  []string{"path\\file", "don't", "normal"},
			expected: "'path\\\\file', 'don''t', 'normal'",
		},
		{
			name:     "空选项列表",
			options:  []string{},
			expected: "",
		},
		{
			name:     "包含空字符串的选项",
			options:  []string{"valid", "", "another"},
			expected: "'valid', 'another'", // 空字符串应该被跳过
		},
		{
			name:     "单个选项",
			options:  []string{"single"},
			expected: "'single'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			formatOptions(&buf, tt.options, escapePwshString)

			result := buf.String()
			if result != tt.expected {
				t.Errorf("formatOptions(%v) = %q, 期望 %q", tt.options, result, tt.expected)
			}
		})
	}
}

// TestPwshTemplateConstants 测试PowerShell模板常量
func TestPwshTemplateConstants(t *testing.T) {
	// 测试标志参数条目格式
	expectedFlagParam := `	@{ Context = "{{.Context}}"; Parameter = "{{.Parameter}}"; ParamType = "{{.ParamType}}"; ValueType = "{{.ValueType}}"; Options = @({{.Options}}) }`
	if PwshFlagParamItem != expectedFlagParam {
		t.Errorf("PwshFlagParamItem 格式错误:\n期望: %s\n实际: %s", expectedFlagParam, PwshFlagParamItem)
	}

	// 测试命令树条目格式
	expectedCmdTree := `	@{ Context = "{{.Context}}"; Options = @({{.Options}}) }`
	if PwshCmdTreeItem != expectedCmdTree {
		t.Errorf("PwshCmdTreeItem 格式错误:\n期望: %s\n实际: %s", expectedCmdTree, PwshCmdTreeItem)
	}

	// 测试函数头部模板
	if !strings.Contains(PwshFunctionHeader, "Register-ArgumentCompleter") {
		t.Error("PwshFunctionHeader 不包含Register-ArgumentCompleter")
	}

	if !strings.Contains(PwshFunctionHeader, "${{.SanitizedName}}_commandName") {
		t.Error("PwshFunctionHeader 不包含命令名称变量")
	}

	if !strings.Contains(PwshFunctionHeader, "${{.SanitizedName}}_cmdTree") {
		t.Error("PwshFunctionHeader 不包含命令树变量")
	}

	if !strings.Contains(PwshFunctionHeader, "${{.SanitizedName}}_flagParams") {
		t.Error("PwshFunctionHeader 不包含标志参数变量")
	}
}

// TestPwshCompletionWithComplexScenario 测试复杂场景的PowerShell补全
func TestPwshCompletionWithComplexScenario(t *testing.T) {
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

	cmdTreeEntries := `	@{ Context = "/start/"; Options = @('--port', '-p', '--daemon', '-d') },
	@{ Context = "/stop/"; Options = @('--force', '-f') },
	@{ Context = "/config/"; Options = @('set', 'get', 'list') },
	@{ Context = "/config/set/"; Options = @('--key', '--value') }`

	programName := "complexapp"

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)

	result := buf.String()

	// 验证复杂场景的各个部分
	tests := []struct {
		name     string
		contains string
	}{
		{"程序名称", "$complexapp_commandName = \"complexapp\""},
		{"根命令选项", `Context = "/"; Options = @('--help', '-h', '--verbose', '-v', '--config', '-c', '--mode', '-m', 'start', 'stop', 'config')`},
		{"子命令树", `Context = "/start/"; Options = @('--port', '-p', '--daemon', '-d')`},
		{"深层命令树", `Context = "/config/set/"; Options = @('--key', '--value')`},
		{"根命令枚举", `Options = @('dev', 'prod', 'test')`},
		{"根命令标志", `Context = "/"; Parameter = "--config"; ParamType = "required"; ValueType = "string"`},
		{"子命令标志", `Context = "/start/"; Parameter = "--port"; ParamType = "required"; ValueType = "string"`},
		{"深层标志", `Context = "/config/set/"; Parameter = "--key"; ParamType = "required"; ValueType = "string"`},
		{"注册补全", "Register-ArgumentCompleter -CommandName ${complexapp_commandName}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(result, tt.contains) {
				t.Errorf("复杂PowerShell补全脚本不包含: %s", tt.contains)
			}
		})
	}
}

// BenchmarkGeneratePwshCompletion 基准测试PowerShell补全生成
func BenchmarkGeneratePwshCompletion(b *testing.B) {
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

	cmdTreeEntries := `	@{ Context = "/sub1/"; Options = @('--flag1', '-f1') },
	@{ Context = "/sub2/"; Options = @('--flag2', '-f2') }`
	programName := "benchapp"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		generatePwshCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)
	}
}

// BenchmarkEscapePwshString 基准测试PowerShell字符串转义
func BenchmarkEscapePwshString(b *testing.B) {
	testStrings := []string{
		"normal string",
		"path\\to\\file with 'quotes' and spaces",
		"simple",
		"complex'string\\with\\many'special\\chars",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_ = escapePwshString(s)
		}
	}
}

// BenchmarkFormatOptions 基准测试选项格式化
func BenchmarkFormatOptions(b *testing.B) {
	options := make([]string, 100)
	for i := 0; i < 100; i++ {
		options[i] = "option" + string(rune('a'+i%26))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		formatOptions(&buf, options, escapePwshString)
	}
}
