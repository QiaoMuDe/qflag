package validator

import (
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
)

// containsIgnoreCase 辅助函数：检查字符串是否包含子字符串（忽略大小写）
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// TestComplexCommandHierarchy 测试复杂的命令层次结构
func TestComplexCommandHierarchy(t *testing.T) {
	// 创建一个复杂的命令树结构
	root := types.NewCmdContext("root", "r", flag.ContinueOnError)

	// 第一层子命令
	cmd1 := types.NewCmdContext("cmd1", "c1", flag.ContinueOnError)
	cmd2 := types.NewCmdContext("cmd2", "c2", flag.ContinueOnError)

	// 第二层子命令
	subcmd1 := types.NewCmdContext("subcmd1", "s1", flag.ContinueOnError)
	subcmd2 := types.NewCmdContext("subcmd2", "s2", flag.ContinueOnError)

	// 第三层子命令
	subsubcmd1 := types.NewCmdContext("subsubcmd1", "ss1", flag.ContinueOnError)

	// 建立层次关系
	cmd1.Parent = root
	cmd2.Parent = root
	subcmd1.Parent = cmd1
	subcmd2.Parent = cmd1
	subsubcmd1.Parent = subcmd1

	// 测试各种验证场景
	testCases := []struct {
		name      string
		parent    *types.CmdContext
		child     *types.CmdContext
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "正常添加第一层子命令",
			parent:    root,
			child:     types.NewCmdContext("newcmd", "n", flag.ContinueOnError),
			shouldErr: false,
		},
		{
			name:      "尝试创建循环：将根命令添加为子命令的子命令",
			parent:    subcmd1,
			child:     root,
			shouldErr: true,
			errMsg:    "cyclic reference detected",
		},
		{
			name:      "尝试创建循环：将父命令添加为孙子命令",
			parent:    subsubcmd1,
			child:     cmd1,
			shouldErr: true,
			errMsg:    "cyclic reference detected",
		},
		{
			name:      "正常添加深层子命令",
			parent:    subsubcmd1,
			child:     types.NewCmdContext("deep", "d", flag.ContinueOnError),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSubCommand(tc.parent, tc.child)

			if tc.shouldErr {
				if err == nil {
					t.Errorf("期望产生错误，但没有错误")
				} else if tc.errMsg != "" && !containsIgnoreCase(err.Error(), tc.errMsg) {
					t.Errorf("期望错误消息包含 '%s'，但得到: %s", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("期望不产生错误，但得到: %s", err.Error())
				}
			}
		})
	}
}

// TestMemoryLeakPrevention 测试内存泄漏预防
func TestMemoryLeakPrevention(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过内存泄漏测试（短测试模式）")
	}

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// 创建大量的命令上下文并进行验证
	for i := 0; i < 1000; i++ {
		parent := types.NewCmdContext(fmt.Sprintf("parent%d", i), fmt.Sprintf("p%d", i), flag.ContinueOnError)

		for j := 0; j < 10; j++ {
			child := types.NewCmdContext(fmt.Sprintf("child%d_%d", i, j), fmt.Sprintf("c%d_%d", i, j), flag.ContinueOnError)
			ValidateSubCommand(parent, child)
		}
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	// 检查内存使用是否在合理范围内
	memIncrease := m2.Alloc - m1.Alloc
	if memIncrease > 50*1024*1024 { // 50MB
		t.Errorf("内存使用增长过多: %d bytes", memIncrease)
	}
}

// TestConcurrentValidation 测试并发验证场景
func TestConcurrentValidation(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	const numGoroutines = 100
	const numOperations = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	// 启动多个goroutine并发执行验证
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				child := types.NewCmdContext(
					fmt.Sprintf("child_%d_%d", goroutineID, j),
					fmt.Sprintf("c_%d_%d", goroutineID, j),
					flag.ContinueOnError,
				)

				err := ValidateSubCommand(parent, child)
				if err != nil {
					errors <- err
				}

				// 添加一些随机延迟来增加竞争条件
				time.Sleep(time.Microsecond * time.Duration(j))
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 收集所有错误
	var allErrors []error
	for err := range errors {
		allErrors = append(allErrors, err)
	}

	// 在并发场景下，不应该有验证错误（因为所有命令名称都是唯一的）
	if len(allErrors) > 0 {
		t.Errorf("并发验证中出现了 %d 个错误，第一个错误: %s", len(allErrors), allErrors[0].Error())
	}
}

// TestConcurrentCycleDetection 测试并发循环检测
func TestConcurrentCycleDetection(t *testing.T) {
	// 创建一个命令链
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

	const numGoroutines = 50
	var wg sync.WaitGroup
	results := make(chan bool, numGoroutines)

	// 并发执行循环检测
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// 测试不同的循环检测场景
			parentIdx := index % 20
			childIdx := (index + 10) % 20

			result := HasCycle(commands[parentIdx], commands[childIdx])
			results <- result
		}(i)
	}

	wg.Wait()
	close(results)

	// 收集结果
	var trueCount, falseCount int
	for result := range results {
		if result {
			trueCount++
		} else {
			falseCount++
		}
	}

	// 验证结果的一致性
	if trueCount == 0 && falseCount == 0 {
		t.Error("没有收到任何循环检测结果")
	}

	t.Logf("循环检测结果: %d个true, %d个false", trueCount, falseCount)
}

// TestExtremeDepthCycle 测试极端深度的循环检测
func TestExtremeDepthCycle(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过极端深度测试（短测试模式）")
	}

	// 创建一个接近最大深度限制的命令链
	const depth = 99
	commands := make([]*types.CmdContext, depth)

	for i := 0; i < depth; i++ {
		commands[i] = types.NewCmdContext(
			fmt.Sprintf("cmd%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)
		if i > 0 {
			commands[i].Parent = commands[i-1]
		}
	}

	// 测试在最大深度边界的循环检测
	start := time.Now()
	result := HasCycle(commands[depth-1], commands[0])
	duration := time.Since(start)

	if !result {
		t.Error("期望在深度链中检测到循环，但返回false")
	}

	// 检查性能：循环检测不应该花费太长时间
	if duration > time.Millisecond*100 {
		t.Errorf("循环检测耗时过长: %v", duration)
	}
}

// TestEdgeCaseValidation 测试边界情况验证
func TestEdgeCaseValidation(t *testing.T) {
	testCases := []struct {
		name        string
		setupParent func() *types.CmdContext
		setupChild  func() *types.CmdContext
		shouldPanic bool
		shouldErr   bool
		errContains string
	}{
		{
			name: "父命令SubCmdMap未初始化",
			setupParent: func() *types.CmdContext {
				ctx := &types.CmdContext{
					LongName:  "parent",
					ShortName: "p",
					SubCmdMap: nil, // 故意设置为nil
				}
				return ctx
			},
			setupChild: func() *types.CmdContext {
				return types.NewCmdContext("child", "c", flag.ContinueOnError)
			},
			shouldPanic: true, // 会panic
		},
		{
			name: "子命令名称包含特殊字符",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				return &types.CmdContext{
					LongName:  "child@#$%",
					ShortName: "c!",
					SubCmdMap: make(map[string]*types.CmdContext),
				}
			},
			shouldErr: false, // 特殊字符应该被允许
		},
		{
			name: "极长的命令名称",
			setupParent: func() *types.CmdContext {
				return types.NewCmdContext("parent", "p", flag.ContinueOnError)
			},
			setupChild: func() *types.CmdContext {
				longName := strings.Repeat("a", 1000)
				return &types.CmdContext{
					LongName:  longName,
					ShortName: "c",
					SubCmdMap: make(map[string]*types.CmdContext),
				}
			},
			shouldErr: false, // 极长名称应该被允许
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parent := tc.setupParent()
			child := tc.setupChild()

			// 捕获可能的panic
			panicked := false
			defer func() {
				if r := recover(); r != nil {
					panicked = true
					if !tc.shouldPanic {
						t.Errorf("意外的panic: %v", r)
					}
				}
			}()

			err := ValidateSubCommand(parent, child)

			if tc.shouldPanic && !panicked {
				t.Error("期望panic，但没有panic")
			}

			if !tc.shouldPanic {
				if tc.shouldErr && err == nil {
					t.Error("期望产生错误，但没有错误")
				} else if !tc.shouldErr && err != nil {
					t.Errorf("期望不产生错误，但得到: %s", err.Error())
				}

				if err != nil && tc.errContains != "" && !containsIgnoreCase(err.Error(), tc.errContains) {
					t.Errorf("期望错误消息包含 '%s'，但得到: %s", tc.errContains, err.Error())
				}
			}
		})
	}
}

// TestStressValidation 压力测试
func TestStressValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过压力测试（短测试模式）")
	}

	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	// 添加大量子命令进行压力测试
	const numCommands = 10000
	start := time.Now()

	for i := 0; i < numCommands; i++ {
		child := types.NewCmdContext(
			fmt.Sprintf("child%d", i),
			fmt.Sprintf("c%d", i),
			flag.ContinueOnError,
		)

		err := ValidateSubCommand(parent, child)
		if err != nil {
			t.Errorf("第 %d 个命令验证失败: %s", i, err.Error())
			break
		}

		// 模拟实际添加到映射中
		parent.SubCmdMap[child.LongName] = child
		if child.ShortName != "" {
			parent.SubCmdMap[child.ShortName] = child
		}
	}

	duration := time.Since(start)
	t.Logf("验证 %d 个命令耗时: %v", numCommands, duration)

	// 性能检查：平均每个命令验证时间不应超过10微秒（更合理的期望）
	avgTime := duration / numCommands
	if avgTime > 10*time.Microsecond {
		t.Errorf("平均验证时间过长: %v", avgTime)
	}
}

// TestRaceConditionDetection 测试竞态条件检测
func TestRaceConditionDetection(t *testing.T) {
	parent := types.NewCmdContext("parent", "p", flag.ContinueOnError)

	const numGoroutines = 20
	var wg sync.WaitGroup

	// 并发添加相同名称的命令，测试竞态条件
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 一半的goroutine尝试添加相同名称的命令
			var child *types.CmdContext
			if id%2 == 0 {
				child = types.NewCmdContext("duplicate", "d", flag.ContinueOnError)
			} else {
				child = types.NewCmdContext(fmt.Sprintf("unique%d", id), fmt.Sprintf("u%d", id), flag.ContinueOnError)
			}

			err := ValidateSubCommand(parent, child)

			// 在没有外部同步的情况下，验证函数本身不应该panic
			if err != nil {
				// 这是预期的，因为会有名称冲突
				t.Logf("Goroutine %d 遇到预期的错误: %s", id, err.Error())
			}
		}(i)
	}

	wg.Wait()
}

// TestCyclicReferenceComplexScenarios 测试复杂的循环引用场景
func TestCyclicReferenceComplexScenarios(t *testing.T) {
	// 场景1：多分支树结构的循环检测
	root := types.NewCmdContext("root", "r", flag.ContinueOnError)
	branch1 := types.NewCmdContext("branch1", "b1", flag.ContinueOnError)
	branch2 := types.NewCmdContext("branch2", "b2", flag.ContinueOnError)
	leaf1 := types.NewCmdContext("leaf1", "l1", flag.ContinueOnError)
	leaf2 := types.NewCmdContext("leaf2", "l2", flag.ContinueOnError)

	// 建立树结构
	branch1.Parent = root
	branch2.Parent = root
	leaf1.Parent = branch1
	leaf2.Parent = branch2

	// 测试跨分支的循环检测
	if HasCycle(leaf1, root) {
		t.Log("正确检测到跨分支循环引用")
	} else {
		t.Error("未能检测到跨分支循环引用")
	}

	// 测试同级节点间的关系（不应该有循环）
	if HasCycle(branch1, branch2) {
		t.Error("错误地检测到同级节点间的循环引用")
	}

	// 场景2：动态修改父子关系后的循环检测
	dynamicChild := types.NewCmdContext("dynamic", "dy", flag.ContinueOnError)
	dynamicChild.Parent = leaf1

	// 现在尝试将root添加为dynamicChild的子命令
	if !HasCycle(dynamicChild, root) {
		t.Error("未能检测到动态修改后的循环引用")
	}
}

// TestValidationWithNilFields 测试包含nil字段的验证
func TestValidationWithNilFields(t *testing.T) {
	// 创建一个部分初始化的命令上下文
	parent := &types.CmdContext{
		LongName:  "parent",
		ShortName: "p",
		SubCmdMap: make(map[string]*types.CmdContext),
		// 其他字段保持nil
	}

	child := &types.CmdContext{
		LongName:  "child",
		ShortName: "c",
		SubCmdMap: make(map[string]*types.CmdContext),
		// 其他字段保持nil
	}

	err := ValidateSubCommand(parent, child)
	if err != nil {
		t.Errorf("部分初始化的命令上下文验证失败: %s", err.Error())
	}
}

// TestPerformanceComparison 性能对比测试
func TestPerformanceComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能对比测试（短测试模式）")
	}

	// 测试不同深度的循环检测性能
	depths := []int{10, 50, 99}

	for _, depth := range depths {
		t.Run(fmt.Sprintf("深度%d", depth), func(t *testing.T) {
			// 创建指定深度的命令链
			commands := make([]*types.CmdContext, depth)
			for i := 0; i < depth; i++ {
				commands[i] = types.NewCmdContext(
					fmt.Sprintf("cmd%d", i),
					fmt.Sprintf("c%d", i),
					flag.ContinueOnError,
				)
				if i > 0 {
					commands[i].Parent = commands[i-1]
				}
			}

			// 测量循环检测性能
			start := time.Now()
			const iterations = 1000

			for i := 0; i < iterations; i++ {
				HasCycle(commands[depth-1], commands[0])
			}

			duration := time.Since(start)
			avgTime := duration / iterations

			t.Logf("深度 %d，平均循环检测时间: %v", depth, avgTime)

			// 性能要求：即使在最大深度下，平均检测时间也不应超过10微秒
			if avgTime > time.Microsecond*10 {
				t.Errorf("深度 %d 的循环检测性能不达标: %v", depth, avgTime)
			}
		})
	}
}
