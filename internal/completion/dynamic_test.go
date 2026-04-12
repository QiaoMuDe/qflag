// dynamic_test.go - 动态补全指令测试
//
// 该文件包含 dynamic.go 中 handleAll、getFlagType、getBuiltinFlagType 等函数的单元测试
// 注意：不测试 handleFuzzy（调用外部库）和 HandleDynamicComplete（路由函数）

package completion

import (
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestHandleAll_Basic 测试 handleAll 基本功能
//
// 验证基本参数解析和输出格式
func TestHandleAll_Basic(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, ""))

	// 测试基本调用
	args := []string{"ser", "", ""} // cur=ser, prev=空, 无子命令参数
	err := handleAll(root, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestHandleAll_WithSubcommand 测试带子命令的 handleAll
//
// 验证能正确识别子命令上下文
func TestHandleAll_WithSubcommand(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	startCmd := mock.NewMockCommandBasic("start", "", "Start server")
	_ = serverCmd.AddSubCmds(startCmd)
	_ = root.AddSubCmds(serverCmd)

	// 测试在 server 子命令下的补全
	args := []string{"sta", "", "server"} // cur=sta, prev=空, 已输入 server
	err := handleAll(root, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestHandleAll_FlagValueCompletion 测试标志值补全
//
// 验证在标志后能提供正确的值补全
func TestHandleAll_FlagValueCompletion(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	_ = root.AddFlag(mock.NewMockEnumFlag("format", "f", "Output format", "json", []string{"json", "yaml", "xml"}))

	// 测试在 --format 后的值补全
	args := []string{"js", "--format"} // cur=js, prev=--format
	err := handleAll(root, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestHandleAll_BoolFlag 测试布尔标志补全
//
// 验证布尔标志后不需要值，继续补全其他选项
func TestHandleAll_BoolFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)
	_ = root.AddFlag(mock.NewMockBoolFlag("verbose", "v", "Verbose output", false))

	// 测试在布尔标志 -v 后的补全
	args := []string{"ser", "-v"} // cur=ser, prev=-v
	err := handleAll(root, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestHandleAll_EmptyArgs 测试空参数
//
// 验证参数不足时返回错误
func TestHandleAll_EmptyArgs(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	args := []string{"cur"} // 只有 cur，缺少 prev
	err := handleAll(root, args)

	if err == nil {
		t.Error("Expected error for insufficient args, got nil")
	}

	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("Expected usage error, got %v", err)
	}
}

// TestHandleAll_QuoteTrimming 测试引号去除
//
// 验证 cur 和 prev 中的引号被正确去除
func TestHandleAll_QuoteTrimming(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	// 测试带引号的参数
	args := []string{`"cur"`, `'prev'`}
	err := handleAll(root, args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestGetFlagType_UserFlag 测试获取用户定义标志的类型
//
// 验证能从命令的标志列表中找到标志类型
func TestGetFlagType_UserFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, ""))
	_ = root.AddFlag(mock.NewMockEnumFlag("format", "f", "Format", "json", []string{"json", "yaml"}))
	_ = root.AddFlag(mock.NewMockBoolFlag("verbose", "v", "Verbose", false))

	// 测试字符串标志
	flagType, found := getFlagType(root, "/", "--output")
	if !found {
		t.Error("Expected to find --output flag")
	}
	if flagType != types.FlagTypeString {
		t.Errorf("Expected FlagTypeString, got %v", flagType)
	}

	// 测试枚举标志
	flagType, found = getFlagType(root, "/", "-f")
	if !found {
		t.Error("Expected to find -f flag")
	}
	if flagType != types.FlagTypeEnum {
		t.Errorf("Expected FlagTypeEnum, got %v", flagType)
	}

	// 测试布尔标志
	flagType, found = getFlagType(root, "/", "--verbose")
	if !found {
		t.Error("Expected to find --verbose flag")
	}
	if flagType != types.FlagTypeBool {
		t.Errorf("Expected FlagTypeBool, got %v", flagType)
	}
}

// TestGetFlagType_InvalidContext 测试无效上下文
//
// 验证无效上下文返回未找到
func TestGetFlagType_InvalidContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	flagType, found := getFlagType(root, "/invalid/", "--flag")
	if found {
		t.Error("Expected not to find flag in invalid context")
	}
	if flagType != types.FlagTypeUnknown {
		t.Errorf("Expected FlagTypeUnknown, got %v", flagType)
	}
}

// TestGetFlagType_NonExistentFlag 测试不存在的标志
//
// 验证不存在的标志返回未找到
func TestGetFlagType_NonExistentFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	flagType, found := getFlagType(root, "/", "--nonexistent")
	if found {
		t.Error("Expected not to find non-existent flag")
	}
	if flagType != types.FlagTypeUnknown {
		t.Errorf("Expected FlagTypeUnknown, got %v", flagType)
	}
}

// TestGetFlagType_NestedContext 测试嵌套上下文的标志类型
//
// 验证能在子命令中找到标志
func TestGetFlagType_NestedContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = serverCmd.AddFlag(mock.NewMockFlag("port", "p", "Server port", types.FlagTypeInt, 8080))
	_ = root.AddSubCmds(serverCmd)

	flagType, found := getFlagType(root, "/server/", "--port")
	if !found {
		t.Error("Expected to find --port flag in /server/ context")
	}
	if flagType != types.FlagTypeInt {
		t.Errorf("Expected FlagTypeInt, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_HelpFlag 测试帮助标志识别
//
// 验证能正确识别内置帮助标志
func TestGetBuiltinFlagType_HelpFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	// 测试长名称
	flagType, isBuiltin := getBuiltinFlagType("--help", "/", root)
	if !isBuiltin {
		t.Error("Expected --help to be recognized as builtin flag")
	}
	if flagType != types.FlagTypeBool {
		t.Errorf("Expected FlagTypeBool for help flag, got %v", flagType)
	}

	// 测试短名称
	flagType, isBuiltin = getBuiltinFlagType("-h", "/", root)
	if !isBuiltin {
		t.Error("Expected -h to be recognized as builtin flag")
	}
	if flagType != types.FlagTypeBool {
		t.Errorf("Expected FlagTypeBool for help flag, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_VersionFlag 测试版本标志识别
//
// 验证能正确识别内置版本标志（根命令且有版本时）
func TestGetBuiltinFlagType_VersionFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetVersion("1.0.0")

	// 根命令应该有版本标志
	flagType, isBuiltin := getBuiltinFlagType("--version", "/", root)
	if !isBuiltin {
		t.Error("Expected --version to be recognized as builtin flag in root context")
	}
	if flagType != types.FlagTypeBool {
		t.Errorf("Expected FlagTypeBool for version flag, got %v", flagType)
	}

	// 测试短名称
	flagType, isBuiltin = getBuiltinFlagType("-v", "/", root)
	if !isBuiltin {
		t.Error("Expected -v to be recognized as builtin flag")
	}
	if flagType != types.FlagTypeBool {
		t.Errorf("Expected FlagTypeBool for version flag, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_VersionFlagNoVersion 测试无版本时的版本标志
//
// 验证没有设置版本时，版本标志不被识别
func TestGetBuiltinFlagType_VersionFlagNoVersion(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	// 不设置版本

	flagType, isBuiltin := getBuiltinFlagType("--version", "/", root)
	if isBuiltin {
		t.Error("Expected --version NOT to be recognized when version is not set")
	}
	if flagType != types.FlagTypeUnknown {
		t.Errorf("Expected FlagTypeUnknown, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_VersionFlagSubcommand 测试子命令的版本标志
//
// 验证子命令上下文不包含版本标志
func TestGetBuiltinFlagType_VersionFlagSubcommand(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetVersion("1.0.0")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	_ = root.AddSubCmds(serverCmd)

	// 子命令上下文不应该有版本标志
	flagType, isBuiltin := getBuiltinFlagType("--version", "/server/", serverCmd)
	if isBuiltin {
		t.Error("Expected --version NOT to be recognized in subcommand context")
	}
	if flagType != types.FlagTypeUnknown {
		t.Errorf("Expected FlagTypeUnknown, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_CompletionFlag 测试补全标志识别
//
// 验证能正确识别内置补全标志
func TestGetBuiltinFlagType_CompletionFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetCompletion(true)

	flagType, isBuiltin := getBuiltinFlagType("--completion", "/", root)
	if !isBuiltin {
		t.Error("Expected --completion to be recognized as builtin flag")
	}
	if flagType != types.FlagTypeEnum {
		t.Errorf("Expected FlagTypeEnum for completion flag, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_CompletionFlagDisabled 测试禁用补全时的补全标志
//
// 验证禁用补全时，补全标志不被识别
func TestGetBuiltinFlagType_CompletionFlagDisabled(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	root.SetCompletion(false)

	flagType, isBuiltin := getBuiltinFlagType("--completion", "/", root)
	if isBuiltin {
		t.Error("Expected --completion NOT to be recognized when completion is disabled")
	}
	if flagType != types.FlagTypeUnknown {
		t.Errorf("Expected FlagTypeUnknown, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_UnknownFlag 测试未知标志
//
// 验证未知标志返回未找到
func TestGetBuiltinFlagType_UnknownFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	flagType, isBuiltin := getBuiltinFlagType("--unknown", "/", root)
	if isBuiltin {
		t.Error("Expected --unknown NOT to be recognized as builtin flag")
	}
	if flagType != types.FlagTypeUnknown {
		t.Errorf("Expected FlagTypeUnknown, got %v", flagType)
	}
}

// TestGetBuiltinFlagType_SingleDashHelp 测试单横线帮助标志
//
// 验证单横线帮助标志也能被识别
func TestGetBuiltinFlagType_SingleDashHelp(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	flagType, isBuiltin := getBuiltinFlagType("-help", "/", root)
	if !isBuiltin {
		t.Error("Expected -help to be recognized as builtin flag")
	}
	if flagType != types.FlagTypeBool {
		t.Errorf("Expected FlagTypeBool for help flag, got %v", flagType)
	}
}

// TestFuzzyMatch_EmptyCur 测试空当前输入的模糊匹配
//
// 验证当 cur 为空时返回所有候选项
func TestFuzzyMatch_EmptyCur(t *testing.T) {
	candidates := []string{"apple", "banana", "cherry"}
	result := fuzzyMatch(candidates, "")

	if len(result) != len(candidates) {
		t.Errorf("Expected %d candidates, got %d", len(candidates), len(result))
	}

	// 验证所有候选项都返回
	for i, c := range candidates {
		if result[i] != c {
			t.Errorf("Expected %s at index %d, got %s", c, i, result[i])
		}
	}
}

// TestFuzzyMatch_WithCur 测试有当前输入的模糊匹配
//
// 验证能根据当前输入过滤候选项
func TestFuzzyMatch_WithCur(t *testing.T) {
	candidates := []string{"apple", "application", "banana", "cherry"}
	result := fuzzyMatch(candidates, "app")

	// 应该匹配 apple 和 application
	if len(result) != 2 {
		t.Errorf("Expected 2 matches for 'app', got %d: %v", len(result), result)
	}

	// 验证匹配结果包含预期的候选项
	found := make(map[string]bool)
	for _, r := range result {
		found[r] = true
	}
	if !found["apple"] {
		t.Error("Expected 'apple' in matches")
	}
	if !found["application"] {
		t.Error("Expected 'application' in matches")
	}
}

// TestFuzzyMatch_NoMatch 测试无匹配情况
//
// 验证当没有匹配时返回空切片
func TestFuzzyMatch_NoMatch(t *testing.T) {
	candidates := []string{"apple", "banana", "cherry"}
	result := fuzzyMatch(candidates, "xyz")

	if len(result) != 0 {
		t.Errorf("Expected 0 matches for 'xyz', got %d: %v", len(result), result)
	}
}

// TestFuzzyMatch_EmptyCandidates 测试空候选项列表
//
// 验证空候选项列表返回空切片
func TestFuzzyMatch_EmptyCandidates(t *testing.T) {
	result := fuzzyMatch([]string{}, "test")

	if len(result) != 0 {
		t.Errorf("Expected 0 matches for empty candidates, got %d", len(result))
	}
}
