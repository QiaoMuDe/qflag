package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestFloat64Var 测试Float64Var方法的功能
func TestFloat64Var(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Float64Var with nil pointer should panic")
			}
		}()
		cmd.Float64Var(nil, "float", "f", 3.14, "test float64 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		var floatFlag flags.Float64Flag
		cmd.Float64Var(&floatFlag, "float", "f", 3.14, "test float64 flag")

		// 测试默认值
		if floatFlag.Get() != 3.14 {
			t.Errorf("default value = %v, want %v", floatFlag.Get(), 3.14)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--float", "6.28"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if floatFlag.Get() != 6.28 {
			t.Errorf("after --float, value = %v, want %v", floatFlag.Get(), 6.28)
		}

		// 测试短标志解析
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
		var floatFlagShort flags.Float64Flag
		cmd.Float64Var(&floatFlagShort, "float-short", "f", 1.5, "test float64 short flag")
		if err := cmd.Parse([]string{"-f", "3.0"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if floatFlagShort.Get() != 3.0 {
			t.Errorf("after -f, value = %v, want %v", floatFlagShort.Get(), 3.0)
		}
	})
}
