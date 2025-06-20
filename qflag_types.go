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

// HelpTemplate 帮助信息模板结构体
type HelpTemplate struct {
	CmdName          string // 命令名称模板
	CmdNameWithShort string // 命令名称带短名称模板
	CmdDescription   string // 命令描述模板
	OptionsHeader    string // 选项头部模板
	Option1          string // 选项模板(带短选项)
	Option2          string // 选项模板(无短选项)
	OptionDefault    string // 选项模板的默认值
	SubCmdsHeader    string // 子命令头部模板
	SubCmd           string // 子命令模板
	SubCmdWithShort  string // 子命令带短名称模板
	NotesHeader      string // 注意事项头部模板
	NoteItem         string // 注意事项项模板
	DefaultNote      string // 默认注意事项
}

// 英文模板实例
var EnglishTemplate = HelpTemplate{
	CmdName:          "Name: %s\n\n",                                                                                                                        // 命令名称模板
	CmdNameWithShort: "Name: %s(%s)\n\n",                                                                                                                    // 命令名称带短名称模板
	CmdDescription:   "Desc: %s\n\n",                                                                                                                        // 命令描述模板
	OptionsHeader:    "Options:\n",                                                                                                                          // 选项头部模板
	Option1:          "  -%s, --%s",                                                                                                                         // 选项模板(带短选项)
	Option2:          "  --%s",                                                                                                                              // 选项模板(无短选项)
	OptionDefault:    "%s%*s%s (default: %s)\n",                                                                                                             // 选项模板默认值
	SubCmdsHeader:    "\nSubCmds:\n",                                                                                                                        // 子命令头部模板
	SubCmd:           "  %s\t%s\n",                                                                                                                          // 子命令模板
	SubCmdWithShort:  "  %s, %s\t%s\n",                                                                                                                      // 子命令模板(带短选项)
	NotesHeader:      "\nNotes:\n",                                                                                                                          // 注意事项头部模板
	NoteItem:         "  %d. %s\n",                                                                                                                          // 注意事项模板
	DefaultNote:      "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.", // 默认注意事项
}

// 中文模板实例
var ChineseTemplate = HelpTemplate{
	CmdName:          "名称: %s\n\n",                  // 命令名称模板
	CmdNameWithShort: "名称: %s(%s)\n\n",              // 命令名称带短名称模板
	CmdDescription:   "描述: %s\n\n",                  // 命令描述模板
	OptionsHeader:    "选项:\n",                       // 选项头部模板
	Option1:          "  -%s, --%s",                 // 选项模板(带短选项)
	Option2:          "  --%s",                      // 选项模板(无短选项)
	OptionDefault:    "%s%*s%s (默认值: %s)\n",         // 选项模板默认值
	SubCmdsHeader:    "\n子命令:\n",                    // 子命令头部模板
	SubCmd:           "  %s\t%s\n",                  // 子命令模板
	SubCmdWithShort:  "  %s, %s\t%s\n",              // 子命令模板(带短选项)
	NotesHeader:      "\n注意事项:\n",                   //注意事项头部模板
	NoteItem:         "  %d、%s\n",                   //注意事项模板
	DefaultNote:      "当长选项和短选项同时使用时，最后指定的选项将优先生效。", //默认注意事项
}
