// Package completion Bash模糊测试
// 本文件包含了Bash补全功能的模糊测试用例，通过随机输入
// 测试补全系统的健壮性和异常处理能力。
package completion

import (
	"bytes"
	"strings"
	"testing"
)

// TestBashFuzzyCompletionGeneration 测试Bash模糊补全脚本生成
func TestBashFuzzyCompletionGeneration(t *testing.T) {
	// 准备测试数据
	params := []FlagParam{
		{
			CommandPath: "/",
			Name:        "--verbose",
			Type:        "optional",
			ValueType:   "bool",
		},
		{
			CommandPath: "/",
			Name:        "--version",
			Type:        "optional",
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
	cmdTreeEntries := "cmd_tree[/build/]=\"--output|--target\"\n"
	programName := "testcli"

	var buf bytes.Buffer
	generateBashCompletion(&buf, params, rootCmdOpts, cmdTreeEntries, programName)

	result := buf.String()

	// 验证模糊补全配置参数是否存在
	t.Run("配置参数检查", func(t *testing.T) {
		expectedConfigs := []string{
			"testcli_FUZZY_COMPLETION_ENABLED=1",
			"testcli_FUZZY_MAX_CANDIDATES=150",
			"testcli_FUZZY_MIN_PATTERN_LENGTH=2",
			"testcli_FUZZY_SCORE_THRESHOLD=30",
			"testcli_FUZZY_MAX_RESULTS=8",
		}

		for _, config := range expectedConfigs {
			if !strings.Contains(result, config) {
				t.Errorf("生成的脚本缺少配置参数: %s", config)
			}
		}
	})

	// 验证核心函数是否存在
	t.Run("核心函数检查", func(t *testing.T) {
		expectedFunctions := []string{
			"_fuzzy_score_fast()",
			"_testcli_fuzzy_score_cached()",
			"_testcli_intelligent_match()",
			"_testcli_completion_debug()",
		}

		for _, function := range expectedFunctions {
			if !strings.Contains(result, function) {
				t.Errorf("生成的脚本缺少核心函数: %s", function)
			}
		}
	})

	// 验证性能优化代码是否存在
	t.Run("性能优化检查", func(t *testing.T) {
		expectedOptimizations := []string{
			"长度预检查",
			"完全匹配检查",
			"字符存在性预检查",
			"缓存命中检查",
			"性能保护: 候选项过多时禁用模糊匹配",
			"分级匹配策略",
		}

		for _, optimization := range expectedOptimizations {
			if !strings.Contains(result, optimization) {
				t.Errorf("生成的脚本缺少性能优化: %s", optimization)
			}
		}
	})

	// 验证数据结构是否正确
	t.Run("数据结构检查", func(t *testing.T) {
		// 检查命令树
		if !strings.Contains(result, "testcli_cmd_tree[/]=\"--help|--verbose|--version|--validate|build|test\"") {
			t.Error("根命令选项未正确生成")
		}

		// 检查标志参数
		if !strings.Contains(result, "testcli_flag_params[\"/|--validate\"]=\"required|enum\"") {
			t.Error("标志参数未正确生成")
		}

		// 检查枚举选项
		if !strings.Contains(result, "testcli_enum_options[\"/|--validate\"]=\"strict|loose|none\"") {
			t.Error("枚举选项未正确生成")
		}
	})

	// 验证脚本结构完整性
	t.Run("脚本结构检查", func(t *testing.T) {
		// 检查shebang
		if !strings.HasPrefix(result, "#!/usr/bin/env bash") {
			t.Error("脚本缺少正确的shebang")
		}

		// 检查主补全函数
		if !strings.Contains(result, "_testcli() {") {
			t.Error("主补全函数未正确生成")
		}

		// 检查函数注册
		if !strings.Contains(result, "complete -F _testcli testcli") {
			t.Error("补全函数注册未正确生成")
		}
	})
}

// TestBashEscapeSpecialChars 测试特殊字符转义功能
func TestBashEscapeSpecialChars(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		desc     string
	}{
		{
			input:    "simple",
			expected: "simple",
			desc:     "普通字符串不需要转义",
		},
		{
			input:    "hello world",
			expected: "hello\\ world",
			desc:     "空格需要转义",
		},
		{
			input:    "test$var",
			expected: "test\\$var",
			desc:     "美元符号需要转义",
		},
		{
			input:    "path/to/file",
			expected: "path/to/file",
			desc:     "路径分隔符不需要转义",
		},
		{
			input:    "cmd|pipe&bg",
			expected: "cmd\\|pipe\\&bg",
			desc:     "管道符和与符号需要转义",
		},
		{
			input:    "test\"quote'mix",
			expected: "test\\\"quote'mix",
			desc:     "双引号需要转义，单引号不需要",
		},
		{
			input:    "complex*?[test]",
			expected: "complex\\*\\?\\[test\\]",
			desc:     "通配符和方括号需要转义",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := escapeSpecialChars(tc.input)
			if result != tc.expected {
				t.Errorf("输入: %q, 期望: %q, 实际: %q", tc.input, tc.expected, result)
			}
		})
	}
}

// TestBashCommandTreeEntry 测试命令树条目生成
func TestBashCommandTreeEntry(t *testing.T) {
	var buf bytes.Buffer
	cmdPath := "/build/"
	cmdOpts := []string{"--output", "--target", "--verbose"}

	generateBashCommandTreeEntry(&buf, cmdPath, cmdOpts, "testprogram")

	result := buf.String()
	expected := "testprogram_cmd_tree[/build/]=\"--output|--target|--verbose\"\n"

	if result != expected {
		t.Errorf("命令树条目生成错误\n期望: %q\n实际: %q", expected, result)
	}
}

// BenchmarkFuzzyScorePerformance 模糊评分性能基准测试
func BenchmarkFuzzyScorePerformance(t *testing.B) {
	// 这个基准测试主要用于文档说明，实际的bash函数无法在Go中直接测试
	// 但可以用来验证我们的算法设计思路

	testCases := []struct {
		pattern   string
		candidate string
		desc      string
	}{
		{"vb", "--verbose", "短模式匹配长选项"},
		{"ver", "--version", "前缀匹配"},
		{"val", "--validate", "前缀匹配"},
		{"bld", "build", "模糊匹配"},
		{"tst", "test", "模糊匹配"},
	}

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		for _, tc := range testCases {
			// 模拟bash中的模糊评分逻辑
			_ = simulateBashFuzzyScore(tc.pattern, tc.candidate)
		}
	}
}

// simulateBashFuzzyScore 模拟bash中的模糊评分算法
// 这个函数用于性能测试和算法验证
func simulateBashFuzzyScore(pattern, candidate string) int {
	patternLen := len(pattern)
	candidateLen := len(candidate)

	// 长度预检查
	if candidateLen < patternLen {
		return 0
	}

	// 完全匹配检查
	if strings.HasPrefix(candidate, pattern) {
		return 100
	}

	// 字符存在性预检查
	patternLower := strings.ToLower(pattern)
	candidateLower := strings.ToLower(candidate)

	for _, char := range patternLower {
		if !strings.ContainsRune(candidateLower, char) {
			return 0
		}
	}

	// 核心匹配算法
	matched := 0
	consecutive := 0
	maxConsecutive := 0
	candidatePos := 0
	startBonus := 0

	if strings.HasPrefix(candidateLower, patternLower) {
		startBonus = 20
	}

	for i, patternChar := range patternLower {
		found := false
		for j := candidatePos; j < len(candidateLower); j++ {
			if rune(candidateLower[j]) == patternChar {
				matched++
				found = true

				if j == candidatePos {
					consecutive++
					if consecutive > maxConsecutive {
						maxConsecutive = consecutive
					}
				} else {
					consecutive = 1
				}

				candidatePos = j + 1
				break
			}
		}

		if !found {
			consecutive = 0
		}

		_ = i // 避免未使用变量警告
	}

	// 评分计算
	baseScore := matched * 60 / patternLen
	consecutiveBonus := maxConsecutive * 20 / patternLen
	lengthPenalty := candidateLen - patternLen
	if lengthPenalty > 10 {
		lengthPenalty = 10
	}

	finalScore := baseScore + consecutiveBonus + startBonus - lengthPenalty

	if finalScore < 0 {
		finalScore = 0
	} else if finalScore > 100 {
		finalScore = 100
	}

	return finalScore
}

// TestSimulateFuzzyScore 测试模拟的模糊评分算法
func TestSimulateFuzzyScore(t *testing.T) {
	testCases := []struct {
		pattern   string
		candidate string
		minScore  int
		desc      string
	}{
		{"ver", "--version", 60, "前缀匹配应该得高分"},
		{"vb", "--verbose", 30, "模糊匹配应该得中等分数"},
		{"xyz", "--verbose", 0, "不匹配应该得0分"},
		{"help", "--help", 70, "包含匹配应该得高分"},
		{"bld", "build", 50, "模糊匹配应该得合理分数"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			score := simulateBashFuzzyScore(tc.pattern, tc.candidate)
			if score < tc.minScore {
				t.Errorf("模式 %q 匹配候选 %q 的分数 %d 低于期望的最低分数 %d",
					tc.pattern, tc.candidate, score, tc.minScore)
			}
		})
	}
}
