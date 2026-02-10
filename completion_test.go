package qflag

import (
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestCompletionSpeed 测试补全脚本生成速度
//
// 参数:
//   - t: 测试实例
func TestCompletionSpeed(t *testing.T) {
	// 使用便捷函数创建测试命令
	cmd := createTestCmd()

	// 测试不同shell类型的补全脚本生成速度
	shells := []string{types.BashShell, types.PwshShell, types.PowershellShell}

	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			// 记录开始时间
			start := time.Now()

			// 生成补全脚本
			script, err := completion.Generate(cmd, shell)
			if err != nil {
				t.Errorf("生成补全脚本失败: %v", err)
				return
			}

			// 计算耗时
			duration := time.Since(start)

			// 验证脚本不为空
			if len(script) == 0 {
				t.Error("生成的补全脚本为空")
				return
			}

			// 输出性能指标
			t.Logf("Shell类型: %s", shell)
			t.Logf("生成时间: %v", duration)
			t.Logf("脚本大小: %d 字节", len(script))
			t.Logf("生成速度: %.2f KB/s", float64(len(script))/duration.Seconds()/1024)
		})
	}
}

// BenchmarkCompletionSpeed 基准测试补全脚本生成速度
//
// 参数:
//   - b: 基准测试实例
func BenchmarkCompletionSpeed(b *testing.B) {
	// 使用便捷函数创建测试命令
	cmd := createTestCmd()

	// 测试不同shell类型的基准性能
	shells := []string{types.BashShell, types.PwshShell, types.PowershellShell}

	for _, shell := range shells {
		b.Run(shell, func(b *testing.B) {
			// 重置计时器
			b.ResetTimer()

			// 运行基准测试
			for i := 0; i < b.N; i++ {
				_, err := completion.Generate(cmd, shell)
				if err != nil {
					b.Fatalf("生成补全脚本失败: %v", err)
				}
			}
		})
	}
}

// 生成完整的补全脚本
//
// 参数:
//   - cmd: 命令接口
//   - shellType: Shell类型 (bash, pwsh, powershell)
//
// 返回值:
//   - 包含所有标志选项和子命令名称的字符串切片
func TestGenerateFullCompletionScript(t *testing.T) {
	// 使用便捷函数创建测试命令
	cmd := createTestCmd()
	//cmd.SetCompletion(true)
	if err := cmd.Parse([]string{}); err != nil {
		t.Errorf("解析命令行参数失败: %v", err)
	}

	//completion.GenAndPrint(cmd, types.BashShell)

	cmd.PrintHelp()
}
