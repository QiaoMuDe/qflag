package cmd

import (
	"flag"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestUint16 测试Uint16方法的功能正确性
func TestUint16(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
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
	cmd := NewCmd("test", "t", flag.ContinueOnError)
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
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint16Var")
		}
	}()
	cmd.Uint16Var(nil, "port", "p", 0, "test")
}

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

// TestUint64 测试Uint64方法的功能正确性
func TestUint64(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
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
	cmd := NewCmd("test", "t", flag.ContinueOnError)
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
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to Uint64Var")
		}
	}()
	cmd.Uint64Var(nil, "number", "n", 0, "test")
}
