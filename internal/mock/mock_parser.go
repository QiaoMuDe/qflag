package mock

import (
	"gitee.com/MM-Q/qflag/internal/types"
)

// MockParser 模拟解析器实现
type MockParser struct {
	parseOnlyCalled     bool
	parseCalled         bool
	parseAndRouteCalled bool
	lastCmd             types.Command
	lastArgs            []string
	parseError          error
	routeError          error
}

// NewMockParser 创建新的模拟解析器
func NewMockParser() *MockParser {
	return &MockParser{
		parseError: nil,
		routeError: nil,
	}
}

// NewMockParserWithError 创建带有错误的模拟解析器
func NewMockParserWithError(parseErr, routeErr error) *MockParser {
	return &MockParser{
		parseError: parseErr,
		routeError: routeErr,
	}
}

// 实现 Parser 接口
func (p *MockParser) ParseOnly(cmd types.Command, args []string) error {
	p.parseOnlyCalled = true
	p.lastCmd = cmd
	p.lastArgs = args

	// 模拟解析过程
	for _, arg := range args {
		if arg == "--error" {
			return p.parseError
		}
	}

	return p.parseError
}

func (p *MockParser) Parse(cmd types.Command, args []string) error {
	p.parseCalled = true
	p.lastCmd = cmd
	p.lastArgs = args

	// 模拟解析过程
	for _, arg := range args {
		if arg == "--error" {
			return p.parseError
		}
	}

	return p.parseError
}

func (p *MockParser) ParseAndRoute(cmd types.Command, args []string) error {
	p.parseAndRouteCalled = true
	p.lastCmd = cmd
	p.lastArgs = args

	// 模拟解析过程
	for _, arg := range args {
		if arg == "--error" {
			return p.parseError
		}
		if arg == "--route-error" {
			return p.routeError
		}
	}

	// 如果有路由错误, 返回路由错误
	if p.routeError != nil {
		return p.routeError
	}

	return p.parseError
}

// 辅助方法, 用于测试验证
func (p *MockParser) WasParseOnlyCalled() bool {
	return p.parseOnlyCalled
}

func (p *MockParser) WasParseCalled() bool {
	return p.parseCalled
}

func (p *MockParser) WasParseAndRouteCalled() bool {
	return p.parseAndRouteCalled
}

func (p *MockParser) GetLastCommand() types.Command {
	return p.lastCmd
}

func (p *MockParser) GetLastArgs() []string {
	return p.lastArgs
}

func (p *MockParser) Reset() {
	p.parseOnlyCalled = false
	p.parseCalled = false
	p.parseAndRouteCalled = false
	p.lastCmd = nil
	p.lastArgs = nil
}
