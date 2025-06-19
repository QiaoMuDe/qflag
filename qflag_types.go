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
	// 英文模板
	cmdNameTemplate          = "Name: %s\n\n"                                                                                                                      // 命令名称
	cmdNameWithShortTemplate = "Name: %s(%s)\n\n"                                                                                                                  // 命令名称和短名称
	cmdDescriptionTemplate   = "Desc: %s\n\n"                                                                                                                      // 命令描述
	optionsHeaderTemplate    = "Options:\n"                                                                                                                        // 选项头部
	optionTemplate1          = "  -%s, --%-*s %s (default: %s)\n"                                                                                                  // 选项模板1
	optionTemplate2          = "  --%-*s %s (default: %s)\n"                                                                                                       // 选项模板2
	subCmdsHeaderTemplate    = "\nSubcommands:\n"                                                                                                                  // 子命令头部
	subCmdTemplate           = "  %s\t%s\n"                                                                                                                        // 子命令模板
	notesHeaderTemplate      = "\nNotes:\n"                                                                                                                        // 注意事项头部
	noteItemTemplate         = "  %d. %s\n"                                                                                                                        // 注意事项项
	defaultNote              = "In the case where both long options and short options are used at the same time, the option specified last shall take precedence." // 默认注意事项

	// 中文模板
	cmdNameTemplateCN          = "名称: %s\n\n"                   // 命令名称（中文）
	cmdNameWithShortTemplateCN = "名称: %s(%s)\n\n"               // 命令名称和短名称（中文）
	cmdDescriptionTemplateCN   = "描述: %s\n\n"                   // 命令描述（中文）
	optionsHeaderTemplateCN    = "选项:\n"                        // 选项头部（中文）
	optionTemplate1CN          = "  -%s, --%-*s %s (默认值: %s)\n" // 选项模板1(中文, 动态宽度)
	optionTemplate2CN          = "  --%-*s %s (默认值: %s)\n"      // 选项模板2(中文, 动态宽度)
	subCmdsHeaderTemplateCN    = "\n子命令:\n"                     // 子命令头部（中文）
	subCmdTemplateCN           = "  %s\t%s\n"                   // 子命令模板（中文）
	notesHeaderTemplateCN      = "\n注意事项:\n"                    // 注意事项头部（中文）
	noteItemTemplateCN         = "  %d、%s\n"                    // 注意事项项（中文）
	defaultNoteCN              = "当长选项和短选项同时使用时，最后指定的选项将优先生效。"  // 默认注意事项（中文）
)
