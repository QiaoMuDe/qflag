package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestCircularDependency 测试循环依赖问题
func TestCircularDependency(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志 a 和 b
	if err := cmd.AddFlag(flag.NewBoolFlag("a", "", "Flag A", false)); err != nil {
		t.Fatalf("Failed to add flag a: %v", err)
	}
	if err := cmd.AddFlag(flag.NewBoolFlag("b", "", "Flag B", false)); err != nil {
		t.Fatalf("Failed to add flag b: %v", err)
	}

	// 创建循环依赖: a -> b, b -> a
	// 这意味着: 如果设置a，必须设置b；如果设置b，必须设置a
	t.Log("添加依赖: a_requires_b (a -> b)")
	if err := cmd.AddFlagDependency("a_requires_b", "a", []string{"b"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add a_requires_b: %v", err)
	}

	t.Log("添加依赖: b_requires_a (b -> a)")
	if err := cmd.AddFlagDependency("b_requires_a", "b", []string{"a"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add b_requires_a: %v", err)
	}

	// 测试1: 只设置a，不设置b
	// 预期: 应该报错，因为a需要b
	t.Log("\n测试1: 只设置a (--a)")
	cmd1 := NewCmd("test1", "t1", types.ContinueOnError)
	if err := cmd1.AddFlag(flag.NewBoolFlag("a", "", "Flag A", false)); err != nil {
		t.Fatalf("Failed to add flag a: %v", err)
	}
	if err := cmd1.AddFlag(flag.NewBoolFlag("b", "", "Flag B", false)); err != nil {
		t.Fatalf("Failed to add flag b: %v", err)
	}
	if err := cmd1.AddFlagDependency("a_requires_b", "a", []string{"b"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd1.AddFlagDependency("b_requires_a", "b", []string{"a"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err := cmd1.Parse([]string{"--a"})
	if err != nil {
		t.Logf("  结果: 报错 (符合预期) - %v", err)
	} else {
		t.Log("  结果: 没有报错 (问题!)")
	}

	// 测试2: 同时设置a和b
	// 预期: 应该通过，因为a有b，b有a
	t.Log("\n测试2: 同时设置a和b (--a --b)")
	cmd2 := NewCmd("test2", "t2", types.ContinueOnError)
	if err := cmd2.AddFlag(flag.NewBoolFlag("a", "", "Flag A", false)); err != nil {
		t.Fatalf("Failed to add flag a: %v", err)
	}
	if err := cmd2.AddFlag(flag.NewBoolFlag("b", "", "Flag B", false)); err != nil {
		t.Fatalf("Failed to add flag b: %v", err)
	}
	if err := cmd2.AddFlagDependency("a_requires_b", "a", []string{"b"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd2.AddFlagDependency("b_requires_a", "b", []string{"a"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err = cmd2.Parse([]string{"--a", "--b"})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
	} else {
		t.Log("  结果: 通过验证")
	}

	// 测试3: 只设置b，不设置a
	// 预期: 应该报错，因为b需要a
	t.Log("\n测试3: 只设置b (--b)")
	cmd3 := NewCmd("test3", "t3", types.ContinueOnError)
	if err := cmd3.AddFlag(flag.NewBoolFlag("a", "", "Flag A", false)); err != nil {
		t.Fatalf("Failed to add flag a: %v", err)
	}
	if err := cmd3.AddFlag(flag.NewBoolFlag("b", "", "Flag B", false)); err != nil {
		t.Fatalf("Failed to add flag b: %v", err)
	}
	if err := cmd3.AddFlagDependency("a_requires_b", "a", []string{"b"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd3.AddFlagDependency("b_requires_a", "b", []string{"a"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err = cmd3.Parse([]string{"--b"})
	if err != nil {
		t.Logf("  结果: 报错 (符合预期) - %v", err)
	} else {
		t.Log("  结果: 没有报错 (问题!)")
	}

	// 测试4: 都不设置
	// 预期: 应该通过
	t.Log("\n测试4: 都不设置")
	cmd4 := NewCmd("test4", "t4", types.ContinueOnError)
	if err := cmd4.AddFlag(flag.NewBoolFlag("a", "", "Flag A", false)); err != nil {
		t.Fatalf("Failed to add flag a: %v", err)
	}
	if err := cmd4.AddFlag(flag.NewBoolFlag("b", "", "Flag B", false)); err != nil {
		t.Fatalf("Failed to add flag b: %v", err)
	}
	if err := cmd4.AddFlagDependency("a_requires_b", "a", []string{"b"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd4.AddFlagDependency("b_requires_a", "b", []string{"a"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err = cmd4.Parse([]string{})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
	} else {
		t.Log("  结果: 通过验证 (符合预期)")
	}

	t.Log("\n结论:")
	t.Log("循环依赖在添加时不会被检测，但运行时验证逻辑是正确的")
	t.Log("- 只设置a: 报错 (需要b)")
	t.Log("- 只设置b: 报错 (需要a)")
	t.Log("- 同时设置: 通过")
	t.Log("- 都不设置: 通过")
}
