// Package types 内置类型测试
// 本文件包含了内置数据类型的单元测试，测试内置标志、配置选项
// 等核心数据类型的定义、初始化和操作功能的正确性。
package types

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestNewBuiltinFlags_基本功能 测试NewBuiltinFlags的基本功能
func TestNewBuiltinFlags_基本功能(t *testing.T) {
	bf := NewBuiltinFlags()

	//nolint:all
	if bf == nil {
		t.Fatal("NewBuiltinFlags返回了nil")
	}

	// 验证所有字段都已初始化
	//nolint:all
	if bf.Help == nil {
		t.Error("Help字段未初始化")
	}

	//nolint:all
	if bf.Version == nil {
		t.Error("Version字段未初始化")
	}

	//nolint:all
	if bf.Completion == nil {
		t.Error("Completion字段未初始化")
	}

	// 验证字段类型
	//nolint:all
	if reflect.TypeOf(bf.Help) != reflect.TypeOf(&flags.BoolFlag{}) {
		t.Error("Help字段类型不正确")
	}

	//nolint:all
	if reflect.TypeOf(bf.Version) != reflect.TypeOf(&flags.BoolFlag{}) {
		t.Error("Version字段类型不正确")
	}

	//nolint:all
	if reflect.TypeOf(bf.Completion) != reflect.TypeOf(&flags.EnumFlag{}) {
		t.Error("Completion字段类型不正确")
	}

	// 验证NameMap初始为空
	count := 0
	//nolint:all
	bf.NameMap.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("NameMap初始应为空, 实际包含 %d 个元素", count)
	}
}

// TestBuiltinFlags_IsBuiltinFlag_基本功能 测试IsBuiltinFlag的基本功能
func TestBuiltinFlags_IsBuiltinFlag_基本功能(t *testing.T) {
	bf := NewBuiltinFlags()

	// 测试空字符串
	if bf.IsBuiltinFlag("") {
		t.Error("空字符串不应该被识别为内置标志")
	}

	// 测试不存在的标志
	if bf.IsBuiltinFlag("nonexistent") {
		t.Error("不存在的标志不应该被识别为内置标志")
	}

	// 添加一些内置标志
	testFlags := []string{"help", "h", "version", "v", "completion"}
	bf.MarkAsBuiltin(testFlags...)

	// 测试存在的标志
	for _, flagName := range testFlags {
		if !bf.IsBuiltinFlag(flagName) {
			t.Errorf("标志 %q 应该被识别为内置标志", flagName)
		}
	}

	// 测试仍然不存在的标志
	if bf.IsBuiltinFlag("still-nonexistent") {
		t.Error("仍然不存在的标志不应该被识别为内置标志")
	}
}

// TestBuiltinFlags_IsBuiltinFlag_边界场景 测试IsBuiltinFlag的边界场景
func TestBuiltinFlags_IsBuiltinFlag_边界场景(t *testing.T) {
	bf := NewBuiltinFlags()

	tests := []struct {
		name        string
		flagName    string
		shouldMark  bool
		expected    bool
		description string
	}{
		{
			name:        "空字符串",
			flagName:    "",
			shouldMark:  false,
			expected:    false,
			description: "空字符串应该返回false",
		},
		{
			name:        "单字符标志",
			flagName:    "h",
			shouldMark:  true,
			expected:    true,
			description: "单字符标志",
		},
		{
			name:        "长标志名",
			flagName:    "very-long-flag-name-with-many-hyphens",
			shouldMark:  true,
			expected:    true,
			description: "很长的标志名",
		},
		{
			name:        "包含数字的标志",
			flagName:    "flag123",
			shouldMark:  true,
			expected:    true,
			description: "包含数字的标志名",
		},
		{
			name:        "包含特殊字符的标志",
			flagName:    "flag_with_underscores",
			shouldMark:  true,
			expected:    true,
			description: "包含下划线的标志名",
		},
		{
			name:        "中文标志名",
			flagName:    "帮助",
			shouldMark:  true,
			expected:    true,
			description: "中文标志名",
		},
		{
			name:        "Unicode标志名",
			flagName:    "🚀flag",
			shouldMark:  true,
			expected:    true,
			description: "包含Unicode字符的标志名",
		},
		{
			name:        "极长标志名",
			flagName:    strings.Repeat("a", 10000),
			shouldMark:  true,
			expected:    true,
			description: "极长的标志名",
		},
		{
			name:        "包含空格的标志名",
			flagName:    "flag with spaces",
			shouldMark:  true,
			expected:    true,
			description: "包含空格的标志名",
		},
		{
			name:        "只有空格的标志名",
			flagName:    "   ",
			shouldMark:  true,
			expected:    true,
			description: "只包含空格的标志名",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 如果需要标记为内置标志，先标记
			if tt.shouldMark {
				bf.MarkAsBuiltin(tt.flagName)
			}

			// 测试IsBuiltinFlag
			result := bf.IsBuiltinFlag(tt.flagName)
			if result != tt.expected {
				t.Errorf("IsBuiltinFlag(%q) = %v, 期望 %v", tt.flagName, result, tt.expected)
			}
		})
	}
}

// TestBuiltinFlags_MarkAsBuiltin_基本功能 测试MarkAsBuiltin的基本功能
func TestBuiltinFlags_MarkAsBuiltin_基本功能(t *testing.T) {
	bf := NewBuiltinFlags()

	// 测试标记单个标志
	bf.MarkAsBuiltin("help")
	if !bf.IsBuiltinFlag("help") {
		t.Error("标记单个标志失败")
	}

	// 测试标记多个标志
	flags := []string{"version", "v", "completion", "c"}
	bf.MarkAsBuiltin(flags...)

	for _, flag := range flags {
		if !bf.IsBuiltinFlag(flag) {
			t.Errorf("标记多个标志失败: %q", flag)
		}
	}

	// 测试重复标记
	bf.MarkAsBuiltin("help") // 重复标记
	if !bf.IsBuiltinFlag("help") {
		t.Error("重复标记后标志丢失")
	}
}

// TestBuiltinFlags_MarkAsBuiltin_边界场景 测试MarkAsBuiltin的边界场景
func TestBuiltinFlags_MarkAsBuiltin_边界场景(t *testing.T) {
	bf := NewBuiltinFlags()

	// 测试空参数列表
	bf.MarkAsBuiltin()
	// 应该不会panic，也不会有任何效果

	// 测试包含空字符串的参数列表
	bf.MarkAsBuiltin("valid", "", "also-valid")
	if !bf.IsBuiltinFlag("valid") {
		t.Error("包含空字符串时，有效标志应该被正确标记")
	}
	if !bf.IsBuiltinFlag("also-valid") {
		t.Error("包含空字符串时，有效标志应该被正确标记")
	}
	// 空字符串本身不应该被标记为内置标志（根据IsBuiltinFlag的逻辑）
	if bf.IsBuiltinFlag("") {
		t.Error("空字符串不应该被标记为内置标志")
	}

	// 测试大量标志
	manyFlags := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		manyFlags[i] = "flag" + string(rune(i))
	}
	bf.MarkAsBuiltin(manyFlags...)

	// 验证部分标志
	testIndices := []int{0, 100, 1000, 5000, 9999}
	for _, idx := range testIndices {
		flagName := manyFlags[idx]
		if !bf.IsBuiltinFlag(flagName) {
			t.Errorf("大量标志中的第%d个标志 %q 未被正确标记", idx, flagName)
		}
	}
}

// TestBuiltinFlags_并发安全性 测试BuiltinFlags的并发安全性
func TestBuiltinFlags_并发安全性(t *testing.T) {
	bf := NewBuiltinFlags()

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// 测试并发标记
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				flagName := "flag_" + string(rune(id)) + "_" + string(rune(j))
				bf.MarkAsBuiltin(flagName)
			}
		}(i)
	}

	// 测试并发查询
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				flagName := "flag_" + string(rune(id)) + "_" + string(rune(j))
				_ = bf.IsBuiltinFlag(flagName)
			}
		}(i)
	}

	// 测试并发混合操作
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				if j%2 == 0 {
					bf.MarkAsBuiltin("concurrent_flag_" + string(rune(id)))
				} else {
					_ = bf.IsBuiltinFlag("concurrent_flag_" + string(rune(id)))
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证并发操作后的状态一致性
	for i := 0; i < numGoroutines; i++ {
		flagName := "concurrent_flag_" + string(rune(i))
		if !bf.IsBuiltinFlag(flagName) {
			t.Errorf("并发操作后标志 %q 丢失", flagName)
		}
	}

	t.Log("并发安全性测试完成")
}

// TestBuiltinFlags_内存使用 测试内存使用情况
func TestBuiltinFlags_内存使用(t *testing.T) {
	bf := NewBuiltinFlags()

	// 添加大量标志
	numFlags := 100000
	for i := 0; i < numFlags; i++ {
		flagName := "memory_test_flag_" + string(rune(i%1000)) + "_" + string(rune(i/1000))
		bf.MarkAsBuiltin(flagName)
	}

	// 验证所有标志都能正确查询
	successCount := 0
	for i := 0; i < numFlags; i++ {
		flagName := "memory_test_flag_" + string(rune(i%1000)) + "_" + string(rune(i/1000))
		if bf.IsBuiltinFlag(flagName) {
			successCount++
		}
	}

	if successCount != numFlags {
		t.Errorf("内存测试失败: 期望 %d 个标志, 实际找到 %d 个", numFlags, successCount)
	}

	t.Logf("内存使用测试完成，成功处理了 %d 个标志", numFlags)
}

// TestBuiltinFlags_极值测试 测试极值情况
func TestBuiltinFlags_极值测试(t *testing.T) {
	bf := NewBuiltinFlags()

	tests := []struct {
		name        string
		flagName    string
		description string
	}{
		{
			name:        "极长标志名",
			flagName:    strings.Repeat("a", 100000),
			description: "测试极长的标志名",
		},
		{
			name:        "单字符标志",
			flagName:    "a",
			description: "测试单字符标志",
		},
		{
			name:        "包含所有ASCII字符",
			flagName:    "!@#$%^&*()_+-={}[]|\\:;\"'<>?,./~`",
			description: "测试包含特殊ASCII字符的标志名",
		},
		{
			name:        "Unicode字符",
			flagName:    "测试标志🚀✨🎉",
			description: "测试Unicode字符标志名",
		},
		{
			name:        "包含换行符",
			flagName:    "flag\nwith\nnewlines",
			description: "测试包含换行符的标志名",
		},
		{
			name:        "包含制表符",
			flagName:    "flag\twith\ttabs",
			description: "测试包含制表符的标志名",
		},
		{
			name:        "只有空格",
			flagName:    "     ",
			description: "测试只包含空格的标志名",
		},
		{
			name:        "混合空白字符",
			flagName:    " \t\n\r ",
			description: "测试混合空白字符的标志名",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 标记为内置标志
			bf.MarkAsBuiltin(tt.flagName)

			// 验证能够正确识别（除了空字符串情况）
			expected := tt.flagName != ""
			if bf.IsBuiltinFlag(tt.flagName) != expected {
				t.Errorf("极值测试失败 %q: 期望 %v, 实际 %v",
					tt.flagName, expected, bf.IsBuiltinFlag(tt.flagName))
			}
		})
	}
}

// TestBuiltinFlags_NameMap_直接操作 测试直接操作NameMap的行为
func TestBuiltinFlags_NameMap_直接操作(t *testing.T) {
	bf := NewBuiltinFlags()

	// 直接向NameMap添加数据
	bf.NameMap.Store("direct_flag", true)
	bf.NameMap.Store("another_flag", "not_bool_value")
	bf.NameMap.Store(123, true) // 非字符串键

	// 测试IsBuiltinFlag的行为
	if !bf.IsBuiltinFlag("direct_flag") {
		t.Error("直接添加到NameMap的字符串标志应该被识别")
	}

	if !bf.IsBuiltinFlag("another_flag") {
		t.Error("直接添加到NameMap的标志应该被识别，无论值的类型")
	}

	// 测试非字符串键不会影响字符串查询
	if bf.IsBuiltinFlag("123") {
		t.Error("非字符串键不应该影响字符串查询")
	}

	// 测试删除操作
	bf.NameMap.Delete("direct_flag")
	if bf.IsBuiltinFlag("direct_flag") {
		t.Error("删除后的标志不应该被识别")
	}
}

// TestBuiltinFlags_字段类型验证 测试字段类型的正确性
func TestBuiltinFlags_字段类型验证(t *testing.T) {
	bf := NewBuiltinFlags()

	// 验证Help字段
	if bf.Help == nil {
		t.Error("Help字段不应该为nil")
	}

	helpType := reflect.TypeOf(bf.Help)
	expectedHelpType := reflect.TypeOf(&flags.BoolFlag{})
	if helpType != expectedHelpType {
		t.Errorf("Help字段类型不正确: 期望 %v, 实际 %v", expectedHelpType, helpType)
	}

	// 验证Version字段
	if bf.Version == nil {
		t.Error("Version字段不应该为nil")
	}

	versionType := reflect.TypeOf(bf.Version)
	expectedVersionType := reflect.TypeOf(&flags.BoolFlag{})
	if versionType != expectedVersionType {
		t.Errorf("Version字段类型不正确: 期望 %v, 实际 %v", expectedVersionType, versionType)
	}

	// 验证Completion字段
	if bf.Completion == nil {
		t.Error("Completion字段不应该为nil")
	}

	completionType := reflect.TypeOf(bf.Completion)
	expectedCompletionType := reflect.TypeOf(&flags.EnumFlag{})
	if completionType != expectedCompletionType {
		t.Errorf("Completion字段类型不正确: 期望 %v, 实际 %v", expectedCompletionType, completionType)
	}

	// 验证NameMap字段 - 使用指针避免copylocks警告
	nameMapType := reflect.TypeOf(&bf.NameMap).Elem()
	expectedNameMapType := reflect.TypeOf((*sync.Map)(nil)).Elem()
	if nameMapType != expectedNameMapType {
		t.Errorf("NameMap字段类型不正确: 期望 %v, 实际 %v", expectedNameMapType, nameMapType)
	}
}

// TestBuiltinFlags_多实例独立性 测试多个BuiltinFlags实例的独立性
func TestBuiltinFlags_多实例独立性(t *testing.T) {
	bf1 := NewBuiltinFlags()
	bf2 := NewBuiltinFlags()

	// 在第一个实例中标记标志
	bf1.MarkAsBuiltin("flag1", "flag2")

	// 在第二个实例中标记不同的标志
	bf2.MarkAsBuiltin("flag3", "flag4")

	// 验证实例间的独立性
	if !bf1.IsBuiltinFlag("flag1") {
		t.Error("bf1应该包含flag1")
	}
	if !bf1.IsBuiltinFlag("flag2") {
		t.Error("bf1应该包含flag2")
	}
	if bf1.IsBuiltinFlag("flag3") {
		t.Error("bf1不应该包含flag3")
	}
	if bf1.IsBuiltinFlag("flag4") {
		t.Error("bf1不应该包含flag4")
	}

	if bf2.IsBuiltinFlag("flag1") {
		t.Error("bf2不应该包含flag1")
	}
	if bf2.IsBuiltinFlag("flag2") {
		t.Error("bf2不应该包含flag2")
	}
	if !bf2.IsBuiltinFlag("flag3") {
		t.Error("bf2应该包含flag3")
	}
	if !bf2.IsBuiltinFlag("flag4") {
		t.Error("bf2应该包含flag4")
	}
}

// TestBuiltinFlags_性能测试 测试性能表现
func TestBuiltinFlags_性能测试(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	bf := NewBuiltinFlags()

	// 预先添加一些标志
	numPreFlags := 1000
	preFlags := make([]string, numPreFlags)
	for i := 0; i < numPreFlags; i++ {
		preFlags[i] = "perf_flag_" + string(rune(48+i%10)) + string(rune(48+(i/10)%10)) + string(rune(48+(i/100)%10))
	}
	bf.MarkAsBuiltin(preFlags...)

	// 测试查询性能
	numQueries := 100000
	existingFlag := preFlags[500] // 使用实际存在的标志
	nonExistingFlag := "non_existing_flag"

	// 测试存在标志的查询性能
	for i := 0; i < numQueries; i++ {
		if !bf.IsBuiltinFlag(existingFlag) {
			t.Errorf("性能测试中存在的标志查询失败，标志: %s", existingFlag)
			break
		}
	}

	// 测试不存在标志的查询性能
	for i := 0; i < numQueries; i++ {
		if bf.IsBuiltinFlag(nonExistingFlag) {
			t.Error("性能测试中不存在的标志查询错误")
		}
	}

	// 测试标记性能
	numMarkOperations := 10000
	for i := 0; i < numMarkOperations; i++ {
		bf.MarkAsBuiltin("mark_perf_flag_" + string(rune(i)))
	}

	t.Logf("性能测试完成: %d 次查询, %d 次标记操作", numQueries*2, numMarkOperations)
}

// TestBuiltinFlags_边界条件组合 测试各种边界条件的组合
func TestBuiltinFlags_边界条件组合(t *testing.T) {
	bf := NewBuiltinFlags()

	// 组合测试：空字符串 + 正常标志 + 特殊字符
	testFlags := []string{
		"",                        // 空字符串
		"normal",                  // 正常标志
		"flag-with-hyphens",       // 包含连字符
		"flag_with_underscores",   // 包含下划线
		"123numeric",              // 以数字开头
		"MixedCase",               // 混合大小写
		"中文标志",                    // 中文
		"🚀emoji",                  // emoji
		strings.Repeat("x", 1000), // 长标志
	}

	// 批量标记
	bf.MarkAsBuiltin(testFlags...)

	// 验证每个标志（除了空字符串）
	for _, flag := range testFlags {
		expected := flag != ""
		if bf.IsBuiltinFlag(flag) != expected {
			t.Errorf("组合测试失败 %q: 期望 %v, 实际 %v",
				flag, expected, bf.IsBuiltinFlag(flag))
		}
	}

	// 测试未标记的类似标志
	similarFlags := []string{
		"Normal",                 // 大小写不同
		"flag-with-hyphens-",     // 多一个字符
		"flag_with_underscores_", // 多一个字符
		"124numeric",             // 数字不同
		"中文标记",                   // 中文不同
		"🎉emoji",                 // emoji不同
	}

	for _, flag := range similarFlags {
		if bf.IsBuiltinFlag(flag) {
			t.Errorf("未标记的类似标志 %q 不应该被识别为内置标志", flag)
		}
	}
}
