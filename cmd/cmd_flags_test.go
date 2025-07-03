package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

// TestStringVar 测试StringVar方法的功能
func TestStringVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("StringVar with nil pointer should panic")
			}
		}()
		cmd.StringVar(nil, "str", "s", "default", "test string flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var strFlag flags.StringFlag
		cmd.StringVar(&strFlag, "str", "st", "default", "test string flag")

		// 测试默认值
		if strFlag.Get() != "default" {
			t.Errorf("default value = %q, want %q", strFlag.Get(), "default")
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--str", "value"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if strFlag.Get() != "value" {
			t.Errorf("after --str, value = %q, want %q", strFlag.Get(), "value")
		}

		// 测试短标志解析
		cmd = NewCommand("test-short", "ts", flag.ContinueOnError)
		var strFlagShort flags.StringFlag
		cmd.StringVar(&strFlagShort, "str-short", "t", "default", "test string short flag")
		if err := cmd.Parse([]string{"-t", "short"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if strFlagShort.Get() != "short" {
			t.Errorf("after -s, value = %q, want %q", strFlagShort.Get(), "short")
		}
	})
}

// TestIntVar 测试IntVar方法的功能
func TestIntVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("IntVar with nil pointer should panic")
			}
		}()
		cmd.IntVar(nil, "int", "i", 123, "test int flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var intFlag flags.IntFlag
		cmd.IntVar(&intFlag, "int", "iv", 123, "test int flag")

		// 测试默认值
		if intFlag.Get() != 123 {
			t.Errorf("default value = %d, want %d", intFlag.Get(), 123)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--int", "456"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if intFlag.Get() != 456 {
			t.Errorf("after --int, value = %d, want %d", intFlag.Get(), 456)
		}

		// 测试短标志解析
		cmd = NewCommand("test", "t", flag.ContinueOnError)
		var intFlagShort flags.IntFlag
		cmd.IntVar(&intFlagShort, "int", "iv", 123, "test int flag")
		if err := cmd.Parse([]string{"-iv", "789"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if intFlagShort.Get() != 789 {
			t.Errorf("after -iv, value = %d, want %d", intFlagShort.Get(), 789)
		}
	})
}

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
		cmd.BoolVar(&boolFlag, "bool", "bl", false, "test bool flag")

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
		cmd = NewCommand("test", "t", flag.ContinueOnError)
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

// TestFloatVar 测试FloatVar方法的功能
func TestFloatVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("FloatVar with nil pointer should panic")
			}
		}()
		cmd.Float64Var(nil, "float", "f", 3.14, "test float flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var floatFlag flags.Float64Flag
		cmd.Float64Var(&floatFlag, "float", "fl", 3.14, "test float flag")

		// 测试默认值
		if floatFlag.Get() != 3.14 {
			t.Errorf("default value = %v, want %v", floatFlag.Get(), 3.14)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--float", "2.718"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if floatFlag.Get() != 2.718 {
			t.Errorf("after --float, value = %v, want %v", floatFlag.Get(), 2.718)
		}

		// 测试短标志解析
		cmd = NewCommand("test", "t", flag.ContinueOnError)
		var floatFlagShort flags.Float64Flag
		cmd.Float64Var(&floatFlagShort, "float-short", "fs", 3.14, "test float short flag")
		if err := cmd.Parse([]string{"-fs", "1.618"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if floatFlagShort.Get() != 1.618 {
			t.Errorf("after -fs, value = %v, want %v", floatFlagShort.Get(), 1.618)
		}
	})
}

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

// TestEnumVar 测试EnumVar方法的功能
func TestEnumVar(t *testing.T) {
	var enumFlag flags.EnumFlag
	var enumFlagShort flags.EnumFlag
	options := []string{"test", "dev", "prod"}
	defaultValue := "test"

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("EnumVar with nil pointer should panic")
			}
		}()
		cmd.EnumVar(nil, "enum", "en", defaultValue, "test enum flag", options)
	})

	// 测试正常功能(长标志)
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		cmd.EnumVar(&enumFlag, "enum", "en", defaultValue, "test enum flag", options)

		// 测试默认值
		if enumFlag.Get() != defaultValue {
			t.Errorf("default value = %q, want %q", enumFlag.Get(), defaultValue)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--enum", "prod"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if enumFlag.Get() != "prod" {
			t.Errorf("after --enum, value = %q, want %q", enumFlag.Get(), "prod")
		}
	})

	// 测试短标志解析
	t.Run("short flag", func(t *testing.T) {
		cmd := NewCommand("test-short", "ts", flag.ContinueOnError)
		cmd.EnumVar(&enumFlagShort, "enum-short", "e", defaultValue, "test enum short flag", options)
		if err := cmd.Parse([]string{"-e", "dev"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if enumFlagShort.Get() != "dev" {
			t.Errorf("after -e, value = %q, want %q", enumFlagShort.Get(), "dev")
		}
	})

	// 测试无效值解析
	t.Run("invalid input", func(t *testing.T) {
		cmd := NewCommand("test-invalid", "ti", flag.ContinueOnError)
		var enumFlagInvalid flags.EnumFlag
		cmd.EnumVar(&enumFlagInvalid, "enum-invalid", "ei", defaultValue, "test enum invalid flag", options)
		if err := cmd.Parse([]string{"--enum-invalid", "invalid"}); err == nil {
			t.Error("Parse with invalid enum value should return error")
		}
	})
}

// TestSliceVar 测试SliceVar方法的功能
func TestSliceVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("SliceVar with nil pointer should panic")
			}
		}()
		cmd.SliceVar(nil, "slice", "s", []string{"a", "b"}, "test slice flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var sliceFlag flags.SliceFlag
		cmd.SliceVar(&sliceFlag, "slice", "sl", []string{"a", "b"}, "test slice flag")

		// 测试默认值
		if !sliceFlag.Contains("a") || !sliceFlag.Contains("b") || sliceFlag.Len() != 2 {
			t.Errorf("default value = %v, want [a b]", sliceFlag.Get())
		}

		// 测试长标志解析（替换逻辑）
		if err := cmd.Parse([]string{"--slice", "c,d,e"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if sliceFlag.Len() != 3 || !sliceFlag.Contains("c") || !sliceFlag.Contains("d") || !sliceFlag.Contains("e") {
			t.Errorf("after --slice, value = %v, want [c d e]", sliceFlag.Get())
		}
	})

	// 测试短标志解析（替换逻辑）
	t.Run("short flag", func(t *testing.T) {
		cmd := NewCommand("test-short-slice", "tss", flag.ContinueOnError)
		var sliceFlagShort flags.SliceFlag
		cmd.SliceVar(&sliceFlagShort, "slice-short", "slss", []string{"a", "b"}, "test slice short flag")
		if err := cmd.Parse([]string{"-slss", "x,y"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if sliceFlagShort.Len() != 2 || !sliceFlagShort.Contains("x") || !sliceFlagShort.Contains("y") {
			t.Errorf("after -slss, value = %v, want [x y]", sliceFlagShort.Get())
		}
	})
}

// TestInt64Var 测试Int64Var方法的功能
func TestInt64Var(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Int64Var with nil pointer should panic")
			}
		}()
		cmd.Int64Var(nil, "int64", "i64", 123456789, "test int64 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var int64Flag flags.Int64Flag
		cmd.Int64Var(&int64Flag, "int64", "i64", 123456789, "test int64 flag")

		// 测试默认值
		if int64Flag.Get() != 123456789 {
			t.Errorf("default value = %d, want %d", int64Flag.Get(), 123456789)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--int64", "987654321"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if int64Flag.Get() != 987654321 {
			t.Errorf("after --int64, value = %d, want %d", int64Flag.Get(), 987654321)
		}

		// 测试短标志解析
		t.Run("short flag", func(t *testing.T) {
			cmd := NewCommand("test-short-int64", "tsi", flag.ContinueOnError)
			var int64FlagShort flags.Int64Flag
			cmd.Int64Var(&int64FlagShort, "int64-short", "i64s", 123456789, "test int64 short flag")
			if err := cmd.Parse([]string{"-i64s", "111222333"}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if int64FlagShort.Get() != 111222333 {
				t.Errorf("after -i64s, value = %d, want %d", int64FlagShort.Get(), 111222333)
			}
		})
	})
}

// TestUint16Var 测试Uint16Var方法的功能
func TestUint16Var(t *testing.T) {
	var uint16Flag flags.Uint16Flag
	var uint16FlagShort flags.Uint16Flag

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Uint16Var with nil pointer should panic")
			}
		}()
		cmd.Uint16Var(nil, "uint16", "u16", 65535, "test uint16 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		cmd.Uint16Var(&uint16Flag, "uint16", "u16", 65535, "test uint16 flag")

		// 测试默认值
		if uint16Flag.Get() != 65535 {
			t.Errorf("default value = %d, want %d", uint16Flag.Get(), 65535)
		}

		// 测试长标志解析
		if err := cmd.Parse([]string{"--uint16", "32768"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if uint16Flag.Get() != 32768 {
			t.Errorf("after --uint16, value = %d, want %d", uint16Flag.Get(), 32768)
		}

		// 测试短标志解析
		t.Run("short flag", func(t *testing.T) {
			cmdShort := NewCommand("test-short-uint16-new", "tsun", flag.ContinueOnError)
			cmdShort.Uint16Var(&uint16FlagShort, "uint16-short", "u16ss", 65535, "test uint16 short flag")
			if err := cmdShort.Parse([]string{"-u16ss", "12345"}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			if uint16FlagShort.Get() != 12345 {
				t.Errorf("after -u16ss, value = %d, want %d", uint16FlagShort.Get(), 12345)
			}
		})

		// 测试无效值解析
		t.Run("invalid input", func(t *testing.T) {
			cmdInvalid := NewCommand("test-uint16-invalid", "tui", flag.ContinueOnError)
			var uint16FlagInvalid flags.Uint16Flag
			cmdInvalid.Uint16Var(&uint16FlagInvalid, "uint16-invalid", "u16i", 65535, "test uint16 invalid flag")

			// 重定向标准输出和错误到缓冲区
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			err := cmdInvalid.Parse([]string{"--uint16-invalid", "65536"})

			// 恢复标准输出和错误
			w.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			// 读取缓冲区内容
			var buf bytes.Buffer
			_, copyErr := io.Copy(&buf, r)
			if copyErr != nil {
				t.Errorf("Failed to copy output: %v", copyErr)
			}

			// 仅在详细模式下打印输出
			if testing.Verbose() {
				t.Logf("Command output: %s", buf.String())
			}

			if err == nil {
				t.Error("Parse with value 65536 should return error")
			}
		})
	})
}

// TestTimeVar 测试TimeVar方法的功能
func TestTimeVar(t *testing.T) {
	// 测试指针为nil的情况
	{
		t.Run("nil pointer", func(t *testing.T) {
			cmd := NewCommand("test", "t", flag.ContinueOnError)
			defer func() {
				if r := recover(); r == nil {
					t.Error("TimeVar with nil pointer should panic")
				}
			}()
			cmd.TimeVar(nil, "time", "tm", time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), "test time flag")
		})
	}

	// 测试正常功能
	{
		t.Run("normal case", func(t *testing.T) {
			cmd := NewCommand("test", "t", flag.ContinueOnError)
			var timeFlag flags.TimeFlag
			defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
			cmd.TimeVar(&timeFlag, "time", "tm", defaultTime, "test time flag")

			// 测试默认值
			if !timeFlag.Get().Equal(defaultTime) {
				t.Errorf("default value = %v, want %v", timeFlag.Get(), defaultTime)
			}

			// 测试长标志解析
			inputTime := "2023-12-31T23:59:59Z"
			if err := cmd.Parse([]string{"--time", inputTime}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			parsedTime, _ := time.Parse(time.RFC3339, inputTime)
			if !timeFlag.Get().Equal(parsedTime) {
				t.Errorf("after --time, value = %v, want %v", timeFlag.Get(), parsedTime)
			}

			// 测试短标志解析
			cmd = NewCommand("test-short", "ts", flag.ContinueOnError)
			var timeFlagShort flags.TimeFlag
			cmd.TimeVar(&timeFlagShort, "time-short", "t", defaultTime, "test time short flag")
			shortInput := "2024-01-01"
			if err := cmd.Parse([]string{"-t", shortInput}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}
			shortParsed, _ := time.Parse("2006-01-02", shortInput)
			if !timeFlagShort.Get().Equal(shortParsed) {
				t.Errorf("after -t, value = %v, want %v", timeFlagShort.Get(), shortParsed)
			}
		})
	}

	// 测试无效格式
	{
		t.Run("invalid format", func(t *testing.T) {
			cmd := NewCommand("test", "t", flag.ContinueOnError)
			var timeFlag flags.TimeFlag
			defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
			cmd.TimeVar(&timeFlag, "time", "tm", defaultTime, "test time flag")

			// 捕获标准错误
			oldStderr := os.Stderr
			_, w, _ := os.Pipe()
			os.Stderr = w

			defer func() {
				w.Close()
				os.Stderr = oldStderr
			}()

			err := cmd.Parse([]string{"--time", "invalid-time"})
			if err == nil {
				t.Fatal("expected error for invalid time format")
			}

			// 验证默认值不变
			if !timeFlag.Get().Equal(defaultTime) {
				t.Errorf("default value changed after invalid input: %v", timeFlag.Get())
			}
		})
	}
}

// TestTime 测试Time方法的功能
func TestTime(t *testing.T) {
	cmd := NewCommand("test", "t", flag.ContinueOnError)
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	flag := cmd.Time("time", "tm", defaultTime, "test time flag")

	// 验证默认值
	if !flag.Get().Equal(defaultTime) {
		t.Errorf("default value = %v, want %v", flag.Get(), defaultTime)
	}

	// 验证解析功能
	inputTime := "2023-06-15T10:30:00+08:00"
	if err := cmd.Parse([]string{"--time", inputTime}); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	parsedTime, _ := time.Parse(time.RFC3339, inputTime)
	if !flag.Get().Equal(parsedTime) {
		t.Errorf("parsed value = %v, want %v", flag.Get(), parsedTime)
	}
}
