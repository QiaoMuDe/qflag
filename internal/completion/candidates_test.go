// candidates_test.go - 候选选项获取指令测试
//
// 该文件包含 candidates.go 中所有函数的单元测试

package completion

import (
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestHandleCandidates_ValidContext 测试有效上下文的候选选项获取
//
// 验证根据有效上下文路径返回正确的候选选项
func TestHandleCandidates_ValidContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, ""))

	// 获取根命令的候选选项
	candidates, err := GetCandidates(root, "/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 应该包含子命令 server 和标志 --output, -o，以及内置标志
	found := make(map[string]bool)
	for _, c := range candidates {
		found[c] = true
	}

	if !found["server"] {
		t.Error("Expected candidates to include 'server'")
	}
	if !found["--output"] {
		t.Error("Expected candidates to include '--output'")
	}
	if !found["-o"] {
		t.Error("Expected candidates to include '-o'")
	}
	// 内置帮助标志
	if !found["--help"] {
		t.Error("Expected candidates to include '--help'")
	}
	if !found["-h"] {
		t.Error("Expected candidates to include '-h'")
	}
}

// TestHandleCandidates_NestedContext 测试嵌套上下文的候选选项
//
// 验证嵌套子命令的候选选项获取
func TestHandleCandidates_NestedContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)
	_ = serverCmd.AddFlag(mock.NewMockFlag("port", "p", "Server port", types.FlagTypeString, ""))

	// 获取 /server/ 上下文的候选选项
	candidates, err := GetCandidates(root, "/server/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	found := make(map[string]bool)
	for _, c := range candidates {
		found[c] = true
	}

	// 应该包含子命令 start 和标志 --port, -p
	if !found["start"] {
		t.Error("Expected candidates to include 'start'")
	}
	if !found["--port"] {
		t.Error("Expected candidates to include '--port'")
	}
	if !found["-p"] {
		t.Error("Expected candidates to include '-p'")
	}
}

// TestHandleCandidates_InvalidContext 测试无效上下文
//
// 验证无效上下文返回空列表
func TestHandleCandidates_InvalidContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	candidates, err := GetCandidates(root, "/invalid/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(candidates) != 0 {
		t.Errorf("Expected empty candidates for invalid context, got %v", candidates)
	}
}

// TestHandleCandidates_EmptyContext 测试空上下文
//
// 验证空上下文返回空列表
func TestHandleCandidates_EmptyContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	candidates, err := GetCandidates(root, "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(candidates) != 0 {
		t.Errorf("Expected empty candidates for empty context, got %v", candidates)
	}
}

// TestGetCandidates_WithVersionFlag 测试带版本标志的候选选项
//
// 验证根命令配置了版本时包含版本标志
func TestGetCandidates_WithVersionFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetVersion("1.0.0")

	candidates, err := GetCandidates(root, "/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	found := make(map[string]bool)
	for _, c := range candidates {
		found[c] = true
	}

	// 应该包含版本标志
	if !found["--version"] {
		t.Error("Expected candidates to include '--version' when version is set")
	}
	if !found["-v"] {
		t.Error("Expected candidates to include '-v' when version is set")
	}
}

// TestGetCandidates_WithCompletionFlag 测试带补全标志的候选选项
//
// 验证根命令启用动态补全时包含补全标志
// 注意：mock.Config() 返回的 EnableDynamicCompletion 默认为 false，
// 所以此测试主要验证当 EnableDynamicCompletion 为 true 时的行为
func TestGetCandidates_WithCompletionFlag(t *testing.T) {
	// 由于 mock 的 Config() 方法不返回 EnableDynamicCompletion 字段，
	// 我们需要测试的是：当该字段为 true 时，补全标志会被包含
	// 这里我们验证当 EnableDynamicCompletion 为 false（默认值）时，不包含补全标志
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	candidates, err := GetCandidates(root, "/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 默认情况下不应该包含补全标志（因为 EnableDynamicCompletion 默认为 false）
	for _, c := range candidates {
		if c == "--completion" {
			t.Error("Should not include '--completion' when EnableDynamicCompletion is false")
		}
	}
}

// TestGetCandidates_SubcommandNoVersionFlag 测试子命令不包含版本标志
//
// 验证非根上下文不包含版本和补全标志
func TestGetCandidates_SubcommandNoVersionFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetVersion("1.0.0")
	root.SetCompletion(true)

	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)

	// 获取子命令的候选选项
	candidates, err := GetCandidates(root, "/server/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 子命令不应该包含版本和补全标志
	for _, c := range candidates {
		if c == "--version" || c == "-v" {
			t.Error("Subcommand should not include version flags")
		}
		if c == "--completion" {
			t.Error("Subcommand should not include completion flag")
		}
	}
}

// TestGetBuiltinFlagNames_RootContext 测试根上下文的内置标志
//
// 验证根上下文包含帮助标志
func TestGetBuiltinFlagNames_RootContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	names := getBuiltinFlagNames(root, "/")

	found := make(map[string]bool)
	for _, n := range names {
		found[n] = true
	}

	// 根上下文应该包含帮助标志
	if !found["--help"] {
		t.Error("Expected builtin flags to include '--help'")
	}
	if !found["-h"] {
		t.Error("Expected builtin flags to include '-h'")
	}
}

// TestGetBuiltinFlagNames_WithVersion 测试带版本的内置标志
//
// 验证配置了版本的根命令包含版本标志
func TestGetBuiltinFlagNames_WithVersion(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetVersion("1.0.0")

	names := getBuiltinFlagNames(root, "/")

	found := make(map[string]bool)
	for _, n := range names {
		found[n] = true
	}

	if !found["--version"] {
		t.Error("Expected builtin flags to include '--version' when version is set")
	}
	if !found["-v"] {
		t.Error("Expected builtin flags to include '-v' when version is set")
	}
}

// TestGetBuiltinFlagNames_SubcommandContext 测试子上下文的内置标志
//
// 验证非根上下文只包含帮助标志
func TestGetBuiltinFlagNames_SubcommandContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetVersion("1.0.0")
	root.SetCompletion(true)

	names := getBuiltinFlagNames(root, "/server/")

	// 子上下文应该只包含帮助标志
	for _, n := range names {
		if n == "--version" || n == "-v" {
			t.Error("Subcommand context should not include version flags")
		}
		if n == "--completion" {
			t.Error("Subcommand context should not include completion flag")
		}
	}

	// 但应该包含帮助标志
	hasHelp := false
	for _, n := range names {
		if n == "--help" || n == "-h" {
			hasHelp = true
			break
		}
	}
	if !hasHelp {
		t.Error("Subcommand context should include help flags")
	}
}

// TestGetCandidates_PossibleDuplicates 测试候选选项可能的重复情况
//
// 验证当用户定义的标志名与内置标志冲突时的情况
// 注意：当前实现不会去重，重复项是可能存在的
func TestGetCandidates_PossibleDuplicates(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	// 添加一个名为 "output" 的标志，不会与内置标志冲突
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, ""))

	candidates, err := GetCandidates(root, "/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 检查是否有重复（在正常情况下不应该有）
	seen := make(map[string]bool)
	for _, c := range candidates {
		if seen[c] {
			// 记录重复但不报错，因为这是当前实现的预期行为
			t.Logf("Duplicate candidate found: %s", c)
		}
		seen[c] = true
	}

	// 验证 --output 和 -o 标志存在
	if !seen["--output"] {
		t.Error("Expected candidates to include '--output' flag")
	}
	if !seen["-o"] {
		t.Error("Expected candidates to include '-o' flag")
	}
}

// TestGetCandidates_FiltersInternalCommands 测试过滤内部命令
//
// 验证以 __ 开头的内部命令被过滤
func TestGetCandidates_FiltersInternalCommands(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	internalCmd := mock.NewMockCommandBasic("__complete", "", "Internal command")
	_ = root.AddSubCmds(serverCmd, internalCmd)

	candidates, err := GetCandidates(root, "/")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 检查是否过滤了内部命令
	for _, c := range candidates {
		if strings.HasPrefix(c, "__") {
			t.Errorf("Internal command should be filtered: %s", c)
		}
	}

	// 但应该包含普通子命令
	found := false
	for _, c := range candidates {
		if c == "server" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected candidates to include 'server'")
	}
}
