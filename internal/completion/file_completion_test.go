// Package completion 文件补全测试
// 本文件包含了文件路径自动补全功能的单元测试，测试文件和
// 目录路径的智能补全功能的正确性和稳定性。
package completion

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestBashFilePathCompletion 测试Bash文件路径补全功能
func TestBashFilePathCompletion(t *testing.T) {
	// 创建包含字符串标志的测试上下文
	ctx := createTestContextWithStringFlags()

	script, err := GenerateShellCompletion(ctx, flags.ShellBash)
	if err != nil {
		t.Fatalf("生成Bash补全脚本失败: %v", err)
	}

	// 验证字符串类型的文件路径补全逻辑
	expectedBashFeatures := []string{
		"string)",
		"COMPREPLY=($(compgen -f -d -- \"$cur\"))",
	}

	for _, feature := range expectedBashFeatures {
		if !strings.Contains(script, feature) {
			t.Errorf("Bash补全脚本缺少文件路径补全功能: %s", feature)
		}
	}

	t.Log("✅ Bash文件路径补全功能验证通过")
}

// TestPowerShellFilePathCompletion 测试PowerShell文件路径补全功能
func TestPowerShellFilePathCompletion(t *testing.T) {
	// 创建包含字符串标志的测试上下文
	ctx := createTestContextWithStringFlags()

	script, err := GenerateShellCompletion(ctx, flags.ShellPowershell)
	if err != nil {
		t.Fatalf("生成PowerShell补全脚本失败: %v", err)
	}

	// 验证字符串类型的文件路径补全逻辑
	expectedPwshFeatures := []string{
		"'string' {",
		"Get-ChildItem",
		"Split-Path",
		"Join-Path",
		"PSIsContainer",
	}

	for _, feature := range expectedPwshFeatures {
		if !strings.Contains(script, feature) {
			t.Errorf("PowerShell补全脚本缺少文件路径补全功能: %s", feature)
		}
	}

	t.Log("✅ PowerShell文件路径补全功能验证通过")
}

// TestFileCompletionWithDifferentFlagTypes 测试不同标志类型的补全行为
func TestFileCompletionWithDifferentFlagTypes(t *testing.T) {
	ctx := createMixedFlagTypesContext()

	tests := []struct {
		name      string
		shellType string
	}{
		{"Bash混合类型测试", flags.ShellBash},
		{"PowerShell混合类型测试", flags.ShellPowershell},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := GenerateShellCompletion(ctx, tt.shellType)
			if err != nil {
				t.Fatalf("生成%s补全脚本失败: %v", tt.shellType, err)
			}

			// 验证不同类型的处理
			switch tt.shellType {
			case flags.ShellBash:
				// Bash应该包含枚举和字符串的不同处理
				if !strings.Contains(script, "enum)") {
					t.Error("Bash脚本缺少枚举类型处理")
				}
				if !strings.Contains(script, "string)") {
					t.Error("Bash脚本缺少字符串类型处理")
				}
				if !strings.Contains(script, "compgen -f") {
					t.Error("Bash脚本缺少文件补全命令")
				}

			case flags.ShellPowershell:
				// PowerShell应该包含枚举和字符串的不同处理
				if !strings.Contains(script, "'enum' {") {
					t.Error("PowerShell脚本缺少枚举类型处理")
				}
				if !strings.Contains(script, "'string' {") {
					t.Error("PowerShell脚本缺少字符串类型处理")
				}
				if !strings.Contains(script, "Get-ChildItem") {
					t.Error("PowerShell脚本缺少文件获取命令")
				}
			}
		})
	}
}

// TestFileCompletionSecurity 测试文件补全的安全性
func TestFileCompletionSecurity(t *testing.T) {
	ctx := createTestContextWithStringFlags()

	tests := []struct {
		name      string
		shellType string
	}{
		{"Bash安全测试", flags.ShellBash},
		{"PowerShell安全测试", flags.ShellPowershell},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := GenerateShellCompletion(ctx, tt.shellType)
			if err != nil {
				t.Fatalf("生成%s补全脚本失败: %v", tt.shellType, err)
			}

			// 检查是否有适当的错误处理
			switch tt.shellType {
			case flags.ShellBash:
				// Bash应该使用安全的compgen命令
				if strings.Contains(script, "eval") || strings.Contains(script, "exec") {
					t.Error("Bash脚本包含不安全的命令执行")
				}

			case flags.ShellPowershell:
				// PowerShell应该有错误处理
				if !strings.Contains(script, "-ErrorAction SilentlyContinue") {
					t.Error("PowerShell脚本缺少错误处理")
				}
				if !strings.Contains(script, "try {") || !strings.Contains(script, "catch {") {
					t.Error("PowerShell脚本缺少异常处理")
				}
			}
		})
	}
}

// TestFileCompletionPerformance 测试文件补全的性能
func TestFileCompletionPerformance(t *testing.T) {
	// 创建临时目录结构用于测试
	tempDir := createTempDirStructure(t)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// 切换到临时目录
	oldDir, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(oldDir)
	}()
	_ = os.Chdir(tempDir)

	ctx := createTestContextWithStringFlags()

	tests := []struct {
		name      string
		shellType string
	}{
		{"Bash性能测试", flags.ShellBash},
		{"PowerShell性能测试", flags.ShellPowershell},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script, err := GenerateShellCompletion(ctx, tt.shellType)
			if err != nil {
				t.Fatalf("生成%s补全脚本失败: %v", tt.shellType, err)
			}

			// 验证脚本大小合理（不应该过大）
			if len(script) > 50000 { // 50KB限制
				t.Errorf("%s补全脚本过大: %d 字节", tt.shellType, len(script))
			}

			// 验证没有明显的性能问题模式
			switch tt.shellType {
			case flags.ShellBash:
				// 检查是否避免了低效的外部命令调用
				inefficientPatterns := []string{"find /", "ls -la", "grep -r"}
				for _, pattern := range inefficientPatterns {
					if strings.Contains(script, pattern) {
						t.Errorf("Bash脚本包含低效模式: %s", pattern)
					}
				}

			case flags.ShellPowershell:
				// 检查是否使用了高效的PowerShell命令
				if strings.Contains(script, "Get-ChildItem -Recurse") {
					t.Error("PowerShell脚本使用了低效的递归搜索")
				}
			}
		})
	}
}

// createTestContextWithStringFlags 创建包含字符串标志的测试上下文
func createTestContextWithStringFlags() *types.CmdContext {
	ctx := types.NewCmdContext("filetest", "", flag.ContinueOnError)

	// 添加各种字符串标志
	stringFlags := []struct {
		long, short, desc string
	}{
		{"config", "c", "配置文件路径"},
		{"output", "o", "输出文件路径"},
		{"input", "i", "输入文件路径"},
		{"directory", "d", "目录路径"},
		{"log", "l", "日志文件路径"},
	}

	for _, sf := range stringFlags {
		var value string
		flag := &flags.StringFlag{}
		_ = flag.Init(sf.long, sf.short, sf.desc, &value)
		meta := &flags.FlagMeta{Flag: flag}
		_ = ctx.FlagRegistry.RegisterFlag(meta)
	}

	return ctx
}

// createMixedFlagTypesContext 创建包含混合标志类型的测试上下文
func createMixedFlagTypesContext() *types.CmdContext {
	ctx := types.NewCmdContext("mixedtest", "", flag.ContinueOnError)

	// 添加字符串标志
	var stringValue string
	stringFlag := &flags.StringFlag{}
	_ = stringFlag.Init("file", "f", "文件路径", &stringValue)
	stringMeta := &flags.FlagMeta{Flag: stringFlag}
	_ = ctx.FlagRegistry.RegisterFlag(stringMeta)

	// 添加枚举标志
	enumFlag := &flags.EnumFlag{}
	_ = enumFlag.Init("format", "fmt", "json", "输出格式", []string{"json", "yaml", "xml", "csv"})
	enumMeta := &flags.FlagMeta{Flag: enumFlag}
	_ = ctx.FlagRegistry.RegisterFlag(enumMeta)

	// 添加布尔标志
	var boolValue bool
	boolFlag := &flags.BoolFlag{}
	_ = boolFlag.Init("verbose", "v", "详细输出", &boolValue)
	boolMeta := &flags.FlagMeta{Flag: boolFlag}
	_ = ctx.FlagRegistry.RegisterFlag(boolMeta)

	// 添加整数标志
	var intValue int
	intFlag := &flags.IntFlag{}
	_ = intFlag.Init("port", "p", "端口号", &intValue)
	intMeta := &flags.FlagMeta{Flag: intFlag}
	_ = ctx.FlagRegistry.RegisterFlag(intMeta)

	return ctx
}

// createTempDirStructure 创建临时目录结构用于测试
func createTempDirStructure(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "completion_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}

	// 创建测试文件和目录
	testStructure := []string{
		"config.yaml",
		"config.json",
		"data.txt",
		"logs/app.log",
		"logs/error.log",
		"docs/readme.md",
		"docs/api.md",
		"bin/app",
		"bin/tool",
	}

	for _, path := range testStructure {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)

		// 创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败 %s: %v", dir, err)
		}

		// 创建文件
		if filepath.Ext(path) != "" || !strings.Contains(path, "/") {
			file, err := os.Create(fullPath)
			if err != nil {
				t.Fatalf("创建文件失败 %s: %v", fullPath, err)
			}
			_ = file.Close()
		}
	}

	return tempDir
}

// BenchmarkFileCompletionGeneration 基准测试文件补全生成
func BenchmarkFileCompletionGeneration(b *testing.B) {
	ctx := createTestContextWithStringFlags()

	b.Run("Bash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := GenerateShellCompletion(ctx, flags.ShellBash)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("PowerShell", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := GenerateShellCompletion(ctx, flags.ShellPowershell)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
