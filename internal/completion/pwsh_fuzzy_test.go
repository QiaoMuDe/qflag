package completion

import (
	"bytes"
	"strings"
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

	// 验证模糊补全配置参数是否存在
	t.Run("配置参数检查", func(t *testing.T) {
		expectedConfigs := []string{
			"$script:testcli_FUZZY_COMPLETION_ENABLED = $true",
			"$script:testcli_FUZZY_MAX_CANDIDATES = 120",
			"$script:testcli_FUZZY_MIN_PATTERN_LENGTH = 2",
			"$script:testcli_FUZZY_SCORE_THRESHOLD = 25",
			"$script:testcli_FUZZY_MAX_RESULTS = 10",
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
			"function Get-testcliFuzzyScoreFast",
			"function Get-testcliFuzzyScoreCached",
			"function Get-testcliIntelligentMatches",
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
			"缓存大小控制",
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
		if !strings.Contains(result, "@{ Context = \"/\"; Options = @('--help', '--verbose', '--version', '--validate', 'build', 'test') }") {
			t.Error("根命令选项未正确生成")
		}

		// 检查标志参数
		if !strings.Contains(result, "@{ Context = \"/\"; Parameter = \"--validate\"; ParamType = \"required\"; ValueType = \"enum\"; Options = @('strict', 'loose', 'none') }") {
			t.Error("标志参数未正确生成")
		}
	})
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

	// 验证算法关键步骤
	t.Run("算法步骤检查", func(t *testing.T) {
		algorithmSteps := []string{
			"$patternLen = $Pattern.Length",
			"$candidateLen = $Candidate.Length",
			"if ($candidateLen -lt $patternLen)",
			"$Candidate.StartsWith($Pattern",
			"$patternLower.ToCharArray()",
			"$candidateLower.IndexOf($char)",
			"$matched = 0",
			"$consecutive = 0",
			"$maxConsecutive = 0",
			"$baseScore = [Math]::Floor(($matched * 60) / $patternLen)",
			"$consecutiveBonus = [Math]::Floor(($maxConsecutive * 20) / $patternLen)",
			"[Math]::Max(0, [Math]::Min(100, $finalScore))",
		}

		for _, step := range algorithmSteps {
			if !strings.Contains(result, step) {
				t.Errorf("生成的脚本缺少算法步骤: %s", step)
			}
		}
	})
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

	// 验证分级匹配策略
	t.Run("分级匹配策略检查", func(t *testing.T) {
		strategies := []string{
			"# 第1级: 精确前缀匹配 (最快，优先级最高)",
			"$option.StartsWith($Pattern, [System.StringComparison]::Ordinal)",
			"# 第2级: 大小写不敏感前缀匹配",
			"$option.StartsWith($Pattern, [System.StringComparison]::OrdinalIgnoreCase)",
			"# 第3级: 模糊匹配 (最慢，仅在必要时使用)",
			"Get-testFuzzyScoreCached -Pattern $Pattern -Candidate $option",
			"# 第4级: 子字符串匹配 (最后的备选方案)",
			"$optionLower.Contains($patternLower)",
		}

		for _, strategy := range strategies {
			if !strings.Contains(result, strategy) {
				t.Errorf("生成的脚本缺少匹配策略: %s", strategy)
			}
		}
	})

	// 验证性能保护机制
	t.Run("性能保护机制检查", func(t *testing.T) {
		protections := []string{
			"if ($totalCandidates -gt $script:test_FUZZY_MAX_CANDIDATES)",
			"# 回退到传统前缀匹配",
			"if ($exactMatches.Count -gt 0 -and $exactMatches.Count -le 12)",
			"if ($script:test_FUZZY_COMPLETION_ENABLED -and $patternLen -ge $script:test_FUZZY_MIN_PATTERN_LENGTH)",
			"if ($score -ge $script:test_FUZZY_SCORE_THRESHOLD)",
			"if ($count -ge $script:test_FUZZY_MAX_RESULTS)",
		}

		for _, protection := range protections {
			if !strings.Contains(result, protection) {
				t.Errorf("生成的脚本缺少性能保护机制: %s", protection)
			}
		}
	})
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

	// 验证缓存机制
	t.Run("缓存机制检查", func(t *testing.T) {
		cacheFeatures := []string{
			"$script:test_fuzzyCache = @{}",
			"$cacheKey = \"$Pattern|$Candidate\"",
			"if ($script:test_fuzzyCache.ContainsKey($cacheKey))",
			"return $script:test_fuzzyCache[$cacheKey]",
			"$script:test_fuzzyCache[$cacheKey] = $score",
			"if ($script:test_fuzzyCache.Count -gt $script:test_FUZZY_CACHE_MAX_SIZE)",
			"$script:test_fuzzyCache.Clear()",
		}

		for _, feature := range cacheFeatures {
			if !strings.Contains(result, feature) {
				t.Errorf("生成的脚本缺少缓存功能: %s", feature)
			}
		}
	})
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

	// 验证调试功能
	t.Run("调试功能检查", func(t *testing.T) {
		debugFeatures := []string{
			"function Get-myappCompletionDebug",
			"Write-Host \"=== myapp PowerShell补全系统诊断 ===\" -ForegroundColor Cyan",
			"Write-Host \"PowerShell版本: $($PSVersionTable.PSVersion)\" -ForegroundColor Green",
			"Write-Host \"模糊补全状态: $(if ($script:myapp_FUZZY_COMPLETION_ENABLED) { '启用' } else { '禁用' })\" -ForegroundColor Green",
			"function Test-myappFuzzyMatch",
			"$score = Get-myappFuzzyScoreFast -Pattern $Pattern -Candidate $Candidate",
			"Write-Host \"模式: '$Pattern' 匹配候选: '$Candidate' 得分: $score\"",
		}

		for _, feature := range debugFeatures {
			if !strings.Contains(result, feature) {
				t.Errorf("生成的脚本缺少调试功能: %s", feature)
			}
		}
	})
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
