package qflag

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// =============================================================================
// IntSlice 方法测试用例
// =============================================================================

func TestCmd_IntSlice_BasicUsage(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试基本创建
	flag := cmd.IntSlice("ports", "p", []int{8080}, "server ports")

	//nolint:all
	if flag == nil {
		t.Fatal("IntSlice returned nil")
	}

	// 验证默认值
	expected := []int{8080}
	//nolint:all
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected default value %v, got %v", expected, actual)
	}

	// 验证标志名
	//nolint:all
	if flag.LongName() != "ports" {
		t.Errorf("Expected long name 'ports', got '%s'", flag.LongName())
	}
	if flag.ShortName() != "p" {
		t.Errorf("Expected short name 'p', got '%s'", flag.ShortName())
	}

	// 验证使用说明
	if flag.Usage() != "server ports" {
		t.Errorf("Expected usage 'server ports', got '%s'", flag.Usage())
	}

	// 验证类型
	if flag.Type() != flags.FlagTypeIntSlice {
		t.Errorf("Expected type %v, got %v", flags.FlagTypeIntSlice, flag.Type())
	}
}

func TestCmd_IntSlice_EmptyDefaults(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试空默认值
	flag := cmd.IntSlice("numbers", "n", []int{}, "number list")

	if flag.Len() != 0 {
		t.Errorf("Expected empty slice, got length %d", flag.Len())
	}

	// 测试nil默认值
	flag2 := cmd.IntSlice("values", "v", nil, "value list")
	if flag2.Len() != 0 {
		t.Errorf("Expected empty slice for nil default, got length %d", flag2.Len())
	}
}

func TestCmd_IntSlice_NoShortName(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试只有长标志名
	flag := cmd.IntSlice("long-only", "", []int{1, 2, 3}, "long only flag")

	if flag.ShortName() != "" {
		t.Errorf("Expected empty short name, got '%s'", flag.ShortName())
	}
	if flag.LongName() != "long-only" {
		t.Errorf("Expected long name 'long-only', got '%s'", flag.LongName())
	}
}

func TestCmd_IntSlice_Registration(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 创建标志
	flag := cmd.IntSlice("ports", "p", []int{8080}, "server ports")

	// 验证标志已注册到注册表 - 通过获取所有标志来验证
	allFlags := cmd.ctx.FlagRegistry.GetFlagMetaList()
	found := false
	for _, meta := range allFlags {
		if meta.Flag == flag {
			found = true
			break
		}
	}

	if !found {
		t.Error("Flag was not registered in the flag registry")
	}
}

// =============================================================================
// IntSliceVar 方法测试用例
// =============================================================================

func TestCmd_IntSliceVar_BasicUsage(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)
	var flag flags.IntSliceFlag

	// 测试基本绑定
	cmd.IntSliceVar(&flag, "ports", "p", []int{8080, 3000}, "server ports")

	// 验证默认值
	expected := []int{8080, 3000}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected default value %v, got %v", expected, actual)
	}

	// 验证标志属性
	if flag.LongName() != "ports" {
		t.Errorf("Expected long name 'ports', got '%s'", flag.LongName())
	}
	if flag.ShortName() != "p" {
		t.Errorf("Expected short name 'p', got '%s'", flag.ShortName())
	}
}

func TestCmd_IntSliceVar_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试nil指针应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil pointer, but didn't panic")
		}
	}()

	cmd.IntSliceVar(nil, "ports", "p", []int{8080}, "server ports")
}

func TestCmd_IntSliceVar_BuiltinFlagConflict(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)
	var flag flags.IntSliceFlag

	// 测试与内置标志冲突应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for builtin flag conflict, but didn't panic")
		}
	}()

	// 假设"help"是内置标志
	cmd.IntSliceVar(&flag, "help", "", []int{}, "help flag")
}

func TestCmd_IntSliceVar_NilDefault(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)
	var flag flags.IntSliceFlag

	// 测试nil默认值处理
	cmd.IntSliceVar(&flag, "numbers", "n", nil, "number list")

	if flag.Len() != 0 {
		t.Errorf("Expected empty slice for nil default, got length %d", flag.Len())
	}
}

// =============================================================================
// Int64Slice 方法测试用例
// =============================================================================

func TestCmd_Int64Slice_BasicUsage(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试基本创建
	flag := cmd.Int64Slice("sizes", "s", []int64{1024, 2048}, "file sizes")

	//nolint:all
	if flag == nil {
		t.Fatal("Int64Slice returned nil")
	}

	// 验证默认值
	expected := []int64{1024, 2048}
	//nolint:all
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected default value %v, got %v", expected, actual)
	}

	// 验证标志名
	//nolint:all
	if flag.LongName() != "sizes" {
		t.Errorf("Expected long name 'sizes', got '%s'", flag.LongName())
	}
	if flag.ShortName() != "s" {
		t.Errorf("Expected short name 's', got '%s'", flag.ShortName())
	}

	// 验证类型
	if flag.Type() != flags.FlagTypeInt64Slice {
		t.Errorf("Expected type %v, got %v", flags.FlagTypeInt64Slice, flag.Type())
	}
}

func TestCmd_Int64Slice_LargeNumbers(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试大数值
	largeNumbers := []int64{9223372036854775807, -9223372036854775808}
	flag := cmd.Int64Slice("big-numbers", "bn", largeNumbers, "big numbers")

	actual := flag.Get()
	if !reflect.DeepEqual(actual, largeNumbers) {
		t.Errorf("Expected large numbers %v, got %v", largeNumbers, actual)
	}
}

func TestCmd_Int64Slice_EmptyDefaults(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试空默认值
	flag := cmd.Int64Slice("numbers", "n", []int64{}, "number list")

	if flag.Len() != 0 {
		t.Errorf("Expected empty slice, got length %d", flag.Len())
	}

	// 测试nil默认值
	flag2 := cmd.Int64Slice("values", "v", nil, "value list")
	if flag2.Len() != 0 {
		t.Errorf("Expected empty slice for nil default, got length %d", flag2.Len())
	}
}

func TestCmd_Int64Slice_Registration(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 创建标志
	flag := cmd.Int64Slice("sizes", "s", []int64{1024}, "file sizes")

	// 验证标志已注册到注册表 - 通过获取所有标志来验证
	allFlags := cmd.ctx.FlagRegistry.GetFlagMetaList()
	found := false
	for _, meta := range allFlags {
		if meta.Flag == flag {
			found = true
			break
		}
	}

	if !found {
		t.Error("Flag was not registered in the flag registry")
	}
}

// =============================================================================
// Int64SliceVar 方法测试用例
// =============================================================================

func TestCmd_Int64SliceVar_BasicUsage(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)
	var flag flags.Int64SliceFlag

	// 测试基本绑定
	cmd.Int64SliceVar(&flag, "sizes", "s", []int64{1024, 2048, 4096}, "file sizes")

	// 验证默认值
	expected := []int64{1024, 2048, 4096}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected default value %v, got %v", expected, actual)
	}

	// 验证标志属性
	if flag.LongName() != "sizes" {
		t.Errorf("Expected long name 'sizes', got '%s'", flag.LongName())
	}
	if flag.ShortName() != "s" {
		t.Errorf("Expected short name 's', got '%s'", flag.ShortName())
	}
}

func TestCmd_Int64SliceVar_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试nil指针应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil pointer, but didn't panic")
		}
	}()

	cmd.Int64SliceVar(nil, "sizes", "s", []int64{1024}, "file sizes")
}

func TestCmd_Int64SliceVar_BuiltinFlagConflict(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)
	var flag flags.Int64SliceFlag

	// 测试与内置标志冲突应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for builtin flag conflict, but didn't panic")
		}
	}()

	// 假设"help"是内置标志
	cmd.Int64SliceVar(&flag, "help", "", []int64{}, "help flag")
}

func TestCmd_Int64SliceVar_NilDefault(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)
	var flag flags.Int64SliceFlag

	// 测试nil默认值处理
	cmd.Int64SliceVar(&flag, "numbers", "n", nil, "number list")

	if flag.Len() != 0 {
		t.Errorf("Expected empty slice for nil default, got length %d", flag.Len())
	}
}

// =============================================================================
// 综合测试和边界测试
// =============================================================================

func TestCmd_SliceFlags_MultipleFlags(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 创建多个不同类型的切片标志
	intFlag := cmd.IntSlice("ports", "p", []int{8080}, "server ports")
	int64Flag := cmd.Int64Slice("sizes", "s", []int64{1024}, "file sizes")

	// 验证它们都被正确创建
	if intFlag == nil || int64Flag == nil {
		t.Fatal("One or more flags were not created")
	}

	// 验证它们有不同的类型
	if intFlag.Type() == int64Flag.Type() {
		t.Error("IntSlice and Int64Slice should have different types")
	}

	// 验证注册表中有两个标志 - 通过获取所有标志验证
	allFlags := cmd.ctx.FlagRegistry.GetFlagMetaList()
	intFlagFound := false
	int64FlagFound := false

	for _, meta := range allFlags {
		if meta.Flag == intFlag {
			intFlagFound = true
		}
		if meta.Flag == int64Flag {
			int64FlagFound = true
		}
	}

	if !intFlagFound {
		t.Error("IntSlice flag was not registered")
	}
	if !int64FlagFound {
		t.Error("Int64Slice flag was not registered")
	}
}

func TestCmd_SliceFlags_NameConflict(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 创建第一个标志
	cmd.IntSlice("numbers", "n", []int{1}, "first flag")

	// 尝试创建同名标志应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for name conflict, but didn't panic")
		}
	}()

	cmd.Int64Slice("numbers", "n", []int64{2}, "second flag")
}

func TestCmd_SliceFlags_OnlyLongName(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试只有长标志名的情况
	intFlag := cmd.IntSlice("int-only", "", []int{1}, "int only")
	int64Flag := cmd.Int64Slice("int64-only", "", []int64{2}, "int64 only")

	if intFlag.ShortName() != "" {
		t.Error("Expected empty short name for int flag")
	}
	if int64Flag.ShortName() != "" {
		t.Error("Expected empty short name for int64 flag")
	}
}

func TestCmd_SliceFlags_OnlyShortName(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试只有短标志名的情况
	intFlag := cmd.IntSlice("", "i", []int{1}, "int short")
	int64Flag := cmd.Int64Slice("", "l", []int64{2}, "int64 short")

	if intFlag.LongName() != "" {
		t.Error("Expected empty long name for int flag")
	}
	if int64Flag.LongName() != "" {
		t.Error("Expected empty long name for int64 flag")
	}
}

func TestCmd_SliceFlags_ConcurrentCreation(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 并发创建标志测试
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			defer func() { done <- true }()

			// 创建不同名称的标志避免冲突
			intFlagName := fmt.Sprintf("int-flag-%d", index)
			int64FlagName := fmt.Sprintf("int64-flag-%d", index)

			intFlag := cmd.IntSlice(intFlagName, "", []int{index}, "int flag")
			int64Flag := cmd.Int64Slice(int64FlagName, "", []int64{int64(index)}, "int64 flag")

			if intFlag == nil || int64Flag == nil {
				t.Errorf("Failed to create flags in goroutine %d", index)
			}
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有标志都被注册 - 通过标志数量验证
	allFlags := cmd.ctx.FlagRegistry.GetFlagMetaList()
	if len(allFlags) < 20 {
		t.Errorf("Expected at least 20 registered flags, got %d", len(allFlags))
	}
}

// =============================================================================
// 功能集成测试
// =============================================================================

func TestCmd_SliceFlags_SetAndGet(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 创建标志
	intFlag := cmd.IntSlice("ports", "p", []int{8080}, "server ports")
	int64Flag := cmd.Int64Slice("sizes", "s", []int64{1024}, "file sizes")

	// 测试设置值
	err := intFlag.Set("80,443,8080")
	if err != nil {
		t.Fatalf("Failed to set int slice: %v", err)
	}

	err = int64Flag.Set("1024,2048,4096")
	if err != nil {
		t.Fatalf("Failed to set int64 slice: %v", err)
	}

	// 验证值
	expectedInt := []int{80, 443, 8080}
	actualInt := intFlag.Get()
	if !reflect.DeepEqual(actualInt, expectedInt) {
		t.Errorf("Expected int slice %v, got %v", expectedInt, actualInt)
	}

	expectedInt64 := []int64{1024, 2048, 4096}
	actualInt64 := int64Flag.Get()
	if !reflect.DeepEqual(actualInt64, expectedInt64) {
		t.Errorf("Expected int64 slice %v, got %v", expectedInt64, actualInt64)
	}
}

func TestCmd_SliceFlags_HelperMethods(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 创建标志
	intFlag := cmd.IntSlice("numbers", "n", []int{1, 2, 3, 4, 5}, "numbers")
	int64Flag := cmd.Int64Slice("values", "v", []int64{10, 20, 30}, "values")

	// 测试辅助方法
	if intFlag.Len() != 5 {
		t.Errorf("Expected int slice length 5, got %d", intFlag.Len())
	}
	if int64Flag.Len() != 3 {
		t.Errorf("Expected int64 slice length 3, got %d", int64Flag.Len())
	}

	if !intFlag.Contains(3) {
		t.Error("Expected int slice to contain 3")
	}
	if !int64Flag.Contains(20) {
		t.Error("Expected int64 slice to contain 20")
	}

	// 测试排序
	_ = intFlag.Set("5,1,3,2,4")
	err := intFlag.Sort()
	if err != nil {
		t.Fatalf("Failed to sort int slice: %v", err)
	}

	expectedSorted := []int{1, 2, 3, 4, 5}
	actualSorted := intFlag.Get()
	if !reflect.DeepEqual(actualSorted, expectedSorted) {
		t.Errorf("Expected sorted slice %v, got %v", expectedSorted, actualSorted)
	}
}

// =============================================================================
// 错误处理和边界测试
// =============================================================================

func TestCmd_SliceFlags_InvalidInput(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	intFlag := cmd.IntSlice("numbers", "n", []int{}, "numbers")
	int64Flag := cmd.Int64Slice("values", "v", []int64{}, "values")

	// 测试无效输入
	tests := []struct {
		name  string
		input string
	}{
		{"empty_string", ""},
		{"non_numeric", "abc,def"},
		{"mixed_valid_invalid", "1,abc,3"},
		{"float_numbers", "1.5,2.5"},
	}

	for _, tt := range tests {
		t.Run("int_"+tt.name, func(t *testing.T) {
			err := intFlag.Set(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", tt.input)
			}
		})

		t.Run("int64_"+tt.name, func(t *testing.T) {
			err := int64Flag.Set(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", tt.input)
			}
		})
	}
}

func TestCmd_SliceFlags_BoundaryValues(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	// 测试边界值
	intFlag := cmd.IntSlice("int-bounds", "", []int{}, "int bounds")
	int64Flag := cmd.Int64Slice("int64-bounds", "", []int64{}, "int64 bounds")

	// 测试int边界值
	err := intFlag.Set("2147483647,-2147483648")
	if err != nil {
		t.Errorf("Failed to set int boundary values: %v", err)
	}

	expectedInt := []int{2147483647, -2147483648}
	actualInt := intFlag.Get()
	if !reflect.DeepEqual(actualInt, expectedInt) {
		t.Errorf("Expected int boundary values %v, got %v", expectedInt, actualInt)
	}

	// 测试int64边界值
	err = int64Flag.Set("9223372036854775807,-9223372036854775808")
	if err != nil {
		t.Errorf("Failed to set int64 boundary values: %v", err)
	}

	expectedInt64 := []int64{9223372036854775807, -9223372036854775808}
	actualInt64 := int64Flag.Get()
	if !reflect.DeepEqual(actualInt64, expectedInt64) {
		t.Errorf("Expected int64 boundary values %v, got %v", expectedInt64, actualInt64)
	}
}

func TestCmd_SliceFlags_CustomDelimiters(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	intFlag := cmd.IntSlice("numbers", "n", []int{}, "numbers")
	int64Flag := cmd.Int64Slice("values", "v", []int64{}, "values")

	// 设置自定义分隔符
	intFlag.SetDelimiters([]string{":"})
	int64Flag.SetDelimiters([]string{";"})

	// 测试自定义分隔符解析
	err := intFlag.Set("1:2:3:4")
	if err != nil {
		t.Errorf("Failed to parse with custom delimiter: %v", err)
	}

	err = int64Flag.Set("10;20;30;40")
	if err != nil {
		t.Errorf("Failed to parse int64 with custom delimiter: %v", err)
	}

	expectedInt := []int{1, 2, 3, 4}
	actualInt := intFlag.Get()
	if !reflect.DeepEqual(actualInt, expectedInt) {
		t.Errorf("Expected %v with custom delimiter, got %v", expectedInt, actualInt)
	}

	expectedInt64 := []int64{10, 20, 30, 40}
	actualInt64 := int64Flag.Get()
	if !reflect.DeepEqual(actualInt64, expectedInt64) {
		t.Errorf("Expected %v with custom delimiter, got %v", expectedInt64, actualInt64)
	}
}

func TestCmd_SliceFlags_SkipEmpty(t *testing.T) {
	cmd := NewCmd("test", "Test command", flag.ExitOnError)

	intFlag := cmd.IntSlice("numbers", "n", []int{}, "numbers")
	int64Flag := cmd.Int64Slice("values", "v", []int64{}, "values")

	// 启用跳过空元素
	intFlag.SetSkipEmpty(true)
	int64Flag.SetSkipEmpty(true)

	// 测试跳过空元素
	err := intFlag.Set("1,,2,,3")
	if err != nil {
		t.Errorf("Failed to parse with skip empty: %v", err)
	}

	err = int64Flag.Set("10,,20,,30")
	if err != nil {
		t.Errorf("Failed to parse int64 with skip empty: %v", err)
	}

	expectedInt := []int{1, 2, 3}
	actualInt := intFlag.Get()
	if !reflect.DeepEqual(actualInt, expectedInt) {
		t.Errorf("Expected %v with skip empty, got %v", expectedInt, actualInt)
	}

	expectedInt64 := []int64{10, 20, 30}
	actualInt64 := int64Flag.Get()
	if !reflect.DeepEqual(actualInt64, expectedInt64) {
		t.Errorf("Expected %v with skip empty, got %v", expectedInt64, actualInt64)
	}
}
