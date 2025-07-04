package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

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
