package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestIntVar 测试IntVar方法的功能
func TestIntVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("IntVar with nil pointer should panic")
			}
		}()
		cmd.IntVar(nil, "int", "i", 42, "test int flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		var intFlag flags.IntFlag
		cmd.IntVar(&intFlag, "int", "i", 42, "test int flag")

		// 测试默认值
		if intFlag.Get() != 42 {
			t.Errorf("default value = %v, want %v", intFlag.Get(), 42)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--int", "100"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if intFlag.Get() != 100 {
			t.Errorf("after --int, value = %v, want %v", intFlag.Get(), 100)
		}

		// 测试短标志解析
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
		var intFlagShort flags.IntFlag
		cmd.IntVar(&intFlagShort, "int-short", "i", -5, "test int short flag")
		if err := cmd.Parse([]string{"-i", "-20"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if intFlagShort.Get() != -20 {
			t.Errorf("after -i, value = %v, want %v", intFlagShort.Get(), -20)
		}
	})

	// 测试边界值
	t.Run("boundary values", func(t *testing.T) {
		cmd := NewCmd("test-boundary", "tb", flag.ContinueOnError)
		var intFlag flags.IntFlag
		cmd.IntVar(&intFlag, "int-boundary", "b", 0, "test int boundary values")

		// 测试最大值
		if err := cmd.Parse([]string{"--int-boundary", "2147483647"}); err != nil {
			t.Fatalf("Parse failed for max value: %v", err)
		}
		if intFlag.Get() != 2147483647 {
			t.Errorf("max value = %v, want %v", intFlag.Get(), 2147483647)
		}

		// 测试最小值
		cmd = NewCmd("test-min", "tm", flag.ContinueOnError)
		var minIntFlag flags.IntFlag
		cmd.IntVar(&minIntFlag, "int-min", "m", 0, "test int min value")
		if err := cmd.Parse([]string{"--int-min", "-2147483648"}); err != nil {
			t.Fatalf("Parse failed for min value: %v", err)
		}
		if minIntFlag.Get() != -2147483648 {
			t.Errorf("min value = %v, want %v", minIntFlag.Get(), -2147483648)
		}
	})
}
