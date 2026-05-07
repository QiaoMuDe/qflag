package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

// TestDependencyChain 测试依赖链传递问题
func TestDependencyChain(t *testing.T) {
	// 场景: ssl -> cert -> key
	// dep1: ssl 需要 cert
	// dep2: cert 需要 key
	// 问题: 设置 ssl 但不设置 cert 和 key 时，应该报错
	//       但设置 ssl 和 cert 但不设置 key 时，是否会检查 key?

	t.Log("=== 测试依赖链传递 ===")
	t.Log("配置:")
	t.Log("  dep1: trigger=ssl, targets=[cert], type=DepRequired")
	t.Log("  dep2: trigger=cert, targets=[key], type=DepRequired")
	t.Log()

	// 测试1: 只设置 ssl
	t.Log("测试1: 只设置 ssl (--ssl)")
	cmd1 := NewCmd("test1", "t1", types.ContinueOnError)
	if err := cmd1.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd1.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd1.AddFlag(flag.NewStringFlag("key", "k", "Key", "")); err != nil {
		t.Fatalf("Failed to add flag key: %v", err)
	}
	if err := cmd1.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd1.AddFlagDependency("cert_requires_key", "cert", []string{"key"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err := cmd1.Parse([]string{"--ssl"})
	if err != nil {
		t.Logf("  结果: 报错 (符合预期) - %v", err)
	} else {
		t.Log("  结果: 没有报错 (问题!)")
	}

	// 测试2: 设置 ssl 和 cert，但不设置 key
	t.Log("\n测试2: 设置 ssl 和 cert，但不设置 key (--ssl --cert cert.pem)")
	cmd2 := NewCmd("test2", "t2", types.ContinueOnError)
	if err := cmd2.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd2.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd2.AddFlag(flag.NewStringFlag("key", "k", "Key", "")); err != nil {
		t.Fatalf("Failed to add flag key: %v", err)
	}
	if err := cmd2.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd2.AddFlagDependency("cert_requires_key", "cert", []string{"key"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err = cmd2.Parse([]string{"--ssl", "--cert", "cert.pem"})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
		t.Log("  分析: 依赖链传递生效，检测到 cert 需要 key")
	} else {
		t.Log("  结果: 通过验证 (依赖链未传递)")
		t.Log("  分析: 只验证了 dep1 (ssl->cert)，没有验证 dep2 (cert->key)")
	}

	// 测试3: 设置 ssl, cert 和 key
	t.Log("\n测试3: 设置 ssl, cert 和 key (--ssl --cert cert.pem --key key.pem)")
	cmd3 := NewCmd("test3", "t3", types.ContinueOnError)
	if err := cmd3.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd3.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd3.AddFlag(flag.NewStringFlag("key", "k", "Key", "")); err != nil {
		t.Fatalf("Failed to add flag key: %v", err)
	}
	if err := cmd3.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd3.AddFlagDependency("cert_requires_key", "cert", []string{"key"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err = cmd3.Parse([]string{"--ssl", "--cert", "cert.pem", "--key", "key.pem"})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
	} else {
		t.Log("  结果: 通过验证 (符合预期)")
	}

	// 测试4: 互斥依赖链
	t.Log("\n=== 测试互斥依赖链 ===")
	t.Log("配置:")
	t.Log("  dep1: trigger=debug, targets=[ssl], type=DepMutex")
	t.Log("  dep2: trigger=ssl, targets=[cert], type=DepRequired")
	t.Log()

	// 设置 debug 和 ssl，但不设置 cert
	t.Log("测试4: 设置 debug 和 ssl，但不设置 cert (--debug --ssl)")
	cmd4 := NewCmd("test4", "t4", types.ContinueOnError)
	if err := cmd4.AddFlag(flag.NewBoolFlag("debug", "d", "Debug mode", false)); err != nil {
		t.Fatalf("Failed to add flag debug: %v", err)
	}
	if err := cmd4.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
		t.Fatalf("Failed to add flag ssl: %v", err)
	}
	if err := cmd4.AddFlag(flag.NewStringFlag("cert", "c", "Certificate", "")); err != nil {
		t.Fatalf("Failed to add flag cert: %v", err)
	}
	if err := cmd4.AddFlagDependency("debug_mutex_ssl", "debug", []string{"ssl"}, types.DepMutex); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}
	if err := cmd4.AddFlagDependency("ssl_requires_cert", "ssl", []string{"cert"}, types.DepRequired); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	err = cmd4.Parse([]string{"--debug", "--ssl"})
	if err != nil {
		t.Logf("  结果: 报错 - %v", err)
		// 检查是哪个依赖触发的错误
		if containsStr(err.Error(), "cannot be used with") {
			t.Log("  分析: dep1 (debug->ssl 互斥) 触发")
		} else if containsStr(err.Error(), "requires flags") {
			t.Log("  分析: dep2 (ssl->cert 必需) 触发")
		}
	} else {
		t.Log("  结果: 通过验证")
	}

	t.Log("\n结论:")
	t.Log("依赖链传递问题是否存在，取决于测试2的结果")
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStrHelper(s, substr))
}

func containsStrHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
