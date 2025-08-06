package flags

import (
	"reflect"
	"testing"
	"time"
)

// TestMapFlag_AdvancedFeatures 测试MapFlag的高级功能
func TestMapFlag_AdvancedFeatures(t *testing.T) {
	t.Run("复杂键值对解析", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(",", "=")

		// 测试包含特殊字符的键值对
		err := flag.Set("key1=value with spaces,key2=value=with=equals")
		if err != nil {
			t.Fatalf("设置复杂键值对失败: %v", err)
		}

		result := flag.Get()
		if result["key1"] != "value with spaces" {
			t.Errorf("key1的值应为 'value with spaces'，实际为 '%s'", result["key1"])
		}
		if result["key2"] != "value=with=equals" {
			t.Errorf("key2的值应为 'value=with=equals'，实际为 '%s'", result["key2"])
		}
	})

	t.Run("多次设置累积效果", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(",", "=")

		// 第一次设置
		err := flag.Set("key1=value1,key2=value2")
		if err != nil {
			t.Fatalf("第一次设置失败: %v", err)
		}

		// 第二次设置，应该累积
		err = flag.Set("key3=value3,key1=newvalue1")
		if err != nil {
			t.Fatalf("第二次设置失败: %v", err)
		}

		result := flag.Get()
		if len(result) != 3 {
			t.Errorf("应该有3个键值对，实际有 %d 个", len(result))
		}

		if result["key1"] != "newvalue1" {
			t.Errorf("key1应该被更新为 'newvalue1'，实际为 '%s'", result["key1"])
		}
	})

	t.Run("自定义分隔符组合", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters("||", "::")

		err := flag.Set("name::张三||age::25||city::北京")
		if err != nil {
			t.Fatalf("自定义分隔符设置失败: %v", err)
		}

		result := flag.Get()
		expected := map[string]string{
			"name": "张三",
			"age":  "25",
			"city": "北京",
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("期望 %v，实际 %v", expected, result)
		}
	})

	t.Run("忽略大小写功能", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(",", "=")
		flag.SetIgnoreCase(true)

		err := flag.Set("Name=张三,AGE=25,CiTy=北京")
		if err != nil {
			t.Fatalf("设置失败: %v", err)
		}

		result := flag.Get()
		// 所有键应该被转换为小写
		expected := map[string]string{
			"name": "张三",
			"age":  "25",
			"city": "北京",
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("期望 %v，实际 %v", expected, result)
		}
	})
}

// TestSliceFlag_AdvancedFeatures 测试SliceFlag的高级功能
func TestSliceFlag_AdvancedFeatures(t *testing.T) {
	t.Run("多分隔符支持", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
			delimiters: []string{",", ";", "|"},
		}

		// 测试不同分隔符
		testCases := []struct {
			input    string
			expected []string
		}{
			{"a,b,c", []string{"a", "b", "c"}},
			{"x;y;z", []string{"x", "y", "z"}},
			{"1|2|3", []string{"1", "2", "3"}},
		}

		for _, tc := range testCases {
			err := flag.Set(tc.input)
			if err != nil {
				t.Errorf("设置 '%s' 失败: %v", tc.input, err)
				continue
			}

			result := flag.Get()
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("输入 '%s'，期望 %v，实际 %v", tc.input, tc.expected, result)
			}
		}
	})

	t.Run("动态修改分隔符", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
		}

		// 设置初始分隔符
		flag.SetDelimiters([]string{","})
		err := flag.Set("a,b,c")
		if err != nil {
			t.Fatalf("设置失败: %v", err)
		}

		// 修改分隔符
		flag.SetDelimiters([]string{";"})
		err = flag.Set("x;y;z")
		if err != nil {
			t.Fatalf("修改分隔符后设置失败: %v", err)
		}

		result := flag.Get()
		expected := []string{"x", "y", "z"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("期望 %v，实际 %v", expected, result)
		}
	})

	t.Run("复杂元素处理", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
			delimiters: []string{","},
			skipEmpty:  false,
		}

		// 测试包含空格和特殊字符的元素
		err := flag.Set("  item1  , item2 with spaces ,item3,")
		if err != nil {
			t.Fatalf("设置失败: %v", err)
		}

		result := flag.Get()
		expected := []string{"item1", "item2 with spaces", "item3", ""}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("期望 %v，实际 %v", expected, result)
		}
	})

	t.Run("获取分隔符", func(t *testing.T) {
		flag := &SliceFlag{}
		customDelimiters := []string{",", ";", "|"}
		flag.SetDelimiters(customDelimiters)

		result := flag.GetDelimiters()
		if !reflect.DeepEqual(result, customDelimiters) {
			t.Errorf("期望分隔符 %v，实际 %v", customDelimiters, result)
		}

		// 修改返回的切片不应影响内部状态
		result[0] = "modified"
		internalDelimiters := flag.GetDelimiters()
		if internalDelimiters[0] == "modified" {
			t.Error("修改返回的分隔符切片不应影响内部状态")
		}
	})
}

// TestTimeFlag_AdvancedFeatures 测试TimeFlag的高级功能
func TestTimeFlag_AdvancedFeatures(t *testing.T) {
	t.Run("自定义输出格式", func(t *testing.T) {
		flag := &TimeFlag{
			BaseFlag: BaseFlag[time.Time]{
				initialValue: time.Time{},
				value:        new(time.Time),
			},
		}

		// 设置时间值
		testTime := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
		err := flag.Set(testTime.Format(time.RFC3339))
		if err != nil {
			t.Fatalf("设置时间失败: %v", err)
		}

		// 测试默认格式
		defaultFormat := flag.String()
		if defaultFormat != testTime.Format(time.RFC3339) {
			t.Errorf("默认格式期望 '%s'，实际 '%s'", testTime.Format(time.RFC3339), defaultFormat)
		}

		// 设置自定义输出格式
		flag.SetOutputFormat("2006-01-02 15:04:05")
		customFormat := flag.String()
		expected := "2024-03-15 14:30:45"
		if customFormat != expected {
			t.Errorf("自定义格式期望 '%s'，实际 '%s'", expected, customFormat)
		}
	})

	t.Run("多种时间格式解析", func(t *testing.T) {
		flag := &TimeFlag{
			BaseFlag: BaseFlag[time.Time]{
				initialValue: time.Time{},
				value:        new(time.Time),
			},
		}

		// 测试多种支持的时间格式
		testCases := []struct {
			input    string
			expected time.Time
		}{
			{"2024-03-15", time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)},
			{"2024-03-15 14:30", time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC)},
			{"2024-03-15 14:30:45", time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)},
			{"2024-03-15T14:30:45Z", time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)},
		}

		for _, tc := range testCases {
			err := flag.Set(tc.input)
			if err != nil {
				t.Errorf("解析时间格式 '%s' 失败: %v", tc.input, err)
				continue
			}

			result := flag.Get()
			if !result.Equal(tc.expected) {
				t.Errorf("时间格式 '%s'，期望 %v，实际 %v", tc.input, tc.expected, result)
			}
		}
	})

	t.Run("并发安全的输出格式设置", func(t *testing.T) {
		flag := &TimeFlag{
			BaseFlag: BaseFlag[time.Time]{
				initialValue: time.Time{},
				value:        new(time.Time),
			},
		}

		testTime := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
		flag.Set(testTime.Format(time.RFC3339))

		// 并发设置和读取输出格式
		done := make(chan bool, 2)

		go func() {
			for i := 0; i < 100; i++ {
				flag.SetOutputFormat("2006-01-02")
				_ = flag.String()
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				flag.SetOutputFormat("15:04:05")
				_ = flag.String()
			}
			done <- true
		}()

		<-done
		<-done
		// 如果没有panic，说明并发访问是安全的
	})
}

// TestEnumFlag_AdvancedFeatures 测试EnumFlag的高级功能
func TestEnumFlag_AdvancedFeatures(t *testing.T) {
	t.Run("动态修改大小写敏感性", func(t *testing.T) {
		flag := &EnumFlag{}
		options := []string{"Apple", "Banana", "Cherry"}

		err := flag.Init("fruit", "f", "Apple", "水果选择", options)
		if err != nil {
			t.Fatalf("初始化失败: %v", err)
		}

		// 默认不区分大小写
		err = flag.Set("apple")
		if err != nil {
			t.Errorf("默认不区分大小写时设置 'apple' 应该成功: %v", err)
		}

		// 设置为区分大小写
		flag.SetCaseSensitive(true)
		err = flag.Set("apple")
		if err == nil {
			t.Error("区分大小写时设置 'apple' 应该失败")
		}

		// 设置正确的大小写
		err = flag.Set("Apple")
		if err != nil {
			t.Errorf("区分大小写时设置 'Apple' 应该成功: %v", err)
		}
	})

	t.Run("获取选项列表", func(t *testing.T) {
		flag := &EnumFlag{}
		options := []string{"red", "green", "blue"}

		err := flag.Init("color", "c", "red", "颜色选择", options)
		if err != nil {
			t.Fatalf("初始化失败: %v", err)
		}

		result := flag.GetOptions()
		if !reflect.DeepEqual(result, options) {
			t.Errorf("期望选项 %v，实际 %v", options, result)
		}

		// 修改返回的切片不应影响内部状态
		result[0] = "modified"
		internalOptions := flag.GetOptions()
		if internalOptions[0] == "modified" {
			t.Error("修改返回的选项切片不应影响内部状态")
		}
	})

	t.Run("空选项列表处理", func(t *testing.T) {
		flag := &EnumFlag{}

		err := flag.Init("any", "a", "default", "任意值", []string{})
		if err != nil {
			t.Fatalf("空选项初始化失败: %v", err)
		}

		// 空选项时应该允许任意值
		err = flag.Set("任意值")
		if err != nil {
			t.Errorf("空选项时应该允许任意值: %v", err)
		}
	})

	t.Run("链式调用支持", func(t *testing.T) {
		flag := &EnumFlag{}
		options := []string{"option1", "option2"}

		err := flag.Init("test", "t", "option1", "测试", options)
		if err != nil {
			t.Fatalf("初始化失败: %v", err)
		}

		// 测试链式调用
		result := flag.SetCaseSensitive(true).SetCaseSensitive(false)
		if result != flag {
			t.Error("SetCaseSensitive应该返回自身以支持链式调用")
		}
	})
}

// TestDurationFlag_AdvancedFeatures 测试DurationFlag的高级功能
func TestDurationFlag_AdvancedFeatures(t *testing.T) {
	t.Run("复杂时间间隔解析", func(t *testing.T) {
		flag := &DurationFlag{
			BaseFlag: BaseFlag[time.Duration]{
				initialValue: 0,
				value:        new(time.Duration),
			},
		}

		testCases := []struct {
			input    string
			expected time.Duration
		}{
			{"1h30m", 90 * time.Minute},
			{"2h45m30s", 2*time.Hour + 45*time.Minute + 30*time.Second},
			{"500ms", 500 * time.Millisecond},
			{"1.5s", 1500 * time.Millisecond},
			{"24h", 24 * time.Hour},
		}

		for _, tc := range testCases {
			err := flag.Set(tc.input)
			if err != nil {
				t.Errorf("解析时间间隔 '%s' 失败: %v", tc.input, err)
				continue
			}

			result := flag.Get()
			if result != tc.expected {
				t.Errorf("时间间隔 '%s'，期望 %v，实际 %v", tc.input, tc.expected, result)
			}
		}
	})

	t.Run("大小写不敏感解析", func(t *testing.T) {
		flag := &DurationFlag{
			BaseFlag: BaseFlag[time.Duration]{
				initialValue: 0,
				value:        new(time.Duration),
			},
		}

		// 测试大写单位
		err := flag.Set("5S")
		if err != nil {
			t.Errorf("大写单位 '5S' 解析失败: %v", err)
		}

		if flag.Get() != 5*time.Second {
			t.Errorf("期望 5 秒，实际 %v", flag.Get())
		}

		// 测试混合大小写
		err = flag.Set("1H30M")
		if err != nil {
			t.Errorf("混合大小写 '1H30M' 解析失败: %v", err)
		}

		expected := time.Hour + 30*time.Minute
		if flag.Get() != expected {
			t.Errorf("期望 %v，实际 %v", expected, flag.Get())
		}
	})
}

// TestSliceFlag_EdgeCases 测试SliceFlag的边界情况
func TestSliceFlag_EdgeCases(t *testing.T) {
	t.Run("空分隔符处理", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
		}

		// 设置空分隔符应该使用默认值
		flag.SetDelimiters([]string{})
		delimiters := flag.GetDelimiters()

		if len(delimiters) == 0 {
			t.Error("空分隔符应该使用默认值")
		}
	})

	t.Run("单个元素无分隔符", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
			delimiters: []string{","},
		}

		err := flag.Set("single_element")
		if err != nil {
			t.Fatalf("设置单个元素失败: %v", err)
		}

		result := flag.Get()
		expected := []string{"single_element"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("期望 %v，实际 %v", expected, result)
		}
	})

	t.Run("移除不存在的元素", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
			delimiters: []string{","},
		}

		flag.Set("a,b,c")

		// 移除不存在的元素
		err := flag.Remove("d")
		if err != nil {
			t.Errorf("移除不存在的元素不应该返回错误: %v", err)
		}

		// 验证原有元素未受影响
		result := flag.Get()
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("移除不存在元素后，期望 %v，实际 %v", expected, result)
		}
	})

	t.Run("排序空切片", func(t *testing.T) {
		flag := &SliceFlag{
			BaseFlag: BaseFlag[[]string]{
				initialValue: []string{},
				value:        new([]string),
			},
		}

		flag.Clear()
		err := flag.Sort()
		if err != nil {
			t.Errorf("排序空切片不应该返回错误: %v", err)
		}

		if flag.Len() != 0 {
			t.Errorf("排序后空切片长度应为0，实际为 %d", flag.Len())
		}
	})
}

// TestMapFlag_EdgeCases 测试MapFlag的边界情况
func TestMapFlag_EdgeCases(t *testing.T) {
	t.Run("键值包含分隔符", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(";", "=") // 使用分号作为键分隔符，避免冲突

		// 值中包含逗号（不是键分隔符）
		err := flag.Set("key1=value,with,commas;key2=another,value")
		if err != nil {
			t.Fatalf("设置包含分隔符的值失败: %v", err)
		}

		result := flag.Get()
		if result["key1"] != "value,with,commas" {
			t.Errorf("期望值 'value,with,commas'，实际 '%s'", result["key1"])
		}
		if result["key2"] != "another,value" {
			t.Errorf("期望值 'another,value'，实际 '%s'", result["key2"])
		}
	})

	t.Run("空键值处理", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters(",", "=")

		// 测试空键
		err := flag.Set("=value")
		if err == nil {
			t.Error("空键应该返回错误")
		}

		// 测试空值
		err = flag.Set("key=")
		if err == nil {
			t.Error("空值应该返回错误")
		}
	})

	t.Run("默认分隔符处理", func(t *testing.T) {
		flag := &MapFlag{}
		flag.SetDelimiters("", "")

		// 应该使用默认分隔符
		err := flag.Set("key=value")
		if err != nil {
			t.Errorf("使用默认分隔符应该成功: %v", err)
		}
	})
}
