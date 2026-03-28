package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestAddRequiredGroup 测试添加必需组
func TestAddRequiredGroup(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host address", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port number", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("username", "U", "Username", "")); err != nil {
		t.Fatalf("Failed to add username flag: %v", err)
	}

	// 添加必需组
	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}

	// 验证必需组已添加
	groups := cmd.RequiredGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 required group, got %d", len(groups))
	}

	group := groups[0]
	if group.Name != "connection" {
		t.Errorf("Expected group name 'connection', got '%s'", group.Name)
	}

	if len(group.Flags) != 2 {
		t.Fatalf("Expected 2 flags in group, got %d", len(group.Flags))
	}

	if group.Flags[0] != "host" || group.Flags[1] != "port" {
		t.Errorf("Expected flags ['host', 'port'], got %v", group.Flags)
	}

	// 验证 Conditional 字段
	if group.Conditional != false {
		t.Errorf("Expected Conditional to be false, got %v", group.Conditional)
	}

	// 测试添加重复的必需组
	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err == nil {
		t.Error("Expected error when adding duplicate required group")
	}

	// 验证必需组在解析时生效
	err := cmd.Parse([]string{"--host", "localhost"})
	if err == nil {
		t.Error("Expected error when not all required flags are set")
	}

	err = cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
	if err != nil {
		t.Errorf("Expected no error when all required flags are set, got: %v", err)
	}
}

// TestAddRequiredGroupWithEmptyFlags 测试添加空标志列表的必需组
func TestAddRequiredGroupWithEmptyFlags(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加空标志列表的必需组
	if err := cmd.AddRequiredGroup("empty", []string{}, false); err == nil {
		t.Error("Expected error when adding required group with empty flags")
	}
}

// TestAddRequiredGroupWithNonExistentFlag 测试添加包含不存在标志的必需组
func TestAddRequiredGroupWithNonExistentFlag(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加一个标志
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host address", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}

	// 添加包含不存在标志的必需组
	if err := cmd.AddRequiredGroup("connection", []string{"host", "nonexistent"}, false); err == nil {
		t.Error("Expected error when adding required group with non-existent flag")
	}
}

// TestGetRequiredGroups 测试获取必需组列表
func TestGetRequiredGroups(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加多个必需组
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("username", "U", "Username", "")); err != nil {
		t.Fatalf("Failed to add username flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("password", "W", "Password", "")); err != nil {
		t.Fatalf("Failed to add password flag: %v", err)
	}

	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}
	if err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, true); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}

	// 获取必需组列表
	groups := cmd.RequiredGroups()
	if len(groups) != 2 {
		t.Fatalf("Expected 2 required groups, got %d", len(groups))
	}

	// 修改返回的列表, 不应影响原始数据
	groups[0].Name = "modified"
	modifiedGroups := cmd.RequiredGroups()
	if modifiedGroups[0].Name == "modified" {
		t.Error("Modifying returned groups should not affect original data")
	}

	// 验证必需组在解析时生效
	err := cmd.Parse([]string{"--host", "localhost"})
	if err == nil {
		t.Error("Expected error when not all required flags are set")
	}

	err = cmd.Parse([]string{"--host", "localhost", "--port", "8080", "--username", "admin", "--password", "secret"})
	if err != nil {
		t.Errorf("Expected no error when all required flags are set, got: %v", err)
	}
}

// TestRemoveRequiredGroup 测试移除必需组
func TestRemoveRequiredGroup(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "p", "Port", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}

	// 添加必需组
	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}

	// 验证必需组已添加
	groups := cmd.RequiredGroups()
	if len(groups) != 1 {
		t.Fatalf("Expected 1 required group, got %d", len(groups))
	}

	// 移除必需组
	if err := cmd.RemoveRequiredGroup("connection"); err != nil {
		t.Fatalf("Failed to remove required group: %v", err)
	}

	// 验证必需组已移除
	groups = cmd.RequiredGroups()
	if len(groups) != 0 {
		t.Fatalf("Expected 0 required groups, got %d", len(groups))
	}

	// 测试移除不存在的必需组
	if err := cmd.RemoveRequiredGroup("nonexistent"); err == nil {
		t.Error("Expected error when removing non-existing group")
	}
}

// TestGetRequiredGroup 测试获取单个必需组
func TestGetRequiredGroup(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "p", "Port", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}

	// 添加必需组
	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, true); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}

	// 获取存在的必需组
	group, found := cmd.GetRequiredGroup("connection")
	if !found {
		t.Fatal("Expected to find required group 'connection'")
	}

	if group.Name != "connection" {
		t.Errorf("Expected group name 'connection', got '%s'", group.Name)
	}

	if len(group.Flags) != 2 {
		t.Fatalf("Expected 2 flags in group, got %d", len(group.Flags))
	}

	// 验证 Conditional 字段
	if group.Conditional != true {
		t.Errorf("Expected Conditional to be true, got %v", group.Conditional)
	}

	// 获取不存在的必需组
	_, found = cmd.GetRequiredGroup("nonexistent")
	if found {
		t.Error("Expected not to find non-existent required group")
	}

	// 验证必需组在解析时生效
	err := cmd.Parse([]string{"--host", "localhost"})
	if err == nil {
		t.Error("Expected error when not all required flags are set")
	}

	err = cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
	if err != nil {
		t.Errorf("Expected no error when all required flags are set, got: %v", err)
	}
}

// TestRequiredGroupValidation 测试必需组验证
func TestRequiredGroupValidation(t *testing.T) {
	// 测试1: 设置所有必需标志, 应该成功
	func() {
		cmd := NewCmd("test1", "t1", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
		if err != nil {
			t.Errorf("Expected no error when all required flags are set, got: %v", err)
		}
	}()

	// 测试2: 只设置部分必需标志, 应该失败
	func() {
		cmd := NewCmd("test2", "t2", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost"})
		if err == nil {
			t.Error("Expected error when not all required flags are set")
		}
	}()

	// 测试3: 不设置任何必需标志, 应该失败
	func() {
		cmd := NewCmd("test3", "t3", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{})
		if err == nil {
			t.Error("Expected error when no required flags are set")
		}
	}()

	// 测试4: 多个必需组, 全部满足, 应该成功
	func() {
		cmd := NewCmd("test4", "t4", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("username", "U", "Username", "")); err != nil {
			t.Fatalf("Failed to add username flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("password", "W", "Password", "")); err != nil {
			t.Fatalf("Failed to add password flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}
		if err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080", "--username", "admin", "--password", "secret"})
		if err != nil {
			t.Errorf("Expected no error when all required groups are satisfied, got: %v", err)
		}
	}()

	// 测试5: 多个必需组, 部分不满足, 应该失败
	func() {
		cmd := NewCmd("test5", "t5", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("username", "U", "Username", "")); err != nil {
			t.Fatalf("Failed to add username flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("password", "W", "Password", "")); err != nil {
			t.Fatalf("Failed to add password flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}
		if err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
		if err == nil {
			t.Error("Expected error when not all required groups are satisfied")
		}
	}()
}

// TestRequiredGroupWithMutexGroup 测试必需组和互斥组组合使用
func TestRequiredGroupWithMutexGroup(t *testing.T) {
	// 测试1: 必需组和互斥组同时使用, 应该成功
	func() {
		cmd := NewCmd("test1", "t1", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("format", "F", "Format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "O", "Output", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}

		// 添加必需组
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		// 添加互斥组
		if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
			t.Fatalf("Failed to add mutex group: %v", err)
		}

		// 满足必需组, 满足互斥组
		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080", "--format", "json"})
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	}()

	// 测试2: 必需组不满足, 应该失败
	func() {
		cmd := NewCmd("test2", "t2", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("format", "F", "Format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}

		// 添加必需组
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--format", "json"})
		if err == nil {
			t.Error("Expected error when required group is not satisfied")
		}
	}()

	// 测试3: 互斥组不满足, 应该失败
	func() {
		cmd := NewCmd("test3", "t3", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("format", "F", "Format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "O", "Output", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}

		// 添加必需组
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		// 添加互斥组
		if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
			t.Fatalf("Failed to add mutex group: %v", err)
		}

		// 满足必需组, 但违反互斥组
		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080", "--format", "json", "--output", "result.txt"})
		if err == nil {
			t.Error("Expected error when mutex group is violated")
		}
	}()
}

// TestRequiredGroupConcurrency 测试必需组的并发安全性
func TestRequiredGroupConcurrency(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "p", "Port", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("username", "u", "Username", "")); err != nil {
		t.Fatalf("Failed to add username flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("password", "P", "Password", "")); err != nil {
		t.Fatalf("Failed to add password flag: %v", err)
	}

	// 并发添加必需组
	done := make(chan bool, 2)

	go func() {
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Errorf("Failed to add required group: %v", err)
		}
		done <- true
	}()

	go func() {
		if err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, true); err != nil {
			t.Errorf("Failed to add required group: %v", err)
		}
		done <- true
	}()

	// 等待两个goroutine完成
	<-done
	<-done

	// 验证两个必需组都已添加
	groups := cmd.RequiredGroups()
	if len(groups) != 2 {
		t.Fatalf("Expected 2 required groups, got %d", len(groups))
	}

	// 验证必需组在解析时生效
	err := cmd.Parse([]string{"--host", "localhost"})
	if err == nil {
		t.Error("Expected error when not all required flags are set")
	}

	err = cmd.Parse([]string{"--host", "localhost", "--port", "8080", "--username", "admin", "--password", "secret"})
	if err != nil {
		t.Errorf("Expected no error when all required flags are set, got: %v", err)
	}
}

// TestRequiredGroupWithCmdOpts 测试通过 CmdOpts 添加必需组
func TestRequiredGroupWithCmdOpts(t *testing.T) {
	// 测试1: 满足必需组
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)

		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}

		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
		if err != nil {
			t.Errorf("Expected no error when all required flags are set, got: %v", err)
		}
	}()

	// 测试2: 不满足必需组
	func() {
		cmd := NewCmd("test2", "t2", types.ContinueOnError)

		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}

		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost"})
		if err == nil {
			t.Error("Expected error when not all required flags are set")
		}
	}()
}

// TestConfigReturnsCopy 测试 Config() 返回副本
func TestConfigReturnsCopy(t *testing.T) {
	// 创建命令
	cmd := NewCmd("test", "t", types.ContinueOnError)

	// 添加标志
	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}

	// 添加必需组
	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}

	// 获取配置
	config := cmd.Config()

	// 修改配置中的必需组
	config.RequiredGroups[0].Name = "modified"

	// 获取配置 again
	newConfig := cmd.Config()

	// 验证修改不影响原始数据
	if newConfig.RequiredGroups[0].Name == "modified" {
		t.Error("Modifying config copy should not affect original data")
	}

	if newConfig.RequiredGroups[0].Name != "connection" {
		t.Errorf("Expected group name 'connection', got '%s'", newConfig.RequiredGroups[0].Name)
	}
}

// TestConditionalRequiredGroup 测试条件性必需组
func TestConditionalRequiredGroup(t *testing.T) {
	// 测试1: 普通必需组，应该按原逻辑工作
	func() {
		cmd := NewCmd("test1", "t1", types.ContinueOnError)
		cmd.String("username", "u", "Username", "")
		cmd.String("password", "pw", "Password", "")
		err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false) // 普通必需组
		if err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err = cmd.Parse([]string{})
		if err == nil {
			t.Error("Expected error when no flags are set in normal required group")
		}
	}()

	// 测试2: 条件性必需组，不使用任何标志，应该成功
	func() {
		cmd := NewCmd("test2", "t2", types.ContinueOnError)
		cmd.String("host", "h2", "Host", "")
		cmd.String("port", "pt", "Port", "")
		err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
		if err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err = cmd.Parse([]string{})
		if err != nil {
			t.Errorf("Expected no error when no flags are used in conditional group, got: %v", err)
		}
	}()

	// 测试3: 条件性必需组，使用部分标志，应该失败
	func() {
		cmd := NewCmd("test3", "t3", types.ContinueOnError)
		cmd.String("host", "h3", "Host", "")
		cmd.String("port", "pt3", "Port", "")
		err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
		if err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err = cmd.Parse([]string{"--host", "localhost"})
		if err == nil {
			t.Error("Expected error when only some flags in conditional group are used")
		}
	}()

	// 测试4: 条件性必需组，使用所有标志，应该成功
	func() {
		cmd := NewCmd("test4", "t4", types.ContinueOnError)
		cmd.String("host", "h4", "Host", "")
		cmd.String("port", "pt4", "Port", "")
		err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
		if err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err = cmd.Parse([]string{"--host", "localhost", "--port", "5432"})
		if err != nil {
			t.Errorf("Expected no error when all flags in conditional group are used, got: %v", err)
		}
	}()

	// 测试5: 普通必需组和条件性必需组混合使用
	func() {
		// 测试5.1: 不设置任何标志，应该失败（因为普通必需组）
		func() {
			cmd := NewCmd("test5_1", "t5_1", types.ContinueOnError)
			cmd.String("username", "u", "Username", "")
			cmd.String("password", "pw", "Password", "")
			cmd.String("host", "h5", "Host", "")
			cmd.String("port", "pt5", "Port", "")
			err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false) // 普通必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{})
			if err == nil {
				t.Error("Expected error when no flags are set in normal required group")
			}
		}()

		// 测试5.2: 只设置普通必需组，应该成功
		func() {
			cmd := NewCmd("test5_2", "t5_2", types.ContinueOnError)
			cmd.String("username", "u", "Username", "")
			cmd.String("password", "pw", "Password", "")
			cmd.String("host", "h5", "Host", "")
			cmd.String("port", "pt5", "Port", "")
			err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false) // 普通必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--username", "admin", "--password", "secret"})
			if err != nil {
				t.Errorf("Expected no error when normal required group is satisfied, got: %v", err)
			}
		}()

		// 测试5.3: 只设置条件性必需组，应该失败
		func() {
			cmd := NewCmd("test5_3", "t5_3", types.ContinueOnError)
			cmd.String("username", "u", "Username", "")
			cmd.String("password", "pw", "Password", "")
			cmd.String("host", "h5", "Host", "")
			cmd.String("port", "pt5", "Port", "")
			err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false) // 普通必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--host", "localhost"})
			if err == nil {
				t.Error("Expected error when only some flags in conditional group are used")
			}
		}()

		// 测试5.4: 设置所有标志，应该成功
		func() {
			cmd := NewCmd("test5_4", "t5_4", types.ContinueOnError)
			cmd.String("username", "u", "Username", "")
			cmd.String("password", "pw", "Password", "")
			cmd.String("host", "h5", "Host", "")
			cmd.String("port", "pt5", "Port", "")
			err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false) // 普通必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--username", "admin", "--password", "secret", "--host", "localhost", "--port", "5432"})
			if err != nil {
				t.Errorf("Expected no error when all flags are set, got: %v", err)
			}
		}()
	}()

	// 测试6: 多个条件性必需组同时使用
	func() {
		// 测试6.1: 不设置任何标志，应该成功
		func() {
			cmd := NewCmd("test6_1", "t6_1", types.ContinueOnError)
			cmd.String("host", "h6", "Host", "")
			cmd.String("port", "pt6", "Port", "")
			cmd.String("url", "u", "URL", "")
			cmd.String("method", "m", "Method", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("http_request", []string{"url", "method"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{})
			if err != nil {
				t.Errorf("Expected no error when no flags are used in conditional groups, got: %v", err)
			}
		}()

		// 测试6.2: 只设置第一个条件性必需组，应该失败
		func() {
			cmd := NewCmd("test6_2", "t6_2", types.ContinueOnError)
			cmd.String("host", "h6", "Host", "")
			cmd.String("port", "pt6", "Port", "")
			cmd.String("url", "u", "URL", "")
			cmd.String("method", "m", "Method", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("http_request", []string{"url", "method"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--host", "localhost"})
			if err == nil {
				t.Error("Expected error when only some flags in first conditional group are used")
			}
		}()

		// 测试6.3: 只设置第二个条件性必需组，应该失败
		func() {
			cmd := NewCmd("test6_3", "t6_3", types.ContinueOnError)
			cmd.String("host", "h6", "Host", "")
			cmd.String("port", "pt6", "Port", "")
			cmd.String("url", "u", "URL", "")
			cmd.String("method", "m", "Method", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("http_request", []string{"url", "method"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--url", "http://example.com"})
			if err == nil {
				t.Error("Expected error when only some flags in second conditional group are used")
			}
		}()

		// 测试6.4: 设置第一个条件性必需组的所有标志，应该成功
		func() {
			cmd := NewCmd("test6_4", "t6_4", types.ContinueOnError)
			cmd.String("host", "h6", "Host", "")
			cmd.String("port", "pt6", "Port", "")
			cmd.String("url", "u", "URL", "")
			cmd.String("method", "m", "Method", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("http_request", []string{"url", "method"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--host", "localhost", "--port", "5432"})
			if err != nil {
				t.Errorf("Expected no error when all flags in first conditional group are used, got: %v", err)
			}
		}()

		// 测试6.5: 设置第二个条件性必需组的所有标志，应该成功
		func() {
			cmd := NewCmd("test6_5", "t6_5", types.ContinueOnError)
			cmd.String("host", "h6", "Host", "")
			cmd.String("port", "pt6", "Port", "")
			cmd.String("url", "u", "URL", "")
			cmd.String("method", "m", "Method", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("http_request", []string{"url", "method"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--url", "http://example.com", "--method", "GET"})
			if err != nil {
				t.Errorf("Expected no error when all flags in second conditional group are used, got: %v", err)
			}
		}()

		// 测试6.6: 设置所有标志，应该成功
		func() {
			cmd := NewCmd("test6_6", "t6_6", types.ContinueOnError)
			cmd.String("host", "h6", "Host", "")
			cmd.String("port", "pt6", "Port", "")
			cmd.String("url", "u", "URL", "")
			cmd.String("method", "m", "Method", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddRequiredGroup("http_request", []string{"url", "method"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.Parse([]string{"--host", "localhost", "--port", "5432", "--url", "http://example.com", "--method", "GET"})
			if err != nil {
				t.Errorf("Expected no error when all flags are set, got: %v", err)
			}
		}()
	}()

	// 测试7: 条件性必需组与互斥组组合使用
	func() {
		// 测试7.1: 不设置任何标志，应该成功
		func() {
			cmd := NewCmd("test7_1", "t7_1", types.ContinueOnError)
			cmd.String("host", "h7", "Host", "")
			cmd.String("port", "pt7", "Port", "")
			cmd.String("format", "f", "Format", "")
			cmd.String("output", "o", "Output", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddMutexGroup("output_format", []string{"format", "output"}, true) // 互斥组
			if err != nil {
				t.Fatalf("Failed to add mutex group: %v", err)
			}

			err = cmd.Parse([]string{})
			if err != nil {
				t.Errorf("Expected no error when no flags are used, got: %v", err)
			}
		}()

		// 测试7.2: 只设置互斥组，应该成功
		func() {
			cmd := NewCmd("test7_2", "t7_2", types.ContinueOnError)
			cmd.String("host", "h7", "Host", "")
			cmd.String("port", "pt7", "Port", "")
			cmd.String("format", "f", "Format", "")
			cmd.String("output", "o", "Output", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddMutexGroup("output_format", []string{"format", "output"}, true) // 互斥组
			if err != nil {
				t.Fatalf("Failed to add mutex group: %v", err)
			}

			err = cmd.Parse([]string{"--format", "json"})
			if err != nil {
				t.Errorf("Expected no error when only mutex group is used, got: %v", err)
			}
		}()

		// 测试7.3: 设置条件性必需组，应该失败
		func() {
			cmd := NewCmd("test7_3", "t7_3", types.ContinueOnError)
			cmd.String("host", "h7", "Host", "")
			cmd.String("port", "pt7", "Port", "")
			cmd.String("format", "f", "Format", "")
			cmd.String("output", "o", "Output", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddMutexGroup("output_format", []string{"format", "output"}, true) // 互斥组
			if err != nil {
				t.Fatalf("Failed to add mutex group: %v", err)
			}

			err = cmd.Parse([]string{"--host", "localhost"})
			if err == nil {
				t.Error("Expected error when only some flags in conditional group are used")
			}
		}()

		// 测试7.4: 设置条件性必需组的所有标志，应该成功
		func() {
			cmd := NewCmd("test7_4", "t7_4", types.ContinueOnError)
			cmd.String("host", "h7", "Host", "")
			cmd.String("port", "pt7", "Port", "")
			cmd.String("format", "f", "Format", "")
			cmd.String("output", "o", "Output", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddMutexGroup("output_format", []string{"format", "output"}, true) // 互斥组
			if err != nil {
				t.Fatalf("Failed to add mutex group: %v", err)
			}

			err = cmd.Parse([]string{"--host", "localhost", "--port", "5432"})
			if err != nil {
				t.Errorf("Expected no error when all flags in conditional group are used, got: %v", err)
			}
		}()

		// 测试7.5: 违反互斥组，应该失败
		func() {
			cmd := NewCmd("test7_5", "t7_5", types.ContinueOnError)
			cmd.String("host", "h7", "Host", "")
			cmd.String("port", "pt7", "Port", "")
			cmd.String("format", "f", "Format", "")
			cmd.String("output", "o", "Output", "")
			err := cmd.AddRequiredGroup("database", []string{"host", "port"}, true) // 条件性必需组
			if err != nil {
				t.Fatalf("Failed to add required group: %v", err)
			}

			err = cmd.AddMutexGroup("output_format", []string{"format", "output"}, true) // 互斥组
			if err != nil {
				t.Fatalf("Failed to add mutex group: %v", err)
			}

			err = cmd.Parse([]string{"--format", "json", "--output", "file.txt"})
			if err == nil {
				t.Error("Expected error when mutex group is violated")
			}
		}()
	}()

	// 测试8: CmdOpts支持条件性必需组
	func() {
		// 测试8.1: 不设置任何标志，应该成功
		func() {
			cmd2 := NewCmd("test8_1", "t8_1", types.ContinueOnError)
			cmd2.String("host", "ht", "Host", "")
			cmd2.String("port", "pt2", "Port", "")

			opts := NewCmdOpts()
			opts.RequiredGroups = []types.RequiredGroup{
				{Name: "database", Flags: []string{"host", "port"}, Conditional: true},
			}

			err := cmd2.ApplyOpts(opts)
			if err != nil {
				t.Fatalf("Failed to apply opts: %v", err)
			}

			// 不设置任何标志，应该成功
			err = cmd2.Parse([]string{})
			if err != nil {
				t.Errorf("Expected no error when no flags are used in conditional group, got: %v", err)
			}
		}()

		// 测试8.2: 只设置部分标志，应该失败
		func() {
			cmd2 := NewCmd("test8_2", "t8_2", types.ContinueOnError)
			cmd2.String("host", "ht", "Host", "")
			cmd2.String("port", "pt2", "Port", "")

			opts := NewCmdOpts()
			opts.RequiredGroups = []types.RequiredGroup{
				{Name: "database", Flags: []string{"host", "port"}, Conditional: true},
			}

			err := cmd2.ApplyOpts(opts)
			if err != nil {
				t.Fatalf("Failed to apply opts: %v", err)
			}

			// 只设置部分标志，应该失败
			err = cmd2.Parse([]string{"--host", "localhost"})
			if err == nil {
				t.Error("Expected error when only some flags in conditional group are used")
			}
		}()
	}()
}
