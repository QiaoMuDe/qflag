package parser

import (
	"errors"
	"testing"

	"gitee.com/MM-Q/qflag/internal/mock"
	"gitee.com/MM-Q/qflag/internal/types"
)

// buildCheckCmd 构造用于测试 checkUnknownFlags 的命令及标志
//
// 返回的命令包含以下标志:
//   - --name / -n  字符串类型（取值）
//   - --port / -p  整数类型（取值）
//   - --ratio / -r 浮点类型（取值）
//   - --num        64位整数类型（取值）
//   - --filter     字符串类型（取值）
//   - --vals       整数切片类型（取值）
//   - --verbose/-v 布尔类型（不取值）
//
// 并注册一个名为 sub 的子命令。
//
// 返回值:
//   - types.Command: 构造好的模拟命令
func buildCheckCmd() types.Command {
	cmd := mock.NewMockCommand("app", "a", "test app")

	flags := []types.Flag{
		mock.NewMockFlag("name", "n", "名称", types.FlagTypeString, ""),
		mock.NewMockFlag("port", "p", "端口", types.FlagTypeInt, 0),
		mock.NewMockFlag("ratio", "r", "比例", types.FlagTypeFloat64, 0.0),
		mock.NewMockFlag("num", "", "大数字", types.FlagTypeInt64, 0),
		mock.NewMockFlag("filter", "", "过滤器", types.FlagTypeString, ""),
		mock.NewMockFlag("vals", "", "数值列表", types.FlagTypeIntSlice, nil),
		mock.NewMockBoolFlag("verbose", "v", "详细输出", false),
	}
	for _, f := range flags {
		if err := cmd.AddFlag(f); err != nil {
			panic(err)
		}
	}

	sub := mock.NewMockSubCommand("sub", "s", "子命令", cmd)
	if err := cmd.AddSubCmds(sub); err != nil {
		panic(err)
	}

	return cmd
}

// TestCheckUnknownFlags 测试 checkUnknownFlags 预扫描逻辑
//
// 覆盖常见的、冷门的、边界的标志解析场景，重点验证:
//   - 非布尔取值标志后跟以 - 开头的值（如负数 -8080）不再被误判为未知标志
//   - 布尔标志不消费下一个参数，其后跟未知标志仍需报错
//   - --flag=value 内联取值（含负值）行为不变
//   - -- 终止符、子命令、未知标志等既有行为保持不变
func TestCheckUnknownFlags(t *testing.T) {
	tests := []struct {
		name           string   // 用例名
		args           []string // 命令行参数
		wantError      bool     // 是否期望返回错误
		wantErrorInput string   // 期望错误中的标志输入（有错误时校验）
	}{
		// ---------- A. 常见（缺陷回归 + 常规） ----------
		{
			name:      "常见:整数空格负值",
			args:      []string{"--port", "-8080"},
			wantError: false,
		},
		{
			name:      "常见:浮点空格负值",
			args:      []string{"--ratio", "-1.5"},
			wantError: false,
		},
		{
			name:      "常见:字符串负号开头值",
			args:      []string{"--name", "-foo"},
			wantError: false,
		},
		{
			name:      "常见:普通正值",
			args:      []string{"--port", "8080"},
			wantError: false,
		},
		{
			name:      "常见:等号形式",
			args:      []string{"--name=hello"},
			wantError: false,
		},
		{
			name:      "常见:布尔单独使用",
			args:      []string{"--verbose"},
			wantError: false,
		},
		// ---------- B. 冷门 / 少见形式 ----------
		{
			name:      "冷门:负值内联",
			args:      []string{"--port=-8080"},
			wantError: false,
		},
		{
			name:      "冷门:负号开头内联",
			args:      []string{"--name=-foo"},
			wantError: false,
		},
		{
			name:      "冷门:短名称取负值",
			args:      []string{"-p", "-8080"},
			wantError: false,
		},
		{
			name:      "冷门:int64 极值负值",
			args:      []string{"--num", "-9223372036854775808"},
			wantError: false,
		},
		{
			name:      "冷门:特殊字符串值",
			args:      []string{"--filter", "-[a-z]+"},
			wantError: false,
		},
		{
			name:      "冷门:切片负值",
			args:      []string{"--vals", "-1,-2"},
			wantError: false,
		},
		{
			name:      "冷门:浮点特殊值 -Inf",
			args:      []string{"--ratio", "-Inf"},
			wantError: false,
		},
		// ---------- C. 边界 / 行为保持 ----------
		{
			name:           "边界:布尔后跟未知标志",
			args:           []string{"--verbose", "-x"},
			wantError:      true,
			wantErrorInput: "-x",
		},
		{
			name:           "边界:消费负值后仍发现未知标志",
			args:           []string{"--port", "-8080", "--unknown"},
			wantError:      true,
			wantErrorInput: "--unknown",
		},
		{
			name:      "边界:取值标志后跟另一合法标志充当值",
			args:      []string{"--port", "--verbose"},
			wantError: false,
		},
		{
			name:      "边界:双横杠终止符后为位置参数",
			args:      []string{"--", "--port"},
			wantError: false,
		},
		{
			name:      "边界:单横杠单独出现视为位置参数",
			args:      []string{"-"},
			wantError: false,
		},
		{
			name:      "边界:布尔后单横杠视为位置参数",
			args:      []string{"--verbose", "-"},
			wantError: false,
		},
		{
			name:      "边界:单横杠作为取值标志的值",
			args:      []string{"--port", "-"},
			wantError: false,
		},
		{
			name:      "边界:双横杠被取值标志消费为值",
			args:      []string{"--port", "--"},
			wantError: false,
		},
		{
			name:      "边界:双横杠作值后继续接位置参数",
			args:      []string{"--port", "--", "8080"},
			wantError: false,
		},
		{
			name:           "边界:双横杠作值后未知标志仍检出",
			args:           []string{"--port", "--", "--unknown"},
			wantError:      true,
			wantErrorInput: "--unknown",
		},
		{
			name:      "边界:遇到子命令名中止扫描",
			args:      []string{"--port", "8080", "sub"},
			wantError: false,
		},
		{
			name:           "边界:纯未知标志回归保护",
			args:           []string{"--nope"},
			wantError:      true,
			wantErrorInput: "--nope",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := buildCheckCmd()
			err := checkUnknownFlags(cmd, tt.args)

			if tt.wantError {
				if err == nil {
					t.Fatalf("期望返回错误，但得到 nil")
				}
				var ufe *types.UnknownFlagError
				if !errors.As(err, &ufe) {
					t.Fatalf("期望 *types.UnknownFlagError，但得到: %v", err)
				}
				if ufe.Input != tt.wantErrorInput {
					t.Fatalf("期望错误输入 '%s'，但得到 '%s'", tt.wantErrorInput, ufe.Input)
				}
				return
			}

			if err != nil {
				t.Fatalf("期望无错误，但得到: %v", err)
			}
		})
	}
}
