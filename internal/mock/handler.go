package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// CustomHandler 自定义处理器, 用于测试
type CustomHandler struct {
	flagType types.BuiltinFlagType
}

// NewCustomHandler 创建自定义处理器
func NewCustomHandler(flagType types.BuiltinFlagType) *CustomHandler {
	return &CustomHandler{
		flagType: flagType,
	}
}

// ShouldRegister 判断是否应该注册此处理器
func (h *CustomHandler) ShouldRegister(cmd types.Command) bool {
	return true
}

// Handle 处理内置标志
func (h *CustomHandler) Handle(cmd types.Command) error {
	return nil
}

// Type 返回处理器类型
func (h *CustomHandler) Type() types.BuiltinFlagType {
	return h.flagType
}
