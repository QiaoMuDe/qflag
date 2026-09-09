package qflag

import (
	"errors"
	"math"
	"os"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/internal/cmd"
	"gitee.com/MM-Q/qflag/internal/types"
)

func TestParser_IntFlag(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	intFlag := c.Int("count", "c", "计数器", 0)

	err := c.Parse([]string{"--count", "42"})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	if intFlag.Get() != 42 {
		t.Errorf("Count flag: expected 42, got %v", intFlag.Get())
	}
}

func TestParser_StringFlag(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	strFlag := c.String("name", "n", "名称", "")

	err := c.Parse([]string{"--name", "test"})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	if strFlag.Get() != "test" {
		t.Errorf("Name flag: expected 'test', got '%v'", strFlag.Get())
	}
}

func TestParser_EnvVarLoading(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	nameFlag := c.String("name", "n", "名称", "default")
	nameFlag.BindEnv("TEST_NAME")
	countFlag := c.Int("count", "c", "计数器", 0)
	countFlag.BindEnv("TEST_COUNT")

	if err := os.Setenv("TEST_NAME", "env-name"); err != nil {
		t.Fatalf("Failed to set TEST_NAME: %v", err)
	}
	if err := os.Setenv("TEST_COUNT", "100"); err != nil {
		t.Fatalf("Failed to set TEST_COUNT: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_NAME"); err != nil {
			t.Logf("Failed to unset TEST_NAME: %v", err)
		}
		if err := os.Unsetenv("TEST_COUNT"); err != nil {
			t.Logf("Failed to unset TEST_COUNT: %v", err)
		}
	}()

	err := c.Parse([]string{})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	if nameFlag.Get() != "env-name" {
		t.Errorf("Name flag: expected 'env-name', got '%v'", nameFlag.Get())
	}

	if countFlag.Get() != 100 {
		t.Errorf("Count flag: expected 100, got %v", countFlag.Get())
	}
}

func TestParser_DurationFlag(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	durationFlag := c.Duration("timeout", "t", "超时时间", time.Second*30)

	err := c.Parse([]string{"--timeout", "1m30s"})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	expected := time.Minute*1 + time.Second*30
	if durationFlag.Get() != expected {
		t.Errorf("Timeout flag: expected %v, got %v", expected, durationFlag.Get())
	}
}

func TestParser_TimeFlag(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	timeFlag := c.Time("start", "s", "开始时间", time.Time{})

	now := time.Now()
	formatted := now.Format(time.RFC3339)
	err := c.Parse([]string{"--start", formatted})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	parsed, _ := time.Parse(time.RFC3339, formatted)
	if !timeFlag.Get().Equal(parsed) {
		t.Errorf("Start flag: expected %v, got %v", parsed, timeFlag.Get())
	}
}

func TestParser_SliceFlags(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	stringSliceFlag := c.StringSlice("paths", "p", "路径列表", nil)
	intSliceFlag := c.IntSlice("ports", "P", "端口列表", nil)

	err := c.Parse([]string{"--paths", "a,b,c", "--ports", "80,443,8080"})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	paths := stringSliceFlag.Get()
	if len(paths) != 3 {
		t.Errorf("Paths: expected 3 items, got %d", len(paths))
	}

	ports := intSliceFlag.Get()
	if len(ports) != 3 {
		t.Errorf("Ports: expected 3 items, got %d", len(ports))
	}
}

func TestParser_SizeFlag(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	sizeFlag := c.Size("max-size", "s", "最大大小", 0)

	err := c.Parse([]string{"--max-size", "1MB"})
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	expected := int64(1000000)
	if sizeFlag.Get() != expected {
		t.Errorf("Max-size flag: expected %d, got %d", expected, sizeFlag.Get())
	}
}

func TestParser_CommandPath(t *testing.T) {
	rootCmd := cmd.NewCmd("app", "a", types.ContinueOnError)
	subCmd := cmd.NewCmd("sub", "s", types.ContinueOnError)
	subSubCmd := cmd.NewCmd("deep", "d", types.ContinueOnError)

	if err := rootCmd.AddSubCmds(subCmd); err != nil {
		t.Fatalf("Failed to add subcommand: %v", err)
	}
	if err := subCmd.AddSubCmds(subSubCmd); err != nil {
		t.Fatalf("Failed to add subcommand: %v", err)
	}

	if subSubCmd.Path() != "app sub deep" {
		t.Errorf("Expected path 'app sub deep', got '%s'", subSubCmd.Path())
	}
}

func TestParser_HelpGeneration(t *testing.T) {
	c := createTestCmd()

	if err := c.Parse([]string{}); err != nil {
		t.Errorf("Parse error: %v", err)
	}

	c.PrintHelp()
}

// TestParser_ResetOnRepeatedParse 测试重复解析时标志重置
// 这个测试验证了在重复解析场景下，标志值和状态会被正确重置
func TestParser_ResetOnRepeatedParse(t *testing.T) {
	// 场景1：测试基本类型的重复解析
	t.Run("BasicTypeReset", func(t *testing.T) {
		c := cmd.NewCmd("test", "t", types.ContinueOnError)
		nameFlag := c.String("name", "n", "名称", "default")
		countFlag := c.Int("count", "c", "计数器", 0)

		// 第一次解析：设置值
		err := c.Parse([]string{"--name", "first", "--count", "42"})
		if err != nil {
			t.Fatalf("First parse error: %v", err)
		}

		if nameFlag.Get() != "first" {
			t.Errorf("First parse: expected name 'first', got '%s'", nameFlag.Get())
		}
		if countFlag.Get() != 42 {
			t.Errorf("First parse: expected count 42, got %d", countFlag.Get())
		}
		if !nameFlag.IsSet() {
			t.Error("First parse: name flag should be set")
		}
		if !countFlag.IsSet() {
			t.Error("First parse: count flag should be set")
		}

		// 第二次解析：不设置值，应该重置到默认值
		err = c.Parse([]string{})
		if err != nil {
			t.Fatalf("Second parse error: %v", err)
		}

		if nameFlag.Get() != "default" {
			t.Errorf("Second parse: expected name 'default', got '%s'", nameFlag.Get())
		}
		if countFlag.Get() != 0 {
			t.Errorf("Second parse: expected count 0, got %d", countFlag.Get())
		}
		if nameFlag.IsSet() {
			t.Error("Second parse: name flag should not be set")
		}
		if countFlag.IsSet() {
			t.Error("Second parse: count flag should not be set")
		}
	})

	// 场景2：测试集合类型的重复解析
	t.Run("SliceTypeReset", func(t *testing.T) {
		c := cmd.NewCmd("test", "t", types.ContinueOnError)
		tagsFlag := c.StringSlice("tags", "t", "标签列表", []string{"default-tag"})

		// 第一次解析：设置值
		err := c.Parse([]string{"--tags", "a,b,c"})
		if err != nil {
			t.Fatalf("First parse error: %v", err)
		}

		tags := tagsFlag.Get()
		if len(tags) != 3 {
			t.Errorf("First parse: expected 3 tags, got %d", len(tags))
		}
		if !tagsFlag.IsSet() {
			t.Error("First parse: tags flag should be set")
		}

		// 第二次解析：不设置值，应该重置到默认值
		err = c.Parse([]string{})
		if err != nil {
			t.Fatalf("Second parse error: %v", err)
		}

		tags = tagsFlag.Get()
		if len(tags) != 1 || tags[0] != "default-tag" {
			t.Errorf("Second parse: expected default tags, got %v", tags)
		}
		if tagsFlag.IsSet() {
			t.Error("Second parse: tags flag should not be set")
		}
	})

	// 场景3：测试环境变量在重复解析中的加载
	t.Run("EnvVarLoadingAfterReset", func(t *testing.T) {
		// 设置环境变量
		if err := os.Setenv("TEST_RESET_NAME", "env-value"); err != nil {
			t.Fatalf("Failed to set env var: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("TEST_RESET_NAME"); err != nil {
				t.Logf("Failed to unset env var: %v", err)
			}
		}()

		c := cmd.NewCmd("test", "t", types.ContinueOnError)
		nameFlag := c.String("name", "n", "名称", "default")
		nameFlag.BindEnv("TEST_RESET_NAME")

		// 第一次解析：命令行参数覆盖环境变量
		err := c.Parse([]string{"--name", "cli-value"})
		if err != nil {
			t.Fatalf("First parse error: %v", err)
		}

		if nameFlag.Get() != "cli-value" {
			t.Errorf("First parse: expected 'cli-value', got '%s'", nameFlag.Get())
		}
		if !nameFlag.IsSet() {
			t.Error("First parse: name flag should be set")
		}

		// 第二次解析：没有命令行参数，应该加载环境变量
		err = c.Parse([]string{})
		if err != nil {
			t.Fatalf("Second parse error: %v", err)
		}

		if nameFlag.Get() != "env-value" {
			t.Errorf("Second parse: expected 'env-value', got '%s'", nameFlag.Get())
		}
		if !nameFlag.IsSet() {
			t.Error("Second parse: name flag should be set by env var")
		}
	})
}

// TestParser_NegativeValues 测试取值标志的负值（以 - 开头）解析
//
// 验证修复后的预扫描不会把值误判为未知标志，且最终能正确解析到负值。
// 涵盖常见的 int、int64、float64、string 类型的空格分隔负值场景。
func TestParser_NegativeValues(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	intFlag := c.Int("count", "c", "计数器", 0)
	int64Flag := c.Int64("big", "b", "大数", 0)
	floatFlag := c.Float64("ratio", "r", "比例", 0.0)
	stringFlag := c.String("name", "n", "名称", "")

	err := c.Parse([]string{"--count", "-8080", "--big", "-9223372036854775808", "--ratio", "-1.5", "--name", "-foo"})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if intFlag.Get() != -8080 {
		t.Errorf("int flag: expected -8080, got %d", intFlag.Get())
	}
	if int64Flag.Get() != -9223372036854775808 {
		t.Errorf("int64 flag: expected -9223372036854775808, got %d", int64Flag.Get())
	}
	if floatFlag.Get() != -1.5 {
		t.Errorf("float64 flag: expected -1.5, got %f", floatFlag.Get())
	}
	if stringFlag.Get() != "-foo" {
		t.Errorf("string flag: expected '-foo', got '%s'", stringFlag.Get())
	}
}

// TestParser_InlineNegative 测试内联形式（--flag=-value）的负值解析
//
// 内联形式原本就正常，本用例用于回归保护，确保修复后行为不变。
func TestParser_InlineNegative(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	intFlag := c.Int("count", "c", "计数器", 0)
	floatFlag := c.Float64("ratio", "r", "比例", 0.0)
	stringFlag := c.String("name", "n", "名称", "")

	err := c.Parse([]string{"--count=-42", "--ratio=-3.14", "--name=-hello-world"})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if intFlag.Get() != -42 {
		t.Errorf("int flag: expected -42, got %d", intFlag.Get())
	}
	if floatFlag.Get() != -3.14 {
		t.Errorf("float64 flag: expected -3.14, got %f", floatFlag.Get())
	}
	if stringFlag.Get() != "-hello-world" {
		t.Errorf("string flag: expected '-hello-world', got '%s'", stringFlag.Get())
	}
}

// TestParser_BoolFollowedByUnknown 验证布尔标志后仍报错未知标志（行为不变）
//
// 布尔标志不消费下一个参数，因此后跟未知标志必须仍然报错。
// 该用例确保本次修复没有破坏原有行为。
func TestParser_BoolFollowedByUnknown(t *testing.T) {
	c := cmd.NewCmd("test", "t", types.ContinueOnError)
	c.Bool("verbose", "v", "详细输出", false)

	err := c.Parse([]string{"--verbose", "-x"})
	if err == nil {
		t.Fatalf("expected unknown flag error, got nil")
	}

	var ufe *types.UnknownFlagError
	if !errors.As(err, &ufe) {
		t.Fatalf("expected *types.UnknownFlagError, got %T: %v", err, err)
	}
	if ufe.Input != "-x" {
		t.Errorf("expected error input '-x', got '%s'", ufe.Input)
	}
}

// TestParser_BoundaryValues 验证各种边界场景的解析
//
// 涵盖 float NaN/Inf、duration 负值、int slice 含负数成员等边界情况。
func TestParser_BoundaryValues(t *testing.T) {
	t.Run("float special values", func(t *testing.T) {
		c := cmd.NewCmd("test", "t", types.ContinueOnError)
		nanFlag := c.Float64("nan", "", "NaN", 0)
		posInfFlag := c.Float64("pinf", "", "+Inf", 0)
		negInfFlag := c.Float64("ninf", "", "-Inf", 0)

		err := c.Parse([]string{"--nan", "NaN", "--pinf", "+Inf", "--ninf", "-Inf"})
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		if !math.IsNaN(nanFlag.Get()) {
			t.Errorf("nan flag: expected NaN, got %f", nanFlag.Get())
		}
		if !math.IsInf(posInfFlag.Get(), 1) {
			t.Errorf("+Inf flag: expected +Inf, got %f", posInfFlag.Get())
		}
		if !math.IsInf(negInfFlag.Get(), -1) {
			t.Errorf("-Inf flag: expected -Inf, got %f", negInfFlag.Get())
		}
	})

	t.Run("negative duration", func(t *testing.T) {
		c := cmd.NewCmd("test", "t", types.ContinueOnError)
		durFlag := c.Duration("offset", "o", "偏移", 0)

		err := c.Parse([]string{"--offset", "-1m30s"})
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		expected := -(time.Minute*1 + time.Second*30)
		if durFlag.Get() != expected {
			t.Errorf("duration: expected %v, got %v", expected, durFlag.Get())
		}
	})

	t.Run("int slice with negative values", func(t *testing.T) {
		c := cmd.NewCmd("test", "t", types.ContinueOnError)
		sliceFlag := c.IntSlice("nums", "n", "数值列表", nil)

		err := c.Parse([]string{"--nums", "-10,5,-3,0"})
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}

		expected := []int{-10, 5, -3, 0}
		got := sliceFlag.Get()
		if len(got) != len(expected) {
			t.Fatalf("expected %d items, got %d", len(expected), len(got))
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("index %d: expected %d, got %d", i, expected[i], got[i])
			}
		}
	})
}
