package parser // import "gitee.com/MM-Q/qflag/internal/parser"

Package parser 环境变量解析和处理 本文件实现了环境变量的解析和处理逻辑，支持从环境变量中读取标志值，
为命令行参数提供环境变量绑定和默认值设置功能。

Package parser 命令行参数解析器 本文件实现了命令行参数的解析逻辑，包括标志解析、参数分离、
子命令识别等核心解析功能，为命令行参数处理提供基础支持。

FUNCTIONS

func LoadEnvVars(ctx *types.CmdContext) error
    LoadEnvVars 从环境变量加载参数值 纯函数设计，不依赖结构体状态

    参数:
      - ctx: 命令上下文

    返回值:
      - error: 错误信息

func ParseCommand(ctx *types.CmdContext, args []string) (err error)
    ParseCommand 解析单个命令的标志和参数

    参数:
      - ctx: 命令上下文
      - args: 命令行参数

    返回值:
      - error: 如果解析失败，返回错误信息

