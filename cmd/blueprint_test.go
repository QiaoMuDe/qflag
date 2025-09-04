package cmd

import (
	"strings"
	"testing"

	"gitee.com/MM-Q/qflag/flags"
)

func TestNewCmdFromBlueprint(t *testing.T) {
	t.Run("Create simple command without flags or subcommands", func(t *testing.T) {
		bp := &CommandBlueprint{
			Name:        "test",
			ShortName:   "t",
			Usage:       "A test command",
			Description: "This is a test command.",
			Action:      func(c *Cmd) error { return nil },
		}

		cmd, err := NewCmdFromBlueprint(bp)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		if cmd == nil {
			t.Fatal("Expected command to be not nil")
		}

		if cmd.LongName() != "test" {
			t.Errorf("Expected long name 'test', got '%s'", cmd.LongName())
		}
		if cmd.ShortName() != "t" {
			t.Errorf("Expected short name 't', got '%s'", cmd.ShortName())
		}
		if cmd.GetUsageSyntax() != "A test command" {
			t.Errorf("Expected usage syntax 'A test command', got '%s'", cmd.GetUsageSyntax())
		}
		if cmd.GetDescription() != "This is a test command." {
			t.Errorf("Expected description 'This is a test command.', got '%s'", cmd.GetDescription())
		}
		if len(cmd.SubCmds()) != 0 {
			t.Errorf("Expected 0 subcommands, got %d", len(cmd.SubCmds()))
		}
		// By default, a 'help' flag is added.
		if cmd.NFlag() != 1 {
			t.Errorf("Expected 1 flag (for help), got %d", cmd.NFlag())
		}
	})

	t.Run("Create command with various flag types", func(t *testing.T) {
		bp := &CommandBlueprint{
			Name: "test-with-flags",
			Flags: []FlagBlueprint{
				{Type: flags.FlagTypeString, Name: "name", ShortName: "n", Usage: "A name", DefaultValue: "default-name"},
				{Type: flags.FlagTypeInt, Name: "age", Usage: "An age", DefaultValue: 30},
				{Type: flags.FlagTypeBool, Name: "verbose", ShortName: "v", DefaultValue: false},
				{Type: flags.FlagTypeSize, Name: "memory", ShortName: "m", Usage: "Memory size", DefaultValue: int64(1024 * 1024)},
			},
		}

		cmd, err := NewCmdFromBlueprint(bp)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		// Expect 4 flags from blueprint + 1 default 'help' flag.
		if cmd.NFlag() != 5 {
			t.Fatalf("Expected 5 flags, but got %d", cmd.NFlag())
		}

		// Check string flag
		nameFlag, ok := cmd.FlagRegistry().GetByName("name")
		if !ok {
			t.Fatal("Flag 'name' not found")
		}
		if nameFlag.GetLongName() != "name" || nameFlag.GetShortName() != "n" || nameFlag.GetUsage() != "A name" || nameFlag.GetDefault() != "default-name" || nameFlag.GetFlagType() != flags.FlagTypeString {
			t.Errorf("String flag properties do not match. Got: %+v", nameFlag)
		}

		// Check int flag
		ageFlag, ok := cmd.FlagRegistry().GetByName("age")
		if !ok {
			t.Fatal("Flag 'age' not found")
		}
		if ageFlag.GetLongName() != "age" || ageFlag.GetDefault() != 30 || ageFlag.GetFlagType() != flags.FlagTypeInt {
			t.Errorf("Int flag properties do not match. Got: %+v", ageFlag)
		}

		// Check bool flag by short name
		verboseFlag, ok := cmd.FlagRegistry().GetByName("v")
		if !ok {
			t.Fatal("Flag 'verbose' (v) not found")
		}
		if verboseFlag.GetLongName() != "verbose" || verboseFlag.GetDefault() != false || verboseFlag.GetFlagType() != flags.FlagTypeBool {
			t.Errorf("Bool flag properties do not match. Got: %+v", verboseFlag)
		}

		// Check size flag
		sizeFlag, ok := cmd.FlagRegistry().GetByName("memory")
		if !ok {
			t.Fatal("Flag 'memory' not found")
		}
		if sizeFlag.GetLongName() != "memory" || sizeFlag.GetShortName() != "m" || sizeFlag.GetDefault() != int64(1024*1024) || sizeFlag.GetFlagType() != flags.FlagTypeSize {
			t.Errorf("Size flag properties do not match. Got: %+v", sizeFlag)
		}
	})

	t.Run("Create command with nested subcommands and flags", func(t *testing.T) {
		bp := &CommandBlueprint{
			Name: "root",
			Flags: []FlagBlueprint{
				{Type: flags.FlagTypeBool, Name: "global"},
			},
			Subcommands: []*CommandBlueprint{
				{
					Name:  "sub1",
					Usage: "First subcommand",
					Flags: []FlagBlueprint{{Type: flags.FlagTypeString, Name: "sub1-flag"}},
				},
				{
					Name:      "sub2",
					ShortName: "s2",
					Subcommands: []*CommandBlueprint{
						{
							Name:  "sub2-child",
							Flags: []FlagBlueprint{{Type: flags.FlagTypeInt, Name: "child-flag"}},
						},
					},
				},
			},
		}

		rootCmd, err := NewCmdFromBlueprint(bp)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		if !rootCmd.FlagExists("global") {
			t.Error("Expected 'global' flag to exist on root command")
		}
		if len(rootCmd.SubCmds()) != 2 {
			t.Fatalf("Expected 2 subcommands, got %d", len(rootCmd.SubCmds()))
		}

		// Check sub1
		sub1, ok := rootCmd.SubCmdMap()["sub1"]
		if !ok {
			t.Fatal("Subcommand 'sub1' not found")
		}
		if sub1.Name() != "sub1" || sub1.GetUsageSyntax() != "First subcommand" || !sub1.FlagExists("sub1-flag") {
			t.Errorf("Properties of subcommand 'sub1' are incorrect")
		}

		// Check sub2 and its child
		sub2, ok := rootCmd.SubCmdMap()["sub2"]
		if !ok {
			t.Fatal("Subcommand 'sub2' not found")
		}
		if len(sub2.SubCmds()) != 1 {
			t.Fatalf("Expected 1 subcommand for 'sub2', got %d", len(sub2.SubCmds()))
		}
		sub2child, ok := sub2.SubCmdMap()["sub2-child"]
		if !ok {
			t.Fatal("Subcommand 'sub2-child' not found")
		}
		if sub2child.Name() != "sub2-child" || !sub2child.FlagExists("child-flag") {
			t.Errorf("Properties of subcommand 'sub2-child' are incorrect")
		}
	})

	t.Run("Error on nil blueprint", func(t *testing.T) {
		_, err := NewCmdFromBlueprint(nil)
		if err == nil {
			t.Fatal("Expected an error for nil blueprint, but got nil")
		}
		expectedErr := "command blueprint cannot be nil"
		if err.Error() != expectedErr {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErr, err.Error())
		}
	})

	t.Run("Error on unsupported flag type", func(t *testing.T) {
		bp := &CommandBlueprint{
			Name:  "test-unsupported",
			Flags: []FlagBlueprint{{Type: flags.FlagTypeUnknown, Name: "bad-flag"}},
		}
		_, err := NewCmdFromBlueprint(bp)
		if err == nil {
			t.Fatal("Expected an error for unsupported flag type, but got nil")
		}
		if !strings.Contains(err.Error(), "unsupported flag type") {
			t.Errorf("Expected error message to contain 'unsupported flag type', but got '%s'", err.Error())
		}
	})
}
