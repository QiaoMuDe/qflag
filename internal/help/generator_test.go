// Package help 帮助信息生成器测试
// 本文件包含了帮助信息生成器的单元测试，测试标志信息、
// 子命令信息等帮助内容的格式化和输出功能的正确性。
package help

import (
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

// TestGenerateHelp_NilContext 测试空上下文情况
func TestGenerateHelp_NilContext(t *testing.T) {
	result := GenerateHelp(nil)
	if result != "" {
		t.Errorf("期望空字符串，但得到: %s", result)
	}
}

// TestGenerateHelp_CustomHelp 测试自定义帮助信息
func TestGenerateHelp_CustomHelp(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.Help = "自定义帮助信息"

	result := GenerateHelp(ctx)
	if result != "自定义帮助信息" {
		t.Errorf("期望 '自定义帮助信息'，但得到: %s", result)
	}
}

// TestGenerateHelp_EnglishTemplate 测试英文模板
func TestGenerateHelp_EnglishTemplate(t *testing.T) {
	ctx := createTestContext("testcmd", "tc")
	ctx.Config.UseChinese = false
	ctx.Config.Desc = "Test command description"

	result := GenerateHelp(ctx)

	// 验证包含英文模板的关键词
	if !strings.Contains(result, "Name: testcmd, tc") {
		t.Error("结果应包含英文格式的命令名称")
	}
	if !strings.Contains(result, "Desc: Test command description") {
		t.Error("结果应包含英文格式的描述")
	}
	if !strings.Contains(result, "Usage: ") {
		t.Error("结果应包含英文格式的用法说明")
	}
}

// TestGenerateHelp_ChineseTemplate 测试中文模板
func TestGenerateHelp_ChineseTemplate(t *testing.T) {
	ctx := createTestContext("测试命令", "测")
	ctx.Config.UseChinese = true
	ctx.Config.Desc = "测试命令描述"

	result := GenerateHelp(ctx)

	// 验证包含中文模板的关键词
	if !strings.Contains(result, "名称: 测试命令, 测") {
		t.Error("结果应包含中文格式的命令名称")
	}
	if !strings.Contains(result, "描述: 测试命令描述") {
		t.Error("结果应包含中文格式的描述")
	}
	if !strings.Contains(result, "用法: ") {
		t.Error("结果应包含中文格式的用法说明")
	}
}

// TestGenerateHelp_WithLogo 测试包含Logo的情况
func TestGenerateHelp_WithLogo(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.LogoText = "测试Logo文本"

	result := GenerateHelp(ctx)
	if !strings.Contains(result, "测试Logo文本") {
		t.Error("结果应包含Logo文本")
	}
}

// TestGenerateHelp_WithModuleHelps 测试包含模块帮助的情况
func TestGenerateHelp_WithModuleHelps(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.ModuleHelps = "模块帮助信息"

	result := GenerateHelp(ctx)
	if !strings.Contains(result, "模块帮助信息") {
		t.Error("结果应包含模块帮助信息")
	}
}

// TestGenerateHelp_WithFlags 测试包含标志的情况
func TestGenerateHelp_WithFlags(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.UseChinese = true

	// 添加测试标志
	addTestFlag(ctx, "verbose", "v", "详细输出", "bool", false)

	result := GenerateHelp(ctx)
	if !strings.Contains(result, "选项:") {
		t.Error("结果应包含选项标题")
	}
	if !strings.Contains(result, " -v, --verbose") {
		t.Error("结果应包含标志信息")
	}
}

// TestGenerateHelp_WithSubCommands 测试包含子命令的情况
func TestGenerateHelp_WithSubCommands(t *testing.T) {
	ctx := createTestContext("main", "m")
	ctx.Config.UseChinese = true

	// 添加子命令
	subCtx := createTestContext("subcmd", "s")
	subCtx.Config.Desc = "子命令描述"
	ctx.SubCmds = append(ctx.SubCmds, subCtx)

	result := GenerateHelp(ctx)
	if !strings.Contains(result, "子命令:") {
		t.Error("结果应包含子命令标题")
	}
	if !strings.Contains(result, "subcmd, s") {
		t.Error("结果应包含子命令信息")
	}
}

// TestGenerateHelp_WithExamples 测试包含示例的情况
func TestGenerateHelp_WithExamples(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.UseChinese = true
	ctx.Config.Examples = []types.ExampleInfo{
		{Desc: "基本用法", Usage: "test --help"},
		{Desc: "详细模式", Usage: "test --verbose"},
	}

	result := GenerateHelp(ctx)
	if !strings.Contains(result, "示例:") {
		t.Error("结果应包含示例标题")
	}
	if !strings.Contains(result, "基本用法") {
		t.Error("结果应包含示例描述")
	}
	if !strings.Contains(result, "test --help") {
		t.Error("结果应包含示例用法")
	}
}

// TestGenerateHelp_WithNotes 测试包含注意事项的情况
func TestGenerateHelp_WithNotes(t *testing.T) {
	ctx := createTestContext("test", "t")
	ctx.Config.UseChinese = true
	ctx.Config.Notes = []string{
		"这是第一个注意事项",
		"这是第二个注意事项",
	}

	result := GenerateHelp(ctx)
	if !strings.Contains(result, "注意事项:") {
		t.Error("结果应包含注意事项标题")
	}
	if !strings.Contains(result, "这是第一个注意事项") {
		t.Error("结果应包含第一个注意事项")
	}
	if !strings.Contains(result, "这是第二个注意事项") {
		t.Error("结果应包含第二个注意事项")
	}
}

// TestHelpTemplate_EnglishTemplate 测试英文模板常量
func TestHelpTemplate_EnglishTemplate(t *testing.T) {
	if EnglishTemplate.CmdName != "Name: %s\n\n" {
		t.Error("英文模板的命令名称格式不正确")
	}
	if EnglishTemplate.UsagePrefix != "Usage: " {
		t.Error("英文模板的用法前缀不正确")
	}
	if EnglishTemplate.OptionsHeader != "Options:\n" {
		t.Error("英文模板的选项标题不正确")
	}
}

// TestHelpTemplate_ChineseTemplate 测试中文模板常量
func TestHelpTemplate_ChineseTemplate(t *testing.T) {
	if ChineseTemplate.CmdName != "名称: %s\n\n" {
		t.Error("中文模板的命令名称格式不正确")
	}
	if ChineseTemplate.UsagePrefix != "用法: " {
		t.Error("中文模板的用法前缀不正确")
	}
	if ChineseTemplate.OptionsHeader != "选项:\n" {
		t.Error("中文模板的选项标题不正确")
	}
}
