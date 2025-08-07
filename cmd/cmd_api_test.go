// Package cmd 命令API测试
// 本文件包含了Cmd结构体API接口的单元测试，测试面向对象API
// 与内部函数式API的适配功能，确保API设计的正确性。
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// TestNewCmd_边界场景 测试NewCmd函数的边界场景
func TestNewCmd_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		longName    string
		shortName   string
		errorMode   flag.ErrorHandling
		expectPanic bool
		description string
	}{
		{
			name:        "正常创建_长短名称都有",
			longName:    "test",
			shortName:   "t",
			errorMode:   flag.ContinueOnError,
			expectPanic: false,
			description: "正常情况下创建命令",
		},
		{
			name:        "只有长名称",
			longName:    "test-long",
			shortName:   "",
			errorMode:   flag.ContinueOnError,
			expectPanic: false,
			description: "只提供长名称",
		},
		{
			name:        "只有短名称",
			longName:    "",
			shortName:   "t",
			errorMode:   flag.ContinueOnError,
			expectPanic: false,
			description: "只提供短名称",
		},
		{
			name:        "长短名称都为空",
			longName:    "",
			shortName:   "",
			errorMode:   flag.ContinueOnError,
			expectPanic: true,
			description: "长短名称都为空字符串应该panic",
		},
		{
			name:        "特殊字符名称",
			longName:    "test-cmd_123",
			shortName:   "t1",
			errorMode:   flag.ContinueOnError,
			expectPanic: false,
			description: "包含特殊字符的名称",
		},
		{
			name:        "中文名称",
			longName:    "测试命令",
			shortName:   "测",
			errorMode:   flag.ContinueOnError,
			expectPanic: false,
			description: "中文命令名称",
		},
		{
			name:        "ExitOnError模式",
			longName:    "test",
			shortName:   "t",
			errorMode:   flag.ExitOnError,
			expectPanic: false,
			description: "退出错误处理模式",
		},
		{
			name:        "PanicOnError模式",
			longName:    "test",
			shortName:   "t",
			errorMode:   flag.PanicOnError,
			expectPanic: false,
			description: "恐慌错误处理模式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd *Cmd
			var panicked bool

			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
						if !tt.expectPanic {
							t.Errorf("意外的panic: %v", r)
						}
					}
				}()
				cmd = NewCmd(tt.longName, tt.shortName, tt.errorMode)
			}()

			if tt.expectPanic {
				if !panicked {
					t.Error("期望panic但没有发生")
				}
				return // panic情况下不继续验证
			}

			if cmd == nil {
				t.Fatal("NewCmd返回了nil")
			}

			if cmd.ctx == nil {
				t.Fatal("命令上下文为nil")
			}

			if cmd.LongName() != tt.longName {
				t.Errorf("长名称不匹配: 期望 %q, 实际 %q", tt.longName, cmd.LongName())
			}

			if cmd.ShortName() != tt.shortName {
				t.Errorf("短名称不匹配: 期望 %q, 实际 %q", tt.shortName, cmd.ShortName())
			}

			// 验证内置help标志是否正确注册
			if !cmd.FlagExists(flags.HelpFlagName) {
				t.Error("内置help标志未正确注册")
			}
		})
	}
}

// TestNew_别名函数 测试New别名函数
func TestNew_别名函数(t *testing.T) {
	cmd1 := New("test", "t", flag.ContinueOnError)
	cmd2 := NewCmd("test", "t", flag.ContinueOnError)

	if cmd1.LongName() != cmd2.LongName() || cmd1.ShortName() != cmd2.ShortName() {
		t.Error("New别名函数与NewCmd行为不一致")
	}
}

// TestAddSubCmd_边界场景 测试AddSubCmd的边界场景
func TestAddSubCmd_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		setupParent func() *Cmd
		setupSubs   func() []*Cmd
		expectError bool
		errorMsg    string
		description string
	}{
		{
			name: "空子命令列表",
			setupParent: func() *Cmd {
				return NewCmd("parent", "p", flag.ContinueOnError)
			},
			setupSubs: func() []*Cmd {
				return []*Cmd{}
			},
			expectError: true,
			errorMsg:    "subCmds list cannot be empty",
			description: "传入空的子命令列表",
		},
		{
			name: "nil子命令",
			setupParent: func() *Cmd {
				return NewCmd("parent", "p", flag.ContinueOnError)
			},
			setupSubs: func() []*Cmd {
				return []*Cmd{nil}
			},
			expectError: true,
			errorMsg:    "subCmd at index 0 cannot be nil",
			description: "传入nil子命令",
		},
		{
			name: "混合nil和正常子命令",
			setupParent: func() *Cmd {
				return NewCmd("parent", "p", flag.ContinueOnError)
			},
			setupSubs: func() []*Cmd {
				return []*Cmd{
					NewCmd("child1", "c1", flag.ContinueOnError),
					nil,
					NewCmd("child2", "c2", flag.ContinueOnError),
				}
			},
			expectError: true,
			errorMsg:    "subCmd at index 1 cannot be nil",
			description: "混合nil和正常子命令",
		},
		{
			name: "重复长名称",
			setupParent: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				child1 := NewCmd("child", "c1", flag.ContinueOnError)
				_ = parent.AddSubCmd(child1) // 先添加一个
				return parent
			},
			setupSubs: func() []*Cmd {
				child2 := NewCmd("child", "c2", flag.ContinueOnError) // 重复长名称
				return []*Cmd{child2}
			},
			expectError: true,
			errorMsg:    "long name 'child' already exists",
			description: "添加重复长名称的子命令",
		},
		{
			name: "重复短名称",
			setupParent: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				child1 := NewCmd("child1", "c", flag.ContinueOnError)
				_ = parent.AddSubCmd(child1) // 先添加一个
				return parent
			},
			setupSubs: func() []*Cmd {
				child2 := NewCmd("child2", "c", flag.ContinueOnError) // 重复短名称
				return []*Cmd{child2}
			},
			expectError: true,
			errorMsg:    "short name 'c' already exists",
			description: "添加重复短名称的子命令",
		},
		{
			name: "正常添加已有父命令的子命令",
			setupParent: func() *Cmd {
				parent := NewCmd("parent", "p", flag.ContinueOnError)
				return parent
			},
			setupSubs: func() []*Cmd {
				// 测试添加一个已经有父命令的子命令（这应该是允许的，会重新设置父命令）
				child := NewCmd("child", "c", flag.ContinueOnError)
				grandparent := NewCmd("grandparent", "gp", flag.ContinueOnError)
				_ = grandparent.AddSubCmd(child) // child现在有了父命令
				return []*Cmd{child}             // 将child添加到新的父命令
			},
			expectError: false,
			errorMsg:    "",
			description: "测试添加已有父命令的子命令（应该允许重新设置父命令）",
		},
		{
			name: "批量添加正常子命令",
			setupParent: func() *Cmd {
				return NewCmd("parent", "p", flag.ContinueOnError)
			},
			setupSubs: func() []*Cmd {
				return []*Cmd{
					NewCmd("child1", "c1", flag.ContinueOnError),
					NewCmd("child2", "c2", flag.ContinueOnError),
					NewCmd("child3", "", flag.ContinueOnError),
					NewCmd("", "c4", flag.ContinueOnError),
				}
			},
			expectError: false,
			description: "批量添加多个正常子命令",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := tt.setupParent()
			subs := tt.setupSubs()

			err := parent.AddSubCmd(subs...)

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

				// 验证子命令是否正确添加
				subCmdMap := parent.SubCmdMap()
				subCmds := parent.SubCmds()

				if len(subCmds) != len(subs) {
					t.Errorf("子命令数量不匹配: 期望 %d, 实际 %d", len(subs), len(subCmds))
				}

				// 验证每个子命令都能通过名称找到
				for _, sub := range subs {
					if sub.LongName() != "" {
						if _, exists := subCmdMap[sub.LongName()]; !exists {
							t.Errorf("长名称 %q 的子命令未找到", sub.LongName())
						}
					}
					if sub.ShortName() != "" {
						if _, exists := subCmdMap[sub.ShortName()]; !exists {
							t.Errorf("短名称 %q 的子命令未找到", sub.ShortName())
						}
					}
				}
			}
		})
	}
}

// TestSubCmdMap_边界场景 测试SubCmdMap的边界场景
func TestSubCmdMap_边界场景(t *testing.T) {
	// 测试空子命令映射
	t.Run("空子命令映射", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		subCmdMap := cmd.SubCmdMap()

		if subCmdMap != nil {
			t.Error("SubCmdMap应该返回nil, 实际不为nil")
		}
	})

	// 测试返回副本而非原始引用
	t.Run("返回副本测试", func(t *testing.T) {
		parent := NewCmd("parent", "p", flag.ContinueOnError)
		child := NewCmd("child", "c", flag.ContinueOnError)

		err := parent.AddSubCmd(child)
		if err != nil {
			t.Fatalf("添加子命令失败: %v", err)
		}

		subCmdMap1 := parent.SubCmdMap()
		subCmdMap2 := parent.SubCmdMap()

		// 修改第一个映射
		delete(subCmdMap1, "child")

		// 验证第二个映射未受影响
		if _, exists := subCmdMap2["child"]; !exists {
			t.Error("SubCmdMap返回的不是副本，外部修改影响了内部状态")
		}
	})
}

// TestSubCmds_边界场景 测试SubCmds的边界场景
func TestSubCmds_边界场景(t *testing.T) {
	// 测试空子命令切片
	t.Run("空子命令切片", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		subCmds := cmd.SubCmds()

		if subCmds != nil {
			t.Errorf("空命令的子命令切片应为空, 实际长度: %d", len(subCmds))
		}
	})

	// 测试返回副本而非原始引用
	t.Run("返回副本测试", func(t *testing.T) {
		parent := NewCmd("parent", "p", flag.ContinueOnError)
		child1 := NewCmd("child1", "c1", flag.ContinueOnError)
		child2 := NewCmd("child2", "c2", flag.ContinueOnError)

		err := parent.AddSubCmd(child1, child2)
		if err != nil {
			t.Fatalf("添加子命令失败: %v", err)
		}

		subCmds1 := parent.SubCmds()
		subCmds2 := parent.SubCmds()

		// 修改第一个切片
		if len(subCmds1) > 0 {
			subCmds1[0] = nil
		}

		// 验证第二个切片未受影响
		if len(subCmds2) == 0 || subCmds2[0] == nil {
			t.Error("SubCmds返回的不是副本，外部修改影响了内部状态")
		}
	})
}

// TestSetEnableCompletion_边界场景 测试SetEnableCompletion的边界场景
func TestSetEnableCompletion_边界场景(t *testing.T) {
	// 测试根命令启用补全
	t.Run("根命令启用补全", func(t *testing.T) {
		cmd := NewCmd("root", "r", flag.ContinueOnError)

		// 启用补全
		cmd.SetEnableCompletion(true)

		// 由于没有公开的getter方法，我们通过内部状态验证
		if !cmd.ctx.Config.EnableCompletion {
			t.Error("根命令启用补全失败")
		}

		// 禁用补全
		cmd.SetEnableCompletion(false)
		if cmd.ctx.Config.EnableCompletion {
			t.Error("根命令禁用补全失败")
		}
	})

	// 测试子命令不能启用补全
	t.Run("子命令不能启用补全", func(t *testing.T) {
		parent := NewCmd("parent", "p", flag.ContinueOnError)
		child := NewCmd("child", "c", flag.ContinueOnError)

		err := parent.AddSubCmd(child)
		if err != nil {
			t.Fatalf("添加子命令失败: %v", err)
		}

		// 尝试在子命令上启用补全
		child.SetEnableCompletion(true)

		// 验证子命令的补全状态未改变
		if child.ctx.Config.EnableCompletion {
			t.Error("子命令不应该能够启用补全")
		}
	})
}

// TestVersionMethods_边界场景 测试版本相关方法的边界场景
func TestVersionMethods_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		description string
	}{
		{
			name:        "正常版本号",
			version:     "1.0.0",
			description: "标准语义化版本号",
		},
		{
			name:        "空版本号",
			version:     "",
			description: "空字符串版本号",
		},
		{
			name:        "复杂版本号",
			version:     "v2.1.3-beta.1+build.123",
			description: "包含预发布和构建信息的版本号",
		},
		{
			name:        "中文版本信息",
			version:     "版本 1.0.0",
			description: "包含中文的版本信息",
		},
		{
			name:        "特殊字符版本",
			version:     "1.0.0-alpha+build_2023.01.01",
			description: "包含特殊字符的版本号",
		},
		{
			name:        "长版本字符串",
			version:     strings.Repeat("1.0.0-", 100) + "final",
			description: "非常长的版本字符串",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)

			// 设置版本
			cmd.SetVersion(tt.version)

			// 获取版本并验证
			gotVersion := cmd.GetVersion()
			if gotVersion != tt.version {
				t.Errorf("版本不匹配: 期望 %q, 实际 %q", tt.version, gotVersion)
			}
		})
	}
}

// TestModuleHelps_边界场景 测试模块帮助相关方法的边界场景
func TestModuleHelps_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		moduleHelps string
		description string
	}{
		{
			name:        "正常模块帮助",
			moduleHelps: "这是模块帮助信息",
			description: "正常的模块帮助文本",
		},
		{
			name:        "空模块帮助",
			moduleHelps: "",
			description: "空的模块帮助",
		},
		{
			name:        "多行模块帮助",
			moduleHelps: "第一行帮助\n第二行帮助\n第三行帮助",
			description: "多行模块帮助信息",
		},
		{
			name:        "包含特殊字符的帮助",
			moduleHelps: "模块帮助: @#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			description: "包含各种特殊字符",
		},
		{
			name:        "长文本帮助",
			moduleHelps: strings.Repeat("这是一个很长的模块帮助信息。", 100),
			description: "非常长的模块帮助文本",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)

			// 设置模块帮助
			cmd.SetModuleHelps(tt.moduleHelps)

			// 获取模块帮助并验证
			gotHelps := cmd.GetModuleHelps()
			if gotHelps != tt.moduleHelps {
				t.Errorf("模块帮助不匹配: 期望 %q, 实际 %q", tt.moduleHelps, gotHelps)
			}
		})
	}
}

// TestLogoText_边界场景 测试Logo文本相关方法的边界场景
func TestLogoText_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		logoText    string
		description string
	}{
		{
			name:        "ASCII艺术Logo",
			logoText:    "  ___  \n /   \\ \n|  o  |\n \\___/ ",
			description: "ASCII艺术风格的Logo",
		},
		{
			name:        "空Logo",
			logoText:    "",
			description: "空的Logo文本",
		},
		{
			name:        "单行Logo",
			logoText:    "MyApp v1.0",
			description: "简单的单行Logo",
		},
		{
			name:        "包含Unicode的Logo",
			logoText:    "🚀 MyApp 🚀\n✨ 版本 1.0 ✨",
			description: "包含Unicode字符的Logo",
		},
		{
			name:        "大型Logo",
			logoText:    strings.Repeat("█", 50) + "\n" + strings.Repeat("█", 50),
			description: "大型Logo文本",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)

			// 设置Logo文本
			cmd.SetLogoText(tt.logoText)

			// 获取Logo文本并验证
			gotLogo := cmd.GetLogoText()
			if gotLogo != tt.logoText {
				t.Errorf("Logo文本不匹配: 期望 %q, 实际 %q", tt.logoText, gotLogo)
			}
		})
	}
}

// TestUseChinese_边界场景 测试中文设置相关方法的边界场景
func TestUseChinese_边界场景(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	// 测试默认值
	defaultUseChinese := cmd.GetUseChinese()
	t.Logf("默认中文设置: %v", defaultUseChinese)

	// 测试设置为true
	cmd.SetUseChinese(true)
	if !cmd.GetUseChinese() {
		t.Error("设置中文为true失败")
	}

	// 测试设置为false
	cmd.SetUseChinese(false)
	if cmd.GetUseChinese() {
		t.Error("设置中文为false失败")
	}

	// 测试多次切换
	for i := 0; i < 10; i++ {
		expected := i%2 == 0
		cmd.SetUseChinese(expected)
		if cmd.GetUseChinese() != expected {
			t.Errorf("第%d次切换失败: 期望 %v, 实际 %v", i, expected, cmd.GetUseChinese())
		}
	}
}

// TestNotes_边界场景 测试备注相关方法的边界场景
func TestNotes_边界场景(t *testing.T) {
	// 测试空备注列表
	t.Run("空备注列表", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		notes := cmd.GetNotes()

		if notes == nil {
			t.Error("GetNotes返回了nil")
		}

		if len(notes) != 0 {
			t.Errorf("新命令应该没有备注, 实际数量: %d", len(notes))
		}
	})

	// 测试添加各种类型的备注
	t.Run("添加各种备注", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		testNotes := []string{
			"正常备注",
			"",
			"包含\n换行符的备注",
			"包含特殊字符的备注: @#$%^&*()",
			"很长很长很长很长很长很长很长很长很长很长的备注信息",
			"中文备注：这是一个中文备注",
			"Unicode备注: 🎉🚀✨",
		}

		// 添加所有备注
		for _, note := range testNotes {
			cmd.AddNote(note)
		}

		// 获取备注并验证
		gotNotes := cmd.GetNotes()
		if len(gotNotes) != len(testNotes) {
			t.Errorf("备注数量不匹配: 期望 %d, 实际 %d", len(testNotes), len(gotNotes))
		}

		for i, expectedNote := range testNotes {
			if i >= len(gotNotes) {
				t.Errorf("缺少第%d个备注", i)
				continue
			}
			if gotNotes[i] != expectedNote {
				t.Errorf("第%d个备注不匹配: 期望 %q, 实际 %q", i, expectedNote, gotNotes[i])
			}
		}
	})

	// 测试返回副本而非原始引用
	t.Run("返回副本测试", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		cmd.AddNote("原始备注")

		notes1 := cmd.GetNotes()
		notes2 := cmd.GetNotes()

		// 修改第一个切片
		if len(notes1) > 0 {
			notes1[0] = "修改后的备注"
		}

		// 验证第二个切片未受影响
		if len(notes2) > 0 && notes2[0] != "原始备注" {
			t.Error("GetNotes返回的不是副本，外部修改影响了内部状态")
		}
	})
}

// TestName_边界场景 测试Name方法的边界场景
func TestName_边界场景(t *testing.T) {
	tests := []struct {
		name         string
		longName     string
		shortName    string
		expectedName string
		expectPanic  bool
		description  string
	}{
		{
			name:         "长短名称都有",
			longName:     "test-long",
			shortName:    "t",
			expectedName: "test-long",
			expectPanic:  false,
			description:  "优先返回长名称",
		},
		{
			name:         "只有长名称",
			longName:     "test-long",
			shortName:    "",
			expectedName: "test-long",
			expectPanic:  false,
			description:  "只有长名称时返回长名称",
		},
		{
			name:         "只有短名称",
			longName:     "",
			shortName:    "t",
			expectedName: "t",
			expectPanic:  false,
			description:  "只有短名称时返回短名称",
		},
		{
			name:        "长短名称都为空",
			longName:    "",
			shortName:   "",
			expectPanic: true,
			description: "长短名称都为空时应该panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd *Cmd
			var panicked bool

			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
						if !tt.expectPanic {
							t.Errorf("意外的panic: %v", r)
						}
					}
				}()
				cmd = NewCmd(tt.longName, tt.shortName, flag.ContinueOnError)
			}()

			if tt.expectPanic {
				if !panicked {
					t.Error("期望panic但没有发生")
				}
				return // panic情况下不继续验证
			}

			if cmd == nil {
				t.Fatal("NewCmd返回了nil")
			}

			gotName := cmd.Name()
			if gotName != tt.expectedName {
				t.Errorf("Name()返回值不匹配: 期望 %q, 实际 %q", tt.expectedName, gotName)
			}
		})
	}
}

// TestDescription_边界场景 测试描述相关方法的边界场景
func TestDescription_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		description string
		testDesc    string
	}{
		{
			name:        "正常描述",
			description: "这是一个测试命令",
			testDesc:    "正常的命令描述",
		},
		{
			name:        "空描述",
			description: "",
			testDesc:    "空的命令描述",
		},
		{
			name:        "多行描述",
			description: "第一行描述\n第二行描述\n第三行描述",
			testDesc:    "多行命令描述",
		},
		{
			name:        "包含特殊字符的描述",
			description: "描述包含特殊字符: @#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			testDesc:    "包含各种特殊字符的描述",
		},
		{
			name:        "长描述",
			description: strings.Repeat("这是一个很长的命令描述。", 50),
			testDesc:    "非常长的命令描述",
		},
		{
			name:        "Unicode描述",
			description: "命令描述包含Unicode: 🎉🚀✨ 中文描述",
			testDesc:    "包含Unicode字符的描述",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)

			// 设置描述
			cmd.SetDescription(tt.description)

			// 获取描述并验证
			gotDesc := cmd.GetDescription()
			if gotDesc != tt.description {
				t.Errorf("描述不匹配: 期望 %q, 实际 %q", tt.description, gotDesc)
			}
		})
	}
}

// TestHelp_边界场景 测试帮助相关方法的边界场景
func TestHelp_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		customHelp  string
		description string
	}{
		{
			name:        "自定义帮助信息",
			customHelp:  "这是自定义的帮助信息",
			description: "设置自定义帮助信息",
		},
		{
			name:        "空帮助信息",
			customHelp:  "",
			description: "空的自定义帮助信息",
		},
		{
			name:        "多行帮助信息",
			customHelp:  "第一行帮助\n第二行帮助\n第三行帮助",
			description: "多行自定义帮助信息",
		},
		{
			name:        "包含格式化的帮助",
			customHelp:  "用法: myapp [选项]\n\n选项:\n  -h, --help    显示帮助信息",
			description: "包含格式化内容的帮助信息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)

			// 设置自定义帮助
			cmd.SetHelp(tt.customHelp)

			// 获取帮助信息
			gotHelp := cmd.GetHelp()

			// 如果设置了自定义帮助，应该返回自定义内容
			if tt.customHelp != "" {
				if !strings.Contains(gotHelp, tt.customHelp) {
					t.Errorf("帮助信息应包含自定义内容: %q", tt.customHelp)
				}
			}
		})
	}
}

// TestLoadHelp_边界场景 测试LoadHelp方法的边界场景
func TestLoadHelp_边界场景(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		setupFile   func() string
		expectError bool
		errorMsg    string
		description string
	}{
		{
			name: "正常加载帮助文件",
			setupFile: func() string {
				filePath := filepath.Join(tmpDir, "help.txt")
				content := "这是从文件加载的帮助信息"
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("创建测试文件失败: %v", err)
				}
				return filePath
			},
			expectError: false,
			description: "正常加载存在的帮助文件",
		},
		{
			name: "空文件路径",
			setupFile: func() string {
				return ""
			},
			expectError: true,
			errorMsg:    "file path cannot be empty",
			description: "传入空的文件路径",
		},
		{
			name: "只包含空白字符的路径",
			setupFile: func() string {
				return "   \t\n   "
			},
			expectError: true,
			errorMsg:    "file path cannot be empty or contain only whitespace",
			description: "传入只包含空白字符的路径",
		},
		{
			name: "不存在的文件",
			setupFile: func() string {
				return filepath.Join(tmpDir, "nonexistent.txt")
			},
			expectError: true,
			errorMsg:    "does not exist",
			description: "尝试加载不存在的文件",
		},
		{
			name: "空文件",
			setupFile: func() string {
				filePath := filepath.Join(tmpDir, "empty.txt")
				err := os.WriteFile(filePath, []byte(""), 0644)
				if err != nil {
					t.Fatalf("创建空测试文件失败: %v", err)
				}
				return filePath
			},
			expectError: false,
			description: "加载空的帮助文件",
		},
		{
			name: "大文件",
			setupFile: func() string {
				filePath := filepath.Join(tmpDir, "large.txt")
				content := strings.Repeat("这是一行很长的帮助信息。\n", 1000)
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("创建大测试文件失败: %v", err)
				}
				return filePath
			},
			expectError: false,
			description: "加载大的帮助文件",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)
			filePath := tt.setupFile()

			err := cmd.LoadHelp(filePath)

			if tt.expectError {
				if err == nil {
					t.Error("期望错误但没有返回错误")
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

				// 验证帮助内容是否正确加载
				if filePath != "" {
					expectedContent, readErr := os.ReadFile(filePath)
					if readErr != nil {
						t.Fatalf("读取测试文件失败: %v", readErr)
					}

					gotHelp := cmd.GetHelp()
					if !strings.Contains(gotHelp, string(expectedContent)) {
						t.Error("加载的帮助内容不正确")
					}
				}
			}
		})
	}
}

// TestUsageSyntax_边界场景 测试用法语法相关方法的边界场景
func TestUsageSyntax_边界场景(t *testing.T) {
	tests := []struct {
		name        string
		usageSyntax string
		description string
	}{
		{
			name:        "正常用法语法",
			usageSyntax: "myapp [选项] <文件>",
			description: "正常的用法语法",
		},
		{
			name:        "空用法语法",
			usageSyntax: "",
			description: "空的用法语法",
		},
		{
			name:        "复杂用法语法",
			usageSyntax: "myapp [全局选项] <命令> [命令选项] [参数...]",
			description: "复杂的用法语法",
		},
		{
			name:        "包含特殊字符的用法",
			usageSyntax: "myapp [-h|--help] [-v|--version] <file1> [file2...]",
			description: "包含特殊字符的用法语法",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)

			// 设置用法语法
			cmd.SetUsageSyntax(tt.usageSyntax)

			// 获取用法语法并验证
			gotUsage := cmd.GetUsageSyntax()
			if gotUsage != tt.usageSyntax {
				t.Errorf("用法语法不匹配: 期望 %q, 实际 %q", tt.usageSyntax, gotUsage)
			}
		})
	}
}

// TestExamples_边界场景 测试示例相关方法的边界场景
func TestExamples_边界场景(t *testing.T) {
	// 测试空示例列表
	t.Run("空示例列表", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		examples := cmd.GetExamples()

		if examples == nil {
			t.Error("GetExamples返回了nil")
		}

		if len(examples) != 0 {
			t.Errorf("新命令应该没有示例, 实际数量: %d", len(examples))
		}
	})

	// 测试添加各种示例
	t.Run("添加各种示例", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		testExamples := []struct {
			desc  string
			usage string
		}{
			{"基本用法", "myapp file.txt"},
			{"", "myapp --help"},
			{"复杂用法", "myapp --config /path/to/config.json --verbose file1.txt file2.txt"},
			{"包含特殊字符", "myapp 'file with spaces.txt'"},
			{"多行用法", "myapp \\\n  --option1 value1 \\\n  --option2 value2"},
		}

		// 添加所有示例
		for _, example := range testExamples {
			cmd.AddExample(example.desc, example.usage)
		}

		// 获取示例并验证
		gotExamples := cmd.GetExamples()
		if len(gotExamples) != len(testExamples) {
			t.Errorf("示例数量不匹配: 期望 %d, 实际 %d", len(testExamples), len(gotExamples))
		}

		for i, expectedExample := range testExamples {
			if i >= len(gotExamples) {
				t.Errorf("缺少第%d个示例", i)
				continue
			}
			if gotExamples[i].Description != expectedExample.desc {
				t.Errorf("第%d个示例描述不匹配: 期望 %q, 实际 %q", i, expectedExample.desc, gotExamples[i].Description)
			}
			if gotExamples[i].Usage != expectedExample.usage {
				t.Errorf("第%d个示例用法不匹配: 期望 %q, 实际 %q", i, expectedExample.usage, gotExamples[i].Usage)
			}
		}
	})
}

// TestArgs_边界场景 测试参数相关方法的边界场景
func TestArgs_边界场景(t *testing.T) {
	// 测试空参数
	t.Run("空参数", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		// 解析空参数
		err := cmd.Parse([]string{})
		if err != nil {
			t.Fatalf("解析空参数失败: %v", err)
		}

		// 验证参数相关方法
		if cmd.NArg() != 0 {
			t.Errorf("空参数的数量应为0, 实际: %d", cmd.NArg())
		}

		args := cmd.Args()
		if len(args) != 0 {
			t.Errorf("空参数列表长度应为0, 实际: %d", len(args))
		}

		if cmd.Arg(0) != "" {
			t.Errorf("索引0的参数应为空字符串, 实际: %q", cmd.Arg(0))
		}

		if cmd.Arg(-1) != "" {
			t.Errorf("负索引的参数应为空字符串, 实际: %q", cmd.Arg(-1))
		}
	})

	// 测试多个参数
	t.Run("多个参数", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		testArgs := []string{"arg1", "arg2", "arg with spaces", "", "arg5"}

		err := cmd.Parse(testArgs)
		if err != nil {
			t.Fatalf("解析参数失败: %v", err)
		}

		// 验证参数数量
		if cmd.NArg() != len(testArgs) {
			t.Errorf("参数数量不匹配: 期望 %d, 实际 %d", len(testArgs), cmd.NArg())
		}

		// 验证参数列表
		gotArgs := cmd.Args()
		if !reflect.DeepEqual(gotArgs, testArgs) {
			t.Errorf("参数列表不匹配: 期望 %v, 实际 %v", testArgs, gotArgs)
		}

		// 验证单个参数访问
		for i, expectedArg := range testArgs {
			gotArg := cmd.Arg(i)
			if gotArg != expectedArg {
				t.Errorf("第%d个参数不匹配: 期望 %q, 实际 %q", i, expectedArg, gotArg)
			}
		}

		// 验证越界访问
		if cmd.Arg(len(testArgs)) != "" {
			t.Error("越界访问应返回空字符串")
		}

		if cmd.Arg(-1) != "" {
			t.Error("负索引访问应返回空字符串")
		}
	})

	// 测试返回副本而非原始引用
	t.Run("返回副本测试", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		testArgs := []string{"arg1", "arg2"}

		err := cmd.Parse(testArgs)
		if err != nil {
			t.Fatalf("解析参数失败: %v", err)
		}

		args1 := cmd.Args()
		args2 := cmd.Args()

		// 修改第一个切片
		if len(args1) > 0 {
			args1[0] = "modified"
		}

		// 验证第二个切片未受影响
		if len(args2) > 0 && args2[0] != "arg1" {
			t.Error("Args返回的不是副本，外部修改影响了内部状态")
		}
	})
}

// TestFlagMethods_边界场景 测试标志相关方法的边界场景
func TestFlagMethods_边界场景(t *testing.T) {
	// 测试NFlag方法
	t.Run("NFlag测试", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		// 初始状态应该有内置的help标志
		initialCount := cmd.NFlag()

		// 添加一些标志
		cmd.String("str", "s", "default", "string flag")
		cmd.Int("int", "i", 0, "int flag")
		cmd.Bool("bool", "b", false, "bool flag")

		// 解析参数以激活标志
		err := cmd.Parse([]string{"--str", "value", "--int", "123", "--bool"})
		if err != nil {
			t.Fatalf("解析参数失败: %v", err)
		}

		// 验证标志数量（应该包括被设置的标志）
		finalCount := cmd.NFlag()
		if finalCount <= initialCount {
			t.Errorf("标志数量应该增加: 初始 %d, 最终 %d", initialCount, finalCount)
		}
	})

	// 测试FlagExists方法
	t.Run("FlagExists测试", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		// 添加一些标志
		cmd.String("string-flag", "s", "default", "string flag")
		cmd.Int("int-flag", "", 0, "int flag without short name")
		cmd.Bool("", "b", false, "bool flag without long name")

		// 测试存在的标志
		if !cmd.FlagExists("string-flag") {
			t.Error("应该找到string-flag")
		}

		if !cmd.FlagExists("s") {
			t.Error("应该找到短标志s")
		}

		if !cmd.FlagExists("int-flag") {
			t.Error("应该找到int-flag")
		}

		if !cmd.FlagExists("b") {
			t.Error("应该找到短标志b")
		}

		// 测试内置help标志
		if !cmd.FlagExists(flags.HelpFlagName) {
			t.Error("应该找到内置help标志")
		}

		if flags.HelpFlagShortName != "" && !cmd.FlagExists(flags.HelpFlagShortName) {
			t.Error("应该找到内置help短标志")
		}

		// 测试不存在的标志
		if cmd.FlagExists("nonexistent") {
			t.Error("不应该找到不存在的标志")
		}

		if cmd.FlagExists("") {
			t.Error("不应该找到空名称的标志")
		}
	})
}

// TestCmdExists_边界场景 测试CmdExists方法的边界场景
func TestCmdExists_边界场景(t *testing.T) {
	parent := NewCmd("parent", "p", flag.ContinueOnError)

	// 测试空子命令列表
	t.Run("空子命令列表", func(t *testing.T) {
		if parent.CmdExists("nonexistent") {
			t.Error("空子命令列表不应该找到任何命令")
		}

		if parent.CmdExists("") {
			t.Error("不应该找到空名称的命令")
		}
	})

	// 添加一些子命令
	child1 := NewCmd("child1", "c1", flag.ContinueOnError)
	child2 := NewCmd("child2", "", flag.ContinueOnError)
	child3 := NewCmd("", "c3", flag.ContinueOnError)

	err := parent.AddSubCmd(child1, child2, child3)
	if err != nil {
		t.Fatalf("添加子命令失败: %v", err)
	}

	// 测试存在的子命令
	t.Run("存在的子命令", func(t *testing.T) {
		if !parent.CmdExists("child1") {
			t.Error("应该找到child1")
		}

		if !parent.CmdExists("c1") {
			t.Error("应该找到短名称c1")
		}

		if !parent.CmdExists("child2") {
			t.Error("应该找到child2")
		}

		if !parent.CmdExists("c3") {
			t.Error("应该找到短名称c3")
		}
	})

	// 测试不存在的子命令
	t.Run("不存在的子命令", func(t *testing.T) {
		if parent.CmdExists("nonexistent") {
			t.Error("不应该找到不存在的命令")
		}

		if parent.CmdExists("") {
			t.Error("不应该找到空名称的命令")
		}

		if parent.CmdExists("child") {
			t.Error("不应该找到部分匹配的命令")
		}
	})
}

// TestIsParsed_边界场景 测试IsParsed方法的边界场景
func TestIsParsed_边界场景(t *testing.T) {
	// 测试未解析状态
	t.Run("未解析状态", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		if cmd.IsParsed() {
			t.Error("新创建的命令不应该处于已解析状态")
		}
	})

	// 测试解析后状态
	t.Run("解析后状态", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		err := cmd.Parse([]string{})
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		if !cmd.IsParsed() {
			t.Error("解析后的命令应该处于已解析状态")
		}
	})

	// 测试ParseFlagsOnly后状态
	t.Run("ParseFlagsOnly后状态", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)

		err := cmd.ParseFlagsOnly([]string{})
		if err != nil {
			t.Fatalf("ParseFlagsOnly失败: %v", err)
		}

		if !cmd.IsParsed() {
			t.Error("ParseFlagsOnly后的命令应该处于已解析状态")
		}
	})
}

// TestSetExitOnBuiltinFlags_边界场景 测试SetExitOnBuiltinFlags方法的边界场景
func TestSetExitOnBuiltinFlags_边界场景(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	// 测试默认值
	defaultValue := cmd.ctx.Config.ExitOnBuiltinFlags
	t.Logf("默认ExitOnBuiltinFlags值: %v", defaultValue)

	// 测试设置为false
	cmd.SetExitOnBuiltinFlags(false)
	if cmd.ctx.Config.ExitOnBuiltinFlags {
		t.Error("设置ExitOnBuiltinFlags为false失败")
	}

	// 测试设置为true
	cmd.SetExitOnBuiltinFlags(true)
	if !cmd.ctx.Config.ExitOnBuiltinFlags {
		t.Error("设置ExitOnBuiltinFlags为true失败")
	}

	// 测试多次切换
	for i := 0; i < 10; i++ {
		expected := i%2 == 0
		cmd.SetExitOnBuiltinFlags(expected)
		if cmd.ctx.Config.ExitOnBuiltinFlags != expected {
			t.Errorf("第%d次切换失败: 期望 %v, 实际 %v", i, expected, cmd.ctx.Config.ExitOnBuiltinFlags)
		}
	}
}

// TestFlagRegistry_边界场景 测试FlagRegistry方法的边界场景
func TestFlagRegistry_边界场景(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	// 获取标志注册表
	registry := cmd.FlagRegistry()
	if registry == nil {
		t.Fatal("FlagRegistry返回了nil")
	}

	// 添加一些标志
	cmd.String("test-flag", "tf", "default", "test flag")

	// 再次获取注册表，应该包含新添加的标志
	registry2 := cmd.FlagRegistry()
	if registry2 == nil {
		t.Fatal("FlagRegistry返回了nil")
	}

	// 验证标志是否在注册表中
	if _, exists := registry2.GetByName("test-flag"); !exists {
		t.Error("新添加的标志应该在注册表中")
	}
}

// TestPrintHelp_边界场景 测试PrintHelp方法的边界场景
func TestPrintHelp_边界场景(t *testing.T) {
	// 重定向标准输出以捕获打印内容
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("创建管道失败: %v", err)
	}
	os.Stdout = w

	// 创建一个goroutine来读取输出
	var output strings.Builder
	done := make(chan bool)
	go func() {
		defer close(done)
		_, _ = io.Copy(&output, r)
	}()

	// 测试打印帮助
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.SetDescription("测试命令描述")
	cmd.SetHelp("自定义帮助信息")

	cmd.PrintHelp()

	// 恢复标准输出
	_ = w.Close()
	os.Stdout = oldStdout
	<-done
	_ = r.Close()

	// 验证输出内容
	outputStr := output.String()
	if !strings.Contains(outputStr, "自定义帮助信息") {
		t.Error("PrintHelp应该输出自定义帮助信息")
	}
}

// TestConcurrency_并发安全测试 测试方法的并发安全性
func TestConcurrency_并发安全测试(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)

	// 并发测试各种setter和getter方法
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// 测试版本设置的并发安全
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				version := fmt.Sprintf("v%d.%d.%d", id, j, time.Now().Nanosecond()%1000)
				cmd.SetVersion(version)
				_ = cmd.GetVersion()
			}
		}(i)
	}

	// 测试描述设置的并发安全
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				desc := fmt.Sprintf("描述_%d_%d", id, j)
				cmd.SetDescription(desc)
				_ = cmd.GetDescription()
			}
		}(i)
	}

	// 测试备注添加的并发安全
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				note := fmt.Sprintf("备注_%d_%d", id, j)
				cmd.AddNote(note)
				_ = cmd.GetNotes()
			}
		}(i)
	}

	// 测试示例添加的并发安全
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				desc := fmt.Sprintf("示例描述_%d_%d", id, j)
				usage := fmt.Sprintf("示例用法_%d_%d", id, j)
				cmd.AddExample(desc, usage)
				_ = cmd.GetExamples()
			}
		}(i)
	}

	wg.Wait()

	// 验证最终状态的一致性
	version := cmd.GetVersion()
	description := cmd.GetDescription()
	notes := cmd.GetNotes()
	examples := cmd.GetExamples()

	t.Logf("并发测试完成 - 版本: %s, 描述: %s, 备注数: %d, 示例数: %d",
		version, description, len(notes), len(examples))
}
