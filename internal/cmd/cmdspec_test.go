package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

func TestNewCmdSpec(t *testing.T) {
	// 测试基本创建
	spec := NewCmdSpec("test", "t")

	if spec.LongName != "test" {
		t.Errorf("Expected LongName 'test', got '%s'", spec.LongName)
	}

	if spec.ShortName != "t" {
		t.Errorf("Expected ShortName 't', got '%s'", spec.ShortName)
	}

	if spec.ErrorHandling != types.ExitOnError {
		t.Errorf("Expected ErrorHandling ExitOnError, got %v", spec.ErrorHandling)
	}

	if spec.UseChinese != false {
		t.Errorf("Expected UseChinese false, got %v", spec.UseChinese)
	}

	if spec.Examples == nil {
		t.Error("Examples map should be initialized")
	}

	if spec.Notes == nil {
		t.Error("Notes slice should be initialized")
	}

	if spec.SubCmds == nil {
		t.Error("SubCmds slice should be initialized")
	}

	if spec.MutexGroups == nil {
		t.Error("MutexGroups slice should be initialized")
	}
}

func TestNewCmdFromSpec(t *testing.T) {
	// 测试基本命令创建
	spec := NewCmdSpec("myapp", "app")
	spec.Desc = "我的应用程序"
	spec.Version = "1.0.0"
	spec.UseChinese = true
	spec.EnvPrefix = "MYAPP"
	spec.RunFunc = func(cmd types.Command) error {
		return nil
	}

	// 添加互斥组
	spec.MutexGroups = []types.MutexGroup{
		{Name: "format", Flags: []string{"json", "xml"}, AllowNone: true},
	}

	cmd, err := NewCmdFromSpec(spec)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}

	// 添加标志
	cmd.String("input", "i", "输入文件", "")
	cmd.Bool("verbose", "v", "详细输出", false)

	if cmd.Name() != "myapp" {
		t.Errorf("Expected command name 'myapp', got '%s'", cmd.Name())
	}

	if cmd.Desc() != "我的应用程序" {
		t.Errorf("Expected description '我的应用程序', got '%s'", cmd.Desc())
	}

	if !cmd.HasRunFunc() {
		t.Error("Command should have run function")
	}

	// 检查标志
	inputFlag, exists := cmd.GetFlag("input")
	if !exists {
		t.Error("Input flag should exist")
	} else if inputFlag.Desc() != "输入文件" {
		t.Errorf("Expected input flag description '输入文件', got '%s'", inputFlag.Desc())
	}

	verboseFlag, exists := cmd.GetFlag("verbose")
	if !exists {
		t.Error("Verbose flag should exist")
	} else if verboseFlag.Desc() != "详细输出" {
		t.Errorf("Expected verbose flag description '详细输出', got '%s'", verboseFlag.Desc())
	}

	// 检查互斥组
	config := cmd.Config()
	if len(config.MutexGroups) != 1 {
		t.Errorf("Expected 1 mutex group, got %d", len(config.MutexGroups))
	}

	if config.MutexGroups[0].Name != "format" {
		t.Errorf("Expected mutex group name 'format', got '%s'", config.MutexGroups[0].Name)
	}
}

func TestNewCmdFromSpecWithSubCommands(t *testing.T) {
	// 创建主命令规格
	mainSpec := NewCmdSpec("main", "m")
	mainSpec.Desc = "主命令"
	mainSpec.RunFunc = func(cmd types.Command) error {
		return nil
	}

	// 创建子命令
	subCmd := NewCmd("sub", "s", types.ExitOnError)
	subCmd.SetDesc("子命令")

	// 添加到主命令
	mainSpec.SubCmds = []types.Command{subCmd}

	// 创建命令
	cmd, err := NewCmdFromSpec(mainSpec)
	if err != nil {
		t.Fatalf("Failed to create command: %v", err)
	}

	// 获取子命令
	retrievedSubCmd, exists := cmd.GetSubCmd("sub")
	if !exists {
		t.Error("Subcommand 'sub' should exist")
	}

	// 检查子命令
	if retrievedSubCmd.Desc() != "子命令" {
		t.Errorf("Expected subcommand description '子命令', got '%s'", retrievedSubCmd.Desc())
	}
}

func TestNewCmdFromSpecNilSpec(t *testing.T) {
	// 测试nil规格
	cmd, err := NewCmdFromSpec(nil)
	if err == nil {
		t.Error("Expected error for nil spec")
	}
	if cmd != nil {
		t.Error("Expected nil command for nil spec")
	}
}

func TestNewCmdFromSpecComplex(t *testing.T) {
	// 复杂命令配置
	complexSpec := NewCmdSpec("complex", "cpx")
	complexSpec.Desc = "复杂命令示例"
	complexSpec.Version = "2.0.0"
	complexSpec.UseChinese = true
	complexSpec.EnvPrefix = "COMPLEX"
	complexSpec.UsageSyntax = "[options] <args>"
	complexSpec.LogoText = "Complex Command v2.0.0"

	// 添加示例
	complexSpec.Examples = map[string]string{
		"基本用法":   "complex --input file.txt",
		"详细模式":   "complex --input file.txt --verbose",
		"输出JSON": "complex --input file.txt --json",
	}

	// 添加注意事项
	complexSpec.Notes = []string{
		"输入文件必须存在",
		"输出目录必须可写",
		"处理大文件时请增加内存限制",
	}

	// 添加子命令
	processCmd := NewCmd("process", "proc", types.ExitOnError)
	processCmd.SetDesc("处理数据")

	validateCmd := NewCmd("validate", "val", types.ExitOnError)
	validateCmd.SetDesc("验证数据")

	complexSpec.SubCmds = []types.Command{processCmd, validateCmd}

	// 创建命令
	cmd, err := NewCmdFromSpec(complexSpec)
	if err != nil {
		t.Fatalf("Failed to create complex command: %v", err)
	}

	// 添加多个标志
	cmd.String("input", "i", "输入文件", "")
	cmd.String("output", "o", "输出文件", "")
	cmd.Bool("verbose", "v", "详细输出", false)
	cmd.Bool("json", "j", "JSON格式输出", false)
	cmd.Bool("xml", "x", "XML格式输出", false)
	cmd.Int("limit", "l", "处理限制", 1000)

	// 验证基本属性
	if cmd.Name() != "complex" {
		t.Errorf("Expected command name 'complex', got '%s'", cmd.Name())
	}

	if cmd.Desc() != "复杂命令示例" {
		t.Errorf("Expected description '复杂命令示例', got '%s'", cmd.Desc())
	}

	// 验证配置
	config := cmd.Config()
	if config.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", config.Version)
	}

	if !config.UseChinese {
		t.Error("Expected UseChinese to be true")
	}

	if config.EnvPrefix != "COMPLEX_" {
		t.Errorf("Expected env prefix 'COMPLEX_', got '%s'", config.EnvPrefix)
	}

	// 验证示例
	if len(config.Example) != 3 {
		t.Errorf("Expected 3 examples, got %d", len(config.Example))
	}

	if config.Example["基本用法"] != "complex --input file.txt" {
		t.Error("Example '基本用法' not set correctly")
	}

	// 验证注意事项
	if len(config.Notes) != 3 {
		t.Errorf("Expected 3 notes, got %d", len(config.Notes))
	}

	if config.Notes[0] != "输入文件必须存在" {
		t.Error("First note not set correctly")
	}

	// 验证标志数量
	flags := cmd.Flags()
	if len(flags) != 6 {
		t.Errorf("Expected 6 flags, got %d", len(flags))
	}

	// 验证子命令
	subCmds := cmd.SubCmds()
	if len(subCmds) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(subCmds))
	}
}
