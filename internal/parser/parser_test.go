// Package parser 解析器测试
// 本文件包含了命令行参数解析器的单元测试，测试参数解析、
// 子命令识别等核心解析功能的正确性和稳定性。
package parser

import (
	"flag"
	"os"
	"testing"

	"gitee.com/MM-Q/qflag/internal/help"
	"gitee.com/MM-Q/qflag/internal/types"
)

// createTestContext 创建测试用的命令上下文
func createTestContext(name string) *types.CmdContext {
	ctx := types.NewCmdContext(name, "", flag.ContinueOnError)

	// 初始化必要的字段
	ctx.Config.UseChinese = false
	ctx.Config.Notes = []string{}

	return ctx
}

// TestParseCommand_基本标志解析 测试基本标志解析功能
func TestParseCommand_基本标志解析(t *testing.T) {
	ctx := createTestContext("测试")

	// 添加一些测试标志
	var stringFlag string
	var intFlag int
	var boolFlag bool

	ctx.FlagSet.StringVar(&stringFlag, "string", "默认值", "字符串标志")
	ctx.FlagSet.IntVar(&intFlag, "int", 0, "整数标志")
	ctx.FlagSet.BoolVar(&boolFlag, "bool", false, "布尔标志")

	// 测试参数
	args := []string{"--string", "测试值", "--int", "42", "--bool", "剩余", "参数"}

	err := ParseCommand(ctx, args)
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	// 验证标志值
	if stringFlag != "测试值" {
		t.Errorf("期望字符串标志为 '测试值'，实际得到 '%s'", stringFlag)
	}

	if intFlag != 42 {
		t.Errorf("期望整数标志为 42，实际得到 %d", intFlag)
	}

	if !boolFlag {
		t.Errorf("期望布尔标志为 true，实际得到 %v", boolFlag)
	}

	// 验证非标志参数
	expectedArgs := []string{"剩余", "参数"}
	if len(ctx.Args) != len(expectedArgs) {
		t.Errorf("期望 %d 个参数，实际得到 %d 个", len(expectedArgs), len(ctx.Args))
	}

	for i, expected := range expectedArgs {
		if i >= len(ctx.Args) || ctx.Args[i] != expected {
			t.Errorf("期望参数[%d] 为 '%s'，实际得到 '%s'", i, expected, ctx.Args[i])
		}
	}
}

// TestParseCommand_空参数 测试空参数处理
func TestParseCommand_空参数(t *testing.T) {
	ctx := createTestContext("测试")

	err := ParseCommand(ctx, []string{})
	if err != nil {
		t.Fatalf("空参数解析失败: %v", err)
	}

	if len(ctx.Args) != 0 {
		t.Errorf("期望没有参数，实际得到 %d 个", len(ctx.Args))
	}
}

// TestParseCommand_仅标志 测试只有标志没有参数的情况
func TestParseCommand_仅标志(t *testing.T) {
	ctx := createTestContext("测试")

	var testFlag string
	ctx.FlagSet.StringVar(&testFlag, "test", "默认值", "测试标志")

	args := []string{"--test", "值"}

	err := ParseCommand(ctx, args)
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	if testFlag != "值" {
		t.Errorf("期望测试标志为 '值'，实际得到 '%s'", testFlag)
	}

	if len(ctx.Args) != 0 {
		t.Errorf("期望没有剩余参数，实际得到 %d 个", len(ctx.Args))
	}
}

// TestParseCommand_仅参数 测试只有参数没有标志的情况
func TestParseCommand_仅参数(t *testing.T) {
	ctx := createTestContext("测试")

	args := []string{"参数1", "参数2", "参数3"}

	err := ParseCommand(ctx, args)
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	if len(ctx.Args) != 3 {
		t.Errorf("期望 3 个参数，实际得到 %d 个", len(ctx.Args))
	}

	expectedArgs := []string{"参数1", "参数2", "参数3"}
	for i, expected := range expectedArgs {
		if i >= len(ctx.Args) || ctx.Args[i] != expected {
			t.Errorf("期望参数[%d] 为 '%s'，实际得到 '%s'", i, expected, ctx.Args[i])
		}
	}
}

// TestParseCommand_无效标志 测试无效标志处理
func TestParseCommand_无效标志(t *testing.T) {
	ctx := createTestContext("测试")

	// 不添加任何标志定义
	args := []string{"--不存在的标志", "值"}

	err := ParseCommand(ctx, args)
	if err == nil {
		t.Fatal("期望 ParseCommand 因无效标志失败，但实际成功了")
	}

	// 验证错误类型
	if err.Error() == "" {
		t.Error("期望非空错误消息")
	}
}

// TestParseCommand_中文注释 测试中文注释添加
func TestParseCommand_中文注释(t *testing.T) {
	ctx := createTestContext("测试")
	ctx.Config.UseChinese = true

	// 清空现有注释
	ctx.Config.Notes = []string{}

	err := ParseCommand(ctx, []string{})
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	// 验证是否添加了中文注释
	if len(ctx.Config.Notes) == 0 {
		t.Error("期望添加中文注释，但没有找到")
	}

	// 验证注释内容是否为中文模板
	expectedNote := help.ChineseTemplate.DefaultNote
	found := false
	for _, note := range ctx.Config.Notes {
		if note == expectedNote {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("期望在注释中找到中文默认注释 '%s'", expectedNote)
	}
}

// TestParseCommand_英文注释 测试英文注释添加
func TestParseCommand_英文注释(t *testing.T) {
	ctx := createTestContext("测试")
	ctx.Config.UseChinese = false

	// 清空现有注释
	ctx.Config.Notes = []string{}

	err := ParseCommand(ctx, []string{})
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	// 验证是否添加了英文注释
	if len(ctx.Config.Notes) == 0 {
		t.Error("期望添加英文注释，但没有找到")
	}

	// 验证注释内容是否为英文模板
	expectedNote := help.EnglishTemplate.DefaultNote
	found := false
	for _, note := range ctx.Config.Notes {
		if note == expectedNote {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("期望在注释中找到英文默认注释 '%s'", expectedNote)
	}
}

// TestParseCommand_环境变量 测试环境变量加载
func TestParseCommand_环境变量(t *testing.T) {
	ctx := createTestContext("测试")

	// 设置环境变量
	envKey := "TEST_FLAG"
	envValue := "环境变量值"
	if err := os.Setenv(envKey, envValue); err != nil {
		t.Fatalf("设置环境变量失败: %v", err)
	}
	defer func() {
		if err := os.Unsetenv(envKey); err != nil {
			t.Logf("清除环境变量失败: %v", err)
		}
	}()

	// 添加对应的标志
	var testFlag string
	ctx.FlagSet.StringVar(&testFlag, "flag", "默认值", "测试标志")

	// 主要测试 LoadEnvVars 函数是否被调用而不出错
	err := ParseCommand(ctx, []string{})
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	// 注意：实际的环境变量处理逻辑在 LoadEnvVars 函数中
	// 这里主要验证函数调用不会出错
}

// TestParseCommand_标志解析错误 测试标志解析错误处理
func TestParseCommand_标志解析错误(t *testing.T) {
	ctx := createTestContext("测试")

	// 添加一个整数标志
	var intFlag int
	ctx.FlagSet.IntVar(&intFlag, "number", 0, "数字标志")

	// 传入无效的整数值
	args := []string{"--number", "不是数字"}

	err := ParseCommand(ctx, args)
	if err == nil {
		t.Fatal("期望 ParseCommand 因无效整数值失败，但实际成功了")
	}

	// 验证错误消息包含标志解析失败的信息
	if err.Error() == "" {
		t.Error("期望非空错误消息")
	}
}

// TestParseCommand_多种标志组合 测试多种标志组合
func TestParseCommand_多种标志组合(t *testing.T) {
	ctx := createTestContext("测试")

	// 添加多种类型的标志
	var (
		stringFlag string
		intFlag    int
		boolFlag   bool
		floatFlag  float64
	)

	ctx.FlagSet.StringVar(&stringFlag, "string", "", "字符串标志")
	ctx.FlagSet.StringVar(&stringFlag, "s", "", "字符串标志简写")
	ctx.FlagSet.IntVar(&intFlag, "int", 0, "整数标志")
	ctx.FlagSet.BoolVar(&boolFlag, "bool", false, "布尔标志")
	ctx.FlagSet.Float64Var(&floatFlag, "float", 0.0, "浮点数标志")

	args := []string{
		"-s", "简写值",
		"--int", "123",
		"--bool",
		"--float", "3.14",
		"剩余", "参数",
	}

	err := ParseCommand(ctx, args)
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	// 验证所有标志值
	if stringFlag != "简写值" {
		t.Errorf("期望字符串标志为 '简写值'，实际得到 '%s'", stringFlag)
	}

	if intFlag != 123 {
		t.Errorf("期望整数标志为 123，实际得到 %d", intFlag)
	}

	if !boolFlag {
		t.Errorf("期望布尔标志为 true，实际得到 %v", boolFlag)
	}

	if floatFlag != 3.14 {
		t.Errorf("期望浮点数标志为 3.14，实际得到 %f", floatFlag)
	}

	// 验证剩余参数
	expectedArgs := []string{"剩余", "参数"}
	if len(ctx.Args) != len(expectedArgs) {
		t.Errorf("期望 %d 个剩余参数，实际得到 %d 个", len(expectedArgs), len(ctx.Args))
	}
}

// TestParseCommand_参数累积 测试参数累积功能
func TestParseCommand_参数累积(t *testing.T) {
	ctx := createTestContext("测试")

	// 预先添加一些参数
	ctx.Args = []string{"现有", "参数"}

	args := []string{"新", "参数"}

	err := ParseCommand(ctx, args)
	if err != nil {
		t.Fatalf("ParseCommand 解析失败: %v", err)
	}

	// 验证参数是否正确累积
	expectedArgs := []string{"现有", "参数", "新", "参数"}
	if len(ctx.Args) != len(expectedArgs) {
		t.Errorf("期望总共 %d 个参数，实际得到 %d 个", len(expectedArgs), len(ctx.Args))
	}

	for i, expected := range expectedArgs {
		if i >= len(ctx.Args) || ctx.Args[i] != expected {
			t.Errorf("期望参数[%d] 为 '%s'，实际得到 '%s'", i, expected, ctx.Args[i])
		}
	}
}

// TestParseCommand_Panic恢复 测试panic恢复机制
func TestParseCommand_Panic恢复(t *testing.T) {
	ctx := createTestContext("测试")

	// 创建一个会导致panic的FlagSet
	ctx.FlagSet = nil

	err := ParseCommand(ctx, []string{})
	if err == nil {
		t.Fatal("期望 ParseCommand 因panic而失败，但实际成功了")
	}

	// 验证错误消息包含panic恢复信息
	if err.Error() == "" {
		t.Error("期望非空错误消息")
	}
}

// BenchmarkParseCommand 性能基准测试
func BenchmarkParseCommand(b *testing.B) {
	ctx := createTestContext("基准测试")

	// 添加一些标志
	var (
		stringFlag string
		intFlag    int
		boolFlag   bool
	)

	ctx.FlagSet.StringVar(&stringFlag, "string", "", "字符串标志")
	ctx.FlagSet.IntVar(&intFlag, "int", 0, "整数标志")
	ctx.FlagSet.BoolVar(&boolFlag, "bool", false, "布尔标志")

	args := []string{"--string", "测试", "--int", "42", "--bool", "参数1", "参数2"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 重置上下文状态
		ctx.Args = ctx.Args[:0]
		ctx.Config.Notes = ctx.Config.Notes[:0]
		ctx.FlagSet = flag.NewFlagSet("基准测试", flag.ContinueOnError)
		ctx.FlagSet.StringVar(&stringFlag, "string", "", "字符串标志")
		ctx.FlagSet.IntVar(&intFlag, "int", 0, "整数标志")
		ctx.FlagSet.BoolVar(&boolFlag, "bool", false, "布尔标志")

		err := ParseCommand(ctx, args)
		if err != nil {
			b.Fatalf("ParseCommand 解析失败: %v", err)
		}
	}
}
