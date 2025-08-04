package cmd

import (
	"bytes"
	"flag"
	"io"
	"os"
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

// TestUint16 测试Uint16方法的功能正确性
func TestUint16(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defaultValue := uint16(32767)
	uint16Flag := cmd.Uint16("port", "p", defaultValue, "uint16 flag test")

	// 验证默认值
	if uint16Flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, uint16Flag.Get())
	}

	// 测试长标志解析
	args := []string{"--port", "8080"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if uint16Flag.Get() != 8080 {
		t.Errorf("Expected parsed value 8080, got %d", uint16Flag.Get())
	}
}

// TestUint16Var 测试Uint16Var方法的功能正确性
func TestUint16Var(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var uint16Flag flags.Uint16Flag
	defaultValue := uint16(1024)
	cmd.Uint16Var(&uint16Flag, "port", "p", defaultValue, "uint16 flag test")

	// 验证默认值
	if uint16Flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, uint16Flag.Get())
	}

	// 测试短标志解析
	args := []string{"-p", "443"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if uint16Flag.Get() != 443 {
		t.Errorf("Expected parsed value 443, got %d", uint16Flag.Get())
	}
}

// TestUint16Var_NilPointer 测试Uint16Var方法传入nil指针时的错误处理
func TestUint16Var_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint16Var")
		}
	}()
	cmd.Uint16Var(nil, "port", "p", 0, "test")
}

// TestUint32 测试Uint32方法的功能正确性
func TestUint32(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defaultValue := uint32(50000)
	uint32Flag := cmd.Uint32("count", "c", defaultValue, "uint32 flag test")

	// 验证默认值
	if uint32Flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, uint32Flag.Get())
	}

	// 测试长标志解析
	args := []string{"--count", "100000"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if uint32Flag.Get() != 100000 {
		t.Errorf("Expected parsed value 100000, got %d", uint32Flag.Get())
	}
}

// TestUint32Var 测试Uint32Var方法的功能正确性
func TestUint32Var(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var uint32Flag flags.Uint32Flag
	defaultValue := uint32(30000)
	cmd.Uint32Var(&uint32Flag, "count", "c", defaultValue, "uint32 flag test")

	// 验证默认值
	if uint32Flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, uint32Flag.Get())
	}

	// 测试短标志解析
	args := []string{"-c", "60000"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if uint32Flag.Get() != 60000 {
		t.Errorf("Expected parsed value 60000, got %d", uint32Flag.Get())
	}
}

// TestUint32Var_NilPointer 测试Uint32Var方法传入nil指针时的错误处理
func TestUint32Var_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint32Var")
		}
	}()
	cmd.Uint32Var(nil, "count", "c", 0, "test")
}

// TestUint64 测试Uint64方法的功能正确性
func TestUint64(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defaultValue := uint64(100)
	uint64Flag := cmd.Uint64("number", "n", defaultValue, "uint64 flag test")

	// 验证默认值
	if uint64Flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, uint64Flag.Get())
	}

	// 测试长标志解析
	args := []string{"--number", "200"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if uint64Flag.Get() != 200 {
		t.Errorf("Expected parsed value 200, got %d", uint64Flag.Get())
	}
}

// TestUint64Var 测试Uint64Var方法的功能正确性
func TestUint64Var(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var uint64Flag flags.Uint64Flag
	defaultValue := uint64(150)
	cmd.Uint64Var(&uint64Flag, "number", "n", defaultValue, "uint64 flag test")

	// 验证默认值
	if uint64Flag.Get() != defaultValue {
		t.Errorf("Expected default value %d, got %d", defaultValue, uint64Flag.Get())
	}

	// 测试短标志解析
	args := []string{"-n", "300"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if uint64Flag.Get() != 300 {
		t.Errorf("Expected parsed value 300, got %d", uint64Flag.Get())
	}
}

// TestUint64Var_NilPointer 测试Uint64Var方法传入nil指针时的错误处理
func TestUint64Var_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint64Var")
		}
	}()
	cmd.Uint64Var(nil, "number", "n", 0, "test")
}

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

// TestIntVarf 测试IntVar方法的功能
func TestIntVarf(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("IntVar with nil pointer should panic")
			}
		}()
		cmd.IntVar(nil, "int", "i", 123, "test int flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd = NewCmd("test", "t", flag.ContinueOnError)
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

// TestFloatVar 测试FloatVar方法的功能
func TestFloatVar(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("FloatVar with nil pointer should panic")
			}
		}()
		cmd.Float64Var(nil, "float", "f", 3.14, "test float flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd = NewCmd("test", "t", flag.ContinueOnError)
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

// TestInt64Varg 测试Int64Var方法的功能
func TestInt64Varf(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Int64Var with nil pointer should panic")
			}
		}()
		cmd.Int64Var(nil, "int64", "i64", 123456789, "test int64 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
			cmd := NewCmd("test-short-int64", "tsi", flag.ContinueOnError)
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

// TestUint16Varf 测试Uint16Var方法的功能
func TestUint16Varf(t *testing.T) {
	var uint16Flag flags.Uint16Flag
	var uint16FlagShort flags.Uint16Flag

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Uint16Var with nil pointer should panic")
			}
		}()
		cmd.Uint16Var(nil, "uint16", "u16", 65535, "test uint16 flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
			cmdShort := NewCmd("test-short-uint16-new", "tsun", flag.ContinueOnError)
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
			cmdInvalid := NewCmd("test-uint16-invalid", "tui", flag.ContinueOnError)
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
			if closeErr := w.Close(); closeErr != nil {
				t.Errorf("failed to close writer: %v", closeErr)
			}
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
