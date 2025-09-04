package cmd

import (
	"flag" // Go standard library for ErrorHandling
	"fmt"

	"gitee.com/MM-Q/qflag/flags"
)

// ActionFunc 定义了命令执行的动作函数类型。
// 当命令被解析并执行时，这个函数会被调用。
type ActionFunc func(c *Cmd) error

// FlagBlueprint 以声明式的方式描述一个命令行标志的配置。
type FlagBlueprint struct {
	// Type 指定了标志的数据类型，例如：flags.FlagTypeString, flags.FlagTypeInt。
	Type flags.FlagType
	// Name 是标志的长名称，例如："config"。
	Name string
	// ShortName 是标志的短名称，例如："c"。
	ShortName string
	// Usage 是该标志的帮助说明文本。
	Usage string
	// DefaultValue 是标志的默认值，其类型应与 Type 字段匹配。
	DefaultValue interface{}
}

// CommandBlueprint 提供了一个完整的、声明式的命令定义蓝图。
// 它允许将一个命令的所有属性（名称、标志、子命令、动作）集中在一个结构体中进行描述。
type CommandBlueprint struct {
	// Name 是命令的长名称。
	Name string
	// ShortName 是命令的短名称。
	ShortName string
	// Usage 是命令的简短用法说明，通常在帮助信息的第一行显示。
	Usage string
	// Description 是对命令功能的更详细描述。
	Description string
	// Flags 是一个包含该命令所有标志定义的切片。
	Flags []FlagBlueprint
	// Subcommands 是该命令的子命令蓝图列表。
	Subcommands []*CommandBlueprint
	// Action 是当该命令被成功解析后要执行的回调函数。
	Action ActionFunc
}

// NewCmdFromBlueprint 根据给定的蓝图递归地创建并配置一个 Cmd 实例及其所有子命令。
// 这个函数是连接声明式蓝图和 qflag 命令式 API 的桥梁。
func NewCmdFromBlueprint(blueprint *CommandBlueprint) (*Cmd, error) {
	if blueprint == nil {
		return nil, fmt.Errorf("command blueprint cannot be nil")
	}

	// 1. 使用蓝图的基本信息创建 Cmd 实例。
	cmd := NewCmd(blueprint.Name, blueprint.ShortName, flag.ContinueOnError)
	cmd.SetDescription(blueprint.Description)
	cmd.SetUsageSyntax(blueprint.Usage)

	// TODO: 将 Action 存储在 Cmd 实例中，以便在 Parse 之后执行。
	// 这需要在 Cmd 结构体中添加一个字段来存储 ActionFunc。
	// cmd.action = blueprint.Action

	// 2. 遍历蓝图中的标志定义，并将其添加到 Cmd 实例中。
	for _, flagBP := range blueprint.Flags {
		if err := addFlagFromBlueprint(cmd, &flagBP); err != nil {
			// 如果添加标志失败，返回一个包含上下文的错误。
			return nil, fmt.Errorf("failed to add flag '%s' to command '%s': %w", flagBP.Name, cmd.Name(), err)
		}
	}

	// 3. 递归地为所有子命令蓝图创建 Cmd 实例，并添加到当前命令。
	for _, subBlueprint := range blueprint.Subcommands {
		subCmd, err := NewCmdFromBlueprint(subBlueprint)
		if err != nil {
			return nil, err // 如果子命令创建失败，则直接向上传递错误。
		}
		if err := cmd.AddSubCmd(subCmd); err != nil {
			return nil, fmt.Errorf("failed to add subcommand '%s' to '%s': %w", subCmd.Name(), cmd.Name(), err)
		}
	}

	return cmd, nil
}

// addFlagFromBlueprint 根据 FlagBlueprint 的定义，调用相应的 Cmd 方法来添加一个标志。
func addFlagFromBlueprint(c *Cmd, blueprint *FlagBlueprint) error {
	// 使用 switch 语句根据标志类型来分派到正确的创建函数。
	switch blueprint.Type {
	case flags.FlagTypeString:
		// 类型断言，获取正确的默认值类型。
		defValue, _ := blueprint.DefaultValue.(string)
		c.String(blueprint.Name, blueprint.ShortName, defValue, blueprint.Usage)

	case flags.FlagTypeInt:
		defValue, _ := blueprint.DefaultValue.(int)
		c.Int(blueprint.Name, blueprint.ShortName, defValue, blueprint.Usage)

	case flags.FlagTypeInt64:
		defValue, _ := blueprint.DefaultValue.(int64)
		c.Int64(blueprint.Name, blueprint.ShortName, defValue, blueprint.Usage)

	case flags.FlagTypeBool:
		defValue, _ := blueprint.DefaultValue.(bool)
		c.Bool(blueprint.Name, blueprint.ShortName, defValue, blueprint.Usage)

	case flags.FlagTypeSize:
		defValue, _ := blueprint.DefaultValue.(int64)
		c.Size(blueprint.Name, blueprint.ShortName, defValue, blueprint.Usage)

	// 可以在这里继续添加对其他标志类型的支持，例如：
	// case flags.FlagTypeDuration:
	// case flags.FlagTypeStringSlice:

	default:
		// 如果遇到不支持的标志类型，则返回错误。
		return fmt.Errorf("unsupported flag type: %v", blueprint.Type)
	}
	return nil
}
