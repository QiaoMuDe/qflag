package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestBoolVar 测试BoolVar方法的功能
func TestBoolVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("BoolVar with nil pointer should panic")
			}
		}()
		cmd.BoolVar(nil, "bool", "b", false, "test bool flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var boolFlag flags.BoolFlag
		cmd.BoolVar(&boolFlag, "bool", "b", false, "test bool flag")

		// 测试默认值
		if boolFlag.Get() != false {
			t.Errorf("default value = %v, want %v", boolFlag.Get(), false)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--bool"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if boolFlag.Get() != true {
			t.Errorf("after --bool, value = %v, want %v", boolFlag.Get(), true)
		}

		// 测试短标志解析
		cmd = NewCommand("test-short", "ts", flag.ContinueOnError)
		var boolFlagShort flags.BoolFlag
		cmd.BoolVar(&boolFlagShort, "bool-short", "b", false, "test bool short flag")
		if err := cmd.Parse([]string{"-b"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if boolFlagShort.Get() != true {
			t.Errorf("after -b, value = %v, want %v", boolFlagShort.Get(), true)
		}
	})
}
