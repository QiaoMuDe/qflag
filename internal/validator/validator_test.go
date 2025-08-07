package validator

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

// TestValidateSubCommand 测试 ValidateSubCommand 函数
func TestValidateSubCommand(t *testing.T) {
	tests := []struct {
		name        string
		setupParent func() *types.CmdContext
		setupChild  func() *types.CmdContext
		wantErr     bool
		errContains string
	}{
		{
			name: "子命令为nil",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return nil
			},
			wantErr:     true,
			errContains: "is nil",
		},
		{
			name: "父命令为nil，应该通过验证",
			setupParent: func() *types.CmdContext {
				return nil
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			wantErr: false,
		},
		{
			name: "正常情况，无冲突",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			wantErr: false,
		},
		{
			name: "长名称冲突",
			setupParent: func() *types.CmdContext {
				parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
				existing := types.NewCmdContext("child", "e", flag.ContinueOnError)
				parent.SubCmdMap["child"] = existing
				return parent
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			wantErr:     true,
			errContains: "long name 'child' already exists",
		},
		{
			name: "短名称冲突",
			setupParent: func() *types.CmdContext {
				parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
				existing := types.NewCmdContext("existing", "c", flag.ContinueOnError)
				parent.SubCmdMap["c"] = existing
				return parent
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			wantErr:     true,
			errContains: "short name 'c' already exists",
		},
		{
			name: "循环引用检测",
			setupParent: func() *types.CmdContext {
				parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
				return parent
			},
			setupChild: func() *types.CmdContext {
				child := types.NewCmdContext("child", "c", flag.ContinueOnError)
				parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
				child.Parent = parent
				return child
			},
			wantErr:     true,
			errContains: "cyclic reference detected",
		},
		{
			name: "空长名称，只有短名称",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("", "c", flag.ContinueOnError)
			},
			wantErr: false,
		},
		{
			name: "空短名称，只有长名称",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "", flag.ContinueOnError)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := tt.setupParent()
			child := tt.setupChild()

			// 特殊处理循环引用测试用例
			if tt.name == "循环引用检测" && parent != nil && child != nil {
				child.Parent = parent
			}

			err := ValidateSubCommand(parent, child)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSubCommand() 期望返回错误，但得到 nil")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ValidateSubCommand() 错误信息 = %v, 期望包含 %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateSubCommand() 期望返回 nil，但得到错误 = %v", err)
				}
			}
		})
	}
}

// TestValidateSubCommand_PanicCase 测试会触发 panic 的情况
func TestValidateSubCommand_PanicCase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("期望触发 panic，但没有发生")
		}
	}()

	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	parent.SubCmdMap = nil // 故意设置为 nil 来触发 panic
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	_ = ValidateSubCommand(parent, child)
}

// TestHasCycleFast 测试 HasCycleFast 函数
func TestHasCycleFast(t *testing.T) {
	tests := []struct {
		name        string
		setupParent func() *types.CmdContext
		setupChild  func() *types.CmdContext
		want        bool
	}{
		{
			name: "父命令为nil",
			setupParent: func() *types.CmdContext {
				return nil
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			want: false,
		},
		{
			name: "子命令为nil",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return nil
			},
			want: false,
		},
		{
			name: "两者都为nil",
			setupParent: func() *types.CmdContext {
				return nil
			},
			setupChild: func() *types.CmdContext {
				return nil
			},
			want: false,
		},
		{
			name: "自引用情况",
			setupParent: func() *types.CmdContext {
				cmd := types.NewCmdContext("cmd", "c", flag.ContinueOnError)
				return cmd
			},
			setupChild: func() *types.CmdContext {
				// 这里会在测试中设置为同一个对象
				return nil
			},
			want: true,
		},
		{
			name: "无循环依赖",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			want: false,
		},
		{
			name: "直接循环依赖",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				child := types.NewCmdContext("child", "c", flag.ContinueOnError)
				parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
				child.Parent = parent
				return child
			},
			want: true,
		},
		{
			name: "间接循环依赖（三层）",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("grandparent", "gp", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				grandparent := types.NewCmdContext("grandparent", "gp", flag.ContinueOnError)
				parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
				child := types.NewCmdContext("child", "c", flag.ContinueOnError)

				parent.Parent = grandparent
				child.Parent = parent

				return child
			},
			want: true,
		},
		{
			name: "深层次无循环",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("newparent", "np", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				root := types.NewCmdContext("root", "r", flag.ContinueOnError)
				level1 := types.NewCmdContext("level1", "l1", flag.ContinueOnError)
				level2 := types.NewCmdContext("level2", "l2", flag.ContinueOnError)
				child := types.NewCmdContext("child", "c", flag.ContinueOnError)

				level1.Parent = root
				level2.Parent = level1
				child.Parent = level2

				return child
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := tt.setupParent()
			child := tt.setupChild()

			// 特殊处理自引用测试用例
			if tt.name == "自引用情况" {
				child = parent
			}

			// 特殊处理直接循环依赖测试用例
			if tt.name == "直接循环依赖" && parent != nil && child != nil {
				child.Parent = parent
			}

			// 特殊处理间接循环依赖测试用例
			if tt.name == "间接循环依赖（三层）" && parent != nil && child != nil {
				child.Parent.Parent = parent
			}

			got := HasCycleFast(parent, child)
			if got != tt.want {
				t.Errorf("HasCycleFast() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetCmdIdentifier 测试 GetCmdIdentifier 函数
func TestGetCmdIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		setupCmd func() *types.CmdContext
		want     string
	}{
		{
			name: "命令为nil",
			setupCmd: func() *types.CmdContext {
				return nil
			},
			want: "<nil>",
		},
		{
			name: "有长名称的命令",
			setupCmd: func() *types.CmdContext {
				return types.NewCmdContext("longname", "s", flag.ContinueOnError)
			},
			want: "longname",
		},
		{
			name: "只有短名称的命令",
			setupCmd: func() *types.CmdContext {
				return types.NewCmdContext("", "short", flag.ContinueOnError)
			},
			want: "short",
		},
		{
			name: "长名称和短名称都有",
			setupCmd: func() *types.CmdContext {
				return types.NewCmdContext("longname", "short", flag.ContinueOnError)
			},
			want: "longname", // 应该返回长名称
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			got := GetCmdIdentifier(cmd)
			if got != tt.want {
				t.Errorf("GetCmdIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHasCycleFast_DepthLimit 测试深度限制功能
func TestHasCycleFast_DepthLimit(t *testing.T) {
	// 创建一个深度超过10的命令链
	root := types.NewCmdContext("root", "r", flag.ContinueOnError)
	current := root

	// 创建12层深度的命令链
	for i := 1; i <= 12; i++ {
		next := types.NewCmdContext("level"+string(rune('0'+i)), "l"+string(rune('0'+i)), flag.ContinueOnError)
		next.Parent = current
		current = next
	}

	// 测试深度限制
	newParent := types.NewCmdContext("newparent", "np", flag.ContinueOnError)
	result := HasCycleFast(newParent, current)

	// 由于深度限制，应该返回false（没有检测到循环）
	if result != false {
		t.Errorf("HasCycleFast() 在深度限制情况下应该返回 false，但得到 %v", result)
	}
}

// TestValidateSubCommand_EdgeCases 测试边界情况
func TestValidateSubCommand_EdgeCases(t *testing.T) {
	t.Run("空字符串名称会panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("期望 NewCmdContext 对于空名称触发 panic，但没有发生")
			}
		}()

		parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
		// 这行应该会 panic
		child := types.NewCmdContext("", "", flag.ContinueOnError)

		// 这行代码不应该执行到，因为上面应该已经 panic 了
		_ = ValidateSubCommand(parent, child)
	})

	t.Run("重复添加相同命令", func(t *testing.T) {
		parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
		child := types.NewCmdContext("child", "c", flag.ContinueOnError)

		// 第一次添加应该成功
		err1 := ValidateSubCommand(parent, child)
		if err1 != nil {
			t.Errorf("第一次验证应该成功，但得到错误: %v", err1)
		}

		// 模拟添加到映射中
		parent.SubCmdMap["child"] = child
		parent.SubCmdMap["c"] = child

		// 第二次添加相同命令应该失败
		err2 := ValidateSubCommand(parent, child)
		if err2 == nil {
			t.Errorf("第二次验证应该失败，但得到 nil")
		}
	})
}

// 辅助函数：检查字符串是否包含子字符串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(substr) > 0 && findSubstring(s, substr)))
}

// 简单的子字符串查找实现
func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
