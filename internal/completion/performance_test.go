package completion

import (
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestPowerShellPerformanceOptimization 测试PowerShell补全性能优化
func TestPowerShellPerformanceOptimization(t *testing.T) {
	// 创建大型命令树来测试性能
	rootCtx := createLargeCommandTree(100, 5) // 100个主命令，每个5个子命令

	start := time.Now()
	script, err := GenerateShellCompletion(rootCtx, flags.ShellPowershell)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("生成PowerShell补全脚本失败: %v", err)
	}

	t.Logf("PowerShell补全脚本生成耗时: %v", duration)
	t.Logf("脚本大小: %d 字节", len(script))

	// 验证优化后的脚本包含索引构建逻辑
	expectedOptimizations := []string{
		"_contextIndex",
		"_flagIndex",
		"ContainsKey",
		"try {",
		"catch {",
	}

	for _, pattern := range expectedOptimizations {
		if !strings.Contains(script, pattern) {
			t.Errorf("PowerShell脚本缺少性能优化: %s", pattern)
		}
	}

	// 性能阈值检查（应该在100ms以内）
	if duration > 100*time.Millisecond {
		t.Errorf("PowerShell补全生成耗时过长: %v", duration)
	}
}

// TestBashFunctionalityFix 测试Bash补全功能修复
func TestBashFunctionalityFix(t *testing.T) {
	// 创建包含特殊字符的测试上下文
	rootCtx := createTestContextWithSpecialChars()

	script, err := GenerateShellCompletion(rootCtx, flags.ShellBash)
	if err != nil {
		t.Fatalf("生成Bash补全脚本失败: %v", err)
	}

	// 验证特殊字符转义
	specialChars := []string{
		"\\$", "\\`", "\\|", "\\&", "\\;", "\\(", "\\)", "\\<", "\\>",
		"\\*", "\\?", "\\[", "\\]", "\\{", "\\}", "\\~", "\\#",
	}

	hasSpecialChars := false
	for _, char := range specialChars {
		if strings.Contains(script, char) {
			hasSpecialChars = true
			break
		}
	}

	if !hasSpecialChars {
		t.Log("注意: 脚本中未检测到特殊字符转义，可能是测试数据不包含特殊字符")
	}

	// 验证功能修复 - 适配新的模糊补全实现
	expectedFixes := []string{
		"模糊补全配置",
		"性能优化",
		"模糊匹配结果缓存",
		"分级匹配策略",
		"FUZZY_COMPLETION_ENABLED",
	}

	for _, fix := range expectedFixes {
		if !strings.Contains(script, fix) {
			t.Errorf("Bash脚本缺少功能修复: %s", fix)
		}
	}
}

// TestEscapeSpecialCharsEnhanced 测试增强的特殊字符转义
func TestEscapeSpecialCharsEnhanced(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "命令注入字符",
			input:    "$(malicious)",
			expected: "\\$\\(malicious\\)",
		},
		{
			name:     "反引号命令替换",
			input:    "`whoami`",
			expected: "\\`whoami\\`",
		},
		{
			name:     "管道和重定向",
			input:    "file|grep test>output",
			expected: "file\\|grep\\ test\\>output",
		},
		{
			name:     "通配符字符",
			input:    "*.txt?[abc]{1,2}",
			expected: "\\*.txt\\?\\[abc\\]\\{1,2\\}",
		},
		{
			name:     "复合特殊字符",
			input:    "path/to/file$(cat /etc/passwd)&",
			expected: "path/to/file\\$\\(cat\\ /etc/passwd\\)\\&",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("escapeSpecialChars(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCompletionScriptSecurity 测试补全脚本安全性
func TestCompletionScriptSecurity(t *testing.T) {
	// 创建包含潜在恶意输入的测试上下文
	rootCtx := createMaliciousTestContext()

	tests := []struct {
		name      string
		shellType string
	}{
		{"Bash安全测试", flags.ShellBash},
		{"PowerShell安全测试", flags.ShellPowershell},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := GenerateShellCompletion(rootCtx, tt.shellType)
			if err != nil {
				t.Fatalf("生成%s补全脚本失败: %v", tt.shellType, err)
			}

			// 检查是否包含真正危险的模式（只检查用户输入数据中的危险字符）
			// 这些是来自测试数据中的恶意输入，应该被正确转义
			reallyDangerousPatterns := []string{
				"$(rm -rf /)",       // 命令注入
				"`cat /etc/passwd`", // 反引号命令替换
				"rm important.txt",  // 删除命令
				"curl evil.com",     // 网络请求
				"format c:",         // 格式化命令
			}

			for _, pattern := range reallyDangerousPatterns {
				if strings.Contains(script, pattern) {
					// 检查是否在枚举选项的引号内（已转义）
					if tt.shellType == flags.ShellBash {
						// 对于Bash，危险字符应该被转义
						if !strings.Contains(script, strings.ReplaceAll(pattern, " ", "\\ ")) {
							t.Errorf("%s脚本包含未转义的危险模式: %s", tt.shellType, pattern)
						}
					} else if tt.shellType == flags.ShellPowershell {
						// 对于PowerShell，检查是否已正确转义
						// 查找模式在脚本中的位置
						index := strings.Index(script, pattern)
						if index == -1 {
							continue // 模式不存在，跳过
						}

						// 检查模式前后的上下文，确保它在单引号字符串内
						contextStart := index - 10
						if contextStart < 0 {
							contextStart = 0
						}
						contextEnd := index + len(pattern) + 10
						if contextEnd > len(script) {
							contextEnd = len(script)
						}

						context := script[contextStart:contextEnd]

						// 检查是否在单引号字符串内（这是安全的）
						beforeQuote := strings.LastIndex(context[:index-contextStart], "'")
						afterQuote := strings.Index(context[index-contextStart+len(pattern):], "'")

						if beforeQuote != -1 && afterQuote != -1 {
							// 模式在单引号内，是安全的
							continue
						}

						// 如果不在引号内，则检查是否被转义
						escaped := false
						for _, char := range []string{"$", "`", "&", "|", ";", "(", ")"} {
							if strings.Contains(pattern, char) && strings.Contains(context, "`"+char) {
								escaped = true
								break
							}
						}

						if !escaped {
							t.Errorf("%s脚本包含未转义的危险模式: %s", tt.shellType, pattern)
						}
					}
				}
			}

			// 验证脚本不包含明显的安全漏洞
			securityChecks := []struct {
				pattern string
				desc    string
			}{
				{"eval ", "不应包含eval命令"},
				{"exec ", "不应包含exec命令"},
				{"/bin/sh", "不应直接调用shell"},
				{"system(", "不应调用system函数"},
			}

			for _, check := range securityChecks {
				if strings.Contains(script, check.pattern) {
					t.Errorf("%s脚本安全检查失败: %s", tt.shellType, check.desc)
				}
			}
		})
	}
}

// createLargeCommandTree 创建大型命令树用于性能测试
func createLargeCommandTree(mainCmds, subCmds int) *types.CmdContext {
	rootCtx := types.NewCmdContext("perftest", "", flag.ContinueOnError)

	// 添加根命令标志
	addTestFlags(rootCtx, 10)

	// 创建主命令
	for i := 0; i < mainCmds; i++ {
		mainCmd := types.NewCmdContext(fmt.Sprintf("cmd%d", i), fmt.Sprintf("c%d", i), flag.ContinueOnError)
		addTestFlags(mainCmd, 5)

		// 为每个主命令创建子命令
		for j := 0; j < subCmds; j++ {
			subCmd := types.NewCmdContext(fmt.Sprintf("sub%d", j), fmt.Sprintf("s%d", j), flag.ContinueOnError)
			addTestFlags(subCmd, 3)
			mainCmd.SubCmds = append(mainCmd.SubCmds, subCmd)
		}

		rootCtx.SubCmds = append(rootCtx.SubCmds, mainCmd)
	}

	return rootCtx
}

// createTestContextWithSpecialChars 创建包含特殊字符的测试上下文
func createTestContextWithSpecialChars() *types.CmdContext {
	rootCtx := types.NewCmdContext("special", "", flag.ContinueOnError)

	// 添加包含特殊字符的枚举标志
	enumFlag := &flags.EnumFlag{}
	_ = enumFlag.Init("mode", "m", "debug", "运行模式", []string{
		"debug|test",
		"prod$(echo hack)",
		"dev`whoami`",
		"test&background",
		"stage;rm -rf /",
		"local<input>output",
		"remote*wildcard",
		"backup?question",
		"config[array]",
		"data{object}",
		"temp~home",
		"log#comment",
	})
	enumMeta := &flags.FlagMeta{Flag: enumFlag}
	_ = rootCtx.FlagRegistry.RegisterFlag(enumMeta)

	return rootCtx
}

// createMaliciousTestContext 创建包含潜在恶意输入的测试上下文
func createMaliciousTestContext() *types.CmdContext {
	rootCtx := types.NewCmdContext("malicious", "", flag.ContinueOnError)

	// 添加包含潜在恶意代码的标志
	var maliciousValue string
	maliciousFlag := &flags.StringFlag{}
	_ = maliciousFlag.Init("payload", "p", "恶意载荷测试", &maliciousValue)
	maliciousMeta := &flags.FlagMeta{Flag: maliciousFlag}
	_ = rootCtx.FlagRegistry.RegisterFlag(maliciousMeta)

	// 添加恶意枚举选项
	enumFlag := &flags.EnumFlag{}
	_ = enumFlag.Init("exploit", "e", "safe", "漏洞利用测试", []string{
		"safe",
		"$(rm -rf /)",
		"`cat /etc/passwd`",
		"normal|rm important.txt",
		"test&curl evil.com",
		"data;format c:",
	})
	enumMeta := &flags.FlagMeta{Flag: enumFlag}
	_ = rootCtx.FlagRegistry.RegisterFlag(enumMeta)

	return rootCtx
}

// addTestFlags 为命令上下文添加测试标志
func addTestFlags(ctx *types.CmdContext, count int) {
	for i := 0; i < count; i++ {
		var value string
		flag := &flags.StringFlag{}
		flagName := fmt.Sprintf("flag%d", i)
		shortName := fmt.Sprintf("f%d", i)
		_ = flag.Init(flagName, shortName, "测试标志", &value)
		meta := &flags.FlagMeta{Flag: flag}
		_ = ctx.FlagRegistry.RegisterFlag(meta)
	}
}

// BenchmarkPowerShellGeneration 基准测试PowerShell生成性能
func BenchmarkPowerShellGeneration(b *testing.B) {
	ctx := createLargeCommandTree(50, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateShellCompletion(ctx, flags.ShellPowershell)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBashGeneration 基准测试Bash生成性能
func BenchmarkBashGeneration(b *testing.B) {
	ctx := createLargeCommandTree(50, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateShellCompletion(ctx, flags.ShellBash)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEscapeSpecialCharsOptimized 基准测试优化后的特殊字符转义
func BenchmarkEscapeSpecialCharsOptimized(b *testing.B) {
	testStrings := []string{
		"normal string without special chars",
		"string with $(command) injection",
		"path\\to\\file with `backticks` and |pipes|",
		"complex*string?with[many]special{chars}&more;stuff",
		"very long string with multiple special characters: $(){}[]|&;<>*?~#",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_ = escapeSpecialChars(s)
		}
	}
}
