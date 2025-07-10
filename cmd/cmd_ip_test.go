package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestIP4Var 测试IP4Var方法的功能
func TestIP4Var(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("IP4Var with nil pointer should panic")
			}
		}()
		cmd.IP4Var(nil, "ip4", "i", "127.0.0.1", "test ip4 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		var ip4Flag flags.IP4Flag
		cmd.IP4Var(&ip4Flag, "ip4", "i", "192.168.1.1", "test ip4 flag")

		// 测试默认值
		if ip4Flag.Get() != "192.168.1.1" {
			t.Errorf("default value = %v, want %v", ip4Flag.Get(), "192.168.1.1")
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--ip4", "10.0.0.1"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if ip4Flag.Get() != "10.0.0.1" {
			t.Errorf("after --ip4, value = %v, want %v", ip4Flag.Get(), "10.0.0.1")
		}

		// 测试短标志解析
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
		var ip4FlagShort flags.IP4Flag
		cmd.IP4Var(&ip4FlagShort, "ip4-short", "i", "172.16.0.1", "test ip4 short flag")
		if err := cmd.Parse([]string{"-i", "172.16.0.2"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if ip4FlagShort.Get() != "172.16.0.2" {
			t.Errorf("after -i, value = %v, want %v", ip4FlagShort.Get(), "172.16.0.2")
		}
	})

	// 测试无效IP地址
	t.Run("invalid ip4", func(t *testing.T) {
		cmd := NewCmd("test-invalid", "ti", flag.ContinueOnError)
		var ip4Flag flags.IP4Flag
		cmd.IP4Var(&ip4Flag, "ip4-invalid", "v", "127.0.0.1", "test invalid ip4")

		// 无效的IPv4地址
		err := cmd.Parse([]string{"--ip4-invalid", "256.0.0.1"})
		if err == nil {
			t.Error("expected error for invalid IPv4 address, got nil")
		}
	})
}

// TestIP6Var 测试IP6Var方法的功能
func TestIP6Var(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("IP6Var with nil pointer should panic")
			}
		}()
		cmd.IP6Var(nil, "ip6", "I", "::1", "test ip6 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		var ip6Flag flags.IP6Flag
		cmd.IP6Var(&ip6Flag, "ip6", "I", "2001:db8::1", "test ip6 flag")

		// 测试默认值
		if ip6Flag.Get() != "2001:db8::1" {
			t.Errorf("default value = %v, want %v", ip6Flag.Get(), "2001:db8::1")
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--ip6", "2001:db8::2"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if ip6Flag.Get() != "2001:db8::2" {
			t.Errorf("after --ip6, value = %v, want %v", ip6Flag.Get(), "2001:db8::2")
		}

		// 测试短标志解析
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
		var ip6FlagShort flags.IP6Flag
		cmd.IP6Var(&ip6FlagShort, "ip6-short", "I", "fe80::1", "test ip6 short flag")
		if err := cmd.Parse([]string{"-I", "fe80::2"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if ip6FlagShort.Get() != "fe80::2" {
			t.Errorf("after -I, value = %v, want %v", ip6FlagShort.Get(), "fe80::2")
		}
	})

	// 测试无效IP地址
	t.Run("invalid ip6", func(t *testing.T) {
		cmd := NewCmd("test-invalid", "ti", flag.ContinueOnError)
		var ip6Flag flags.IP6Flag
		cmd.IP6Var(&ip6Flag, "ip6-invalid", "V", "::1", "test invalid ip6")

		// 无效的IPv6地址
		err := cmd.Parse([]string{"--ip6-invalid", "2001:db8::gibberish"})
		if err == nil {
			t.Error("expected error for invalid IPv6 address, got nil")
		}
	})
}
