package help

import (
	"flag"
	"reflect"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

// TestSortWithShortNamePriority 测试通用排序比较函数
func TestSortWithShortNamePriority(t *testing.T) {
	tests := []struct {
		name      string
		aHasShort bool
		bHasShort bool
		aName     string
		bName     string
		aShort    string
		bShort    string
		expected  bool
	}{
		{
			name:      "a有短名称，b没有短名称",
			aHasShort: true,
			bHasShort: false,
			aName:     "apple",
			bName:     "banana",
			aShort:    "a",
			bShort:    "",
			expected:  true,
		},
		{
			name:      "a没有短名称，b有短名称",
			aHasShort: false,
			bHasShort: true,
			aName:     "apple",
			bName:     "banana",
			aShort:    "",
			bShort:    "b",
			expected:  false,
		},
		{
			name:      "都有短名称，按长名称排序",
			aHasShort: true,
			bHasShort: true,
			aName:     "apple",
			bName:     "banana",
			aShort:    "a",
			bShort:    "b",
			expected:  true,
		},
		{
			name:      "都有短名称，长名称相同，按短名称排序",
			aHasShort: true,
			bHasShort: true,
			aName:     "test",
			bName:     "test",
			aShort:    "a",
			bShort:    "b",
			expected:  true,
		},
		{
			name:      "都没有短名称，按长名称排序",
			aHasShort: false,
			bHasShort: false,
			aName:     "zebra",
			bName:     "apple",
			aShort:    "",
			bShort:    "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortWithShortNamePriority(
				tt.aHasShort, tt.bHasShort,
				tt.aName, tt.bName,
				tt.aShort, tt.bShort,
			)
			if result != tt.expected {
				t.Errorf("期望 %v，但得到 %v", tt.expected, result)
			}
		})
	}
}

// TestSortFlags 测试标志排序功能
func TestSortFlags(t *testing.T) {
	flags := []flagInfo{
		{longFlag: "zebra", shortFlag: "", usage: "Z flag", defValue: "z", typeStr: "<string>"},
		{longFlag: "apple", shortFlag: "a", usage: "A flag", defValue: "a", typeStr: "<string>"},
		{longFlag: "banana", shortFlag: "", usage: "B flag", defValue: "b", typeStr: "<string>"},
		{longFlag: "cherry", shortFlag: "c", usage: "C flag", defValue: "c", typeStr: "<string>"},
	}

	sortFlags(flags)

	// 验证排序结果：有短名称的优先，然后按长名称排序
	expected := []string{"apple", "cherry", "banana", "zebra"}
	for i, flag := range flags {
		if flag.longFlag != expected[i] {
			t.Errorf("位置 %d: 期望 %s，但得到 %s", i, expected[i], flag.longFlag)
		}
	}
}

// TestSortSubCommands 测试子命令排序功能
func TestSortSubCommands(t *testing.T) {
	subCmds := []*types.CmdContext{
		createTestSubCommand("zebra", ""),
		createTestSubCommand("apple", "a"),
		createTestSubCommand("banana", ""),
		createTestSubCommand("cherry", "c"),
	}

	sortSubCommands(subCmds)

	// 验证排序结果：有短名称的优先，然后按长名称排序
	expected := []string{"apple", "cherry", "banana", "zebra"}
	for i, subCmd := range subCmds {
		if subCmd.LongName != expected[i] {
			t.Errorf("位置 %d: 期望 %s，但得到 %s", i, expected[i], subCmd.LongName)
		}
	}
}

// TestSortByNamePriority 测试通用排序函数
func TestSortByNamePriority(t *testing.T) {
	items := []NamedItem{
		flagInfoItem{flagInfo{longFlag: "zebra", shortFlag: ""}},
		flagInfoItem{flagInfo{longFlag: "apple", shortFlag: "a"}},
		flagInfoItem{flagInfo{longFlag: "banana", shortFlag: ""}},
		flagInfoItem{flagInfo{longFlag: "cherry", shortFlag: "c"}},
	}

	sortByNamePriority(items)

	// 验证排序结果
	expected := []string{"apple", "cherry", "banana", "zebra"}
	for i, item := range items {
		if item.GetLongName() != expected[i] {
			t.Errorf("位置 %d: 期望 %s，但得到 %s", i, expected[i], item.GetLongName())
		}
	}
}

// TestFlagInfoItem_NamedItem 测试 flagInfoItem 实现 NamedItem 接口
func TestFlagInfoItem_NamedItem(t *testing.T) {
	flag := flagInfoItem{
		flagInfo{
			longFlag:  "test",
			shortFlag: "t",
			usage:     "test flag",
			defValue:  "default",
			typeStr:   "<string>",
		},
	}

	if flag.GetLongName() != "test" {
		t.Errorf("期望长名称为 'test'，但得到 '%s'", flag.GetLongName())
	}

	if flag.GetShortName() != "t" {
		t.Errorf("期望短名称为 't'，但得到 '%s'", flag.GetShortName())
	}
}

// TestSubCmdItem_NamedItem 测试 subCmdItem 实现 NamedItem 接口
func TestSubCmdItem_NamedItem(t *testing.T) {
	ctx := createTestSubCommand("testcmd", "tc")
	subCmd := subCmdItem{ctx}

	if subCmd.GetLongName() != "testcmd" {
		t.Errorf("期望长名称为 'testcmd'，但得到 '%s'", subCmd.GetLongName())
	}

	if subCmd.GetShortName() != "tc" {
		t.Errorf("期望短名称为 'tc'，但得到 '%s'", subCmd.GetShortName())
	}
}

// TestSortFlags_EmptySlice 测试空切片排序
func TestSortFlags_EmptySlice(t *testing.T) {
	var flags []flagInfo
	sortFlags(flags) // 不应该崩溃
	if len(flags) != 0 {
		t.Error("空切片排序后长度应该为0")
	}
}

// TestSortSubCommands_EmptySlice 测试空子命令切片排序
func TestSortSubCommands_EmptySlice(t *testing.T) {
	var subCmds []*types.CmdContext
	sortSubCommands(subCmds) // 不应该崩溃
	if len(subCmds) != 0 {
		t.Error("空切片排序后长度应该为0")
	}
}

// TestSortFlags_SingleItem 测试单个标志排序
func TestSortFlags_SingleItem(t *testing.T) {
	flags := []flagInfo{
		{longFlag: "single", shortFlag: "s", usage: "Single flag", defValue: "s", typeStr: "<string>"},
	}

	original := make([]flagInfo, len(flags))
	copy(original, flags)

	sortFlags(flags)

	if !reflect.DeepEqual(flags, original) {
		t.Error("单个元素排序后应该保持不变")
	}
}

// TestSortSubCommands_SingleItem 测试单个子命令排序
func TestSortSubCommands_SingleItem(t *testing.T) {
	subCmds := []*types.CmdContext{
		createTestSubCommand("single", "s"),
	}

	original := subCmds[0]
	sortSubCommands(subCmds)

	if subCmds[0] != original {
		t.Error("单个元素排序后应该保持不变")
	}
}

// TestSortFlags_SameLongName 测试相同长名称的标志排序
func TestSortFlags_SameLongName(t *testing.T) {
	flags := []flagInfo{
		{longFlag: "test", shortFlag: "z", usage: "Test Z", defValue: "z", typeStr: "<string>"},
		{longFlag: "test", shortFlag: "a", usage: "Test A", defValue: "a", typeStr: "<string>"},
	}

	sortFlags(flags)

	// 长名称相同时，应该按短名称排序
	if flags[0].shortFlag != "a" || flags[1].shortFlag != "z" {
		t.Error("相同长名称时应该按短名称排序")
	}
}

// createTestSubCommand 创建测试用的子命令上下文
func createTestSubCommand(longName, shortName string) *types.CmdContext {
	return types.NewCmdContext(longName, shortName, flag.ContinueOnError)
}
