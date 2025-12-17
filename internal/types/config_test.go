// Package types 配置类型测试
// 本文件包含了配置数据类型的单元测试，测试帮助信息设置、
// 版本信息等配置数据的定义和管理功能的正确性。
package types

import (
	"reflect"
	"strings"
	"testing"
)

// TestNewCmdConfig_基本功能 测试NewCmdConfig的基本功能
func TestNewCmdConfig_基本功能(t *testing.T) {
	config := NewCmdConfig()

	//nolint:all
	if config == nil {
		t.Fatal("NewCmdConfig返回了nil")
	}

	// 验证默认值
	//nolint:all
	if config.Version != "" {
		t.Errorf("Version默认值应为空字符串, 实际: %q", config.Version)
	}

	if config.Description != "" {
		t.Errorf("Description默认值应为空字符串, 实际: %q", config.Description)
	}

	if config.Help != "" {
		t.Errorf("Help默认值应为空字符串, 实际: %q", config.Help)
	}

	if config.UsageSyntax != "" {
		t.Errorf("UsageSyntax默认值应为空字符串, 实际: %q", config.UsageSyntax)
	}

	if config.ModuleHelps != "" {
		t.Errorf("ModuleHelps默认值应为空字符串, 实际: %q", config.ModuleHelps)
	}

	if config.LogoText != "" {
		t.Errorf("LogoText默认值应为空字符串, 实际: %q", config.LogoText)
	}

	if config.Notes == nil {
		t.Error("Notes应该初始化为空切片而不是nil")
	}

	if len(config.Notes) != 0 {
		t.Errorf("Notes初始长度应为0, 实际: %d", len(config.Notes))
	}

	if config.Examples == nil {
		t.Error("Examples应该初始化为空切片而不是nil")
	}

	if len(config.Examples) != 0 {
		t.Errorf("Examples初始长度应为0, 实际: %d", len(config.Examples))
	}

	if config.UseChinese != false {
		t.Errorf("UseChinese默认值应为false, 实际: %v", config.UseChinese)
	}

	if config.NoBuiltinExit != false {
		t.Errorf("NoBuiltinExit默认值应为false, 实际: %v", config.NoBuiltinExit)
	}

	if config.Completion != false {
		t.Errorf("Completion默认值应为false, 实际: %v", config.Completion)
	}
}

// TestCmdConfig_字段赋值 测试CmdConfig各字段的赋值
func TestCmdConfig_字段赋值(t *testing.T) {
	config := NewCmdConfig()

	// 测试字符串字段
	testCases := []struct {
		fieldName string
		setValue  string
		getValue  func() string
	}{
		{
			fieldName: "Version",
			setValue:  "1.0.0",
			getValue:  func() string { return config.Version },
		},
		{
			fieldName: "Description",
			setValue:  "测试描述",
			getValue:  func() string { return config.Description },
		},
		{
			fieldName: "Help",
			setValue:  "帮助信息",
			getValue:  func() string { return config.Help },
		},
		{
			fieldName: "UsageSyntax",
			setValue:  "myapp [选项] <文件>",
			getValue:  func() string { return config.UsageSyntax },
		},
		{
			fieldName: "ModuleHelps",
			setValue:  "模块帮助",
			getValue:  func() string { return config.ModuleHelps },
		},
		{
			fieldName: "LogoText",
			setValue:  "ASCII Logo",
			getValue:  func() string { return config.LogoText },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			// 使用反射设置值
			v := reflect.ValueOf(config).Elem()
			field := v.FieldByName(tc.fieldName)
			if !field.IsValid() {
				t.Fatalf("字段 %s 不存在", tc.fieldName)
			}
			field.SetString(tc.setValue)

			// 验证值是否正确设置
			gotValue := tc.getValue()
			if gotValue != tc.setValue {
				t.Errorf("%s 值不匹配: 期望 %q, 实际 %q", tc.fieldName, tc.setValue, gotValue)
			}
		})
	}

	// 测试布尔字段
	boolTests := []struct {
		fieldName string
		setValue  bool
		getValue  func() bool
	}{
		{
			fieldName: "UseChinese",
			setValue:  true,
			getValue:  func() bool { return config.UseChinese },
		},
		{
			fieldName: "NoBuiltinExit",
			setValue:  false,
			getValue:  func() bool { return config.NoBuiltinExit },
		},
		{
			fieldName: "Completion",
			setValue:  true,
			getValue:  func() bool { return config.Completion },
		},
	}

	for _, tc := range boolTests {
		t.Run(tc.fieldName, func(t *testing.T) {
			// 使用反射设置值
			v := reflect.ValueOf(config).Elem()
			field := v.FieldByName(tc.fieldName)
			if !field.IsValid() {
				t.Fatalf("字段 %s 不存在", tc.fieldName)
			}
			field.SetBool(tc.setValue)

			// 验证值是否正确设置
			gotValue := tc.getValue()
			if gotValue != tc.setValue {
				t.Errorf("%s 值不匹配: 期望 %v, 实际 %v", tc.fieldName, tc.setValue, gotValue)
			}
		})
	}
}

// TestCmdConfig_切片操作 测试Notes和Examples切片的操作
func TestCmdConfig_切片操作(t *testing.T) {
	config := NewCmdConfig()

	// 测试Notes切片操作
	t.Run("Notes切片操作", func(t *testing.T) {
		// 添加备注
		testNotes := []string{
			"第一个备注",
			"",
			"包含\n换行符的备注",
			"包含特殊字符的备注: @#$%^&*()",
			"很长很长很长很长很长很长很长很长很长很长的备注信息",
			"中文备注：这是一个中文备注",
			"Unicode备注: 🎉🚀✨",
		}

		config.Notes = append(config.Notes, testNotes...)

		// 验证Notes
		if len(config.Notes) != len(testNotes) {
			t.Errorf("Notes长度不匹配: 期望 %d, 实际 %d", len(testNotes), len(config.Notes))
		}

		for i, expectedNote := range testNotes {
			if i >= len(config.Notes) {
				t.Errorf("缺少第%d个备注", i)
				continue
			}
			if config.Notes[i] != expectedNote {
				t.Errorf("第%d个备注不匹配: 期望 %q, 实际 %q", i, expectedNote, config.Notes[i])
			}
		}

		// 测试清空Notes
		config.Notes = []string{}
		if len(config.Notes) != 0 {
			t.Errorf("清空后Notes长度应为0, 实际: %d", len(config.Notes))
		}
	})

	// 测试Examples切片操作
	t.Run("Examples切片操作", func(t *testing.T) {
		testExamples := []ExampleInfo{
			{Description: "基本用法", Usage: "myapp file.txt"},
			{Description: "", Usage: "myapp --help"},
			{Description: "复杂用法", Usage: "myapp --config /path/to/config.json --verbose file1.txt file2.txt"},
			{Description: "包含特殊字符", Usage: "myapp 'file with spaces.txt'"},
			{Description: "多行用法", Usage: "myapp \\\n  --option1 value1 \\\n  --option2 value2"},
			{Description: "中文示例", Usage: "myapp --配置 配置文件.json"},
			{Description: "Unicode示例", Usage: "myapp 🚀 --emoji ✨"},
		}

		// 添加所有示例
		config.Examples = append(config.Examples, testExamples...)

		// 验证Examples
		if len(config.Examples) != len(testExamples) {
			t.Errorf("Examples长度不匹配: 期望 %d, 实际 %d", len(testExamples), len(config.Examples))
		}

		for i, expectedExample := range testExamples {
			if i >= len(config.Examples) {
				t.Errorf("缺少第%d个示例", i)
				continue
			}
			if config.Examples[i].Description != expectedExample.Description {
				t.Errorf("第%d个示例描述不匹配: 期望 %q, 实际 %q", i, expectedExample.Description, config.Examples[i].Description)
			}
			if config.Examples[i].Usage != expectedExample.Usage {
				t.Errorf("第%d个示例用法不匹配: 期望 %q, 实际 %q", i, expectedExample.Usage, config.Examples[i].Usage)
			}
		}

		// 测试清空Examples
		config.Examples = []ExampleInfo{}
		if len(config.Examples) != 0 {
			t.Errorf("清空后Examples长度应为0, 实际: %d", len(config.Examples))
		}
	})
}

// TestExampleInfo_结构体 测试ExampleInfo结构体
func TestExampleInfo_结构体(t *testing.T) {
	tests := []struct {
		name        string
		description string
		usage       string
		testDesc    string
	}{
		{
			name:        "正常示例",
			description: "基本用法示例",
			usage:       "myapp input.txt",
			testDesc:    "正常的示例信息",
		},
		{
			name:        "空描述",
			description: "",
			usage:       "myapp --help",
			testDesc:    "描述为空的示例",
		},
		{
			name:        "空用法",
			description: "空用法示例",
			usage:       "",
			testDesc:    "用法为空的示例",
		},
		{
			name:        "都为空",
			description: "",
			usage:       "",
			testDesc:    "描述和用法都为空的示例",
		},
		{
			name:        "多行描述",
			description: "第一行描述\n第二行描述\n第三行描述",
			usage:       "myapp --multi-line",
			testDesc:    "多行描述的示例",
		},
		{
			name:        "多行用法",
			description: "复杂命令示例",
			usage:       "myapp \\\n  --option1 value1 \\\n  --option2 value2 \\\n  input.txt",
			testDesc:    "多行用法的示例",
		},
		{
			name:        "特殊字符",
			description: "包含特殊字符: @#$%^&*()",
			usage:       "myapp --special '@#$%^&*()'",
			testDesc:    "包含特殊字符的示例",
		},
		{
			name:        "Unicode字符",
			description: "Unicode示例: 🎉🚀✨",
			usage:       "myapp --emoji '🎉🚀✨'",
			testDesc:    "包含Unicode字符的示例",
		},
		{
			name:        "极长文本",
			description: strings.Repeat("很长的描述。", 100),
			usage:       strings.Repeat("myapp --very-long-option ", 50),
			testDesc:    "极长文本的示例",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			example := ExampleInfo{
				Description: tt.description,
				Usage:       tt.usage,
			}

			if example.Description != tt.description {
				t.Errorf("Description不匹配: 期望 %q, 实际 %q", tt.description, example.Description)
			}

			if example.Usage != tt.usage {
				t.Errorf("Usage不匹配: 期望 %q, 实际 %q", tt.usage, example.Usage)
			}
		})
	}
}

// TestCmdConfig_极值测试 测试极值情况
func TestCmdConfig_极值测试(t *testing.T) {
	config := NewCmdConfig()

	// 测试极长字符串
	extremelyLongString := strings.Repeat("a", 100000)

	config.Version = extremelyLongString
	if config.Version != extremelyLongString {
		t.Error("极长Version字符串设置失败")
	}

	config.Description = extremelyLongString
	if config.Description != extremelyLongString {
		t.Error("极长Description字符串设置失败")
	}

	// 测试包含所有ASCII字符的字符串
	allASCII := ""
	for i := 32; i <= 126; i++ {
		allASCII += string(rune(i))
	}

	config.Help = allASCII
	if config.Help != allASCII {
		t.Error("包含所有ASCII字符的Help字符串设置失败")
	}

	// 测试Unicode字符串
	unicodeString := "测试🎉🚀✨中文和emoji混合内容"
	config.LogoText = unicodeString
	if config.LogoText != unicodeString {
		t.Error("Unicode字符串设置失败")
	}

	// 测试大量Notes
	for i := 0; i < 10000; i++ {
		config.Notes = append(config.Notes, "note")
	}
	if len(config.Notes) != 10000 {
		t.Errorf("大量Notes添加失败: 期望 10000, 实际 %d", len(config.Notes))
	}

	// 测试大量Examples
	for i := 0; i < 5000; i++ {
		config.Examples = append(config.Examples, ExampleInfo{
			Description: "example",
			Usage:       "usage",
		})
	}
	if len(config.Examples) != 5000 {
		t.Errorf("大量Examples添加失败: 期望 5000, 实际 %d", len(config.Examples))
	}
}

// TestCmdConfig_内存使用 测试内存使用情况
func TestCmdConfig_内存使用(t *testing.T) {
	// 创建大量配置实例
	configs := make([]*CmdConfig, 1000)
	for i := 0; i < 1000; i++ {
		configs[i] = NewCmdConfig()

		// 添加一些数据
		configs[i].Version = "1.0.0"
		configs[i].Description = "测试描述"
		configs[i].Notes = append(configs[i].Notes, "note1", "note2", "note3")
		configs[i].Examples = append(configs[i].Examples,
			ExampleInfo{Description: "desc", Usage: "usage"})
	}

	// 验证所有配置都正确创建
	for i, config := range configs {
		if config == nil {
			t.Error("配置不应为nil")
			return
		}
		if len(config.Notes) != 3 {
			t.Errorf("第%d个配置Notes数量不正确", i)
		}
		if len(config.Examples) != 1 {
			t.Errorf("第%d个配置Examples数量不正确", i)
		}
	}

	// 清理引用
	for i := range configs {
		configs[i] = nil
	}
	// 清空配置切片
	_ = configs[:0]

	t.Log("内存使用测试完成")
}

// TestCmdConfig_字段完整性 测试所有字段的完整性
func TestCmdConfig_字段完整性(t *testing.T) {
	config := NewCmdConfig()

	// 使用反射检查所有字段
	v := reflect.ValueOf(config).Elem()
	typ := v.Type()

	expectedFields := map[string]reflect.Kind{
		"Version":       reflect.String,
		"Description":   reflect.String,
		"Help":          reflect.String,
		"UsageSyntax":   reflect.String,
		"ModuleHelps":   reflect.String,
		"LogoText":      reflect.String,
		"Notes":         reflect.Slice,
		"Examples":      reflect.Slice,
		"UseChinese":    reflect.Bool,
		"NoBuiltinExit": reflect.Bool,
		"Completion":    reflect.Bool,
	}

	// 检查所有期望的字段是否存在
	for expectedField, expectedKind := range expectedFields {
		field := v.FieldByName(expectedField)
		if !field.IsValid() {
			t.Errorf("缺少字段: %s", expectedField)
			continue
		}

		if field.Kind() != expectedKind {
			t.Errorf("字段 %s 类型不正确: 期望 %v, 实际 %v",
				expectedField, expectedKind, field.Kind())
		}
	}

	// 检查是否有意外的字段
	for i := 0; i < v.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if _, exists := expectedFields[fieldName]; !exists {
			t.Errorf("发现意外字段: %s", fieldName)
		}
	}

	t.Logf("字段完整性检查完成，共检查了 %d 个字段", len(expectedFields))
}
