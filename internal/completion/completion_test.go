// Package completion 自动补全测试
// 本文件包含了自动补全系统的单元测试，测试标志补全、
// 参数值补全等功能，为用户提供便捷的命令行交互体验。
package completion

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestCompletionPerformance 测试补全脚本生成性能
func TestCompletionPerformance(t *testing.T) {
	// ---------- 阶段 1：构建命令树 ----------
	buildStart := time.Now()
	rootCtx := createTestContext("root", "r")

	// 统计节点、flag 数量
	var (
		totalCmds  = 1 // root
		totalFlags int
	)

	// 创建复杂的命令树结构
	for i := 0; i < 10; i++ {
		parentCtx := createTestContext(fmt.Sprintf("sub%d", i), fmt.Sprintf("s%d", i))
		rootCtx.SubCmds = append(rootCtx.SubCmds, parentCtx)
		if rootCtx.SubCmdMap == nil {
			rootCtx.SubCmdMap = make(map[string]*types.CmdContext)
		}
		rootCtx.SubCmdMap[parentCtx.LongName] = parentCtx
		rootCtx.SubCmdMap[parentCtx.ShortName] = parentCtx
		totalCmds++

		// 为每个父命令添加20个标志
		for j := 0; j < 20; j++ {
			var value string
			flag := &flags.StringFlag{}
			flagName := fmt.Sprintf("option%d", j)
			shortName := fmt.Sprintf("o%d", j)
			_ = flag.Init(flagName, shortName, "测试选项", &value)
			meta := &flags.FlagMeta{Flag: flag}
			_ = parentCtx.FlagRegistry.RegisterFlag(meta)
			totalFlags++
		}

		// 为每个父命令创建5个子命令
		for k := 0; k < 5; k++ {
			childCtx := createTestContext(fmt.Sprintf("sub%d-grand%d", i, k), fmt.Sprintf("g%d", k))
			parentCtx.SubCmds = append(parentCtx.SubCmds, childCtx)
			if parentCtx.SubCmdMap == nil {
				parentCtx.SubCmdMap = make(map[string]*types.CmdContext)
			}
			parentCtx.SubCmdMap[childCtx.LongName] = childCtx
			parentCtx.SubCmdMap[childCtx.ShortName] = childCtx
			totalCmds++

			// 为每个子命令添加15个标志
			for l := 0; l < 15; l++ {
				var value int
				flag := &flags.IntFlag{}
				flagName := fmt.Sprintf("param%d", l)
				shortName := fmt.Sprintf("p%d", l)
				_ = flag.Init(flagName, shortName, "测试参数", &value)
				meta := &flags.FlagMeta{Flag: flag}
				_ = childCtx.FlagRegistry.RegisterFlag(meta)
				totalFlags++
			}
		}
	}
	buildDuration := time.Since(buildStart)
	t.Logf("构建命令树耗时: %v, 命令数量=%d, 标志数量=%d", buildDuration, totalCmds, totalFlags)

	// ---------- 阶段 2：内存基线 ----------
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	// ---------- 阶段 3：Bash 补全 ----------
	t.Run("bash", func(t *testing.T) {
		start := time.Now()
		script, err := GenerateShellCompletion(rootCtx, flags.ShellBash)
		if err != nil {
			t.Fatalf("bash 生成失败: %v", err)
		}
		genDuration := time.Since(start)
		t.Logf("bash 生成 %d 字节耗时: %v", len(script), genDuration)

		// 阈值：50 ms 以内
		if genDuration > 50*time.Millisecond {
			t.Errorf("bash 生成耗时过长: %v", genDuration)
		}
	})

	// ---------- 阶段 4：PowerShell 补全 ----------
	t.Run("pwsh", func(t *testing.T) {
		start := time.Now()
		script, err := GenerateShellCompletion(rootCtx, flags.ShellPowershell)
		if err != nil {
			t.Fatalf("pwsh 生成失败: %v", err)
		}
		genDuration := time.Since(start)
		t.Logf("pwsh 生成 %d 字节耗时: %v", len(script), genDuration)

		// 阈值：75 ms 以内
		if genDuration > 75*time.Millisecond {
			t.Errorf("pwsh 生成耗时过长: %v", genDuration)
		}
	})

	// ---------- 阶段 5：内存增量 ----------
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	allocKB := (after.Alloc - before.Alloc) / 1024
	t.Logf("内存增量: %d KB", allocKB)
}

// TestCompletionBash 测试Bash自动补全生成
func TestCompletionBash(t *testing.T) {
	// 创建根命令上下文
	rootCtx := createTestContext("root", "r")

	// 创建子命令
	cmd1Ctx := createTestContext("cmd1", "c1")
	cmd2Ctx := createTestContext("cmd2", "c2")

	// 为cmd1添加字符串标志
	var strValue string
	strFlag := &flags.StringFlag{}
	_ = strFlag.Init("str", "s", "test string", &strValue)
	strMeta := &flags.FlagMeta{Flag: strFlag}
	_ = cmd1Ctx.FlagRegistry.RegisterFlag(strMeta)

	// 为cmd2添加整数标志
	var intValue int
	intFlag := &flags.IntFlag{}
	_ = intFlag.Init("int", "i", "test int", &intValue)
	intMeta := &flags.FlagMeta{Flag: intFlag}
	_ = cmd2Ctx.FlagRegistry.RegisterFlag(intMeta)

	// 构建命令树
	rootCtx.SubCmds = []*types.CmdContext{cmd1Ctx}
	rootCtx.SubCmdMap = map[string]*types.CmdContext{
		"cmd1": cmd1Ctx,
		"c1":   cmd1Ctx,
	}

	cmd1Ctx.SubCmds = []*types.CmdContext{cmd2Ctx}
	cmd1Ctx.SubCmdMap = map[string]*types.CmdContext{
		"cmd2": cmd2Ctx,
		"c2":   cmd2Ctx,
	}

	// 生成Bash补全脚本
	script, err := GenerateShellCompletion(rootCtx, flags.ShellBash)
	if err != nil {
		t.Fatalf("生成Bash补全脚本失败: %v", err)
	}

	// 验证生成的脚本包含预期内容
	expectedContents := []string{
		"#!/usr/bin/env bash",
		"declare -A completion.test.exe_cmd_tree",
		"declare -A completion.test.exe_flag_params",
		"completion.test.exe_cmd_tree[/]",
		"completion.test.exe_cmd_tree[/cmd1/]",
		"completion.test.exe_cmd_tree[/c1/]",
		"complete -F _completion.test.exe completion.test.exe",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(script, expected) {
			t.Errorf("PowerShell补全脚本不包含预期内容: %s", expected)
		}
	}

	fmt.Println(script)
}

// TestCompletionPwsh 测试PowerShell自动补全生成
func TestCompletionPwsh(t *testing.T) {
	// 创建根命令上下文
	rootCtx := createTestContext("root", "r")

	// 创建子命令
	cmd1Ctx := createTestContext("cmd1", "c1")
	cmd2Ctx := createTestContext("cmd2", "c2")

	// 为cmd1添加字符串标志
	var strValue string
	strFlag := &flags.StringFlag{}
	_ = strFlag.Init("str", "s", "test string", &strValue)
	strMeta := &flags.FlagMeta{Flag: strFlag}
	_ = cmd1Ctx.FlagRegistry.RegisterFlag(strMeta)

	// 为cmd2添加整数标志
	var intValue int
	intFlag := &flags.IntFlag{}
	_ = intFlag.Init("int", "i", "test int", &intValue)
	intMeta := &flags.FlagMeta{Flag: intFlag}
	_ = cmd2Ctx.FlagRegistry.RegisterFlag(intMeta)

	// 构建命令树
	rootCtx.SubCmds = []*types.CmdContext{cmd1Ctx}
	rootCtx.SubCmdMap = map[string]*types.CmdContext{
		"cmd1": cmd1Ctx,
		"c1":   cmd1Ctx,
	}

	cmd1Ctx.SubCmds = []*types.CmdContext{cmd2Ctx}
	cmd1Ctx.SubCmdMap = map[string]*types.CmdContext{
		"cmd2": cmd2Ctx,
		"c2":   cmd2Ctx,
	}

	// 生成PowerShell补全脚本
	script, err := GenerateShellCompletion(rootCtx, flags.ShellPowershell)
	if err != nil {
		t.Fatalf("生成PowerShell补全脚本失败: %v", err)
	}

	// 验证生成的脚本包含预期内容
	expectedContents := []string{
		"Register-ArgumentCompleter",
		"$completion.test_cmdTree = @(",
		"$completion.test_flagParams = @(",
		"Context = \"/cmd1/\"",
		"Context = \"/c1/\"",
		"Parameter = \"--string\"",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(script, expected) {
			t.Errorf("PowerShell补全脚本不包含预期内容: %s", expected)
		}
	}

	fmt.Println(script)
}

// createTestContext 创建测试用的命令上下文
func createTestContext(longName, shortName string) *types.CmdContext {
	ctx := types.NewCmdContext(longName, shortName, flag.ContinueOnError)

	// 添加一些测试标志
	var stringValue string
	stringFlag := &flags.StringFlag{}
	_ = stringFlag.Init("string", "s", "字符串标志", &stringValue)
	stringMeta := &flags.FlagMeta{Flag: stringFlag}
	_ = ctx.FlagRegistry.RegisterFlag(stringMeta)

	var boolValue bool
	boolFlag := &flags.BoolFlag{}
	_ = boolFlag.Init("verbose", "v", "详细输出", &boolValue)
	boolMeta := &flags.FlagMeta{Flag: boolFlag}
	_ = ctx.FlagRegistry.RegisterFlag(boolMeta)

	enumFlag := &flags.EnumFlag{}
	_ = enumFlag.Init("mode", "m", "debug", "运行模式", []string{"debug", "release", "test"})
	enumMeta := &flags.FlagMeta{Flag: enumFlag}
	_ = ctx.FlagRegistry.RegisterFlag(enumMeta)

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
				"declare -A completion.test.exe_cmd_tree",
				"declare -A completion.test.exe_flag_params",
				"complete -F _completion.test.exe completion.test.exe",
			},
		},
		{
			name:      "生成PowerShell补全脚本",
			shellType: "powershell",
			wantErr:   false,
			contains: []string{
				"Register-ArgumentCompleter",
				"$completion.test_cmdTree = @(",
				"$completion.test_flagParams = @(",
			},
		},
		{
			name:      "生成Pwsh补全脚本",
			shellType: "pwsh",
			wantErr:   false,
			contains: []string{
				"Register-ArgumentCompleter",
				"$completion.test_cmdTree = @(",
				"$completion.test_flagParams = @(",
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
				"_cmd_tree[/sub1/]",
				"_cmd_tree[/s1/]",
				"_cmd_tree[/sub2/]",
				"_cmd_tree[/s2/]",
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

// TestCompletionNone 测试无补全模式
func TestCompletionNone(t *testing.T) {
	// 测试场景1: ShellNone模式下不应生成任何补全脚本
	t.Run("no_completion_script_generated", func(t *testing.T) {
		// 创建命令上下文
		ctx := createTestContext("test", "t")

		// 生成补全脚本（使用不支持的shell类型）
		script, err := GenerateShellCompletion(ctx, "none")
		if err != nil {
			t.Fatalf("生成补全脚本失败: %v", err)
		}

		// 验证没有生成补全脚本内容
		if script != "" {
			t.Errorf("不支持的shell类型不应生成补全脚本，实际输出: %q", script)
		}
	})

	// 测试场景2: 验证空命令上下文的处理
	t.Run("empty_context_handling", func(t *testing.T) {
		// 测试nil上下文
		_, err := GenerateShellCompletion(nil, flags.ShellBash)
		if err == nil {
			t.Error("nil上下文应该返回错误")
		}

		// 测试空的标志注册表
		ctx := createTestContext("test", "t")
		ctx.FlagRegistry = nil
		_, err = GenerateShellCompletion(ctx, flags.ShellBash)
		if err == nil {
			t.Error("nil标志注册表应该返回错误")
		}
	})
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
		_ = flag.Init(flagName, shortName, "测试标志", &value)
		meta := &flags.FlagMeta{Flag: flag}
		_ = ctx.FlagRegistry.RegisterFlag(meta)
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
		_ = flag.Init(flagName, shortName, "测试标志", &value)
		meta := &flags.FlagMeta{Flag: flag}
		_ = ctx.FlagRegistry.RegisterFlag(meta)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collectFlagParameters(ctx)
	}
}
