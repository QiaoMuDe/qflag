package cmd

import (
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

func TestNewCmdOpts(t *testing.T) {
	// 测试基本创建
	opts := NewCmdOpts()

	// 验证字段初始化
	if opts.Examples == nil {
		t.Error("Examples map should be initialized")
	}

	if opts.Notes == nil {
		t.Error("Notes slice should be initialized")
	}

	if opts.SubCmds == nil {
		t.Error("SubCmds slice should be initialized")
	}

	if opts.MutexGroups == nil {
		t.Error("MutexGroups slice should be initialized")
	}

	// 验证零值
	if opts.Desc != "" {
		t.Errorf("Expected empty Desc, got '%s'", opts.Desc)
	}

	if opts.Version != "" {
		t.Errorf("Expected empty Version, got '%s'", opts.Version)
	}

	if opts.UseChinese != false {
		t.Errorf("Expected UseChinese false, got %v", opts.UseChinese)
	}
}

func TestCmd_ApplyOpts_BasicProperties(t *testing.T) {
	// 测试基本属性设置
	cmd := NewCmd("test", "t", types.ExitOnError)
	opts := &CmdOpts{
		Desc:    "测试命令",
		Version: "1.0.0",
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	if cmd.Desc() != "测试命令" {
		t.Errorf("Expected Desc '测试命令', got '%s'", cmd.Desc())
	}

	if cmd.Config().Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", cmd.Config().Version)
	}
}

func TestCmd_ApplyOpts_PartialConfig(t *testing.T) {
	// 测试部分配置
	cmd := NewCmd("test", "t", types.ExitOnError)
	cmd.SetDesc("原始描述")

	opts := &CmdOpts{
		Version: "1.0.0",
		// 不设置 Desc
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	if cmd.Desc() != "原始描述" {
		t.Errorf("Expected Desc '原始描述', got '%s'", cmd.Desc())
	}

	if cmd.Config().Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", cmd.Config().Version)
	}
}

func TestCmd_ApplyOpts_NilOpts(t *testing.T) {
	// 测试 nil 选项
	cmd := NewCmd("test", "t", types.ExitOnError)
	var opts *CmdOpts

	err := cmd.ApplyOpts(opts)
	if err == nil {
		t.Fatal("Expected error for nil opts, got nil")
	}

	if err.(*types.Error).Code != "INVALID_CMDOPTS" {
		t.Errorf("Expected error code 'INVALID_CMDOPTS', got '%s'", err.(*types.Error).Code)
	}
}

func TestCmd_ApplyOpts_AllFields(t *testing.T) {
	// 测试所有字段
	cmd := NewCmd("test", "t", types.ExitOnError)

	opts := &CmdOpts{
		// 基本属性
		Desc: "测试命令",
		RunFunc: func(c types.Command) error {
			return nil
		},

		// 配置选项
		Version:     "1.0.0",
		UseChinese:  true,
		EnvPrefix:   "TEST",
		UsageSyntax: "test [options]",
		LogoText:    "Test Logo",

		// 示例和说明
		Examples: map[string]string{
			"示例1": "test --help",
			"示例2": "test --version",
		},
		Notes: []string{
			"注意1",
			"注意2",
		},

		// 子命令和互斥组
		MutexGroups: []types.MutexGroup{
			{
				Name:      "format",
				Flags:     []string{"json", "xml"},
				AllowNone: false,
			},
		},
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	// 验证基本属性
	if cmd.Desc() != "测试命令" {
		t.Errorf("Expected Desc '测试命令', got '%s'", cmd.Desc())
	}

	// 验证配置选项
	if cmd.Config().Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", cmd.Config().Version)
	}

	if cmd.Config().UseChinese != true {
		t.Errorf("Expected UseChinese true, got %v", cmd.Config().UseChinese)
	}

	if cmd.Config().EnvPrefix != "TEST_" {
		t.Errorf("Expected EnvPrefix 'TEST_', got '%s'", cmd.Config().EnvPrefix)
	}

	if cmd.Config().UsageSyntax != "test [options]" {
		t.Errorf("Expected UsageSyntax 'test [options]', got '%s'", cmd.Config().UsageSyntax)
	}

	if cmd.Config().LogoText != "Test Logo" {
		t.Errorf("Expected LogoText 'Test Logo', got '%s'", cmd.Config().LogoText)
	}

	// 验证示例和说明
	if len(cmd.Config().Example) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(cmd.Config().Example))
	}

	if cmd.Config().Example["示例1"] != "test --help" {
		t.Errorf("Expected example 'test --help', got '%s'", cmd.Config().Example["示例1"])
	}

	if cmd.Config().Example["示例2"] != "test --version" {
		t.Errorf("Expected example 'test --version', got '%s'", cmd.Config().Example["示例2"])
	}

	if len(cmd.Config().Notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(cmd.Config().Notes))
	}

	if cmd.Config().Notes[0] != "注意1" {
		t.Errorf("Expected note '注意1', got '%s'", cmd.Config().Notes[0])
	}

	if cmd.Config().Notes[1] != "注意2" {
		t.Errorf("Expected note '注意2', got '%s'", cmd.Config().Notes[1])
	}

	// 验证互斥组
	groups := cmd.GetMutexGroups()
	if len(groups) != 1 {
		t.Errorf("Expected 1 mutex group, got %d", len(groups))
	}

	if groups[0].Name != "format" {
		t.Errorf("Expected mutex group name 'format', got '%s'", groups[0].Name)
	}

	if len(groups[0].Flags) != 2 {
		t.Errorf("Expected 2 flags, got %d", len(groups[0].Flags))
	}

	if groups[0].Flags[0] != "json" {
		t.Errorf("Expected flag 'json', got '%s'", groups[0].Flags[0])
	}

	if groups[0].Flags[1] != "xml" {
		t.Errorf("Expected flag 'xml', got '%s'", groups[0].Flags[1])
	}

	if groups[0].AllowNone != false {
		t.Errorf("Expected AllowNone false, got %v", groups[0].AllowNone)
	}
}

func TestCmd_ApplyOpts_SubCmds(t *testing.T) {
	// 测试子命令
	cmd := NewCmd("test", "t", types.ExitOnError)

	subCmd1 := NewCmd("sub1", "s1", types.ExitOnError)
	subCmd2 := NewCmd("sub2", "s2", types.ExitOnError)

	opts := &CmdOpts{
		SubCmds: []types.Command{subCmd1, subCmd2},
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	if !cmd.HasSubCmd("sub1") {
		t.Error("Expected subcommand 'sub1' to be added")
	}

	if !cmd.HasSubCmd("sub2") {
		t.Error("Expected subcommand 'sub2' to be added")
	}
}

func TestCmd_ApplyOpts_ConcurrentSafety(t *testing.T) {
	// 测试并发安全
	cmd := NewCmd("test", "t", types.ExitOnError)
	opts := &CmdOpts{
		Desc: "测试命令",
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cmd.ApplyOpts(opts)
		}()
	}
	wg.Wait()

	if cmd.Desc() != "测试命令" {
		t.Errorf("Expected Desc '测试命令', got '%s'", cmd.Desc())
	}
}

func TestCmd_ApplyOpts_EmptyValues(t *testing.T) {
	// 测试空值不会被应用
	cmd := NewCmd("test", "t", types.ExitOnError)
	cmd.SetDesc("原始描述")
	cmd.SetVersion("0.0.0")

	opts := &CmdOpts{
		// 空值不应该覆盖原有值
		Desc:    "",
		Version: "",
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	if cmd.Desc() != "原始描述" {
		t.Errorf("Expected Desc '原始描述', got '%s'", cmd.Desc())
	}

	if cmd.Config().Version != "0.0.0" {
		t.Errorf("Expected Version '0.0.0', got '%s'", cmd.Config().Version)
	}
}
