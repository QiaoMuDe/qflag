package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestAddMutexGroup 测试添加互斥组
func TestAddMutexGroup(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
		t.Fatalf("Failed to add format flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
		t.Fatalf("Failed to add output flag: %v", err)
	}

	// 添加互斥组
	cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)

	// 验证互斥组已添加
	groups := cmd.GetMutexGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 mutex group, got %d", len(groups))
	}

	group := groups[0]
	if group.Name != "output_format" {
		t.Errorf("Expected group name 'output_format', got '%s'", group.Name)
	}

	if len(group.Flags) != 2 {
		t.Fatalf("Expected 2 flags in group, got %d", len(group.Flags))
	}

	if group.Flags[0] != "format" || group.Flags[1] != "output" {
		t.Errorf("Expected flags ['format', 'output'], got %v", group.Flags)
	}

	if !group.AllowNone {
		t.Errorf("Expected AllowNone to be true, got false")
	}
}

// TestGetMutexGroups 测试获取互斥组列表
func TestGetMutexGroups(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加多个互斥组
	cmd.AddMutexGroup("group1", []string{"flag1", "flag2"}, true)
	cmd.AddMutexGroup("group2", []string{"flag3", "flag4"}, false)

	// 获取互斥组列表
	groups := cmd.GetMutexGroups()
	if len(groups) != 2 {
		t.Fatalf("Expected 2 mutex groups, got %d", len(groups))
	}

	// 修改返回的列表, 不应影响原始数据
	groups[0].Name = "modified"
	modifiedGroups := cmd.GetMutexGroups()
	if modifiedGroups[0].Name == "modified" {
		t.Error("Modifying returned groups should not affect original data")
	}
}

// TestRemoveMutexGroup 测试移除互斥组
func TestRemoveMutexGroup(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加互斥组
	cmd.AddMutexGroup("group1", []string{"flag1", "flag2"}, true)
	cmd.AddMutexGroup("group2", []string{"flag3", "flag4"}, false)

	// 移除存在的互斥组
	if !cmd.RemoveMutexGroup("group1") {
		t.Error("Expected RemoveMutexGroup to return true for existing group")
	}

	// 验证互斥组已移除
	groups := cmd.GetMutexGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 mutex group after removal, got %d", len(groups))
	}

	if groups[0].Name != "group2" {
		t.Errorf("Expected remaining group name 'group2', got '%s'", groups[0].Name)
	}

	// 尝试移除不存在的互斥组
	if cmd.RemoveMutexGroup("nonexistent") {
		t.Error("Expected RemoveMutexGroup to return false for non-existing group")
	}
}

// TestMutexGroupValidation 测试互斥组验证
func TestMutexGroupValidation(t *testing.T) {
	// 测试1: 不设置任何标志, 应该成功 (AllowNone=true)
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)

		err := cmd.Parse([]string{})
		if err != nil {
			t.Errorf("Expected no error when no flags are set with AllowNone=true, got: %v", err)
		}
	}()

	// 测试2: 只设置一个标志, 应该成功
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)

		err := cmd.Parse([]string{"--format", "json"})
		if err != nil {
			t.Errorf("Expected no error when only one flag is set, got: %v", err)
		}
	}()

	// 测试3: 设置互斥的两个标志, 应该失败
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)

		err := cmd.Parse([]string{"--format", "json", "--output", "result.txt"})
		if err == nil {
			t.Error("Expected error when both mutually exclusive flags are set")
		}
	}()

	// 测试4: 添加不允许为空的互斥组
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("input", "i", "Input file", "")); err != nil {
			t.Fatalf("Failed to add input flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("source", "s", "Source file", "")); err != nil {
			t.Fatalf("Failed to add source flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)
		cmd.AddMutexGroup("input_source", []string{"input", "source"}, false)

		// 不设置任何标志, 应该失败 (因为第二个互斥组不允许为空)
		err := cmd.Parse([]string{})
		if err == nil {
			t.Error("Expected error when no flag is set for group with AllowNone=false")
		}
	}()

	// 测试5: 只设置一个标志, 应该成功
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("input", "i", "Input file", "")); err != nil {
			t.Fatalf("Failed to add input flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("source", "s", "Source file", "")); err != nil {
			t.Fatalf("Failed to add source flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)
		cmd.AddMutexGroup("input_source", []string{"input", "source"}, false)

		err := cmd.Parse([]string{"--input", "file.txt"})
		if err != nil {
			t.Errorf("Expected no error when one flag is set for group with AllowNone=false, got: %v", err)
		}
	}()
}

// TestMutexGroupWithNonExistentFlags 测试包含不存在标志的互斥组
func TestMutexGroupWithNonExistentFlags(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
		t.Fatalf("Failed to add format flag: %v", err)
	}

	// 添加包含不存在标志的互斥组
	cmd.AddMutexGroup("test_group", []string{"format", "nonexistent"}, true)

	// 设置存在的标志, 应该成功
	err := cmd.Parse([]string{"--format", "json"})
	if err != nil {
		t.Errorf("Expected no error when existing flag is set, got: %v", err)
	}
}

// TestMultipleMutexGroups 测试多个互斥组
func TestMultipleMutexGroups(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
		t.Fatalf("Failed to add format flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
		t.Fatalf("Failed to add output flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("input", "i", "Input file", "")); err != nil {
		t.Fatalf("Failed to add input flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("source", "s", "Source file", "")); err != nil {
		t.Fatalf("Failed to add source flag: %v", err)
	}

	// 添加两个互斥组
	cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)
	cmd.AddMutexGroup("input_source", []string{"input", "source"}, false)

	// 测试1: 每个互斥组设置一个标志, 应该成功
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("input", "i", "Input file", "")); err != nil {
			t.Fatalf("Failed to add input flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("source", "s", "Source file", "")); err != nil {
			t.Fatalf("Failed to add source flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)
		cmd.AddMutexGroup("input_source", []string{"input", "source"}, false)

		err := cmd.Parse([]string{"--format", "json", "--input", "file.txt"})
		if err != nil {
			t.Errorf("Expected no error when one flag from each group is set, got: %v", err)
		}
	}()

	// 测试2: 一个互斥组设置多个标志, 应该失败
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("input", "i", "Input file", "")); err != nil {
			t.Fatalf("Failed to add input flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("source", "s", "Source file", "")); err != nil {
			t.Fatalf("Failed to add source flag: %v", err)
		}
		cmd.AddMutexGroup("output_format", []string{"format", "output"}, true)
		cmd.AddMutexGroup("input_source", []string{"input", "source"}, false)

		err := cmd.Parse([]string{"--format", "json", "--output", "result.txt", "--input", "file.txt"})
		if err == nil {
			t.Error("Expected error when multiple flags from one mutex group are set")
		}
	}()
}

// TestMutexGroupConcurrency 测试互斥组的并发安全性
func TestMutexGroupConcurrency(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 并发添加互斥组
	done := make(chan bool, 2)

	go func() {
		cmd.AddMutexGroup("group1", []string{"flag1", "flag2"}, true)
		done <- true
	}()

	go func() {
		cmd.AddMutexGroup("group2", []string{"flag3", "flag4"}, false)
		done <- true
	}()

	// 等待两个goroutine完成
	<-done
	<-done

	// 验证两个互斥组都已添加
	groups := cmd.GetMutexGroups()
	if len(groups) != 2 {
		t.Fatalf("Expected 2 mutex groups, got %d", len(groups))
	}
}
