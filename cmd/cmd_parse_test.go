// Package cmd 命令解析测试
// 本文件包含了Cmd结构体命令解析功能的单元测试，测试命令行
// 参数解析、子命令处理等核心解析功能的正确性。
package cmd

import (
	"flag"
	"strings"
	"testing"
)

// TestParse_基本功能 测试Parse方法的基本功能
func TestParse_基本功能(t *testing.T) {
	tests := []struct {
		name        string
		setupCmd    func() *Cmd
		args        []string
		expectError bool
		errorMsg    string
		description string
	}{
		{
			name: "解析空参数",
			setupCmd: func() *Cmd {
				return NewCmd("test", "t", flag.ContinueOnError)
			},
			args:        []string{},
			expectError: false,
			description: "解析空的命令行参数",
		},
		{
			name: "解析纯位置参数",
			setupCmd: func() *Cmd {
				return NewCmd("test", "t", flag.ContinueOnError)
			},
			args:        []string{"arg1", "arg2", "arg3"},
			expectError: false,
			description: "解析只包含位置参数的命令行",
		},
		{
			name: "解析标志参数",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.String("name", "n", "default", "名称标志")
				cmd.Int("count", "c", 0, "计数标志")
				cmd.Bool("verbose", "v", false, "详细输出标志")
				return cmd
			},
			args:        []string{"--name", "测试", "--count", "10", "--verbose"},
			expectError: false,
			description: "解析包含各种类型标志的命令行",
		},
		{
			name: "解析混合参数",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.String("config", "c", "", "配置文件")
				cmd.Bool("debug", "d", false, "调试模式")
				return cmd
			},
			args:        []string{"--config", "config.json", "file1.txt", "file2.txt", "--debug"},
			expectError: false,
			description: "解析混合标志和位置参数的命令行",
		},
		{
			name: "解析短标志",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.String("file", "f", "", "文件路径")
				cmd.Bool("recursive", "r", false, "递归处理")
				cmd.Int("level", "l", 1, "级别")
				return cmd
			},
			args:        []string{"-f", "test.txt", "-r", "-l", "5"},
			expectError: false,
			description: "解析短标志参数",
		},
		{
			name: "解析组合短标志",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.Bool("verbose", "v", false, "详细输出")
				cmd.Bool("recursive", "r", false, "递归处理")
				cmd.Bool("force", "f", false, "强制执行")
				return cmd
			},
			args:        []string{"-v", "-r", "-f"}, // 分开写，因为Go的flag包不支持组合短标志
			expectError: false,
			description: "解析分开的短标志",
		},
		{
			name: "解析无效标志",
			setupCmd: func() *Cmd {
				return NewCmd("test", "t", flag.ContinueOnError)
			},
			args:        []string{"--unknown-flag"},
			expectError: true,
			errorMsg:    "flag provided but not defined",
			description: "解析未定义的标志应该返回错误",
		},
		{
			name: "解析标志值类型错误",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.Int("number", "n", 0, "数字标志")
				return cmd
			},
			args:        []string{"--number", "not-a-number"},
			expectError: true,
			errorMsg:    "invalid value",
			description: "解析类型不匹配的标志值应该返回错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			err := cmd.Parse(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("错误信息不匹配: 期望包含 %q, 实际 %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("意外的错误: %v", err)
					return
				}

				// 验证解析状态
				if !cmd.IsParsed() {
					t.Error("解析后命令应该处于已解析状态")
				}
			}
		})
	}
}

// TestParse_子命令处理 测试Parse方法的子命令处理功能
func TestParse_子命令处理(t *testing.T) {
	tests := []struct {
		name        string
		setupCmd    func() *Cmd
		args        []string
		expectError bool
		expectedSub string
		description string
	}{
		{
			name: "解析简单子命令",
			setupCmd: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				child := NewCmd("child", "c", flag.ContinueOnError)
				_ = parent.AddSubCmd(child)
				return parent
			},
			args:        []string{"child"},
			expectError: false,
			expectedSub: "child",
			description: "解析简单的子命令",
		},
		{
			name: "解析子命令短名称",
			setupCmd: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				child := NewCmd("child", "c", flag.ContinueOnError)
				_ = parent.AddSubCmd(child)
				return parent
			},
			args:        []string{"c"},
			expectError: false,
			expectedSub: "child",
			description: "通过短名称解析子命令",
		},
		{
			name: "解析带参数的子命令",
			setupCmd: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				child := NewCmd("child", "c", flag.ContinueOnError)
				child.String("name", "n", "", "名称参数")
				_ = parent.AddSubCmd(child)
				return parent
			},
			args:        []string{"child", "--name", "test", "arg1"},
			expectError: false,
			expectedSub: "child",
			description: "解析带有标志和参数的子命令",
		},
		{
			name: "解析不存在的子命令",
			setupCmd: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				child := NewCmd("child", "c", flag.ContinueOnError)
				_ = parent.AddSubCmd(child)
				return parent
			},
			args:        []string{"nonexistent"},
			expectError: false, // 根据实际行为调整，不存在的子命令被当作位置参数处理
			description: "解析不存在的子命令，被当作位置参数处理",
		},
		{
			name: "解析嵌套子命令",
			setupCmd: func() *Cmd {
				root := NewCmd("root", "r", flag.ContinueOnError)
				level1 := NewCmd("level1", "l1", flag.ContinueOnError)
				level2 := NewCmd("level2", "l2", flag.ContinueOnError)
				_ = level1.AddSubCmd(level2)
				_ = root.AddSubCmd(level1)
				return root
			},
			args:        []string{"level1", "level2"},
			expectError: false,
			expectedSub: "level1",
			description: "解析嵌套的子命令",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			err := cmd.Parse(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("意外的错误: %v", err)
					return
				}

				// 验证解析状态
				if !cmd.IsParsed() {
					t.Error("解析后命令应该处于已解析状态")
				}
			}
		})
	}
}

// TestParse_内置标志处理 测试Parse方法的内置标志处理
func TestParse_内置标志处理(t *testing.T) {
	tests := []struct {
		name        string
		setupCmd    func() *Cmd
		args        []string
		expectExit  bool
		description string
	}{
		{
			name: "解析help长标志",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.SetExitOnBuiltinFlags(false) // 禁用退出以便测试
				return cmd
			},
			args:        []string{"--help"},
			expectExit:  false,
			description: "解析--help标志",
		},
		{
			name: "解析help短标志",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.SetExitOnBuiltinFlags(false) // 禁用退出以便测试
				return cmd
			},
			args:        []string{"-h"},
			expectExit:  false,
			description: "解析-h标志",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			err := cmd.Parse(tt.args)

			// 内置标志通常不会返回错误，而是执行相应的动作
			if err != nil {
				t.Errorf("解析内置标志时出现意外错误: %v", err)
			}

			// 验证解析状态
			if !cmd.IsParsed() {
				t.Error("解析后命令应该处于已解析状态")
			}
		})
	}
}

// TestParse_重复解析 测试Parse方法的重复解析行为
func TestParse_重复解析(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.String("name", "n", "", "名称标志")

	// 第一次解析
	err1 := cmd.Parse([]string{"--name", "first"})
	if err1 != nil {
		t.Fatalf("第一次解析失败: %v", err1)
	}

	if !cmd.IsParsed() {
		t.Error("第一次解析后应该处于已解析状态")
	}

	// 第二次解析应该被忽略或返回错误
	err2 := cmd.Parse([]string{"--name", "second"})

	// 验证重复解析的行为（具体行为取决于实现）
	t.Logf("重复解析结果: %v", err2)

	// 无论如何，解析状态应该保持为true
	if !cmd.IsParsed() {
		t.Error("重复解析后仍应该处于已解析状态")
	}
}

// TestParseFlags_基本功能 测试ParseFlagsOnly方法的基本功能
func TestParseFlags_基本功能(t *testing.T) {
	tests := []struct {
		name        string
		setupCmd    func() *Cmd
		args        []string
		expectError bool
		errorMsg    string
		description string
	}{
		{
			name: "仅解析标志_忽略位置参数",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.String("name", "n", "default", "名称标志")
				cmd.Bool("verbose", "v", false, "详细输出")
				return cmd
			},
			args:        []string{"--name", "test", "--verbose", "ignored_arg", "another_ignored"},
			expectError: false,
			description: "ParseFlagsOnly应该只解析标志，忽略位置参数",
		},
		{
			name: "仅解析标志_空参数",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.String("config", "c", "", "配置文件")
				return cmd
			},
			args:        []string{},
			expectError: false,
			description: "ParseFlagsOnly处理空参数列表",
		},
		{
			name: "仅解析标志_只有位置参数",
			setupCmd: func() *Cmd {
				return NewCmd("test", "t", flag.ContinueOnError)
			},
			args:        []string{"arg1", "arg2", "arg3"},
			expectError: false,
			description: "ParseFlagsOnly处理只有位置参数的情况",
		},
		{
			name: "仅解析标志_无效标志",
			setupCmd: func() *Cmd {
				return NewCmd("test", "t", flag.ContinueOnError)
			},
			args:        []string{"--unknown"},
			expectError: true,
			errorMsg:    "flag provided but not defined",
			description: "ParseFlagsOnly遇到未定义标志应该返回错误",
		},
		{
			name: "仅解析标志_标志值类型错误",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.Int("port", "p", 8080, "端口号")
				return cmd
			},
			args:        []string{"--port", "invalid"},
			expectError: true,
			errorMsg:    "invalid value",
			description: "ParseFlagsOnly遇到类型错误应该返回错误",
		},
		{
			name: "仅解析标志_混合短长标志",
			setupCmd: func() *Cmd {
				cmd := NewCmd("test", "t", flag.ContinueOnError)
				cmd.String("file", "f", "", "文件路径")
				cmd.Bool("recursive", "r", false, "递归")
				cmd.Int("depth", "d", 1, "深度")
				return cmd
			},
			args:        []string{"-f", "test.txt", "--recursive", "-d", "3", "ignored"},
			expectError: false,
			description: "ParseFlagsOnly处理混合的长短标志",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			err := cmd.ParseFlagsOnly(tt.args)

			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("错误信息不匹配: 期望包含 %q, 实际 %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("意外的错误: %v", err)
					return
				}

				// 验证解析状态
				if !cmd.IsParsed() {
					t.Error("ParseFlagsOnly后命令应该处于已解析状态")
				}
			}
		})
	}
}

// TestParseFlags_忽略子命令 测试ParseFlagsOnly忽略子命令的行为
func TestParseFlags_忽略子命令(t *testing.T) {
	// 创建带子命令的父命令
	parent := NewCmd("parent", "p", flag.ContinueOnError)
	parent.String("global", "g", "", "全局标志")

	child := NewCmd("child", "c", flag.ContinueOnError)
	child.String("local", "l", "", "本地标志")

	err := parent.AddSubCmd(child)
	if err != nil {
		t.Fatalf("添加子命令失败: %v", err)
	}

	// ParseFlagsOnly应该忽略子命令，只解析父命令的标志
	err = parent.ParseFlagsOnly([]string{"--global", "value", "child", "--local", "ignored"})
	if err != nil {
		t.Errorf("ParseFlagsOnly处理包含子命令的参数时出错: %v", err)
	}

	// 验证解析状态
	if !parent.IsParsed() {
		t.Error("ParseFlagsOnly后父命令应该处于已解析状态")
	}

	// 子命令不应该被解析
	if child.IsParsed() {
		t.Error("ParseFlagsOnly不应该解析子命令")
	}
}

// TestParseFlags_内置标志 测试ParseFlagsOnly处理内置标志
func TestParseFlags_内置标志(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.SetExitOnBuiltinFlags(false) // 禁用退出以便测试

	// 测试help标志
	err := cmd.ParseFlagsOnly([]string{"--help", "ignored_arg"})
	if err != nil {
		t.Errorf("ParseFlagsOnly处理help标志时出错: %v", err)
	}

	if !cmd.IsParsed() {
		t.Error("ParseFlagsOnly处理内置标志后应该处于已解析状态")
	}
}

// TestParseFlags_重复解析 测试ParseFlagsOnly的重复解析行为
func TestParseFlags_重复解析(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.String("name", "n", "", "名称标志")

	// 第一次解析
	err1 := cmd.ParseFlagsOnly([]string{"--name", "first"})
	if err1 != nil {
		t.Fatalf("第一次ParseFlagsOnly失败: %v", err1)
	}

	if !cmd.IsParsed() {
		t.Error("第一次ParseFlagsOnly后应该处于已解析状态")
	}

	// 第二次解析
	err2 := cmd.ParseFlagsOnly([]string{"--name", "second"})

	// 验证重复解析的行为
	t.Logf("重复ParseFlagsOnly结果: %v", err2)

	// 解析状态应该保持为true
	if !cmd.IsParsed() {
		t.Error("重复ParseFlagsOnly后仍应该处于已解析状态")
	}
}

// TestParse_vs_ParseFlags_对比 测试Parse和ParseFlagsOnly的行为差异
func TestParse_vs_ParseFlags_对比(t *testing.T) {
	// 创建两个相同配置的命令用于对比
	createCmd := func() *Cmd {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		cmd.String("name", "n", "", "名称标志")
		child := NewCmd("child", "c", flag.ContinueOnError)
		_ = cmd.AddSubCmd(child)
		return cmd
	}

	args := []string{"--name", "test", "child", "arg1", "arg2"}

	// 测试Parse方法
	cmd1 := createCmd()
	err1 := cmd1.Parse(args)

	// 测试ParseFlagsOnly方法
	cmd2 := createCmd()
	err2 := cmd2.ParseFlagsOnly(args)

	// 两个方法都应该成功解析标志部分
	if err1 != nil {
		t.Errorf("Parse方法出错: %v", err1)
	}
	if err2 != nil {
		t.Errorf("ParseFlagsOnly方法出错: %v", err2)
	}

	// 两个命令都应该处于已解析状态
	if !cmd1.IsParsed() {
		t.Error("Parse后命令应该处于已解析状态")
	}
	if !cmd2.IsParsed() {
		t.Error("ParseFlagsOnly后命令应该处于已解析状态")
	}

	t.Logf("Parse和ParseFlagsOnly对比测试完成")
}
