// enum_test.go - 枚举值补全指令测试
//
// 该文件包含 enum.go 中所有函数的单元测试

package completion

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestHandleEnum_ValidEnumFlag 测试有效枚举标志
//
// 验证根据有效上下文和枚举标志名返回正确的枚举值
func TestHandleEnum_ValidEnumFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	levelFlag := mock.NewMockEnumFlag("log-level", "l", "Log level", "info", []string{"debug", "info", "warn", "error"})
	_ = root.AddFlag(levelFlag)

	// 获取枚举值
	values, err := GetEnumValues(root, "/", "--log-level")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 4 {
		t.Errorf("Expected 4 enum values, got %d: %v", len(values), values)
	}

	// 验证所有枚举值都存在
	expected := map[string]bool{
		"debug": false,
		"info":  false,
		"warn":  false,
		"error": false,
	}
	for _, v := range values {
		if _, ok := expected[v]; ok {
			expected[v] = true
		}
	}
	for v, found := range expected {
		if !found {
			t.Errorf("Expected enum value '%s' not found", v)
		}
	}
}

// TestHandleEnum_ShortFlagName 测试短标志名
//
// 验证使用短标志名也能获取枚举值
func TestHandleEnum_ShortFlagName(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	levelFlag := mock.NewMockEnumFlag("log-level", "l", "Log level", "info", []string{"debug", "info", "warn", "error"})
	_ = root.AddFlag(levelFlag)

	// 使用短标志名获取枚举值
	values, err := GetEnumValues(root, "/", "-l")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 4 {
		t.Errorf("Expected 4 enum values, got %d: %v", len(values), values)
	}
}

// TestHandleEnum_InvalidContext 测试无效上下文
//
// 验证无效上下文返回空列表
func TestHandleEnum_InvalidContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	levelFlag := mock.NewMockEnumFlag("log-level", "l", "Log level", "info", []string{"debug", "info", "warn", "error"})
	_ = root.AddFlag(levelFlag)

	values, err := GetEnumValues(root, "/invalid/", "--log-level")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected empty values for invalid context, got %v", values)
	}
}

// TestHandleEnum_InvalidFlagName 测试无效标志名
//
// 验证无效标志名返回空列表
func TestHandleEnum_InvalidFlagName(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	levelFlag := mock.NewMockEnumFlag("log-level", "l", "Log level", "info", []string{"debug", "info", "warn", "error"})
	_ = root.AddFlag(levelFlag)

	values, err := GetEnumValues(root, "/", "--invalid-flag")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected empty values for invalid flag, got %v", values)
	}
}

// TestHandleEnum_NonEnumFlag 测试非枚举类型标志
//
// 验证非枚举类型标志返回空列表
func TestHandleEnum_NonEnumFlag(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	stringFlag := mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, "")
	_ = root.AddFlag(stringFlag)

	values, err := GetEnumValues(root, "/", "--output")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected empty values for non-enum flag, got %v", values)
	}
}

// TestHandleEnum_NestedContext 测试嵌套上下文的枚举标志
//
// 验证嵌套子命令中的枚举标志能正确获取枚举值
func TestHandleEnum_NestedContext(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	serverCmd := mock.NewMockCommandBasic("server", "", "Server management")
	modeFlag := mock.NewMockEnumFlag("mode", "m", "Server mode", "http", []string{"http", "https", "grpc"})
	_ = serverCmd.AddFlag(modeFlag)
	_ = root.AddSubCmds(serverCmd)

	values, err := GetEnumValues(root, "/server/", "--mode")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 enum values, got %d: %v", len(values), values)
	}

	// 验证枚举值
	expected := map[string]bool{
		"http":  false,
		"https": false,
		"grpc":  false,
	}
	for _, v := range values {
		if _, ok := expected[v]; ok {
			expected[v] = true
		}
	}
	for v, found := range expected {
		if !found {
			t.Errorf("Expected enum value '%s' not found", v)
		}
	}
}

// TestFindFlagByName_LongName 测试长标志名查找
//
// 验证通过长标志名（带 -- 前缀）能正确找到标志
func TestFindFlagByName_LongName(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	outputFlag := mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, "")
	_ = root.AddFlag(outputFlag)

	found := findFlagByName(root, "--output")
	if found == nil {
		t.Error("Expected to find flag by long name")
	}
	if found.LongName() != "output" {
		t.Errorf("Expected flag name 'output', got '%s'", found.LongName())
	}
}

// TestFindFlagByName_ShortName 测试短标志名查找
//
// 验证通过短标志名（带 - 前缀）能正确找到标志
func TestFindFlagByName_ShortName(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	outputFlag := mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, "")
	_ = root.AddFlag(outputFlag)

	found := findFlagByName(root, "-o")
	if found == nil {
		t.Error("Expected to find flag by short name")
	}
	if found.ShortName() != "o" {
		t.Errorf("Expected flag short name 'o', got '%s'", found.ShortName())
	}
}

// TestFindFlagByName_WithEqualsSuffix 测试带等号后缀的标志名
//
// 验证带 = 后缀的标志名能正确查找
func TestFindFlagByName_WithEqualsSuffix(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	outputFlag := mock.NewMockFlag("output", "o", "Output file", types.FlagTypeString, "")
	_ = root.AddFlag(outputFlag)

	found := findFlagByName(root, "--output=")
	if found == nil {
		t.Error("Expected to find flag with equals suffix")
	}
	if found.LongName() != "output" {
		t.Errorf("Expected flag name 'output', got '%s'", found.LongName())
	}
}

// TestFindFlagByName_NotFound 测试未找到标志
//
// 验证查找不存在的标志返回 nil
func TestFindFlagByName_NotFound(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	found := findFlagByName(root, "--nonexistent")
	if found != nil {
		t.Error("Expected nil for non-existent flag")
	}
}

// TestFindFlagByName_NoShortName 测试无短名称的标志
//
// 验证只有长名称的标志通过短名称查找返回 nil
func TestFindFlagByName_NoShortName(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	verboseFlag := mock.NewMockFlag("verbose", "", "Verbose", types.FlagTypeBool, false)
	_ = root.AddFlag(verboseFlag)

	// 通过长名称应该找到
	found := findFlagByName(root, "--verbose")
	if found == nil {
		t.Error("Expected to find flag by long name")
	}

	// 通过短名称（-）不应该找到，因为该标志没有短名称
	// 注意：这里传入 "-" 会尝试匹配空字符串的短名称
	found = findFlagByName(root, "-")
	// 由于 verboseFlag 的 short 是空字符串，不会被匹配
	if found != nil {
		t.Error("Expected nil when searching by empty short name prefix")
	}
}

// TestGetEnumValues_EnumType 测试枚举类型获取值
//
// 验证枚举类型标志能正确返回枚举值
func TestGetEnumValues_EnumType(t *testing.T) {
	enumFlag := mock.NewMockEnumFlag("level", "l", "Level", "medium", []string{"low", "medium", "high"})

	values := getEnumValues(enumFlag)

	if len(values) != 3 {
		t.Errorf("Expected 3 enum values, got %d", len(values))
	}
}

// TestGetEnumValues_NonEnumType 测试非枚举类型
//
// 验证非枚举类型标志返回空切片
func TestGetEnumValues_NonEnumType(t *testing.T) {
	stringFlag := mock.NewMockFlag("name", "n", "Name", types.FlagTypeString, "")

	values := getEnumValues(stringFlag)

	if len(values) != 0 {
		t.Errorf("Expected empty values for non-enum flag, got %v", values)
	}
}

// TestGetEnumValues_EmptyEnum 测试空枚举值
//
// 验证没有设置枚举值的枚举类型标志返回空切片
func TestGetEnumValues_EmptyEnum(t *testing.T) {
	// 使用 NewMockEnumFlag 但传入空切片，创建没有枚举值的枚举标志
	emptyEnumFlag := mock.NewMockEnumFlag("empty", "e", "Empty enum", "", []string{})

	values := getEnumValues(emptyEnumFlag)

	if len(values) != 0 {
		t.Errorf("Expected empty values for empty enum, got %v", values)
	}
}

// TestHandleEnum_EmptyFlagName 测试空标志名
//
// 验证空标志名返回空列表
func TestHandleEnum_EmptyFlagName(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")

	values, err := GetEnumValues(root, "/", "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 0 {
		t.Errorf("Expected empty values for empty flag name, got %v", values)
	}
}

// TestHandleEnum_MultipleFlags 测试多个标志中查找
//
// 验证在多个标志中能正确找到目标枚举标志
func TestHandleEnum_MultipleFlags(t *testing.T) {
	root := mock.NewMockCommandBasic("myapp", "", "Test application")
	_ = root.AddFlag(mock.NewMockFlag("output", "o", "Output", types.FlagTypeString, ""))
	_ = root.AddFlag(mock.NewMockEnumFlag("format", "f", "Format", "json", []string{"json", "yaml", "xml"}))
	_ = root.AddFlag(mock.NewMockFlag("verbose", "v", "Verbose", types.FlagTypeBool, false))

	values, err := GetEnumValues(root, "/", "--format")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(values) != 3 {
		t.Errorf("Expected 3 enum values, got %d", len(values))
	}

	// 验证找到了正确的标志的枚举值
	expected := map[string]bool{
		"json": false,
		"yaml": false,
		"xml":  false,
	}
	for _, v := range values {
		if _, ok := expected[v]; ok {
			expected[v] = true
		}
	}
	for v, found := range expected {
		if !found {
			t.Errorf("Expected enum value '%s' not found", v)
		}
	}
}
