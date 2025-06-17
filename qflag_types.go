package qflag

// 定义非法字符集常量, 防止非法的标志名称
const invalidFlagChars = " !@#$%^&*(){}[]|\\;:'\"<>,.?/"

// 标志类型
type FlagType int

const (
	FlagTypeInt      FlagType = iota + 1 // 整数类型
	FlagTypeString                       // 字符串类型
	FlagTypeBool                         // 布尔类型
	FlagTypeFloat                        // 浮点数类型
	FlagTypeSlice                        // 切片类型
	FlagTypeEnum                         // 枚举类型
	FlagTypeDuration                     // 时间间隔类型
)

// 帮助信息模板常量
const (
	cmdNameTemplate          = "Command: %s\n\n"                 // 命令名称
	cmdNameWithShortTemplate = "Command: %s(%s)\n\n"             // 命令名称和短名称
	cmdDescriptionTemplate   = "Description: %s\n\n"             // 命令描述
	optionsHeaderTemplate    = "Options:\n"                      // 选项头部
	optionTemplate1          = "  -%s, --%s\t%s (default: %s)\n" // 选项模板1
	optionTemplate2          = "  --%s\t%s (default: %s)\n"      // 选项模板2
	subCmdsHeaderTemplate    = "\nSubcommands:\n"                // 子命令头部
	subCmdTemplate           = "  %s\t%s\n"                      // 子命令模板
	notesHeaderTemplate      = "\nNotes:\n"                      // 注意事项头部
	noteItemTemplate         = "  %d、%s\n"                       // 注意事项项
)
