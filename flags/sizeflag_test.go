package flags

import (
	"fmt"
	"math"
	"testing"
)

// TestSizeFlagInit 测试初始化功能
func TestSizeFlagInit(t *testing.T) {
	tests := []struct {
		name        string
		longName    string
		shortName   string
		defValue    int64
		usage       string
		expectError bool
	}{
		{
			name:        "正常初始化",
			longName:    "size",
			shortName:   "s",
			defValue:    1024,
			usage:       "File size",
			expectError: false,
		},
		{
			name:        "零值初始化",
			longName:    "size",
			shortName:   "s",
			defValue:    0,
			usage:       "File size",
			expectError: false,
		},
		{
			name:        "负值初始化",
			longName:    "size",
			shortName:   "s",
			defValue:    -1024,
			usage:       "File size",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			err := flag.Init(tt.longName, tt.shortName, tt.defValue, tt.usage)

			if tt.expectError && err == nil {
				t.Errorf("期望错误但没有发生")
			}
			if !tt.expectError && err != nil {
				t.Errorf("不期望错误但发生了: %v", err)
			}
			if err == nil && flag.Get() != tt.defValue {
				t.Errorf("默认值设置错误，期望 %d，得到 %d", tt.defValue, flag.Get())
			}
		})
	}
}

// TestSizeFlagInitOnce 测试 sync.Once 功能
func TestSizeFlagInitOnce(t *testing.T) {
	flag := &SizeFlag{}

	// 第一次初始化
	err1 := flag.Init("size", "s", 1024, "File size")
	if err1 != nil {
		t.Fatalf("第一次初始化失败: %v", err1)
	}

	// 第二次初始化（应该被忽略）
	err2 := flag.Init("size", "s", 2048, "Different size")
	if err2 != nil {
		t.Errorf("第二次初始化应该被忽略但返回错误: %v", err2)
	}

	// 验证值没有被第二次初始化改变
	if flag.Get() != 1024 {
		t.Errorf("值被第二次初始化改变了，期望 1024，得到 %d", flag.Get())
	}
}

// TestSizeFlagBasicUnits 测试基本单位解析
func TestSizeFlagBasicUnits(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		// 字节单位
		{"字节 - b", "1024b", 1024},
		{"字节 - B", "1024B", 1024},
		{"字节 - byte", "1024byte", 1024},
		{"字节 - bytes", "1024bytes", 1024},

		// 二进制单位 (1024)
		{"KiB", "1KiB", 1024},
		{"MiB", "1MiB", 1024 * 1024},
		{"GiB", "1GiB", 1024 * 1024 * 1024},
		{"TiB", "1TiB", 1024 * 1024 * 1024 * 1024},
		{"PiB", "1PiB", 1024 * 1024 * 1024 * 1024 * 1024},

		// 十进制单位 (1000)
		{"KB", "1KB", 1000},
		{"MB", "1MB", 1000 * 1000},
		{"GB", "1GB", 1000 * 1000 * 1000},
		{"TB", "1TB", 1000 * 1000 * 1000 * 1000},
		{"PB", "1PB", 1000 * 1000 * 1000 * 1000 * 1000},

		// 简写单位 (默认二进制)
		{"K", "1K", 1024},
		{"M", "1M", 1024 * 1024},
		{"G", "1G", 1024 * 1024 * 1024},
		{"T", "1T", 1024 * 1024 * 1024 * 1024},
		{"P", "1P", 1024 * 1024 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")

			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("解析 %s 失败: %v", tt.input, err)
			}

			if flag.Get() != tt.expected {
				t.Errorf("解析 %s 错误，期望 %d，得到 %d", tt.input, tt.expected, flag.Get())
			}
		})
	}
}

// TestSizeFlagDecimalValues 测试小数值解析
func TestSizeFlagDecimalValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"小数 KB", "1.5KB", 1500},
		{"小数 MB", "2.5MB", 2500000},
		{"小数 GB", "1.5GB", 1500000000},
		{"小数 KiB", "1.5KiB", 1536},    // 1.5 * 1024
		{"小数 MiB", "2.5MiB", 2621440}, // 2.5 * 1024 * 1024
		{"零点几", "0.5GB", 500000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")
			flag.SetAllowDecimal(true)

			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("解析 %s 失败: %v", tt.input, err)
			}

			if flag.Get() != tt.expected {
				t.Errorf("解析 %s 错误，期望 %d，得到 %d", tt.input, tt.expected, flag.Get())
			}
		})
	}
}

// TestSizeFlagNegativeValues 测试负数值解析
func TestSizeFlagNegativeValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"负数 KB", "-1KB", -1000},
		{"负数 MB", "-2MB", -2000000},
		{"负数 GB", "-1GB", -1000000000},
		{"负数小数", "-1.5GB", -1500000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")
			flag.SetAllowNegative(true)

			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("解析 %s 失败: %v", tt.input, err)
			}

			if flag.Get() != tt.expected {
				t.Errorf("解析 %s 错误，期望 %d，得到 %d", tt.input, tt.expected, flag.Get())
			}
		})
	}
}

// TestSizeFlagZeroValue 测试零值处理
func TestSizeFlagZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"零值特例", "0", 0},
		{"零字节", "0B", 0},
		{"零KB", "0KB", 0},
		{"零MB", "0MB", 0},
		{"空字符串", "", 0},
		{"空白字符", "   ", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 100, "test") // 设置非零默认值

			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("解析 %s 失败: %v", tt.input, err)
			}

			if flag.Get() != tt.expected {
				t.Errorf("解析 %s 错误，期望 %d，得到 %d", tt.input, tt.expected, flag.Get())
			}
		})
	}
}

// TestSizeFlagCaseInsensitive 测试大小写不敏感
func TestSizeFlagCaseInsensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"大写 KB", "1KB", 1000},
		{"小写 kb", "1kb", 1000},
		{"混合 Kb", "1Kb", 1000},
		{"混合 kB", "1kB", 1000},
		{"大写 GIB", "1GIB", 1024 * 1024 * 1024},
		{"小写 gib", "1gib", 1024 * 1024 * 1024},
		{"大写 BYTES", "1024BYTES", 1024},
		{"小写 bytes", "1024bytes", 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")

			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("解析 %s 失败: %v", tt.input, err)
			}

			if flag.Get() != tt.expected {
				t.Errorf("解析 %s 错误，期望 %d，得到 %d", tt.input, tt.expected, flag.Get())
			}
		})
	}
}

// TestSizeFlagSpaceHandling 测试空格处理
func TestSizeFlagSpaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"数字和单位之间有空格", "1 GB", 1000000000},
		{"前后有空格", "  1GB  ", 1000000000},
		{"多个空格", "1   GB", 1000000000},
		{"制表符", "1\tGB", 1000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")

			err := flag.Set(tt.input)
			if err != nil {
				t.Fatalf("解析 %s 失败: %v", tt.input, err)
			}

			if flag.Get() != tt.expected {
				t.Errorf("解析 %s 错误，期望 %d，得到 %d", tt.input, tt.expected, flag.Get())
			}
		})
	}
}

// TestSizeFlagInvalidInputs 测试无效输入
func TestSizeFlagInvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"纯数字（非零）", "1024"},
		{"纯数字（大数）", "999999"},
		{"无效单位", "1XB"},
		{"无效格式", "abc"},
		{"只有单位", "GB"},
		{"只有数字和空格", "1024 "},
		{"多个数字", "1 2 GB"},
		{"无效字符", "1@GB"},
		{"空单位", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")

			err := flag.Set(tt.input)
			if err == nil {
				t.Errorf("期望 %s 解析失败但成功了", tt.input)
			}
		})
	}
}

// TestSizeFlagDecimalRestriction 测试小数限制
func TestSizeFlagDecimalRestriction(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 0, "test")
	flag.SetAllowDecimal(false) // 不允许小数

	tests := []string{"1.5GB", "2.0MB", "0.5KB"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("期望 %s 解析失败（不允许小数）但成功了", input)
			}
		})
	}
}

// TestSizeFlagNegativeRestriction 测试负数限制
func TestSizeFlagNegativeRestriction(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 0, "test")
	flag.SetAllowNegative(false) // 不允许负数

	tests := []string{"-1GB", "-500MB", "-1KB"}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			err := flag.Set(input)
			if err == nil {
				t.Errorf("期望 %s 解析失败（不允许负数）但成功了", input)
			}
		})
	}
}

// TestSizeFlagOverflow 测试溢出处理
func TestSizeFlagOverflow(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"超大正数", "9999999999999999999999PB"},
		{"超大负数", "-9999999999999999999999PB"},
		{"接近溢出", "9223372036854775807B"}, // 这个应该成功
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", 0, "test")
			flag.SetAllowNegative(true)

			err := flag.Set(tt.input)
			// 前两个应该失败，第三个应该成功
			if tt.name == "接近溢出" {
				if err != nil {
					t.Errorf("期望 %s 解析成功但失败了: %v", tt.input, err)
				}
			} else {
				if err == nil {
					t.Errorf("期望 %s 解析失败（溢出）但成功了", tt.input)
				}
			}
		})
	}
}

// TestSizeFlagString 测试字符串格式化
func TestSizeFlagString(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{"零值", 0, "0B"},
		{"字节", 512, "512B"},
		{"1KB", 1000, "1KB"},
		{"1.5KB", 1500, "1.5KB"},
		{"1MB", 1000000, "1MB"},
		{"1.2MB", 1200000, "1.2MB"},
		{"1GB", 1000000000, "1GB"},
		{"1.5GB", 1500000000, "1.5GB"},
		{"1TB", 1000000000000, "1TB"},
		{"负数", -1000, "-1KB"},
		{"负数小数", -1500, "-1.5KB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := &SizeFlag{}
			flag.Init("size", "s", tt.value, "test")

			result := flag.String()
			if result != tt.expected {
				t.Errorf("格式化 %d 错误，期望 %s，得到 %s", tt.value, tt.expected, result)
			}
		})
	}
}

// TestSizeFlagGetters 测试各种获取方法
func TestSizeFlagGetters(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 2048, "test") // 2KiB

	tests := []struct {
		name     string
		method   func() interface{}
		expected interface{}
	}{
		{"GetBytes", func() interface{} { return flag.GetBytes() }, int64(2048)},
		{"GetKiB", func() interface{} { return flag.GetKiB() }, 2.0},
		{"GetMiB", func() interface{} { return flag.GetMiB() }, 2.0 / 1024},
		{"IsZero", func() interface{} { return flag.IsZero() }, false},
		{"IsPositive", func() interface{} { return flag.IsPositive() }, true},
		{"IsNegative", func() interface{} { return flag.IsNegative() }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()

			// 对于浮点数，使用近似比较
			if expectedFloat, ok := tt.expected.(float64); ok {
				if resultFloat, ok := result.(float64); ok {
					if math.Abs(resultFloat-expectedFloat) > 1e-9 {
						t.Errorf("%s 错误，期望 %v，得到 %v", tt.name, tt.expected, result)
					}
				} else {
					t.Errorf("%s 返回类型错误", tt.name)
				}
			} else {
				if result != tt.expected {
					t.Errorf("%s 错误，期望 %v，得到 %v", tt.name, tt.expected, result)
				}
			}
		})
	}
}

// TestSizeFlagZeroGetters 测试零值的获取方法
func TestSizeFlagZeroGetters(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 0, "test")

	if !flag.IsZero() {
		t.Error("IsZero() 应该返回 true")
	}
	if flag.IsPositive() {
		t.Error("IsPositive() 应该返回 false")
	}
	if flag.IsNegative() {
		t.Error("IsNegative() 应该返回 false")
	}
}

// TestSizeFlagNegativeGetters 测试负值的获取方法
func TestSizeFlagNegativeGetters(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", -1024, "test")

	if flag.IsZero() {
		t.Error("IsZero() 应该返回 false")
	}
	if flag.IsPositive() {
		t.Error("IsPositive() 应该返回 false")
	}
	if !flag.IsNegative() {
		t.Error("IsNegative() 应该返回 true")
	}
}

// TestSizeFlagType 测试类型返回
func TestSizeFlagType(t *testing.T) {
	flag := &SizeFlag{}
	if flag.Type() != FlagTypeSize {
		t.Errorf("Type() 错误，期望 %v，得到 %v", FlagTypeSize, flag.Type())
	}
}

// TestSizeFlagChainedSetters 测试链式设置
func TestSizeFlagChainedSetters(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 0, "test")

	// 测试链式调用
	result := flag.SetAllowDecimal(true).SetAllowNegative(true)
	if result != flag {
		t.Error("链式调用应该返回自身")
	}

	// 验证设置生效
	if !flag.GetAllowDecimal() {
		t.Error("AllowDecimal 设置失败")
	}
	if !flag.GetAllowNegative() {
		t.Error("AllowNegative 设置失败")
	}
}

// TestSizeFlagConcurrency 测试并发安全
func TestSizeFlagConcurrency(t *testing.T) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 0, "test")

	// 启动多个 goroutine 同时设置值
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(val int) {
			err := flag.Set(fmt.Sprintf("%dMB", val))
			if err != nil {
				t.Errorf("并发设置失败: %v", err)
			}
			done <- true
		}(i + 1)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证最终值是有效的
	value := flag.Get()
	if value < 1000000 || value > 10000000 { // 1MB 到 10MB 之间
		t.Errorf("并发测试后的值异常: %d", value)
	}
}

// BenchmarkSizeFlagSet 性能测试
func BenchmarkSizeFlagSet(b *testing.B) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 0, "test")

	inputs := []string{"1GB", "512MB", "1.5TB", "2048KB", "1PB"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		err := flag.Set(input)
		if err != nil {
			b.Fatalf("性能测试失败: %v", err)
		}
	}
}

// BenchmarkSizeFlagString 字符串格式化性能测试
func BenchmarkSizeFlagString(b *testing.B) {
	flag := &SizeFlag{}
	flag.Init("size", "s", 1500000000, "test") // 1.5GB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = flag.String()
	}
}
