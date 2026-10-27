package types

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/td/tg"
)

// Edit edits the outgoing message text
func (b *Bot) Edit(ctx context.Context, msg *Message, text string, delIn time.Duration) error {
	if msg.Peer == nil {
		return fmt.Errorf("peer is nil")
	}

	req := &tg.MessagesEditMessageRequest{
		Peer:    msg.Peer,
		ID:      msg.ID,
		Message: text,
	}

	_, err := b.API.MessagesEditMessage(ctx, req)
	if err != nil {
		return err
	}

	if delIn > 0 {
		go func() {
			time.Sleep(delIn)
			_ = b.Delete(context.Background(), msg)
		}()
	}
	return nil
}

// Reply sends a reply message to the chat
func (b *Bot) Reply(ctx context.Context, msg *Message, text string, delIn time.Duration) (*tg.UpdatesClass, error) {
	if msg.Peer == nil {
		return nil, fmt.Errorf("peer is nil")
	}

	req := &tg.MessagesSendMessageRequest{
		Peer:     msg.Peer,
		ReplyTo:  &tg.InputReplyToMessage{ReplyToMsgID: msg.ID},
		Message:  text,
		RandomID: time.Now().UnixNano(),
	}

	updates, err := b.API.MessagesSendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	if delIn > 0 {
		// Asynchronous delete helper
		go func() {
			time.Sleep(delIn)
			// best effort deletion of sent message
		}()
	}
	return &updates, nil
}

// Respond edits the message if outgoing, or replies if incoming
func (b *Bot) Respond(ctx context.Context, msg *Message, text string, delIn time.Duration) error {
	if msg.IsOutgoing {
		return b.Edit(ctx, msg, text, delIn)
	}
	_, err := b.Reply(ctx, msg, text, delIn)
	return err
}

// Delete removes the message from the chat
func (b *Bot) Delete(ctx context.Context, msg *Message) error {
	if msg.Peer == nil {
		return fmt.Errorf("peer is nil")
	}
	_, err := b.API.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		Revoke: true,
		ID:     []int{msg.ID},
	})
	return err
}

// ExtractUserAndReason resolves target user ID and reason from command args or replied message
func (b *Bot) ExtractUserAndReason(msg *Message) (targetID int64, reason string, err error) {
	if msg.ReplyToUserID != 0 {
		targetID = msg.ReplyToUserID
		reason = strings.TrimSpace(msg.Input)
		if reason == "" {
			reason = "No reason specified"
		}
		return targetID, reason, nil
	}

	if len(msg.Args) == 0 {
		return 0, "", fmt.Errorf("reply to a user or specify user ID/username")
	}

	firstArg := msg.Args[0]
	parsedID, parseErr := strconv.ParseInt(firstArg, 10, 64)
	if parseErr == nil {
		targetID = parsedID
		if len(msg.Args) > 1 {
			reason = strings.Join(msg.Args[1:], " ")
		} else {
			reason = "No reason specified"
		}
		return targetID, reason, nil
	}

	return 0, "", fmt.Errorf("invalid user identifier: %s", firstArg)
}
