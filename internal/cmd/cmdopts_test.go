package cmd

import (
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

func TestNewCmdOpts(t *testing.T) {
	opts := NewCmdOpts()

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
	cmd := NewCmd("test", "t", types.ExitOnError)
	opts := &CmdOpts{
		Desc:    "test command",
		Version: "1.0.0",
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	if cmd.Desc() != "test command" {
		t.Errorf("Expected Desc 'test command', got '%s'", cmd.Desc())
	}

	if cmd.Config().Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", cmd.Config().Version)
	}
}

func TestCmd_ApplyOpts_PartialConfig(t *testing.T) {
	cmd := NewCmd("test", "t", types.ExitOnError)
	cmd.SetDesc("original desc")

	opts := &CmdOpts{
		Version: "1.0.0",
	}

	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts failed: %v", err)
	}

	if cmd.Desc() != "original desc" {
		t.Errorf("Expected Desc 'original desc', got '%s'", cmd.Desc())
	}

	if cmd.Config().Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", cmd.Config().Version)
	}
}

func TestCmd_ApplyOpts_NilOpts(t *testing.T) {
	cmd := NewCmd("test", "t", types.ExitOnError)
	var opts *CmdOpts

	err := cmd.ApplyOpts(opts)
	if err == nil {
		t.Fatal("Expected error for nil opts, got nil")
	}
}

func TestCmd_ApplyOpts_AllFields(t *testing.T) {
	cmd := NewCmd("test", "t", types.ExitOnError)

	if err := cmd.AddFlag(flag.NewStringFlag("json", "j", "JSON format", "")); err != nil {
		t.Fatalf("Failed to add json flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("xml", "x", "XML format", "")); err != nil {
		t.Fatalf("Failed to add xml flag: %v", err)
	}

	opts := &CmdOpts{
		Desc: "test command",
		RunFunc: func(c types.Command) error {
			return nil
		},
		Version:     "1.0.0",
		UseChinese:  true,
		EnvPrefix:   "TEST",
		UsageSyntax: "test [options]",
		LogoText:    "Test Logo",
		Examples: map[string]string{
			"example1": "test --help",
			"example2": "test --version",
		},
		Notes: []string{
			"note1",
			"note2",
		},
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

	if cmd.Desc() != "test command" {
		t.Errorf("Expected Desc 'test command', got '%s'", cmd.Desc())
	}

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

	if len(cmd.Config().Example) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(cmd.Config().Example))
	}

	if cmd.Config().Example["example1"] != "test --help" {
		t.Errorf("Expected example 'test --help', got '%s'", cmd.Config().Example["example1"])
	}

	if len(cmd.Config().Notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(cmd.Config().Notes))
	}

	groups := cmd.Config().MutexGroups
	if len(groups) != 1 {
		t.Errorf("Expected 1 mutex group, got %d", len(groups))
	}

	if groups[0].Name != "format" {
		t.Errorf("Expected mutex group name 'format', got '%s'", groups[0].Name)
	}
}

func TestCmd_ApplyOpts_SubCmds(t *testing.T) {
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
	cmd := NewCmd("test", "t", types.ExitOnError)
	opts := &CmdOpts{
		Desc: "test command",
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := cmd.ApplyOpts(opts)
			if err != nil {
				t.Errorf("ApplyOpts failed: %v", err)
			}
		}()
	}
	wg.Wait()
}
