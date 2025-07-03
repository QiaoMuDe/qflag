package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestIP6Var 测试IP6Var方法的功能
func TestIP6Var(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("IP6Var with nil pointer should panic")
			}
		}()
		cmd.IP6Var(nil, "ip6", "I", "::1", "test ip6 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
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
		cmd = NewCommand("test-short", "ts", flag.ContinueOnError)
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
		cmd := NewCommand("test-invalid", "ti", flag.ContinueOnError)
		var ip6Flag flags.IP6Flag
		cmd.IP6Var(&ip6Flag, "ip6-invalid", "V", "::1", "test invalid ip6")

		// 无效的IPv6地址
		err := cmd.Parse([]string{"--ip6-invalid", "2001:db8::gibberish"})
		if err == nil {
			t.Error("expected error for invalid IPv6 address, got nil")
		}
	})
}
