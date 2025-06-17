// 定义错误常量
package qflag

// 命令行解析相关错误常量
const (
	ErrFlagParseFailed       = "Parameter parsing error"  // 全局实例标志解析错误
	ErrSubCommandParseFailed = "Subcommand parsing error" // 子命令标志解析错误
	ErrPanicRecovered        = "panic recovered"          // 恐慌捕获错误
)
