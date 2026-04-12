// context_test.go - 上下文计算指令测试
//
// 该文件包含 context.go 中所有函数的单元测试

package completion

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestCalculateContext_EmptyArgs 测试空参数情况
//
// 验证当只提供程序名时，上下文应为根目录 "/"
func TestCalculateContext_EmptyArgs(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	// 只有程序名，没有子命令
	tokens := []string{"myapp"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	if result.Context != "/" {
		t.Errorf("Expected context '/', got '%s'", result.Context)
	}
	if result.Depth != 0 {
		t.Errorf("Expected depth 0, got %d", result.Depth)
	}
	if result.CurrentCmd != "myapp" {
		t.Errorf("Expected current cmd 'myapp', got '%s'", result.CurrentCmd)
	}
}

// TestCalculateContext_SingleSubcommand 测试单级子命令
//
// 验证单级子命令的上下文计算
func TestCalculateContext_SingleSubcommand(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)

	tokens := []string{"myapp", "server"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	if result.Context != "/server/" {
		t.Errorf("Expected context '/server/', got '%s'", result.Context)
	}
	if result.Depth != 1 {
		t.Errorf("Expected depth 1, got %d", result.Depth)
	}
	if result.CurrentCmd != "server" {
		t.Errorf("Expected current cmd 'server', got '%s'", result.CurrentCmd)
	}
	if result.ParentContext != "/" {
		t.Errorf("Expected parent context '/', got '%s'", result.ParentContext)
	}
}

// TestCalculateContext_NestedSubcommands 测试嵌套子命令
//
// 验证多级嵌套子命令的上下文计算
func TestCalculateContext_NestedSubcommands(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)

	tokens := []string{"myapp", "server", "start"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	if result.Context != "/server/start/" {
		t.Errorf("Expected context '/server/start/', got '%s'", result.Context)
	}
	if result.Depth != 2 {
		t.Errorf("Expected depth 2, got %d", result.Depth)
	}
	if result.CurrentCmd != "start" {
		t.Errorf("Expected current cmd 'start', got '%s'", result.CurrentCmd)
	}
	if result.ParentContext != "/server/" {
		t.Errorf("Expected parent context '/server/', got '%s'", result.ParentContext)
	}
}

// TestCalculateContext_WithFlags 测试遇到标志时的行为
//
// 验证当遇到以 "-" 开头的标志时，上下文构建应停止
func TestCalculateContext_WithFlags(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)
	_ = serverCmd.AddFlag(mock.NewMockFlag("port", "p", "Server port", types.FlagTypeString, ""))

	tokens := []string{"myapp", "server", "--port", "8080"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	// 遇到 --port 后应该停止，上下文应为 /server/
	if result.Context != "/server/" {
		t.Errorf("Expected context '/server/', got '%s'", result.Context)
	}
	if !result.IsFlagContext {
		t.Error("Expected IsFlagContext to be true")
	}
	if result.FlagsStartIndex != 2 {
		t.Errorf("Expected FlagsStartIndex 2, got %d", result.FlagsStartIndex)
	}
}

// TestCalculateContext_FlagWithEquals 测试等号赋值标志
//
// 验证 --flag=value 格式的标志应被跳过，继续解析后续内容
func TestCalculateContext_FlagWithEquals(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)
	_ = serverCmd.AddFlag(mock.NewMockFlag("config", "c", "Config file", types.FlagTypeString, ""))

	tokens := []string{"myapp", "server", "--config=test.yaml", "start"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	// --config=test.yaml 应该被跳过，继续解析到 start
	if result.Context != "/server/start/" {
		t.Errorf("Expected context '/server/start/', got '%s'", result.Context)
	}
}

// TestCalculateContext_DoubleDashTerminator 测试双横线终止符
//
// 验证 -- 后的内容应继续解析为位置参数
func TestCalculateContext_DoubleDashTerminator(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)

	// 使用 -- 来停止标志解析，后面跟着一个普通的位置参数
	tokens := []string{"myapp", "server", "--", "start"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	// -- 后的 start 应该被当作位置参数（子命令）继续解析
	if result.Context != "/server/start/" {
		t.Errorf("Expected context '/server/start/', got '%s'", result.Context)
	}
	// -- 不是标志，start 也不是标志，所以 IsFlagContext 应该是 false
	if result.IsFlagContext {
		t.Error("Expected IsFlagContext to be false after --")
	}
}

// TestCalculateContext_InvalidSubcommand 测试无效子命令
//
// 验证遇到无效子命令时应停止在上一级上下文
func TestCalculateContext_InvalidSubcommand(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)

	tokens := []string{"myapp", "server", "invalid-cmd"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	// 遇到无效子命令应该停止在 /server/
	if result.Context != "/server/" {
		t.Errorf("Expected context '/server/', got '%s'", result.Context)
	}
	if result.Depth != 1 {
		t.Errorf("Expected depth 1, got %d", result.Depth)
	}
}

// TestCalculateContext_PartialInput 测试部分输入
//
// 验证光标位置在中间的计算
func TestCalculateContext_PartialInput(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)

	// 用户输入了 "myapp server"，正在输入 start
	tokens := []string{"myapp", "server", "star"}
	cursorPos := 2 // 光标在 server 后

	result := CalculateContext(root, tokens, cursorPos)

	if result.Context != "/server/" {
		t.Errorf("Expected context '/server/', got '%s'", result.Context)
	}
}

// TestGetSubCommandNames 测试获取子命令名称
//
// 验证子命令名称列表的获取和过滤（隐藏命令应被过滤）
func TestGetSubCommandNames(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	clientCmd := mock.NewMockCommandBasic("client", "", "Client management")
	internalCmd := mock.NewMockCommandBasic("__internal", "", "Internal command")
	internalCmd.SetHidden(true) // 设置为隐藏命令

	_ = root.AddSubCmds(serverCmd, clientCmd, internalCmd)

	names := getSubCommandNames(root)

	// 应该包含 server 和 client，但不包含 __internal
	if len(names) != 2 {
		t.Errorf("Expected 2 subcommands, got %d: %v", len(names), names)
	}

	hasServer := false
	hasClient := false
	for _, name := range names {
		if name == "server" {
			hasServer = true
		}
		if name == "client" {
			hasClient = true
		}
		if name == "__internal" {
			t.Error("Should not include __internal command")
		}
	}

	if !hasServer {
		t.Error("Should include 'server' command")
	}
	if !hasClient {
		t.Error("Should include 'client' command")
	}
}

// TestGetSubCommandNames_Empty 测试无子命令情况
//
// 验证没有子命令时返回空切片
func TestGetSubCommandNames_Empty(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	names := getSubCommandNames(root)

	if len(names) != 0 {
		t.Errorf("Expected 0 subcommands, got %d", len(names))
	}
}

// TestGetFlagNames 测试获取标志名称
//
// 验证标志名称列表包含长短名称
func TestGetFlagNames(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, ""))
	_ = root.AddFlag(mock.NewMockFlag("count", "c", "Count", types.FlagTypeInt, 0))
	_ = root.AddFlag(mock.NewMockFlag("verbose", "", "Verbose", types.FlagTypeBool, false)) // 无短名称

	names := getFlagNames(root)

	// 应该有: --output, -o, --count, -c, --verbose (5个)
	if len(names) != 5 {
		t.Errorf("Expected 5 flag names, got %d: %v", len(names), names)
	}

	expected := map[string]bool{
		"--output":  false,
		"-o":        false,
		"--count":   false,
		"-c":        false,
		"--verbose": false,
	}

	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("Expected flag '%s' not found in %v", name, names)
		}
	}
}

// TestGetFlagNames_Empty 测试无标志情况
//
// 验证没有标志时返回空切片
func TestGetFlagNames_Empty(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	names := getFlagNames(root)

	if len(names) != 0 {
		t.Errorf("Expected 0 flag names, got %d", len(names))
	}
}

// TestFindCommandByContext_Root 测试查找根命令
//
// 验证根上下文返回根命令
func TestFindCommandByContext_Root(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	found := findCommandByContext(root, "/")

	if found == nil {
		t.Error("Expected to find root command, got nil")
	}
	if found.Name() != "myapp" {
		t.Errorf("Expected command 'myapp', got '%s'", found.Name())
	}
}

// TestFindCommandByContext_SingleLevel 测试查找单级子命令
//
// 验证单级上下文的命令查找
func TestFindCommandByContext_SingleLevel(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)

	found := findCommandByContext(root, "/server/")

	if found == nil {
		t.Error("Expected to find server command, got nil")
	}
	if found.Name() != "server" {
		t.Errorf("Expected command 'server', got '%s'", found.Name())
	}
}

// TestFindCommandByContext_Nested 测试查找嵌套子命令
//
// 验证多级嵌套上下文的命令查找
func TestFindCommandByContext_Nested(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)

	found := findCommandByContext(root, "/server/start/")

	if found == nil {
		t.Error("Expected to find start command, got nil")
	}
	if found.Name() != "start" {
		t.Errorf("Expected command 'start', got '%s'", found.Name())
	}
}

// TestFindCommandByContext_Invalid 测试无效上下文
//
// 验证无效上下文返回 nil
func TestFindCommandByContext_Invalid(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	found := findCommandByContext(root, "/invalid/")

	if found != nil {
		t.Errorf("Expected nil for invalid context, got '%s'", found.Name())
	}
}

// TestFindCommandByContext_Empty 测试空上下文
//
// 验证空字符串上下文返回 nil
func TestFindCommandByContext_Empty(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	found := findCommandByContext(root, "")

	if found != nil {
		t.Errorf("Expected nil for empty context, got '%s'", found.Name())
	}
}

// TestContextResult_SubCommands 测试上下文中获取子命令列表
//
// 验证 CalculateContext 返回的子命令列表
func TestContextResult_SubCommands(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	clientCmd := mock.NewMockCommandBasic("client", "", "Client management")
	_ = root.AddSubCmds(serverCmd, clientCmd)

	tokens := []string{"myapp"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	if len(result.SubCommands) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(result.SubCommands))
	}
}

// TestContextResult_Flags 测试上下文中获取标志列表
//
// 验证 CalculateContext 返回的标志列表
func TestContextResult_Flags(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, ""))
	_ = root.AddFlag(mock.NewMockFlag("verbose", "v", "Verbose", types.FlagTypeBool, false))

	tokens := []string{"myapp"}
	cursorPos := len(tokens)

	result := CalculateContext(root, tokens, cursorPos)

	// 应该有 4 个: --output, -o, --verbose, -v
	if len(result.Flags) != 4 {
		t.Errorf("Expected 4 flags, got %d: %v", len(result.Flags), result.Flags)
	}
}
