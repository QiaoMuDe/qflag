package cmd

import (
	"flag"
	"gitee.com/MM-Q/qflag/flags"
	"testing"
)

func TestString(t *testing.T) {
	cmd := NewCommand("test", "t", flag.ContinueOnError)
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
	cmd := NewCommand("test", "t", flag.ContinueOnError)
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
	cmd := NewCommand("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to StringVar")
		}
	}()
	cmd.StringVar(nil, "str", "s", "default", "test")
}
