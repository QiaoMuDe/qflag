package cmd

import (
	"flag"
	"os"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

func TestTime(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	timeFlag := cmd.Time("time", "t", defaultTime, "time flag test")

	// 测试默认值
	if !timeFlag.Get().Equal(defaultTime) {
		t.Errorf("Expected default time %v, got %v", defaultTime, timeFlag.Get())
	}

	// 解析参数
	args := []string{"--time", "2023-12-31T23:59:59Z"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	if !timeFlag.Get().Equal(parsedTime) {
		t.Errorf("Expected parsed time %v, got %v", parsedTime, timeFlag.Get())
	}
}

func TestTimeVar(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var timeFlag flags.TimeFlag
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	cmd.TimeVar(&timeFlag, "time", "t", defaultTime, "time flag test")

	if !timeFlag.Get().Equal(defaultTime) {
		t.Errorf("Expected default time %v, got %v", defaultTime, timeFlag.Get())
	}

	args := []string{"-t", "2024-01-01T12:00:00+08:00"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, "2024-01-01T12:00:00+08:00")
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	if !timeFlag.Get().Equal(parsedTime) {
		t.Errorf("Expected parsed time %v, got %v", parsedTime, timeFlag.Get())
	}
}

func TestTimeVar_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to TimeVar")
		}
	}()
	cmd.TimeVar(nil, "time", "t", time.Time{}, "test")
}

// TestDurationVar 测试DurationVar方法的功能
func TestDurationVar(t *testing.T) {
	var durationFlagShort flags.DurationFlag

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("DurationVar with nil pointer should panic")
			}
		}()
		cmd.DurationVar(nil, "duration", "dur", time.Second*5, "test duration flag")
	})

	// 测试正常功能(长标志)
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd := NewCmd("test-short", "ts", flag.ContinueOnError)
		cmd.DurationVar(&durationFlagShort, "duration-short", "d", time.Second*5, "test duration short flag")
		if err := cmd.Parse([]string{"-d", "2m"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if durationFlagShort.Get() != time.Minute*2 {
			t.Errorf("after -d, value = %v, want %v", durationFlagShort.Get(), time.Minute*2)
		}
	})
}

// TestDurationVarf 测试DurationVar方法的功能
func TestDurationVarf(t *testing.T) {
	var durationFlagShort flags.DurationFlag

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("DurationVar with nil pointer should panic")
			}
		}()
		cmd.DurationVar(nil, "duration", "dur", time.Second*5, "test duration flag")
	})

	// 测试正常功能(长标志)
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd := NewCmd("test-short", "ts", flag.ContinueOnError)
		cmd.DurationVar(&durationFlagShort, "duration-short", "d", time.Second*5, "test duration short flag")
		if err := cmd.Parse([]string{"-d", "2m"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if durationFlagShort.Get() != time.Minute*2 {
			t.Errorf("after -d, value = %v, want %v", durationFlagShort.Get(), time.Minute*2)
		}
	})
}

// TestTimeVarf 测试TimeVar方法的功能
func TestTimeVarf(t *testing.T) {
	// 测试指针为nil的情况
	{
		t.Run("nil pointer", func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)
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
			cmd := NewCmd("test", "t", flag.ContinueOnError)
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
			cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
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
			cmd := NewCmd("test", "t", flag.ContinueOnError)
			var timeFlag flags.TimeFlag
			defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
			cmd.TimeVar(&timeFlag, "time", "tm", defaultTime, "test time flag")

			// 捕获标准错误
			oldStderr := os.Stderr
			_, w, _ := os.Pipe()
			os.Stderr = w

			defer func() {
				if err := w.Close(); err != nil {
					t.Errorf("failed to close writer: %v", err)
				}
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

// TestTimef 测试Time方法的功能
func TestTimef(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
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
