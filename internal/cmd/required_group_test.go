package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

func TestAddRequiredGroup(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host address", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port number", "")); err != nil {
		t.Fatalf("Failed to add port flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("username", "U", "Username", "")); err != nil {
		t.Fatalf("Failed to add username flag: %v", err)
	}

	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
		t.Fatalf("Failed to add required group: %v", err)
	}

	groups := cmd.Config().RequiredGroups
	if len(groups) != 1 {
		t.Fatalf("Expected 1 required group, got %d", len(groups))
	}

	group := groups[0]
	if group.Name != "connection" {
		t.Errorf("Expected group name 'connection', got '%s'", group.Name)
	}

	if len(group.Flags) != 2 {
		t.Fatalf("Expected 2 flags in group, got %d", len(group.Flags))
	}

	if group.Flags[0] != "host" || group.Flags[1] != "port" {
		t.Errorf("Expected flags ['host', 'port'], got %v", group.Flags)
	}

	if group.Conditional != false {
		t.Errorf("Expected Conditional to be false, got %v", group.Conditional)
	}

	if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err == nil {
		t.Error("Expected error when adding duplicate required group")
	}

	err := cmd.Parse([]string{"--host", "localhost"})
	if err == nil {
		t.Error("Expected error when not all required flags are set")
	}

	err = cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
	if err != nil {
		t.Errorf("Expected no error when all required flags are set, got: %v", err)
	}
}

func TestAddRequiredGroupWithEmptyFlags(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddRequiredGroup("empty", []string{}, false); err == nil {
		t.Error("Expected error when adding required group with empty flags")
	}
}

func TestAddRequiredGroupWithNonExistentFlag(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host address", "")); err != nil {
		t.Fatalf("Failed to add host flag: %v", err)
	}

	if err := cmd.AddRequiredGroup("connection", []string{"host", "nonexistent"}, false); err == nil {
		t.Error("Expected error when adding required group with non-existent flag")
	}
}

func TestRequiredGroupValidation(t *testing.T) {
	func() {
		cmd := NewCmd("test1", "t1", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080"})
		if err != nil {
			t.Errorf("Expected no error when all required flags are set, got: %v", err)
		}
	}()

	func() {
		cmd := NewCmd("test2", "t2", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost"})
		if err == nil {
			t.Error("Expected error when not all required flags are set")
		}
	}()

	func() {
		cmd := NewCmd("test3", "t3", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{})
		if err == nil {
			t.Error("Expected error when no required flags are set")
		}
	}()

	func() {
		cmd := NewCmd("test4", "t4", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("host", "H", "Host", "")); err != nil {
			t.Fatalf("Failed to add host flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("port", "P", "Port", "")); err != nil {
			t.Fatalf("Failed to add port flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("username", "U", "Username", "")); err != nil {
			t.Fatalf("Failed to add username flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("password", "W", "Password", "")); err != nil {
			t.Fatalf("Failed to add password flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("connection", []string{"host", "port"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}
		if err := cmd.AddRequiredGroup("auth", []string{"username", "password"}, false); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--host", "localhost", "--port", "8080", "--username", "admin", "--password", "secret"})
		if err != nil {
			t.Errorf("Expected no error when all required flags are set, got: %v", err)
		}
	}()
}

func TestConditionalRequiredGroup(t *testing.T) {
	func() {
		cmd := NewCmd("test1", "t1", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
			t.Fatalf("Failed to add ssl flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("cert", "c", "Certificate file", "")); err != nil {
			t.Fatalf("Failed to add cert flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("key", "k", "Key file", "")); err != nil {
			t.Fatalf("Failed to add key flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("ssl_config", []string{"cert", "key"}, true); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{})
		if err != nil {
			t.Errorf("Expected no error when no flags are set, got: %v", err)
		}
	}()

	func() {
		cmd := NewCmd("test2", "t2", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
			t.Fatalf("Failed to add ssl flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("cert", "c", "Certificate file", "")); err != nil {
			t.Fatalf("Failed to add cert flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("key", "k", "Key file", "")); err != nil {
			t.Fatalf("Failed to add key flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("ssl_config", []string{"cert", "key"}, true); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--ssl"})
		if err != nil {
			t.Errorf("Expected no error when only ssl is set, got: %v", err)
		}
	}()

	func() {
		cmd := NewCmd("test3", "t3", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
			t.Fatalf("Failed to add ssl flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("cert", "c", "Certificate file", "")); err != nil {
			t.Fatalf("Failed to add cert flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("key", "k", "Key file", "")); err != nil {
			t.Fatalf("Failed to add key flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("ssl_config", []string{"cert", "key"}, true); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--cert", "cert.pem"})
		if err == nil {
			t.Error("Expected error when only cert is set in conditional group")
		}
	}()

	func() {
		cmd := NewCmd("test4", "t4", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewBoolFlag("ssl", "s", "Use SSL", false)); err != nil {
			t.Fatalf("Failed to add ssl flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("cert", "c", "Certificate file", "")); err != nil {
			t.Fatalf("Failed to add cert flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("key", "k", "Key file", "")); err != nil {
			t.Fatalf("Failed to add key flag: %v", err)
		}
		if err := cmd.AddRequiredGroup("ssl_config", []string{"cert", "key"}, true); err != nil {
			t.Fatalf("Failed to add required group: %v", err)
		}

		err := cmd.Parse([]string{"--cert", "cert.pem", "--key", "key.pem"})
		if err != nil {
			t.Errorf("Expected no error when all conditional flags are set, got: %v", err)
		}
	}()
}
