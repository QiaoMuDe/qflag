package types

import (
	"sync"

	"gitee.com/MM-Q/qflag/flags"
)

// BuiltinFlags 内置标志结构体
type BuiltinFlags struct {
	Help       *flags.BoolFlag // 标志-帮助
	Version    *flags.BoolFlag // 标志-版本
	Completion *flags.EnumFlag // 标志-自动完成
	NameMap    sync.Map        // 内置标志名称映射
}

// NewBuiltinFlags 创建内置标志实例
func NewBuiltinFlags() *BuiltinFlags {
	return &BuiltinFlags{
		Help:       &flags.BoolFlag{}, // 标志-帮助
		Version:    &flags.BoolFlag{}, // 标志-版本
		Completion: &flags.EnumFlag{}, // 标志-自动完成
		NameMap:    sync.Map{},        // 内置标志名称映射
	}
}

// IsBuiltinFlag 检查是否为内置标志
func (bf *BuiltinFlags) IsBuiltinFlag(name string) bool {
	_, exists := bf.NameMap.Load(name)
	return exists
}

// MarkAsBuiltin 标记为内置标志
func (bf *BuiltinFlags) MarkAsBuiltin(names ...string) {
	for _, name := range names {
		bf.NameMap.Store(name, true)
	}
}
