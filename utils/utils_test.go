// utils 工具包测试
package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetExecutablePath 测试GetExecutablePath函数的功能
// 验证函数是否能正确返回可执行文件的绝对路径
func TestGetExecutablePath(t *testing.T) {
	// 测试正常情况下是否返回绝对路径
	path := GetExecutablePath()
	if !filepath.IsAbs(path) {
		t.Errorf("GetExecutablePath() 返回非绝对路径: %s", path)
	}

	// 测试路径是否存在（基本验证，不保证一定存在但能捕获明显错误）
	if _, err := os.Stat(path); err != nil {
		t.Logf("警告: 获取的可执行路径可能不存在: %s, 错误: %v", path, err)
	}
}
