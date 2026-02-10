package registry

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestNewRegistry 测试泛型注册表创建
func TestNewRegistry(t *testing.T) {
	reg := NewRegistry[string]()

	// 检查注册表是否正确初始化
	if len(reg.items) != 0 {
		t.Error("NewRegistry() should initialize empty items map")
	}

	if len(reg.nameIndex) != 0 {
		t.Error("NewRegistry() should initialize empty nameIndex map")
	}
}

// TestRegistry_Register 测试注册功能
func TestRegistry_Register(t *testing.T) {
	reg := NewRegistry[string]()

	// 测试正常注册
	err := reg.Register("value", "test", "")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试重复注册
	err = reg.Register("value2", "test", "")
	if err != types.ErrFlagAlreadyExists {
		t.Errorf("Register() error = %v, want %v", err, types.ErrFlagAlreadyExists)
	}

	// 测试空名称注册
	err = reg.Register("value", "", "")
	if err == nil {
		t.Error("Register() with empty names should return error")
	}
}

// TestRegistry_Unregister 测试注销功能
func TestRegistry_Unregister(t *testing.T) {
	reg := NewRegistry[string]()

	// 注册一个项
	if err := reg.Register("value", "test", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试正常注销
	err := reg.Unregister("test")
	if err != nil {
		t.Errorf("Unregister() error = %v", err)
	}

	// 验证项已移除
	_, exists := reg.Get("test")
	if exists {
		t.Error("Unregister() should remove the item")
	}

	// 测试注销不存在的项
	err = reg.Unregister("nonexistent")
	if err != types.ErrFlagNotFound {
		t.Errorf("Unregister() error = %v, want %v", err, types.ErrFlagNotFound)
	}
}

// TestRegistry_Get 测试获取功能
func TestRegistry_Get(t *testing.T) {
	reg := NewRegistry[string]()

	// 注册一个项
	if err := reg.Register("value", "test", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试获取存在的项
	value, exists := reg.Get("test")
	if !exists {
		t.Error("Get() should find the item")
	}
	if value != "value" {
		t.Errorf("Get() value = %v, want %v", value, "value")
	}

	// 测试获取不存在的项
	_, exists = reg.Get("nonexistent")
	if exists {
		t.Error("Get() should not find nonexistent item")
	}
}

// TestRegistry_GetByShortName 测试通过短名称获取功能
func TestRegistry_GetByShortName(t *testing.T) {
	reg := NewRegistry[*mock.MockFlag]()

	// 创建一个带短名称的标志
	flag := mock.NewMockFlag("test", "t", "Test flag", types.FlagTypeBool, false)
	if err := reg.Register(flag, "test", "t"); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试通过短名称获取
	value, exists := reg.Get("t")
	if !exists {
		t.Error("GetByShortName() should find the item")
	}
	if value != flag {
		t.Error("GetByShortName() should return the correct flag")
	}

	// 测试获取不存在的短名称
	_, exists = reg.Get("x")
	if exists {
		t.Error("GetByShortName() should not find nonexistent short name")
	}
}

// TestRegistry_List 测试列出所有项
func TestRegistry_List(t *testing.T) {
	reg := NewRegistry[string]()

	// 注册多个项
	if err := reg.Register("value1", "test1", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if err := reg.Register("value2", "test2", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if err := reg.Register("value3", "test3", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 获取列表
	items := reg.List()

	if len(items) != 3 {
		t.Errorf("List() returned %d items, want 3", len(items))
	}

	// 验证所有项都在列表中
	itemSet := make(map[string]bool)
	for _, item := range items {
		itemSet[item] = true
	}

	if !itemSet["value1"] || !itemSet["value2"] || !itemSet["value3"] {
		t.Error("List() should include all registered items")
	}
}

// TestRegistry_Has 测试检查项是否存在
func TestRegistry_Has(t *testing.T) {
	reg := NewRegistry[string]()

	// 注册一个项
	if err := reg.Register("value", "test", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试存在的项
	if !reg.Has("test") {
		t.Error("Has() should return true for existing item")
	}

	// 测试不存在的项
	if reg.Has("nonexistent") {
		t.Error("Has() should return false for nonexistent item")
	}
}

// TestRegistry_Count 测试计数功能
func TestRegistry_Count(t *testing.T) {
	reg := NewRegistry[string]()

	// 初始计数应为0
	if reg.Count() != 0 {
		t.Errorf("Count() = %d, want 0", reg.Count())
	}

	// 添加项后计数应增加
	if err := reg.Register("value1", "test1", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if reg.Count() != 1 {
		t.Errorf("Count() = %d, want 1", reg.Count())
	}

	if err := reg.Register("value2", "test2", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if reg.Count() != 2 {
		t.Errorf("Count() = %d, want 2", reg.Count())
	}

	// 移除项后计数应减少
	if err := reg.Unregister("test1"); err != nil {
		t.Errorf("Unregister() error = %v", err)
	}
	if reg.Count() != 1 {
		t.Errorf("Count() = %d, want 1", reg.Count())
	}
}

// TestRegistry_Clear 测试清空功能
func TestRegistry_Clear(t *testing.T) {
	reg := NewRegistry[string]()

	// 添加多个项
	if err := reg.Register("value1", "test1", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if err := reg.Register("value2", "test2", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if err := reg.Register("value3", "test3", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 验证项已添加
	if reg.Count() != 3 {
		t.Errorf("Before Clear(), Count() = %d, want 3", reg.Count())
	}

	// 清空注册表
	reg.Clear()

	// 验证注册表已清空
	if reg.Count() != 0 {
		t.Errorf("After Clear(), Count() = %d, want 0", reg.Count())
	}

	if reg.Has("test1") || reg.Has("test2") || reg.Has("test3") {
		t.Error("Clear() should remove all items")
	}
}

// TestRegistry_Range 测试遍历功能
func TestRegistry_Range(t *testing.T) {
	reg := NewRegistry[string]()

	// 添加多个项
	if err := reg.Register("value1", "test1", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if err := reg.Register("value2", "test2", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}
	if err := reg.Register("value3", "test3", ""); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 遍历所有项
	visited := make(map[string]bool)
	count := 0

	reg.Range(func(name string, item string) bool {
		visited[name] = true
		count++

		// 验证项的值
		if name == "test1" && item != "value1" {
			t.Errorf("Range() item mismatch for test1: got %v, want %v", item, "value1")
		}
		if name == "test2" && item != "value2" {
			t.Errorf("Range() item mismatch for test2: got %v, want %v", item, "value2")
		}
		if name == "test3" && item != "value3" {
			t.Errorf("Range() item mismatch for test3: got %v, want %v", item, "value3")
		}

		return true // 继续遍历
	})

	// 验证所有项都被访问
	if count != 3 {
		t.Errorf("Range() visited %d items, want 3", count)
	}

	if !visited["test1"] || !visited["test2"] || !visited["test3"] {
		t.Error("Range() should visit all items")
	}

	// 测试提前终止遍历
	earlyCount := 0
	reg.Range(func(name string, item string) bool {
		earlyCount++
		return false // 提前终止
	})

	if earlyCount != 1 {
		t.Errorf("Range() with early termination visited %d items, want 1", earlyCount)
	}
}

// TestNewFlagRegistry 测试标志注册表创建
func TestNewFlagRegistry(t *testing.T) {
	flagReg := NewFlagRegistry()

	if flagReg == nil {
		t.Error("NewFlagRegistry() should not return nil")
	}

	// 验证初始状态
	if flagReg.Count() != 0 {
		t.Errorf("NewFlagRegistry() Count() = %d, want 0", flagReg.Count())
	}
}

// TestFlagRegistryImpl_Register 测试标志注册
func TestFlagRegistryImpl_Register(t *testing.T) {
	flagReg := NewFlagRegistry()

	// 创建测试标志
	testFlag := flag.NewBoolFlag("test", "t", "Test flag", false)

	// 测试正常注册
	err := flagReg.Register(testFlag)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 验证标志已注册
	_, exists := flagReg.Get("test")
	if !exists {
		t.Error("Register() should register the flag")
	}

	// 测试注册nil标志
	err = flagReg.Register(nil)
	if err == nil {
		t.Error("Register() with nil flag should return error")
	}

	// 测试注册空名称标志
	emptyFlag := mock.NewMockFlag("", "", "Empty name flag", types.FlagTypeBool, false)
	err = flagReg.Register(emptyFlag)
	if err == nil {
		t.Error("Register() with empty name and short name should return error")
	}

	// 测试重复注册
	err = flagReg.Register(testFlag)
	if err != types.ErrFlagAlreadyExists {
		t.Errorf("Register() duplicate error = %v, want %v", err, types.ErrFlagAlreadyExists)
	}

	// 测试只有短名称的标志
	shortOnlyFlag := mock.NewMockBoolFlag("", "s", "Short only flag", false)
	if err := flagReg.Register(shortOnlyFlag); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 验证短名称可以获取
	flag, exists := flagReg.Get("s")
	if !exists {
		t.Error("Register() should register flag with short name")
	}
	if flag != shortOnlyFlag {
		t.Error("Register() should return correct flag for short name")
	}
}

// TestFlagRegistryImpl_GetByShortName 测试通过短名称获取标志
func TestFlagRegistryImpl_GetByShortName(t *testing.T) {
	flagReg := NewFlagRegistry()

	// 创建测试标志
	testFlag := flag.NewBoolFlag("test", "t", "Test flag", false)
	if err := flagReg.Register(testFlag); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试通过短名称获取
	flag, exists := flagReg.Get("t")
	if !exists {
		t.Error("GetByShortName() should find the flag")
	}
	if flag != testFlag {
		t.Error("GetByShortName() should return the correct flag")
	}

	// 测试获取不存在的短名称
	_, exists = flagReg.Get("x")
	if exists {
		t.Error("GetByShortName() should not find nonexistent short name")
	}
}

// TestNewCmdRegistry 测试命令注册表创建
func TestNewCmdRegistry(t *testing.T) {
	cmdReg := NewCmdRegistry()

	if cmdReg == nil {
		t.Error("NewCmdRegistry() should not return nil")
	}

	// 验证初始状态
	if cmdReg.Count() != 0 {
		t.Errorf("NewCmdRegistry() Count() = %d, want 0", cmdReg.Count())
	}
}

// TestCmdRegistryImpl_Register 测试命令注册
func TestCmdRegistryImpl_Register(t *testing.T) {
	cmdReg := NewCmdRegistry()

	// 创建测试命令
	testCmd := mock.NewMockCommand("test", "t", "Test command")

	// 测试正常注册
	err := cmdReg.Register(testCmd)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 验证命令已注册
	_, exists := cmdReg.Get("test")
	if !exists {
		t.Error("Register() should register the command")
	}

	// 测试注册nil命令
	err = cmdReg.Register(nil)
	if err == nil {
		t.Error("Register() with nil command should return error")
	}

	// 测试注册空名称命令
	emptyCmd := mock.NewMockCommand("", "", "Empty name command")
	err = cmdReg.Register(emptyCmd)
	if err == nil {
		t.Error("Register() with empty name and short name should return error")
	}

	// 测试重复注册
	err = cmdReg.Register(testCmd)
	if err != types.ErrFlagAlreadyExists {
		t.Errorf("Register() duplicate error = %v, want %v", err, types.ErrFlagAlreadyExists)
	}

	// 测试只有短名称的命令
	shortOnlyCmd := mock.NewMockCommand("", "s", "Short only command")
	if err := cmdReg.Register(shortOnlyCmd); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 验证短名称可以获取
	cmd, exists := cmdReg.Get("s")
	if !exists {
		t.Error("Register() should register command with short name")
	}
	if cmd != shortOnlyCmd {
		t.Error("Register() should return correct command for short name")
	}
}

// TestCmdRegistryImpl_GetByShortName 测试通过短名称获取命令
func TestCmdRegistryImpl_GetByShortName(t *testing.T) {
	cmdReg := NewCmdRegistry()

	// 创建测试命令
	testCmd := mock.NewMockCommand("test", "t", "Test command")
	if err := cmdReg.Register(testCmd); err != nil {
		t.Errorf("Register() error = %v", err)
	}

	// 测试通过短名称获取
	cmd, exists := cmdReg.Get("t")
	if !exists {
		t.Error("GetByShortName() should find the command")
	}
	if cmd != testCmd {
		t.Error("GetByShortName() should return the correct command")
	}

	// 测试获取不存在的短名称
	_, exists = cmdReg.Get("x")
	if exists {
		t.Error("GetByShortName() should not find nonexistent short name")
	}
}
