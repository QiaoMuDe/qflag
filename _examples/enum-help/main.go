// Package main 提供 EnumHelp 枚举帮助生成示例
//
// 本示例展示了 qflag.EnumHelp 函数的使用效果：
//   - 纯枚举值对齐
//   - 带描述的枚举选项对齐
//   - 混合 ASCII 与中文的显示宽度对齐
//   - 空值选项过滤
package main

import (
	"fmt"

	"gitee.com/MM-Q/qflag"
)

func main() {
	// 演示：纯值无描述
	fmt.Println("=== 示例 1: 纯值无描述 ===")
	options1 := []string{"debug", "info", "warn", "error"}
	result1 := qflag.EnumHelp("日志级别", options1, "\t")
	fmt.Println(result1)
	fmt.Println()

	// 演示：带描述对齐
	fmt.Println("=== 示例 2: 带描述枚举选项 ===")
	options2 := []string{
		"debug: 调试输出，最详细",
		"info:  普通信息输出",
		"warn:  警告提示",
		"error: 错误输出",
	}
	result2 := qflag.EnumHelp("日志级别", options2, "\t")
	fmt.Println(result2)
	fmt.Println()

	// 演示：混合中文与英文（验证显示宽度对齐）
	fmt.Println("=== 示例 3: 混合中文与英文（双宽对齐） ===")
	options3 := []string{
		"debug:   调试模式",
		"release: 发布模式",
		"auto:    自动模式",
		"手动:    手动模式",
	}
	result3 := qflag.EnumHelp("运行模式", options3, "\t")
	fmt.Println(result3)
	fmt.Println()

	// 演示：中文枚举值（每个值都是中文，验证对齐）
	fmt.Println("=== 示例 4: 中文枚举值 ===")
	options4 := []string{
		"调试:  调试输出",
		"信息:  普通信息",
		"警告:  警告提示",
		"错误:  错误输出",
		"跟踪:  详细跟踪",
	}
	result4 := qflag.EnumHelp("日志级别", options4, "\t")
	fmt.Println(result4)
	fmt.Println()

	// 演示：混合宽度（验证宽度计算正确）
	fmt.Println("=== 示例 5: 混合宽度枚举值 ===")
	options5 := []string{
		"开: 启用该功能",
		"关: 禁用该功能",
		"自动: 根据环境自动判断",
	}
	result5 := qflag.EnumHelp("功能状态", options5, "\t")
	fmt.Println(result5)
	fmt.Println()

	// 演示：包含空值选项会被自动跳过
	fmt.Println("=== 示例 6: 自动跳过空值选项 ===")
	options6 := []string{
		"auto: 自动模式",
		"", // 空值，会被跳过
		"manual: 手动模式",
		": debug", // 只有描述没有值，会被跳过
		"debug: 调试模式",
	}
	result6 := qflag.EnumHelp("运行模式（含空值）", options6, "\t")
	fmt.Println(result6)
}
