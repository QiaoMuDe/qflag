package cmd

import (
	"flag"
	"net/url"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

func TestURL(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defaultURL := "https://example.com"
	urlFlag := cmd.URL("url", "u", defaultURL, "url flag test")

	// 测试默认值
	if urlFlag.Get() != defaultURL {
		t.Errorf("Expected default URL '%s', got '%s'", defaultURL, urlFlag.Get())
	}

	// 解析参数
	args := []string{"--url", "https://test.com/path?query=1"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expectedURL, _ := url.Parse("https://test.com/path?query=1")
	if urlFlag.Get() != expectedURL.String() {
		t.Errorf("Expected parsed URL '%s', got '%s'", expectedURL.String(), urlFlag.Get())
	}
}

func TestURLVar(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	var urlFlag flags.URLFlag
	defaultURL := "https://example.com"
	cmd.URLVar(&urlFlag, "url", "u", defaultURL, "url flag test")

	if urlFlag.Get() != defaultURL {
		t.Errorf("Expected default URL '%s', got '%s'", defaultURL, urlFlag.Get())
	}

	args := []string{"-u", "https://short.com"}
	if err := cmd.Parse(args); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if urlFlag.Get() != "https://short.com" {
		t.Errorf("Expected parsed URL 'https://short.com', got '%s'", urlFlag.Get())
	}
}

func TestURLVar_NilPointer(t *testing.T) {
	cmd := NewCmd("test", "t", flag.ContinueOnError)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to URLVar")
		}
	}()
	cmd.URLVar(nil, "url", "u", "https://example.com", "test")
}
