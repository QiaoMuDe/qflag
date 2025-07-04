package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestInt64Var 测试Int64Var方法的功能
func TestInt64Var(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Int64Var with nil pointer should panic")
			}
		}()
		cmd.Int64Var(nil, "int64", "i", 42, "test int64 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		var int64Flag flags.Int64Flag
		cmd.Int64Var(&int64Flag, "int64", "i", 42, "test int64 flag")

		// 测试默认值
		if int64Flag.Get() != 42 {
			t.Errorf("default value = %v, want %v", int64Flag.Get(), 42)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--int64", "100"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if int64Flag.Get() != 100 {
			t.Errorf("after --int64, value = %v, want %v", int64Flag.Get(), 100)
		}

		// 测试短标志解析
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
		var int64FlagShort flags.Int64Flag
		cmd.Int64Var(&int64FlagShort, "int64-short", "i", -5, "test int64 short flag")
		if err := cmd.Parse([]string{"-i", "-20"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if int64FlagShort.Get() != -20 {
			t.Errorf("after -i, value = %v, want %v", int64FlagShort.Get(), -20)
		}
	})

	// 测试边界值
	t.Run("boundary values", func(t *testing.T) {
		cmd := NewCmd("test-boundary", "tb", flag.ContinueOnError)
		var int64Flag flags.Int64Flag
		cmd.Int64Var(&int64Flag, "int64-boundary", "b", 0, "test int64 boundary values")

		// 测试最大值
		if err := cmd.Parse([]string{"--int64-boundary", "9223372036854775807"}); err != nil {
			t.Fatalf("Parse failed for max value: %v", err)
		}
		if int64Flag.Get() != 9223372036854775807 {
			t.Errorf("max value = %v, want %v", int64Flag.Get(), 9223372036854775807)
		}

		// 测试最小值
		cmd = NewCmd("test-min", "tm", flag.ContinueOnError)
		var minInt64Flag flags.Int64Flag
		cmd.Int64Var(&minInt64Flag, "int64-min", "m", 0, "test int64 min value")
		if err := cmd.Parse([]string{"--int64-min", "-9223372036854775808"}); err != nil {
			t.Fatalf("Parse failed for min value: %v", err)
		}
		if minInt64Flag.Get() != -9223372036854775808 {
			t.Errorf("min value = %v, want %v", minInt64Flag.Get(), -9223372036854775808)
		}
	})
}
