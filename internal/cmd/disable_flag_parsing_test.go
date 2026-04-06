package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

/*
本测试文件验证禁用标志解析功能 (DisableFlagParsing) 的各种场景，包括：

1. 基本禁用功能验证：
   - 验证当 DisableFlagParsing = true 时，所有参数都作为位置参数
   - 验证 --flag 形式的参数不被解析，而是作为位置参数保留
   - 验证 -f 形式的参数不被解析，而是作为位置参数保留

2. 子命令路由验证：
   - 验证禁用标志解析不影响子命令路由
   - 验证子命令可以正常被识别和执行
   - 验证子命令的参数也被正确处理

3. 嵌套子命令验证：
   - 验证多层嵌套子命令在禁用标志解析时的行为
   - 验证每一层都可以独立控制是否禁用标志解析

4. ParseOnly/Parse/ParseAndRoute 一致性验证：
   - 验证三种解析方法在禁用标志解析时的行为一致
   - 确保参数正确传递和设置

测试覆盖了所有主要功能路径，确保禁用标志解析功能的正确性。
*/

// TestDisableFlagParsing_Basic 测试基本的禁用标志解析功能
func TestDisableFlagParsing_Basic(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantParsed bool
		wantArgs   []string
	}{
		{
			name:       "普通参数",
			args:       []string{"arg1", "arg2", "arg3"},
			wantParsed: true,
			wantArgs:   []string{"arg1", "arg2", "arg3"},
		},
		{
			name:       "双横线参数被当作位置参数",
			args:       []string{"--flag", "value", "--another"},
			wantParsed: true,
			wantArgs:   []string{"--flag", "value", "--another"},
		},
		{
			name:       "单横线参数被当作位置参数",
			args:       []string{"-f", "value", "-a"},
			wantParsed: true,
			wantArgs:   []string{"-f", "value", "-a"},
		},
		{
			name:       "混合参数",
			args:       []string{"--flag", "-f", "pos1", "--", "pos2"},
			wantParsed: true,
			wantArgs:   []string{"--flag", "-f", "pos1", "--", "pos2"},
		},
		{
			name:       "空参数",
			args:       []string{},
			wantParsed: true,
			wantArgs:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建命令并禁用标志解析
			cmd := NewCmd("test", "t", types.ContinueOnError)
			cmd.SetDisableFlagParsing(true)

			// 添加一个标志（应该不会被解析）
			_ = cmd.String("flag", "f", "test flag", "default")

			// 使用 ParseOnly 解析
			err := cmd.ParseOnly(tt.args)
			if err != nil {
				t.Fatalf("ParseOnly() error = %v", err)
			}

			// 验证 parsed 状态
			if cmd.IsParsed() != tt.wantParsed {
				t.Errorf("IsParsed() = %v, want %v", cmd.IsParsed(), tt.wantParsed)
			}

			// 验证 args
			gotArgs := cmd.Args()
			if len(gotArgs) != len(tt.wantArgs) {
				t.Errorf("Args() = %v, want %v", gotArgs, tt.wantArgs)
			} else {
				for i := range gotArgs {
					if gotArgs[i] != tt.wantArgs[i] {
						t.Errorf("Args()[%d] = %v, want %v", i, gotArgs[i], tt.wantArgs[i])
					}
				}
			}

			// 验证标志没有被设置（因为禁用了标志解析）
			flag, _ := cmd.FlagRegistry().Get("flag")
			if flag.IsSet() {
				t.Error("Flag should not be set when flag parsing is disabled")
			}
		})
	}
}

// TestDisableFlagParsing_WithSubCommand 测试禁用标志解析时子命令路由
func TestDisableFlagParsing_WithSubCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantSubCmd  string // 期望路由到的子命令
		wantSubArgs []string
	}{
		{
			name:        "子命令正常路由",
			args:        []string{"sub", "arg1", "arg2"},
			wantSubCmd:  "sub",
			wantSubArgs: []string{"arg1", "arg2"},
		},
		{
			name:        "带标志形式的子命令路由",
			args:        []string{"sub", "--flag", "value"},
			wantSubCmd:  "sub",
			wantSubArgs: []string{"--flag", "value"},
		},
		{
			name:        "无子命令只有位置参数",
			args:        []string{"--flag", "value", "arg1"},
			wantSubCmd:  "",
			wantSubArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建父命令并禁用标志解析
			parent := NewCmd("parent", "p", types.ContinueOnError)
			parent.SetDisableFlagParsing(true)

			// 创建子命令
			sub := NewCmd("sub", "s", types.ContinueOnError)
			sub.SetDisableFlagParsing(true)
			if err := parent.AddSubCmds(sub); err != nil {
				t.Fatalf("AddSubCmds() error = %v", err)
			}

			// 使用 Parse 解析（会路由子命令）
			err := parent.Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if tt.wantSubCmd == "" {
				// 没有子命令，验证父命令的参数
				gotArgs := parent.Args()
				if len(gotArgs) != len(tt.args) {
					t.Errorf("Parent Args() = %v, want %v", gotArgs, tt.args)
				}
			} else {
				// 验证子命令被正确路由且参数传递正确
				if !sub.IsParsed() {
					t.Error("Subcommand should be parsed")
				}
				gotArgs := sub.Args()
				if len(gotArgs) != len(tt.wantSubArgs) {
					t.Errorf("Subcommand Args() = %v, want %v", gotArgs, tt.wantSubArgs)
				} else {
					for i := range gotArgs {
						if gotArgs[i] != tt.wantSubArgs[i] {
							t.Errorf("Subcommand Args()[%d] = %v, want %v", i, gotArgs[i], tt.wantSubArgs[i])
						}
					}
				}
			}
		})
	}
}

// TestDisableFlagParsing_NestedSubCommands 测试嵌套子命令
func TestDisableFlagParsing_NestedSubCommands(t *testing.T) {
	// 创建三层命令结构：parent -> child -> grandchild
	parent := NewCmd("parent", "p", types.ContinueOnError)
	parent.SetDisableFlagParsing(true) // 父命令禁用标志解析

	child := NewCmd("child", "c", types.ContinueOnError)
	child.SetDisableFlagParsing(false) // 子命令不禁用

	grandchild := NewCmd("grandchild", "g", types.ContinueOnError)
	grandchild.SetDisableFlagParsing(true) // 孙命令禁用

	if err := parent.AddSubCmds(child); err != nil {
		t.Fatalf("parent.AddSubCmds() error = %v", err)
	}
	if err := child.AddSubCmds(grandchild); err != nil {
		t.Fatalf("child.AddSubCmds() error = %v", err)
	}

	// 在孙命令添加一个标志
	_ = grandchild.String("flag", "f", "test flag", "default")

	// 测试：parent child grandchild --flag value
	args := []string{"child", "grandchild", "--flag", "value"}
	err := parent.Parse(args)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// 验证孙命令被正确路由
	if !grandchild.IsParsed() {
		t.Error("Grandchild should be parsed")
	}

	// 验证孙命令的参数（因为孙命令禁用了标志解析，--flag 应该作为位置参数）
	wantArgs := []string{"--flag", "value"}
	gotArgs := grandchild.Args()
	if len(gotArgs) != len(wantArgs) {
		t.Errorf("Grandchild Args() = %v, want %v", gotArgs, wantArgs)
	} else {
		for i := range gotArgs {
			if gotArgs[i] != wantArgs[i] {
				t.Errorf("Grandchild Args()[%d] = %v, want %v", i, gotArgs[i], wantArgs[i])
			}
		}
	}

	// 验证标志没有被设置
	flag, _ := grandchild.FlagRegistry().Get("flag")
	if flag.IsSet() {
		t.Error("Flag should not be set when flag parsing is disabled")
	}
}

// TestDisableFlagParsing_ParseAndRoute 测试 ParseAndRoute 方法
func TestDisableFlagParsing_ParseAndRoute(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantRun    bool // 是否期望执行父命令运行函数
		wantSubRun bool // 是否期望执行子命令运行函数
	}{
		{
			name:    "无子命令执行当前命令",
			args:    []string{"--flag", "value"},
			wantRun: true,
		},
		{
			name:       "有子命令执行子命令",
			args:       []string{"sub", "--flag"},
			wantRun:    false,
			wantSubRun: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建父命令并禁用标志解析
			parentRun := false
			parent := NewCmd("parent", "p", types.ContinueOnError)
			parent.SetDisableFlagParsing(true)
			parent.SetRun(func(cmd types.Command) error {
				parentRun = true
				return nil
			})

			// 创建子命令
			subRun := false
			sub := NewCmd("sub", "s", types.ContinueOnError)
			sub.SetDisableFlagParsing(true)
			sub.SetRun(func(cmd types.Command) error {
				subRun = true
				return nil
			})
			if err := parent.AddSubCmds(sub); err != nil {
				t.Fatalf("AddSubCmds() error = %v", err)
			}

			// 使用 ParseAndRoute
			err := parent.ParseAndRoute(tt.args)
			if err != nil {
				t.Fatalf("ParseAndRoute() error = %v", err)
			}

			if parentRun != tt.wantRun {
				t.Errorf("Parent Run = %v, want %v", parentRun, tt.wantRun)
			}
			if subRun != tt.wantSubRun {
				t.Errorf("Sub Run = %v, want %v", subRun, tt.wantSubRun)
			}
		})
	}
}

// TestDisableFlagParsing_NotDisabled 测试不禁用标志解析时的正常行为
func TestDisableFlagParsing_NotDisabled(t *testing.T) {
	// 创建命令（不禁用标志解析）
	cmd := NewCmd("test", "t", types.ContinueOnError)
	cmd.SetDisableFlagParsing(false)

	// 添加一个标志
	flag := cmd.String("flag", "f", "test flag", "default")

	// 解析带标志的参数
	args := []string{"--flag", "value", "pos1"}
	err := cmd.ParseOnly(args)
	if err != nil {
		t.Fatalf("ParseOnly() error = %v", err)
	}

	// 验证标志被正确解析
	if !flag.IsSet() {
		t.Error("Flag should be set")
	}

	// 验证标志值
	flagValue := flag.Get()
	if flagValue != "value" {
		t.Errorf("Flag value = %v, want value", flagValue)
	}

	// 验证位置参数正确
	wantArgs := []string{"pos1"}
	gotArgs := cmd.Args()
	if len(gotArgs) != len(wantArgs) {
		t.Errorf("Args() = %v, want %v", gotArgs, wantArgs)
	}
}

// TestDisableFlagParsing_MethodConsistency 测试三种解析方法的一致性
func TestDisableFlagParsing_MethodConsistency(t *testing.T) {
	args := []string{"--flag", "value", "arg1"}

	// 测试 ParseOnly
	t.Run("ParseOnly", func(t *testing.T) {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDisableFlagParsing(true)

		err := cmd.ParseOnly(args)
		if err != nil {
			t.Fatalf("ParseOnly() error = %v", err)
		}

		if !cmd.IsParsed() {
			t.Error("Command should be parsed")
		}

		gotArgs := cmd.Args()
		if len(gotArgs) != len(args) {
			t.Errorf("Args() = %v, want %v", gotArgs, args)
		}
	})

	// 测试 Parse（无子命令）
	t.Run("Parse", func(t *testing.T) {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDisableFlagParsing(true)

		err := cmd.Parse(args)
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		if !cmd.IsParsed() {
			t.Error("Command should be parsed")
		}

		gotArgs := cmd.Args()
		if len(gotArgs) != len(args) {
			t.Errorf("Args() = %v, want %v", gotArgs, args)
		}
	})

	// 测试 ParseAndRoute（无子命令）
	t.Run("ParseAndRoute", func(t *testing.T) {
		runCalled := false
		cmd := NewCmd("test", "t", types.ContinueOnError)
		cmd.SetDisableFlagParsing(true)
		cmd.SetRun(func(c types.Command) error {
			runCalled = true
			return nil
		})

		err := cmd.ParseAndRoute(args)
		if err != nil {
			t.Fatalf("ParseAndRoute() error = %v", err)
		}

		if !cmd.IsParsed() {
			t.Error("Command should be parsed")
		}

		if !runCalled {
			t.Error("Run function should be called")
		}

		gotArgs := cmd.Args()
		if len(gotArgs) != len(args) {
			t.Errorf("Args() = %v, want %v", gotArgs, args)
		}
	})
}

// TestDisableFlagParsing_ViaOpts 测试通过 CmdOpts 设置禁用标志解析
func TestDisableFlagParsing_ViaOpts(t *testing.T) {
	// 通过 CmdOpts 创建命令并设置禁用标志解析
	cmd := NewCmd("test", "t", types.ContinueOnError)
	opts := &CmdOpts{
		Desc:               "test command",
		DisableFlagParsing: true,
	}
	err := cmd.ApplyOpts(opts)
	if err != nil {
		t.Fatalf("ApplyOpts() error = %v", err)
	}

	// 验证禁用标志解析已设置
	if !cmd.IsDisableFlagParsing() {
		t.Error("DisableFlagParsing should be true")
	}

	// 添加标志
	_ = cmd.String("flag", "f", "test flag", "default")

	// 解析参数
	args := []string{"--flag", "value"}
	err = cmd.ParseOnly(args)
	if err != nil {
		t.Fatalf("ParseOnly() error = %v", err)
	}

	// 验证标志没有被解析
	flag, _ := cmd.FlagRegistry().Get("flag")
	if flag.IsSet() {
		t.Error("Flag should not be set when flag parsing is disabled")
	}

	// 验证参数作为位置参数保留
	gotArgs := cmd.Args()
	if len(gotArgs) != 2 || gotArgs[0] != "--flag" || gotArgs[1] != "value" {
		t.Errorf("Args() = %v, want [--flag value]", gotArgs)
	}
}
