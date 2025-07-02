package flags

import (
	"reflect"
	"testing"
)

// TestSliceFlag_BasicParsing 测试基本切片解析功能
func TestSliceFlag_BasicParsing(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	// 测试正常分割
	if err := flag.Set("a,b,c"); err != nil {
		t.Errorf("Set failed: %v", err)
	}
	result := flag.Get()
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// 测试无分隔符情况
	if err := flag.Set("d"); err != nil {
		t.Errorf("Set failed: %v", err)
	}
	result = flag.Get()
	expected = []string{"d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	result = flag.Get()
	expected = []string{"d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestSliceFlag_SkipEmpty 测试空元素过滤功能
func TestSliceFlag_SkipEmpty(t *testing.T) {
	// 测试SkipEmpty=true情况
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
		skipEmpty:  true,
	}

	if err := flag.Set("a,,b,,c"); err != nil {
		t.Errorf("Set failed: %v", err)
	}
	result := flag.Get()
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// 测试SkipEmpty=false情况
	flag = &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
		skipEmpty:  false,
	}

	if err := flag.Set("a,,b,,c"); err != nil {
		t.Errorf("Set failed: %v", err)
	}
	result = flag.Get()
	expected = []string{"a", "", "b", "", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestSliceFlag_LenAndContains 测试Len和Contains方法
func TestSliceFlag_LenAndContains(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{"x", "y"},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	// 测试Len
	if flag.Len() != 2 {
		t.Errorf("Expected length 2, got %d", flag.Len())
	}

	// 设置值后测试
	if err := flag.Set("a,b,c"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if flag.Len() != 3 {
		t.Errorf("Expected length 3, got %d", flag.Len())
	}

	// 测试Contains
	if !flag.Contains("b") {
		t.Error("Should contain 'b'")
	}
	if flag.Contains("z") {
		t.Error("Should not contain 'z'")
	}
}

// TestSliceFlag_ClearAndRemove 测试Clear和Remove方法
func TestSliceFlag_ClearAndRemove(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	// 设置初始值
	if err := flag.Set("a,b,c,d"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// 测试Remove
	if err := flag.Remove("b"); err != nil {
		t.Errorf("Remove failed: %v", err)
	}
	result := flag.Get()
	expected := []string{"a", "c", "d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After Remove, expected %v, got %v", expected, result)
	}

	// 测试Clear
	if err := flag.Clear(); err != nil {
		t.Errorf("Clear failed: %v", err)
	}
	if flag.Len() != 0 {
		t.Errorf("After Clear, expected length 0, got %d", flag.Len())
	}
}

// TestSliceFlag_Sort 测试Sort方法
func TestSliceFlag_Sort(t *testing.T) {
	flag := &SliceFlag{
		BaseFlag: BaseFlag[[]string]{
			initialValue: []string{},
			value:        new([]string),
		},
		delimiters: []string{","},
	}

	if err := flag.Set("c,a,b"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if err := flag.Sort(); err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	result := flag.Get()
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("After Sort, expected %v, got %v", expected, result)
	}
}
