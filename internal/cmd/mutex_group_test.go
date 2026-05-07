package cmd

import (
	"testing"

	"gitee.com/MM-Q/qflag/internal/flag"
	"gitee.com/MM-Q/qflag/internal/types"
)

func TestAddMutexGroup(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
		t.Fatalf("Failed to add format flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
		t.Fatalf("Failed to add output flag: %v", err)
	}

	if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
		t.Fatalf("Failed to add mutex group: %v", err)
	}

	groups := cmd.Config().MutexGroups
	if len(groups) != 1 {
		t.Fatalf("Expected 1 mutex group, got %d", len(groups))
	}

	group := groups[0]
	if group.Name != "output_format" {
		t.Errorf("Expected group name 'output_format', got '%s'", group.Name)
	}

	if len(group.Flags) != 2 {
		t.Fatalf("Expected 2 flags in group, got %d", len(group.Flags))
	}

	if group.Flags[0] != "format" || group.Flags[1] != "output" {
		t.Errorf("Expected flags ['format', 'output'], got %v", group.Flags)
	}

	if !group.AllowNone {
		t.Errorf("Expected AllowNone to be true, got false")
	}

	if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err == nil {
		t.Error("Expected error when adding duplicate mutex group")
	}
}

func TestAddMutexGroupWithEmptyFlags(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddMutexGroup("empty", []string{}, true); err == nil {
		t.Error("Expected error when adding mutex group with empty flags")
	}
}

func TestAddMutexGroupWithNonExistentFlag(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
		t.Fatalf("Failed to add format flag: %v", err)
	}

	if err := cmd.AddMutexGroup("test_group", []string{"format", "nonexistent"}, true); err == nil {
		t.Error("Expected error when adding mutex group with non-existent flag")
	}
}

func TestMutexGroupValidation(t *testing.T) {
	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
			t.Fatalf("Failed to add mutex group: %v", err)
		}

		err := cmd.Parse([]string{})
		if err != nil {
			t.Errorf("Expected no error when no flags are set with AllowNone=true, got: %v", err)
		}
	}()

	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
			t.Fatalf("Failed to add mutex group: %v", err)
		}

		err := cmd.Parse([]string{"--format", "json"})
		if err != nil {
			t.Errorf("Expected no error when only one flag is set, got: %v", err)
		}
	}()

	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
			t.Fatalf("Failed to add mutex group: %v", err)
		}

		err := cmd.Parse([]string{"--format", "json", "--output", "file.txt"})
		if err == nil {
			t.Error("Expected error when both mutex flags are set")
		}
	}()

	func() {
		cmd := NewCmd("test", "t", types.ContinueOnError)
		if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "json")); err != nil {
			t.Fatalf("Failed to add format flag: %v", err)
		}
		if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
			t.Fatalf("Failed to add output flag: %v", err)
		}
		if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, false); err != nil {
			t.Fatalf("Failed to add mutex group: %v", err)
		}

		err := cmd.Parse([]string{})
		if err == nil {
			t.Error("Expected error when no flags are set with AllowNone=false")
		}
	}()
}

func TestMutexGroupWithMultipleGroups(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddFlag(flag.NewBoolFlag("json", "j", "JSON format", false)); err != nil {
		t.Fatalf("Failed to add json flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewBoolFlag("xml", "x", "XML format", false)); err != nil {
		t.Fatalf("Failed to add xml flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewBoolFlag("compress", "c", "Compress output", false)); err != nil {
		t.Fatalf("Failed to add compress flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewBoolFlag("encrypt", "e", "Encrypt output", false)); err != nil {
		t.Fatalf("Failed to add encrypt flag: %v", err)
	}

	if err := cmd.AddMutexGroup("format", []string{"json", "xml"}, true); err != nil {
		t.Fatalf("Failed to add format mutex group: %v", err)
	}
	if err := cmd.AddMutexGroup("security", []string{"compress", "encrypt"}, true); err != nil {
		t.Fatalf("Failed to add security mutex group: %v", err)
	}

	groups := cmd.Config().MutexGroups
	if len(groups) != 2 {
		t.Fatalf("Expected 2 mutex groups, got %d", len(groups))
	}

	err := cmd.Parse([]string{"--json", "--compress"})
	if err != nil {
		t.Errorf("Expected no error when using flags from different groups, got: %v", err)
	}

	cmd2 := NewCmd("test2", "t2", types.ContinueOnError)
	if err := cmd2.AddFlag(flag.NewBoolFlag("json", "j", "JSON format", false)); err != nil {
		t.Fatalf("Failed to add json flag: %v", err)
	}
	if err := cmd2.AddFlag(flag.NewBoolFlag("xml", "x", "XML format", false)); err != nil {
		t.Fatalf("Failed to add xml flag: %v", err)
	}
	if err := cmd2.AddMutexGroup("format", []string{"json", "xml"}, true); err != nil {
		t.Fatalf("Failed to add format mutex group: %v", err)
	}

	err = cmd2.Parse([]string{"--json", "--xml"})
	if err == nil {
		t.Error("Expected error when using flags from same mutex group")
	}
}

func TestMutexGroupWithEnvVar(t *testing.T) {
	cmd := NewCmd("test", "t", types.ContinueOnError)

	if err := cmd.AddFlag(flag.NewStringFlag("format", "f", "Output format", "")); err != nil {
		t.Fatalf("Failed to add format flag: %v", err)
	}
	if err := cmd.AddFlag(flag.NewStringFlag("output", "o", "Output file", "")); err != nil {
		t.Fatalf("Failed to add output flag: %v", err)
	}

	if err := cmd.AddMutexGroup("output_format", []string{"format", "output"}, true); err != nil {
		t.Fatalf("Failed to add mutex group: %v", err)
	}

	groups := cmd.Config().MutexGroups
	if len(groups) != 1 {
		t.Fatalf("Expected 1 mutex group, got %d", len(groups))
	}
}
