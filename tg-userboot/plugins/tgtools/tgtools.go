package tgtools

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/gotd/td/tg"
	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all Telegram tool commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "ping",
		Aliases:     []string{"p"},
		Category:    "TGTools",
		Description: "Check bot latency and response time",
		Usage:       ".ping",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handlePing,
	})

	d.Register(&types.Command{
		Name:        "alive",
		Aliases:     []string{"status"},
		Category:    "TGTools",
		Description: "Show userbot uptime and system status",
		Usage:       ".alive",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleAlive,
	})

	d.Register(&types.Command{
		Name:        "id",
		Aliases:     []string{"chatid"},
		Category:    "TGTools",
		Description: "Get user ID or chat ID",
		Usage:       ".id [reply or blank]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleID,
	})

	d.Register(&types.Command{
		Name:        "info",
		Aliases:     []string{"whois", "user"},
		Category:    "TGTools",
		Description: "Fetch detailed Telegram user information",
		Usage:       ".info [reply or username/id]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleInfo,
	})

	d.Register(&types.Command{
		Name:        "leave",
		Aliases:     []string{"kickme"},
		Category:    "TGTools",
		Description: "Leave the current group or channel",
		Usage:       ".leave",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleLeave,
	})

	d.Register(&types.Command{
		Name:        "block",
		Category:    "TGTools",
		Description: "Block a user in private messages",
		Usage:       ".block [reply or user_id]",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleBlock,
	})

	d.Register(&types.Command{
		Name:        "unblock",
		Category:    "TGTools",
		Description: "Unblock a user",
		Usage:       ".unblock [reply or user_id]",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleUnblock,
	})
}

func handlePing(ctx context.Context, b *types.Bot, msg *types.Message) error {
	start := time.Now()
	if err := b.Respond(ctx, msg, "⚡ <code>Pinging...</code>", 0); err != nil {
		return err
	}
	latency := time.Since(start).Milliseconds()
	return b.Edit(ctx, msg, fmt.Sprintf("🏓 <b>Pong!</b> <code>%d ms</code>", latency), 0)
}

func handleAlive(ctx context.Context, b *types.Bot, msg *types.Message) error {
	uptime := time.Since(b.Start).Round(time.Second)
	ownerName := "Owner"
	if b.Self != nil {
		ownerName = b.Self.FirstName
	}

	text := fmt.Sprintf(`🚀 <b>TG-Userboot is Alive!</b>
━━━━━━━━━━━━━━━━━
👤 <b>Master:</b> <code>%s</code>
⏱️ <b>Uptime:</b> <code>%s</code>
⚙️ <b>Architecture:</b> <code>%s/%s</code>
🐹 <b>Go Version:</b> <code>%s</code>
📦 <b>Loaded Commands:</b> <code>%d</code>
✨ <b>Core:</b> <code>gotd/td (Pure Go MTProto)</code>
━━━━━━━━━━━━━━━━━`,
		ownerName,
		uptime.String(),
		runtime.GOOS, runtime.GOARCH,
		runtime.Version(),
		len(b.Commands),
	)

	return b.Respond(ctx, msg, text, 0)
}

func handleID(ctx context.Context, b *types.Bot, msg *types.Message) error {
	var chatIDStr string
	if msg.Peer != nil {
		switch p := msg.Peer.(type) {
		case *tg.InputPeerUser:
			chatIDStr = fmt.Sprintf("<code>%d</code>", p.UserID)
		case *tg.InputPeerChat:
			chatIDStr = fmt.Sprintf("<code>-%d</code>", p.ChatID)
		case *tg.InputPeerChannel:
			chatIDStr = fmt.Sprintf("<code>-100%d</code>", p.ChannelID)
		default:
			chatIDStr = "<code>Unknown</code>"
		}
	}

	res := fmt.Sprintf("💬 <b>Chat ID:</b> %s", chatIDStr)
	if msg.ReplyToUserID != 0 {
		res += fmt.Sprintf("\n👤 <b>Replied User ID:</b> <code>%d</code>", msg.ReplyToUserID)
	}
	res += fmt.Sprintf("\n🆔 <b>Message ID:</b> <code>%d</code>", msg.ID)

	return b.Respond(ctx, msg, res, 0)
}

func handleInfo(ctx context.Context, b *types.Bot, msg *types.Message) error {
	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil && msg.Peer != nil {
		if p, ok := msg.Peer.(*tg.InputPeerUser); ok {
			targetID = p.UserID
		}
	}

	if targetID == 0 {
		return b.Respond(ctx, msg, "❌ <i>Please reply to a user or provide user ID/username.</i>", 5*time.Second)
	}

	users, err := b.API.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUser{UserID: targetID}})
	if err != nil || len(users) == 0 {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Could not fetch user: %v", err), 5*time.Second)
	}

	u, ok := users[0].(*tg.User)
	if !ok {
		return b.Respond(ctx, msg, "❌ User not found", 5*time.Second)
	}

	username := "None"
	if u.Username != "" {
		username = "@" + u.Username
	}

	status := "Unknown"
	if u.Bot {
		status = "🤖 Bot"
	} else if u.Verified {
		status = "✅ Verified User"
	} else {
		status = "👤 User"
	}

	infoText := fmt.Sprintf(`👤 <b>User Info</b>
━━━━━━━━━━━━━━━━━
🆔 <b>User ID:</b> <code>%d</code>
📛 <b>Name:</b> <code>%s %s</code>
🌐 <b>Username:</b> %s
🏷️ <b>Type:</b> %s
🛡️ <b>Restricted:</b> <code>%t</code>
━━━━━━━━━━━━━━━━━`,
		u.ID,
		u.FirstName, u.LastName,
		username,
		status,
		u.Restricted,
	)

	return b.Respond(ctx, msg, infoText, 0)
}

func handleLeave(ctx context.Context, b *types.Bot, msg *types.Message) error {
	_ = b.Respond(ctx, msg, "👋 <i>Leaving chat...</i>", 3*time.Second)
	switch p := msg.Peer.(type) {
	case *tg.InputPeerChannel:
		_, err := b.API.ChannelsLeaveChannel(ctx, &tg.InputChannel{
			ChannelID: p.ChannelID,
		})
		return err
	case *tg.InputPeerChat:
		_, err := b.API.MessagesDeleteChatUser(ctx, &tg.MessagesDeleteChatUserRequest{
			ChatID: p.ChatID,
			UserID: &tg.InputUserSelf{},
		})
		return err
	default:
		return b.Respond(ctx, msg, "❌ Can only leave groups or channels.", 5*time.Second)
	}
}

func handleBlock(ctx context.Context, b *types.Bot, msg *types.Message) error {
	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, "❌ Reply to a user or specify user ID to block.", 5*time.Second)
	}

	_, err = b.API.ContactsBlock(ctx, &tg.ContactsBlockRequest{
		ID: &tg.InputPeerUser{UserID: targetID},
	})
	if err != nil {
		return err
	}
	return b.Respond(ctx, msg, fmt.Sprintf("🚫 <b>Blocked user:</b> <code>%d</code>", targetID), 5*time.Second)
}

func handleUnblock(ctx context.Context, b *types.Bot, msg *types.Message) error {
	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, "❌ Reply to a user or specify user ID to unblock.", 5*time.Second)
	}

	_, err = b.API.ContactsUnblock(ctx, &tg.ContactsUnblockRequest{
		ID: &tg.InputPeerUser{UserID: targetID},
	})
	if err != nil {
		return err
	}
	return b.Respond(ctx, msg, fmt.Sprintf("✅ <b>Unblocked user:</b> <code>%d</code>", targetID), 5*time.Second)
}
