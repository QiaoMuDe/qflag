package qflag

import (
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
