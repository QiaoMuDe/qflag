package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestTargetFlagNotExists 测试目标标志不存在时的错误处理
func TestTargetFlagNotExists(t *testing.T) {
	t.Log("=== 测试目标标志不存在时的错误处理 ===")
	t.Log()

	// 场景1: 正常情况（目标标志存在）
	t.Log("场景1: 正常情况（目标标志存在）")
	cmd1 := NewCmd("test1", "t1", types.ContinueOnError)
	if err := cmd1.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd1.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd1.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err := cmd1.Parse([]string{"--ssl"})
	if err != nil {
		t.Logf("  结果: 报错 (符合预期) - %v", err)
	} else {
		t.Log("  结果: 通过")
	}

	// 场景2: 手动修改配置，添加一个不存在的目标标志
	t.Log("\n场景2: 目标标志不存在（通过手动修改配置模拟）")
	cmd2 := NewCmd("test2", "t2", types.ContinueOnError)
	if err := cmd2.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd2.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd2.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	// 手动添加一个不存在的目标标志到依赖关系中
	// 这模拟了配置被错误修改的情况
	config := cmd2.Config()
	if len(config.FlagDependencies) > 0 {
		config.FlagDependencies[0].Targets = []string{"nonexistent"} // 只保留不存在的标志
	}

	t.Log("  修改后的依赖关系:")
	t.Logf("    Trigger: %s", config.FlagDependencies[0].Trigger)
	t.Logf("    Targets: %v", config.FlagDependencies[0].Targets)
	t.Log("  注意: 'nonexistent' 标志不存在")
	t.Log()

	err = cmd2.Parse([]string{"--ssl"})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
		if containsStr(err.Error(), "nonexistent") {
			t.Log("  分析: 错误信息中提到了不存在的标志")
		} else if containsStr(err.Error(), "not found") || containsStr(err.Error(), "invalid") {
			t.Log("  分析: 正确检测到标志不存在")
		} else {
			t.Log("  分析: 错误信息可能不清晰")
		}
	} else {
		t.Log("  结果: 通过验证 (问题! 应该报错)")
	}

	// 场景3: 触发标志存在，但目标标志不存在（互斥依赖）
	t.Log("\n场景3: 互斥依赖中目标标志不存在")
	cmd3 := NewCmd("test3", "t3", types.ContinueOnError)
	if err := cmd3.AddFlag(flag.NewBoolFlag("debug", "d", "Debug mode", false)); err != nil {
		t.Fatalf("Failed to add flag debug: %v", err)
	}
	if err := cmd3.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd3.AddFlagDependency("debug_mutex_ssl", "debug", []string{"ssl"}, types.DepMutex); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	// 手动修改，添加不存在的目标
	config3 := cmd3.Config()
	if len(config3.FlagDependencies) > 0 {
		config3.FlagDependencies[0].Targets = append(config3.FlagDependencies[0].Targets, "nonexistent")
	}

	t.Log("  修改后的依赖关系:")
	t.Logf("    Trigger: %s", config3.FlagDependencies[0].Trigger)
	t.Logf("    Targets: %v", config3.FlagDependencies[0].Targets)
	t.Log()

	err = cmd3.Parse([]string{"--debug"})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
	} else {
		t.Log("  结果: 通过验证")
	}

	// 场景4: 检查 validateFlagDependencies 中的代码逻辑
	t.Log("\n场景4: 代码逻辑分析")
	t.Log("在 validateFlagDependencies 中:")
	t.Log("  displayName := p.flagDisplayNames[target]")
	t.Log("  如果 target 不在 flagDisplayNames 中，displayName 为空字符串")
	t.Log()
	t.Log("  对于 DepRequired:")
	t.Log("    if !setFlags[target] { // target 未设置")
	t.Log("      displayName := p.flagDisplayNames[target] // 可能为空")
	t.Log("      missingFlags = append(missingFlags, displayName) // 添加空字符串")
	t.Log("    }")
	t.Log()
	t.Log("  对于 DepMutex:")
	t.Log("    if setFlags[target] { // target 已设置")
	t.Log("      // 不会执行，因为不存在的标志不可能被设置")
	t.Log("    }")
	t.Log()
	t.Log("  结论:")
	t.Log("    - DepRequired: 可能添加空字符串到错误列表")
	t.Log("    - DepMutex: 不会触发，因为不存在的标志不会被设置")
}

// TestTriggerFlagNotExists 测试触发标志不存在的情况
func TestTriggerFlagNotExists(t *testing.T) {
	t.Log("\n=== 测试触发标志不存在 ===")

	cmd := NewCmd("test", "t", types.ContinueOnError)
	if err := cmd.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	// 手动修改，将触发标志改为不存在的标志
	config := cmd.Config()
	if len(config.FlagDependencies) > 0 {
		config.FlagDependencies[0].Trigger = "nonexistent"
	}

	t.Log("修改后的依赖关系:")
	t.Logf("  Trigger: %s (不存在)", config.FlagDependencies[0].Trigger)
	t.Logf("  Targets: %v", config.FlagDependencies[0].Targets)
	t.Log()

	// 设置 cert，但不设置 ssl（因为 ssl 不存在）
	err := cmd.Parse([]string{"--cert", "cert.pem"})
	if err != nil {
		t.Logf("结果: 报错 - %v", err)
	} else {
		t.Log("结果: 通过验证")
		t.Log("分析: 触发标志不存在，所以依赖检查被跳过")
	}
}
