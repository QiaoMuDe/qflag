# EnumHelp 文档更新计划

## Summary

本项目新增了一个公开函数 `EnumHelp(desc string, options []string, indent string) string`（位于根包 `gitee.com/MM-Q/qflag`，`enum_help.go`），用于生成枚举标志的帮助描述文本，支持 `值` 与 `值: 描述` 两种形态、CJK 双宽对齐、空值选项过滤。

本计划梳理并拟定同时更新 **4 处文档**，使该新函数在公开 API 文档、项目入口文档、AI 记忆点、CLI 技能包中均得到一致收录。`_examples/enum-help/main.go` 示例已存在，无需动作。

## Current State Analysis（现状）

- **根 `APIDOC.md`**：手工维护的 godoc 风格 Markdown，分「变量 / 函数 / 类型」三节。「函数」节条目按**字母序**排列（AddSubCmdFrom → AddSubCmds → ApplyOpts → Parse → ParseAndRoute …）。目前没有 EnumHelp 条目。插入点：`### ApplyOpts`（L171）段之后、`### Parse`（L192）之前。
- **`README.md`**：面向用户的特性 + 示例文档，未逐函数列举，但含「📚 API 文档概述」章节（L367）指向 `APIDOC.md` 与 `_examples/`。适合在此处加一句提点 EnumHelp 及示例。
- **`AGENTS.md`**：AI agent 项目分析报告，含「记忆点」维护规范（编号 1 最旧 → 10 最新，满 10 条时删除最旧→顺移→末尾追加）。当前**仅有 1 条记忆点**（初始架构基线），未达上限。
- **`qflag-cli/references/FLAG_USAGE.md`**：CLI 技能包内的标志使用语法指南，含「### 6. 枚举标志（Enum）」章节（L222–L236），可在此补充 EnumHelp 用法说明。
- **`_examples/enum-help/main.go`**：已创建并能运行（上一轮已实现），不作为本次改动对象。

## Proposed Changes（拟定改动）

### 1. 根 `APIDOC.md` — 新增 EnumHelp 函数条目【必改】

- **位置**：在「函数」表字母序的 `### Parse`（L192）之前，即 `### ApplyOpts` 条目（L171–L191）之后插入。
- **内容**：严格仿照现有条目格式（`### 函数名` / ``` `func` 签名``` / 中文说明 / `**参数:**` / `**返回值:**` / `**功能说明:**`），编写：

```
### EnumHelp

func EnumHelp(desc string, options []string, indent string) string

生成 enum 标志完整的帮助描述文本

**参数:**
- `desc`: 选项开头的描述文案（结尾无需换行，函数会自动追加）
- `options`: 选项列表，每项为 "值" 或 "值: 描述"
- `indent`: 每行选项的缩进前缀（如 "\t"），保持项目统一风格

**返回值:**
- `string`: 拼好的完整 enum desc 字符串

**功能说明:**
- 支持 "值" → "[值]" 与 "值: 描述" → "[值]   - 描述" 两种形态
- 按终端显示宽度对齐（全角/CJK 占 2 列，其余占 1 列）
- 自动跳过空值选项，避免输出空盒子
```

### 2. `README.md` — 提及 EnumHelp 与示例【必改（用户选全量）】

- **位置**：「📚 API 文档概述」章节（L367 附近）或紧随其后的「更多文档」列表。
- **内容**：新增一句/一条，说明 EnumHelp 是枚举帮助文本生成辅助函数，并指向 `_examples/enum-help`。例如在 API 文档概述段落末尾补充：

> 枚举标志帮助文本可使用 `EnumHelp` 便捷生成，示例见 [`_examples/enum-help/`](_examples/enum-help/)。

- **注意**：README 现状不逐函数列举，因此**不做**完整 API 条目，仅提点入口，避免风格偏差。

### 3. `AGENTS.md` — 追加记忆点【必改（用户选全量）】

- **做法**：当前仅 1 条记忆点，未达 10 条上限，故**在末尾追加**「记忆点 2」，不触发滚动删除。内容记录 EnumHelp 的引入与行为要点。
- **拟定内容**：

> **记忆点 2**（新增 EnumHelp）：
> - 新增公开辅助函数 `EnumHelp`（根包 `enum_help.go`），生成枚举标志帮助文本，支持 `值` 与 `值: 描述` 两种形态。
> - 对齐按终端显示宽度计算（CJK/全角计 2 列，`isWideRune`），修复中英混排对齐失真。
> - 自动跳过空值选项；解析约定：值本身不能含冒号（冒号仅作值与描述分隔），首个冒号拆分。
> - 配套测试位于 `enum_help_test.go`，示例位于 `_examples/enum-help/`。

- **注意**：为保留重要的「初始架构基线」记忆点，本次用「追加」而非「删除」，符合未满上限时的维护语义；若后续超过 10 条再按滚动规则处理。

### 4. `qflag-cli/references/FLAG_USAGE.md` — 枚举章节补充【必改（用户选全量）】

- **位置**：「### 6. 枚举标志（Enum）」章节末尾（L236 下方、L238 `## 标志定义示例` 之前）。
- **内容**：补充一段 EnumHelp 用法说明（面向用 qflag 开发 CLI 的工程师），例如：

```
枚举标志帮助文本可通过 `EnumHelp` 便捷生成：

qflag.EnumHelp("运行模式", []string{"debug: 调试模式", "release: 发布模式"}, "\t")
```

可注明：对齐按显示宽度（CJK 双宽）、空值选项自动跳过、值为空白的选项被过滤。

## Assumptions & Decisions（假设与决策）

- **范围**：按用户选择采用全量（APIDOC + README + AGENTS + qflag-cli），`_examples` 已具备、不重复创建。
- **APIDOC 插入点**：以字母序 `ApplyOpts` 与 `Parse` 之间为准，与现有排序一致性。
- **AGENTS.md 追加策略**：因未满上限，采用追加「记忆点 2」，避免误删重要的初始基线记忆点。
- **doc2md**：经探查未接入本项目文档生成，APIDOC 为纯手工维护，故本次为直接编辑。
- **CHANGELOG**：项目中无集中 CHANGELOG 文件，无需维护。

## Verification（验证）

1. `go build ./...` 通过（文档改动不影响编译逻辑，但确认无破坏）。
2. `go test ./...` 全量通过（enum_help 相关测试仍为绿）。
3. 打开 `APIDOC.md` 确认 `EnumHelp` 条目出现在字母序正确位置、格式与邻近条目一致。
4. 打开 `README.md`、`AGENTS.md`、`qflag-cli/references/FLAG_USAGE.md` 核对新增内容与排版。
5. `go run ./_examples/enum-help` 可正常运行（示例回归确认）。