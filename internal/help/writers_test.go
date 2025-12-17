// Package help 帮助信息输出器测试
// 本文件包含了帮助信息输出器的单元测试，测试多种输出格式和
// 样式的帮助信息展示功能的正确性和稳定性。
package help

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestWriteModuleHelps 测试写入模块帮助信息
func TestWriteModuleHelps(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *types.CmdContext
		moduleHelps string
		expected    string
	}{
		{
			name:        "空上下文",
			ctx:         nil,
			moduleHelps: "",
			expected:    "",
		},
		{
			name:        "空模块帮助",
			ctx:         createTestContext("test", "t"),
			moduleHelps: "",
			expected:    "",
		},
		{
			name:        "有模块帮助",
			ctx:         createTestContext("test", "t"),
			moduleHelps: "模块帮助信息",
			expected:    "\n模块帮助信息\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if tt.ctx != nil {
				tt.ctx.Config.ModuleHelps = tt.moduleHelps
			}
			writeModuleHelps(tt.ctx, &buf)
			result := buf.String()
			if result != tt.expected {
				t.Errorf("期望 '%s'，但得到 '%s'", tt.expected, result)
			}
		})
	}
}

// TestWriteLogoText 测试写入Logo信息
func TestWriteLogoText(t *testing.T) {
	tests := []struct {
		name     string
		ctx      *types.CmdContext
		logoText string
		expected string
	}{
		{
			name:     "空上下文",
			ctx:      nil,
			logoText: "",
			expected: "",
		},
		{
			name:     "空Logo",
			ctx:      createTestContext("test", "t"),
			logoText: "",
			expected: "",
		},
		{
			name:     "有Logo",
			ctx:      createTestContext("test", "t"),
			logoText: "测试Logo",
			expected: "测试Logo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if tt.ctx != nil {
				tt.ctx.Config.LogoText = tt.logoText
			}
			writeLogoText(tt.ctx, &buf)
			result := buf.String()
			if result != tt.expected {
				t.Errorf("期望 '%s'，但得到 '%s'", tt.expected, result)
			}
		})
	}
}

// TestWriteCommandHeader 测试写入命令头部信息
func TestWriteCommandHeader(t *testing.T) {
	tests := []struct {
		name        string
		longName    string
		shortName   string
		description string
		template    HelpTemplate
		expected    string
	}{
		{
			name:      "长短名称都有",
			longName:  "testcmd",
			shortName: "tc",
			template:  ChineseTemplate,
			expected:  "名称: testcmd, tc\n\n",
		},
		{
			name:      "只有长名称",
			longName:  "testcmd",
			shortName: "",
			template:  ChineseTemplate,
			expected:  "名称: testcmd\n\n",
		},
		{
			name:      "只有短名称",
			longName:  "",
			shortName: "tc",
			template:  ChineseTemplate,
			expected:  "名称: tc\n\n",
		},
		{
			name:        "有描述信息",
			longName:    "testcmd",
			shortName:   "tc",
			description: "测试命令",
			template:    ChineseTemplate,
			expected:    "名称: testcmd, tc\n\n描述: 测试命令\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := createTestContext(tt.longName, tt.shortName)
			ctx.Config.Desc = tt.description
			writeCommandHeader(ctx, tt.template, &buf)
			result := buf.String()
			if result != tt.expected {
				t.Errorf("期望 '%s'，但得到 '%s'", tt.expected, result)
			}
		})
	}
}

// TestWriteUsageLine 测试写入用法说明
func TestWriteUsageLine(t *testing.T) {
	tests := []struct {
		name         string
		ctx          *types.CmdContext
		usageSyntax  string
		hasSubCmds   bool
		template     HelpTemplate
		expectedPart string
	}{
		{
			name:         "自定义用法语法",
			ctx:          createTestContext("test", "t"),
			usageSyntax:  "test [options] <file>",
			template:     ChineseTemplate,
			expectedPart: "用法: test [options] <file>",
		},
		{
			name:         "主命令有子命令",
			ctx:          createTestContext("main", "m"),
			hasSubCmds:   true,
			template:     ChineseTemplate,
			expectedPart: "用法: main [全局选项] [子命令] [选项]",
		},
		{
			name:         "主命令无子命令",
			ctx:          createTestContext("main", "m"),
			hasSubCmds:   false,
			template:     ChineseTemplate,
			expectedPart: "用法: main [选项]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if tt.usageSyntax != "" {
				tt.ctx.Config.UsageSyntax = tt.usageSyntax
			}
			if tt.hasSubCmds {
				subCtx := createTestContext("sub", "s")
				tt.ctx.SubCmds = append(tt.ctx.SubCmds, subCtx)
			}
			writeUsageLine(tt.ctx, tt.template, &buf)
			result := buf.String()
			if !strings.Contains(result, tt.expectedPart) {
				t.Errorf("结果应包含 '%s'，但得到 '%s'", tt.expectedPart, result)
			}
		})
	}
}

// TestWriteOptions 测试写入选项信息
func TestWriteOptions(t *testing.T) {
	ctx := createTestContext("test", "t")

	// 添加测试标志
	addTestFlag(ctx, "verbose", "v", "详细输出", "bool", false)
	addTestFlag(ctx, "output", "", "输出文件", "string", "default.txt")
	addTestFlag(ctx, "", "q", "安静模式", "bool", false)

	var buf bytes.Buffer
	writeOptions(ctx, ChineseTemplate, &buf)
	result := buf.String()

	// 验证包含选项标题
	if !strings.Contains(result, "选项:") {
		t.Error("结果应包含选项标题")
	}

	// 验证包含各种标志格式
	if !strings.Contains(result, " -v, --verbose") {
		t.Error("结果应包含长短标志格式")
	}
	if !strings.Contains(result, "--output") {
		t.Error("结果应包含仅长标志格式")
	}
	if !strings.Contains(result, "-q") {
		t.Error("结果应包含仅短标志格式")
	}
}

// TestWriteSubCmds 测试写入子命令信息
func TestWriteSubCmds(t *testing.T) {
	ctx := createTestContext("main", "m")

	// 添加子命令
	subCtx1 := createTestContext("subcmd1", "s1")
	subCtx1.Config.Desc = "第一个子命令"
	subCtx2 := createTestContext("subcmd2", "")
	subCtx2.Config.Desc = "第二个子命令"

	ctx.SubCmds = append(ctx.SubCmds, subCtx1, subCtx2)

	var buf bytes.Buffer
	writeSubCmds(ctx, ChineseTemplate, &buf)
	result := buf.String()

	// 验证包含子命令标题
	if !strings.Contains(result, "子命令:") {
		t.Error("结果应包含子命令标题")
	}

	// 验证包含子命令信息
	if !strings.Contains(result, "subcmd1, s1") {
		t.Error("结果应包含带短名称的子命令")
	}
	if !strings.Contains(result, "subcmd2") {
		t.Error("结果应包含仅长名称的子命令")
	}
	if !strings.Contains(result, "第一个子命令") {
		t.Error("结果应包含子命令描述")
	}
}

// TestWriteExamples 测试写入示例信息
func TestWriteExamples(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.Examples = []types.ExampleInfo{
		{Desc: "基本用法", Usage: "test --help"},
		{Desc: "详细模式", Usage: "test --verbose"},
	}

	var buf bytes.Buffer
	writeExamples(ctx, ChineseTemplate, &buf)
	result := buf.String()

	// 验证包含示例标题
	if !strings.Contains(result, "示例:") {
		t.Error("结果应包含示例标题")
	}

	// 验证包含示例内容
	if !strings.Contains(result, "1、基本用法") {
		t.Error("结果应包含第一个示例")
	}
	if !strings.Contains(result, "test --help") {
		t.Error("结果应包含示例用法")
	}
}

// TestWriteNotes 测试写入注意事项
func TestWriteNotes(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.Notes = []string{
		"这是第一个注意事项",
		"这是第二个注意事项",
	}

	var buf bytes.Buffer
	writeNotes(ctx, ChineseTemplate, &buf)
	result := buf.String()

	// 验证包含注意事项标题
	if !strings.Contains(result, "注意事项:") {
		t.Error("结果应包含注意事项标题")
	}

	// 验证包含注意事项内容
	if !strings.Contains(result, "1、这是第一个注意事项") {
		t.Error("结果应包含第一个注意事项")
	}
	if !strings.Contains(result, "2、这是第二个注意事项") {
		t.Error("结果应包含第二个注意事项")
	}
}

// TestCollectFlags 测试收集标志信息
func TestCollectFlags(t *testing.T) {
	tests := []struct {
		name     string
		ctx      *types.CmdContext
		expected int
	}{
		{
			name:     "空上下文",
			ctx:      nil,
			expected: 0,
		},
		{
			name:     "无标志",
			ctx:      createTestContext("test", "t"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := collectFlags(tt.ctx)
			if len(flags) != tt.expected {
				t.Errorf("期望 %d 个标志，但得到 %d 个", tt.expected, len(flags))
			}
		})
	}

	// 测试有标志的情况
	ctx := createTestContext("test", "t")
	addTestFlag(ctx, "verbose", "v", "详细输出", "bool", false)
	addTestFlag(ctx, "output", "", "输出文件", "string", "default.txt")

	flags := collectFlags(ctx)
	if len(flags) != 2 {
		t.Errorf("期望 2 个标志，但得到 %d 个", len(flags))
	}

	// 验证标志信息
	found := false
	for _, flag := range flags {
		if flag.longFlag == "verbose" && flag.shortFlag == "v" {
			found = true
			if flag.usage != "详细输出" {
				t.Error("标志用法信息不正确")
			}
			break
		}
	}
	if !found {
		t.Error("未找到预期的verbose标志")
	}
}

// TestCalculateMaxWidth 测试计算最大宽度
func TestCalculateMaxWidth(t *testing.T) {
	tests := []struct {
		name     string
		flags    []flagInfo
		expected int
	}{
		{
			name:     "空标志列表",
			flags:    []flagInfo{},
			expected: 0,
		},
		{
			name: "单个长短标志",
			flags: []flagInfo{
				{longFlag: "verbose", shortFlag: "v", typeStr: "<bool>"},
			},
			expected: 22, // "--verbose, -v <bool>" 的实际长度
		},
		{
			name: "仅长标志",
			flags: []flagInfo{
				{longFlag: "output", shortFlag: "", typeStr: "<string>"},
			},
			expected: 19, // "--output <string>" 的实际长度
		},
		{
			name: "仅短标志",
			flags: []flagInfo{
				{longFlag: "", shortFlag: "q", typeStr: "<bool>"},
			},
			expected: 11, // "-q <bool>" 的实际长度
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMaxWidth(tt.flags)
			if result != tt.expected {
				t.Errorf("期望宽度 %d，但得到 %d", tt.expected, result)
			}
		})
	}
}

// TestGetFullCommandPath 测试获取完整命令路径
func TestGetFullCommandPath(t *testing.T) {
	tests := []struct {
		name     string
		ctx      *types.CmdContext
		expected string
	}{
		{
			name:     "空上下文",
			ctx:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFullCommandPath(tt.ctx)
			if result != tt.expected {
				t.Errorf("期望 '%s'，但得到 '%s'", tt.expected, result)
			}
		})
	}

	// 测试根命令
	rootCtx := createTestContext("root", "r")
	result := getFullCommandPath(rootCtx)
	expected := "root"
	if result != expected {
		t.Errorf("根命令路径期望 '%s'，但得到 '%s'", expected, result)
	}

	// 测试子命令
	subCtx := createTestContext("sub", "s")
	subCtx.Parent = rootCtx
	result = getFullCommandPath(subCtx)
	expected = "root sub"
	if result != expected {
		t.Errorf("子命令路径期望 '%s'，但得到 '%s'", expected, result)
	}
}

// TestWriteOptions_NilBuffer 测试空缓冲区
func TestWriteOptions_NilBuffer(t *testing.T) {
	ctx := createTestContext("test", "t")
	writeOptions(ctx, ChineseTemplate, nil) // 不应该崩溃
}

// TestWriteSubCmds_NilBuffer 测试空缓冲区
func TestWriteSubCmds_NilBuffer(t *testing.T) {
	ctx := createTestContext("test", "t")
	writeSubCmds(ctx, ChineseTemplate, nil) // 不应该崩溃
}

// TestWriteExamples_EmptyExamples 测试空示例列表
func TestWriteExamples_EmptyExamples(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.Examples = []types.ExampleInfo{}

	var buf bytes.Buffer
	writeExamples(ctx, ChineseTemplate, &buf)
	result := buf.String()

	if result != "" {
		t.Error("空示例列表应该不输出任何内容")
	}
}

// TestWriteNotes_EmptyNotes 测试空注意事项列表
func TestWriteNotes_EmptyNotes(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.Notes = []string{}

	var buf bytes.Buffer
	writeNotes(ctx, ChineseTemplate, &buf)
	result := buf.String()

	if result != "" {
		t.Error("空注意事项列表应该不输出任何内容")
	}
}

// TestWriteOptions_NoFlags 测试无标志情况
func TestWriteOptions_NoFlags(t *testing.T) {
	ctx := createTestContext("test", "t")

	var buf bytes.Buffer
	writeOptions(ctx, ChineseTemplate, &buf)
	result := buf.String()

	if result != "" {
		t.Error("无标志时应该不输出任何内容")
	}
}

// TestWriteSubCmds_NoSubCmds 测试无子命令情况
func TestWriteSubCmds_NoSubCmds(t *testing.T) {
	ctx := createTestContext("test", "t")

	var buf bytes.Buffer
	writeSubCmds(ctx, ChineseTemplate, &buf)
	result := buf.String()

	if result != "" {
		t.Error("无子命令时应该不输出任何内容")
	}
}

// TestCollectFlags_WithDurationFlag 测试Duration类型标志的收集
func TestCollectFlags_WithDurationFlag(t *testing.T) {
	ctx := createTestContext("test", "t")

	// 创建Duration标志
	durationFlag := &flags.DurationFlag{}
	currentDuration := new(time.Duration)
	*currentDuration = 5 * time.Second
	if err := durationFlag.Init("timeout", "t", "超时时间", currentDuration); err != nil {
		t.Fatalf("初始化Duration标志失败: %v", err)
	}

	// 注册标志
	if err := ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: durationFlag}); err != nil {
		t.Fatalf("注册标志失败: %v", err)
	}

	flags := collectFlags(ctx)
	if len(flags) != 1 {
		t.Errorf("期望 1 个标志，但得到 %d 个", len(flags))
	}

	if flags[0].defValue != "5s" {
		t.Errorf("Duration标志默认值期望 '5s'，但得到 '%s'", flags[0].defValue)
	}
}
