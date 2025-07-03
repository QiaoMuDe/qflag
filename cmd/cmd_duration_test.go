package cmd

import (
	"flag"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// TestDurationVar 测试DurationVar方法的功能
func TestDurationVar(t *testing.T) {
	var durationFlagShort flags.DurationFlag

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("DurationVar with nil pointer should panic")
			}
		}()
		cmd.DurationVar(nil, "duration", "dur", time.Second*5, "test duration flag")
	})

	// 测试正常功能(长标志)
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var durationFlag flags.DurationFlag
		cmd.DurationVar(&durationFlag, "duration", "dur", time.Second*5, "test duration flag")

		// 测试默认值
		if durationFlag.Get() != time.Second*5 {
			t.Errorf("default value = %v, want %v", durationFlag.Get(), time.Second*5)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--duration", "10s"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if durationFlag.Get() != time.Second*10 {
			t.Errorf("after --duration, value = %v, want %v", durationFlag.Get(), time.Second*10)
		}
	})

	// 测试短标志解析
	t.Run("short flag", func(t *testing.T) {
		cmd := NewCommand("test-short", "ts", flag.ContinueOnError)
		cmd.DurationVar(&durationFlagShort, "duration-short", "d", time.Second*5, "test duration short flag")
		if err := cmd.Parse([]string{"-d", "2m"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if durationFlagShort.Get() != time.Minute*2 {
			t.Errorf("after -d, value = %v, want %v", durationFlagShort.Get(), time.Minute*2)
		}
	})
}
