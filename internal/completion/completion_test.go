package completion

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// createTestContext 创建测试用的命令上下文
func createTestContext(longName, shortName string) *types.CmdContext {
	ctx := types.NewCmdContext(longName, shortName, flag.ContinueOnError)

	// 添加一些测试标志
	var stringValue string
	stringFlag := &flags.StringFlag{}
	stringFlag.Init("string", "s", "字符串标志", &stringValue)
	stringMeta := &flags.FlagMeta{Flag: stringFlag}
	ctx.FlagRegistry.RegisterFlag(stringMeta)

	var boolValue bool
	boolFlag := &flags.BoolFlag{}
	boolFlag.Init("verbose", "v", "详细输出", &boolValue)
	boolMeta := &flags.FlagMeta{Flag: boolFlag}
	ctx.FlagRegistry.RegisterFlag(boolMeta)

	enumFlag := &flags.EnumFlag{}
	enumFlag.Init("mode", "m", "debug", "运行模式", []string{"debug", "release", "test"})
	enumMeta := &flags.FlagMeta{Flag: enumFlag}
	ctx.FlagRegistry.RegisterFlag(enumMeta)

	return ctx
}

// TestGenerateShellCompletion 测试生成shell补全脚本
func TestGenerateShellCompletion(t *testing.T) {
	tests := []struct {
		name      string
		shellType string
		wantErr   bool
		contains  []string // 期望包含的字符串
	}{
		{
			name:      "生成Bash补全脚本",
			shellType: "bash",
			wantErr:   false,
			contains: []string{
				"declare -A cmd_tree",
				"declare -A flag_params",
				"complete -F",
			},
		},
		{
			name:      "生成PowerShell补全脚本",
			shellType: "powershell",
			wantErr:   false,
			contains: []string{
				"Register-ArgumentCompleter",
				"_cmdTree = @(",
				"_flagParams = @(",
			},
		},
		{
			name:      "生成Pwsh补全脚本",
			shellType: "pwsh",
			wantErr:   false,
			contains: []string{
				"Register-ArgumentCompleter",
				"_cmdTree = @(",
				"_flagParams = @(",
			},
		},
		{
			name:      "不支持的shell类型",
			shellType: "unsupported",
			wantErr:   false, // 函数不会报错，但生成的脚本为空
			contains:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("myapp", "")

			result, err := GenerateShellCompletion(ctx, tt.shellType)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateShellCompletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 检查是否包含期望的字符串
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("GenerateShellCompletion() 结果不包含期望的字符串: %s", expected)
				}
			}
		})
	}
}

// TestGenerateShellCompletionWithSubCommands 测试带子命令的补全脚本生成
func TestGenerateShellCompletionWithSubCommands(t *testing.T) {
	// 创建根命令
	rootCtx := createTestContext("myapp", "")

	// 创建子命令
	subCtx1 := createTestContext("start", "s")
	subCtx2 := createTestContext("stop", "st")

	// 添加子命令到根命令
	rootCtx.SubCmds = []*types.CmdContext{subCtx1, subCtx2}
	rootCtx.SubCmdMap = map[string]*types.CmdContext{
		"start": subCtx1,
		"s":     subCtx1,
		"stop":  subCtx2,
		"st":    subCtx2,
	}

	result, err := GenerateShellCompletion(rootCtx, flags.ShellBash)
	if err != nil {
		t.Fatalf("GenerateShellCompletion() error = %v", err)
	}

	// 检查子命令是否包含在结果中
	expectedSubCmds := []string{"/start/", "/s/", "/stop/", "/st/"}
	for _, subCmd := range expectedSubCmds {
		if !strings.Contains(result, subCmd) {
			t.Errorf("补全脚本不包含子命令: %s", subCmd)
		}
	}
}

// TestValidateCompletionGeneration 测试补全生成验证
func TestValidateCompletionGeneration(t *testing.T) {
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
				ctx := createTestContext("sub", "")
				parent := createTestContext("parent", "")
				ctx.Parent = parent
				return ctx
			}(),
			wantErr: true,
			errMsg:  "not a root command",
		},
		{
			name: "标志注册表为nil",
			ctx: func() *types.CmdContext {
				ctx := createTestContext("test", "")
				ctx.FlagRegistry = nil
				return ctx
			}(),
			wantErr: true,
			errMsg:  "flag registry is nil",
		},
		{
			name:    "有效的根命令",
			ctx:     createTestContext("valid", ""),
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

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateCompletionGeneration() error = %v, 期望包含 %v", err, tt.errMsg)
			}
		})
	}
}

// TestCollectCompletionOptions 测试收集补全选项
func TestCollectCompletionOptions(t *testing.T) {
	tests := []struct {
		name     string
		ctx      *types.CmdContext
		expected []string
	}{
		{
			name:     "nil上下文",
			ctx:      nil,
			expected: nil,
		},
		{
			name: "标志注册表为nil",
			ctx: func() *types.CmdContext {
				ctx := createTestContext("test", "")
				ctx.FlagRegistry = nil
				return ctx
			}(),
			expected: nil,
		},
		{
			name:     "正常命令",
			ctx:      createTestContext("test", ""),
			expected: []string{"--string", "-s", "--verbose", "-v", "--mode", "-m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectCompletionOptions(tt.ctx)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("collectCompletionOptions() = %v, 期望 nil", result)
				}
				return
			}

			// 检查所有期望的选项是否都存在
			for _, expected := range tt.expected {
				found := false
				for _, actual := range result {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("collectCompletionOptions() 缺少期望的选项: %s", expected)
				}
			}
		})
	}
}

// TestCollectFlagParameters 测试收集标志参数
func TestCollectFlagParameters(t *testing.T) {
	ctx := createTestContext("test", "")

	params := collectFlagParameters(ctx)

	// 验证参数数量（每个标志有长短名，所以应该有6个参数）
	expectedCount := 6 // string(2) + verbose(2) + mode(2)
	if len(params) != expectedCount {
		t.Errorf("collectFlagParameters() 返回 %d 个参数, 期望 %d", len(params), expectedCount)
	}

	// 验证特定参数
	foundStringFlag := false
	foundEnumFlag := false
	foundBoolFlag := false

	for _, param := range params {
		switch param.Name {
		case "--string":
			foundStringFlag = true
			if param.Type != "required" || param.ValueType != "string" {
				t.Errorf("字符串标志参数类型错误: Type=%s, ValueType=%s", param.Type, param.ValueType)
			}
		case "--mode":
			foundEnumFlag = true
			if param.Type != "required" || param.ValueType != "enum" {
				t.Errorf("枚举标志参数类型错误: Type=%s, ValueType=%s", param.Type, param.ValueType)
			}
			expectedOptions := []string{"debug", "release", "test"}
			if len(param.EnumOptions) != len(expectedOptions) {
				t.Errorf("枚举选项数量错误: got %d, want %d", len(param.EnumOptions), len(expectedOptions))
			}
		case "--verbose":
			foundBoolFlag = true
			if param.Type != "none" || param.ValueType != "bool" {
				t.Errorf("布尔标志参数类型错误: Type=%s, ValueType=%s", param.Type, param.ValueType)
			}
		}
	}

	if !foundStringFlag {
		t.Error("未找到字符串标志参数")
	}
	if !foundEnumFlag {
		t.Error("未找到枚举标志参数")
	}
	if !foundBoolFlag {
		t.Error("未找到布尔标志参数")
	}
}

// TestTraverseCommandTree 测试遍历命令树
func TestTraverseCommandTree(t *testing.T) {
	// 创建子命令
	subCtx1 := createTestContext("sub1", "s1")
	subCtx2 := createTestContext("sub2", "s2")

	// 为子命令添加子子命令
	subSubCtx := createTestContext("subsub", "ss")
	subCtx1.SubCmds = []*types.CmdContext{subSubCtx}

	cmdContexts := []*types.CmdContext{subCtx1, subCtx2}

	tests := []struct {
		name      string
		shellType string
		contains  []string
	}{
		{
			name:      "Bash命令树遍历",
			shellType: flags.ShellBash,
			contains: []string{
				"cmd_tree[/sub1/]",
				"cmd_tree[/s1/]",
				"cmd_tree[/sub2/]",
				"cmd_tree[/s2/]",
			},
		},
		{
			name:      "PowerShell命令树遍历",
			shellType: flags.ShellPowershell,
			contains: []string{
				"Context = \"/sub1/\"",
				"Context = \"/s1/\"",
				"Context = \"/sub2/\"",
				"Context = \"/s2/\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			traverseCommandTree(&buf, "/", cmdContexts, tt.shellType)

			result := buf.String()
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("命令树遍历结果不包含: %s", expected)
				}
			}
		})
	}
}

// TestCompletionNotes 测试补全注意事项
func TestCompletionNotes(t *testing.T) {
	// 测试中文注意事项
	if len(CompletionNotesCN) == 0 {
		t.Error("中文补全注意事项为空")
	}

	// 测试英文注意事项
	if len(CompletionNotesEN) == 0 {
		t.Error("英文补全注意事项为空")
	}

	// 验证注意事项内容
	expectedCNKeywords := []string{"Windows", "PowerShell", "Linux", "bash"}
	for _, keyword := range expectedCNKeywords {
		found := false
		for _, note := range CompletionNotesCN {
			if strings.Contains(note, keyword) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("中文注意事项中未找到关键词: %s", keyword)
		}
	}
}

// TestCompletionExamples 测试补全示例
func TestCompletionExamples(t *testing.T) {
	// 测试中文示例
	if len(CompletionExamplesCN) == 0 {
		t.Error("中文补全示例为空")
	}

	// 测试英文示例
	if len(CompletionExamplesEN) == 0 {
		t.Error("英文补全示例为空")
	}

	// 验证示例格式
	for i, example := range CompletionExamplesCN {
		if example.Description == "" {
			t.Errorf("中文示例 %d 描述为空", i)
		}
		if example.Usage == "" {
			t.Errorf("中文示例 %d 用法为空", i)
		}
		if !strings.Contains(example.Usage, "%s") {
			t.Errorf("中文示例 %d 用法不包含占位符 %%s", i)
		}
	}
}

// BenchmarkGenerateShellCompletion 基准测试补全脚本生成
func BenchmarkGenerateShellCompletion(b *testing.B) {
	ctx := createTestContext("benchmark", "")

	// 添加更多标志以模拟真实场景
	for i := 0; i < 10; i++ {
		var value string
		flag := &flags.StringFlag{}
		flagName := strings.Repeat("flag", i+1)
		shortName := string(rune('a' + i))
		flag.Init(flagName, shortName, "测试标志", &value)
		meta := &flags.FlagMeta{Flag: flag}
		ctx.FlagRegistry.RegisterFlag(meta)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateShellCompletion(ctx, flags.ShellBash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCollectFlagParameters 基准测试收集标志参数
func BenchmarkCollectFlagParameters(b *testing.B) {
	ctx := createTestContext("benchmark", "")

	// 添加大量标志
	for i := 0; i < 100; i++ {
		var value string
		flag := &flags.StringFlag{}
		flagName := strings.Repeat("flag", i+1)
		shortName := string(rune('a' + i%26))
		flag.Init(flagName, shortName, "测试标志", &value)
		meta := &flags.FlagMeta{Flag: flag}
		ctx.FlagRegistry.RegisterFlag(meta)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collectFlagParameters(ctx)
	}
}
