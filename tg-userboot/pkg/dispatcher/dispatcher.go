package dispatcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gotd/td/tg"
	"tg-userboot/pkg/types"
)

// Dispatcher manages command registration and routing
type Dispatcher struct {
	mu       sync.RWMutex
	commands map[string]*types.Command
	aliases  map[string]string
	bot      *types.Bot
}

// NewDispatcher creates a new Dispatcher instance
func NewDispatcher(b *types.Bot) *Dispatcher {
	d := &Dispatcher{
		commands: make(map[string]*types.Command),
		aliases:  make(map[string]string),
		bot:      b,
	}
	b.Commands = d.commands
	return d
}

// Register registers a command and its aliases
func (d *Dispatcher) Register(cmd *types.Command) {
	d.mu.Lock()
	defer d.mu.Unlock()

	name := strings.ToLower(cmd.Name)
	d.commands[name] = cmd

	for _, alias := range cmd.Aliases {
		d.aliases[strings.ToLower(alias)] = name
	}
}

// HandleUpdate handles raw updates from Telegram MTProto
func (d *Dispatcher) HandleUpdate(ctx context.Context, u tg.UpdateClass) error {
	var msg *tg.Message

	switch update := u.(type) {
	case *tg.UpdateNewMessage:
		if m, ok := update.Message.(*tg.Message); ok {
			msg = m
		}
	case *tg.UpdateNewChannelMessage:
		if m, ok := update.Message.(*tg.Message); ok {
			msg = m
		}
	case *tg.UpdateEditMessage:
		if m, ok := update.Message.(*tg.Message); ok {
			msg = m
		}
	case *tg.UpdateEditChannelMessage:
		if m, ok := update.Message.(*tg.Message); ok {
			msg = m
		}
	}

	if msg == nil || msg.Message == "" {
		return nil
	}

	return d.ProcessMessage(ctx, msg)
}

// ProcessMessage parses and executes a command if message matches trigger
func (d *Dispatcher) ProcessMessage(ctx context.Context, raw *tg.Message) error {
	text := strings.TrimSpace(raw.Message)
	if text == "" {
		return nil
	}

	ownerTrigger := d.bot.Config.CmdTrigger
	sudoTrigger := d.bot.Config.SudoTrigger

	var trigger string
	var isOwnerTrigger bool
	var isSudoTrigger bool

	if strings.HasPrefix(text, ownerTrigger) {
		trigger = ownerTrigger
		isOwnerTrigger = true
	} else if sudoTrigger != "" && strings.HasPrefix(text, sudoTrigger) {
		trigger = sudoTrigger
		isSudoTrigger = true
	} else {
		return nil
	}

	// Permission checks
	var senderID int64
	if from, ok := raw.GetFromID(); ok {
		switch peer := from.(type) {
		case *tg.PeerUser:
			senderID = peer.UserID
		}
	} else if raw.Out && d.bot.Self != nil {
		senderID = d.bot.Self.ID
	}

	isOwner := raw.Out || (d.bot.Self != nil && senderID == d.bot.Self.ID) || (d.bot.Config.OwnerID != 0 && senderID == d.bot.Config.OwnerID)
	isSudo := false
	for _, sudoID := range d.bot.Config.SudoUsers {
		if senderID == sudoID {
			isSudo = true
			break
		}
	}

	if isOwnerTrigger && !isOwner {
		return nil
	}
	if isSudoTrigger && !isOwner && !isSudo {
		return nil
	}

	// Parse command name and arguments
	body := strings.TrimPrefix(text, trigger)
	parts := strings.Fields(body)
	if len(parts) == 0 {
		return nil
	}

	cmdName := strings.ToLower(parts[0])
	args := parts[1:]
	input := strings.TrimSpace(strings.TrimPrefix(body, parts[0]))

	flags := make(map[string]bool)
	cleanArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 && !isNumber(arg) {
			flags[strings.TrimPrefix(arg, "-")] = true
		} else {
			cleanArgs = append(cleanArgs, arg)
		}
	}

	// Resolve command from aliases if needed
	d.mu.RLock()
	resolvedName := cmdName
	if actual, ok := d.aliases[cmdName]; ok {
		resolvedName = actual
	}
	cmd, exists := d.commands[resolvedName]
	d.mu.RUnlock()

	if !exists {
		return nil
	}

	// Check command-specific permissions
	if cmd.OwnerOnly && !isOwner {
		return nil
	}
	if !cmd.SudoAllowed && isSudo && !isOwner {
		return nil
	}

	// Resolve Peer
	peer := raw.GetPeerID()
	var inputPeer tg.InputPeerClass
	switch p := peer.(type) {
	case *tg.PeerUser:
		inputPeer = &tg.InputPeerUser{UserID: p.UserID}
	case *tg.PeerChat:
		inputPeer = &tg.InputPeerChat{ChatID: p.ChatID}
	case *tg.PeerChannel:
		inputPeer = &tg.InputPeerChannel{ChannelID: p.ChannelID}
	}

	replyToMsgID := 0
	if raw.ReplyTo != nil {
		if replyHeader, ok := raw.ReplyTo.(*tg.MessageReplyHeader); ok {
			replyToMsgID = replyHeader.ReplyToMsgID
		}
	}

	wrappedMsg := &types.Message{
		ID:           raw.ID,
		SenderID:     senderID,
		IsOutgoing:   raw.Out,
		Text:         text,
		RawText:      raw.Message,
		Trigger:      trigger,
		Cmd:          cmdName,
		Input:        input,
		Args:         cleanArgs,
		Flags:        flags,
		ReplyToMsgID: replyToMsgID,
		Peer:         inputPeer,
		RawMessage:   raw,
	}

	// Execute command safely with recover
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] in command %s: %v", cmdName, r)
				_ = d.bot.Respond(ctx, wrappedMsg, fmt.Sprintf("⚠️ <b>Internal Error</b> in <code>%s</code>: %v", cmdName, r), 10*time.Second)
			}
		}()

		if err := cmd.Handler(ctx, d.bot, wrappedMsg); err != nil {
			log.Printf("[ERROR] Command %s failed: %v", cmdName, err)
			_ = d.bot.Respond(ctx, wrappedMsg, fmt.Sprintf("❌ <b>Error:</b> %v", err), 10*time.Second)
		}
	}()

	return nil
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		start = 1
	}
	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return len(s) > start
}
