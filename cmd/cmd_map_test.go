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
		cmd := NewCommand("test", "t", flag.ContinueOnError)
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
		cmd := NewCommand("test", "t", flag.ContinueOnError)
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
		cmd = NewCommand("test-short", "ts", flag.ContinueOnError)
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
		cmd := NewCommand("test-delimiters", "td", flag.ContinueOnError)
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
