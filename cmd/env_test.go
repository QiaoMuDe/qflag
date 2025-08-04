package cmd

import (
	"flag"
	"os"
	"testing"
	"time"
)

// TestLoadEnvVars 测试环境变量加载功能的各种场景
func TestLoadEnvVars(t *testing.T) {
	// 测试环境变量覆盖默认值
	t.Run("environment variable overrides default value", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		strVal := cmd.String("str-flag", "sf", "default", "测试字符串标志").BindEnv("TEST_STR_ENV")

		// 设置环境变量
		if err := os.Setenv("TEST_STR_ENV", "env_value"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			_ = os.Unsetenv("TEST_STR_ENV")
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		if strVal.Get() != "env_value" {
			t.Errorf("期望获取环境变量值 'env_value', 实际获取 '%s'", strVal.Get())
		}
	})

	// 测试枚举类型环境变量
	t.Run("enum type environment variable", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		enumVal := cmd.Enum("mode", "m", "test", "测试枚举标志", []string{"debug", "test", "prod"}).BindEnv("TEST_ENUM_ENV")

		// 设置有效环境变量
		if err := os.Setenv("TEST_ENUM_ENV", "prod"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			_ = os.Unsetenv("TEST_ENUM_ENV")
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		if enumVal.Get() != "prod" {
			t.Errorf("期望枚举值 'prod', 实际获取 '%s'", enumVal.Get())
		}

		// 测试无效枚举值
		if setEnvErr := os.Setenv("TEST_ENUM_ENV", "invalid"); setEnvErr != nil {
			t.Fatalf("设置环境变量失败: %v", setEnvErr)
		}

		err = cmd.loadEnvVars()
		if err == nil {
			t.Error("期望解析无效枚举值时返回错误, 但未返回错误")
		}
	})

	// 测试命令行参数优先级高于环境变量
	t.Run("command line argument has higher priority than environment variable", func(t *testing.T) {
		cmd := NewCmd("test", "t", flag.ContinueOnError)
		intVal := cmd.Int("int-flag", "", 10, "测试整数标志").BindEnv("TEST_INT_ENV")

		// 设置环境变量
		if err := os.Setenv("TEST_INT_ENV", "20"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			_ = os.Unsetenv("TEST_INT_ENV")
		}()

		// 设置命令行参数
		if err := cmd.Parse([]string{"--int-flag", "30"}); err != nil {
			t.Error("解析错误")
		}

		// 验证结果
		if intVal.Get() != 30 {
			t.Errorf("期望命令行参数值 30, 实际获取 %d", intVal.Get())
		}
	})

	// 测试环境变量未设置时使用默认值
	t.Run("use default value when environment variable not set", func(t *testing.T) {
		cmd := NewCmd("start", "s", flag.ContinueOnError)
		boolVal := cmd.Bool("bool-flag", "", false, "测试布尔标志").BindEnv("TEST_BOOL_ENV")

		// 确保环境变量未设置
		err := os.Unsetenv("TEST_BOOL_ENV")
		if err != nil {
			t.Fatalf("Failed to unset TEST_BOOL_ENV: %v", err)
		}

		// 加载环境变量
		err = cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		if boolVal.Get() != false {
			t.Errorf("期望默认值 false, 实际获取 %t", boolVal.Get())
		}
	})

	// 测试环境变量类型转换
	t.Run("environment variable type conversion", func(t *testing.T) {
		cmd := NewCmd("test", "", flag.ContinueOnError)
		durationVal := cmd.Duration("duration-flag", "", time.Second*10, "测试时长标志").BindEnv("TEST_DURATION_ENV")

		// 设置环境变量
		if err := os.Setenv("TEST_DURATION_ENV", "30s"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_DURATION_ENV"); err != nil {
				t.Errorf("Failed to unset TEST_DURATION_ENV: %v", err)
			}
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		if durationVal.Get() != time.Second*30 {
			t.Errorf("期望时长 30s, 实际获取 %v", durationVal.Get())
		}
	})

	// 测试映射类型环境变量
	t.Run("map type environment variable", func(t *testing.T) {
		cmd := NewCmd("test", "", flag.ContinueOnError)
		mapVal := cmd.Map("map-flag", "", map[string]string{"default": "value"}, "测试映射标志").BindEnv("TEST_MAP_ENV")

		// 设置环境变量
		if err := os.Setenv("TEST_MAP_ENV", "key1=val1,key2=val2"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_MAP_ENV"); err != nil {
				t.Errorf("Failed to unset TEST_MAP_ENV: %v", err)
			}
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		result := mapVal.Get()
		if result["key1"] != "val1" || result["key2"] != "val2" {
			t.Errorf("期望映射值 {key1:val1, key2:val2}, 实际获取 %v", result)
		}
	})

	// 测试切片类型环境变量
	t.Run("slice type environment variable", func(t *testing.T) {
		cmd := NewCmd("test", "", flag.ContinueOnError)
		sliceVal := cmd.Slice("slice-flag", "", []string{"default"}, "测试切片标志").BindEnv("TEST_SLICE_ENV")

		// 设置环境变量
		if err := os.Setenv("TEST_SLICE_ENV", "item1,item2,item3"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_SLICE_ENV"); err != nil {
				t.Errorf("Failed to unset TEST_SLICE_ENV: %v", err)
			}
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		result := sliceVal.Get()
		if len(result) != 3 || result[0] != "item1" || result[1] != "item2" || result[2] != "item3" {
			t.Errorf("期望切片值 [item1,item2,item3], 实际获取 %v", result)
		}
	})

	// 测试无符号整数类型环境变量
	t.Run("uint64 type environment variable", func(t *testing.T) {
		cmd := NewCmd("test", "", flag.ContinueOnError)
		uint64Val := cmd.Uint64("uint64-flag", "", 100, "测试Uint64标志").BindEnv("TEST_UINT64_ENV")

		// 设置环境变量
		if err := os.Setenv("TEST_UINT64_ENV", "18446744073709551615"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_UINT64_ENV"); err != nil {
				t.Errorf("Failed to unset TEST_UINT64_ENV: %v", err)
			}
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		if uint64Val.Get() != 18446744073709551615 {
			t.Errorf("期望Uint64值 18446744073709551615, 实际获取 %d", uint64Val.Get())
		}
	})

	// 测试Time类型环境变量
	t.Run("time type environment variable", func(t *testing.T) {
		cmd := NewCmd("test", "", flag.ContinueOnError)
		defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		timeVal := cmd.Time("time-flag", "", defaultTime, "测试Time标志").BindEnv("TEST_TIME_ENV")

		// 设置环境变量 (RFC3339格式)
		envTime := "2024-05-20T15:30:45Z"
		if err := os.Setenv("TEST_TIME_ENV", envTime); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_TIME_ENV"); err != nil {
				t.Errorf("Failed to unset TEST_TIME_ENV: %v", err)
			}
		}()

		// 加载环境变量
		err := cmd.loadEnvVars()
		if err != nil {
			t.Fatalf("加载环境变量失败: %v", err)
		}

		parsedTime := timeVal.Get()
		expectedTime, _ := time.Parse(time.RFC3339, envTime)
		if !parsedTime.Equal(expectedTime) {
			t.Errorf("期望Time值 '%s', 实际获取 '%s'", expectedTime.Format(time.RFC3339), parsedTime.Format(time.RFC3339))
		}
	})

	// 测试环境变量解析错误
	t.Run("environment variable parsing error", func(t *testing.T) {
		cmd := NewCmd("test", "", flag.ContinueOnError)
		_ = cmd.Float64("float-flag", "", 3.14, "测试浮点数标志").BindEnv("TEST_FLOAT_ENV")

		// 设置无效的环境变量值
		if err := os.Setenv("TEST_FLOAT_ENV", "not_a_float"); err != nil {
			t.Fatalf("设置环境变量失败: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_FLOAT_ENV"); err != nil {
				t.Errorf("Failed to unset TEST_FLOAT_ENV: %v", err)
			}
		}()

		// 加载环境变量并验证错误
		err := cmd.loadEnvVars()
		if err == nil {
			t.Error("期望解析环境变量时返回错误, 但未返回错误")
		} else if err.Error() != "validation failed: Failed to load environment variables: validation failed: Failed to parse environment variable TEST_FLOAT_ENV for flag float-flag: validation failed: failed to parse float64 value: strconv.ParseFloat: parsing \"not_a_float\": invalid syntax" {
			t.Errorf("期望特定错误信息, 实际获取: %v", err)
		}
	})
}
