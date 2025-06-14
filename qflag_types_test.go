package qflag

import (
	"testing"
)

// TestCmd_Name 测试Cmd的Name方法
func TestCmd_Name(t *testing.T) {
	cmd := &Cmd{name: "testcmd"}
	if cmd.Name() != "testcmd" {
		t.Errorf("Name() = %q, want %q", cmd.Name(), "testcmd")
	}
}

// TestCmd_ShortName 测试Cmd的ShortName方法
func TestCmd_ShortName(t *testing.T) {
	cmd := &Cmd{shortName: "tc"}
	if cmd.ShortName() != "tc" {
		t.Errorf("ShortName() = %q, want %q", cmd.ShortName(), "tc")
	}
}

// TestCmd_Description 测试Cmd的Description和SetDescription方法
func TestCmd_Description(t *testing.T) {
	cmd := &Cmd{}
	desc := "test description"
	cmd.SetDescription(desc)
	if cmd.Description() != desc {
		t.Errorf("Description() = %q, want %q", cmd.Description(), desc)
	}
}

// TestCmd_Usage 测试Cmd的Usage和SetUsage方法
func TestCmd_Usage(t *testing.T) {
	cmd := &Cmd{}
	usage := "test usage"
	cmd.SetUsage(usage)
	if cmd.Usage() != usage {
		t.Errorf("Usage() = %q, want %q", cmd.Usage(), usage)
	}
}

// TestCmd_Args 测试Cmd的Args方法
func TestCmd_Args(t *testing.T) {
	args := []string{"arg1", "arg2"}
	cmd := &Cmd{args: args}
	result := cmd.Args()
	// 检查长度
	if len(result) != len(args) {
		t.Fatalf("Args() length = %d, want %d", len(result), len(args))
	}
	// 检查每个元素
	for i, arg := range args {
		if result[i] != arg {
			t.Errorf("Args()[%d] = %q, want %q", i, result[i], arg)
		}
	}
}

// TestCmd_Arg 测试Cmd的Arg方法
func TestCmd_Arg(t *testing.T) {
	cmd := &Cmd{args: []string{"arg0", "arg1", "arg2"}}
	tests := []struct {
		name string
		i    int
		want string
	}{{
		name: "valid index 0",
		i:    0,
		want: "arg0",
	}, {
		name: "valid index 1",
		i:    1,
		want: "arg1",
	}, {
		name: "index out of range",
		i:    3,
		want: "",
	}, {
		name: "negative index",
		i:    -1,
		want: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.Arg(tt.i); got != tt.want {
				t.Errorf("Arg(%d) = %q, want %q", tt.i, got, tt.want)
			}
		})
	}
}

// TestIntFlag_Interface 验证IntFlag实现了Flag接口
func TestIntFlag_Interface(t *testing.T) {
	var f Flag = &IntFlag{}
	_ = f // 若编译通过，则说明实现了接口
}

// TestIntFlag_Methods 测试IntFlag的各种方法
func TestIntFlag_Methods(t *testing.T) {
	defValue := 100
	f := &IntFlag{
		name:      "intflag",
		shortName: "i",
		usage:     "integer flag test",
		defValue:  defValue,
	}
	if f.Name() != "intflag" {
		t.Errorf("IntFlag.Name() = %q, want %q", f.Name(), "intflag")
	}
	if f.ShortName() != "i" {
		t.Errorf("IntFlag.ShortName() = %q, want %q", f.ShortName(), "i")
	}
	if f.Usage() != "integer flag test" {
		t.Errorf("IntFlag.Usage() = %q, want %q", f.Usage(), "integer flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("IntFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeInt {
		t.Errorf("IntFlag.Type() = %v, want %v", f.Type(), FlagTypeInt)
	}
}

// TestStringFlag_Interface 验证StringFlag实现了Flag接口
func TestStringFlag_Interface(t *testing.T) {
	var f Flag = &StringFlag{}
	_ = f
}

// TestStringFlag_Methods 测试StringFlag的各种方法
func TestStringFlag_Methods(t *testing.T) {
	defValue := "default string"
	f := &StringFlag{
		name:      "strflag",
		shortName: "s",
		usage:     "string flag test",
		defValue:  defValue,
	}
	if f.Name() != "strflag" {
		t.Errorf("StringFlag.Name() = %q, want %q", f.Name(), "strflag")
	}
	if f.ShortName() != "s" {
		t.Errorf("StringFlag.ShortName() = %q, want %q", f.ShortName(), "s")
	}
	if f.Usage() != "string flag test" {
		t.Errorf("StringFlag.Usage() = %q, want %q", f.Usage(), "string flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("StringFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeString {
		t.Errorf("StringFlag.Type() = %v, want %v", f.Type(), FlagTypeString)
	}
}

// TestBoolFlag_Interface 验证BoolFlag实现了Flag接口
func TestBoolFlag_Interface(t *testing.T) {
	var f Flag = &BoolFlag{}
	_ = f
}

// TestBoolFlag_Methods 测试BoolFlag的各种方法
func TestBoolFlag_Methods(t *testing.T) {
	defValue := true
	f := &BoolFlag{
		name:      "boolflag",
		shortName: "b",
		usage:     "bool flag test",
		defValue:  defValue,
	}
	if f.Name() != "boolflag" {
		t.Errorf("BoolFlag.Name() = %q, want %q", f.Name(), "boolflag")
	}
	if f.ShortName() != "b" {
		t.Errorf("BoolFlag.ShortName() = %q, want %q", f.ShortName(), "b")
	}
	if f.Usage() != "bool flag test" {
		t.Errorf("BoolFlag.Usage() = %q, want %q", f.Usage(), "bool flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("BoolFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeBool {
		t.Errorf("BoolFlag.Type() = %v, want %v", f.Type(), FlagTypeBool)
	}
}

// TestFloatFlag_Interface 验证FloatFlag实现了Flag接口
func TestFloatFlag_Interface(t *testing.T) {
	var f Flag = &FloatFlag{}
	_ = f
}

// TestFloatFlag_Methods 测试FloatFlag的各种方法
func TestFloatFlag_Methods(t *testing.T) {
	defValue := 3.14
	f := &FloatFlag{
		name:      "floatflag",
		shortName: "f",
		usage:     "float flag test",
		defValue:  defValue,
	}
	if f.Name() != "floatflag" {
		t.Errorf("FloatFlag.Name() = %q, want %q", f.Name(), "floatflag")
	}
	if f.ShortName() != "f" {
		t.Errorf("FloatFlag.ShortName() = %q, want %q", f.ShortName(), "f")
	}
	if f.Usage() != "float flag test" {
		t.Errorf("FloatFlag.Usage() = %q, want %q", f.Usage(), "float flag test")
	}
	if f.GetDefault() != defValue {
		t.Errorf("FloatFlag.GetDefault() = %v, want %v", f.GetDefault(), defValue)
	}
	if f.Type() != FlagTypeFloat {
		t.Errorf("FloatFlag.Type() = %v, want %v", f.Type(), FlagTypeFloat)
	}
}
