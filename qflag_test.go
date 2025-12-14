// Package qflag 测试文件
// 展示如何使用Root全局实例注册和解析命令行参数的基本示例
package qflag

import (
	"testing"
)

// TestGlobalRootInstance 测试并展示如何使用全局Root实例
// 这是一个简洁的示例，展示了qflag包的基本使用方式
func TestGlobalRootInstance(t *testing.T) {
	// 重置Root实例，确保从干净状态开始测试
	defer func() {
		Root = nil
	}()

	// 使用全局Root实例注册标志
	name := Root.String("name", "n", "default", "用户名")
	age := Root.Int("age", "a", 18, "年龄")
	verbose := Root.Bool("verbose", "v", false, "详细模式")

	// 直接解析测试参数，避免修改全局os.Args
	testArgs := []string{
		"--name", "testuser",
		"--age", "30",
		"--verbose",
	}

	err := Root.Parse(testArgs)
	if err != nil {
		t.Fatalf("解析命令行参数失败: %v", err)
	}

	// 验证解析结果
	if name.Get() != "testuser" {
		t.Errorf("期望name为'testuser'，实际为'%s'", name.Get())
	}
	if age.Get() != 30 {
		t.Errorf("期望age为30，实际为%d", age.Get())
	}
	if verbose.Get() != true {
		t.Errorf("期望verbose为true，实际为false")
	}

	t.Logf("全局Root实例测试成功！")
	t.Logf("  用户名: %s", name.Get())
	t.Logf("  年龄: %d", age.Get())
	t.Logf("  详细模式: %v", verbose.Get())
}
