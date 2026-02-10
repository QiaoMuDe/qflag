package utils

import (
	"strings"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestGetCmdName 测试 GetCmdName 函数
func TestGetCmdName(t *testing.T) {
	helper := mock.NewTestHelper()

	tests := []struct {
		name     string
		cmd      types.Command
		expected string
	}{
		{
			name:     "有长名和短名",
			cmd:      helper.CreateMockCommandWithFlags("test-cmd", "t", "Test command"),
			expected: "test-cmd, t\n",
		},
		{
			name:     "只有名称",
			cmd:      helper.CreateMockCommandWithFlags("test", "", "Test command"),
			expected: "test\n",
		},
		// 注意: GetCmdName 函数没有处理 cmd 为 nil 的情况, 所以不测试这种情况
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCmdName(tt.cmd)
			if result != tt.expected {
				t.Errorf("GetCmdName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestFormatDefaultValue 测试 FormatDefaultValue 函数
func TestFormatDefaultValue(t *testing.T) {
	tests := []struct {
		name     string
		flagType types.FlagType
		defValue any
		expected string
	}{
		{
			name:     "nil 值",
			flagType: types.FlagTypeString,
			defValue: nil,
			expected: "",
		},
		{
			name:     "空字符串",
			flagType: types.FlagTypeString,
			defValue: "",
			expected: `""`,
		},
		{
			name:     "非空字符串",
			flagType: types.FlagTypeString,
			defValue: "hello",
			expected: "hello",
		},
		{
			name:     "布尔值 true",
			flagType: types.FlagTypeBool,
			defValue: true,
			expected: "true",
		},
		{
			name:     "布尔值 false",
			flagType: types.FlagTypeBool,
			defValue: false,
			expected: "false",
		},
		{
			name:     "整数",
			flagType: types.FlagTypeInt,
			defValue: 42,
			expected: "42",
		},
		{
			name:     "int64",
			flagType: types.FlagTypeInt64,
			defValue: int64(64),
			expected: "64",
		},
		{
			name:     "uint8",
			flagType: types.FlagTypeUint8,
			defValue: uint8(8),
			expected: "8",
		},
		{
			name:     "uint16",
			flagType: types.FlagTypeUint16,
			defValue: uint16(16),
			expected: "16",
		},
		{
			name:     "uint32",
			flagType: types.FlagTypeUint32,
			defValue: uint32(32),
			expected: "32",
		},
		{
			name:     "uint64",
			flagType: types.FlagTypeUint64,
			defValue: uint64(64),
			expected: "64",
		},
		{
			name:     "浮点数",
			flagType: types.FlagTypeFloat64,
			defValue: 3.14159,
			expected: "3.14",
		},
		{
			name:     "空字符串数组",
			flagType: types.FlagTypeStringSlice,
			defValue: []string{},
			expected: "[]",
		},
		{
			name:     "非空字符串数组",
			flagType: types.FlagTypeStringSlice,
			defValue: []string{"a", "b", "c"},
			expected: "[a b c]",
		},
		{
			name:     "空整数数组",
			flagType: types.FlagTypeIntSlice,
			defValue: []int{},
			expected: "[]",
		},
		{
			name:     "非空整数数组",
			flagType: types.FlagTypeIntSlice,
			defValue: []int{1, 2, 3},
			expected: "[1 2 3]",
		},
		{
			name:     "空int64数组",
			flagType: types.FlagTypeInt64Slice,
			defValue: []int64{},
			expected: "[]",
		},
		{
			name:     "非空int64数组",
			flagType: types.FlagTypeInt64Slice,
			defValue: []int64{10, 20, 30},
			expected: "[10 20 30]",
		},
		{
			name:     "时间间隔",
			flagType: types.FlagTypeDuration,
			defValue: 5 * time.Second,
			expected: "5s",
		},
		{
			name:     "零时间",
			flagType: types.FlagTypeTime,
			defValue: time.Time{},
			expected: `""`,
		},
		{
			name:     "非零时间",
			flagType: types.FlagTypeTime,
			defValue: time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC),
			expected: "2023-05-15 14:30:00",
		},
		{
			name:     "空字符串映射",
			flagType: types.FlagTypeMap,
			defValue: map[string]string{},
			expected: "{}",
		},
		{
			name:     "非空字符串映射",
			flagType: types.FlagTypeMap,
			defValue: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:     "枚举类型",
			flagType: types.FlagTypeEnum,
			defValue: "option1",
			expected: "option1",
		},
		{
			name:     "大小类型",
			flagType: types.FlagTypeSize,
			defValue: 1024,
			expected: "1024 bytes",
		},
		{
			name:     "未知类型",
			flagType: types.FlagTypeUnknown,
			defValue: struct{ Name string }{Name: "test"},
			expected: "{test}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDefaultValue(tt.flagType, tt.defValue)

			// 对于 map 类型, 特殊处理, 不检查顺序
			if tt.flagType == types.FlagTypeMap && tt.expected == "" {
				// 检查结果是否包含所有键值对
				if !strings.Contains(result, "key1=value1") || !strings.Contains(result, "key2=value2") {
					t.Errorf("FormatDefaultValue() = %v, expected to contain key1=value1 and key2=value2", result)
				}
				return
			}

			// 对于其他类型, 检查精确匹配
			if result != tt.expected {
				t.Errorf("FormatDefaultValue() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCalcOptionMaxWidth 测试 CalcOptionMaxWidth 函数
func TestCalcOptionMaxWidth(t *testing.T) {
	tests := []struct {
		name     string
		options  []types.OptionInfo
		expected int
	}{
		{
			name:     "空选项列表",
			options:  []types.OptionInfo{},
			expected: 0,
		},
		{
			name: "单个选项",
			options: []types.OptionInfo{
				{NamePart: "--help"},
			},
			expected: 6,
		},
		{
			name: "多个选项, 不同长度",
			options: []types.OptionInfo{
				{NamePart: "-h"},
				{NamePart: "--verbose"},
				{NamePart: "--config"},
			},
			expected: 9, // "--verbose" 是最长的, 长度为9
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalcOptionMaxWidth(tt.options)
			if result != tt.expected {
				t.Errorf("CalcOptionMaxWidth() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCalcSubCmdMaxLen 测试 CalcSubCmdMaxLen 函数
func TestCalcSubCmdMaxLen(t *testing.T) {
	tests := []struct {
		name     string
		subCmds  []types.SubCmdInfo
		expected int
	}{
		{
			name:     "空子命令列表",
			subCmds:  []types.SubCmdInfo{},
			expected: 0,
		},
		{
			name: "单个子命令",
			subCmds: []types.SubCmdInfo{
				{Name: "help"},
			},
			expected: 4,
		},
		{
			name: "多个子命令, 不同长度",
			subCmds: []types.SubCmdInfo{
				{Name: "ls"},
				{Name: "list"},
				{Name: "show"},
			},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalcSubCmdMaxLen(tt.subCmds)
			if result != tt.expected {
				t.Errorf("CalcSubCmdMaxLen() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestSortOptions 测试 SortOptions 函数
func TestSortOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []types.OptionInfo
		expected []types.OptionInfo
	}{
		{
			name:     "空选项列表",
			options:  []types.OptionInfo{},
			expected: []types.OptionInfo{},
		},
		{
			name: "长短都有 > 仅长选项 > 仅短选项",
			options: []types.OptionInfo{
				{NamePart: "-h"},
				{NamePart: "--verbose"},
				{NamePart: "-v, --version"},
			},
			expected: []types.OptionInfo{
				{NamePart: "-v, --version"},
				{NamePart: "--verbose"},
				{NamePart: "-h"},
			},
		},
		{
			name: "同组内按首字母排序",
			options: []types.OptionInfo{
				{NamePart: "--zebra"},
				{NamePart: "--apple"},
				{NamePart: "--banana"},
			},
			expected: []types.OptionInfo{
				{NamePart: "--apple"},
				{NamePart: "--banana"},
				{NamePart: "--zebra"},
			},
		},
		{
			name: "大小写不敏感排序",
			options: []types.OptionInfo{
				{NamePart: "--Zebra"},
				{NamePart: "--apple"},
				{NamePart: "--Banana"},
			},
			expected: []types.OptionInfo{
				{NamePart: "--apple"},
				{NamePart: "--Banana"},
				{NamePart: "--Zebra"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制一份选项以避免修改原始数据
			options := make([]types.OptionInfo, len(tt.options))
			copy(options, tt.options)

			SortOptions(options)

			// 检查长度是否一致
			if len(options) != len(tt.expected) {
				t.Errorf("SortOptions() length = %v, expected %v", len(options), len(tt.expected))
				return
			}

			// 检查每个元素是否一致
			for i := range options {
				if options[i].NamePart != tt.expected[i].NamePart {
					t.Errorf("SortOptions() at index %d = %v, expected %v", i, options[i].NamePart, tt.expected[i].NamePart)
				}
			}
		})
	}
}

// TestSortSubCmds 测试 SortSubCmds 函数
func TestSortSubCmds(t *testing.T) {
	tests := []struct {
		name     string
		subCmds  []types.SubCmdInfo
		expected []types.SubCmdInfo
	}{
		{
			name:     "空子命令列表",
			subCmds:  []types.SubCmdInfo{},
			expected: []types.SubCmdInfo{},
		},
		{
			name: "长短名都有 > 仅单一名字",
			subCmds: []types.SubCmdInfo{
				{Name: "help"},
				{Name: "ls, list"},
				{Name: "show"},
			},
			expected: []types.SubCmdInfo{
				{Name: "ls, list"},
				{Name: "help"},
				{Name: "show"},
			},
		},
		{
			name: "同组内按首字母排序",
			subCmds: []types.SubCmdInfo{
				{Name: "zebra"},
				{Name: "apple"},
				{Name: "banana"},
			},
			expected: []types.SubCmdInfo{
				{Name: "apple"},
				{Name: "banana"},
				{Name: "zebra"},
			},
		},
		{
			name: "大小写不敏感排序",
			subCmds: []types.SubCmdInfo{
				{Name: "Zebra"},
				{Name: "apple"},
				{Name: "Banana"},
			},
			expected: []types.SubCmdInfo{
				{Name: "apple"},
				{Name: "Banana"},
				{Name: "Zebra"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 复制一份子命令以避免修改原始数据
			subCmds := make([]types.SubCmdInfo, len(tt.subCmds))
			copy(subCmds, tt.subCmds)

			SortSubCmds(subCmds)

			// 检查长度是否一致
			if len(subCmds) != len(tt.expected) {
				t.Errorf("SortSubCmds() length = %v, expected %v", len(subCmds), len(tt.expected))
				return
			}

			// 检查每个元素是否一致
			for i := range subCmds {
				if subCmds[i].Name != tt.expected[i].Name {
					t.Errorf("SortSubCmds() at index %d = %v, expected %v", i, subCmds[i].Name, tt.expected[i].Name)
				}
			}
		})
	}
}

// TestValidateFlagName 测试 ValidateFlagName 函数
func TestValidateFlagName(t *testing.T) {
	helper := mock.NewTestHelper()

	tests := []struct {
		name      string
		cmd       types.Command
		longName  string
		shortName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "空命令",
			cmd:       nil,
			longName:  "help",
			shortName: "h",
			wantErr:   true,
			errMsg:    "cmd cannot be nil",
		},
		{
			name:      "空长名和短名",
			cmd:       helper.CreateMockCommandWithFlags("test", "", "Test command"),
			longName:  "",
			shortName: "",
			wantErr:   true,
			errMsg:    "flag name cannot be empty",
		},
		{
			name:      "有效的长名和短名",
			cmd:       helper.CreateMockCommandWithFlags("test", "", "Test command"),
			longName:  "help",
			shortName: "h",
			wantErr:   false,
		},
		{
			name: "重复的长名",
			cmd: func() types.Command {
				cmd := helper.CreateMockCommandWithFlags("test", "", "Test command")
				// 添加一个已存在的标志
				helpFlag := helper.CreateMockBoolFlag("help", "h", "Help flag", false)
				if err := cmd.AddFlag(helpFlag); err != nil {
					panic(err) // 在测试中, 如果添加标志失败, 应该立即 panic
				}
				return cmd
			}(),
			longName:  "help",
			shortName: "i",
			wantErr:   true,
			errMsg:    "flag '--help' already exists",
		},
		{
			name: "重复的短名",
			cmd: func() types.Command {
				cmd := helper.CreateMockCommandWithFlags("test", "", "Test command")
				// 添加一个已存在的标志
				helpFlag := helper.CreateMockBoolFlag("someflag", "h", "Some flag", false)
				if err := cmd.AddFlag(helpFlag); err != nil {
					panic(err) // 在测试中, 如果添加标志失败, 应该立即 panic
				}
				return cmd
			}(),
			longName:  "info",
			shortName: "h",
			wantErr:   true,
			errMsg:    "flag '-h' already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFlagName(tt.cmd, tt.longName, tt.shortName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFlagName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateFlagName() error = %v, expected to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestFormatFlagName 测试 FormatFlagName 函数
func TestFormatFlagName(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		expected  string
	}{
		{
			name:      "长名和短名都有",
			longName:  "help",
			shortName: "h",
			expected:  "-h, --help",
		},
		{
			name:      "只有长名",
			longName:  "verbose",
			shortName: "",
			expected:  "--verbose",
		},
		{
			name:      "只有短名",
			longName:  "",
			shortName: "v",
			expected:  "-v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFlagName(tt.longName, tt.shortName)
			if result != tt.expected {
				t.Errorf("FormatFlagName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestToStrSlice 测试 ToStrSlice 函数
func TestToStrSlice(t *testing.T) {
	tests := []struct {
		name     string
		slice    any
		expected []string
		wantErr  bool
	}{
		{
			name:     "nil 切片",
			slice:    nil,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "非切片类型",
			slice:    "not a slice",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "字符串切片",
			slice:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "整数切片",
			slice:    []int{1, 2, 3},
			expected: []string{"1", "2", "3"},
			wantErr:  false,
		},
		{
			name:     "浮点数切片",
			slice:    []float64{1.1, 2.2, 3.3},
			expected: []string{"1.1", "2.2", "3.3"},
			wantErr:  false,
		},
		{
			name:     "布尔切片",
			slice:    []bool{true, false, true},
			expected: []string{"true", "false", "true"},
			wantErr:  false,
		},
		{
			name:     "混合类型切片",
			slice:    []interface{}{1, "two", 3.0, true},
			expected: []string{"1", "two", "3", "true"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToStrSlice(tt.slice)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToStrSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(result) != len(tt.expected) {
					t.Errorf("ToStrSlice() length = %v, expected %v", len(result), len(tt.expected))
					return
				}
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("ToStrSlice() at index %d = %v, expected %v", i, result[i], tt.expected[i])
					}
				}
			}
		})
	}
}

// TestFormatSize 测试 FormatSize 函数
func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected string
	}{
		{
			name:     "负数",
			size:     -1,
			expected: "0B",
		},
		{
			name:     "零",
			size:     0,
			expected: "0B",
		},
		{
			name:     "小于1KB",
			size:     512,
			expected: "512B",
		},
		{
			name:     "正好1KB",
			size:     1000,
			expected: "1.00KB",
		},
		{
			name:     "小于1MB",
			size:     1500,
			expected: "1.50KB",
		},
		{
			name:     "正好1MB",
			size:     1000000,
			expected: "1.00MB",
		},
		{
			name:     "小于1GB",
			size:     1500000,
			expected: "1.50MB",
		},
		{
			name:     "正好1GB",
			size:     1000000000,
			expected: "1.00GB",
		},
		{
			name:     "小于1TB",
			size:     1500000000,
			expected: "1.50GB",
		},
		{
			name:     "正好1TB",
			size:     1000000000000,
			expected: "1.00TB",
		},
		{
			name:     "小于1PB",
			size:     1500000000000,
			expected: "1.50TB",
		},
		{
			name:     "正好1PB",
			size:     1000000000000000,
			expected: "1.00PB",
		},
		{
			name:     "大于1PB",
			size:     1500000000000000,
			expected: "1.50PB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSize(tt.size)
			if result != tt.expected {
				t.Errorf("FormatSize() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// 辅助函数: 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
