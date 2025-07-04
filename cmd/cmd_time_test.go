package cmd

import (
	"flag"
	"testing"
	"time"

	"gitee.com/MM-Q/qflag/flags"
)

func TestTime(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	timeFlag := cmd.Time("time", "t", defaultTime, "time flag test")

	// 测试默认值
	if !timeFlag.Get().Equal(defaultTime) {
		t.Errorf("Expected default time %v, got %v", defaultTime, timeFlag.Get())
	}

	// 解析参数
	args := []string{"--time", "2023-12-31T23:59:59Z"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	if !timeFlag.Get().Equal(parsedTime) {
		t.Errorf("Expected parsed time %v, got %v", parsedTime, timeFlag.Get())
	}
}

func TestTimeVar(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var timeFlag flags.TimeFlag
	defaultTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	cmd.TimeVar(&timeFlag, "time", "t", defaultTime, "time flag test")

	if !timeFlag.Get().Equal(defaultTime) {
		t.Errorf("Expected default time %v, got %v", defaultTime, timeFlag.Get())
	}

	args := []string{"-t", "2024-01-01T12:00:00+08:00"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, "2024-01-01T12:00:00+08:00")
	if err != nil {
		t.Fatalf("Failed to parse test time: %v", err)
	}

	if !timeFlag.Get().Equal(parsedTime) {
		t.Errorf("Expected parsed time %v, got %v", parsedTime, timeFlag.Get())
	}
}

func TestTimeVar_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to TimeVar")
		}
	}()
	cmd.TimeVar(nil, "time", "t", time.Time{}, "test")
}
