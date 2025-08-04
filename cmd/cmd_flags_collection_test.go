package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestMapVar 测试MapVar方法的功能
func TestMapVar(t *testing.T) {

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("MapVar with nil pointer should panic")
			}
		}()
		cmd.MapVar(nil, "map", "m", map[string]string{"key": "val"}, "test map flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		defaultMap := map[string]string{"default": "value"}
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		var mapFlag flags.MapFlag
		cmd.MapVar(&mapFlag, "map", "m", defaultMap, "test map flag")

		// 测试默认值
		if val, ok := mapFlag.Get()["default"]; !ok || val != "value" {
			t.Errorf("default value = %v, want %v", mapFlag.Get(), defaultMap)
		}

		// 测试长标志解析（单个键值对）
		if err := cmd.Parse([]string{"--map", "name=test"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if val, ok := mapFlag.Get()["name"]; !ok || val != "test" {
			t.Errorf("after --map, value = %v, want {name:test}", mapFlag.Get())
		}

		// 测试短标志解析（多个键值对）
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
		var mapFlagShort flags.MapFlag
		cmd.MapVar(&mapFlagShort, "map-short", "m", nil, "test map short flag")
		if err := cmd.Parse([]string{"-m", "key1=val1,key2=val2"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		result := mapFlagShort.Get()
		if result["key1"] != "val1" || result["key2"] != "val2" {
			t.Errorf("after -m, value = %v, want {key1:val1, key2:val2}", result)
		}
	})

	// 测试自定义分隔符
	t.Run("custom delimiters", func(t *testing.T) {
		cmd := NewCmd("test-delimiters", "td", flag.ContinueOnError)
		var mapFlag flags.MapFlag
		cmd.MapVar(&mapFlag, "map-delim", "d", nil, "test map delimiters")
		mapFlag.SetDelimiters(flags.FlagSplitSemicolon, flags.FlagKVColon)

		if err := cmd.Parse([]string{"--map-delim", "key1:val1;key2:val2"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		result := mapFlag.Get()
		if result["key1"] != "val1" || result["key2"] != "val2" {
			t.Errorf("after custom delimiters, value = %v, want {key1:val1, key2:val2}", result)
		}
	})
}

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

// TestSliceVarf 测试SliceVar方法的功能
func TestSliceVarf(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("SliceVar with nil pointer should panic")
			}
		}()
		cmd.SliceVar(nil, "slice", "s", []string{"a", "b"}, "test slice flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd := NewCmd("test-short-slice", "tss", flag.ContinueOnError)
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
