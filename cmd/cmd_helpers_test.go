package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// TestSetVersionAndGetVersion 测试版本设置和获取功能
func TestSetVersionAndGetVersion(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	version := "v1.0.0"
	cmd.SetVersion(version)

	if cmd.GetVersion() != version {
		t.Errorf("GetVersion() = %q, want %q", cmd.GetVersion(), version)
	}
}

// TestAddNoteAndGetNotes 测试备注添加和获取功能
func TestAddNoteAndGetNotes(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	note1 := "Test note 1"
	note2 := "Test note 2"

	cmd.AddNote(note1)
	cmd.AddNote(note2)

	notes := cmd.GetNotes()
	if len(notes) != 3 {
		t.Fatalf("GetNotes() returned %d notes, want 3", len(notes))
	}

	// 0 是默认的内置备注
	if notes[1] != note1 || notes[2] != note2 {
		t.Errorf("GetNotes() = %v, want [%q, %q]", notes, note1, note2)
	}
}

// TestLoadHelp 测试从文件加载帮助信息功能
func TestLoadHelp(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	content := "Test help content"

	// 创建临时文件
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "help.txt")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试正常加载
	if err := cmd.LoadHelp(filePath); err != nil {
		t.Errorf("LoadHelp() returned error: %v", err)
	}

	if cmd.GetHelp() != content {
		t.Errorf("GetHelp() = %q, want %q", cmd.GetHelp(), content)
	}

	// 测试文件不存在情况
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.txt")
	err := cmd.LoadHelp(nonExistentPath)
	if err == nil {
		t.Error("LoadHelp() expected error for nonexistent file, got nil")
	}
}

// TestFlagExists 测试标志存在性检测功能
func TestFlagExists(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.String("test-flag", "t", "default", "test flag")

	// 测试存在的标志
	if !cmd.FlagExists("test-flag") {
		t.Error("FlagExists() returned false for existing long flag")
	}
	if !cmd.FlagExists("t") {
		t.Error("FlagExists() returned false for existing short flag")
	}

	// 测试不存在的标志
	if cmd.FlagExists("nonexistent") {
		t.Error("FlagExists() returned true for nonexistent flag")
	}
}

// TestArgsMethods 测试参数获取相关方法
func TestArgsMethods(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	args := []string{"arg1", "arg2", "arg3"}

	// 解析参数
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	// 测试NArg
	if cmd.NArg() != 3 {
		t.Errorf("NArg() = %d, want 3", cmd.NArg())
	}

	// 测试Arg
	if cmd.Arg(0) != "arg1" || cmd.Arg(1) != "arg2" || cmd.Arg(2) != "arg3" {
		t.Errorf("Arg() returned unexpected values: %q, %q, %q", cmd.Arg(0), cmd.Arg(1), cmd.Arg(2))
	}

	// 测试越界Arg
	if cmd.Arg(3) != "" {
		t.Errorf("Arg(3) = %q, want empty string", cmd.Arg(3))
	}

	// 测试Args
	allArgs := cmd.Args()
	if len(allArgs) != 3 || allArgs[0] != "arg1" || allArgs[1] != "arg2" || allArgs[2] != "arg3" {
		t.Errorf("Args() = %v, want %v", allArgs, args)
	}
}
