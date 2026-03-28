// Package main 演示错误格式化功能
package main

import (
	"fmt"

	"gitee.com/MM-Q/qflag"
)

func main() {
	fmt.Println("=== QFlag 错误格式化示例 ===")
	fmt.Println()

	// 示例 1: 无原始错误
	fmt.Println("【示例 1】无原始错误:")
	err1 := qflag.NewError("INVALID_VALUE", "invalid port value", nil)
	fmt.Printf("  原始错误: %v\n", err1)
	fmt.Printf("  格式化后: %s\n", qflag.FormatError(err1))
	fmt.Println()

	// 示例 2: 有原始错误
	fmt.Println("【示例 2】有原始错误:")
	originalErr := fmt.Errorf("port must be between 1 and 65535")
	err2 := qflag.NewError("VALIDATION_FAILED", "port validation failed", originalErr)
	fmt.Printf("  原始错误: %v\n", err2)
	fmt.Printf("  格式化后: %s\n", qflag.FormatError(err2))
	fmt.Println()

	// 示例 3: 包装错误
	fmt.Println("【示例 3】包装错误:")
	err3 := qflag.WrapError(err1, "PARSE_ERROR", "failed to parse config")
	fmt.Printf("  原始错误: %v\n", err3)
	fmt.Printf("  格式化后: %s\n", qflag.FormatError(err3))
	fmt.Println()

	// 示例 4: 多层包装错误
	fmt.Println("【示例 4】多层包装错误:")
	baseErr := fmt.Errorf("connection refused")
	layer1 := qflag.NewError("DB_CONNECTION", "database connection failed", baseErr)
	layer2 := qflag.WrapError(layer1, "APP_ERROR", "application startup failed")
	fmt.Printf("  原始错误: %v\n", layer2)
	fmt.Printf("  格式化后: %s\n", qflag.FormatError(layer2))
	fmt.Println()

	// 示例 5: 解析错误包装
	fmt.Println("【示例 5】解析错误包装:")
	parseErr := fmt.Errorf("invalid syntax")
	err5 := qflag.WrapParseError(parseErr, "int", "abc")
	fmt.Printf("  原始错误: %v\n", err5)
	fmt.Printf("  格式化后: %s\n", qflag.FormatError(err5))
	fmt.Println()

	// 示例 6: 非 *Error 类型
	fmt.Println("【示例 6】非 *Error 类型:")
	stdErr := fmt.Errorf("standard error")
	fmt.Printf("  原始错误: %v\n", stdErr)
	fmt.Printf("  格式化后: %s\n", qflag.FormatError(stdErr))
	fmt.Println()

	// 示例 7: nil 错误
	fmt.Println("【示例 7】nil 错误:")
	result := qflag.FormatError(nil)
	fmt.Printf("  格式化结果: '%s' (空字符串)\n", result)
	fmt.Println()

	// 示例 8: 使用预定义错误
	fmt.Println("【示例 8】使用预定义错误:")
	fmt.Printf("  ErrInvalidFlagType: %s\n", qflag.FormatError(qflag.ErrInvalidFlagType))
	fmt.Printf("  ErrFlagNotFound: %s\n", qflag.FormatError(qflag.ErrFlagNotFound))
	fmt.Printf("  ErrParseFailed: %s\n", qflag.FormatError(qflag.ErrParseFailed))
	fmt.Println()

	// 示例 9: 实际应用场景 - 模拟命令行解析错误
	fmt.Println("【示例 9】实际应用场景:")
	simulateParseError()
	fmt.Println()
}

// simulateParseError 模拟一个实际的解析错误场景
func simulateParseError() {
	// 创建命令
	cmd := qflag.NewCmd("myapp", "m", qflag.ContinueOnError)

	// 添加标志
	port := cmd.Int("port", "p", "端口号", 8080)

	// 模拟解析无效值
	err := port.Set("invalid_port")
	if err != nil {
		// 包装错误
		wrappedErr := qflag.WrapError(err, "CLI_PARSE_ERROR", "failed to parse command line arguments")
		fmt.Printf("  解析错误: %s\n", qflag.FormatError(wrappedErr))
	}
}
