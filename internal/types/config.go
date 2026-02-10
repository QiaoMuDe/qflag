package types

// MutexGroup 互斥组定义
//
// MutexGroup 定义了一组互斥的标志, 其中最多只能有一个被设置。
// 当用户设置了互斥组中的多个标志时, 解析器会返回错误。
//
// 字段说明:
//   - Name: 互斥组名称, 用于错误提示和标识
//   - Flags: 互斥组中的标志名称列表
//   - AllowNone: 是否允许一个都不设置
//
// 使用场景:
//   - 输出格式互斥 (如 --json 和 --xml 不能同时使用)
//   - 操作模式互斥 (如 --create 和 --update 不能同时使用)
//   - 必选选项 (如必须指定 --file 或 --url 中的一个)
type MutexGroup struct {
	Name      string   // 互斥组名称, 用于错误提示和标识
	Flags     []string // 互斥组中的标志名称列表
	AllowNone bool     // 是否允许一个都不设置
}

// CmdConfig 命令配置类型
type CmdConfig struct {
	Version     string            // 版本号
	UseChinese  bool              // 是否使用中文
	EnvPrefix   string            // 环境变量前缀
	UsageSyntax string            // 命令使用语法
	Example     map[string]string // 示例使用, key为描述, value为示例命令
	Notes       []string          // 注意事项
	LogoText    string            // 命令logo文本
	MutexGroups []MutexGroup      // 互斥组列表
}

// NewCmdConfig 创建新的命令配置
//
// 返回值:
//   - *CmdConfig: 新创建的 CmdConfig 实例, 初始化为零值
func NewCmdConfig() *CmdConfig {
	return &CmdConfig{
		Version:     "",
		UseChinese:  false,
		EnvPrefix:   "",
		UsageSyntax: "",
		Example:     map[string]string{},
		Notes:       []string{},
		LogoText:    "",
		MutexGroups: []MutexGroup{},
	}
}
