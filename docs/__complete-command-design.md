# __complete 子命令设计方案

## 1. 背景与目标

### 1.1 当前问题

目前 qflag 的自动补全功能完全依赖 Shell 脚本实现（bash.tmpl / pwsh.tmpl），存在以下问题：

- **重复实现**：每个 Shell 都要单独实现模糊匹配、枚举处理等逻辑
- **维护困难**：修改算法需要同时改 bash 和 pwsh 两套脚本
- **功能受限**：Shell 脚本能力有限，复杂逻辑难以实现
- **测试困难**：难以对 Shell 脚本进行单元测试

### 1.2 设计目标

引入 `__complete` 隐藏子命令，将补全核心逻辑迁移到 Go 代码中：

- 实现跨平台统一的补全逻辑
- 简化 Shell 脚本，降低维护成本
- 支持更复杂的补全场景
- 便于单元测试

## 2. 总体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Shell 层 (简化)                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Bash      │  │   PowerShell│  │   其他 Shell        │  │
│  │  补全函数    │  │  补全函数    │  │   ...               │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
│         │                │                    │             │
│         └────────────────┼────────────────────┘             │
│                          │                                  │
│                          ▼                                  │
│              yourapp __complete <指令> [参数]                │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go 核心层                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                 __complete 子命令                    │   │
│  │              ┌─────────────┐                        │   │
│  │              │    fuzzy    │                        │   │
│  │              │   模糊匹配   │                        │   │
│  │              └─────────────┘                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              补全算法实现 (Go)                       │   │
│  │  - 模糊匹配算法                                       │   │
│  │  - 高性能评分系统                                     │   │
│  │  - 缓存机制                                           │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## 3. 命令设计

### 3.1 基本用法

```bash
yourapp __complete <指令> [参数...]
```

### 3.2 第一阶段：模糊匹配指令 (fuzzy)

| 指令 | 功能描述 | 参数 | 输出格式 |
|------|----------|------|----------|
| `fuzzy` | 执行模糊匹配 | `<模式> <候选1> [候选2] ...` | 每行一个匹配结果：匹配到的字符串 |

#### 3.2.1 fuzzy 指令详细说明

执行模糊匹配算法，返回按匹配质量排序的匹配字符串列表。

```bash
# 基本用法
yourapp __complete fuzzy "po" "port" "path" "pod" "proxy" "server"
# 输出:
# port
# pod
# proxy
# path
```

**参数说明：**
- 第一个参数：匹配模式（pattern）
- 后续参数：候选列表

**输出格式：**
```
匹配字符串
```
- 每行一个匹配结果
- 按匹配质量降序排列
- 只输出匹配的字符串，不包含分数

**匹配算法：**
使用 `github.com/sahilm/fuzzy` 包提供的 Find 算法：
- 第一个字符匹配奖励
- 分隔符后匹配奖励（如 `-` `_` `.` 等）
- 驼峰命名匹配奖励
- 相邻字符匹配奖励
- 未匹配前导字符惩罚

**使用场景：**
- Shell 补全函数中用于模糊匹配命令名、标志名
- 替代 Shell 脚本中复杂的模糊匹配逻辑

## 4. 实现设计

### 4.1 目录结构

```
internal/completion/
├── completion.go              # 现有：补全脚本生成
├── bash_completion.go         # 现有：Bash 补全
├── pwsh_completion.go         # 现有：PowerShell 补全
├── templates/
│   ├── bash.tmpl              # 现有：Bash 模板
│   ├── pwsh.tmpl              # 现有：PowerShell 模板
│   ├── bash_v2.tmpl            # 新增：Bash 模板 (v2)
│   └── pwsh_v2.tmpl            # 新增：PowerShell 模板 (v2)
└── cmdcomplete/               # 新增：__complete 命令实现
    └── fuzzy.go               # 模糊匹配指令实现
```

### 4.2 fuzzy 指令实现

文件：`internal/completion/cmdcomplete/fuzzy.go`

```go
// Package cmdcomplete 实现 __complete 子命令的核心逻辑
package cmdcomplete

import (
	"fmt"

	"github.com/sahilm/fuzzy"
)

// HandleFuzzy 处理 fuzzy 指令
//
// 参数:
//   - args: 参数列表，第一个是模式，后面是候选列表
//
// 返回值:
//   - error: 处理错误
//
// 输出格式: 每行一个匹配结果（按匹配质量降序）
func HandleFuzzy(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("用法: __complete fuzzy <模式> <候选1> [候选2] ...")
	}

	pattern := args[0]
	candidates := args[1:]

	// 使用 go-kit fuzzy 包执行模糊匹配
	matches := fuzzy.Find(pattern, candidates)

	// 输出匹配结果（只输出匹配的字符串）
	for _, match := range matches {
		fmt.Println(match.Str)
	}

	return nil
}
```

### 4.3 路由实现

文件：`internal/completion/cmdcomplete/router.go`

```go
// Package cmdcomplete 实现 __complete 子命令的核心逻辑
package cmdcomplete

import "fmt"

// Handle 处理 __complete 子命令的路由
//
// 参数:
//   - instruction: 指令名称
//   - params: 指令参数列表
//
// 返回值:
//   - error: 处理错误
func Handle(instruction string, params []string) error {
	switch instruction {
	case "fuzzy":
		return HandleFuzzy(params)
	default:
		return fmt.Errorf("未知指令: %s", instruction)
	}
}
```

### 4.4 在 builtin.go 中调用

修改 `internal/cmd/builtin.go`，负责参数校验和拆分：

```go
package cmd

import (
	"fmt"

	"gitee.com/MM-Q/qflag/internal/completion/cmdcomplete"
	"gitee.com/MM-Q/qflag/internal/types"
)

const completeCmdName = "__complete"

// createCompleteCmd 创建动态补全子命令
func createCompleteCmd() types.Command {
	cmd := NewCmd(completeCmdName, "", types.ExitOnError)
	cmd.SetDesc("内部命令：执行动态补全")
	cmd.SetHidden(true)
	cmd.SetDisableFlagParsing(true) // 禁用标志解析

	cmd.SetRun(func(c types.Command) error {
		args := c.Args()

		// 校验参数数量
		if len(args) < 2 {
			return fmt.Errorf("用法: __complete <指令> [参数...]")
		}

		// 拆分参数：第一个是指令，后续是参数列表
		instruction := args[0]
		params := args[1:]

		// 传递给 completion 包处理
		return cmdcomplete.Handle(instruction, params)
	})

	return cmd
}
```

## 5. 新的 Shell 模板

### 5.1 新增模板文件

- `internal/completion/templates/bash_v2.tmpl` - 基于现有模板，替换模糊匹配逻辑
- `internal/completion/templates/pwsh_v2.tmpl` - 基于现有模板，替换模糊匹配逻辑

### 5.2 模板修改策略

新模板基于现有模板复制，**仅修改模糊匹配相关部分**：

1. **保留内容**：
   - 配置参数（模糊补全开关、阈值等）
   - 静态数据定义（命令树、标志参数、枚举选项）
   - 缓存机制
   - 补全主流程逻辑

2. **替换内容**：
   - 删除内置的 `_fuzzy_score_fast` / `Get-FuzzyScoreFast` 函数
   - 删除内置的 `_fuzzy_score_cached` / `Get-FuzzyScoreCached` 函数
   - 将模糊匹配调用改为 `__complete fuzzy` 子命令

### 5.3 Bash 模板修改示例

**原逻辑（删除）：**
```bash
# ==================== 模糊匹配核心算法 ====================
_{{.ProgramName}}_fuzzy_score_fast() {
    local pattern="$1"
    local candidate="$2"
    # ... 复杂的评分算法 ...
}

_{{.ProgramName}}_fuzzy_score_cached() {
    local pattern="$1"
    local candidate="$2"
    # ... 缓存逻辑 ...
    _{{.ProgramName}}_fuzzy_score_fast "$pattern" "$candidate"
}
```

**新逻辑（替换）：**
```bash
# ==================== 模糊匹配（调用 Go 实现）====================
_{{.ProgramName}}_fuzzy_match() {
    local pattern="$1"
    shift
    local candidates=("$@")
    
    # 调用 __complete fuzzy 子命令执行模糊匹配
    {{.ProgramName}} __complete fuzzy "$pattern" "${candidates[@]}"
}
```

**模糊匹配调用处修改：**
```bash
# 原逻辑：使用内置函数评分
for opt in "${opts_arr[@]}"; do
    local score
    score=$(_{{.ProgramName}}_fuzzy_score_cached "$pattern" "$opt")
    if [[ $score -ge ${{.ProgramName}}_FUZZY_SCORE_THRESHOLD ]]; then
        scored_matches+=("$score:$opt")
    fi
done

# 新逻辑：调用 __complete fuzzy
local fuzzy_results
fuzzy_results=$(_{{.ProgramName}}_fuzzy_match "$pattern" "${opts_arr[@]}")
while IFS= read -r match; do
    [[ -n "$match" ]] && scored_matches+=("$match")
done <<< "$fuzzy_results"
```

### 5.4 PowerShell 模板修改示例

**原逻辑（删除）：**
```powershell
# ==================== 模糊匹配核心算法 ====================
function Get-{{.SanitizedName}}FuzzyScoreFast {
    param([string]$Pattern, [string]$Candidate)
    # ... 复杂的评分算法 ...
}

function Get-{{.SanitizedName}}FuzzyScoreCached {
    param([string]$Pattern, [string]$Candidate)
    # ... 缓存逻辑 ...
    Get-{{.SanitizedName}}FuzzyScoreFast $Pattern $Candidate
}
```

**新逻辑（替换）：**
```powershell
# ==================== 模糊匹配（调用 Go 实现）====================
function Get-{{.SanitizedName}}FuzzyMatch {
    param(
        [string]$Pattern,
        [string[]]$Candidates
    )
    
    # 调用 __complete fuzzy 子命令执行模糊匹配
    & {{.ProgramName}} __complete fuzzy $Pattern $Candidates
}
```

**模糊匹配调用处修改：**
```powershell
# 原逻辑：使用内置函数评分
foreach ($opt in $optsArr) {
    $score = Get-{{.SanitizedName}}FuzzyScoreCached $pattern $opt
    if ($score -ge $script:{{.SanitizedName}}_FUZZY_SCORE_THRESHOLD) {
        $scoredMatches += "$score`:$opt"
    }
}

# 新逻辑：调用 __complete fuzzy
$fuzzyResults = Get-{{.SanitizedName}}FuzzyMatch $pattern $optsArr
foreach ($match in $fuzzyResults) {
    if ($match) { $scoredMatches += $match }
}
```

## 6. 实现步骤

### 阶段一：核心实现

1. [ ] 创建 `internal/completion/cmdcomplete/` 目录
2. [ ] 实现 `fuzzy.go` - 模糊匹配算法
3. [ ] 添加 `fuzzy_test.go` - 单元测试
4. [ ] 修改 `builtin.go` - 集成 fuzzy 指令

### 阶段二：模板创建

1. [ ] 创建 `bash_v2.tmpl` - 简化版 Bash 模板
2. [ ] 创建 `pwsh_v2.tmpl` - 简化版 PowerShell 模板
3. [ ] 在 `completion.go` 中添加新模板支持

### 阶段三：测试验证

1. [ ] 测试 fuzzy 指令功能
2. [ ] 对比新旧模板输出
3. [ ] 性能测试

## 7. 注意事项

### 7.1 向后兼容

- 现有模板（bash.tmpl / pwsh.tmpl）保持不变
- 新模板（bash_v2.tmpl / pwsh_v2.tmpl）使用新逻辑
- 用户可选择使用哪种模板

### 7.2 性能考虑

- `__complete fuzzy` 执行要快（< 10ms）
- 避免频繁调用（Shell 层应缓存结果）

### 7.3 错误处理

- fuzzy 指令返回清晰的错误信息
- Shell 脚本要能优雅处理失败情况

## 8. 总结

第一阶段实现 `fuzzy` 指令，将模糊匹配逻辑从 Shell 脚本迁移到 Go 代码中：

1. **统一算法**：一套 Go 代码服务所有平台
2. **简化脚本**：Shell 脚本只需调用指令
3. **便于测试**：可以写单元测试验证算法
4. **易于扩展**：后续可添加更多指令

