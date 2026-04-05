package flag

import (
	"os"
	"strings"
	"testing"
)

func TestAutoBindEnv(t *testing.T) {
	tests := []struct {
		name      string
		longName  string
		shortName string
		wantEnv   string
		wantPanic bool
	}{
		{
			name:      "长名称转为大写",
			longName:  "host",
			shortName: "h",
			wantEnv:   "HOST",
			wantPanic: false,
		},
		{
			name:      "复杂标志名-连字符",
			longName:  "db-host",
			shortName: "",
			wantEnv:   "DB-HOST",
			wantPanic: false,
		},
		{
			name:      "复杂标志名-下划线",
			longName:  "db_port",
			shortName: "",
			wantEnv:   "DB_PORT",
			wantPanic: false,
		},
		{
			name:      "复杂标志名-驼峰",
			longName:  "dbUserName",
			shortName: "",
			wantEnv:   "DBUSERNAME",
			wantPanic: false,
		},
		{
			name:      "没有长名称应panic",
			longName:  "",
			shortName: "h",
			wantEnv:   "",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := NewStringFlag(tt.longName, tt.shortName, "测试标志", "default")

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("期望 panic 但没有发生")
					}
				}()
			}

			flag.AutoBindEnv()

			if !tt.wantPanic {
				if got := flag.GetEnvVar(); got != tt.wantEnv {
					t.Errorf("AutoBindEnv() 环境变量 = %v, 期望 %v", got, tt.wantEnv)
				}
			}
		})
	}
}

func TestAutoBindEnvIntegration(t *testing.T) {
	// 设置环境变量
	if err := os.Setenv("HOST", "192.168.1.1"); err != nil {
		t.Fatalf("设置环境变量 HOST 失败: %v", err)
	}
	if err := os.Setenv("DB-HOST", "192.168.1.100"); err != nil {
		t.Fatalf("设置环境变量 DB-HOST 失败: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("HOST"); err != nil {
			t.Logf("清理环境变量 HOST 失败: %v", err)
		}
		if err := os.Unsetenv("DB-HOST"); err != nil {
			t.Logf("清理环境变量 DB-HOST 失败: %v", err)
		}
	}()

	tests := []struct {
		name      string
		longName  string
		shortName string
		envValue  string
		wantValue string
	}{
		{
			name:      "基础环境变量绑定",
			longName:  "host",
			shortName: "h",
			envValue:  "HOST",
			wantValue: "192.168.1.1",
		},
		{
			name:      "复杂环境变量绑定",
			longName:  "db-host",
			shortName: "",
			envValue:  "DB-HOST",
			wantValue: "192.168.1.100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := NewStringFlag(tt.longName, tt.shortName, "测试标志", "default")
			flag.AutoBindEnv()

			// 验证环境变量绑定
			if got := flag.GetEnvVar(); got != strings.ToUpper(tt.longName) {
				t.Errorf("环境变量绑定错误, 期望 %s, 得到 %s", strings.ToUpper(tt.longName), got)
			}
		})
	}
}

func TestBindEnvVsAutoBindEnv(t *testing.T) {
	// 测试 BindEnv 和 AutoBindEnv 的区别
	flag1 := NewStringFlag("host", "h", "主机地址", "localhost")
	flag1.BindEnv("DATABASE_HOST")
	if got := flag1.GetEnvVar(); got != "DATABASE_HOST" {
		t.Errorf("BindEnv 错误, 期望 DATABASE_HOST, 得到 %s", got)
	}

	flag2 := NewStringFlag("host", "h", "主机地址", "localhost")
	flag2.AutoBindEnv()
	if got := flag2.GetEnvVar(); got != "HOST" {
		t.Errorf("AutoBindEnv 错误, 期望 HOST, 得到 %s", got)
	}
}
