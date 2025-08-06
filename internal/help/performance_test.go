package help

import (
	"fmt"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/types"
)

// BenchmarkGenerateHelp_LargeScale 基准测试：大规模标志和子命令的帮助信息生成性能
func BenchmarkGenerateHelp_LargeScale(b *testing.B) {
	// 创建包含大量标志和子命令的测试上下文
	ctx := createLargeScaleContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateHelp(ctx)
	}
}

// BenchmarkGenerateHelp_MediumScale 基准测试：中等规模的帮助信息生成性能
func BenchmarkGenerateHelp_MediumScale(b *testing.B) {
	ctx := createMediumScaleContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateHelp(ctx)
	}
}

// BenchmarkGenerateHelp_SmallScale 基准测试：小规模的帮助信息生成性能
func BenchmarkGenerateHelp_SmallScale(b *testing.B) {
	ctx := createSmallScaleContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateHelp(ctx)
	}
}

// TestGenerateHelp_LargeScale_Performance 测试大规模情况下的性能表现
func TestGenerateHelp_LargeScale_Performance(t *testing.T) {
	ctx := createLargeScaleContext()

	// 测试生成时间
	start := time.Now()
	result := GenerateHelp(ctx)
	duration := time.Since(start)

	// 验证结果不为空
	if result == "" {
		t.Error("大规模测试：生成的帮助信息不应为空")
	}

	// 验证性能要求（应在合理时间内完成）
	maxDuration := 100 * time.Millisecond
	if duration > maxDuration {
		t.Errorf("大规模测试：生成帮助信息耗时过长，期望小于 %v，实际耗时 %v", maxDuration, duration)
	}

	t.Logf("大规模测试：生成帮助信息耗时 %v，结果长度 %d 字符", duration, len(result))
}

// TestGenerateHelp_ScalabilityComparison 测试不同规模下的性能对比
func TestGenerateHelp_ScalabilityComparison(t *testing.T) {
	testCases := []struct {
		name        string
		createCtx   func() *types.CmdContext
		maxDuration time.Duration
	}{
		{
			name:        "小规模(10个标志,5个子命令)",
			createCtx:   createSmallScaleContext,
			maxDuration: 10 * time.Millisecond,
		},
		{
			name:        "中等规模(50个标志,20个子命令)",
			createCtx:   createMediumScaleContext,
			maxDuration: 50 * time.Millisecond,
		},
		{
			name:        "大规模(200个标志,100个子命令)",
			createCtx:   createLargeScaleContext,
			maxDuration: 100 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.createCtx()

			start := time.Now()
			result := GenerateHelp(ctx)
			duration := time.Since(start)

			if result == "" {
				t.Errorf("%s：生成的帮助信息不应为空", tc.name)
			}

			if duration > tc.maxDuration {
				t.Errorf("%s：生成耗时 %v 超过预期的 %v", tc.name, duration, tc.maxDuration)
			}

			t.Logf("%s：耗时 %v，结果长度 %d", tc.name, duration, len(result))
		})
	}
}

// TestGenerateHelp_MemoryUsage 测试内存使用情况
func TestGenerateHelp_MemoryUsage(t *testing.T) {
	ctx := createLargeScaleContext()

	// 多次生成帮助信息，检查是否有内存泄漏
	const iterations = 100

	start := time.Now()
	for i := 0; i < iterations; i++ {
		result := GenerateHelp(ctx)
		if result == "" {
			t.Errorf("第 %d 次生成失败", i+1)
			break
		}
	}
	duration := time.Since(start)

	avgDuration := duration / iterations
	t.Logf("内存测试：%d 次生成平均耗时 %v", iterations, avgDuration)

	// 验证平均时间在合理范围内
	if avgDuration > 10*time.Millisecond {
		t.Errorf("平均生成时间过长：%v", avgDuration)
	}
}

// createSmallScaleContext 创建小规模测试上下文
func createSmallScaleContext() *types.CmdContext {
	ctx := createTestContext("smallapp", "sa")
	ctx.Config.UseChinese = true
	ctx.Config.Description = "小规模测试应用"

	// 添加10个标志
	for i := 0; i < 10; i++ {
		addTestFlag(ctx, fmt.Sprintf("flag%d", i), fmt.Sprintf("f%d", i),
			fmt.Sprintf("标志%d的描述", i), "string", fmt.Sprintf("default%d", i))
	}

	// 添加5个子命令
	for i := 0; i < 5; i++ {
		subCtx := createTestContext(fmt.Sprintf("subcmd%d", i), fmt.Sprintf("s%d", i))
		subCtx.Config.Description = fmt.Sprintf("子命令%d的描述", i)
		ctx.SubCmds = append(ctx.SubCmds, subCtx)
	}

	// 添加示例和注意事项
	ctx.Config.Examples = []types.ExampleInfo{
		{Description: "基本用法", Usage: "smallapp --flag1=value1"},
		{Description: "使用子命令", Usage: "smallapp subcmd1 --help"},
	}
	ctx.Config.Notes = []string{
		"这是一个小规模测试应用",
		"用于性能测试",
	}

	return ctx
}

// createMediumScaleContext 创建中等规模测试上下文
func createMediumScaleContext() *types.CmdContext {
	ctx := createTestContext("mediumapp", "ma")
	ctx.Config.UseChinese = true
	ctx.Config.Description = "中等规模测试应用"
	ctx.Config.LogoText = "中等规模测试应用 Logo"
	ctx.Config.ModuleHelps = "模块帮助信息"

	// 添加50个标志
	for i := 0; i < 50; i++ {
		flagType := "string"
		defValue := interface{}(fmt.Sprintf("default%d", i))

		// 混合不同类型的标志
		switch i % 4 {
		case 0:
			flagType = "bool"
			defValue = false
		case 1:
			flagType = "int"
			defValue = i
		case 2:
			flagType = "string"
			defValue = fmt.Sprintf("value%d", i)
		}

		addTestFlag(ctx, fmt.Sprintf("flag%d", i), fmt.Sprintf("f%d", i%26),
			fmt.Sprintf("这是标志%d的详细描述信息", i), flagType, defValue)
	}

	// 添加20个子命令
	for i := 0; i < 20; i++ {
		subCtx := createTestContext(fmt.Sprintf("subcmd%d", i), fmt.Sprintf("s%d", i%26))
		subCtx.Config.Description = fmt.Sprintf("这是子命令%d的详细描述信息", i)

		// 为子命令也添加一些标志
		for j := 0; j < 3; j++ {
			addTestFlag(subCtx, fmt.Sprintf("subflag%d", j), fmt.Sprintf("sf%d", j),
				fmt.Sprintf("子命令标志%d", j), "string", fmt.Sprintf("subdefault%d", j))
		}

		ctx.SubCmds = append(ctx.SubCmds, subCtx)
	}

	// 添加多个示例
	for i := 0; i < 5; i++ {
		ctx.Config.Examples = append(ctx.Config.Examples, types.ExampleInfo{
			Description: fmt.Sprintf("示例%d：演示功能%d", i+1, i+1),
			Usage:       fmt.Sprintf("mediumapp --flag%d=value%d subcmd%d", i, i, i),
		})
	}

	// 添加多个注意事项
	for i := 0; i < 3; i++ {
		ctx.Config.Notes = append(ctx.Config.Notes, fmt.Sprintf("注意事项%d：这是重要的使用说明", i+1))
	}

	return ctx
}

// createLargeScaleContext 创建大规模测试上下文
func createLargeScaleContext() *types.CmdContext {
	ctx := createTestContext("largeapp", "la")
	ctx.Config.UseChinese = true
	ctx.Config.Description = "大规模测试应用，用于测试在大量标志和子命令情况下的性能表现"
	ctx.Config.LogoText = `
    ██╗      █████╗ ██████╗  ██████╗ ███████╗     █████╗ ██████╗ ██████╗ 
    ██║     ██╔══██╗██╔══██╗██╔════╝ ██╔════╝    ██╔══██╗██╔══██╗██╔══██╗
    ██║     ███████║██████╔╝██║  ███╗█████╗      ███████║██████╔╝██████╔╝
    ██║     ██╔══██║██╔══██╗██║   ██║██╔══╝      ██╔══██║██╔═══╝ ██╔═══╝ 
    ███████╗██║  ██║██║  ██║╚██████╔╝███████╗    ██║  ██║██║     ██║     
    ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝    ╚═╝  ╚═╝╚═╝     ╚═╝     
    `
	ctx.Config.ModuleHelps = "这是一个复杂的模块帮助信息，包含了详细的使用说明和配置指南"

	// 添加200个标志
	for i := 0; i < 200; i++ {
		var flagType string
		var defValue interface{}

		// 更多样化的标志类型
		switch i % 6 {
		case 0:
			flagType = "bool"
			defValue = i%2 == 0
		case 1:
			flagType = "int"
			defValue = i * 10
		case 2:
			flagType = "string"
			defValue = fmt.Sprintf("这是一个较长的默认值字符串_%d", i)
		case 3:
			flagType = "string"
			defValue = fmt.Sprintf("path/to/file_%d.txt", i)
		case 4:
			flagType = "int"
			defValue = i * 100
		case 5:
			flagType = "bool"
			defValue = i%3 == 0
		}

		longName := fmt.Sprintf("very-long-flag-name-%d", i)
		shortName := ""
		if i%3 == 0 { // 只有部分标志有短名称
			shortName = fmt.Sprintf("f%d", i%100)
		}

		usage := fmt.Sprintf("这是标志%d的详细使用说明，包含了完整的参数描述和使用场景说明", i)

		addTestFlag(ctx, longName, shortName, usage, flagType, defValue)
	}

	// 添加100个子命令
	for i := 0; i < 100; i++ {
		subCtx := createTestContext(
			fmt.Sprintf("very-long-subcommand-name-%d", i),
			fmt.Sprintf("s%d", i%50),
		)
		subCtx.Config.Description = fmt.Sprintf(
			"这是子命令%d的详细描述信息，包含了完整的功能说明和使用指南，用于演示大规模场景下的性能表现", i)

		// 为每个子命令添加一些标志
		for j := 0; j < 5; j++ {
			addTestFlag(subCtx,
				fmt.Sprintf("sub-flag-%d-%d", i, j),
				fmt.Sprintf("sf%d", j),
				fmt.Sprintf("子命令%d的标志%d", i, j),
				"string",
				fmt.Sprintf("sub-default-%d-%d", i, j))
		}

		ctx.SubCmds = append(ctx.SubCmds, subCtx)
	}

	// 添加大量示例
	for i := 0; i < 10; i++ {
		ctx.Config.Examples = append(ctx.Config.Examples, types.ExampleInfo{
			Description: fmt.Sprintf("复杂示例%d：演示高级功能%d的使用方法", i+1, i+1),
			Usage: fmt.Sprintf("largeapp --very-long-flag-name-%d=value%d very-long-subcommand-name-%d --sub-flag-%d-0=subvalue",
				i*10, i, i*5, i*5),
		})
	}

	// 添加大量注意事项
	for i := 0; i < 8; i++ {
		ctx.Config.Notes = append(ctx.Config.Notes,
			fmt.Sprintf("重要注意事项%d：这是一个详细的使用说明，包含了重要的配置信息和最佳实践建议", i+1))
	}

	return ctx
}

// BenchmarkSortFlags_LargeScale 基准测试：大规模标志排序性能
func BenchmarkSortFlags_LargeScale(b *testing.B) {
	// 创建大量标志信息
	flags := make([]flagInfo, 1000)
	for i := 0; i < 1000; i++ {
		flags[i] = flagInfo{
			longFlag:  fmt.Sprintf("flag-%d", i),
			shortFlag: fmt.Sprintf("f%d", i%100),
			usage:     fmt.Sprintf("标志%d的使用说明", i),
			defValue:  fmt.Sprintf("default%d", i),
			typeStr:   "<string>",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次都复制一份进行排序，避免影响下次测试
		testFlags := make([]flagInfo, len(flags))
		copy(testFlags, flags)
		sortFlags(testFlags)
	}
}

// BenchmarkSortSubCommands_LargeScale 基准测试：大规模子命令排序性能
func BenchmarkSortSubCommands_LargeScale(b *testing.B) {
	// 创建大量子命令
	subCmds := make([]*types.CmdContext, 500)
	for i := 0; i < 500; i++ {
		subCmds[i] = createTestContext(fmt.Sprintf("subcmd-%d", i), fmt.Sprintf("s%d", i%100))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 每次都复制一份进行排序
		testSubCmds := make([]*types.CmdContext, len(subCmds))
		copy(testSubCmds, subCmds)
		sortSubCommands(testSubCmds)
	}
}
