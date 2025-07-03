package cmd

import (
	"flag"
	"gitee.com/MM-Q/qflag/flags"
	"testing"
)

// TestUint64 测试Uint64方法的功能正确性
func TestUint64(t *testing.T) {
	cmd := NewCommand("test", "t", flag.ContinueOnError)
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
	cmd := NewCommand("test", "t", flag.ContinueOnError)
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
	cmd := NewCommand("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint64Var")
		}
	}()
	cmd.Uint64Var(nil, "number", "n", 0, "test")
}
