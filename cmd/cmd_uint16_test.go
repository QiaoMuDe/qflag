package cmd

import (
	"flag"
	"gitee.com/MM-Q/qflag/flags"
	"testing"
)

// TestUint16 测试Uint16方法的功能正确性
func TestUint16(t *testing.T) {
	cmd := NewCommand("test", "t", flag.ContinueOnError)
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
	cmd := NewCommand("test", "t", flag.ContinueOnError)
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
	cmd := NewCommand("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint16Var")
		}
	}()
	cmd.Uint16Var(nil, "port", "p", 0, "test")
}
