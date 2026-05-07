// Package main 提供标志构造函数使用示例
//
// 本示例展示了如何使用 qflag 项目中的各种标志构造函数。
// 主要演示:
//   - 基础类型标志 (String, Int, Bool)
//   - 数值类型标志 (Float64, Uint, Uint16)
//   - 特殊类型标志 (Enum, Duration, Time, Size)
//   - 集合类型标志 (StringSlice, IntSlice, Map)
//   - 枚举标志和辅助函数的使用
package main

import (
	"fmt"
	"time"

	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
	"gitee.com/MM-Q/qflag/internal/utils"
)

func main() {
	// 创建根命令
	root := cmd.NewCmd("flag-constructors", "fc", types.ContinueOnError)
	root.SetDesc("标志构造函数使用示例")

	// 基础类型标志
	name := root.String("name", "n", "用户名", "guest")
	age := root.Int("age", "a", "年龄", 18)
	enabled := root.Bool("enabled", "e", "是否启用", false)

	// 数值类型标志
	score := root.Float64("score", "s", "分数", 0.0)
	count := root.Uint("count", "c", "计数", 0)
	port := root.Uint16("port", "p", "端口号", 8080)

	// 特殊类型标志
	mode := root.Enum("mode", "m", "运行模式", "debug", []string{"debug", "release", "test"})
	timeout := root.Duration("timeout", "t", "超时时间", 30*time.Second)
	startTime := root.Time("start", "st", "开始时间", time.Now())
	maxSize := root.Size("max-size", "ms", "最大大小", 1024*1024) // 1MB

	// 集合类型标志
	tags := root.StringSlice("tags", "tg", "标签列表", []string{})
	scores := root.IntSlice("scores", "sc", "分数列表", []int{})
	config := root.Map("config", "cfg", "配置项", map[string]string{})

	// 解析命令行参数
	if err := root.Parse(nil); err != nil {
		fmt.Printf("参数解析失败: %v\n", err)
		return
	}

	// 打印标志的默认值
	fmt.Println("=== 标志默认值 ===")
	fmt.Printf("用户名: %s\n", name.Get())
	fmt.Printf("年龄: %d\n", age.Get())
	fmt.Printf("是否启用: %t\n", enabled.Get())
	fmt.Printf("分数: %.2f\n", score.Get())
	fmt.Printf("计数: %d\n", count.Get())
	fmt.Printf("端口号: %d\n", port.Get())
	fmt.Printf("运行模式: %s\n", mode.Get())
	fmt.Printf("超时时间: %v\n", timeout.Get())
	fmt.Printf("开始时间: %v\n", startTime.Get())
	fmt.Printf("最大大小: %d bytes\n", maxSize.Get())
	fmt.Printf("标签列表: %v\n", tags.Get())
	fmt.Printf("分数列表: %v\n", scores.Get())
	fmt.Printf("配置项: %v\n", config.Get())

	// 演示枚举辅助函数的使用
	fmt.Println("\n=== 枚举辅助函数演示 ===")

	// 1. 整数切片转换为字符串切片
	intLevels := []int{0, 1, 2, 3}
	strLevels, err := utils.ToStrSlice(intLevels)
	if err != nil {
		fmt.Printf("转换整数切片失败: %v\n", err)
	} else {
		fmt.Printf("整数切片 %v 转换为字符串切片: %v\n", intLevels, strLevels)
	}

	// 2. 混合类型切片
	mixedValues := []interface{}{"enabled", "disabled", true, false, 1, 0}
	mixedStr, err := utils.ToStrSlice(mixedValues)
	if err != nil {
		fmt.Printf("转换混合类型切片失败: %v\n", err)
	} else {
		fmt.Printf("混合类型切片 %v 转换为字符串切片: %v\n", mixedValues, mixedStr)
	}

	// 3. 创建自定义枚举标志
	levelFlag := flag.NewEnumFlag("level", "l", "日志级别", "0", strLevels)
	allowedValues := levelFlag.GetAllowedValues()
	fmt.Printf("日志级别枚举标志允许的值: %v\n", allowedValues)

	// 4. 测试枚举标志的设置和获取
	fmt.Println("\n=== 枚举标志设置测试 ===")
	testValues := []string{"0", "1", "2", "3", "invalid"}
	for _, val := range testValues {
		if err := levelFlag.Set(val); err != nil {
			fmt.Printf("设置值 '%s' 失败: %v\n", val, err)
		} else {
			fmt.Printf("设置值 '%s' 成功, 当前值: '%s'\n", val, levelFlag.Get())
		}
	}

	// 演示集合类型标志的使用
	fmt.Println("\n=== 集合类型标志演示 ===")

	// 设置字符串切片
	if err := tags.Set("tag1,tag2,tag3"); err != nil {
		fmt.Printf("设置标签列表失败: %v\n", err)
	} else {
		fmt.Printf("标签列表: %v\n", tags.Get())
	}

	// 设置整数切片
	if err := scores.Set("90,85,95"); err != nil {
		fmt.Printf("设置分数列表失败: %v\n", err)
	} else {
		fmt.Printf("分数列表: %v\n", scores.Get())
	}

	// 设置映射
	if err := config.Set("key1=value1,key2=value2"); err != nil {
		fmt.Printf("设置配置项失败: %v\n", err)
	} else {
		fmt.Printf("配置项: %v\n", config.Get())
	}

	// 演示标志类型信息
	fmt.Println("\n=== 标志类型信息 ===")
	fmt.Printf("name 标志类型: %v\n", name.Type())
	fmt.Printf("age 标志类型: %v\n", age.Type())
	fmt.Printf("enabled 标志类型: %v\n", enabled.Type())
	fmt.Printf("score 标志类型: %v\n", score.Type())
	fmt.Printf("count 标志类型: %v\n", count.Type())
	fmt.Printf("port 标志类型: %v\n", port.Type())
	fmt.Printf("mode 标志类型: %v\n", mode.Type())
	fmt.Printf("timeout 标志类型: %v\n", timeout.Type())
	fmt.Printf("startTime 标志类型: %v\n", startTime.Type())
	fmt.Printf("maxSize 标志类型: %v\n", maxSize.Type())
	fmt.Printf("tags 标志类型: %v\n", tags.Type())
	fmt.Printf("scores 标志类型: %v\n", scores.Type())
	fmt.Printf("config 标志类型: %v\n", config.Type())

	// 演示标志设置状态
	fmt.Println("\n=== 标志设置状态 ===")
	fmt.Printf("name 是否设置: %t\n", name.IsSet())
	fmt.Printf("age 是否设置: %t\n", age.IsSet())
	fmt.Printf("enabled 是否设置: %t\n", enabled.IsSet())
	fmt.Printf("tags 是否设置: %t\n", tags.IsSet())
	fmt.Printf("scores 是否设置: %t\n", scores.IsSet())
	fmt.Printf("config 是否设置: %t\n", config.IsSet())
}
