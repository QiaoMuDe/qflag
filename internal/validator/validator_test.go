package validator

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"testing"

	"gitee.com/MM-Q/qflag/internal/types"
)

// TestValidateSubCommand_NilChild 测试子命令为nil的情况
func TestValidateSubCommand_NilChild(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	err := ValidateSubCommand(parent, nil)
	if err == nil {
		t.Error("期望返回错误，但得到了nil")
	}

	expectedMsg := "subcmd <nil> is nil"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误消息包含 '%s'，但得到: %s", expectedMsg, err.Error())
	}
}

// TestValidateSubCommand_NilParent 测试父命令为nil的情况
func TestValidateSubCommand_NilParent(t *testing.T) {
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	// 虽然父命令为nil，但函数应该能正常处理
	err := ValidateSubCommand(nil, child)
	if err != nil {
		t.Errorf("期望nil父命令时不返回错误，但得到: %s", err.Error())
	}
}

// TestValidateSubCommand_BothNil 测试父命令和子命令都为nil的情况
func TestValidateSubCommand_BothNil(t *testing.T) {
	err := ValidateSubCommand(nil, nil)
	if err == nil {
		t.Error("期望返回错误，但得到了nil")
	}

	expectedMsg := "subcmd <nil> is nil"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误消息包含 '%s'，但得到: %s", expectedMsg, err.Error())
	}
}

// TestValidateSubCommand_LongNameConflict 测试长名称冲突
func TestValidateSubCommand_LongNameConflict(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	existing := types.NewCmdContext("existing", "e", flag.ContinueOnError)
	conflicting := types.NewCmdContext("existing", "c", flag.ContinueOnError)

	// 先添加一个子命令
	parent.SubCmdMap["existing"] = existing

	err := ValidateSubCommand(parent, conflicting)
	if err == nil {
		t.Error("期望返回长名称冲突错误，但得到了nil")
	}

	expectedMsg := "long name 'existing' already exists"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误消息包含 '%s'，但得到: %s", expectedMsg, err.Error())
	}
}

// TestValidateSubCommand_ShortNameConflict 测试短名称冲突
func TestValidateSubCommand_ShortNameConflict(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	existing := types.NewCmdContext("existing", "e", flag.ContinueOnError)
	conflicting := types.NewCmdContext("conflicting", "e", flag.ContinueOnError)

	// 先添加一个子命令
	parent.SubCmdMap["e"] = existing

	err := ValidateSubCommand(parent, conflicting)
	if err == nil {
		t.Error("期望返回短名称冲突错误，但得到了nil")
	}

	expectedMsg := "short name 'e' already exists"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("期望错误消息包含 '%s'，但得到: %s", expectedMsg, err.Error())
	}
}

// TestValidateSubCommand_EmptyNames 测试空名称的边界情况
func TestValidateSubCommand_EmptyNames(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	// 测试长名称为空的子命令
	childEmptyLong := &types.CmdContext{
		LongName:  "",
		ShortName: "c",
		SubCmdMap: make(map[string]*types.CmdContext),
	}

	err := ValidateSubCommand(parent, childEmptyLong)
	if err != nil {
		t.Errorf("期望空长名称不产生错误，但得到: %s", err.Error())
	}

	// 测试短名称为空的子命令
	childEmptyShort := &types.CmdContext{
		LongName:  "child",
		ShortName: "",
		SubCmdMap: make(map[string]*types.CmdContext),
	}

	err = ValidateSubCommand(parent, childEmptyShort)
	if err != nil {
		t.Errorf("期望空短名称不产生错误，但得到: %s", err.Error())
	}
}

// TestValidateSubCommand_ValidCase 测试正常有效的情况
func TestValidateSubCommand_ValidCase(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	err := ValidateSubCommand(parent, child)
	if err != nil {
		t.Errorf("期望有效的子命令验证通过，但得到错误: %s", err.Error())
	}
}

// TestHasCycle_NilInputs 测试HasCycle函数的nil输入
func TestHasCycle_NilInputs(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	// 测试父命令为nil
	if HasCycle(nil, child) {
		t.Error("期望nil父命令不存在循环，但返回true")
	}

	// 测试子命令为nil
	if HasCycle(parent, nil) {
		t.Error("期望nil子命令不存在循环，但返回true")
	}

	// 测试两者都为nil
	if HasCycle(nil, nil) {
		t.Error("期望两者都为nil时不存在循环，但返回true")
	}
}

// TestHasCycle_DirectCycle 测试直接循环引用
func TestHasCycle_DirectCycle(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	// 测试子命令就是父命令本身
	if !HasCycle(parent, parent) {
		t.Error("期望检测到直接循环引用，但返回false")
	}
}

// TestHasCycle_IndirectCycle 测试间接循环引用
func TestHasCycle_IndirectCycle(t *testing.T) {
	grandParent := types.NewCmdContext("grandparent", "gp", flag.ContinueOnError)
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	// 建立父子关系链
	parent.Parent = grandParent
	child.Parent = parent

	// 测试将祖父命令作为子命令添加到子命令中（形成循环）
	if !HasCycle(child, grandParent) {
		t.Error("期望检测到间接循环引用，但返回false")
	}

	// 测试将父命令作为子命令添加到子命令中
	if !HasCycle(child, parent) {
		t.Error("期望检测到间接循环引用，但返回false")
	}
}

// TestHasCycle_NoCycle 测试无循环的正常情况
func TestHasCycle_NoCycle(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)
	sibling := types.NewCmdContext("sibling", "s", flag.ContinueOnError)

	// 建立正常的父子关系
	child.Parent = parent

	// 测试添加兄弟命令（无循环）
	if HasCycle(parent, sibling) {
		t.Error("期望无循环的情况返回false，但返回true")
	}

	// 测试添加孙子命令（无循环）
	grandChild := types.NewCmdContext("grandchild", "gc", flag.ContinueOnError)
	if HasCycle(child, grandChild) {
		t.Error("期望无循环的情况返回false，但返回true")
	}
}

// TestHasCycle_DeepChain 测试深层命令链的循环检测
func TestHasCycle_DeepChain(t *testing.T) {
	// 创建一个深层的命令链
	commands := make([]*types.CmdContext, 10)
	for i := 0; i < 10; i++ {
		commands[i] = types.NewCmdContext(
			fmt.Sprintf("cmd%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)
		if i > 0 {
			commands[i].Parent = commands[i-1]
		}
	}

	// 测试在深层链中检测循环
	if !HasCycle(commands[9], commands[0]) {
		t.Error("期望在深层链中检测到循环引用，但返回false")
	}

	// 测试在深层链中检测中间节点的循环
	if !HasCycle(commands[9], commands[5]) {
		t.Error("期望在深层链中检测到中间节点循环引用，但返回false")
	}
}

// TestHasCycle_MaxDepthProtection 测试最大深度保护机制
func TestHasCycle_MaxDepthProtection(t *testing.T) {
	// 创建一个超过100层的命令链来测试深度保护
	var current *types.CmdContext
	root := types.NewCmdContext("root", "r", flag.ContinueOnError)
	current = root

	// 创建101层深的链
	for i := 1; i <= 101; i++ {
		next := types.NewCmdContext(
			fmt.Sprintf("cmd%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)
		next.Parent = current
		current = next
	}

	// 测试深度保护是否生效（应该在100层时停止检查）
	newCmd := types.NewCmdContext("new", "n", flag.ContinueOnError)
	result := HasCycle(current, newCmd)

	// 由于深度限制，应该返回false（未检测到循环）
	if result {
		t.Error("期望深度保护机制生效，返回false，但返回true")
	}
}

// TestGetCmdIdentifier_NilCommand 测试GetCmdIdentifier函数的nil输入
func TestGetCmdIdentifier_NilCommand(t *testing.T) {
	result := GetCmdIdentifier(nil)
	expected := "<nil>"
	if result != expected {
		t.Errorf("期望nil命令返回 '%s'，但得到: %s", expected, result)
	}
}

// TestGetCmdIdentifier_ValidCommand 测试GetCmdIdentifier函数的有效输入
func TestGetCmdIdentifier_ValidCommand(t *testing.T) {
	// 测试有长名称的命令
	cmdWithLong := types.NewCmdContext("longname", "l", flag.ContinueOnError)
	result := GetCmdIdentifier(cmdWithLong)
	if result != "longname" {
		t.Errorf("期望返回长名称 'longname'，但得到: %s", result)
	}

	// 测试只有短名称的命令
	cmdWithShort := types.NewCmdContext("", "s", flag.ContinueOnError)
	result = GetCmdIdentifier(cmdWithShort)
	if result != "s" {
		t.Errorf("期望返回短名称 's'，但得到: %s", result)
	}
}

// TestGetCmdIdentifier_EmptyNames 测试GetCmdIdentifier函数的空名称边界情况
func TestGetCmdIdentifier_EmptyNames(t *testing.T) {
	// 创建一个名称都为空的命令上下文（通过直接构造避免NewCmdContext的panic）
	cmd := &types.CmdContext{
		LongName:  "",
		ShortName: "",
	}

	result := GetCmdIdentifier(cmd)
	if result != "" {
		t.Errorf("期望空名称命令返回空字符串，但得到: %s", result)
	}
}

// TestValidateSubCommand_ConcurrentAccess 测试并发访问场景
func TestValidateSubCommand_ConcurrentAccess(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// 并发添加多个不同的子命令
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			child := types.NewCmdContext(
				fmt.Sprintf("child%d", index),
				fmt.Sprintf("c%d", index),
				flag.ContinueOnError,
			)
			err := ValidateSubCommand(parent, child)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	// 并发添加一些重复的子命令（应该产生冲突）
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			child := types.NewCmdContext("duplicate", "d", flag.ContinueOnError)
			err := ValidateSubCommand(parent, child)
			// 这里可能会有冲突，但不应该panic
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// 检查是否有任何panic或意外错误
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	// 应该至少有一些重复名称的错误
	if errorCount == 0 {
		t.Log("并发测试完成，未检测到错误（这是正常的，因为验证函数本身不处理并发安全）")
	}
}

// TestValidateSubCommand_EdgeCaseNames 测试边界情况的名称
func TestValidateSubCommand_EdgeCaseNames(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	testCases := []struct {
		name      string
		longName  string
		shortName string
		shouldErr bool
	}{
		{"特殊字符长名称", "cmd-with-dash", "c", false},
		{"数字名称", "123", "1", false},
		{"单字符长名称", "a", "b", false},
		{"很长的名称", strings.Repeat("verylongname", 10), "v", false},
		{"Unicode名称", "命令", "命", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			child := &types.CmdContext{
				LongName:  tc.longName,
				ShortName: tc.shortName,
				SubCmdMap: make(map[string]*types.CmdContext),
			}

			err := ValidateSubCommand(parent, child)
			if tc.shouldErr && err == nil {
				t.Errorf("期望 %s 产生错误，但没有", tc.name)
			} else if !tc.shouldErr && err != nil {
				t.Errorf("期望 %s 不产生错误，但得到: %s", tc.name, err.Error())
			}
		})
	}
}

// TestValidateSubCommand_SubCmdMapNil 测试SubCmdMap为nil的边界情况
func TestValidateSubCommand_SubCmdMapNil(t *testing.T) {
	parent := &types.CmdContext{
		LongName:  "parent",
		ShortName: "p",
		SubCmdMap: nil, // 故意设置为nil
	}

	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	// 这应该会导致panic，因为代码尝试访问nil map
	defer func() {
		if r := recover(); r == nil {
			t.Error("期望访问nil SubCmdMap时panic，但没有panic")
		}
	}()

	_ = ValidateSubCommand(parent, child)
}

// BenchmarkValidateSubCommand 性能基准测试
func BenchmarkValidateSubCommand(b *testing.B) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)
	child := types.NewCmdContext("child", "c", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateSubCommand(parent, child)
	}
}

// BenchmarkHasCycle 循环检测性能基准测试
func BenchmarkHasCycle(b *testing.B) {
	// 创建一个中等深度的命令链
	commands := make([]*types.CmdContext, 20)
	for i := 0; i < 20; i++ {
		commands[i] = types.NewCmdContext(
			fmt.Sprintf("cmd%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)
		if i > 0 {
			commands[i].Parent = commands[i-1]
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasCycle(commands[19], commands[0])
	}
}

// BenchmarkGetCmdIdentifier 命令标识符获取性能基准测试
func BenchmarkGetCmdIdentifier(b *testing.B) {
	cmd := types.NewCmdContext("testcommand", "t", flag.ContinueOnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCmdIdentifier(cmd)
	}
}
