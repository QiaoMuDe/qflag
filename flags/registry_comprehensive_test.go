// Package flags 标志注册表综合测试
// 本文件包含了标志注册表的全面测试用例，测试复杂场景下的
// 标志注册、冲突检测、批量操作等高级功能的正确性。
package flags

import (
	"sync"
	"testing"
)

// TestFlagMeta 测试FlagMeta结构体的功能
func TestFlagMeta(t *testing.T) {
	// 创建一个测试标志
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			longName:     "test",
			shortName:    "t",
			initialValue: "default",
			usage:        "测试标志",
			value:        new(string),
			initialized:  true,
		},
	}

	meta := &FlagMeta{Flag: flag}

	t.Run("获取长名称", func(t *testing.T) {
		if meta.GetLongName() != "test" {
			t.Errorf("期望长名称为 'test'，实际为 '%s'", meta.GetLongName())
		}
	})

	t.Run("获取短名称", func(t *testing.T) {
		if meta.GetShortName() != "t" {
			t.Errorf("期望短名称为 't'，实际为 '%s'", meta.GetShortName())
		}
	})

	t.Run("获取名称优先长名称", func(t *testing.T) {
		if meta.GetName() != "test" {
			t.Errorf("期望名称为 'test'，实际为 '%s'", meta.GetName())
		}
	})

	t.Run("获取用法说明", func(t *testing.T) {
		if meta.GetUsage() != "测试标志" {
			t.Errorf("期望用法说明为 '测试标志'，实际为 '%s'", meta.GetUsage())
		}
	})

	t.Run("获取标志类型", func(t *testing.T) {
		if meta.GetFlagType() != FlagTypeString {
			t.Errorf("期望标志类型为 %d，实际为 %d", FlagTypeString, meta.GetFlagType())
		}
	})

	t.Run("获取默认值", func(t *testing.T) {
		if meta.GetDefault() != "default" {
			t.Errorf("期望默认值为 'default'，实际为 '%v'", meta.GetDefault())
		}
	})

	t.Run("获取标志对象", func(t *testing.T) {
		if meta.GetFlag() != flag {
			t.Error("GetFlag()应该返回原始标志对象")
		}
	})
}

// TestFlagMeta_OnlyShortName 测试只有短名称的FlagMeta
func TestFlagMeta_OnlyShortName(t *testing.T) {
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			longName:     "",
			shortName:    "s",
			initialValue: "default",
			usage:        "短名称标志",
			value:        new(string),
			initialized:  true,
		},
	}

	meta := &FlagMeta{Flag: flag}

	if meta.GetName() != "s" {
		t.Errorf("只有短名称时GetName()应返回 's'，实际为 '%s'", meta.GetName())
	}
}

// TestNewFlagRegistry 测试创建新的标志注册表
func TestNewFlagRegistry(t *testing.T) {
	registry := NewFlagRegistry()

	if registry == nil {
		t.Fatal("NewFlagRegistry()不应返回nil")
	}

	if registry.GetFlagMetaCount() != 0 {
		t.Errorf("新注册表应该为空，实际有 %d 个标志", registry.GetFlagMetaCount())
	}

	if registry.GetLongFlagsCount() != 0 {
		t.Errorf("新注册表长标志数量应为0，实际为 %d", registry.GetLongFlagsCount())
	}

	if registry.GetShortFlagsCount() != 0 {
		t.Errorf("新注册表短标志数量应为0，实际为 %d", registry.GetShortFlagsCount())
	}
}

// TestFlagRegistry_RegisterFlag 测试标志注册功能
func TestFlagRegistry_RegisterFlag(t *testing.T) {
	registry := NewFlagRegistry()

	t.Run("正常注册标志", func(t *testing.T) {
		flag := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "test",
				shortName:    "t",
				initialValue: "default",
				usage:        "测试标志",
				value:        new(string),
				initialized:  true,
			},
		}
		meta := &FlagMeta{Flag: flag}

		err := registry.RegisterFlag(meta)
		if err != nil {
			t.Fatalf("注册标志失败: %v", err)
		}

		if registry.GetFlagMetaCount() != 1 {
			t.Errorf("注册后应有1个标志，实际有 %d 个", registry.GetFlagMetaCount())
		}
	})

	t.Run("注册重复长名称", func(t *testing.T) {
		flag1 := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "duplicate",
				shortName:    "d1",
				initialValue: "default1",
				usage:        "重复标志1",
				value:        new(string),
				initialized:  true,
			},
		}
		flag2 := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "duplicate",
				shortName:    "d2",
				initialValue: "default2",
				usage:        "重复标志2",
				value:        new(string),
				initialized:  true,
			},
		}

		registry := NewFlagRegistry()
		err := registry.RegisterFlag(&FlagMeta{Flag: flag1})
		if err != nil {
			t.Fatalf("第一次注册失败: %v", err)
		}

		err = registry.RegisterFlag(&FlagMeta{Flag: flag2})
		if err == nil {
			t.Error("注册重复长名称应该返回错误")
		}
	})

	t.Run("注册重复短名称", func(t *testing.T) {
		flag1 := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "first",
				shortName:    "f",
				initialValue: "default1",
				usage:        "第一个标志",
				value:        new(string),
				initialized:  true,
			},
		}
		flag2 := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "second",
				shortName:    "f",
				initialValue: "default2",
				usage:        "第二个标志",
				value:        new(string),
				initialized:  true,
			},
		}

		registry := NewFlagRegistry()
		err := registry.RegisterFlag(&FlagMeta{Flag: flag1})
		if err != nil {
			t.Fatalf("第一次注册失败: %v", err)
		}

		err = registry.RegisterFlag(&FlagMeta{Flag: flag2})
		if err == nil {
			t.Error("注册重复短名称应该返回错误")
		}
	})

	t.Run("注册空名称标志", func(t *testing.T) {
		flag := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "",
				shortName:    "",
				initialValue: "default",
				usage:        "空名称标志",
				value:        new(string),
				initialized:  true,
			},
		}

		registry := NewFlagRegistry()
		err := registry.RegisterFlag(&FlagMeta{Flag: flag})
		if err == nil {
			t.Error("注册空名称标志应该返回错误")
		}
	})

	t.Run("注册包含非法字符的标志", func(t *testing.T) {
		flag := &StringFlag{
			BaseFlag: BaseFlag[string]{
				longName:     "test@flag",
				shortName:    "t",
				initialValue: "default",
				usage:        "非法字符标志",
				value:        new(string),
				initialized:  true,
			},
		}

		registry := NewFlagRegistry()
		err := registry.RegisterFlag(&FlagMeta{Flag: flag})
		if err == nil {
			t.Error("注册包含非法字符的标志应该返回错误")
		}
	})
}

// TestFlagRegistry_GetMethods 测试标志查找方法
func TestFlagRegistry_GetMethods(t *testing.T) {
	registry := NewFlagRegistry()

	// 注册测试标志
	flag := &StringFlag{
		BaseFlag: BaseFlag[string]{
			longName:     "verbose",
			shortName:    "v",
			initialValue: "false",
			usage:        "详细输出",
			value:        new(string),
			initialized:  true,
		},
	}
	meta := &FlagMeta{Flag: flag}
	err := registry.RegisterFlag(meta)
	if err != nil {
		t.Fatalf("注册标志失败: %v", err)
	}

	t.Run("通过长名称查找", func(t *testing.T) {
		foundMeta, exists := registry.GetByLong("verbose")
		if !exists {
			t.Error("应该能通过长名称找到标志")
		}
		if foundMeta != meta {
			t.Error("找到的标志元数据应该与注册的相同")
		}

		_, exists = registry.GetByLong("nonexistent")
		if exists {
			t.Error("不存在的长名称应该返回false")
		}
	})

	t.Run("通过短名称查找", func(t *testing.T) {
		foundMeta, exists := registry.GetByShort("v")
		if !exists {
			t.Error("应该能通过短名称找到标志")
		}
		if foundMeta != meta {
			t.Error("找到的标志元数据应该与注册的相同")
		}

		_, exists = registry.GetByShort("x")
		if exists {
			t.Error("不存在的短名称应该返回false")
		}
	})

	t.Run("通过名称查找", func(t *testing.T) {
		// 测试长名称查找
		foundMeta, exists := registry.GetByName("verbose")
		if !exists {
			t.Error("应该能通过长名称找到标志")
		}
		if foundMeta != meta {
			t.Error("找到的标志元数据应该与注册的相同")
		}

		// 测试短名称查找
		foundMeta, exists = registry.GetByName("v")
		if !exists {
			t.Error("应该能通过短名称找到标志")
		}
		if foundMeta != meta {
			t.Error("找到的标志元数据应该与注册的相同")
		}

		// 测试不存在的名称
		_, exists = registry.GetByName("nonexistent")
		if exists {
			t.Error("不存在的名称应该返回false")
		}
	})
}

// TestFlagRegistry_GetCollections 测试获取集合的方法
func TestFlagRegistry_GetCollections(t *testing.T) {
	registry := NewFlagRegistry()

	// 注册多个测试标志
	flags := []*FlagMeta{
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "help", shortName: "h", usage: "帮助", value: new(string), initialized: true}}},
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "version", shortName: "v", usage: "版本", value: new(string), initialized: true}}},
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "output", shortName: "", usage: "输出", value: new(string), initialized: true}}},
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "", shortName: "q", usage: "安静模式", value: new(string), initialized: true}}},
	}

	for _, flag := range flags {
		err := registry.RegisterFlag(flag)
		if err != nil {
			t.Fatalf("注册标志失败: %v", err)
		}
	}

	t.Run("获取标志元数据列表", func(t *testing.T) {
		metaList := registry.GetFlagMetaList()
		if len(metaList) != 4 {
			t.Errorf("应该有4个标志元数据，实际有 %d 个", len(metaList))
		}
	})

	t.Run("获取所有标志映射", func(t *testing.T) {
		allFlags := registry.GetFlagNameMap()
		expectedCount := 6 // help, h, version, v, output, q
		if len(allFlags) != expectedCount {
			t.Errorf("应该有 %d 个标志映射，实际有 %d 个", expectedCount, len(allFlags))
		}

		// 验证特定标志存在
		if _, exists := allFlags["help"]; !exists {
			t.Error("应该包含 'help' 标志")
		}
		if _, exists := allFlags["h"]; !exists {
			t.Error("应该包含 'h' 标志")
		}
	})

	t.Run("获取长标志映射", func(t *testing.T) {
		longFlags := registry.GetLongFlagMap()
		expectedCount := 3 // help, version, output
		if len(longFlags) != expectedCount {
			t.Errorf("应该有 %d 个长标志，实际有 %d 个", expectedCount, len(longFlags))
		}

		if _, exists := longFlags["help"]; !exists {
			t.Error("长标志映射应该包含 'help'")
		}
		if _, exists := longFlags["h"]; exists {
			t.Error("长标志映射不应该包含短名称 'h'")
		}
	})

	t.Run("获取短标志映射", func(t *testing.T) {
		shortFlags := registry.GetShortFlagMap()
		expectedCount := 3 // h, v, q
		if len(shortFlags) != expectedCount {
			t.Errorf("应该有 %d 个短标志，实际有 %d 个", expectedCount, len(shortFlags))
		}

		if _, exists := shortFlags["h"]; !exists {
			t.Error("短标志映射应该包含 'h'")
		}
		if _, exists := shortFlags["help"]; exists {
			t.Error("短标志映射不应该包含长名称 'help'")
		}
	})
}

// TestFlagRegistry_CountMethods 测试计数方法
func TestFlagRegistry_CountMethods(t *testing.T) {
	registry := NewFlagRegistry()

	// 注册测试标志
	flags := []*FlagMeta{
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "help", shortName: "h", usage: "帮助", value: new(string), initialized: true}}},
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "version", shortName: "", usage: "版本", value: new(string), initialized: true}}},
		{Flag: &StringFlag{BaseFlag: BaseFlag[string]{longName: "", shortName: "q", usage: "安静模式", value: new(string), initialized: true}}},
	}

	for _, flag := range flags {
		err := registry.RegisterFlag(flag)
		if err != nil {
			t.Fatalf("注册标志失败: %v", err)
		}
	}

	t.Run("标志元数据计数", func(t *testing.T) {
		if registry.GetFlagMetaCount() != 3 {
			t.Errorf("应该有3个标志元数据，实际有 %d 个", registry.GetFlagMetaCount())
		}
	})

	t.Run("长标志计数", func(t *testing.T) {
		if registry.GetLongFlagsCount() != 2 {
			t.Errorf("应该有2个长标志，实际有 %d 个", registry.GetLongFlagsCount())
		}
	})

	t.Run("短标志计数", func(t *testing.T) {
		if registry.GetShortFlagsCount() != 2 {
			t.Errorf("应该有2个短标志，实际有 %d 个", registry.GetShortFlagsCount())
		}
	})

	t.Run("所有标志计数", func(t *testing.T) {
		if registry.GetAllFlagsCount() != 4 {
			t.Errorf("应该有4个标志（长+短），实际有 %d 个", registry.GetAllFlagsCount())
		}
	})
}

// TestFlagRegistry_ConcurrentAccess 测试并发访问安全性
func TestFlagRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewFlagRegistry()
	var wg sync.WaitGroup
	numGoroutines := 50

	// 并发注册标志
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			flag := &StringFlag{
				BaseFlag: BaseFlag[string]{
					longName:     "flag" + string(rune('0'+index%10)),
					shortName:    string(rune('a' + index%26)),
					initialValue: "default",
					usage:        "并发测试标志",
					value:        new(string),
					initialized:  true,
				},
			}
			meta := &FlagMeta{Flag: flag}
			_ = registry.RegisterFlag(meta) // 可能会因为重复名称而失败，这是正常的
		}(i)
	}

	// 并发读取标志
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			registry.GetByLong("flag" + string(rune('0'+index%10)))
			registry.GetByShort(string(rune('a' + index%26)))
			registry.GetFlagMetaList()
			registry.GetFlagMetaCount()
		}(i)
	}

	wg.Wait()
	// 如果没有panic或死锁，说明并发访问是安全的
}

// TestFlagRegistry_ValidateFlagName 测试标志名称验证
func TestFlagRegistry_ValidateFlagName(t *testing.T) {
	registry := NewFlagRegistry()

	testCases := []struct {
		name        string
		longName    string
		shortName   string
		shouldError bool
		description string
	}{
		{"正常标志", "normal", "n", false, "正常的标志名称应该通过验证"},
		{"包含空格", "with space", "w", true, "包含空格的标志名称应该失败"},
		{"包含特殊字符", "with@symbol", "w", true, "包含@符号的标志名称应该失败"},
		{"包含感叹号", "with!", "w", true, "包含感叹号的标志名称应该失败"},
		{"数字开头", "123flag", "1", false, "数字开头的标志名称应该允许"},
		{"下划线", "with_underscore", "_", false, "包含下划线的标志名称应该允许"},
		{"连字符", "with-dash", "-", false, "包含连字符的标志名称应该允许"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag := &StringFlag{
				BaseFlag: BaseFlag[string]{
					longName:     tc.longName,
					shortName:    tc.shortName,
					initialValue: "default",
					usage:        tc.description,
					value:        new(string),
					initialized:  true,
				},
			}
			meta := &FlagMeta{Flag: flag}

			err := registry.RegisterFlag(meta)
			if tc.shouldError && err == nil {
				t.Errorf("%s: 应该返回错误但没有", tc.description)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("%s: 不应该返回错误但返回了: %v", tc.description, err)
			}
		})
	}
}
