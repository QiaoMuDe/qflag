package flags

import (
	"reflect"
	"testing"
)

// =============================================================================
// IntSliceFlag 测试用例
// =============================================================================

func TestIntSliceFlag_BasicParsing(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试基本解析
	err = flag.Set("1,2,3,4")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestIntSliceFlag_MultipleDelimiters(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{"comma", "10,20,30", []int{10, 20, 30}},
		{"semicolon", "100;200;300", []int{100, 200, 300}},
		{"pipe", "1000|2000|3000", []int{1000, 2000, 3000}},
		{"single_value", "42", []int{42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("Set failed for %s: %v", tt.name, err)
			}

			actual := flag.Get()
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestIntSliceFlag_WithSpaces(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试带空格的输入
	err = flag.Set(" 1 , 2 , 3 , 4 ")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestIntSliceFlag_SkipEmpty(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 启用跳过空元素
	flag.SetSkipEmpty(true)

	err = flag.Set("1,,2,,3")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int{1, 2, 3}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestIntSliceFlag_InvalidValues(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{"empty_string", ""},
		{"non_integer", "abc"},
		{"mixed_valid_invalid", "1,abc,3"},
		{"float_number", "1.5,2.5"},
		{"overflow", "999999999999999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flag.Set(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", tt.input)
			}
		})
	}
}

func TestIntSliceFlag_NegativeNumbers(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	err = flag.Set("-1,-2,3,-4")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int{-1, -2, 3, -4}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestIntSliceFlag_HelperMethods(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{1, 2, 3}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试 Len
	if flag.Len() != 3 {
		t.Errorf("Expected length 3, got %d", flag.Len())
	}

	// 测试 Contains
	if !flag.Contains(2) {
		t.Error("Expected to contain 2")
	}
	if flag.Contains(5) {
		t.Error("Expected not to contain 5")
	}

	// 测试 Remove
	err = flag.Remove(2)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	expected := []int{1, 3}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After remove, expected %v, got %v", expected, actual)
	}

	// 测试 Clear
	err = flag.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	if flag.Len() != 0 {
		t.Errorf("After clear, expected length 0, got %d", flag.Len())
	}
}

func TestIntSliceFlag_Sort(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{3, 1, 4, 1, 5}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	err = flag.Sort()
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	expected := []int{1, 1, 3, 4, 5}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After sort, expected %v, got %v", expected, actual)
	}
}

func TestIntSliceFlag_String(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{1, 2, 3}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	expected := "1,2,3"
	actual := flag.String()
	if actual != expected {
		t.Errorf("Expected string %s, got %s", expected, actual)
	}
}

func TestIntSliceFlag_Type(t *testing.T) {
	flag := &IntSliceFlag{}
	if flag.Type() != FlagTypeIntSlice {
		t.Errorf("Expected type %v, got %v", FlagTypeIntSlice, flag.Type())
	}
}

func TestIntSliceFlag_CustomDelimiters(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 设置自定义分隔符
	flag.SetDelimiters([]string{":"})

	err = flag.Set("1:2:3:4")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	// 测试获取分隔符
	delimiters := flag.GetDelimiters()
	if len(delimiters) != 1 || delimiters[0] != ":" {
		t.Errorf("Expected delimiters [\":\"], got %v", delimiters)
	}
}

// =============================================================================
// Int64SliceFlag 测试用例
// =============================================================================

func TestInt64SliceFlag_BasicParsing(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试基本解析
	err = flag.Set("1,2,3,4")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int64{1, 2, 3, 4}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestInt64SliceFlag_LargeNumbers(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试大数值
	err = flag.Set("9223372036854775807,-9223372036854775808")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int64{9223372036854775807, -9223372036854775808}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestInt64SliceFlag_MultipleDelimiters(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected []int64
	}{
		{"comma", "10,20,30", []int64{10, 20, 30}},
		{"semicolon", "100;200;300", []int64{100, 200, 300}},
		{"pipe", "1000|2000|3000", []int64{1000, 2000, 3000}},
		{"single_value", "42", []int64{42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("Set failed for %s: %v", tt.name, err)
			}

			actual := flag.Get()
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestInt64SliceFlag_InvalidValues(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{"empty_string", ""},
		{"non_integer", "abc"},
		{"mixed_valid_invalid", "1,abc,3"},
		{"float_number", "1.5,2.5"},
		{"overflow", "99999999999999999999999999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flag.Set(tt.input)
			if err == nil {
				t.Errorf("Expected error for input %s, but got none", tt.input)
			}
		})
	}
}

func TestInt64SliceFlag_HelperMethods(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{1, 2, 3}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试 Len
	if flag.Len() != 3 {
		t.Errorf("Expected length 3, got %d", flag.Len())
	}

	// 测试 Contains
	if !flag.Contains(2) {
		t.Error("Expected to contain 2")
	}
	if flag.Contains(5) {
		t.Error("Expected not to contain 5")
	}

	// 测试 Remove
	err = flag.Remove(2)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	expected := []int64{1, 3}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After remove, expected %v, got %v", expected, actual)
	}

	// 测试 Clear
	err = flag.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	if flag.Len() != 0 {
		t.Errorf("After clear, expected length 0, got %d", flag.Len())
	}
}

func TestInt64SliceFlag_Sort(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{3, 1, 4, 1, 5}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	err = flag.Sort()
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	expected := []int64{1, 1, 3, 4, 5}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("After sort, expected %v, got %v", expected, actual)
	}
}

func TestInt64SliceFlag_String(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{1, 2, 3}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	expected := "1,2,3"
	actual := flag.String()
	if actual != expected {
		t.Errorf("Expected string %s, got %s", expected, actual)
	}
}

func TestInt64SliceFlag_Type(t *testing.T) {
	flag := &Int64SliceFlag{}
	if flag.Type() != FlagTypeInt64Slice {
		t.Errorf("Expected type %v, got %v", FlagTypeInt64Slice, flag.Type())
	}
}

func TestInt64SliceFlag_SkipEmpty(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 启用跳过空元素
	flag.SetSkipEmpty(true)

	err = flag.Set("1,,2,,3")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	expected := []int64{1, 2, 3}
	actual := flag.Get()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// =============================================================================
// 边界测试和并发测试
// =============================================================================

func TestIntSliceFlag_ConcurrentAccess(t *testing.T) {
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", []int{1, 2, 3}, "integer list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 并发读取测试
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 100; j++ {
				_ = flag.Get()
				_ = flag.Contains(2)
				_ = flag.Len()
				_ = flag.String()
			}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestInt64SliceFlag_ConcurrentAccess(t *testing.T) {
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", []int64{1, 2, 3}, "int64 list")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 并发读取测试
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 100; j++ {
				_ = flag.Get()
				_ = flag.Contains(2)
				_ = flag.Len()
				_ = flag.String()
			}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestIntSliceFlag_EdgeCases(t *testing.T) {
	// 测试nil默认值
	flag := &IntSliceFlag{}
	err := flag.Init("numbers", "n", nil, "integer list")
	if err != nil {
		t.Fatalf("Init with nil default failed: %v", err)
	}

	if flag.Len() != 0 {
		t.Errorf("Expected empty slice for nil default, got length %d", flag.Len())
	}

	// 测试空分隔符
	flag.SetDelimiters([]string{})
	delimiters := flag.GetDelimiters()
	if len(delimiters) == 0 {
		t.Error("Expected default delimiters when setting empty delimiters")
	}
}

func TestInt64SliceFlag_EdgeCases(t *testing.T) {
	// 测试nil默认值
	flag := &Int64SliceFlag{}
	err := flag.Init("numbers", "n", nil, "int64 list")
	if err != nil {
		t.Fatalf("Init with nil default failed: %v", err)
	}

	if flag.Len() != 0 {
		t.Errorf("Expected empty slice for nil default, got length %d", flag.Len())
	}

	// 测试空分隔符
	flag.SetDelimiters([]string{})
	delimiters := flag.GetDelimiters()
	if len(delimiters) == 0 {
		t.Error("Expected default delimiters when setting empty delimiters")
	}
}
