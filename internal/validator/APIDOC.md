package validator // import "gitee.com/MM-Q/qflag/internal/validator"

internal/subcmd/validator.go Package validator 内部验证器实现
本文件实现了内部使用的验证器功能，提供命令和标志的验证逻辑， 包括循环引用检测、命名冲突检查等内部验证机制。

FUNCTIONS

func GetCmdIdentifier(cmd *types.CmdContext) string
    GetCmdIdentifier 获取命令的标识字符串，用于错误信息

    参数：
      - cmd: 命令对象

    返回：
      - 命令标识字符串, 如果为空则返回 <nil>

func HasCycleFast(parent, child *types.CmdContext) bool
    HasCycleFast 快速检测父命令和子命令之间是否存在循环依赖

    核心原理： 1. 只检查child的父链向上遍历，避免复杂的子树遍历 2. 利用CLI工具命令层级浅的特点（通常<10层） 3.
    时间复杂度从O(n²)优化到O(d)，其中d是命令深度

    参数:
      - parent: 待添加的父命令上下文
      - child: 待添加的子命令上下文

    返回值:
      - bool: true表示存在循环依赖，false表示安全

    使用场景：
      - 在AddSubCmd函数中调用，防止添加会造成循环依赖的子命令

func ValidateSubCommand(parent, child *types.CmdContext) error
    ValidateSubCommand 验证单个子命令的有效性

    参数：
      - parent: 当前上下文实例
      - child: 待添加的上下文实例

    返回值：
      - error: 验证失败时返回的错误信息, 否则返回nil

