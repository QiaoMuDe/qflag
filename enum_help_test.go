package qflag

import (
	"strings"
	"testing"
)

// TestEnumHelp 测试 EnumHelp 函数生成 enum 标志帮助描述文本的逻辑。
//
// 场景覆盖:
//   - 无选项列表（应原样返回 desc）
//   - 无描述的纯值选项 → "[值]"
//   - 带描述的 "值: 描述" 选项 → "[值]   - 描述"
//   - 混合形态（部分带描述、部分不带）
//   - 中英文描述（验证宽度计算按显示宽度而非 rune 数或字节）
//   - 多字节（中文）值等宽对齐（CJK 双宽感知）
//   - 混合 ASCII 与 CJK 值的显示宽度对齐
//   - 空值选项被跳过
//   - 末尾无多余换行（便于 default 值衔接）
func TestEnumHelp(t *testing.T) {
	tests := []struct {
		name    string
		desc    string
		options []string
		indent  string
		want    string
	}{
		{
			name:    "nil 选项列表返回 desc",
			desc:    "运行模式",
			options: nil,
			indent:  "\t\t",
			want:    "运行模式",
		},
		{
			name:    "空选项切片返回 desc",
			desc:    "运行模式",
			options: []string{},
			indent:  "\t\t",
			want:    "运行模式",
		},
		{
			name:    "无描述纯值选项等宽对齐",
			desc:    "运行模式",
			options: []string{"auto", "manual", "debug"},
			indent:  "\t\t\t",
			want: "运行模式\n" +
				"\t\t\t[auto  ]\n" +
				"\t\t\t[manual]\n" +
				"\t\t\t[debug ]",
		},
		{
			name:    "带描述选项并对齐闭合括号",
			desc:    "运行模式",
			options: []string{"auto: 自动模式", "manual: 手动模式", "debug: 调试模式"},
			indent:  "\t\t\t\t\t",
			want: "运行模式\n" +
				"\t\t\t\t\t[auto  ]   - 自动模式\n" +
				"\t\t\t\t\t[manual]   - 手动模式\n" +
				"\t\t\t\t\t[debug ]   - 调试模式",
		},
		{
			name:    "混合形态选项",
			desc:    "日志级别",
			options: []string{"debug", "info: 信息", "warn: 警告"},
			indent:  "\t\t",
			want: "日志级别\n" +
				"\t\t[debug]\n" +
				"\t\t[info ]   - 信息\n" +
				"\t\t[warn ]   - 警告",
		},
		{
			name:    "多字节值等宽对齐",
			desc:    "模式",
			options: []string{"键盘", "鼠标", "键"},
			indent:  "\t",
			want: "模式\n" +
				"\t[键盘]\n" +
				"\t[鼠标]\n" +
				"\t[键  ]",
		},
		{
			name:    "混合 ASCII 与 CJK 值显示宽度对齐",
			desc:    "混合值",
			options: []string{"A", "键盘", "BC"},
			indent:  "\t",
			want: "混合值\n" +
				"\t[A   ]\n" +
				"\t[键盘]\n" +
				"\t[BC  ]",
		},
		{
			name:    "空值选项被跳过",
			desc:    "运行模式",
			options: []string{"auto", "   ", ": desc"},
			indent:  "\t",
			want: "运行模式\n" +
				"\t[auto]",
		},
		{
			name:    "全部选项为空值时返回 desc",
			desc:    "运行模式",
			options: []string{"   ", ": 描述"},
			indent:  "\t",
			want:    "运行模式",
		},
		{
			name:    "英文描述单行",
			desc:    "Mode",
			options: []string{"auto", "manual"},
			indent:  "\t\t",
			want: "Mode\n" +
				"\t\t[auto  ]\n" +
				"\t\t[manual]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnumHelp(tt.desc, tt.options, tt.indent)
			if got != tt.want {
				t.Errorf("EnumHelp() 不匹配:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

// TestEnumHelp_ContainsDesc 验证生成的帮助文本以 desc 开头且包含全部选项及描述。
func TestEnumHelp_ContainsDesc(t *testing.T) {
	desc := "运行模式"
	got := EnumHelp(desc, []string{"auto: 自动模式", "manual: 手动模式"}, "\t\t")

	if !strings.HasPrefix(got, desc) {
		t.Errorf("期望结果以 desc(%q) 开头，实际: %q", desc, got)
	}
	if !strings.Contains(got, "auto") || !strings.Contains(got, "manual") {
		t.Errorf("期望结果包含全部选项值，实际: %q", got)
	}
	if !strings.Contains(got, "自动模式") || !strings.Contains(got, "手动模式") {
		t.Errorf("期望结果包含选项描述，实际: %q", got)
	}
}

// TestEnumHelp_NoTrailingNewline 验证最后一行结尾不追加换行（便于 default 值衔接）。
func TestEnumHelp_NoTrailingNewline(t *testing.T) {
	got := EnumHelp("运行模式", []string{"auto", "manual"}, "\t")
	if strings.HasSuffix(got, "\n") {
		t.Errorf("期望结果末尾无换行，实际以换行结尾: %q", got)
	}
}
