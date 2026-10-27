package files

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers file management commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "spoiler",
		Category:    "Files",
		Description: "Wrap text in spoiler tags",
		Usage:       ".spoiler [text or reply]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleSpoiler,
	})
}

func handleSpoiler(ctx context.Context, b *types.Bot, msg *types.Message) error {
	text := strings.TrimSpace(msg.Input)
	if text == "" && msg.ReplyToText != "" {
		text = msg.ReplyToText
	}

	if text == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.spoiler [text]</code>", 3*time.Second)
	}

	output := fmt.Sprintf("<spoiler>%s</spoiler>", text)
	return b.Respond(ctx, msg, output, 0)
}
