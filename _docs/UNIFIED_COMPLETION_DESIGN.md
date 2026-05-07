# QFlag 统一补全池设计方案

> **文档版本**: 1.0  
> **设计日期**: 2026-04-05  
> **设计目标**: 优化 Shell 补全体验，实现标志、命令、文件的统一模糊补全

---

## 一、设计背景与问题

### 1.1 现有问题

当前补全模板存在以下体验问题：

| 问题 | 现象 | 影响 |
|------|------|------|
| 路径补全触发困难 | 必须输入 `./` 或包含 `/` 才能触发路径补全 | 用户难以发现文件补全功能 |
| 标志与文件冲突 | 输入 `conf` 时，优先匹配 `--config` 标志而非 `config.json` 文件 | 经常需要手动输入完整文件名 |
| 类型判断复杂 | 需要预先判断用户想要补全什么类型 | 逻辑复杂，容易误判 |
| 上下文切换不流畅 | 标志后、命令后、参数后的补全逻辑不一致 | 学习成本高 |

### 1.2 设计目标

- **零学习成本**：用户无需了解触发规则，按 Tab 即可得到期望结果
- **智能排序**：最相关的结果自动排在前面
- **统一体验**：无论当前上下文如何，补全行为一致
- **模糊友好**：支持拼音、缩写、子串等多种模糊匹配

---

## 二、核心设计思想

### 2.1 统一补全池

将所有可能的补全候选放入一个统一的候选池，通过智能评分排序返回：

```
┌─────────────────────────────────────────────────────────────┐
│                     统一补全池                               │
├─────────────────────────────────────────────────────────────┤
│  文件候选        │  config.json, main.go, README.md, ...    │
│  标志候选        │  --config, --verbose, --output, ...       │
│  子命令候选      │  build, run, test, init, ...              │
│  枚举值候选      │  json, yaml, xml, ... (仅枚举标志后)      │
└─────────────────────────────────────────────────────────────┘
                           ↓
                    统一模糊匹配评分
                           ↓
              排序返回最相关的候选列表
```

### 2.2 候选来源优先级

| 优先级 | 候选类型 | 说明 |
|--------|---------|------|
| P0 | 当前目录文件 | 用户最可能想要输入的内容 |
| P1 | 当前上下文标志 | 当前命令可用的标志 |
| P1 | 当前上下文子命令 | 当前命令可用的子命令 |
| P2 | 枚举值 | 仅在枚举类型标志后显示 |
| P3 | 上级目录文件 | 可选，防止候选过多 |

### 2.3 评分算法

综合评分 = 基础匹配分 + 类型加成 + 位置加成

#### 基础匹配分 (0-100)

| 匹配类型 | 分数 | 说明 |
|---------|------|------|
| 精确前缀匹配 | 100 | `conf` → `config.json` |
| 大小写不敏感前缀 | 95 | `conf` → `Config.json` |
| 单词边界匹配 | 90 | `cj` → `config.json` |
| 连续子串匹配 | 80 | `fig` → `config.json` |
| 分散字符匹配 | 60 | `cnj` → `config.json` |
| 模糊匹配 | 30-50 | 根据相似度计算 |

#### 类型加成

| 类型 | 加成 | 原因 |
|------|------|------|
| 文件 | +5 | 用户通常想要操作文件 |
| 目录 | +3 | 目录也是常见操作目标 |
| 标志 | +0 | 标准优先级 |
| 子命令 | +0 | 标准优先级 |

#### 位置加成

| 位置 | 加成 | 说明 |
|------|------|------|
| 起始位置匹配 | +10 | 从开头匹配更相关 |
| 近期使用 | +5 | 记忆用户习惯（可选） |

---

## 三、详细设计方案

### 3.1 渐进式混合补全流程

采用**渐进式 + 双策略混合**设计：
- **渐进式**：根据输入内容逐步缩小候选范围
- **双策略**：明显路径特征走专用路径补全，其他情况走统一补全池

```
用户按 Tab
    ↓
1. 解析当前上下文
   ├── 获取当前命令路径
   ├── 获取前一个参数
   └── 获取当前输入词 (cur)
    ↓
2. 判断输入类型
   │
   ├── 带有明显路径特征? (包含 / 或 ./ 或 ../ 或 ~)
   │   └── 是 → 使用专用路径补全
   │       └── 直接返回匹配的文件/目录列表
   │
   └── 否 → 使用渐进式统一补全池
       ↓
       3. 收集候选池
          ├── 收集当前目录文件候选（最多100个）
          ├── 收集当前上下文标志候选
          ├── 收集当前上下文子命令候选
          └── 收集枚举值候选（如适用）
       ↓
       4. 模糊匹配过滤
          ├── 根据输入词过滤候选池
          ├── 只保留匹配度 > 0 的候选
          └── 得到相关候选列表
       ↓
       5. 评分排序
          ├── 对每个相关候选计算匹配分数
          ├── 应用类型加成（文件+5分）
          └── 按分数降序排序
       ↓
       6. 返回结果
          └── 返回所有相关候选（已模糊匹配过滤并排序）
       ↓
       6. 评分排序
          ├── 对每个相关候选计算匹配分数
          ├── 应用类型加成（文件+5分）
          └── 按分数降序排序
       ↓
       6. 返回结果
          └── 返回所有相关候选（已模糊匹配过滤并排序）
```

### 3.2 路径特征判断规则

| 特征 | 示例 | 处理方式 |
|------|------|---------|
| 包含 `/` | `config/`, `/home/user` | 专用路径补全 |
| 以 `./` 开头 | `./config` | 专用路径补全 |
| 以 `../` 开头 | `../config` | 专用路径补全 |
| 以 `~` 开头 | `~/config` | 专用路径补全 |
| 无上述特征 | `config`, `main` | 统一补全池 |

**说明**：
- 专用路径补全：使用原有的 `compgen -f -d` 或 `Get-ChildItem` 逻辑
- 统一补全池：将文件、标志、命令混在一起模糊匹配
- 文件候选仅收集**当前目录**（不包含上级目录）

### 3.3 候选收集策略

#### 文件候选收集（统一补全池场景）

```bash
# 收集范围（仅当前目录）
- 当前目录下的文件和子目录
- 不包含上级目录（../）
- 不包含其他路径

# 限制策略
- 最大收集数量: 100 个
- 隐藏文件: 默认不显示（以 . 开头）
- 目录标记: 目录名后加 /

# 注意：带有路径特征的输入不走此逻辑，直接走专用路径补全
```

#### 标志候选收集

```bash
# 收集范围
- 当前命令上下文中注册的所有标志
- 包括长名称和短名称

# 特殊处理
- 已使用的标志: 降低优先级（可选）
- 互斥标志: 如果互斥组中一个已使用，其他降低优先级
```

#### 子命令候选收集

```bash
# 收集范围
- 当前命令下注册的所有子命令

# 特殊处理
- 无
```

#### 枚举值候选收集

```bash
# 触发条件
- 前一个参数是枚举类型的标志

# 收集范围
- 该枚举标志定义的所有允许值
```

### 3.4 返回策略

**核心原则**：先模糊匹配过滤，再评分排序，返回所有相关候选

```
流程：
全部候选（文件最多100个） → 模糊匹配筛选 → 评分排序 → 返回所有相关候选
```

| 场景 | 处理流程 | 说明 |
|------|---------|------|
| 空输入 | 返回全部候选 | 展示所有可能性 |
| 有输入 | 模糊匹配 → 评分排序 → 返回所有相关 | 只返回匹配的内容 |

**说明**：
- **文件数量已限制**：候选收集阶段最多 100 个文件，无需额外限制
- **模糊匹配过滤**：只有匹配度 > 0 的候选才会被返回
- **Bash**：显示所有相关候选列表
- **PowerShell**：循环遍历所有相关候选

### 3.5 过滤规则

| 场景 | 过滤规则 |
|------|---------|
| 输入带有路径特征 (`/` `./` `../` `~`) | 不走统一补全池，直接使用专用路径补全 |
| 输入以 `-` 开头 | 只保留标志候选，过滤文件和命令 |
| 前一个参数是枚举标志 | 优先枚举值，同时保留文件候选 |
| 前一个参数是其他标志 | 优先文件候选 |
| 默认情况 | 保留所有类型（文件+标志+命令） |

---

## 四、模板实现方案

### 4.1 Bash 模板修改

#### 新增函数

```bash
# 收集所有候选到统一池
_{{.ProgramName}}_collect_all_candidates() {
    local context="$1"
    local cur="$2"
    local candidates=()
    
    # 1. 收集文件候选
    while IFS= read -r -d '' file; do
        local basename=$(basename "$file")
        # 目录加 / 后缀
        if [[ -d "$file" ]]; then
            candidates+=("$basename/")
        else
            candidates+=("$basename")
        fi
    done < <(find . -maxdepth 1 -name "$cur*" -not -name ".*" -print0 2>/dev/null | head -z -n 100)
    
    # 2. 收集标志候选
    local context_flags="${{ .ProgramName }}_flag_params[$context]"
    for flag in ${context_flags[@]}; do
        candidates+=("$flag")
    done
    
    # 3. 收集子命令候选
    local context_cmds="${{ .ProgramName }}_cmd_tree[$context]"
    for cmd in ${context_cmds[@]}; do
        # 排除标志（以 - 开头）
        if [[ "$cmd" != -* ]]; then
            candidates+=("$cmd")
        fi
    done
    
    # 4. 如果前一个参数是枚举标志，收集枚举值
    local prev_flag_type="${{ .ProgramName }}_flag_types[$prev]"
    if [[ "$prev_flag_type" == "enum" ]]; then
        local enum_values="${{ .ProgramName }}_enum_options[$prev]"
        for val in ${enum_values[@]}; do
            candidates+=("$val")
        done
    fi
    
    printf '%s\n' "${candidates[@]}"
}

# 统一评分排序
_{{.ProgramName}}_score_and_sort() {
    local pattern="$1"
    shift
    local candidates=("$@")
    local scored=()
    
    for candidate in "${candidates[@]}"; do
        local score=0
        local type_bonus=0
        
        # 计算基础匹配分
        score=$(_{{.ProgramName}}_fuzzy_score_fast "$pattern" "$candidate")
        
        # 类型加成
        if [[ -f "$candidate" ]]; then
            type_bonus=5
        elif [[ -d "$candidate" ]]; then
            type_bonus=3
        fi
        
        # 起始位置加成
        local candidate_lower="${candidate,,}"
        local pattern_lower="${pattern,,}"
        if [[ "$candidate_lower" == "$pattern_lower"* ]]; then
            type_bonus=$((type_bonus + 10))
        fi
        
        local final_score=$((score + type_bonus))
        scored+=("$final_score:$candidate")
    done
    
    # 排序并返回
    printf '%s\n' "${scored[@]}" | sort -t: -k1 -nr | cut -d: -f2 | head -20
}
```

#### 修改主补全函数（渐进式）

```bash
_{{.ProgramName}}_complete() {
    local cur prev words cword context
    COMPREPLY=()
    
    # 获取补全参数
    _get_comp_words_by_ref -n =: cur prev words cword
    
    # ===== 策略1: 路径特征检测 =====
    # 如果输入带有明显路径特征，使用专用路径补全
    if [[ "$cur" == *"/"* || "$cur" == "./"* || "$cur" == "../"* || "$cur" == "~"* ]]; then
        COMPREPLY=($(compgen -f -d -- "$cur"))
        return 0
    fi
    
    # ===== 策略2: 渐进式统一补全池 =====
    # 计算当前上下文
    context=$(_{{.ProgramName}}_get_context "${words[@]}")
    
    # 收集所有候选
    local all_candidates=($(_{{.ProgramName}}_collect_all_candidates "$context" "$cur"))
    
    # 根据输入过滤（类型过滤）
    local type_filtered_candidates=()
    if [[ "$cur" == -* ]]; then
        # 输入以 - 开头，只保留标志
        for cand in "${all_candidates[@]}"; do
            [[ "$cand" == -* ]] && type_filtered_candidates+=("$cand")
        done
    else
        type_filtered_candidates=("${all_candidates[@]}")
    fi
    
    # 模糊匹配过滤
    local fuzzy_filtered_candidates=()
    if [[ ${#cur} -eq 0 ]]; then
        # 空输入：跳过模糊匹配，使用全部
        fuzzy_filtered_candidates=("${type_filtered_candidates[@]}")
    else
        # 有输入：先进行模糊匹配，只保留相关候选
        for cand in "${type_filtered_candidates[@]}"; do
            local match_score=$(_{{.ProgramName}}_fuzzy_match "$cur" "$cand")
            if [[ $match_score -gt 0 ]]; then
                fuzzy_filtered_candidates+=("$cand")
            fi
        done
    fi
    
    # 评分排序并返回所有相关候选
    if [[ ${#fuzzy_filtered_candidates[@]} -gt 0 ]]; then
        local sorted=($(_{{.ProgramName}}_score_and_sort "$cur" "${fuzzy_filtered_candidates[@]}"))
        # 返回所有相关候选（不限制数量，已在收集阶段限制文件最多100个）
        COMPREPLY=("${sorted[@]}")
    fi
    
    return 0
}
```

### 4.2 PowerShell 模板修改

#### 新增函数

```powershell
# 收集所有候选到统一池
function Get-{{.SanitizedName}}AllCandidates {
    param(
        [string]$Context,
        [string]$Pattern,
        [string]$PrevElement
    )
    
    $candidates = [System.Collections.ArrayList]::new()
    
    # 1. 收集文件候选
    try {
        $files = Get-ChildItem -Path "." -Name -Filter "$Pattern*" -ErrorAction SilentlyContinue | 
                Select-Object -First 100
        foreach ($file in $files) {
            $fullPath = Join-Path "." $file
            if (Test-Path -Path $fullPath -PathType Container) {
                [void]$candidates.Add("$file/")
            } else {
                [void]$candidates.Add($file)
            }
        }
    } catch {
        # 忽略文件访问错误
    }
    
    # 2. 收集标志候选
    $flags = ${{$.SanitizedName}}_flagParams | 
             Where-Object { $_.Context -eq $Context } | 
             Select-Object -ExpandProperty Parameter
    $candidates.AddRange($flags)
    
    # 3. 收集子命令候选
    $cmds = ${{$.SanitizedName}}_cmdTree | 
            Where-Object { $_.Context -eq $Context -and $_.Options -notlike "-*" } | 
            Select-Object -ExpandProperty Options
    $candidates.AddRange($cmds)
    
    # 4. 如果前一个参数是枚举标志，收集枚举值
    $prevFlag = ${{$.SanitizedName}}_flagParams | 
                Where-Object { $_.Context -eq $Context -and $_.Parameter -eq $PrevElement }
    if ($prevFlag -and $prevFlag.ValueType -eq 'enum') {
        $candidates.AddRange($prevFlag.Options)
    }
    
    return $candidates.ToArray()
}

# 统一评分排序
function Get-{{.SanitizedName}}ScoredCompletions {
    param(
        [string]$Pattern,
        [array]$Candidates
    )
    
    $scored = [System.Collections.ArrayList]::new()
    
    foreach ($candidate in $Candidates) {
        $score = 0
        $typeBonus = 0
        
        # 基础匹配分
        $score = Get-{{.SanitizedName}}FuzzyScoreCached -Pattern $Pattern -Candidate $candidate
        
        # 类型加成
        $fullPath = Join-Path "." $candidate
        if (Test-Path -Path $fullPath -PathType Leaf) {
            $typeBonus = 5
        } elseif (Test-Path -Path $fullPath -PathType Container) {
            $typeBonus = 3
        }
        
        # 起始位置加成
        if ($candidate.ToLower().StartsWith($Pattern.ToLower())) {
            $typeBonus += 10
        }
        
        $finalScore = $score + $typeBonus
        [void]$scored.Add([PSCustomObject]@{
            Candidate = $candidate
            Score = $finalScore
        })
    }
    
    # 排序并返回
    return $scored | Sort-Object Score -Descending | 
           Select-Object -First 20 | 
           Select-Object -ExpandProperty Candidate
}
```

#### 修改主补全逻辑（渐进式）

```powershell
$scriptBlock = {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    try {
        # 解析令牌
        $tokens = $commandAst.CommandElements | ForEach-Object { $_.Extent.Text }
        $currentIndex = $tokens.Count - 1
        $prevElement = if ($currentIndex -ge 1) { $tokens[$currentIndex - 1] } else { $null }
        
        # ===== 策略1: 路径特征检测 =====
        # 如果输入带有明显路径特征，使用专用路径补全
        if ($wordToComplete -match '[/\~]' -or 
            $wordToComplete -match '^\.\.?[/\]') {
            return Get-{{.SanitizedName}}PathCompletions -WordToComplete $wordToComplete
        }
        
        # ===== 策略2: 渐进式统一补全池 =====
        # 计算上下文
        $context = Get-{{.SanitizedName}}Context -Tokens $tokens
        
        # 收集所有候选
        $allCandidates = Get-{{.SanitizedName}}AllCandidates `
            -Context $context `
            -Pattern $wordToComplete `
            -PrevElement $prevElement
        
        # 根据输入过滤（类型过滤）
        $typeFiltered = if ($wordToComplete.StartsWith('-')) {
            $allCandidates | Where-Object { $_ -like '-*' }
        } else {
            $allCandidates
        }
        
        # 模糊匹配过滤
        $fuzzyFiltered = if ($wordToComplete.Length -eq 0) {
            # 空输入：跳过模糊匹配，使用全部
            $typeFiltered
        } else {
            # 有输入：先进行模糊匹配，只保留相关候选
            $typeFiltered | Where-Object { 
                $score = Get-{{.SanitizedName}}FuzzyMatchScore -Pattern $wordToComplete -Candidate $_
                $score -gt 0
            }
        }
        
        # 评分排序并返回所有相关候选
        if ($fuzzyFiltered.Count -gt 0) {
            $sorted = Get-{{.SanitizedName}}ScoredCompletions `
                -Pattern $wordToComplete `
                -Candidates $fuzzyFiltered
            
            # 返回所有相关候选（不限制数量，已在收集阶段限制文件最多100个）
            return $sorted
        }
        
        return @()
    } catch {
        Write-Debug "补全错误: $($_.Exception.Message)"
        return @()
    }
}
```

---

## 五、性能优化

### 5.1 缓存策略

| 缓存内容 | 缓存时长 | 说明 |
|---------|---------|------|
| 文件列表 | 1秒 | 目录内容不会频繁变化 |
| 模糊评分结果 | 会话级 | 相同输入的评分结果不变 |
| 候选池 | 实时 | 上下文可能随时变化 |

### 5.2 数量限制策略

**唯一限制**：候选收集阶段限制文件数量

| 限制项 | 数值 | 原因 |
|--------|------|------|
| 最大文件候选 | 100 | 防止大目录卡顿 |

**说明**：
- 文件候选在收集阶段已限制最多 100 个
- 模糊匹配后返回所有相关候选，不再额外限制
- 标志和命令候选数量通常较少，无需限制

**Bash 场景**：
- 显示所有相关候选列表
- 用户可通过继续输入进一步过滤

**PowerShell 场景**：
- 循环遍历所有相关候选
- 用户可通过继续输入快速定位

### 5.3 其他限制

| 限制项 | 数值 | 原因 |
|--------|------|------|
| 最大文件候选 | 100 | 防止大目录卡顿 |
| 最大缓存条目 | 500 | 内存保护 |

### 5.4 异步处理（可选）

对于文件收集较慢的场景，可以考虑：

```
1. 先返回标志和命令候选（快速）
2. 异步收集文件候选
3. 文件收集完成后，刷新补全列表
```

---

## 六、测试用例

### 6.1 渐进式补全场景

```bash
# 场景1: 空输入（全量展示）
$ myapp <Tab>
Bash 显示:    （所有子命令 + 所有标志 + 当前目录文件）
               build, run, test, --config, --verbose, config.json, main.go, ...
PowerShell:   循环遍历所有候选

# 场景2: 短输入（快速收敛）
$ myapp c<Tab>
Bash 显示:    config.json, --config, --conf-level, configure, connect
               （最多 20 个相关候选）
PowerShell:   跳到下一个匹配 "c" 的候选

# 场景3: 较长输入（精准定位）
$ myapp conf<Tab>
Bash 显示:    config.json, --config, configure
               （最多 10 个高质量匹配）
PowerShell:   跳到下一个匹配 "conf" 的候选

# 场景4: 精确输入（唯一匹配）
$ myapp confi<Tab>
Bash 显示:    config.json
PowerShell:   自动填充 "config.json"
```

### 6.2 其他场景

```bash
# 场景5: 输入标志前缀  
$ myapp --conf<Tab>
期望: --config, --conf-level 显示（文件不显示，因为以 - 开头）

# 场景6: 枚举标志后
$ myapp --format js<Tab>
期望: json, yaml, xml 显示（枚举值），同时显示 js* 文件

# 场景7: 路径特征输入
$ myapp ./conf<Tab>
期望: 直接补全 ./config.json, ./config.yaml（专用路径补全）

# 场景8: 子命令补全
$ myapp bu<Tab>
期望: build, bundle 等子命令显示

# 场景9: 大目录
$ myapp a<Tab>  # 当前目录有1000+文件
期望: 快速返回前20个匹配文件 + 标志 + 命令（短输入限制）

# 场景10: 无匹配
$ myapp xyz<Tab>
期望: 返回空列表，不报错
```

---

## 七、实施计划

### 7.1 阶段划分

| 阶段 | 内容 | 预计时间 |
|------|------|---------|
| Phase 1 | 修改 Bash 模板，实现基础统一补全 | 2h |
| Phase 2 | 修改 PowerShell 模板，实现基础统一补全 | 2h |
| Phase 3 | 优化评分算法，调整类型加成 | 1h |
| Phase 4 | 性能优化，添加缓存 | 1h |
| Phase 5 | 测试验证 | 2h |

### 7.2 回滚方案

保留原始模板文件，如果新方案有问题可以快速回滚：

```bash
# 备份
bash.tmpl.backup
pwsh.tmpl.backup
```

---

## 八、附录

### 8.1 术语表

| 术语 | 说明 |
|------|------|
| 统一补全池 | 将文件、标志、命令等所有候选放入同一个集合 |
| 模糊匹配 | 不完全精确匹配，允许字符跳跃、子串等 |
| 上下文 | 当前命令的路径，如 `/` 或 `/build/` |
| 候选 | 可能被补全的选项 |

### 8.2 参考实现

- [fzf](https://github.com/junegunn/fzf): 模糊匹配算法参考
- [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions): 智能补全参考

---

*文档结束*
