package types

import (
	"testing"
)

func TestExtractUserAndReasonFromArgs(t *testing.T) {
	b := &Bot{}

	msgNumeric := &Message{
		Args: []string{"12345678", "Spamming", "in", "chat"},
	}

	targetID, reason, err := b.ExtractUserAndReason(msgNumeric)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if targetID != 12345678 {
		t.Errorf("expected targetID 12345678, got %d", targetID)
	}
	if reason != "Spamming in chat" {
		t.Errorf("expected reason 'Spamming in chat', got '%s'", reason)
	}

	msgReplied := &Message{
		ReplyToUserID: 87654321,
		Input:         "Violation of rules",
	}

	targetID2, reason2, err := b.ExtractUserAndReason(msgReplied)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if targetID2 != 87654321 {
		t.Errorf("expected targetID 87654321, got %d", targetID2)
	}
	if reason2 != "Violation of rules" {
		t.Errorf("expected reason 'Violation of rules', got '%s'", reason2)
	}
}
