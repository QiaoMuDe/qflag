package qflag

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"testing"
)

// TestGenerateHelpInfo_BasicCommand 测试基本命令的帮助信息生成
func TestGenerateHelpInfo_BasicCommand(t *testing.T) {
	cmd := NewCmd("testcmd", "tc", flag.ContinueOnError)
	cmd.SetUseChinese(true)

	helpInfo := generateHelpInfo(cmd)

	// 验证命令名称和描述
	if !strings.Contains(helpInfo, "testcmd(tc)") {
		t.Errorf("帮助信息未包含命令名称, 实际输出: %s", helpInfo)
	}
}

// TestGenerateHelpInfo_WithOptions 测试带选项的命令帮助信息
func TestGenerateHelpInfo_WithOptions(t *testing.T) {
	cmd := NewCmd("testcmd", "tc", flag.ContinueOnError)
	cmd.String("config", "c", "/etc/config.json", "配置文件路径")

	helpInfo := generateHelpInfo(cmd)

	// 验证选项部分
	if !strings.Contains(helpInfo, "--config, -c") {
		t.Errorf("帮助信息未包含选项, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "配置文件路径") {
		t.Errorf("帮助信息未包含选项描述, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "/etc/config.json") {
		t.Errorf("帮助信息未包含默认值, 实际输出: %s", helpInfo)
	}
}

// TestGenerateHelpInfo_WithSubCommands 测试带子命令的帮助信息
func TestGenerateHelpInfo_WithSubCommands(t *testing.T) {
	cmd := NewCmd("parent", "p", flag.ContinueOnError)
	subCmd1 := NewCmd("child1", "c1", flag.ContinueOnError)
	subCmd1.SetDescription("First child command")
	subCmd2 := NewCmd("child2", "", flag.ContinueOnError)
	subCmd2.SetDescription("Second child command without short name")

	_ = cmd.AddSubCmd(subCmd1, subCmd2)
	cmd.SetUseChinese(true)
	helpInfo := generateHelpInfo(cmd)

	// 验证子命令部分
	if !strings.Contains(helpInfo, "子命令:") {
		t.Errorf("帮助信息未包含子命令标题, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "child1, c1") {
		t.Errorf("帮助信息未包含带短名称的子命令, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "child2") {
		t.Errorf("帮助信息未包含无子名称的子命令, 实际输出: %s", helpInfo)
	}
}

// TestGenerateHelpInfo_WithExamples 测试带示例的命令帮助信息
func TestGenerateHelpInfo_WithExamples(t *testing.T) {
	cmd := NewCmd("testcmd", "tc", flag.ContinueOnError)
	cmd.SetUseChinese(true)

	cmd.AddExample(ExampleInfo{
		Description: "基本用法",
		Usage:       "testcmd --config /custom.json",
	})
	cmd.AddExample(ExampleInfo{
		Description: "详细输出",
		Usage:       "testcmd -v",
	})

	helpInfo := generateHelpInfo(cmd)

	// 当使用-v选项运行测试时打印生成的帮助信息
	// 在较旧版本的 Go 中，t.Verbose() 方法不存在，可以通过获取测试标志 -test.v 的值来判断是否开启详细输出
	if testing.Verbose() {
		fmt.Println(helpInfo)
	}

	// 验证示例部分
	if !strings.Contains(helpInfo, "示例:") {
		t.Errorf("帮助信息未包含示例标题, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "1、基本用法") {
		t.Errorf("帮助信息未包含第一个示例描述, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "testcmd --config /custom.json") {
		t.Errorf("帮助信息未包含第一个示例用法, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "2、详细输出") {
		t.Errorf("帮助信息未包含第二个示例描述, 实际输出: %s", helpInfo)
	}
}

// TestGenerateHelpInfo_EnglishLanguage 测试英文环境下的帮助信息
func TestGenerateHelpInfo_EnglishLanguage(t *testing.T) {
	cmd := NewCmd("testcmd", "tc", flag.ContinueOnError)
	cmd.SetUseChinese(false)
	cmd.SetDescription("English test command")
	cmd.AddNote("Important note for English users")

	helpInfo := generateHelpInfo(cmd)

	// 验证英文模板内容
	if !strings.Contains(helpInfo, "Name: testcmd(tc)") {
		t.Errorf("帮助信息未包含英文名称, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "Desc: English test command") {
		t.Errorf("帮助信息未包含英文描述, 实际输出: %s", helpInfo)
	}
	if !strings.Contains(helpInfo, "Notes:") {
		t.Errorf("帮助信息未包含英文注意事项标题, 实际输出: %s", helpInfo)
	}
}

// TestSortWithShortNamePriority 测试子命令排序逻辑
func TestSortWithShortNamePriority(t *testing.T) {
	// 创建测试用例: 有短名称的应排在前面, 按长名称字母序排列
	subCmds := []*Cmd{
		NewCmd("banana", "b", flag.ContinueOnError),
		NewCmd("apple", "a", flag.ContinueOnError),
		NewCmd("cherry", "", flag.ContinueOnError),
	}

	// 执行排序
	sortedSubCmds := make([]*Cmd, len(subCmds))
	copy(sortedSubCmds, subCmds)
	sort.Slice(sortedSubCmds, func(i, j int) bool {
		a, b := sortedSubCmds[i], sortedSubCmds[j]
		return sortWithShortNamePriority(
			a.ShortName() != "",
			b.ShortName() != "",
			a.LongName(),
			b.LongName(),
			a.ShortName(),
			b.ShortName(),
		)
	})

	// 验证排序结果: apple(a) -> banana(b) -> cherry
	if sortedSubCmds[0].LongName() != "apple" {
		t.Errorf("排序错误, 第一个子命令应为apple, 实际为%s", sortedSubCmds[0].LongName())
	}
	if sortedSubCmds[1].LongName() != "banana" {
		t.Errorf("排序错误, 第二个子命令应为banana, 实际为%s", sortedSubCmds[1].LongName())
	}
	if sortedSubCmds[2].LongName() != "cherry" {
		t.Errorf("排序错误, 第三个子命令应为cherry, 实际为%s", sortedSubCmds[2].LongName())
	}
}

// TestSetLogoTextAndModuleHelps 测试设置Logo文本和自定义模块帮助信息
func TestSetLogoTextAndModuleHelps(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	cmd.SetUseChinese(true)

	loggo := `________      ________          ___  __       
|\  _____\    |\   ____\        |\  \|\  \     
\ \  \__/     \ \  \___|        \ \  \/  /|_   
 \ \   __\     \ \  \            \ \   ___  \  
  \ \  \_|      \ \  \____        \ \  \\ \  \ 
   \ \__\        \ \_______\       \ \__\\ \__\
    \|__|         \|_______|        \|__| \|__|
                FCK CLI Test Logo Text               
`

	cmd.SetLogoText(loggo)

	cmd.SetModuleHelps("testMode:\n\tThis is a test module helps\t测试")

	helpInfo := generateHelpInfo(cmd)
	// 如果是-v运行测试，则打印帮助信息
	if testing.Verbose() {
		fmt.Println(helpInfo)
	}

	// 验证Logo文本
	if !strings.Contains(helpInfo, "Test Logo Text") {
		t.Errorf("帮助信息未包含Logo文本, 实际输出: %s", helpInfo)
	}
}
