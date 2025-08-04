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
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("BoolVar with nil pointer should panic")
			}
		}()
		cmd.BoolVar(nil, "bool", "b", false, "test bool flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
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

// TestEnumVar 测试EnumVar方法的功能
func TestEnumVar(t *testing.T) {
	options := []string{"option1", "option2", "option3"}

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test-enum-nil", "ten", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("EnumVar with nil pointer should panic")
			}
		}()
		cmd.EnumVar(nil, "enumq-nil", "q", "option1", "test enum nil flag", options)
	})

	// 测试正常功能(长标志)
	t.Run("normal case", func(t *testing.T) {
		cmd1 := NewCmd("test-enum-normal", "tenm", flag.ContinueOnError)
		var enumFlag flags.EnumFlag
		cmd1.EnumVar(&enumFlag, "enumq-test-normal", "m", "option1", "test enum normal flag", options)

		// 测试默认值
		if enumFlag.Get() != "option1" {
			t.Errorf("default value = %v, want %v", enumFlag.Get(), "option1")
		}

		// 测试长标志解析-有效选项
		if err := cmd1.Parse([]string{"--enumq-test-normal", "option2"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if enumFlag.Get() != "option2" {
			t.Errorf("after --enumtest-normal, value = %v, want %v", enumFlag.Get(), "option2")
		}

		// 测试长标志解析-无效选项
		var enumFlagInvalid flags.EnumFlag
		cmdInvalid := NewCmd("test-enum-invalid", "tein", flag.ContinueOnError)
		cmdInvalid.EnumVar(&enumFlagInvalid, "enumq-invalid", "a", "option1", "test enum invalid flag", options)
		if err := cmdInvalid.Parse([]string{"--enumq-invalid", "invalid"}); err == nil {
			t.Error("Parse should fail with invalid enum value")
		}
	})

	// 测试短标志解析
	t.Run("short flag", func(t *testing.T) {
		var enumFlagShort flags.EnumFlag
		cmd := NewCmd("test-enum-short", "tes", flag.ContinueOnError)
		cmd.EnumVar(&enumFlagShort, "enumq-short", "b", "option1", "test enum short flag", options)
		if err := cmd.Parse([]string{"--enumq-short", "option3"}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if enumFlagShort.Get() != "option3" {
			t.Errorf("after --enumtest-short, value = %v, want %v", enumFlagShort.Get(), "option3")
		}
	})
}

func TestString(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	strFlag := cmd.String("str", "s", "default", "string flag test")

	// 测试默认值
	if strFlag.Get() != "default" {
		t.Errorf("Expected default value 'default', got %s", strFlag.Get())
	}

	// 解析参数
	args := []string{"--str", "testvalue"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if strFlag.Get() != "testvalue" {
		t.Errorf("Expected 'testvalue', got %s", strFlag.Get())
	}
}

func TestStringVar(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var strFlag flags.StringFlag
	cmd.StringVar(&strFlag, "str", "s", "default", "string flag test")

	if strFlag.Get() != "default" {
		t.Errorf("Expected default value 'default', got %s", strFlag.Get())
	}

	args := []string{"-s", "shortvalue"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if strFlag.Get() != "shortvalue" {
		t.Errorf("Expected 'shortvalue', got %s", strFlag.Get())
	}
}

func TestStringVar_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to StringVar")
		}
	}()
	cmd.StringVar(nil, "str", "s", "default", "test")
}

// TestStringVarf 测试StringVar方法的功能
func TestStringVarf(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("StringVar with nil pointer should panic")
			}
		}()
		cmd.StringVar(nil, "str", "s", "default", "test string flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd = NewCmd("test-short", "ts", flag.ContinueOnError)
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

// TestBoolVarf 测试BoolVar方法的功能
func TestBoolVarf(t *testing.T) {
	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("BoolVar with nil pointer should panic")
			}
		}()
		cmd.BoolVar(nil, "bool", "b", false, "test bool flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd = NewCmd("test", "t", flag.ContinueOnError)
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

// TestEnumVarf 测试EnumVar方法的功能
func TestEnumVarf(t *testing.T) {
	var enumFlag flags.EnumFlag
	var enumFlagShort flags.EnumFlag
	options := []string{"test", "dev", "prod"}
	defaultValue := "test"

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("EnumVar with nil pointer should panic")
			}
		}()
		cmd.EnumVar(nil, "enum", "en", defaultValue, "test enum flag", options)
	})

	// 测试正常功能(长标志)
	t.Run("normal case", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
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
		cmd := NewCmd("test-short", "ts", flag.ContinueOnError)
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
		cmd := NewCmd("test-invalid", "ti", flag.ContinueOnError)
		var enumFlagInvalid flags.EnumFlag
		cmd.EnumVar(&enumFlagInvalid, "enum-invalid", "ei", defaultValue, "test enum invalid flag", options)
		if err := cmd.Parse([]string{"--enum-invalid", "invalid"}); err == nil {
			t.Error("Parse with invalid enum value should return error")
		}
	})
}
