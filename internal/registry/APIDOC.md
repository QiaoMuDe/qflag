package registry // import "gitee.com/MM-Q/qflag/internal/registry"

Package registry 内部注册表管理 本文件实现了内部组件的注册表管理功能，提供统一的组件注册、 查找和管理机制，支持模块化的架构设计。

FUNCTIONS

func RegisterFlag(ctx *types.CmdContext, flag flags.Flag, longName, shortName string) error
    RegisterFlag 注册标志 纯函数设计，通过参数传递所有必要信息

func ValidateFlagNames(ctx *types.CmdContext, longName, shortName string) error
    ValidateFlagNames 验证标志名称

    参数：
      - ctx: 命令上下文
      - longName: 长标志名称
      - shortName: 短标志名称

    返回：
      - error: 如果标志名称无效或已存在，则返回错误；否则返回 nil

