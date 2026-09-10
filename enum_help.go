// Package qflag 提供跨模块共享的通用工具函数。
// 本文件包含枚举标志帮助文本生成的便捷函数。
package qflag

import (
	"strings"
)

// EnumHelp 生成 enum 标志完整的帮助描述文本。
// 支持两种形态的选项：
//   - 无描述：  "值"       → "[值]"
//   - 有描述：  "值: 描述" → "[值]   - 描述"
//
// 解析约定:
//   - 值（value）本身不能包含冒号，冒号仅用作值与描述的分隔符（取第一个冒号拆分）
//   - 值为空白的选项（纯空格或冒号开头）会被忽略，避免产生空盒子
//   - 对齐按终端显示宽度计算：全角/CJK 字符占 2 列，其余占 1 列
//
// 参数:
//   - desc:   选项开头的描述文案（结尾无需换行，函数会自动追加）
//   - options: 选项列表，每项为 "值" 或 "值: 描述"
//   - indent:  每行选项的缩进前缀（如 "\t\t\t\t\t"），保持项目统一风格
//
// 返回:
//   - string: 拼好的完整 enum desc 字符串（desc + 每行缩进与换行的选项列表）
func EnumHelp(desc string, options []string, indent string) string {
	if len(options) == 0 {
		return desc
	}

	// 逐项解析出"值"与"描述"，并计算方框内"值"的最大显示宽度
	values := make([]string, 0, len(options))
	descs := make([]string, 0, len(options))
	maxWidth := 0
	for _, opt := range options {
		value := opt
		d := ""
		if idx := strings.Index(opt, ":"); idx >= 0 {
			value = strings.TrimSpace(opt[:idx])
			d = strings.TrimSpace(opt[idx+1:])
		}
		// 跳过空值选项（纯空白或冒号开头），避免输出空盒子污染对齐
		if strings.TrimSpace(value) == "" {
			continue
		}
		values = append(values, value)
		descs = append(descs, d)
		if w := runeDisplayWidth(value); w > maxWidth {
			maxWidth = w
		}
	}

	// 所有选项均为空值时，无内容可展示，直接返回 desc
	if len(values) == 0 {
		return desc
	}

	// 拼接描述与逐行选项：缩进 + 方框 + 可选描述（末行不加换行，避免 default 值另起一行）
	var buf strings.Builder
	buf.WriteString(desc)
	buf.WriteString("\n")
	for i, v := range values {
		buf.WriteString(indent)
		buf.WriteString("[")
		buf.WriteString(v)
		buf.WriteString(strings.Repeat(" ", maxWidth-runeDisplayWidth(v)))
		buf.WriteString("]")
		if descs[i] != "" {
			buf.WriteString("   - ")
			buf.WriteString(descs[i])
		}
		if i != len(values)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// runeDisplayWidth 计算字符串在等宽终端中的显示宽度。
// 全角/CJK 字符计为 2 列，其余计为 1 列，用于帮助文本的列对齐。
func runeDisplayWidth(s string) int {
	w := 0
	for _, r := range s {
		if isWideRune(r) {
			w += 2
		} else {
			w++
		}
	}
	return w
}

// isWideRune 判断该字符是否为宽字符（在等宽终端中占用 2 列）。
// 覆盖 CJK 统一表意文字、全角形式、谚文音节等常见宽字符区间。
func isWideRune(r rune) bool {
	switch {
	// 谚文 Jamo / 谚文音节
	case r >= 0x1100 && r <= 0x115F,
		r >= 0xAC00 && r <= 0xD7A3:
		return true
	// CJK 统一表意文字及扩展（含部首、拼音注、谚文兼容等）
	case r >= 0x2E80 && r <= 0x303E,
		r >= 0x3041 && r <= 0x33FF,
		r >= 0x3400 && r <= 0x4DBF,
		r >= 0x4E00 && r <= 0x9FFF,
		r >= 0xA000 && r <= 0xA4CF,
		r >= 0xF900 && r <= 0xFAFF,
		r >= 0x20000 && r <= 0x2FFFD,
		r >= 0x30000 && r <= 0x3FFFD:
		return true
	// 竖排标点、CJK 兼容标点
	case r >= 0xFE10 && r <= 0xFE19,
		r >= 0xFE30 && r <= 0xFE6F:
		return true
	// 全角形式（含全角 ASCII，如 ０-９、Ａ-Ｚ）
	case r >= 0xFF00 && r <= 0xFF60,
		r >= 0xFFE0 && r <= 0xFFE6:
		return true
	}
	return false
}
