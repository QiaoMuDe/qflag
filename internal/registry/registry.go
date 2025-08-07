// Package registry 内部注册表管理
// 本文件实现了内部组件的注册表管理功能，提供统一的组件注册、
// 查找和管理机制，支持模块化的架构设计。
package registry

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/flags"
	"gitee.com/MM-Q/qflag/internal/types"
)

// RegisterFlag 注册标志
// 纯函数设计，通过参数传递所有必要信息
func RegisterFlag(ctx *types.CmdContext, flag flags.Flag, longName, shortName string) error {
	// 验证标志名称
	if err := ValidateFlagNames(ctx, longName, shortName); err != nil {
		return err
	}

	// 注册到注册表
	return ctx.FlagRegistry.RegisterFlag(&flags.FlagMeta{Flag: flag})
}

// ValidateFlagNames 验证标志名称
//
// 参数：
//   - ctx: 命令上下文
//   - longName: 长标志名称
//   - shortName: 短标志名称
//
// 返回：
//   - error: 如果标志名称无效或已存在，则返回错误；否则返回 nil
func ValidateFlagNames(ctx *types.CmdContext, longName, shortName string) error {
	// 检查长短标志名是否同时为空
	if longName == "" && shortName == "" {
		return fmt.Errorf("flag long name and short name cannot both be empty")
	}

	// 检查长标志相关逻辑
	if longName != "" {
		if err := validateSingleFlagName(ctx, longName, "long name"); err != nil {
			return err
		}
	}

	// 检查短标志相关逻辑
	if shortName != "" {
		if err := validateSingleFlagName(ctx, shortName, "short name"); err != nil {
			return err
		}
	}

	return nil
}

// validateSingleFlagName 验证单个标志名称
//
// 参数：
//   - ctx: 命令上下文
//   - name: 标志名称
//   - nameType: 标志名称类型（长名称或短名称）
//
// 返回：
//   - error: 如果标志名称无效或已存在，则返回错误；否则返回 nil
func validateSingleFlagName(ctx *types.CmdContext, name, nameType string) error {
	// 检查名称是否包含非法字符
	if strings.ContainsAny(name, flags.InvalidFlagChars) {
		return fmt.Errorf("the flag %s '%s' contains illegal characters", nameType, name)
	}

	// 检查标志是否为内置标志
	if ok := ctx.BuiltinFlags.IsBuiltinFlag(name); ok {
		return fmt.Errorf("flag %s '%s' is reserved", nameType, name)
	}

	return nil
}
