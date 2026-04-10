# __complete context 指令设计方案

> 完全替代静态命令树的动态上下文计算方案

## 1. 设计目标

### 1.1 核心目标

- **消除静态命令树**：不再需要在 Shell 脚本中生成静态的 `cmd_tree` 数组
- **动态查询**：实时通过 Go 层的 `cmdRegistry` 查询命令树
- **支持动态子命令**：可以支持运行时动态注册的子命令
- **简化模板**：Shell 脚本只负责传递参数，所有逻辑在 Go 层实现

### 1.2 预期效果

**当前方案（静态命令树）：**
```bash
# 脚本生成时写入，体积大且固定
declare -A cmd_tree
cmd_tree["/"]="server client --help"
cmd_tree["/server/"]="start stop --port"
cmd_tree["/server/start/"]="--config --daemon"
```

**新方案（动态查询）：**
```bash
# 实时查询，无需静态数据
context=$(myapp __complete context 2 "myapp" "server" "start")
# 输出: /server/start/
```

---

## 2. 指令设计

### 2.1 设计原则

为了**完全替代**原脚本中的静态命令树逻辑，需要新增 **两个指令**：

| 指令 | 功能 | 替代原脚本的哪部分 |
|------|------|------------------|
| `context` | 根据输入参数计算上下文路径 | 替代上下文计算循环 |
| `candidates` | 根据上下文返回可用选项 | 替代查 `cmd_tree` 静态数组 |

### 2.2 context 指令

**用途**：计算当前命令行所处的上下文路径

**用法**：
```bash
__complete context <arg0> <arg1> ... <argN>
```

**参数说明**：
- `arg0`: 程序名（通常是 os.Args[0]）
- `arg1...argN`: 用户输入的参数
- **cursorPos 由 Go 层自动计算**：使用 `len(args)` 作为光标位置

**输出**：
```bash
__complete context "myapp" "server" "start"
# 输出: /server/start/
```

**算法逻辑**：
```
输入: ["myapp", "server", "start", "--port", "8080"]
自动计算 cursorPos = 5

遍历过程:
  i=1: token="server"   → 有效子命令 → 上下文=/server/
  i=2: token="start"    → 有效子命令 → 上下文=/server/start/
  i=3: token="--port"   → 遇到标志   → 停止遍历
  
最终结果: /server/start/
```

### 2.3 candidates 指令

**用途**：根据上下文路径返回该上下文下所有可用的补全选项

**用法**：
```bash
__complete candidates <context>
```

**参数说明**：
- `context`: 上下文路径（由 `context` 指令返回）

**输出**：
```bash
__complete candidates "/server/start/"
# 输出（空格分隔）: start stop --port --host --help
```

**返回内容**：
- 该上下文下的所有子命令名称
- 该上下文下的所有标志名称（包括长短名称）

### 2.4 enum 指令（枚举值补全）

**用途**：根据上下文路径和标志名，返回该枚举标志的所有可选值

**用法**：
```bash
__complete enum <context> <flag-name>
```

**参数说明**：
- `context`: 上下文路径（由 `context` 指令返回）
- `flag-name`: 枚举类型标志的名称（长名称或短名称）

**输出**：
```bash
# 示例：获取 --log-level 标志的枚举值
__complete enum "/server/start/" "--log-level"
# 输出（空格分隔）: debug info warn error fatal

# 示例：获取 -o 标志的枚举值（短名称）
__complete enum "/server/start/" "-o"
# 输出（空格分隔）: json yaml text
```

**使用场景**：
- 当用户输入 `--log-level=` 或 `--log-level ` 时，Shell 脚本调用此指令获取可补全的枚举值
- 替代原脚本中的静态 `enum_options` 数组

### 2.5 三个指令的配合使用

Shell 脚本中的完整流程：

```bash
# 1. 计算上下文路径
context=$({{.ProgramName}} __complete context "${words[@]}")
# 输出: /server/start/

# 2. 检测是否是枚举值补全场景
if [[ "$cur" == --*"="* ]]; then
    # 提取标志名和当前值
    flag_name="${cur%%=*}="
    current_value="${cur#*=}"
    
    # 3a. 获取枚举值
    enum_values=$({{.ProgramName}} __complete enum "$context" "$flag_name")
    COMPREPLY=($(compgen -W "$enum_values" -- "$current_value"))
else
    # 3b. 普通补全：获取候选选项
    candidates=$({{.ProgramName}} __complete candidates "$context")
    COMPREPLY=($(compgen -W "$candidates" -- "$cur"))
fi
```

---

## 3. Go 层实现

### 3.1 文件位置

```
internal/completion/
├── context.go          # 新增：上下文计算指令实现
├── candidates.go       # 新增：候选选项获取指令实现
├── enum.go             # 新增：枚举值补全指令实现
└── context_test.go     # 新增：单元测试
```

### 3.2 核心代码实现

```go
// context.go
package completion

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// ContextResult 上下文计算结果
type ContextResult struct {
	// 基础信息
	Context   string `json:"context"`   // 上下文路径，如 "/server/start/"
	Command   string `json:"command"`   // 当前命令名
	Depth     int    `json:"depth"`     // 嵌套深度

	// 当前命令信息
	CurrentCmd    string   `json:"currentCmd"`    // 当前命令名称
	CurrentDesc   string   `json:"currentDesc"`   // 当前命令描述

	// 可用选项
	SubCommands []string `json:"subCommands"` // 可用子命令列表
	Flags       []string `json:"flags"`       // 可用标志列表（长短名称）

	// 上下文状态
	IsFlagContext   bool `json:"isFlagContext"`   // 是否处于标志上下文
	FlagsStartIndex int  `json:"flagsStartIndex"` // 标志开始的位置（-1 表示无）

	// 父上下文
	ParentContext string `json:"parentContext"` // 父上下文路径
}

// HandleContext 处理 context 指令
// 这是 __complete 子命令的入口点
//
// 参数:
//   - root: 根命令实例，通过 cmdRegistry 查询子命令
//   - args: 命令行参数，格式: [--json] <arg0> <arg1> ...
//
// 返回值:
//   - error: 处理过程中的错误
//
// 示例:
//   HandleContext(root, []string{"myapp", "server", "start"})
//   HandleContext(root, []string{"--json", "myapp", "server", "--port"})
func HandleContext(root types.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: __complete context [--json] <arg0> [arg1] ...")
	}

	// 解析选项
	jsonOutput := false
	argOffset := 0

	if args[0] == "--json" {
		jsonOutput = true
		argOffset = 1
	}

	// 检查参数数量
	if len(args) <= argOffset {
		return fmt.Errorf("缺少参数")
	}

	// 获取 tokens（包括程序名）
	tokens := args[argOffset:]

	// 自动计算 cursorPos：使用 tokens 的长度
	cursorPos := len(tokens)

	// 计算上下文
	result := CalculateContext(root, tokens, cursorPos)

	// 输出结果
	if jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSON 序列化失败: %v", err)
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(result.Context)
	}

	return nil
}

// CalculateContext 计算当前上下文
// 这是核心算法，完全替代 Shell 脚本中的静态命令树查询
//
// 算法逻辑:
//   1. 从索引 1 开始遍历 tokens（跳过程序名）
//   2. 遇到以 "-" 开头的 token，标记为标志上下文并停止
//   3. 在 cmdRegistry 中查找子命令
//   4. 找到则更新上下文，继续遍历
//   5. 未找到则停止遍历，保持当前上下文
//
// 参数:
//   - root: 根命令实例
//   - tokens: 完整的命令行参数列表（包括程序名）
//   - cursorPos: 当前光标位置
//
// 返回值:
//   - *ContextResult: 上下文计算结果
func CalculateContext(root types.Command, tokens []string, cursorPos int) *ContextResult {
	result := &ContextResult{
		Context:         "/",
		Command:         root.Name(),
		Depth:           0,
		CurrentCmd:      root.Name(),
		CurrentDesc:     root.Desc(),
		SubCommands:     []string{},
		Flags:           []string{},
		IsFlagContext:   false,
		FlagsStartIndex: -1,
		ParentContext:   "",
	}

	currentCmd := root

	// 从索引 1 开始遍历（跳过程序名 arg0）
	for i := 1; i < cursorPos && i < len(tokens); i++ {
		token := tokens[i]

		// 规则 1: 遇到标志，停止上下文构建
		if strings.HasPrefix(token, "-") {
			result.IsFlagContext = true
			result.FlagsStartIndex = i
			break
		}

		// 规则 2: 在注册表中查找子命令
		subCmd, found := currentCmd.GetSubCmd(token)
		if !found {
			// 不是有效的子命令，停止遍历
			break
		}

		// 规则 3: 更新上下文
		result.ParentContext = result.Context
		result.Context += token + "/"
		result.Depth++
		result.CurrentCmd = token
		result.CurrentDesc = subCmd.Desc()
		currentCmd = subCmd
	}

	// 获取当前命令的可用选项
	result.SubCommands = getSubCommandNames(currentCmd)
	result.Flags = getFlagNames(currentCmd)

	return result
}

// getSubCommandNames 获取子命令名称列表
//
// 参数:
//   - cmd: 命令实例
//
// 返回值:
//   - []string: 子命令名称列表
func getSubCommandNames(cmd types.Command) []string {
	subCmds := cmd.SubCmds()
	names := make([]string, len(subCmds))
	for i, subCmd := range subCmds {
		names[i] = subCmd.Name()
	}
	return names
}

// getFlagNames 获取标志名称列表（包括长短名称）
//
// 参数:
//   - cmd: 命令实例
//
// 返回值:
//   - []string: 标志名称列表（长名称和短名称）
func getFlagNames(cmd types.Command) []string {
	flags := cmd.Flags()
	names := make([]string, 0, len(flags)*2)

	for _, flag := range flags {
		// 添加长名称
		names = append(names, flag.LongName())

		// 添加短名称（如果有）
		if flag.ShortName() != "" {
			names = append(names, "-"+flag.ShortName())
		}
	}

	return names
}

// findCommandByContext 根据上下文路径查找命令
//
// 参数:
//   - root: 根命令
//   - context: 上下文路径，如 "/server/start/"
//
// 返回值:
//   - types.Command: 找到的命令，如果未找到则返回 nil
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

### 3.3 candidates 指令实现

创建 `internal/completion/candidates.go`：

```go
// candidates.go
package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleCandidates 处理 candidates 指令
//
// 参数:
//   - root: 根命令实例
//   - args: [context]
//
// 返回值:
//   - error: 处理错误
//
// 示例:
//   HandleCandidates(root, []string{"/server/start/"})
func HandleCandidates(root types.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("用法: __complete candidates <context>")
	}

	context := args[0]

	// 根据上下文查找命令
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		// 无效的上下文，返回空
		return nil
	}

	// 收集所有候选选项
	var candidates []string

	// 添加子命令
	candidates = append(candidates, getSubCommandNames(cmd)...)

	// 添加标志
	candidates = append(candidates, getFlagNames(cmd)...)

	// 输出（空格分隔）
	fmt.Println(strings.Join(candidates, " "))

	return nil
}
```

### 3.4 enum 指令实现

创建 `internal/completion/enum.go`：

```go
// enum.go
package completion

import (
	"fmt"
	"strings"

	"gitee.com/MM-Q/qflag/internal/types"
)

// HandleEnum 处理 enum 指令
//
// 参数:
//   - root: 根命令实例
//   - args: [context, flag-name]
//
// 返回值:
//   - error: 处理错误
//
// 示例:
//   HandleEnum(root, []string{"/server/start/", "--log-level"})
//   HandleEnum(root, []string{"/server/start/", "-o"})
func HandleEnum(root types.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: __complete enum <context> <flag-name>")
	}

	context := args[0]
	flagName := args[1]

	// 根据上下文查找命令
	cmd := findCommandByContext(root, context)
	if cmd == nil {
		// 无效的上下文，返回空
		return nil
	}

	// 查找标志
	var targetFlag types.Flag
	for _, flag := range cmd.Flags() {
		// 匹配长名称（去掉可能的 "=" 后缀）
		if flag.LongName() == flagName || flag.LongName()+"=" == flagName {
			targetFlag = flag
			break
		}
		// 匹配短名称（去掉可能的 "=" 后缀）
		shortName := "-" + flag.ShortName()
		if shortName == flagName || shortName+"=" == flagName {
			targetFlag = flag
			break
		}
	}

	if targetFlag == nil {
		// 标志不存在，返回空
		return nil
	}

	// 获取枚举值
	enumValues := getEnumValues(targetFlag)
	if len(enumValues) == 0 {
		// 不是枚举类型或没有枚举值，返回空
		return nil
	}

	// 输出（空格分隔）
	fmt.Println(strings.Join(enumValues, " "))

	return nil
}

// getEnumValues 获取标志的枚举值列表
//
// 参数:
//   - flag: 标志实例
//
// 返回值:
//   - []string: 枚举值列表，如果不是枚举类型则返回空切片
func getEnumValues(flag types.Flag) []string {
	// 检查是否是枚举类型标志
	// 这里假设有一个 EnumValues() 方法或可以通过其他方式获取枚举值
	// 具体实现取决于 QFlag 的标志类型定义

	// 示例实现：通过类型断言或接口检查
	if enumFlag, ok := flag.(interface{ EnumValues() []string }); ok {
		return enumFlag.EnumValues()
	}

	// 如果不是枚举类型，返回空切片
	return []string{}
}
```

### 3.5 集成到 dynamic.go

```go
// dynamic.go 中添加路由

const (
	InstructionFuzzy      = "fuzzy"
	InstructionContext    = "context"     // 新增
	InstructionCandidates = "candidates"  // 新增
	InstructionEnum       = "enum"        // 新增
)

// HandleDynamicComplete 处理 __complete 子命令
func HandleDynamicComplete(root types.Command, instruction string, params []string) error {
	switch instruction {
	case InstructionFuzzy:
		return handleFuzzy(params)
	case InstructionContext:
		return HandleContext(root, params)
	case InstructionCandidates:
		return HandleCandidates(root, params)
	case InstructionEnum:
		return HandleEnum(root, params)
	default:
		return fmt.Errorf("未知指令: %s", instruction)
	}
}
```

---

## 4. Shell 脚本简化

### 4.1 简化后的 Bash 模板

```bash
#!/usr/bin/env bash

# ==================== 主补全函数 ====================
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
	if [[ "$cur" == *"/"* || "$cur" == *"."* || "$cur" == *"~"* ]]; then
		COMPREPLY=($(compgen -f -d -- "$cur"))
		return 0
	fi

	# 2. 【核心】调用 Go 层计算上下文（替代静态 cmd_tree 循环）
	local context
	context=$({{.ProgramName}} __complete context "${words[@]}")
	
	# 检查上下文计算是否成功
	[[ -z "$context" ]] && return 1

	# 3. 【核心】调用 Go 层获取候选选项（替代查 cmd_tree 数组）
	local candidates
	candidates=$({{.ProgramName}} __complete candidates "$context")

	# 4. 检测是否是枚举值补全场景
	if [[ "$cur" == --*"="* ]]; then
		# 提取标志名和当前值
		flag_name="${cur%%=*}="
		current_value="${cur#*=}"
		
		# 4a. 获取枚举值（替代查 enum_options 数组）
		enum_values=$({{.ProgramName}} __complete enum "$context" "$flag_name")
		COMPREPLY=($(compgen -W "$enum_values" -- "$current_value"))
	else
		# 4b. 普通补全：模糊匹配
		if [[ -n "$candidates" ]]; then
			COMPREPLY=($(compgen -W "$candidates" -- "$cur"))
		fi
	fi

	return 0
}

# 注册补全函数
complete -F _{{.ProgramName}}_complete {{.ProgramName}}
```

### 4.2 简化后的 PowerShell 模板

```powershell
# 智能补全主函数
$scriptBlock = {
    param($wordToComplete, $commandAst, $cursorPosition)

    try {
        # 解析 tokens
        $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
        if (-not $tokens -or $tokens.Count -eq 0) { return @() }

        $currentIndex = $tokens.Count - 1

        # 1. 路径补全快速路径
        if ($wordToComplete -match '[/\~\.]' -or $wordToComplete -like './*') {
            return @(Get-ChildItem -Path "$wordToComplete*" | Select-Object -ExpandProperty Name)
        }

        # 2. 【核心】调用 Go 层计算上下文（替代静态 cmdTree 循环）
        $context = & {{.ProgramName}} __complete context @tokens
        
        if (-not $context) { return @() }

        # 3. 检测是否是枚举值补全场景
        if ($wordToComplete -match '^--[^=]+=') {
            # 3a. 枚举值补全（替代查 enumOptions 数组）
            $parts = $wordToComplete -split '=', 2
            $flagName = $parts[0] + "="
            $currentValue = $parts[1]
            
            $enumValues = & {{.ProgramName}} __complete enum $context $flagName
            $matches = & {{.ProgramName}} __complete fuzzy $currentValue ($enumValues -split ' ')
            return $matches
        } else {
            # 3b. 普通补全：调用 Go 层获取候选选项（替代查 cmdTree 数组）
            $candidates = & {{.ProgramName}} __complete candidates $context
            $allOptions = $candidates -split ' '
            
            if ($allOptions.Count -gt 0) {
                $matches = & {{.ProgramName}} __complete fuzzy $wordToComplete $allOptions
                return $matches
            }
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

---

## 5. 对比：静态 vs 动态

### 5.1 模板体积对比

| 项目 | 静态方案 | 动态方案 | 减少 |
|------|---------|---------|------|
| Bash 模板 | ~150 行 | ~30 行 | 80% |
| PowerShell 模板 | ~250 行 | ~35 行 | 86% |
| 生成数据 | cmd_tree, flag_params, enum_options | 无 | 100% |

### 5.2 功能对比

| 功能 | 静态方案 | 动态方案 |
|------|---------|---------|
| 子命令补全 | ✅ 支持 | ✅ 支持 |
| 标志补全 | ✅ 支持 | ✅ 支持 |
| 动态子命令 | ❌ 不支持 | ✅ 支持 |
| 运行时注册 | ❌ 不支持 | ✅ 支持 |
| 脚本体积 | 随命令树增大 | 恒定 |
| 维护成本 | 高（需同步更新） | 低（自动同步） |

### 5.3 性能对比

| 场景 | 静态方案 | 动态方案 |
|------|---------|---------|
| 补全延迟 | < 1ms | < 10ms |
| 内存占用 | 低（无进程创建） | 中（进程创建开销） |
| 首次调用 | 快 | 稍慢（进程启动） |
| 后续调用 | 快 | 快（进程缓存） |

---

## 6. 实现步骤

### 阶段一：Go 层实现

1. [ ] 创建 `internal/completion/context.go`
   - 实现 `ContextResult` 结构体
   - 实现 `HandleContext` 函数
   - 实现 `CalculateContext` 核心算法
   - 实现辅助函数 `getSubCommandNames`, `getFlagNames`

2. [ ] 创建 `internal/completion/candidates.go`
   - 实现 `HandleCandidates` 函数
   - 实现 `findCommandByContext` 辅助函数

3. [ ] 创建 `internal/completion/enum.go`
   - 实现 `HandleEnum` 函数
   - 实现 `getEnumValues` 辅助函数

4. [ ] 修改 `internal/completion/dynamic.go`
   - 添加 `InstructionContext` 常量
   - 添加 `InstructionCandidates` 常量
   - 添加 `InstructionEnum` 常量
   - 在 `HandleDynamicComplete` 中添加路由

5. [ ] 创建 `internal/completion/context_test.go`
   - 单元测试：`CalculateContext` 各种场景
   - 单元测试：`HandleCandidates` 各种场景
   - 单元测试：`HandleEnum` 各种场景
   - 边界测试：空输入、无效输入
   - 性能测试：大命令树的处理速度

### 阶段二：集成到 builtin

1. [ ] 修改 `internal/cmd/builtin.go`
   - 确保 `createCompleteCmd` 正确传递 root 命令
   - 测试 `__complete context` 指令可用性

### 阶段三：模板更新

1. [ ] 创建 `bash_dynamic_v2.tmpl`
   - 移除静态 cmd_tree 定义
   - 移除静态 enum_options 定义
   - 使用 `__complete context` 计算上下文
   - 使用 `__complete candidates` 获取候选选项
   - 使用 `__complete enum` 获取枚举值

2. [ ] 创建 `pwsh_dynamic_v2.tmpl`
   - 移除静态 cmdTree 定义
   - 移除静态 enumOptions 定义
   - 使用 `__complete context` 计算上下文
   - 使用 `__complete candidates` 获取候选选项
   - 使用 `__complete enum` 获取枚举值

### 阶段四：测试验证

1. [ ] 功能测试
   - 多级子命令上下文计算
   - 标志上下文识别
   - 无效子命令处理

2. [ ] 集成测试
   - 完整补全流程
   - Bash 和 PowerShell 一致性

3. [ ] 性能测试
   - 确保延迟 < 10ms
   - 大命令树（100+ 子命令）测试

---

## 7. 使用示例

### 7.1 命令行测试

**测试 context 指令：**

```bash
# 测试上下文计算
myapp __complete context "myapp" "server" "start"
# 输出: /server/start/

# 测试根上下文
myapp __complete context "myapp"
# 输出: /

# 测试标志上下文
myapp __complete context "myapp" "server" "--port"
# 输出: /server/
```

**测试 candidates 指令：**

```bash
# 获取根命令的候选选项
myapp __complete candidates "/"
# 输出: server client --help -h

# 获取子命令的候选选项
myapp __complete candidates "/server/"
# 输出: start stop --port --host --help -p -h

# 获取深层子命令的候选选项
myapp __complete candidates "/server/start/"
# 输出: --config --daemon --help -c -d -h
```

**测试 enum 指令：**

```bash
# 获取 --log-level 标志的枚举值
myapp __complete enum "/server/start/" "--log-level"
# 输出: debug info warn error fatal

# 获取短名称标志的枚举值
myapp __complete enum "/server/start/" "-o"
# 输出: json yaml text

# 获取带 "=" 的标志枚举值
myapp __complete enum "/server/start/" "--format="
# 输出: table json yaml

# 无效标志返回空
myapp __complete enum "/server/start/" "--invalid"
# 输出: （空）

# 非枚举标志返回空
myapp __complete enum "/server/start/" "--port"
# 输出: （空）
```

**测试不完整输入：**

```bash
myapp __complete context "myapp" "serv"
# 输出: / （因为 "serv" 不是有效子命令）
```

### 7.2 程序中使用

```go
package main

import (
	"fmt"
	"gitee.com/MM-Q/qflag/internal/completion"
	"gitee.com/MM-Q/qflag/internal/cmd"
)

func main() {
	root := cmd.NewCmd("myapp", "", cmd.ExitOnError)
	
	// 添加子命令...
	serverCmd := cmd.NewCmd("server", "", cmd.ExitOnError)
	root.AddSubCmds(serverCmd)
	
	// 测试上下文计算
	result := completion.CalculateContext(root, 
		[]string{"myapp", "server", "start"}, 3)
	
	fmt.Printf("上下文: %s\n", result.Context)
	
	// 测试候选选项获取
	candidates, _ := completion.GetCandidates(root, result.Context)
	fmt.Printf("候选选项: %v\n", candidates)
	
	// 测试枚举值获取
	enumValues, _ := completion.GetEnumValues(root, result.Context, "--log-level")
	fmt.Printf("枚举值: %v\n", enumValues)
}
```

---

## 8. 注意事项

### 8.1 性能优化

- **进程创建开销**：每次调用 `__complete` 都会创建新进程
- **优化建议**：
  - 保持逻辑简单，避免复杂计算
  - 后续可考虑常驻进程或缓存机制

### 8.2 错误处理

- **命令未找到**：返回根上下文 `/`
- **无效参数**：返回清晰错误信息
- **注册表为空**：返回空列表而非错误

### 8.3 向后兼容

- 保留原有静态模板作为 fallback
- 新增 `EnableDynamicCompletion` 配置选项
- 默认使用动态方案，可选静态方案

---

## 9. 总结

`__complete context`、`__complete candidates` 和 `__complete enum` 三个指令的设计实现了：

1. **完全替代静态数据结构**：
   - `context` 指令替代原脚本中的上下文计算循环
   - `candidates` 指令替代查 `cmd_tree` 静态数组
   - `enum` 指令替代查 `enum_options` 静态数组

2. **支持动态子命令**：运行时注册的子命令也能正确补全

3. **支持动态枚举值**：运行时定义的枚举值也能正确补全

4. **简化 Shell 脚本**：模板体积减少 80%-86%
   - Bash: ~150 行 → ~30 行
   - PowerShell: ~250 行 → ~35 行

5. **统一跨平台逻辑**：Bash 和 PowerShell 使用相同的 Go 层算法

6. **纯文本输出**：Shell 脚本无需解析 JSON，直接处理空格分隔的字符串

7. **易于测试**：可以为上下文计算、候选获取和枚举值查询编写单元测试

这是实现完全动态补全的核心方案，三个指令配合使用可以完全消除 Shell 脚本中的静态数据定义。

---

## 10. 变更记录

### v1.3 (2026-04-10)
- **新增**: `enum` 指令用于枚举值补全，替代静态 `enum_options` 数组
- **完善**: 三个指令（`context` + `candidates` + `enum`）形成完整的动态补全方案
- **更新**: Shell 模板示例增加枚举值补全场景处理

### v1.2 (2026-04-10)
- **重构**: 明确需要两个指令 `context` + `candidates` 才能完全替代静态命令树
- **简化**: 移除 JSON 输出，改为纯文本空格分隔，Shell 处理更简单
- **优化**: 更新 Shell 模板示例，展示两个指令的配合使用

### v1.1 (2026-04-10)
- **优化**: cursorPos 由 Go 层自动计算，Shell 脚本无需传递
- **简化**: Shell 调用方式从 `__complete context <cursorPos> <args...>` 简化为 `__complete context <args...>`
- **优势**: 进一步简化 Shell 脚本，跨平台逻辑更统一

### v1.0 (2026-04-10)
- **初始版本**: 完整的 `__complete context` 指令设计方案
- **特性**: 动态上下文计算，替代静态命令树

---

*文档版本: 1.3*
*更新日期: 2026-04-10*
