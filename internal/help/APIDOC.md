package help // import "gitee.com/MM-Q/qflag/internal/help"

Package help 帮助信息生成器 本文件实现了命令行帮助信息的自动生成功能，包括标志列表、用法说明、 子命令信息等帮助内容的格式化和输出。

Package help 帮助信息排序和组织 本文件实现了帮助信息的排序和组织功能，包括标志排序、子命令排序等， 确保帮助信息以合理的顺序展示给用户。

Package help 测试辅助工具 本文件提供了帮助信息模块的测试辅助函数和工具， 用于支持帮助信息生成和格式化功能的单元测试。

Package help 帮助信息输出和格式化 本文件实现了帮助信息的输出和格式化功能，包括不同格式的帮助信息写入器， 支持多种输出格式和样式的帮助信息展示。

CONSTANTS

const (
	// 默认最大宽度，当计算失败时使用
	DefaultMaxWidth = 30

	// 描述信息与选项之间的间距
	DescriptionPadding = 5

	// 子命令名称分隔符长度 (", " 的长度)
	SubCmdSeparatorLen = 2

	// 子命令对齐额外空格数
	SubCmdAlignSpaces = 5

	// 最小填充空格数
	MinPadding = 1
)
    帮助信息格式化常量


VARIABLES

var ChineseTemplate = HelpTemplate{
	CmdName:              "名称: %s\n\n",
	UsagePrefix:          "用法: ",
	UsageSubCmd:          " [子命令]",
	UsageInfoWithOptions: " [选项]\n\n",
	UsageGlobalOptions:   " [全局选项]",
	CmdNameWithShort:     "名称: %s, %s\n\n",
	CmdDescription:       "描述: %s\n\n",
	OptionsHeader:        "选项:\n",
	Option1:              "  --%s, -%s %s",
	Option2:              "  --%s %s",
	Option3:              "  -%s %s",
	OptionDefault:        "%s%*s%s (默认值: %s)\n",
	SubCmdsHeader:        "\n子命令:\n",
	SubCmd:               "  %s\t%s\n",
	SubCmdWithShort:      "  %s, %s\t%s\n",
	NotesHeader:          "\n注意事项:\n",
	NoteItem:             "  %d、%s\n",
	DefaultNote:          "当长选项和短选项同时使用时，最后指定的选项将优先生效。",
	ExamplesHeader:       "\n示例:\n",
	ExampleItem:          "  %d、%s\n     %s\n",
}
    中文模板实例

var EnglishTemplate = HelpTemplate{
	CmdName:              "Name: %s\n\n",
	UsagePrefix:          "Usage: ",
	UsageSubCmd:          " [subcmd]",
	UsageInfoWithOptions: " [options]\n\n",
	UsageGlobalOptions:   " [global options]",
	CmdNameWithShort:     "Name: %s, %s\n\n",
	CmdDescription:       "Desc: %s\n\n",
	OptionsHeader:        "Options:\n",
	Option1:              "  --%s, -%s %s",
	Option2:              "  --%s %s",
	Option3:              "  -%s %s",
	OptionDefault:        "%s%*s%s (default: %s)\n",
	SubCmdsHeader:        "\nSubCmds:\n",
	SubCmd:               "  %s\t%s\n",
	SubCmdWithShort:      "  %s, %s\t%s\n",
	NotesHeader:          "\nNotes:\n",
	NoteItem:             "  %d. %s\n",
	DefaultNote:          "In the case where both long options and short options are used at the same time,\n the option specified last shall take precedence.",
	ExamplesHeader:       "\nExamples:\n",
	ExampleItem:          "  %d. %s\n     %s\n",
}
    英文模板实例


FUNCTIONS

func GenerateHelp(ctx *types.CmdContext) string
    GenerateHelp 生成帮助信息 纯函数设计，不依赖任何结构体状态


TYPES

type HelpTemplate struct {
	CmdName              string // 命令名称模板
	CmdNameWithShort     string // 命令名称带短名称模板
	CmdDescription       string // 命令描述模板
	UsagePrefix          string // 用法说明前缀模板
	UsageSubCmd          string // 用法说明子命令模板
	UsageInfoWithOptions string // 带选项的用法说明信息模板
	UsageGlobalOptions   string // 全局选项部分
	OptionsHeader        string // 选项头部模板
	Option1              string // 选项模板(带短选项)
	Option2              string // 选项模板(无短选项)
	Option3              string // 选项模板(无长选项)
	OptionDefault        string // 选项模板的默认值
	SubCmdsHeader        string // 子命令头部模板
	SubCmd               string // 子命令模板
	SubCmdWithShort      string // 子命令带短名称模板
	NotesHeader          string // 注意事项头部模板
	NoteItem             string // 注意事项项模板
	DefaultNote          string // 默认注意事项
	ExamplesHeader       string // 示例信息头部模板
	ExampleItem          string // 示例信息项模板
}
    HelpTemplate 帮助信息模板结构体

type NamedItem interface {
	GetLongName() string
	GetShortName() string
}
    NamedItem 表示具有长名称和短名称的项目接口

