package qerr // import "gitee.com/MM-Q/qflag/qerr"


CONSTANTS

const (
	ErrFlagParseFailed       = "Parameter parsing error"  // 全局实例标志解析错误
	ErrSubCommandParseFailed = "Subcommand parsing error" // 子命令标志解析错误
	ErrPanicRecovered        = "panic recovered"          // 恐慌捕获错误
	ErrValidationFailed      = "Validation failed"        // 参数验证失败错误
)
    命令行解析相关错误常量


FUNCTIONS

func JoinErrors(errors []error) error
    JoinErrors 将错误切片合并为单个错误，并去除重复错误

func NewValidationError(message string) error
    NewValidationError 创建一个新的验证错误

func NewValidationErrorf(format string, v ...interface{}) error
    NewValidationErrorf 创建一个格式化的验证错误

