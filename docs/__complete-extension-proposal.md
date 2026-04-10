# __complete 子命令扩展方案

> 将补全逻辑进一步整合到 Go 层的完整设计方案

## 1. 设计原则

### 1.1 核心洞察

`__complete` 作为内部隐藏子命令，天然拥有访问父命令所有数据的能力：
- 直接访问 `cmdRegistry` - 子命令注册表
- 直接访问 `flagRegistry` - 标志注册表
- 直接遍历命令树结构
- 无需 Shell 传递命令树数据

### 1.2 设计目标

- **Shell 脚本极简**: 只负责参数解析和调用指令
- **逻辑完全统一**: 所有复杂逻辑在 Go 层实现
- **跨平台一致**: Bash 和 PowerShell 行为完全一致
- **易于扩展**: 后续可轻松添加新的补全指令

---

## 2. 指令体系设计

### 2.1 指令总览

```
__complete <指令> [参数...]

指令列表:
├── fuzzy <pattern> <candidates...>     # 已存在：模糊匹配
├── context <cursorPos> <args...>       # 新增：计算当前上下文
├── candidates <context> [usedFlags...] # 新增：获取候选选项
├── path <wordToComplete> [opts...]     # 新增：路径补全
├── enum <context> <flagName>           # 新增：获取枚举值
└── suggest <context> <word> [--json]   # 新增：智能建议
```

### 2.2 详细指令设计

#### 2.2.1 fuzzy 指令（已存在）

```bash
# 用法
__complete fuzzy <匹配模式> <候选1> [候选2] ...

# 示例
__complete fuzzy "po" "port" "path" "pod" "proxy"

# 输出（每行一个，按匹配质量排序）
port
pod
proxy
path
```

#### 2.2.2 context 指令（新增）

```bash
# 基础用法 - 只返回上下文路径
__complete context <cursorPos> <arg1> [arg2] ...

# 示例
__complete context 2 "myapp" "subcmd" "--help"

# 输出
/subcmd/


# 详细用法 - 返回 JSON 格式详细信息
__complete context --json <cursorPos> <arg1> [arg2] ...

# 示例输出
{
  "context": "/subcmd/",
  "command": "subcmd",
  "parentContext": "/",
  "depth": 1,
  "currentCmd": "subcmd",
  "availableSubCommands": ["child1", "child2"],
  "availableFlags": ["--help", "--verbose", "--output"],
  "isFlagContext": false,
  "flagsStartIndex": -1
}
```

**参数说明：**
- `cursorPos`: 当前光标位置（0-based，0 表示命令名）
- `args`: 完整的命令行参数列表

**实现逻辑：**
1. 从 `cursorPos` 开始向前遍历 `args`
2. 遇到以 `-` 开头的参数，标记 `isFlagContext = true`
3. 在 `cmdRegistry` 中查找子命令
4. 构建上下文路径 `/cmd1/cmd2/`
5. 返回当前上下文信息

#### 2.2.3 candidates 指令（新增）

```bash
# 用法
__complete candidates <context> [wordToComplete] [usedFlags...]

# 示例 - 获取 /subcmd/ 上下文的所有候选
__complete candidates "/subcmd/" "ver" "--output"

# 输出（JSON 格式）
{
  "commands": ["version", "verify"],
  "flags": ["--verbose", "--version"],
  "positional": [],
  "all": ["version", "verify", "--verbose", "--version"]
}
```

**功能说明：**
- 根据上下文返回所有可用的补全候选
- 自动过滤已使用的标志
- 支持前缀过滤（`wordToComplete`）

#### 2.2.4 path 指令（新增）

```bash
# 用法
__complete path <wordToComplete> [--type=<file|dir|all>]

# 示例
__complete path "/home/user/"
__complete path "./src/" --type=file

# 输出（每行一个路径）
/home/user/documents/
/home/user/downloads/
/home/user/file.txt
```

**功能说明：**
- 统一的跨平台路径补全
- 支持文件/目录过滤
- 自动处理路径分隔符（Windows `/` 和 `\`）

#### 2.2.5 enum 指令（新增）

```bash
# 用法
__complete enum <context> <flagName>

# 示例
__complete enum "/" "--output"

# 输出（每行一个枚举值）
json
yaml
table
```

**功能说明：**
- 返回指定标志的所有枚举值
- 支持动态枚举（基于其他标志的值）

#### 2.2.6 suggest 指令（新增）

```bash
# 用法
__complete suggest <context> <wordToComplete> [--json]

# 示例
__complete suggest "/subcmd/" "ver" --json

# 输出
{
  "matches": [
    {"value": "version", "type": "command", "description": "Show version information"},
    {"value": "--verbose", "type": "flag", "description": "Enable verbose output"},
    {"value": "--version", "type": "flag", "description": "Show version"}
  ],
  "context": "/subcmd/"
}
```

**功能说明：**
- 返回带描述信息的智能建议
- 自动分类（命令、标志、枚举值等）

---

## 3. Go 层实现设计

### 3.1 文件结构

```
internal/completion/
├── completion.go              # 主入口和脚本生成
├── bash_completion.go         # Bash 脚本生成
├── pwsh_completion.go         # PowerShell 脚本生成
├── dynamic.go                 # __complete 子命令路由
├── fuzzy.go                   # 模糊匹配实现
├── context.go                 # 新增：上下文计算
├── candidates.go              # 新增：候选选项获取
├── path.go                    # 新增：路径补全
├── enum.go                    # 新增：枚举值获取
├── suggest.go                 # 新增：智能建议
└── templates/
    ├── bash_dynamic.tmpl      # 动态 Bash 模板（简化版）
    └── pwsh_dynamic.tmpl      # 动态 PowerShell 模板（简化版）
```

### 3.2 核心实现

#### 3.2.1 指令路由（dynamic.go）

```go
package completion

import (
	"fmt"
	"gitee.com/MM-Q/qflag/internal/cmd"
)

// 指令常量
const (
	InstructionFuzzy     = "fuzzy"
	InstructionContext   = "context"
	InstructionCandidates = "candidates"
	InstructionPath      = "path"
	InstructionEnum      = "enum"
	InstructionSuggest   = "suggest"
)

// HandleDynamicComplete 处理 __complete 子命令
// 参数:
//   - root: 根命令实例，用于访问注册表
//   - instruction: 指令名称
//   - params: 指令参数列表
// 返回值:
//   - error: 处理错误
func HandleDynamicComplete(root *cmd.Cmd, instruction string, params []string) error {
	switch instruction {
	case InstructionFuzzy:
		return handleFuzzy(params)
	case InstructionContext:
		return handleContext(root, params)
	case InstructionCandidates:
		return handleCandidates(root, params)
	case InstructionPath:
		return handlePath(params)
	case InstructionEnum:
		return handleEnum(root, params)
	case InstructionSuggest:
		return handleSuggest(root, params)
	default:
		return fmt.Errorf("未知指令: %s", instruction)
	}
}
```

#### 3.2.2 上下文计算（context.go）

```go
package completion

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitee.com/MM-Q/qflag/internal/cmd"
)

// ContextInfo 上下文信息结构
type ContextInfo struct {
	Context              string   `json:"context"`
	Command              string   `json:"command"`
	ParentContext        string   `json:"parentContext"`
	Depth                int      `json:"depth"`
	CurrentCmd           string   `json:"currentCmd"`
	AvailableSubCommands []string `json:"availableSubCommands"`
	AvailableFlags       []string `json:"availableFlags"`
	IsFlagContext        bool     `json:"isFlagContext"`
	FlagsStartIndex      int      `json:"flagsStartIndex"`
}

// handleContext 处理 context 指令
// 参数:
//   - root: 根命令实例
//   - args: [cursorPos, arg1, arg2, ...]
// 返回值:
//   - error: 处理错误
func handleContext(root *cmd.Cmd, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: __complete context [--json] <cursorPos> [arg1] [arg2] ...")
	}

	// 解析参数
	jsonOutput := false
	argOffset := 0

	if args[0] == "--json" {
		jsonOutput = true
		argOffset = 1
	}

	if len(args) <= argOffset {
		return fmt.Errorf("缺少 cursorPos 参数")
	}

	cursorPos, err := strconv.Atoi(args[argOffset])
	if err != nil {
		return fmt.Errorf("cursorPos 必须是整数: %v", err)
	}

	tokens := args[argOffset+1:]

	// 计算上下文
	info := calculateContext(root, tokens, cursorPos)

	// 输出结果
	if jsonOutput {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println(info.Context)
	}

	return nil
}

// calculateContext 计算当前上下文
func calculateContext(root *cmd.Cmd, tokens []string, cursorPos int) *ContextInfo {
	info := &ContextInfo{
		Context:         "/",
		ParentContext:   "",
		Depth:           0,
		CurrentCmd:      root.Name(),
		IsFlagContext:   false,
		FlagsStartIndex: -1,
	}

	currentCmd := root

	// 遍历 tokens 直到光标位置
	for i := 0; i < cursorPos && i < len(tokens); i++ {
		token := tokens[i]

		// 遇到标志，记录位置并停止
		if strings.HasPrefix(token, "-") {
			info.IsFlagContext = true
			info.FlagsStartIndex = i
			break
		}

		// 查找子命令
		subCmd, err := currentCmd.LookupSubCmd(token)
		if err != nil {
			// 不是子命令，可能是位置参数
			break
		}

		// 更新上下文
		info.ParentContext = info.Context
		info.Context += token + "/"
		info.Depth++
		info.CurrentCmd = token
		currentCmd = subCmd
	}

	// 获取可用的子命令
	info.AvailableSubCommands = getSubCommandNames(currentCmd)

	// 获取可用的标志
	info.AvailableFlags = getFlagNames(currentCmd)

	return info
}

// getSubCommandNames 获取子命令名称列表
func getSubCommandNames(c *cmd.Cmd) []string {
	names := []string{}
	// 通过 cmdRegistry 获取所有子命令
	// 具体实现取决于 cmd.Cmd 的 API
	return names
}

// getFlagNames 获取标志名称列表
func getFlagNames(c *cmd.Cmd) []string {
	names := []string{}
	// 通过 flagRegistry 获取所有标志
	// 具体实现取决于 cmd.Cmd 的 API
	return names
}
```

#### 3.2.3 候选选项获取（candidates.go）

```go
package completion

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/cmd"
)

// CandidatesInfo 候选选项信息
type CandidatesInfo struct {
	Commands   []string `json:"commands"`
	Flags      []string `json:"flags"`
	Positional []string `json:"positional"`
	All        []string `json:"all"`
}

// handleCandidates 处理 candidates 指令
// 参数:
//   - root: 根命令实例
//   - args: [context, wordToComplete, usedFlag1, usedFlag2, ...]
func handleCandidates(root *cmd.Cmd, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: __complete candidates <context> [wordToComplete] [usedFlags...]")
	}

	context := args[0]
	wordToComplete := ""
	usedFlags := []string{}

	if len(args) > 1 {
		wordToComplete = args[1]
	}
	if len(args) > 2 {
		usedFlags = args[2:]
	}

	// 根据上下文找到对应命令
	currentCmd := findCommandByContext(root, context)
	if currentCmd == nil {
		return fmt.Errorf("无效的上下文: %s", context)
	}

	info := &CandidatesInfo{
		Commands:   []string{},
		Flags:      []string{},
		Positional: []string{},
		All:        []string{},
	}

	// 获取子命令
	for _, name := range getSubCommandNames(currentCmd) {
		if strings.HasPrefix(name, wordToComplete) {
			info.Commands = append(info.Commands, name)
			info.All = append(info.All, name)
		}
	}

	// 获取标志（过滤已使用的）
	usedFlagsSet := make(map[string]bool)
	for _, f := range usedFlags {
		usedFlagsSet[f] = true
	}

	for _, name := range getFlagNames(currentCmd) {
		if !usedFlagsSet[name] && strings.HasPrefix(name, wordToComplete) {
			info.Flags = append(info.Flags, name)
			info.All = append(info.All, name)
		}
	}

	// 输出 JSON
	data, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(data))

	return nil
}

// findCommandByContext 根据上下文路径查找命令
func findCommandByContext(root *cmd.Cmd, context string) *cmd.Cmd {
	if context == "/" {
		return root
	}

	parts := strings.Split(strings.Trim(context, "/"), "/")
	current := root

	for _, part := range parts {
		subCmd, err := current.LookupSubCmd(part)
		if err != nil {
			return nil
		}
		current = subCmd
	}

	return current
}
```

#### 3.2.4 路径补全（path.go）

```go
package completion

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// handlePath 处理 path 指令
// 参数:
//   - args: [wordToComplete, --type=file|dir|all]
func handlePath(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: __complete path <wordToComplete> [--type=<file|dir|all>]")
	}

	wordToComplete := args[0]
	pathType := "all" // file, dir, all

	// 解析可选参数
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--type=") {
			pathType = strings.TrimPrefix(arg, "--type=")
		}
	}

	// 获取目录和文件
	dir := filepath.Dir(wordToComplete)
	if dir == "." {
		dir = ""
	}

	prefix := filepath.Base(wordToComplete)

	// 列出目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil // 静默处理错误
	}

	for _, entry := range entries {
		name := entry.Name()

		// 前缀匹配
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		// 类型过滤
		isDir := entry.IsDir()
		if pathType == "file" && isDir {
			continue
		}
		if pathType == "dir" && !isDir {
			continue
		}

		// 构建完整路径
		fullPath := filepath.Join(dir, name)
		if isDir {
			fullPath += string(filepath.Separator)
		}

		fmt.Println(fullPath)
	}

	return nil
}
```

#### 3.2.5 枚举值获取（enum.go）

```go
package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/cmd"
)

// handleEnum 处理 enum 指令
// 参数:
//   - root: 根命令实例
//   - args: [context, flagName]
func handleEnum(root *cmd.Cmd, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: __complete enum <context> <flagName>")
	}

	context := args[0]
	flagName := args[1]

	// 根据上下文找到对应命令
	currentCmd := findCommandByContext(root, context)
	if currentCmd == nil {
		return fmt.Errorf("无效的上下文: %s", context)
	}

	// 查找标志
	flag, err := currentCmd.LookupFlag(flagName)
	if err != nil {
		return fmt.Errorf("标志不存在: %s", flagName)
	}

	// 获取枚举值
	// 具体实现取决于 flag 的 API
	// 假设 flag 有 GetEnumOptions() 方法
	if enumFlag, ok := flag.(EnumFlag); ok {
		for _, option := range enumFlag.GetEnumOptions() {
			fmt.Println(option)
		}
	}

	return nil
}
```

#### 3.2.6 智能建议（suggest.go）

```go
package completion

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/cmd"
)

// Suggestion 单个建议项
type Suggestion struct {
	Value       string `json:"value"`
	Type        string `json:"type"` // command, flag, enum, file
	Description string `json:"description"`
}

// SuggestInfo 建议信息
type SuggestInfo struct {
	Matches []Suggestion `json:"matches"`
	Context string       `json:"context"`
}

// handleSuggest 处理 suggest 指令
// 参数:
//   - root: 根命令实例
//   - args: [context, wordToComplete, --json]
func handleSuggest(root *cmd.Cmd, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: __complete suggest <context> <wordToComplete> [--json]")
	}

	context := args[0]
	wordToComplete := args[1]
	jsonOutput := false

	// 检查是否有 --json 参数
	for _, arg := range args[2:] {
		if arg == "--json" {
			jsonOutput = true
		}
	}

	// 根据上下文找到对应命令
	currentCmd := findCommandByContext(root, context)
	if currentCmd == nil {
		return fmt.Errorf("无效的上下文: %s", context)
	}

	info := &SuggestInfo{
		Matches: []Suggestion{},
		Context: context,
	}

	// 收集匹配的建议
	// 1. 子命令
	for _, name := range getSubCommandNames(currentCmd) {
		if strings.HasPrefix(name, wordToComplete) {
			info.Matches = append(info.Matches, Suggestion{
				Value:       name,
				Type:        "command",
				Description: getCommandDescription(currentCmd, name),
			})
		}
	}

	// 2. 标志
	for _, name := range getFlagNames(currentCmd) {
		if strings.HasPrefix(name, wordToComplete) {
			info.Matches = append(info.Matches, Suggestion{
				Value:       name,
				Type:        "flag",
				Description: getFlagDescription(currentCmd, name),
			})
		}
	}

	// 输出
	if jsonOutput {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
	} else {
		for _, match := range info.Matches {
			fmt.Println(match.Value)
		}
	}

	return nil
}

// getCommandDescription 获取命令描述
func getCommandDescription(c *cmd.Cmd, name string) string {
	// 具体实现取决于 cmd.Cmd 的 API
	return ""
}

// getFlagDescription 获取标志描述
func getFlagDescription(c *cmd.Cmd, name string) string {
	// 具体实现取决于 cmd.Cmd 的 API
	return ""
}
```

---

## 4. Shell 模板简化方案

### 4.1 简化后的 PowerShell 模板

```powershell
# 智能补全主函数 - 调用 Go 层的 __complete 指令
$scriptBlock = {
    param($wordToComplete, $commandAst, $cursorPosition)

    try {
        # 解析 tokens
        $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
        if (-not $tokens -or $tokens.Count -eq 0) { return @() }

        $currentIndex = $tokens.Count - 1

        # 1. 路径补全快速路径
        if ($wordToComplete -match '[/\~\.]' -or $wordToComplete -like './*') {
            return & {{.ProgramName}} __complete path "$wordToComplete"
        }

        # 2. 计算上下文（调用 Go 层）
        $contextInfo = & {{.ProgramName}} __complete context --json $currentIndex @tokens | ConvertFrom-Json
        $context = $contextInfo.context

        # 3. 获取候选选项（调用 Go 层）
        $candidates = & {{.ProgramName}} __complete candidates $context $wordToComplete | ConvertFrom-Json

        # 4. 模糊匹配（调用 Go 层）
        if ($candidates.all.Count -gt 0) {
            $matches = & {{.ProgramName}} __complete fuzzy $wordToComplete $candidates.all
            return $matches
        }

        return @()
    }
    catch {
        Write-Debug "补全错误: $($_.Exception.Message)"
        return @()
    }
}

# 注册补全
Register-ArgumentCompleter -CommandName {{.ProgramName}} -ScriptBlock $scriptBlock
```

### 4.2 简化后的 Bash 模板

```bash
# 主补全函数
_{{.ProgramName}}_complete() {
    local cur prev words cword
    COMPREPLY=()

    # 获取补全参数
    if declare -F _get_comp_words_by_ref >/dev/null 2>&1; then
        _get_comp_words_by_ref -n =: cur prev words cword
    else
        words=("${COMP_WORDS[@]}")
        cword=$COMP_CWORD
        cur="${words[cword]}"
    fi

    # 输入验证
    [[ $cword -lt 0 || ${#words[@]} -eq 0 ]] && return 1

    # 1. 路径补全快速路径
    if [[ "$cur" == *"/"* || "$cur" == *"."* ]]; then
        COMPREPLY=($({{.ProgramName}} __complete path "$cur"))
        return 0
    fi

    # 2. 计算上下文（调用 Go 层）
    local context
    context=$({{.ProgramName}} __complete context "$cword" "${words[@]}")

    # 3. 获取候选选项（调用 Go 层）
    local candidates_json
    candidates_json=$({{.ProgramName}} __complete candidates "$context" "$cur")

    # 解析候选（简化处理，实际可用 jq）
    local candidates
    candidates=$(echo "$candidates_json" | grep -o '"all": \[[^]]*\]' | sed 's/.*\[\(.*\)\].*/\1/' | tr ',' '\n' | tr -d '" ')

    # 4. 模糊匹配（调用 Go 层）
    if [[ -n "$candidates" ]]; then
        COMPREPLY=($({{.ProgramName}} __complete fuzzy "$cur" $candidates))
    fi

    return 0
}

# 注册补全
complete -F _{{.ProgramName}}_complete {{.ProgramName}}
```

---

## 5. 实现步骤

### 阶段一：核心指令实现

1. [ ] 创建 `internal/completion/context.go`
2. [ ] 创建 `internal/completion/candidates.go`
3. [ ] 创建 `internal/completion/path.go`
4. [ ] 修改 `internal/completion/dynamic.go` 添加路由
5. [ ] 修改 `internal/cmd/builtin.go` 传递 root 命令

### 阶段二：模板更新

1. [ ] 创建简化的 `bash_dynamic_v2.tmpl`
2. [ ] 创建简化的 `pwsh_dynamic_v2.tmpl`
3. [ ] 更新 `completion.go` 支持新模板

### 阶段三：测试验证

1. [ ] 单元测试：各指令功能测试
2. [ ] 集成测试：完整补全流程测试
3. [ ] 性能测试：确保调用延迟 < 10ms

---

## 6. 预期效果

### 6.1 代码量对比

| 组件 | 当前（Shell 实现） | 新方案（Go 实现） | 减少 |
|------|-------------------|------------------|------|
| Bash 模板 | ~150 行 | ~30 行 | 80% |
| PowerShell 模板 | ~250 行 | ~40 行 | 84% |
| Go 代码 | ~100 行（fuzzy） | ~400 行（全部） | - |
| **总计** | ~500 行 | ~470 行 | 逻辑更集中 |

### 6.2 优势总结

1. **逻辑统一**: 所有复杂逻辑在 Go 层实现，跨平台行为一致
2. **易于测试**: 可以为每个指令编写单元测试
3. **易于扩展**: 添加新指令只需修改 Go 代码
4. **Shell 极简**: Shell 脚本只负责参数传递和调用
5. **功能增强**: 支持更复杂的补全场景（标志依赖、动态子命令等）

---

## 7. 附录

### 7.1 指令速查表

| 指令 | 用途 | 复杂度 |
|------|------|--------|
| `fuzzy` | 模糊匹配 | 已实现 |
| `context` | 计算上下文 | 中 |
| `candidates` | 获取候选 | 中 |
| `path` | 路径补全 | 低 |
| `enum` | 枚举值 | 低 |
| `suggest` | 智能建议 | 高 |

### 7.2 依赖 API 清单

项目已提供的 API 完全满足需求：

#### 7.2.1 API 对照表

| 方案中需要的 API | 项目现有 API | 状态 |
|-----------------|-------------|------|
| `LookupSubCmd(name)` | `GetSubCmd(name) (types.Command, bool)` | ✅ 完全匹配 |
| `GetSubCmds()` | `SubCmds() []types.Command` | ✅ 完全匹配 |
| `GetSubCmdNames()` | 需简单封装 | 🟡 基于 `SubCmds()` 实现 |
| `LookupFlag(name)` | `GetFlag(name) (types.Flag, bool)` | ✅ 完全匹配 |
| `GetFlags()` | `Flags() []types.Flag` | ✅ 完全匹配 |
| `GetFlagNames()` | 需简单封装 | 🟡 基于 `Flags()` 实现 |
| `GetDescription()` | `Desc() string` | ✅ 完全匹配 |

#### 7.2.2 需要封装的辅助函数

```go
// 获取子命令名称列表
func getSubCommandNames(c types.Command) []string {
    cmds := c.SubCmds()
    names := make([]string, len(cmds))
    for i, cmd := range cmds {
        names[i] = cmd.Name()
    }
    return names
}

// 获取标志名称列表（包括长短名称）
func getFlagNames(c types.Command) []string {
    flags := c.Flags()
    names := []string{}
    for _, flag := range flags {
        names = append(names, flag.LongName())
        if flag.ShortName() != "" {
            names = append(names, "-"+flag.ShortName())
        }
    }
    return names
}

// 根据上下文路径查找命令
func findCommandByContext(root types.Command, context string) types.Command {
    if context == "/" {
        return root
    }
    
    parts := strings.Split(strings.Trim(context, "/"), "/")
    current := root
    
    for _, part := range parts {
        subCmd, found := current.GetSubCmd(part)
        if !found {
            return nil
        }
        current = subCmd
    }
    
    return current
}
```

#### 7.2.3 types.Flag 接口方法（可直接使用）

```go
// 名称相关
LongName() string          // 长名称，如 "--output"
ShortName() string         // 短名称，如 "o"，可能为空
Name() string              // 与 LongName() 相同

// 描述
Desc() string              // 标志描述

// 类型和枚举
Type() FlagType            // 标志类型
EnumValues() []string      // 枚举值列表（如果是枚举类型）
```

#### 7.2.4 types.Command 接口方法（可直接使用）

```go
// 子命令管理
GetSubCmd(name string) (Command, bool)  // 获取子命令
SubCmds() []Command                      // 获取所有子命令
HasSubCmd(name string) bool              // 检查子命令是否存在

// 标志管理
GetFlag(name string) (Flag, bool)        // 获取标志
Flags() []Flag                           // 获取所有标志

// 基本属性
Name() string                            // 命令名称
Desc() string                            // 命令描述
```

---

## 8. 接口使用示例

### 8.1 遍历命令树获取所有子命令

```go
func collectAllSubCommands(cmd types.Command, prefix string) []string {
    var result []string
    
    subCmds := cmd.SubCmds()
    for _, subCmd := range subCmds {
        fullName := prefix + subCmd.Name()
        result = append(result, fullName)
        
        // 递归收集子命令的子命令
        children := collectAllSubCommands(subCmd, fullName+" ")
        result = append(result, children...)
    }
    
    return result
}
```

### 8.2 获取命令的所有可用标志

```go
func getAllFlags(cmd types.Command) map[string]types.Flag {
    flags := make(map[string]types.Flag)
    
    for _, flag := range cmd.Flags() {
        flags[flag.LongName()] = flag
        if flag.ShortName() != "" {
            flags["-"+flag.ShortName()] = flag
        }
    }
    
    return flags
}
```

### 8.3 检查标志是否为枚举类型并获取可选值

```go
func getEnumOptions(cmd types.Command, flagName string) ([]string, error) {
    flag, found := cmd.GetFlag(flagName)
    if !found {
        return nil, fmt.Errorf("标志不存在: %s", flagName)
    }
    
    if flag.Type() != types.FlagTypeEnum {
        return nil, fmt.Errorf("标志不是枚举类型: %s", flagName)
    }
    
    return flag.EnumValues(), nil
}
```

---

*文档版本: 1.1*
*更新日期: 2026-04-10*
*更新内容: 补充 API 对照表和接口使用示例*
