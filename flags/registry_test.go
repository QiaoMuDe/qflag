package flags

import (
	"testing"
)

// TestFlagTypeToString 测试FlagType到字符串的转换功能
func TestFlagTypeToString(t *testing.T) {
	// 测试不带括号的情况
	testCases := []struct {
		flagType FlagType
		expected string
	}{
		{FlagTypeInt, "int"},
		{FlagTypeInt64, "int64"},
		{FlagTypeUint16, "uint16"},
		{FlagTypeUint32, "uint32"},
		{FlagTypeUint64, "uint64"},
		{FlagTypeString, "string"},
		{FlagTypeBool, "bool"},
		{FlagTypeFloat64, "float64"},
		{FlagTypeEnum, "enum"},
		{FlagTypeDuration, "duration"},
		{FlagTypeTime, "time"},
		{FlagTypeMap, "map"},
		{FlagTypeSlice, "slice"},
		{FlagTypeUnknown, "unknown"},
		{FlagType(999), "unknown"}, // 测试未知类型
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := FlagTypeToString(tc.flagType, false)
			if result != tc.expected {
				t.Errorf("FlagTypeToString(%v, false) = %s, 期望 %s", tc.flagType, result, tc.expected)
			}
		})
	}
}

// TestFlagTypeToStringWithBrackets 测试带括号的FlagType转换
func TestFlagTypeToStringWithBrackets(t *testing.T) {
	testCases := []struct {
		flagType FlagType
		expected string
	}{
		{FlagTypeInt, "<int>"},
		{FlagTypeString, "<string>"},
		{FlagTypeBool, ""}, // 布尔类型特殊处理
		{FlagTypeFloat64, "<float64>"},
		{FlagTypeUnknown, "<unknown>"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			result := FlagTypeToString(tc.flagType, true)
			if result != tc.expected {
				t.Errorf("FlagTypeToString(%v, true) = %s, 期望 %s", tc.flagType, result, tc.expected)
			}
		})
	}
}

// TestConstants 测试常量定义
func TestConstants(t *testing.T) {
	// 测试标志名称常量
	if HelpFlagName != "help" {
		t.Errorf("HelpFlagName 应为 'help', 实际为 '%s'", HelpFlagName)
	}
	if HelpFlagShortName != "h" {
		t.Errorf("HelpFlagShortName 应为 'h', 实际为 '%s'", HelpFlagShortName)
	}
	if VersionFlagLongName != "version" {
		t.Errorf("VersionFlagLongName 应为 'version', 实际为 '%s'", VersionFlagLongName)
	}
	if VersionFlagShortName != "v" {
		t.Errorf("VersionFlagShortName 应为 'v', 实际为 '%s'", VersionFlagShortName)
	}

	// 测试分隔符常量
	if FlagSplitComma != "," {
		t.Errorf("FlagSplitComma 应为 ',', 实际为 '%s'", FlagSplitComma)
	}
	if FlagSplitSemicolon != ";" {
		t.Errorf("FlagSplitSemicolon 应为 ';', 实际为 '%s'", FlagSplitSemicolon)
	}
	if FlagKVEqual != "=" {
		t.Errorf("FlagKVEqual 应为 '=', 实际为 '%s'", FlagKVEqual)
	}
}

// TestShellSlice 测试Shell类型切片
func TestShellSlice(t *testing.T) {
	expectedShells := []string{ShellNone, ShellBash, ShellPowershell, ShellPwsh}

	if len(ShellSlice) != len(expectedShells) {
		t.Errorf("ShellSlice 长度应为 %d, 实际为 %d", len(expectedShells), len(ShellSlice))
	}

	for i, expected := range expectedShells {
		if i >= len(ShellSlice) || ShellSlice[i] != expected {
			t.Errorf("ShellSlice[%d] 应为 '%s', 实际为 '%s'", i, expected, ShellSlice[i])
		}
	}
}

// TestFlagSplitSlice 测试标志分隔符切片
func TestFlagSplitSlice(t *testing.T) {
	expectedDelimiters := []string{FlagSplitComma, FlagSplitSemicolon, FlagSplitPipe, FlagKVColon}

	if len(FlagSplitSlice) != len(expectedDelimiters) {
		t.Errorf("FlagSplitSlice 长度应为 %d, 实际为 %d", len(expectedDelimiters), len(FlagSplitSlice))
	}

	for i, expected := range expectedDelimiters {
		if i >= len(FlagSplitSlice) || FlagSplitSlice[i] != expected {
			t.Errorf("FlagSplitSlice[%d] 应为 '%s', 实际为 '%s'", i, expected, FlagSplitSlice[i])
		}
	}
}

// TestInvalidFlagChars 测试非法字符常量
func TestInvalidFlagChars(t *testing.T) {
	expectedChars := " !@#$%^&*(){}[]|\\;:'\"<>,.?/"
	if InvalidFlagChars != expectedChars {
		t.Errorf("InvalidFlagChars 应为 '%s', 实际为 '%s'", expectedChars, InvalidFlagChars)
	}
}
