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

// TestURLFlag_BasicFunctionality 测试URLFlag的基本功能
func TestURLFlag_BasicFunctionality(t *testing.T) {
	flag := &URLFlag{
		BaseFlag: BaseFlag[string]{
			initialValue: "https://example.com",
			value:        new(string),
		},
	}

	// 测试默认值
	if flag.GetDefault() != "https://example.com" {
		t.Errorf("默认值应为https://example.com, 实际为%s", flag.GetDefault())
	}

	// 测试设置有效URL
	validURLs := []string{
		"https://www.google.com",
		"http://localhost:8080",
		"https://api.github.com/users",
		"ftp://files.example.com/path",
		"https://example.com:443/path?query=value#fragment",
	}

	for _, url := range validURLs {
		if err := flag.Set(url); err != nil {
			t.Errorf("设置有效URL %s 失败: %v", url, err)
		}
		if flag.Get() != url {
			t.Errorf("设置URL后值不匹配, 期望%s, 实际%s", url, flag.Get())
		}
	}

	// 测试重置功能
	flag.Reset()
	if flag.Get() != "https://example.com" {
		t.Errorf("重置后应返回默认值https://example.com, 实际为%s", flag.Get())
	}
}

// TestURLFlag_InvalidValue 测试设置无效URL
func TestURLFlag_InvalidValue(t *testing.T) {
	flag := &URLFlag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	invalidURLs := []string{
		"",                          // 空字符串
		"not-a-url",                 // 无协议的普通字符串
		"www.example.com",           // 缺少协议
		"://example.com",            // 空协议
		"http://",                   // 只有协议无主机
		"http:// invalid space.com", // 包含空格的无效URL
		"http://[invalid-ipv6",      // 无效的IPv6格式
		"javascript:alert('xss')",   // 可能的XSS攻击
	}

	for _, url := range invalidURLs {
		if err := flag.Set(url); err == nil {
			t.Errorf("设置无效URL %s 应返回错误", url)
		}
	}
}

// TestURLFlag_Type 验证Type()方法返回正确的标志类型
func TestURLFlag_Type(t *testing.T) {
	flag := &URLFlag{}
	if flag.Type() != FlagTypeURL {
		t.Errorf("Type()应返回%d, 实际返回%d", FlagTypeURL, flag.Type())
	}
}

// TestURLFlag_String 测试String()方法
func TestURLFlag_String(t *testing.T) {
	flag := &URLFlag{
		BaseFlag: BaseFlag[string]{
			initialValue: "https://test.com",
			value:        new(string),
		},
	}

	// 初始化后测试String()方法
	*flag.value = flag.initialValue
	if flag.String() != "https://test.com" {
		t.Errorf("String()应返回https://test.com, 实际返回%s", flag.String())
	}

	// 设置新值后测试String()方法
	_ = flag.Set("http://new-url.com")
	if flag.String() != "http://new-url.com" {
		t.Errorf("设置新值后String()应返回http://new-url.com, 实际返回%s", flag.String())
	}
}

// TestURLFlag_ConcurrentAccess 测试并发访问安全性
func TestURLFlag_ConcurrentAccess(t *testing.T) {
	flag := &URLFlag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	urls := []string{
		"https://example1.com",
		"https://example2.com",
		"https://example3.com",
		"https://example4.com",
		"https://example5.com",
	}

	// 启动多个goroutine并发设置URL
	done := make(chan bool, len(urls))
	for _, url := range urls {
		go func(u string) {
			defer func() { done <- true }()
			if err := flag.Set(u); err != nil {
				t.Errorf("并发设置URL %s 失败: %v", u, err)
			}
		}(url)
	}

	// 等待所有goroutine完成
	for i := 0; i < len(urls); i++ {
		<-done
	}

	// 验证最终值是有效的URL之一
	finalValue := flag.Get()
	found := false
	for _, url := range urls {
		if finalValue == url {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("并发访问后的最终值 %s 不在预期的URL列表中", finalValue)
	}
}

// TestURLFlag_SpecialCases 测试特殊情况
func TestURLFlag_SpecialCases(t *testing.T) {
	flag := &URLFlag{
		BaseFlag: BaseFlag[string]{
			value: new(string),
		},
	}

	// 测试包含特殊字符的URL
	specialURLs := []string{
		"https://example.com/path?param=value&other=123",
		"https://user:pass@example.com:8080/path",
		"https://example.com/path#fragment",
		"https://example.com/path%20with%20spaces",
	}

	for _, url := range specialURLs {
		if err := flag.Set(url); err != nil {
			t.Errorf("设置特殊URL %s 失败: %v", url, err)
		}
	}

	// 测试相对路径（应该失败）
	relativePaths := []string{
		"/relative/path",
		"./relative/path",
		"../relative/path",
		"relative/path",
	}

	for _, path := range relativePaths {
		if err := flag.Set(path); err == nil {
			t.Errorf("设置相对路径 %s 应该失败", path)
		}
	}
}
