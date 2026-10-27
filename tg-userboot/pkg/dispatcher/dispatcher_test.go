package dispatcher

import (
	"context"
	"testing"

	"tg-userboot/pkg/types"
)

func TestDispatcherRegisterAndLookup(t *testing.T) {
	bot := &types.Bot{
		Config: &types.Config{
			CmdTrigger:  ".",
			SudoTrigger: "!",
		},
	}
	disp := NewDispatcher(bot)

	called := false
	cmd := &types.Command{
		Name:        "testcmd",
		Aliases:     []string{"tc", "talias"},
		Category:    "Test",
		Description: "A unit test command",
		Handler: func(ctx context.Context, b *types.Bot, msg *types.Message) error {
			called = true
			return nil
		},
	}

	disp.Register(cmd)

	if len(disp.commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(disp.commands))
	}
	if len(disp.aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(disp.aliases))
	}
	if disp.aliases["tc"] != "testcmd" {
		t.Errorf("expected alias 'tc' to resolve to 'testcmd', got %s", disp.aliases["tc"])
	}
	if disp.aliases["talias"] != "testcmd" {
		t.Errorf("expected alias 'talias' to resolve to 'testcmd', got %s", disp.aliases["talias"])
	}

	_ = called
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"-456", true},
		{"+789", true},
		{"abc", false},
		{"-", false},
		{"", false},
		{"12a3", false},
	}

	for _, tc := range tests {
		res := isNumber(tc.input)
		if res != tc.expected {
			t.Errorf("isNumber(%q) = %v; want %v", tc.input, res, tc.expected)
		}
	}
}
