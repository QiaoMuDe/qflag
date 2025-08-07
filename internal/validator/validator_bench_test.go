package validator

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

// BenchmarkValidateSubCommand 基准测试 ValidateSubCommand 函数
func BenchmarkValidateSubCommand(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateSubCommand(parent, child)
	}
}

// BenchmarkValidateSubCommand_WithConflict 基准测试有冲突的情况
func BenchmarkValidateSubCommand_WithConflict(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	existing := types.NewCmdContext("child", "e", flag.ContinueOnError)
	parent.SubCmdMap["child"] = existing

	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateSubCommand(parent, child)
	}
}

// BenchmarkHasCycleFast 基准测试 HasCycleFast 函数
func BenchmarkHasCycleFast(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycleFast(parent, child)
	}
}

// BenchmarkHasCycleFast_WithCycle 基准测试有循环的情况
func BenchmarkHasCycleFast_WithCycle(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)
	child.Parent = parent

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycleFast(parent, child)
	}
}

// BenchmarkHasCycleFast_DeepChain 基准测试深层命令链
func BenchmarkHasCycleFast_DeepChain(b *testing.B) {
	// 创建深度为5的命令链
	root := types.NewCmdContext("root", "r", flag.ContinueOnError)
	current := root

	for i := 1; i <= 5; i++ {
		next := types.NewCmdContext("level"+string(rune('0'+i)), "l"+string(rune('0'+i)), flag.ContinueOnError)
		next.Parent = current
		current = next
	}

	parent := types.NewCmdContext("newparent", "np", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycleFast(parent, current)
	}
}

// BenchmarkGetCmdIdentifier 基准测试 GetCmdIdentifier 函数
func BenchmarkGetCmdIdentifier(b *testing.B) {
	cmd := types.NewCmdContext("longname", "s", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCmdIdentifier(cmd)
	}
}

// BenchmarkGetCmdIdentifier_Nil 基准测试 nil 情况
func BenchmarkGetCmdIdentifier_Nil(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCmdIdentifier(nil)
	}
}

// BenchmarkValidateSubCommand_MultipleSubCommands 基准测试多个子命令的情况
func BenchmarkValidateSubCommand_MultipleSubCommands(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	// 添加多个已存在的子命令
	for i := 0; i < 100; i++ {
		existing := types.NewCmdContext("cmd"+string(rune('0'+i%10)), "c"+string(rune('0'+i%10)), flag.ContinueOnError)
		parent.SubCmdMap["cmd"+string(rune('0'+i%10))] = existing
	}

	child := types.NewCmdContext("newchild", "nc", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateSubCommand(parent, child)
	}
}
