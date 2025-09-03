package flags

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// SizeUnit 大小单位定义
type SizeUnit struct {
	Name       string // 单位名称
	Multiplier int64  // 乘数
	IsBinary   bool   // 是否为二进制单位
}

// sizeUnits 单位映射表 - 按长度排序，优先匹配长单位避免歧义
var sizeUnits = map[string]int64{
	// 字节单位
	"bytes": 1,
	"byte":  1,
	"b":     1,

	// 二进制单位 (1024) - IEC标准
	"pib": 1024 * 1024 * 1024 * 1024 * 1024,
	"tib": 1024 * 1024 * 1024 * 1024,
	"gib": 1024 * 1024 * 1024,
	"mib": 1024 * 1024,
	"kib": 1024,

	// 十进制单位 (1000) - SI标准
	"pb": 1000 * 1000 * 1000 * 1000 * 1000,
	"tb": 1000 * 1000 * 1000 * 1000,
	"gb": 1000 * 1000 * 1000,
	"mb": 1000 * 1000,
	"kb": 1000,

	// 简写单位 (默认二进制)
	"p": 1024 * 1024 * 1024 * 1024 * 1024,
	"t": 1024 * 1024 * 1024 * 1024,
	"g": 1024 * 1024 * 1024,
	"m": 1024 * 1024,
	"k": 1024,
}

// unitOrder 单位匹配顺序 - 长单位优先，避免歧义
var unitOrder = []string{
	"bytes", "byte", "pib", "tib", "gib", "mib", "kib",
	"pb", "tb", "gb", "mb", "kb", "p", "t", "g", "m", "k", "b",
}

// SizeFlag 大小标志结构体
type SizeFlag struct {
	BaseFlag[int64]
	allowDecimal  bool // 是否允许小数
	allowNegative bool // 是否允许负数
	mu            sync.RWMutex
}

// Init 初始化大小标志
func (f *SizeFlag) Init(longName, shortName string, defValue int64, usage string) {
	f.value = &defValue
	f.BaseFlag.Init(longName, shortName, usage, f.value)
	f.allowDecimal = true   // 默认允许小数
	f.allowNegative = false // 默认不允许负数
}

// Set 设置标志值
func (f *SizeFlag) Set(value string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 处理空值
	if strings.TrimSpace(value) == "" {
		return nil
	}

	// 解析大小值
	size, err := f.parseSize(value)
	if err != nil {
		return fmt.Errorf("invalid size format: %v", err)
	}

	*f.value = size
	f.isSet = true
	return nil
}

// parseSize 解析大小字符串 - 使用后缀匹配方式
func (f *SizeFlag) parseSize(input string) (int64, error) {
	// 1. 清理输入并转换为小写
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return 0, fmt.Errorf("empty input")
	}

	// 2. 尝试匹配已知单位后缀
	numStr, multiplier, err := f.extractNumberAndUnit(input)
	if err != nil {
		return 0, err
	}

	// 3. 解析数字部分
	num, err := f.parseNumber(numStr)
	if err != nil {
		return 0, err
	}

	// 4. 检查负数
	if num < 0 && !f.allowNegative {
		return 0, fmt.Errorf("negative values not allowed")
	}

	// 5. 计算最终大小，检查溢出
	return f.calculateSize(num, multiplier)
}

// extractNumberAndUnit 提取数字和单位
func (f *SizeFlag) extractNumberAndUnit(input string) (string, int64, error) {
	// 按顺序检查单位后缀，长单位优先
	for _, unit := range unitOrder {
		if strings.HasSuffix(input, unit) {
			// 提取数字部分
			numStr := strings.TrimSpace(input[:len(input)-len(unit)])
			if numStr == "" {
				return "", 0, fmt.Errorf("missing number before unit '%s'", unit)
			}

			// 检查数字部分是否纯净（不包含其他字母）
			if f.containsLetters(numStr) {
				return "", 0, fmt.Errorf("invalid format: number contains letters")
			}

			return numStr, sizeUnits[unit], nil
		}
	}

	// 没有匹配到单位，检查是否为纯数字
	if f.containsLetters(input) {
		return "", 0, fmt.Errorf("invalid format: contains unrecognized units or mixed content")
	}

	// 纯数字，默认为字节
	return input, 1, nil
}

// containsLetters 检查字符串是否包含字母
func (f *SizeFlag) containsLetters(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

// parseNumber 解析数字部分
func (f *SizeFlag) parseNumber(numStr string) (float64, error) {
	// 检查是否包含小数点
	if strings.Contains(numStr, ".") {
		if !f.allowDecimal {
			return 0, fmt.Errorf("decimal values not allowed")
		}
		return strconv.ParseFloat(numStr, 64)
	}

	// 解析整数
	intVal, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", numStr)
	}
	return float64(intVal), nil
}

// calculateSize 计算最终大小，检查溢出
func (f *SizeFlag) calculateSize(num float64, multiplier int64) (int64, error) {
	// 检查乘法溢出
	const maxFloat64ForInt64 = float64(1<<63 - 1)
	const minFloat64ForInt64 = float64(-1 << 63)

	result := num * float64(multiplier)

	// 检查结果是否在int64范围内
	if result > maxFloat64ForInt64 {
		return 0, fmt.Errorf("size too large: overflow")
	}
	if result < minFloat64ForInt64 {
		return 0, fmt.Errorf("size too small: underflow")
	}

	return int64(result), nil
}

// String 返回标志的字符串表示
func (f *SizeFlag) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	size := f.Get()

	// 处理零值
	if size == 0 {
		return "0B"
	}

	return f.formatSize(size)
}

// formatSize 格式化大小为人类可读的字符串
func (f *SizeFlag) formatSize(bytes int64) string {
	if bytes < 0 {
		return fmt.Sprintf("-%s", f.formatSize(-bytes))
	}

	units := []struct {
		name string
		size int64
	}{
		{"PiB", 1024 * 1024 * 1024 * 1024 * 1024},
		{"TiB", 1024 * 1024 * 1024 * 1024},
		{"GiB", 1024 * 1024 * 1024},
		{"MiB", 1024 * 1024},
		{"KiB", 1024},
		{"B", 1},
	}

	for _, unit := range units {
		if bytes >= unit.size {
			value := float64(bytes) / float64(unit.size)
			if value == float64(int64(value)) {
				return fmt.Sprintf("%.0f%s", value, unit.name)
			}
			return fmt.Sprintf("%.1f%s", value, unit.name)
		}
	}

	return fmt.Sprintf("%dB", bytes)
}

// Type 返回标志类型
func (f *SizeFlag) Type() FlagType {
	return FlagTypeSize
}

// SetAllowDecimal 设置是否允许小数
func (f *SizeFlag) SetAllowDecimal(allow bool) *SizeFlag {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowDecimal = allow
	return f
}

// SetAllowNegative 设置是否允许负数
func (f *SizeFlag) SetAllowNegative(allow bool) *SizeFlag {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowNegative = allow
	return f
}

// GetBytes 获取字节数
func (f *SizeFlag) GetBytes() int64 {
	return f.Get()
}

// GetKiB 获取KiB数
func (f *SizeFlag) GetKiB() float64 {
	return float64(f.Get()) / 1024
}

// GetMiB 获取MiB数
func (f *SizeFlag) GetMiB() float64 {
	return float64(f.Get()) / (1024 * 1024)
}

// GetGiB 获取GiB数
func (f *SizeFlag) GetGiB() float64 {
	return float64(f.Get()) / (1024 * 1024 * 1024)
}

// GetTiB 获取TiB数
func (f *SizeFlag) GetTiB() float64 {
	return float64(f.Get()) / (1024 * 1024 * 1024 * 1024)
}

// GetPiB 获取PiB数
func (f *SizeFlag) GetPiB() float64 {
	return float64(f.Get()) / (1024 * 1024 * 1024 * 1024 * 1024)
}

// IsZero 检查是否为零
func (f *SizeFlag) IsZero() bool {
	return f.Get() == 0
}

// IsPositive 检查是否为正数
func (f *SizeFlag) IsPositive() bool {
	return f.Get() > 0
}

// IsNegative 检查是否为负数
func (f *SizeFlag) IsNegative() bool {
	return f.Get() < 0
}

// GetAllowDecimal 获取是否允许小数设置
func (f *SizeFlag) GetAllowDecimal() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.allowDecimal
}

// GetAllowNegative 获取是否允许负数设置
func (f *SizeFlag) GetAllowNegative() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.allowNegative
}

// Clone 克隆标志
func (f *SizeFlag) Clone() Flag {
	f.mu.RLock()
	defer f.mu.RUnlock()

	clone := &SizeFlag{
		allowDecimal:  f.allowDecimal,
		allowNegative: f.allowNegative,
	}

	// 复制BaseFlag的字段
	clone.BaseFlag.longName = f.BaseFlag.longName
	clone.BaseFlag.shortName = f.BaseFlag.shortName
	clone.BaseFlag.initialValue = f.BaseFlag.initialValue
	clone.BaseFlag.usage = f.BaseFlag.usage
	clone.BaseFlag.validator = f.BaseFlag.validator
	clone.BaseFlag.initialized = f.BaseFlag.initialized
	clone.BaseFlag.isSet = f.BaseFlag.isSet
	clone.BaseFlag.envVar = f.BaseFlag.envVar

	// 创建新的值指针
	if f.BaseFlag.value != nil {
		newValue := *f.BaseFlag.value
		clone.BaseFlag.value = &newValue
	}

	return clone
}
