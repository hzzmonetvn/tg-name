package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/tg"
	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all group administration commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "ban",
		Aliases:     []string{"sban"},
		Category:    "Admin",
		Description: "Ban a user from the chat",
		Usage:       ".ban [reply or user_id] [reason]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handleBan,
	})

	d.Register(&types.Command{
		Name:        "unban",
		Category:    "Admin",
		Description: "Unban a user in the chat",
		Usage:       ".unban [reply or user_id]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handleUnban,
	})

	d.Register(&types.Command{
		Name:        "kick",
		Aliases:     []string{"skick", "dkick"},
		Category:    "Admin",
		Description: "Kick a user from the chat",
		Usage:       ".kick [reply or user_id] [reason]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handleKick,
	})

	d.Register(&types.Command{
		Name:        "mute",
		Aliases:     []string{"smute"},
		Category:    "Admin",
		Description: "Mute a user in the chat",
		Usage:       ".mute [reply or user_id] [reason]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handleMute,
	})

	d.Register(&types.Command{
		Name:        "unmute",
		Category:    "Admin",
		Description: "Unmute a user in the chat",
		Usage:       ".unmute [reply or user_id]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handleUnmute,
	})

	d.Register(&types.Command{
		Name:        "del",
		Aliases:     []string{"purge"},
		Category:    "Admin",
		Description: "Delete replied message",
		Usage:       ".del [reply]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handleDel,
	})

	d.Register(&types.Command{
		Name:        "pin",
		Aliases:     []string{"unpin"},
		Category:    "Admin",
		Description: "Pin or unpin replied message",
		Usage:       ".pin [reply]",
		OwnerOnly:   true,
		SudoAllowed: true,
		Handler:     handlePin,
	})
}

func handleBan(ctx context.Context, b *types.Bot, msg *types.Message) error {
	ch, ok := msg.Peer.(*tg.InputPeerChannel)
	if !ok {
		return b.Respond(ctx, msg, "❌ Ban can only be used in supergroups/channels.", 5*time.Second)
	}

	targetID, reason, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ %v", err), 5*time.Second)
	}

	_, err = b.API.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel: &tg.InputChannel{ChannelID: ch.ChannelID},
		Participant: &tg.InputPeerUser{UserID: targetID},
		BannedRights: tg.ChatBannedRights{
			ViewMessages: true,
			SendMessages: true,
			SendMedia:    true,
			SendStickers: true,
			SendGifs:     true,
			SendGames:    true,
			SendInline:   true,
			EmbedLinks:   true,
			UntilDate:    0, // Forever
		},
	})
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to ban: %v", err), 5*time.Second)
	}

	return b.Respond(ctx, msg, fmt.Sprintf("🔨 <b>Banned:</b> <code>%d</code>\n📝 <b>Reason:</b> %s", targetID, reason), 0)
}

func handleUnban(ctx context.Context, b *types.Bot, msg *types.Message) error {
	ch, ok := msg.Peer.(*tg.InputPeerChannel)
	if !ok {
		return b.Respond(ctx, msg, "❌ Can only be used in supergroups/channels.", 5*time.Second)
	}

	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ %v", err), 5*time.Second)
	}

	_, err = b.API.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel:      &tg.InputChannel{ChannelID: ch.ChannelID},
		Participant:  &tg.InputPeerUser{UserID: targetID},
		BannedRights: tg.ChatBannedRights{},
	})
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to unban: %v", err), 5*time.Second)
	}

	return b.Respond(ctx, msg, fmt.Sprintf("🔓 <b>Unbanned:</b> <code>%d</code>", targetID), 5*time.Second)
}

func handleKick(ctx context.Context, b *types.Bot, msg *types.Message) error {
	ch, ok := msg.Peer.(*tg.InputPeerChannel)
	if !ok {
		return b.Respond(ctx, msg, "❌ Kick can only be used in supergroups/channels.", 5*time.Second)
	}

	targetID, reason, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ %v", err), 5*time.Second)
	}

	// Kick by temporarily banning and then lifting the ban
	_, err = b.API.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel: &tg.InputChannel{ChannelID: ch.ChannelID},
		Participant: &tg.InputPeerUser{UserID: targetID},
		BannedRights: tg.ChatBannedRights{
			ViewMessages: true,
			UntilDate:    int(time.Now().Add(60 * time.Second).Unix()),
		},
	})
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to kick: %v", err), 5*time.Second)
	}

	return b.Respond(ctx, msg, fmt.Sprintf("👢 <b>Kicked:</b> <code>%d</code>\n📝 <b>Reason:</b> %s", targetID, reason), 0)
}

func handleMute(ctx context.Context, b *types.Bot, msg *types.Message) error {
	ch, ok := msg.Peer.(*tg.InputPeerChannel)
	if !ok {
		return b.Respond(ctx, msg, "❌ Mute can only be used in supergroups/channels.", 5*time.Second)
	}

	targetID, reason, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ %v", err), 5*time.Second)
	}

	_, err = b.API.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel: &tg.InputChannel{ChannelID: ch.ChannelID},
		Participant: &tg.InputPeerUser{UserID: targetID},
		BannedRights: tg.ChatBannedRights{
			SendMessages: true,
			SendMedia:    true,
			SendStickers: true,
			SendGifs:     true,
			SendGames:    true,
			SendInline:   true,
			EmbedLinks:   true,
			UntilDate:    0,
		},
	})
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to mute: %v", err), 5*time.Second)
	}

	return b.Respond(ctx, msg, fmt.Sprintf("🔇 <b>Muted:</b> <code>%d</code>\n📝 <b>Reason:</b> %s", targetID, reason), 0)
}

func handleUnmute(ctx context.Context, b *types.Bot, msg *types.Message) error {
	ch, ok := msg.Peer.(*tg.InputPeerChannel)
	if !ok {
		return b.Respond(ctx, msg, "❌ Can only be used in supergroups/channels.", 5*time.Second)
	}

	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ %v", err), 5*time.Second)
	}

	_, err = b.API.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel:      &tg.InputChannel{ChannelID: ch.ChannelID},
		Participant:  &tg.InputPeerUser{UserID: targetID},
		BannedRights: tg.ChatBannedRights{},
	})
	if err != nil {
		return b.Respond(ctx, msg, fmt.Sprintf("❌ Failed to unmute: %v", err), 5*time.Second)
	}

	return b.Respond(ctx, msg, fmt.Sprintf("🔊 <b>Unmuted:</b> <code>%d</code>", targetID), 5*time.Second)
}

func handleDel(ctx context.Context, b *types.Bot, msg *types.Message) error {
	if msg.ReplyToMsgID == 0 {
		return b.Respond(ctx, msg, "❌ Reply to a message to delete it.", 5*time.Second)
	}

	targetIDs := []int{msg.ReplyToMsgID}
	if msg.IsOutgoing {
		targetIDs = append(targetIDs, msg.ID)
	}

	_, err := b.API.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		Revoke: true,
		ID:     targetIDs,
	})
	return err
}

func handlePin(ctx context.Context, b *types.Bot, msg *types.Message) error {
	if msg.ReplyToMsgID == 0 {
		return b.Respond(ctx, msg, "❌ Reply to a message to pin/unpin it.", 5*time.Second)
	}

	silent := msg.Flags["s"]
	isUnpin := msg.Cmd == "unpin"

	_, err := b.API.MessagesUpdatePinnedMessage(ctx, &tg.MessagesUpdatePinnedMessageRequest{
		Silent: silent,
		Unpin:  isUnpin,
		Peer:   msg.Peer,
		ID:     msg.ReplyToMsgID,
	})
	if err != nil {
		return err
	}
	if isUnpin {
		return b.Respond(ctx, msg, "📌 <i>Unpinned message.</i>", 3*time.Second)
	}
	return b.Respond(ctx, msg, "📌 <i>Pinned message.</i>", 3*time.Second)
}
