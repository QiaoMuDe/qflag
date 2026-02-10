package builtin

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

func TestBuiltinFlagManager_Basic(t *testing.T) {
	// 测试帮助标志处理器
	helpHandler := &HelpHandler{}
	if helpHandler.Type() != types.HelpFlag {
		t.Errorf("Expected HelpFlag, got %v", helpHandler.Type())
	}

	// 测试版本标志处理器
	versionHandler := &VersionHandler{}
	if versionHandler.Type() != types.VersionFlag {
		t.Errorf("Expected VersionFlag, got %v", versionHandler.Type())
	}
}

func TestHelpHandler_Basic(t *testing.T) {
	// 创建模拟命令
	c := mock.NewMockCommandBasic("test", "t", "Test command")

	// 创建帮助处理器
	handler := &HelpHandler{}

	// 测试ShouldRegister总是返回true
	if !handler.ShouldRegister(c) {
		t.Error("HelpHandler.ShouldRegister should always return true")
	}

	// 测试Type方法
	if handler.Type() != types.HelpFlag {
		t.Errorf("Expected HelpFlag, got %v", handler.Type())
	}
}

func TestVersionHandler_Basic(t *testing.T) {
	// 创建模拟命令
	c := mock.NewMockCommandBasic("test", "t", "Test command")

	// 创建版本处理器
	handler := &VersionHandler{}

	// 测试没有版本信息时不应该注册
	if handler.ShouldRegister(c) {
		t.Error("VersionHandler.ShouldRegister should return false when no version is set")
	}

	// 设置版本信息
	c.SetVersion("1.0.0")

	// 测试有版本信息时应该注册
	if !handler.ShouldRegister(c) {
		t.Error("VersionHandler.ShouldRegister should return true when version is set")
	}

	// 测试Type方法
	if handler.Type() != types.VersionFlag {
		t.Errorf("Expected VersionFlag, got %v", handler.Type())
	}
}
