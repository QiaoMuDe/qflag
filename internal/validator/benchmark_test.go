package validator

import (
	"flag"
	"fmt"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

// BenchmarkValidateSubCommand_Simple 简单验证的基准测试
func BenchmarkValidateSubCommand_Simple(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ValidateSubCommand(parent, child)
		if err != nil {
			b.Fatalf("验证子命令失败: %v", err)
		}
	}
}

// BenchmarkValidateSubCommand_WithConflict 有冲突的验证基准测试
func BenchmarkValidateSubCommand_WithConflict(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	existing := types.NewCmdContext("existing", "e", flag.ContinueOnError)
	conflicting := types.NewCmdContext("existing", "c", flag.ContinueOnError)

	// 预先添加一个命令
	parent.SubCmdMap["existing"] = existing

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ValidateSubCommand(parent, conflicting)
		if err == nil {
			b.Fatal("期望验证冲突的子命令时返回错误")
		}
	}
}

// BenchmarkValidateSubCommand_LargeSubCmdMap 大量子命令映射的基准测试
func BenchmarkValidateSubCommand_LargeSubCmdMap(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	// 预先添加1000个子命令
	for i := 0; i < 1000; i++ {
		cmd := types.NewCmdContext(fmt.Sprintf("cmd%d", i), fmt.Sprintf("c%d", i), flag.ContinueOnError)
		parent.SubCmdMap[cmd.LongName] = cmd
		parent.SubCmdMap[cmd.ShortName] = cmd
	}

	child := types.NewCmdContext("newchild", "nc", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ValidateSubCommand(parent, child)
		if err != nil {
			b.Fatalf("验证子命令失败: %v", err)
		}
	}
}

// BenchmarkHasCycle_NoCycle 无循环情况的基准测试
func BenchmarkHasCycle_NoCycle(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycle(parent, child)
	}
}

// BenchmarkHasCycle_DirectCycle 直接循环的基准测试
func BenchmarkHasCycle_DirectCycle(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycle(parent, parent)
	}
}

// BenchmarkHasCycle_DeepChain 深层链的基准测试
func BenchmarkHasCycle_DeepChain(b *testing.B) {
	// 创建不同深度的命令链进行基准测试
	depths := []int{10, 50, 99}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("深度%d", depth), func(b *testing.B) {
			commands := make([]*types.CmdContext, depth)
			for i := 0; i < depth; i++ {
				commands[i] = types.NewCmdContext(
					fmt.Sprintf("cmd%d", i),
					fmt.Sprintf("c%d", i),
					flag.ContinueOnError,
				)
				if i > 0 {
					commands[i].Parent = commands[i-1]
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				HasCycle(commands[depth-1], commands[0])
			}
		})
	}
}

// BenchmarkHasCycle_MaxDepth 最大深度保护的基准测试
func BenchmarkHasCycle_MaxDepth(b *testing.B) {
	// 创建超过100层的命令链
	const depth = 150
	commands := make([]*types.CmdContext, depth)

	for i := 0; i < depth; i++ {
		commands[i] = types.NewCmdContext(
			fmt.Sprintf("cmd%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)
		if i > 0 {
			commands[i].Parent = commands[i-1]
		}
	}

	newCmd := types.NewCmdContext("new", "n", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycle(commands[depth-1], newCmd)
	}
}

// BenchmarkGetCmdIdentifier_LongName 长名称获取的基准测试
func BenchmarkGetCmdIdentifier_LongName(b *testing.B) {
	cmd := types.NewCmdContext("verylongcommandname", "v", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCmdIdentifier(cmd)
	}
}

// BenchmarkGetCmdIdentifier_ShortName 短名称获取的基准测试
func BenchmarkGetCmdIdentifier_ShortName(b *testing.B) {
	cmd := types.NewCmdContext("", "s", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCmdIdentifier(cmd)
	}
}

// BenchmarkGetCmdIdentifier_Nil nil命令的基准测试
func BenchmarkGetCmdIdentifier_Nil(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCmdIdentifier(nil)
	}
}

// BenchmarkValidateSubCommand_Parallel 并行验证的基准测试
func BenchmarkValidateSubCommand_Parallel(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			child := types.NewCmdContext(fmt.Sprintf("child%d", i), fmt.Sprintf("c%d", i), flag.ContinueOnError)
			err := ValidateSubCommand(parent, child)
			if err != nil {
				b.Fatalf("验证子命令失败: %v", err)
			}
			i++
		}
	})
}

// BenchmarkHasCycle_Parallel 并行循环检测的基准测试
func BenchmarkHasCycle_Parallel(b *testing.B) {
	// 创建一个中等深度的命令链
	commands := make([]*types.CmdContext, 20)
	for i := 0; i < 20; i++ {
		commands[i] = types.NewCmdContext(
			fmt.Sprintf("cmd%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)
		if i > 0 {
			commands[i].Parent = commands[i-1]
		}
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			parentIdx := i % 20
			childIdx := (i + 10) % 20
			HasCycle(commands[parentIdx], commands[childIdx])
			i++
		}
	})
}

// BenchmarkComplexValidationScenario 复杂验证场景的基准测试
func BenchmarkComplexValidationScenario(b *testing.B) {
	// 创建一个复杂的命令层次结构
	root := types.NewCmdContext("root", "r", flag.ContinueOnError)

	// 添加多层子命令
	for i := 0; i < 10; i++ {
		level1 := types.NewCmdContext(fmt.Sprintf("l1_%d", i), fmt.Sprintf("l1_%d", i), flag.ContinueOnError)
		level1.Parent = root
		root.SubCmdMap[level1.LongName] = level1

		for j := 0; j < 5; j++ {
			level2 := types.NewCmdContext(fmt.Sprintf("l2_%d_%d", i, j), fmt.Sprintf("l2_%d_%d", i, j), flag.ContinueOnError)
			level2.Parent = level1
			level1.SubCmdMap[level2.LongName] = level2
		}
	}

	newChild := types.NewCmdContext("newchild", "nc", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ValidateSubCommand(root, newChild)
		if err != nil {
			b.Fatalf("验证子命令失败: %v", err)
		}
	}
}
