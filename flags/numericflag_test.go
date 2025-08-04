package flags

import (
	"fmt"
	"testing"
)

// TestIntFlag_BasicFunctionality 测试IntFlag的基本功能
func TestIntFlag_BasicFunctionality(t *testing.T) {
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			initialValue: 0,
			value:        new(int),
		},
	}

	// 测试默认值
	if flag.Get() != 0 {
		t.Error("默认值应为0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []int{10, -5, 0, 1000}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%d失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%d, 实际值%d", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0 {
		t.Error("重置后应返回默认值0")
	}
}

// TestIntFlag_RangeValidation 测试整数范围验证
func TestIntFlag_RangeValidation(t *testing.T) {
	flag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			initialValue: 5,
			value:        new(int),
		},
	}

	// 设置范围为1-10
	flag.SetRange(1, 10)

	// 测试有效范围内的值
	validValues := []int{1, 5, 10}
	for _, val := range validValues {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置有效值%d失败: %v", val, err)
		}
	}

	// 测试超出范围的值
	invalidValues := []int{0, 11, -5, 100}
	for _, val := range invalidValues {
		if err := flag.Set(fmt.Sprint(val)); err == nil {
			t.Errorf("设置无效值%d应返回错误", val)
		}
	}
}

// TestIntFlag_Type 验证Type()方法返回正确的标志类型
func TestIntFlag_Type(t *testing.T) {
	flag := &IntFlag{}
	if flag.Type() != FlagTypeInt {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeInt, flag.Type())
	}
}

// TestInt64Flag_BasicFunctionality 测试Int64Flag的基本功能
func TestInt64Flag_BasicFunctionality(t *testing.T) {
	flag := &Int64Flag{
		BaseFlag: BaseFlag[int64]{
			initialValue: 0,
			value:        new(int64),
		},
	}

	// 测试默认值
	if flag.Get() != 0 {
		t.Error("默认值应为0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []int64{100, -50, 0, 9223372036854775807}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%d失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%d, 实际值%d", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0 {
		t.Error("重置后应返回默认值0")
	}
}

// TestInt64Flag_RangeValidation 测试64位整数范围验证
func TestInt64Flag_RangeValidation(t *testing.T) {
	flag := &Int64Flag{
		BaseFlag: BaseFlag[int64]{
			initialValue: 100,
			value:        new(int64),
		},
	}

	// 设置范围为-1000到1000
	flag.SetRange(-1000, 1000)

	// 测试有效范围内的值
	validValues := []int64{-1000, 0, 500, 1000}
	for _, val := range validValues {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置有效值%d失败: %v", val, err)
		}
	}

	// 测试超出范围的值
	invalidValues := []int64{-1001, 1001, -9223372036854775808, 9223372036854775807}
	for _, val := range invalidValues {
		if err := flag.Set(fmt.Sprint(val)); err == nil {
			t.Errorf("设置无效值%d应返回错误", val)
		}
	}
}

// TestInt64Flag_Type 验证Type()方法返回正确的标志类型
func TestInt64Flag_Type(t *testing.T) {
	flag := &Int64Flag{}
	if flag.Type() != FlagTypeInt64 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeInt64, flag.Type())
	}
}

// TestFloatFlag_BasicFunctionality 测试FloatFlag的基本功能
func TestFloatFlag_BasicFunctionality(t *testing.T) {
	flag := &Float64Flag{
		BaseFlag: BaseFlag[float64]{
			initialValue: 0.0,
			value:        new(float64),
		},
	}

	// 测试默认值
	if flag.Get() != 0.0 {
		t.Error("默认值应为0.0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []float64{3.14, -2.5, 0.0, 100.0}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%.2f失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%.2f, 实际值%.2f", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0.0 {
		t.Error("重置后应返回默认值0.0")
	}
}

// TestFloatFlag_Type 验证Type()方法返回正确的标志类型
func TestFloatFlag_Type(t *testing.T) {
	flag := &Float64Flag{}
	if flag.Type() != FlagTypeFloat64 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeFloat64, flag.Type())
	}
}

// TestUint16Flag_ValidParsing 测试有效的uint16值解析
func TestUint16Flag_ValidParsing(t *testing.T) {
	flag := &Uint16Flag{
		BaseFlag: BaseFlag[uint16]{
			initialValue: 0,
			value:        new(uint16),
		},
	}

	// 测试用例集合
	testCases := []struct {
		input    string
		expected uint16
	}{{
		input:    "0",
		expected: 0,
	}, {
		input:    "65535",
		expected: 65535,
	}, {
		input:    "32768",
		expected: 32768,
	}, {
		input:    "1000",
		expected: 1000,
	}}

	for _, tc := range testCases {
		if err := flag.Set(tc.input); err != nil {
			t.Errorf("解析 %s 失败: %v", tc.input, err)
			continue
		}
		if flag.Get() != tc.expected {
			t.Errorf("%s 期望 %d, 实际 %d", tc.input, tc.expected, flag.Get())
		}
	}
}

// TestUint16Flag_InvalidParsing 测试无效的uint16值解析
func TestUint16Flag_InvalidParsing(t *testing.T) {
	flag := &Uint16Flag{
		BaseFlag: BaseFlag[uint16]{
			initialValue: 0,
			value:        new(uint16),
		},
	}

	// 测试用例集合
	invalidInputs := []string{
		"-1",    // 负值
		"65536", // 超出最大值
		"abc",   // 非数字
		"12.34", // 浮点数
		" 123 ", // 带空格
	}

	for _, input := range invalidInputs {
		if err := flag.Set(input); err == nil {
			t.Errorf("解析无效值 '%s' 应返回错误", input)
		}
	}
}

// TestUint16Flag_TypeAndString 测试类型和字符串表示
func TestUint16Flag_TypeAndString(t *testing.T) {
	flag := &Uint16Flag{}

	// 测试Type()
	if flag.Type() != FlagTypeUint16 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeUint16, flag.Type())
	}

	// 测试String()
	if err := flag.Set("1234"); err != nil {
		t.Fatalf("设置值失败: %v", err)
	}
	if flag.String() != "1234" {
		t.Errorf("String()期望'1234', 实际'%s'", flag.String())
	}
}

// TestUint32Flag_BasicFunctionality 测试Uint32Flag的基本功能
func TestUint32Flag_BasicFunctionality(t *testing.T) {
	flag := &Uint32Flag{
		BaseFlag: BaseFlag[uint32]{
			initialValue: 0,
			value:        new(uint32),
		},
	}

	// 测试默认值
	if flag.Get() != 0 {
		t.Error("默认值应为0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []uint32{100, 4294967295, 0, 2147483648}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%d失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%d, 实际值%d", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0 {
		t.Error("重置后应返回默认值0")
	}
}

// TestUint32Flag_InvalidValue 测试设置无效值
func TestUint32Flag_InvalidValue(t *testing.T) {
	flag := &Uint32Flag{
		BaseFlag: BaseFlag[uint32]{
			value: new(uint32),
		},
	}

	invalidValues := []string{"4294967296", "-1", "abc"}
	for _, val := range invalidValues {
		if err := flag.Set(val); err == nil {
			t.Errorf("设置无效值%s应返回错误", val)
		}
	}
}

// TestUint32Flag_Type 验证Type()方法返回正确的标志类型
func TestUint32Flag_Type(t *testing.T) {
	flag := &Uint32Flag{}
	if flag.Type() != FlagTypeUint32 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeUint32, flag.Type())
	}
}

// TestUint64Flag_BasicFunctionality 测试Uint64Flag的基本功能
func TestUint64Flag_BasicFunctionality(t *testing.T) {
	flag := &Uint64Flag{
		BaseFlag: BaseFlag[uint64]{
			initialValue: 0,
			value:        new(uint64),
		},
	}

	// 测试默认值
	if flag.Get() != 0 {
		t.Error("默认值应为0")
	}
	if flag.IsSet() {
		t.Error("未设置值时IsSet应返回false")
	}

	// 测试设置有效值
	testCases := []uint64{100, 18446744073709551615, 0, 9223372036854775808}
	for _, val := range testCases {
		if err := flag.Set(fmt.Sprint(val)); err != nil {
			t.Errorf("设置值%d失败: %v", val, err)
		}
		if flag.Get() != val {
			t.Errorf("期望值%d, 实际值%d", val, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.IsSet() {
		t.Error("重置后IsSet应返回false")
	}
	if flag.Get() != 0 {
		t.Error("重置后应返回默认值0")
	}
}

// TestUint64Flag_InvalidValue 测试设置无效值
func TestUint64Flag_InvalidValue(t *testing.T) {
	flag := &Uint64Flag{
		BaseFlag: BaseFlag[uint64]{
			value: new(uint64),
		},
	}

	invalidValues := []string{"18446744073709551616", "-1", "not_a_number"}
	for _, val := range invalidValues {
		if err := flag.Set(val); err == nil {
			t.Errorf("设置无效值%s应返回错误", val)
		}
	}
}

// TestUint64Flag_Type 验证Type()方法返回正确的标志类型
func TestUint64Flag_Type(t *testing.T) {
	flag := &Uint64Flag{}
	if flag.Type() != FlagTypeUint64 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeUint64, flag.Type())
	}
}

// TestBaseFlag_GetPointer 验证GetPointer()方法的基本功能和指针访问有效性
func TestBaseFlag_GetPointer(t *testing.T) {
	// 1. 测试整数类型标志的指针行为
	intFlag := &IntFlag{
		BaseFlag: BaseFlag[int]{
			initialValue: 10,
			value:        nil,
		},
	}

	// 未设置值时指针应为nil
	if ptr := intFlag.GetPointer(); ptr != nil {
		t.Error("IntFlag未设置值时, GetPointer()应返回nil")
	}

	// 设置值后验证指针有效性
	if err := intFlag.Set("20"); err != nil {
		t.Fatalf("设置IntFlag值失败: %v", err)
	}

	ptr := intFlag.GetPointer()
	if ptr == nil {
		t.Fatal("IntFlag设置值后, GetPointer()不应返回nil")
	}

	if *ptr != 20 {
		t.Errorf("IntFlag指针值错误, 期望20, 实际%d", *ptr)
	}

	// 通过指针修改值并验证
	*ptr = 30
	if intFlag.Get() != 30 {
		t.Errorf("通过指针修改值失败, 期望30, 实际%d", intFlag.Get())
	}

	// 2. 测试字符串类型标志的指针行为
	strFlag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			initialValue: "default",
		},
	}

	if err := strFlag.Set("test"); err != nil {
		t.Fatalf("设置StringFlag值失败: %v", err)
	}

	*strFlag.GetPointer() = "modified"
	if strFlag.Get() != "modified" {
		t.Errorf("StringFlag指针修改失败, 期望'modified', 实际'%s'", strFlag.Get())
	}

	// 3. 测试默认值场景（值未显式设置时）
	defaultFlag := &BoolFlag{
		BaseFlag: BaseFlag[bool]{
			initialValue: true,
			value:        nil,
		},
	}

	// 未设置值时指针应为nil, Get()应返回默认值
	if ptr := defaultFlag.GetPointer(); ptr != nil {
		t.Error("BoolFlag未设置值时, GetPointer()应返回nil")
	}
	if defaultFlag.Get() != true {
		t.Error("BoolFlag未设置值时, Get()应返回默认值true")
	}
}
