package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

// TestPathVar 测试PathVar方法的功能
func TestPathVar(t *testing.T) {

	// 测试指针为nil的情况
	t.Run("nil pointer", func(t *testing.T) {
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		defer func() {
			if r := recover(); r == nil {
				t.Error("PathVar with nil pointer should panic")
			}
		}()
		cmd.PathVar(nil, "path", "p", "./test", "test path flag")
	})

	// 测试正常功能
	t.Run("normal case", func(t *testing.T) {
		defaultPath := "./default"
		cmd := NewCommand("test", "t", flag.ContinueOnError)
		var pathFlag flags.PathFlag
		cmd.PathVar(&pathFlag, "path", "p", defaultPath, "test path flag")

		// 测试默认值
		absDefault, err := filepath.Abs("./default")
		if err != nil {
			t.Fatalf("failed to get absolute path: %v", err)
		}
		if pathFlag.GetDefault() != absDefault {
			t.Errorf("default value = %v, want %v", pathFlag.Get(), absDefault)
		}

		// 测试短标志解析（使用临时目录确保路径存在）
		cmd = NewCommand("test-short", "ts", flag.ContinueOnError)
		var pathFlagShort flags.PathFlag
		tempDir := t.TempDir()
		cmd.PathVar(&pathFlagShort, "path-short", "p", "/", "test path short flag")
		if err := cmd.Parse([]string{"-p", tempDir}); err != nil {
			t.Fatalf("Parse failed: %v", err)
		}
		if pathFlagShort.Get() != tempDir {
			t.Errorf("after -p, value = %v, want %v", pathFlagShort.Get(), tempDir)
		}
	})

	// 测试路径验证（存在性检查）
	t.Run("path validation", func(t *testing.T) {
		// 创建临时目录用于测试
		tempDir := t.TempDir()
		validPath := filepath.Join(tempDir, "valid.txt")
		_ = os.WriteFile(validPath, []byte("test"), 0644)
		invalidPath := filepath.Join(tempDir, "nonexistent.txt")

		// 测试有效路径
		cmdValid := NewCommand("test-valid", "tv", flag.ContinueOnError)
		var validPathFlag flags.PathFlag
		cmdValid.PathVar(&validPathFlag, "valid-path", "v", "", "test valid path")
		if err := cmdValid.Parse([]string{"--valid-path", validPath}); err != nil {
			t.Fatalf("Parse valid path failed: %v", err)
		}

		// 测试无效路径（假设PathFlag有存在性验证）
		cmdInvalid := NewCommand("test-invalid", "ti", flag.ContinueOnError)
		var invalidPathFlag flags.PathFlag
		cmdInvalid.PathVar(&invalidPathFlag, "invalid-path", "i", "", "test invalid path")
		err := cmdInvalid.Parse([]string{"--invalid-path", invalidPath})
		if err == nil {
			t.Error("expected error for nonexistent path, got nil")
		}
	})
}
