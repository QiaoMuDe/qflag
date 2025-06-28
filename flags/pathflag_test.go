package flags

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPathFlag_BasicNormalization 测试路径规范化功能
func TestPathFlag_BasicNormalization(t *testing.T) {
	flag := &PathFlag{}
	if err := flag.Init("path", "p", ".", "test path flag"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 获取当前工作目录作为预期结果
	expected, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// 验证初始值会被规范化
	initialValue := flag.Get()
	if !strings.HasPrefix(initialValue, expected) {
		t.Errorf("Expected initial value starting with %q, got %q", expected, initialValue)
	}

	// 测试相对路径转换
	if setErr := flag.Set("../"); setErr != nil {
		t.Fatalf("Set failed: %v", setErr)
	}

	absPath, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	if flag.Get() != absPath {
		t.Errorf("Expected %q, got %q", absPath, flag.Get())
	}
}

// TestPathFlag_ExistingPath 测试存在的路径验证
func TestPathFlag_ExistingPath(t *testing.T) {
	flag := &PathFlag{}
	if err := flag.Init("path", "p", ".", "test path flag"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 创建临时目录
	tempDir := t.TempDir()

	// 测试存在的目录
	if err := flag.Set(tempDir); err != nil {
		t.Errorf("Expected valid path %q, got error: %v", tempDir, err)
	}

	// 创建临时文件
	f, err := os.CreateTemp(tempDir, "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	f.Close()

	// 测试存在的文件
	if err := flag.Set(f.Name()); err != nil {
		t.Errorf("Expected valid path %q, got error: %v", f.Name(), err)
	}
}

// TestPathFlag_NonExistingPath 测试不存在的路径验证
func TestPathFlag_NonExistingPath(t *testing.T) {
	flag := &PathFlag{}
	if err := flag.Init("path", "p", ".", "test path flag"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试不存在的路径
	invalidPath := filepath.Join(t.TempDir(), "nonexistentfile.txt")
	if err := flag.Set(invalidPath); err == nil {
		t.Errorf("Expected error for non-existent path %q, got nil", invalidPath)
	} else if !strings.Contains(err.Error(), "path does not exist") {
		t.Errorf("Expected 'path does not exist' error, got: %v", err)
	}
}

// TestPathFlag_EmptyPath 测试空路径处理
func TestPathFlag_EmptyPath(t *testing.T) {
	flag := &PathFlag{}
	if err := flag.Init("path", "p", ".", "test path flag"); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试空字符串输入
	if err := flag.Set(""); err == nil {
		t.Error("Expected error for empty path, got nil")
	} else if !strings.Contains(err.Error(), "path cannot be empty") {
		t.Errorf("Expected 'path cannot be empty' error, got: %v", err)
	}
}
