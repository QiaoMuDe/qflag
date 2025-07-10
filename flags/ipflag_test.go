package flags

import "testing"

// TestIP4Flag_BasicFunctionality 测试IP4Flag的基本功能
func TestIP4Flag_BasicFunctionality(t *testing.T) {
	flag := &IP4Flag{
		BaseFlag: BaseFlag[string]{
			initialValue: "127.0.0.1",
			value:        new(string),
		},
	}

	// 测试默认值
	if flag.GetDefault() != "127.0.0.1" {
		t.Errorf("默认值应为127.0.0.1, 实际为%s", flag.Get())
	}

	// 测试设置有效IPv4地址
	validIPs := []string{"8.8.8.8", "192.168.1.1", "0.0.0.0", "255.255.255.255"}
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
	if flag.Get() != "127.0.0.1" {
		t.Errorf("重置后应返回默认值127.0.0.1, 实际为%s", flag.Get())
	}
}

// TestIP4Flag_InvalidValue 测试设置无效IPv4地址
func TestIP4Flag_InvalidValue(t *testing.T) {
	flag := &IP4Flag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	invalidIPs := []string{
		"", "256.0.0.1", "192.168.1", "192.168.1.1.1",
		"fe80::1", "not.an.ip", "192.168.1.a",
	}
	for _, ip := range invalidIPs {
		if err := flag.Set(ip); err == nil {
			t.Errorf("设置无效IP %s 应返回错误", ip)
		}
	}
}

// TestIP4Flag_Type 验证Type()方法返回正确的标志类型
func TestIP4Flag_Type(t *testing.T) {
	flag := &IP4Flag{}
	if flag.Type() != FlagTypeIP4 {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeIP4, flag.Type())
	}
}

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
