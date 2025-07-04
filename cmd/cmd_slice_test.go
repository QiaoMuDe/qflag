package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestSliceVar 测试SliceVar方法的功能
func TestSliceVar(t *testing.T) {

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("SliceVar with nil pointer should panic")
			}
		}()
		cmd.SliceVar(nil, "slice", "s", []string{"a"}, "test slice flag")
	})

	// 测试默认值
	t.Run("default value", func(t *testing.T) {
		// 非nil默认值
		t.Run("non-nil default", func(t *testing.T) {
			defaultSlice := []string{"default1", "default2"}
			cmd := NewCmd("test", "t", flag.ContinueOnError)
			var sliceFlag flags.SliceFlag
			cmd.SliceVar(&sliceFlag, "slice", "s", defaultSlice, "test slice default")

			result := sliceFlag.Get()
			if len(result) != 2 || result[0] != "default1" || result[1] != "default2" {
				t.Errorf("default value = %v, want %v", result, defaultSlice)
			}
		})

		// nil默认值（应初始化为空切片）
		t.Run("nil default", func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)
			var sliceFlag flags.SliceFlag
			cmd.SliceVar(&sliceFlag, "slice-nil", "n", nil, "test nil default")

			result := sliceFlag.Get()
			if len(result) != 0 {
				t.Errorf("nil default should be empty slice, got %v", result)
			}
		})
	})

	// 测试标志解析
	t.Run("flag parsing", func(t *testing.T) {
		// 长标志单值解析
		t.Run("long flag single value", func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)
			var sliceFlag flags.SliceFlag
			cmd.SliceVar(&sliceFlag, "slice", "s", nil, "test long flag")

			if err := cmd.Parse([]string{"--slice", "value1"}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			result := sliceFlag.Get()
			if len(result) != 1 || result[0] != "value1" {
				t.Errorf("after --slice, value = %v, want [value1]", result)
			}
		})

		// 长标志多值解析（逗号分隔）
		t.Run("long flag multiple values", func(t *testing.T) {
			cmd := NewCmd("test", "t", flag.ContinueOnError)
			var sliceFlag flags.SliceFlag
			cmd.SliceVar(&sliceFlag, "slice", "s", nil, "test long flag multiple values")

			if err := cmd.Parse([]string{"--slice", "value1,value2,value3"}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			result := sliceFlag.Get()
			expected := []string{"value1", "value2", "value3"}
			if len(result) != len(expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(expected))
			}
			for i, v := range expected {
				if result[i] != v {
					t.Errorf("index %d: got %s, want %s", i, result[i], v)
				}
			}
		})

		// 短标志解析
		t.Run("short flag", func(t *testing.T) {
			cmd := NewCmd("test-short", "ts", flag.ContinueOnError)
			var sliceFlag flags.SliceFlag
			cmd.SliceVar(&sliceFlag, "slice-short", "s", nil, "test short flag")

			if err := cmd.Parse([]string{"-s", "a,b,c"}); err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			result := sliceFlag.Get()
			expected := []string{"a", "b", "c"}
			if len(result) != len(expected) {
				t.Fatalf("length mismatch: got %d, want %d", len(result), len(expected))
			}
			for i, v := range expected {
				if result[i] != v {
					t.Errorf("index %d: got %s, want %s", i, result[i], v)
				}
			}
		})

		// 空输入处理
		t.Run("empty input", func(t *testing.T) {
			cmd := NewCmd("test-empty", "te", flag.ContinueOnError)
			var sliceFlag flags.SliceFlag
			cmd.SliceVar(&sliceFlag, "slice-empty", "e", []string{"default"}, "test empty input")

			if err := cmd.Parse([]string{"--slice-empty", ""}); err == nil {
				t.Fatal("expected error when parsing empty input, got nil")
			}
		})
	})
}
