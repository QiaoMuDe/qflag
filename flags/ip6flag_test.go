package flags

import (
	"testing"
)

// TestIP6Flag_BasicFunctionality 测试IP6Flag的基本功能
func TestIP6Flag_BasicFunctionality(t *testing.T) {
	flag := &IP6Flag{
		BaseFlag: BaseFlag[string]{
			initialValue: "::1",
			value:        new(string),
		},
	}

	// 测试默认值
	if flag.GetDefault() != "::1" {
		t.Errorf("默认值应为::1, 实际为%s", flag.GetDefault())
	}

	// 测试设置有效IPv6地址
	validIPs := []string{"2001:db8::1", "fe80::1", "::", "2001:db8:85a3::8a2e:370:7334"}
	for _, ip := range validIPs {
		if err := flag.Set(ip); err != nil {
			t.Errorf("设置有效IP %s 失败: %v", ip, err)
		}
		if flag.Get() != ip {
			t.Errorf("设置IP后值不匹配, 期望%s, 实际%s", ip, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.Get() != "::1" {
		t.Errorf("重置后应返回默认值::1, 实际为%s", flag.Get())
	}
}

// TestIP6Flag_InvalidValue 测试设置无效IPv6地址
func TestIP6Flag_InvalidValue(t *testing.T) {
	flag := &IP6Flag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	invalidIPs := []string{
		"", "192.168.1.1", "2001:db8::g", "2001:db8::1::2",
		"not.an.ipv6", "2001:db8:85a3:0000:0000:8a2e:0370:7334:7334",
	}
	for _, ip := range invalidIPs {
		if err := flag.Set(ip); err == nil {
			t.Errorf("设置无效IP %s 应返回错误", ip)
		}
	}
}

// TestIP6Flag_Type 验证Type()方法返回正确的标志类型
func TestIP6Flag_Type(t *testing.T) {
	flag := &IP6Flag{}
	if flag.Type() != FlagTypeIP6 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeIP6, flag.Type())
	}
}
