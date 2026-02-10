# 帮助信息生成器设计方案

## 1. 概述

帮助信息生成器用于根据命令配置自动生成格式良好的帮助文档, 支持中英文切换、自动对齐、模板渲染等功能。

## 2. 设计原则

本设计遵循以下原则: 
- **直接使用现有接口**: 直接从 Cmd 和 CmdConfig 读取数据, 不创建中间数据拷贝
- **最小化新类型**: 仅定义 1-2 个必要类型, 避免过度抽象
- **配置简洁**: 只保留真正需要的配置选项
- **复用现有数据结构**: Examples 和 Notes 直接复用 CmdConfig 中的定义

## 3. 数据结构

### 3.1 HelpGenerator 帮助生成器接口

```go
// HelpGenerator 帮助生成器接口
type HelpGenerator interface {
    // Generate 生成帮助信息
    Generate(cmd types.Cmd) string
}
```

### 3.2 TextGenerator 文本帮助生成器

```go
// TextGenerator 文本帮助生成器
type TextGenerator struct {
    templates map[string]string  // 模板缓存
    language  string             // 当前语言: "zh" | "en"
    showLogo  bool                // 是否显示 Logo
}
```

**配置说明**: 

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| templates | map[string]string | 内置中英文模板 | 模板缓存 |
| language | string | "zh" | 当前语言 |
| showLogo | bool | true | 是否显示 Logo |

## 4. 模板设计

### 4.1 中文模板 (helpTemplateCN)

```go
const helpTemplateCN = `
名称:
    {{.Name}}

描述:
    {{.Description}}

{{if and .ShowLogo .LogoText}}
{{.LogoText}}
{{end}}

用法:
    {{.UsageSyntax}}

选项:
{{range .Flags}}
    {{.Name}}  {{.Description}}{{if .DefaultValue}} (默认值: {{.DefaultValue}}){{end}}
{{end}}
{{if .SubCmds}}
子命令:
{{range .SubCmds}}
    {{.Name}}  {{.Description}}
{{end}}
{{end}}
{{if .Examples}}
示例:
{{range $title, $content := .Examples}}
    {{$title}}
       {{$content}}
{{end}}
{{end}}
{{if .Notes}}
注意事项:
{{range $index, $note := .Notes}}
    {{add $index 1}}、{{$note}}
{{end}}
{{end}}
`
```

### 4.2 英文模板 (helpTemplateEN)

```go
const helpTemplateEN = `
Name:
    {{.Name}}

Description:
    {{.Description}}

{{if and .ShowLogo .LogoText}}
{{.LogoText}}
{{end}}

Usage:
    {{.UsageSyntax}}

Options:
{{range .Flags}}
    {{.Name}}  {{.Description}}{{if .DefaultValue}} (default: {{.DefaultValue}}){{end}}
{{end}}
{{if .SubCmds}}
SubCmds:
{{range .SubCmds}}
    {{.Name}}  {{.Description}}
{{end}}
{{end}}
{{if .Examples}}
Examples:
{{range $title, $content := .Examples}}
    {{$title}}
       {{$content}}
{{end}}
{{end}}
{{if .Notes}}
Notes:
{{range $index, $note := .Notes}}
    {{add $index 1}}. {{$note}}
{{end}}
{{end}}
`
```

### 4.3 模板数据来源

模板直接使用 Cmd 和 CmdConfig 中的数据, 无需创建中间类型: 

| 模板变量 | 数据来源 | 说明 |
|----------|----------|------|
| .Name | cmd.Name() | 命令名称 |
| .Description | cmd.Desc() | 命令描述 |
| .LogoText | cmd.Config().LogoText | Logo 文本 |
| .ShowLogo | generator.showLogo | 是否显示 Logo |
| .UsageSyntax | cmd.Config().UsageSyntax | 用法语法 |
| .Flags | cmd.Flags() | 标志列表 |
| .SubCmds | cmd.SubCmds() | 子命令列表 |
| .Examples | cmd.Config().Example | 示例列表 |
| .Notes | cmd.Config().Notes | 注意事项列表 |

## 5. 核心实现

### 5.1 创建帮助生成器

```go
// NewTextGenerator 创建文本帮助生成器
func NewTextGenerator() HelpGenerator {
    return &TextGenerator{
        templates: map[string]string{
            "zh": helpTemplateCN,
            "en": helpTemplateEN,
        },
        language:  "zh",
        showLogo: true,
    }
}
```

### 5.2 生成帮助信息

```go
// Generate 生成帮助信息
func (g *TextGenerator) Generate(cmd types.Cmd) string {
    data := map[string]interface{}{
        "Name":        cmd.Name(),
        "Description": cmd.Desc(),
        "UsageSyntax": cmd.Config().UsageSyntax,
        "LogoText":    cmd.Config().LogoText,
        "ShowLogo":    g.showLogo,
        "Language":    g.language,
    }

    // 处理标志列表, 添加名称格式化
    flags := cmd.Flags()
    formattedFlags := make([]FormattedFlag, len(flags))
    for i, f := range flags {
        formattedFlags[i] = FormattedFlag{
            Flag:         f,
            Name:         formatFlagName(f),
            Description:  f.Desc(),
            DefaultValue: formatDefaultValue(f),
        }
    }
    data["Flags"] = formattedFlags

    // 处理子命令列表
    SubCmds := cmd.SubCmds()
    formattedSubCmds := make([]FormattedSubCmd, len(SubCmds))
    for i, sub := range SubCmds {
        formattedSubCmds[i] = FormattedSubCmd{
            Cmd:    sub,
            Name:       formatSubCmdName(sub),
            Description: sub.Desc(),
        }
    }
    data["SubCmds"] = formattedSubCmds

    // 复用 CmdConfig 中的数据
    data["Examples"] = cmd.Config().Example
    data["Notes"] = cmd.Config().Notes

    return g.renderTemplate(data)
}
```

### 5.3 格式化辅助结构

```go
// FormattedFlag 已格式化的标志信息
type FormattedFlag struct {
    types.Flag
    Name         string // 格式化后的名称
    Description  string
    DefaultValue string
}

// FormattedSubCmd 已格式化的子命令信息
type FormattedSubCmd struct {
    types.Cmd
    Name        string // 格式化后的名称
    Description string
}
```

这两个类型仅用于模板渲染时的数据格式化, 不影响核心数据流。

### 5.4 格式辅助函数

```go
// formatFlagName 格式化选项名称
func formatFlagName(f types.Flag) string {
    if f.ShortName() != "" && f.LongName() != "" {
        return fmt.Sprintf("-%s, --%s", f.ShortName(), f.LongName())
    } else if f.LongName() != "" {
        return fmt.Sprintf("--%s", f.LongName())
    }
    return fmt.Sprintf("-%s", f.ShortName())
}

// formatSubCmdName 格式化子命令名称
func formatSubCmdName(cmd types.Cmd) string {
    if cmd.ShortName() != "" {
        return fmt.Sprintf("%s, %s", cmd.LongName(), cmd.ShortName())
    }
    return cmd.LongName()
}

// formatDefaultValue 格式化默认值
func formatDefaultValue(f types.Flag) string {
    def := f.GetDefault()
    if def == nil {
        return ""
    }
    return fmt.Sprintf("%v", def)
}
```

### 5.5 模板渲染

```go
// renderTemplate 渲染模板
func (g *TextGenerator) renderTemplate(data map[string]interface{}) string {
    tmpl, ok := g.templates[g.language]
    if !ok {
        tmpl = g.templates["zh"]
    }

    funcMap := template.FuncMap{
        "add": func(a, b int) int { return a + b },
    }

    t := template.New("help").Funcs(funcMap)
    t = template.Must(t.Parse(tmpl))

    var buf bytes.Buffer
    if err := t.Execute(&buf, data); err != nil {
        return fmt.Sprintf("Error generating help: %v", err)
    }

    return strings.TrimSpace(buf.String())
}
```

### 5.6 配置方法

```go
// SetLanguage 设置语言
func (g *TextGenerator) SetLanguage(lang string) {
    g.language = lang
}

// SetShowLogo 设置是否显示 Logo
func (g *TextGenerator) SetShowLogo(show bool) {
    g.showLogo = show
}
```

## 6. 与 Cmd 集成

### 6.1 修改 Cmd

```go
type Cmd struct {
    // ... 现有字段
    helpGen help.HelpGenerator // 添加帮助生成器
}

// 修改 NewCmd
func NewCmd(longName, shortName string) *Cmd {
    return &Cmd{
        // ... 其他初始化
        helpGen: help.NewTextGenerator(),
    }
}

// 修改 Help 方法
func (c *Cmd) Help() string {
    return c.helpGen.Generate(c)
}

// 添加设置帮助生成器的方法
func (c *Cmd) SetHelpGenerator(helpGen help.HelpGenerator) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.helpGen = helpGen
}
```

## 7. 使用示例

### 7.1 基本使用

```go
cmd := cmd.NewCmd("fck", "")
cmd.SetDesc("多功能文件处理工具集")
cmd.SetUsageSyntax("fck [全局选项] [子命令] [选项]")
cmd.SetLogoText(`
     ________      ________          ___  __
    |\  _____\    |\   ____\        |\  \|\  \
    \ \  \__/     \ \  \___|        \ \  \/  /|_
     \ \   __\     \ \  \            \ \   ___  \
      \ \  \_|      \ \  \____        \ \  \\ \  \
       \ \__\        \ \_______\       \ \__\\ \__\
        \|__|         \|_______|        \|__| \|__|`)

cmd.AddExample("Windows 临时启用",
    `D:\AppData\GoPath\bin\fck.exe --completion powershell | Out-String | Invoke-Expression`)

cmd.AddNote("各子命令有独立帮助文档, 可通过-h参数查看")

fmt.Println(cmd.Help())
```

### 7.2 自定义配置

```go
helpGen := help.NewTextGenerator()
helpGen.SetLanguage("en")
helpGen.SetShowLogo(false)

cmd.SetHelpGenerator(helpGen)
```

## 8. 输出示例

### 8.1 中文输出

```
名称:
    fck.exe

描述:
    多功能文件处理工具集, 提供文件哈希计算、大小统计、查找和校验等实用功能

     ________      ________          ___  __
    |\  _____\    |\   ____\        |\  \|\  \
    \ \  \__/     \ \  \___|        \ \  \/  /|_
     \ \   __\     \ \  \            \ \   ___  \
      \ \  \_|      \ \  \____        \ \  \\ \  \
       \ \__\        \ \_______\       \ \__\\ \__\
        \|__|         \|_______|        \|__| \|__|

用法:
    fck.exe [全局选项] [子命令] [选项]

选项:
    -h, --help  Show help (默认值: false)
    -v, --version  Show version (默认值: false)

子命令:
    check, c  文件校验工具, 对比指定目录A和目录B的文件差异...
    find, f  文件目录查找工具, 在指定目录及其子目录中...

示例:
    Windows 临时启用
       D:\AppData\GoPath\bin\fck.exe --completion powershell | Out-String | Invoke-Expression

注意事项:
    1、各子命令有独立帮助文档, 可通过-h参数查看...
    2、所有路径参数支持Windows和Unix风格
```

### 8.2 英文输出

```
Name:
    fck.exe

Description:
    Multi-functional file processing toolkit...

Usage:
    fck.exe [global options] [subcmd] [options]

Options:
    -h, --help  Show help (default: false)
    -v, --version  Show version (default: false)

SubCmds:
    check, c  File checksum tool...
    find, f  File search tool...

Examples:
    Temporary enable on Windows
       D:\AppData\GoPath\bin\fck.exe --completion powershell | Out-String | Invoke-Expression

Notes:
    1. Each subcmd has its own help documentation...
```

## 9. 文件结构

```
internal/help/
├── help.go            # 接口定义和模板
└── README.md          # 本文档
```

## 10. 设计对比

### 10.1 类型数量对比

| 方面 | 原设计 | 简化后 |
|------|--------|--------|
| 新增类型数量 | 7 个 | 2 个 |
| 配置字段数 | 8 个 | 3 个 |
| 数据重新封装 | 是 | 否 |

### 10.2 简化要点

1. **移除 HelpData**: 直接使用 Cmd 和 CmdConfig 的数据
2. **移除 HelpOption**: 模板中直接遍历 Flags
3. **移除 HelpSubCmd**: 模板中直接遍历 SubCmds
4. **复用 CmdConfig**: Examples 和 Notes 直接复用
5. **简化 HelpConfig**: 只保留 language、showLogo 等必要配置

## 11. 后续优化

- [ ] 支持 Markdown 格式输出
- [ ] 支持彩色终端输出
- [ ] 支持自定义模板
- [ ] 支持按宽度自动换行
