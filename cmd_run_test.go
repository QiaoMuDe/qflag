package qflag

import (
	"errors"
	"testing"
)

// TestCmd_Run 测试手动执行Run函数的基本功能
func TestCmd_Run(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		setupRun   func(cmd *Cmd)
		executeRun bool // 是否手动执行Run函数
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "成功手动执行Run函数",
			args: []string{"--value", "test"},
			setupRun: func(cmd *Cmd) {
				value := cmd.String("value", "v", "", "测试值")
				cmd.Run = func(c *Cmd) error {
					if value.Get() != "test" {
						t.Errorf("期望值 'test', 实际值 '%s'", value.Get())
					}
					return nil
				}
			},
			executeRun: true,
			wantErr:    false,
		},
		{
			name: "手动执行Run函数返回错误",
			args: []string{"--value", "test"},
			setupRun: func(cmd *Cmd) {
				cmd.String("value", "v", "", "测试值")
				cmd.Run = func(c *Cmd) error {
					return errors.New("运行错误")
				}
			},
			executeRun: true,
			wantErr:    true,
			wantErrMsg: "运行错误",
		},
		{
			name: "不执行Run函数（只解析）",
			args: []string{"--value", "test"},
			setupRun: func(cmd *Cmd) {
				cmd.String("value", "v", "", "测试值")
				cmd.Run = func(c *Cmd) error {
					t.Error("Run函数不应该被执行")
					return nil
				}
			},
			executeRun: false, // 不执行Run函数
			wantErr:    false,
		},
		{
			name: "解析错误时不应该手动执行Run函数",
			args: []string{"--unknown-flag"},
			setupRun: func(cmd *Cmd) {
				cmd.Run = func(c *Cmd) error {
					t.Error("Run函数不应该被执行")
					return nil
				}
			},
			executeRun: true,
			wantErr:    true, // 应该有解析错误
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmd("test", "t", ContinueOnError)

			// 设置Run函数
			if tt.setupRun != nil {
				tt.setupRun(cmd)
			}

			// 执行解析（现在只解析，不自动执行Run函数）
			err := cmd.Parse(tt.args)

			// 如果解析出错，直接返回
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			// 如果需要手动执行Run函数
			if tt.executeRun && cmd.Run != nil {
				runErr := cmd.Run(cmd)
				if (runErr != nil) != tt.wantErr {
					t.Errorf("Run() error = %v, wantErr %v", runErr, tt.wantErr)
					return
				}
				if tt.wantErr && tt.wantErrMsg != "" && runErr.Error() != tt.wantErrMsg {
					t.Errorf("Run() error message = %v, want %v", runErr.Error(), tt.wantErrMsg)
				}
			}
		})
	}
}

// TestCmd_RunWithSubCommand 测试手动执行子命令的Run函数
func TestCmd_RunWithSubCommand(t *testing.T) {
	// 创建根命令
	rootCmd := NewCmd("app", "a", ContinueOnError)

	// 创建子命令
	subCmd := NewCmd("sub", "s", ContinueOnError)
	value := subCmd.String("value", "v", "default", "测试值")

	// 设置子命令的Run函数
	executed := false
	subCmd.Run = func(c *Cmd) error {
		executed = true
		if value.Get() != "custom" {
			t.Errorf("期望值 'custom', 实际值 '%s'", value.Get())
		}
		return nil
	}

	// 添加子命令
	if err := rootCmd.AddSubCmd(subCmd); err != nil {
		t.Errorf("AddSubCmd() error = %v", err)
	}

	// 解析包含子命令的参数（只解析，不执行）
	args := []string{"sub", "--value", "custom"}
	err := rootCmd.Parse(args)

	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 手动执行子命令的Run函数
	if subCmd.Run != nil {
		if err := subCmd.Run(subCmd); err != nil {
			t.Errorf("子命令Run函数执行错误: %v", err)
		}
	}

	if !executed {
		t.Error("子命令的Run函数没有被执行")
	}
}

// TestCmd_RunErrorPropagation 测试手动执行Run函数时错误的传播
func TestCmd_RunErrorPropagation(t *testing.T) {
	cmd := NewCmd("test", "t", ContinueOnError)
	cmd.Bool("flag", "f", false, "测试标志")

	// 设置会返回错误的Run函数
	runError := errors.New("运行失败")
	cmd.Run = func(c *Cmd) error {
		return runError
	}

	// 解析参数（只解析，不执行）
	err := cmd.Parse([]string{"--flag"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
		return
	}

	// 手动执行Run函数，应该返回Run函数的错误
	if cmd.Run != nil {
		runErr := cmd.Run(cmd)
		if runErr != runError {
			t.Errorf("Run() error = %v, want %v", runErr, runError)
		}
	}
}
