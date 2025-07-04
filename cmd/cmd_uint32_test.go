package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

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
