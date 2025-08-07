// Package completion 转义性能测试
// 本文件包含了补全系统中字符转义功能的性能基准测试，
// 测试特殊字符处理的性能表现，为优化提供数据支持。
package completion

import (
	"testing"
)

// BenchmarkBashEscapeSpecialChars 基准测试Bash转义函数性能
func BenchmarkBashEscapeSpecialChars(b *testing.B) {
	testStrings := []string{
		"normal string without special chars",
		"string with $(command) injection",
		"path\\to\\file with `backticks` and |pipes|",
		"complex*string?with[many]special{chars}&more;stuff",
		"very long string with multiple special characters: $(){}[]|&;<>*?~#",
		"mixed 'quotes' and \"double quotes\" with spaces",
		"file.txt > output.log 2>&1",
		"grep 'pattern' file.txt | sort | uniq",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_ = escapeSpecialChars(s)
		}
	}
}

// BenchmarkPwshEscapeString 基准测试PowerShell转义函数性能
func BenchmarkPwshEscapeString(b *testing.B) {
	testStrings := []string{
		"normal string without special chars",
		"path\\to\\file with 'quotes' and spaces",
		"simple",
		"complex'string\\with\\many'special\\chars",
		"PowerShell $variables and `backticks`",
		"file & process | pipeline ; commands",
		"redirect < input > output",
		"parentheses (grouping) and \"quotes\"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			_ = escapePwshString(s)
		}
	}
}

// BenchmarkBashEscapeVaryingLengths 测试不同长度字符串的Bash转义性能
func BenchmarkBashEscapeVaryingLengths(b *testing.B) {
	testCases := []struct {
		name string
		str  string
	}{
		{"短字符串", "test$var"},
		{"中等字符串", "grep 'pattern' file.txt | sort | uniq > output.log"},
		{"长字符串", "very long command with many special characters: $(){}[]|&;<>*?~# and more content to test performance with longer strings that contain multiple special characters"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = escapeSpecialChars(tc.str)
			}
		})
	}
}

// BenchmarkPwshEscapeVaryingLengths 测试不同长度字符串的PowerShell转义性能
func BenchmarkPwshEscapeVaryingLengths(b *testing.B) {
	testCases := []struct {
		name string
		str  string
	}{
		{"短字符串", "test$var"},
		{"中等字符串", "Get-Process | Where-Object {$_.Name -like 'chrome*'}"},
		{"长字符串", "very long PowerShell command with many special characters: $variables `backticks` 'quotes' \"double quotes\" & | ; < > ( ) and more content to test performance"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = escapePwshString(tc.str)
			}
		})
	}
}

// BenchmarkBashEscapeSpecialCharDensity 测试不同特殊字符密度的Bash转义性能
func BenchmarkBashEscapeSpecialCharDensity(b *testing.B) {
	testCases := []struct {
		name string
		str  string
	}{
		{"无特殊字符", "normalstringwithoutspecialchars"},
		{"低密度", "normal string with few $ special chars"},
		{"中密度", "string with $more |special &chars; and *wildcards?"},
		{"高密度", "$|&;*?[]{}~#<>()\\\"` all special chars"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = escapeSpecialChars(tc.str)
			}
		})
	}
}

// BenchmarkPwshEscapeSpecialCharDensity 测试不同特殊字符密度的PowerShell转义性能
func BenchmarkPwshEscapeSpecialCharDensity(b *testing.B) {
	testCases := []struct {
		name string
		str  string
	}{
		{"无特殊字符", "normalstringwithoutspecialchars"},
		{"低密度", "normal string with few $ special chars"},
		{"中密度", "string with $more |special &chars; and 'quotes'"},
		{"高密度", "$`'\"&|;<>()\\all special chars"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = escapePwshString(tc.str)
			}
		})
	}
}

// TestBashEscapeMapCompleteness 测试Bash转义映射表的完整性
func TestBashEscapeMapCompleteness(t *testing.T) {
	// 验证所有预期的特殊字符都在映射表中
	expectedChars := []rune{'\\', '"', ' ', '$', '`', '|', '&', ';', '(', ')', '<', '>', '*', '?', '[', ']', '{', '}', '~', '#'}

	for _, char := range expectedChars {
		if _, exists := bashEscapeMap[char]; !exists {
			t.Errorf("Bash转义映射表缺少字符: %c", char)
		}
	}

	t.Logf("Bash转义映射表包含 %d 个特殊字符", len(bashEscapeMap))
}

// TestPwshEscapeMapCompleteness 测试PowerShell转义映射表的完整性
func TestPwshEscapeMapCompleteness(t *testing.T) {
	// 验证所有预期的特殊字符都在映射表中
	expectedChars := []byte{'\'', '\\', '$', '`', '"', '&', '|', ';', '<', '>', '(', ')'}

	for _, char := range expectedChars {
		if _, exists := pwshEscapeMap[char]; !exists {
			t.Errorf("PowerShell转义映射表缺少字符: %c", char)
		}
	}

	t.Logf("PowerShell转义映射表包含 %d 个特殊字符", len(pwshEscapeMap))
}

// TestEscapeFunctionConsistency 测试转义函数的一致性
func TestEscapeFunctionConsistency(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		bashExpected string
		pwshExpected string
	}{
		{
			name:         "普通字符串",
			input:        "normal",
			bashExpected: "normal",
			pwshExpected: "normal",
		},
		{
			name:         "包含美元符号",
			input:        "test$var",
			bashExpected: "test\\$var",
			pwshExpected: "test`$var",
		},
		{
			name:         "包含反斜杠",
			input:        "path\\file",
			bashExpected: "path\\\\file",
			pwshExpected: "path\\\\file",
		},
		{
			name:         "包含管道符",
			input:        "cmd|grep",
			bashExpected: "cmd\\|grep",
			pwshExpected: "cmd`|grep",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bashResult := escapeSpecialChars(tc.input)
			pwshResult := escapePwshString(tc.input)

			if bashResult != tc.bashExpected {
				t.Errorf("Bash转义结果不匹配: 输入=%q, 期望=%q, 实际=%q", tc.input, tc.bashExpected, bashResult)
			}

			if pwshResult != tc.pwshExpected {
				t.Errorf("PowerShell转义结果不匹配: 输入=%q, 期望=%q, 实际=%q", tc.input, tc.pwshExpected, pwshResult)
			}
		})
	}
}
