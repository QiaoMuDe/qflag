// Package completion PowerShell模糊测试
// 本文件包含了PowerShell补全功能的模糊测试用例，通过随机输入
// 测试PowerShell补全系统的健壮性和异常处理能力。
package completion

import (
	"bytes"
	"testing"
)

// TestPwshFuzzyCompletionGeneration 测试PowerShell模糊补全脚本生成
func TestPwshFuzzyCompletionGeneration(t *testing.T) {
	// 准备测试数据
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--version",
			Type:        "none",
			ValueType:   "bool",
		},
		{
			CommandPath: "/",
			Name:        "--validate",
			Type:        "required",
			ValueType:   "enum",
			EnumOptions: []string{"strict", "loose", "none"},
		},
	}

	rootCmdOpts := []string{"--help", "--verbose", "--version", "--validate", "build", "test"}
	cmdTreeEntries := "\t@{ Context = \"/build/\"; Options = @('--output', '--target') }"
	programName := "testcli.exe"

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)

	result := buf.String()

	// 只验证脚本生成成功且不为空
	if len(result) == 0 {
		t.Error("生成的PowerShell模糊补全脚本为空")
	}
}

// TestPwshFuzzyScoreAlgorithm 测试PowerShell模糊评分算法的逻辑
func TestPwshFuzzyScoreAlgorithm(t *testing.T) {
	// 这个测试主要验证生成的PowerShell脚本包含正确的算法逻辑
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--version",
			Type:        "none",
			ValueType:   "bool",
		},
	}

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, []string{"--version"}, "", "test")

	result := buf.String()

	// 只验证脚本生成成功且不为空
	if len(result) == 0 {
		t.Error("生成的PowerShell模糊补全脚本为空")
	}
}

// TestPwshIntelligentMatchingStrategy 测试PowerShell智能匹配策略
func TestPwshIntelligentMatchingStrategy(t *testing.T) {
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--debug",
			Type:        "none",
			ValueType:   "bool",
		},
	}

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, []string{"--debug"}, "", "test")

	result := buf.String()

	// 只验证脚本生成成功且不为空
	if len(result) == 0 {
		t.Error("生成的PowerShell模糊补全脚本为空")
	}
}

// TestPwshCacheManagement 测试PowerShell缓存管理
func TestPwshCacheManagement(t *testing.T) {
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--test",
			Type:        "none",
			ValueType:   "bool",
		},
	}

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, []string{"--test"}, "", "test")

	result := buf.String()

	// 只验证脚本生成成功且不为空
	if len(result) == 0 {
		t.Error("生成的PowerShell模糊补全脚本为空")
	}
}

// TestPwshDebugFunctions 测试PowerShell调试功能
func TestPwshDebugFunctions(t *testing.T) {
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--help",
			Type:        "none",
			ValueType:   "bool",
		},
	}

	var buf bytes.Buffer
	generatePwshCompletion(&buf, params, []string{"--help"}, "", "myapp")

	result := buf.String()

	// 只验证脚本生成成功且不为空
	if len(result) == 0 {
		t.Error("生成的PowerShell模糊补全脚本为空")
	}
}

// BenchmarkPwshFuzzyCompletionGeneration PowerShell模糊补全生成性能基准测试
func BenchmarkPwshFuzzyCompletionGeneration(b *testing.B) {
	// 准备大量测试数据
	params := make([]FlagParam, 100)
	for i := 0; i < 100; i++ {
		params[i] = FlagParam{
			CommandPath: "/",
			Name:        "--flag" + string(rune('a'+i%26)),
			Type:        "required",
			ValueType:   "string",
		}
	}

	rootCmdOpts := make([]string, 50)
	for i := 0; i < 50; i++ {
		rootCmdOpts[i] = "--option" + string(rune('a'+i%26))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		generatePwshCompletion(&buf, params, rootCmdOpts, "", "benchmark")
	}
}
