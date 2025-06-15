package qflag

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"testing"
)

// TestStringFlagLong 测试字符串类型长标志的注册和解析
func TestStringFlagLong(t *testing.T) {
	// 完全重定向标准输出到缓冲区
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 不输出任何捕获的内容
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Errorf("ReadFrom error: %v", err)
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "string-flag"
	defValue := "default"
	usage := "test string flag"

	// 测试String方法(仅长标志)
	f := cmd.String(flagName, "sf", defValue, usage)
	if f == nil {
		t.Error("String() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "test-value"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != "test-value" {
		t.Errorf("String flag value = %q, want %q", *f.value, "test-value")
	}
}

// TestStringFlagShort 测试字符串类型短标志的注册和解析
func TestStringFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "s"
	defValue := "default"
	usage := "test string flag"

	// 测试String方法(仅短标志)
	f := cmd.String("sf", shortName, defValue, usage)
	if f == nil {
		t.Error("String() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "test-value"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != "test-value" {
		t.Errorf("String flag value = %q, want %q", *f.value, "test-value")
	}
}

// TestIntFlagLong 测试整数类型长标志的注册和解析
func TestIntFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "int-flag"
	defValue := 100
	usage := "test int flag"

	// 测试Int方法(仅长标志)
	f := cmd.Int(flagName, "if", defValue, usage)
	if f == nil {
		t.Error("Int() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "200"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 200 {
		t.Errorf("Int flag value = %d, want %d", *f.value, 200)
	}
}

// TestIntFlagShort 测试整数类型短标志的注册和解析
func TestIntFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "i"
	defValue := 100
	usage := "test int flag"

	// 测试Int方法(仅短标志)
	f := cmd.Int("ci", shortName, defValue, usage)
	if f == nil {
		t.Error("Int() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "200"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 200 {
		t.Errorf("Int flag value = %d, want %d", *f.value, 200)
	}
}

// TestBoolFlagLong 测试布尔类型长标志的注册和解析
func TestBoolFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "bool-flag"
	defValue := false
	usage := "test bool flag"

	// 测试Bool方法(仅长标志)
	f := cmd.Bool(flagName, "bl", defValue, usage)
	if f == nil {
		t.Error("Bool() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != true {
		t.Errorf("Bool flag value = %v, want %v", *f.value, true)
	}
}

// TestBoolFlagShort 测试布尔类型短标志的注册和解析
func TestBoolFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "b"
	defValue := false
	usage := "test bool flag"

	// 测试Bool方法(仅短标志)
	f := cmd.Bool("ct", shortName, defValue, usage)
	if f == nil {
		t.Error("Bool() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != true {
		t.Errorf("Bool flag value = %v, want %v", *f.value, true)
	}
}

// TestFloatFlagLong 测试浮点数类型长标志的注册和解析
func TestFloatFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "float-flag"
	defValue := 3.14
	usage := "test float flag"

	// 测试Float方法(仅长标志)
	f := cmd.Float(flagName, "ff", defValue, usage)
	if f == nil {
		t.Error("Float() returned nil")
	}

	// 测试长标志解析
	err := cmd.Parse([]string{"--" + flagName, "6.28"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 6.28 {
		t.Errorf("Float flag value = %v, want %v", *f.value, 6.28)
	}
}

// TestFloatFlagShort 测试浮点数类型短标志的注册和解析
func TestFloatFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "f"
	defValue := 3.14
	usage := "test float flag"

	// 测试Float方法(仅短标志)
	f := cmd.Float("cf", shortName, defValue, usage)
	if f == nil {
		t.Error("Float() returned nil")
	}

	// 测试短标志解析
	err := cmd.Parse([]string{"-" + shortName, "6.28"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	if *f.value != 6.28 {
		t.Errorf("Float flag value = %v, want %v", *f.value, 6.28)
	}
}

// TestSliceFlagLong 测试字符串切片类型长标志的注册和解析
func TestSliceFlagLong(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	flagName := "slice-flag"
	defValue := []string{"default1", "default2"}
	usage := "test slice flag"

	// 测试Slice方法(长标志)
	f := cmd.Slice(flagName, "sf", defValue, usage)
	if f == nil {
		t.Error("Slice() returned nil")
	}

	// 测试长标志解析(多个值)
	err := cmd.Parse([]string{"--" + flagName, "value1", "--" + flagName, "value2"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	expected := []string{"value1", "value2"}
	if !reflect.DeepEqual(*f.value, expected) {
		t.Errorf("Slice flag value = %q, want %q", *f.value, expected)
	}
}

// TestSliceFlagShort 测试字符串切片类型短标志的注册和解析
func TestSliceFlagShort(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	shortName := "s"
	defValue := []string{"default"}
	usage := "test slice flag"

	// 测试Slice方法(短标志)
	f := cmd.Slice("slice", shortName, defValue, usage)
	if f == nil {
		t.Error("Slice() returned nil")
	}

	// 测试短标志解析(多个值)
	err := cmd.Parse([]string{"-" + shortName, "value1", "-" + shortName, "value2", "-" + shortName, "value3"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证值
	expected := []string{"value1", "value2", "value3"}
	if !reflect.DeepEqual(*f.value, expected) {
		t.Errorf("Slice flag value = %q, want %q", *f.value, expected)
	}
}

// TestSliceFlagDefault 测试字符串切片类型标志的默认值
func TestSliceFlagDefault(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r); err != nil {
				t.Errorf("ReadFrom error: %v", err)
			} else {
				t.Logf("Captured output:\n%s", buf.String())
			}
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defValue := []string{"default1", "default2"}

	// 测试Slice方法
	f := cmd.Slice("slice-def", "sd", defValue, "test slice default value")
	if f == nil {
		t.Error("Slice() returned nil")
	}

	// 不传递任何参数，使用默认值
	err := cmd.Parse([]string{})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	// 验证默认值
	if !reflect.DeepEqual(*f.value, defValue) {
		t.Errorf("Slice flag default value = %q, want %q", *f.value, defValue)
	}
}

// TestParseError 测试参数解析错误
func TestParseError(t *testing.T) {
	// 捕获标准输出和错误输出
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr
	defer func() {
		wOut.Close()
		wErr.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr

		// 只有在-v模式或测试失败时输出
		if testing.Verbose() || t.Failed() {
			outBuf := new(bytes.Buffer)
			errBuf := new(bytes.Buffer)
			if _, err := outBuf.ReadFrom(rOut); err != nil {
				t.Errorf("ReadFrom stdout error: %v", err)
			}
			if _, err := errBuf.ReadFrom(rErr); err != nil {
				t.Errorf("ReadFrom stderr error: %v", err)
			}
			t.Logf("Captured stdout:\n%s", outBuf.String())
			t.Logf("Captured stderr:\n%s", errBuf.String())
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.Int("int-flag", "i", 0, "test int flag")

	// 测试无效参数
	err := cmd.Parse([]string{"--int-flag", "not-a-number"})
	if err == nil {
		t.Error("Parse() should return error for invalid input")
	}
}

// TestHelpFlag 测试帮助标志
func TestHelpFlag(t *testing.T) {
	// 捕获标准输出
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		// 只在测试失败时输出捕获的内容
		if t.Failed() {
			var buf bytes.Buffer
			buf.ReadFrom(r)
			t.Logf("Captured output:\n%s", buf.String())
		}
	}()

	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.String("string-flag", "s", "", "test string flag")

	// 测试帮助标志
	err := cmd.Parse([]string{"--help"})
	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}
}
