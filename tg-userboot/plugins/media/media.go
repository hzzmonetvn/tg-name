package media

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all media and visual utility commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "makeqr",
		Aliases:     []string{"qr", "genqr"},
		Category:    "Media",
		Description: "Generate a QR code from text or link",
		Usage:       ".makeqr [text or link]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleMakeQR,
	})
}

func handleMakeQR(ctx context.Context, b *types.Bot, msg *types.Message) error {
	data := strings.TrimSpace(msg.Input)
	if data == "" && msg.ReplyToText != "" {
		data = msg.ReplyToText
	}

	if data == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.makeqr [text or URL]</code>", 5*time.Second)
	}

	qrURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=400x400&data=%s", url.QueryEscape(data))
	output := fmt.Sprintf("📱 <b>Generated QR Code</b>\n\n🔗 <b>Payload:</b> <code>%s</code>\n🖼️ <a href=\"%s\">View / Download QR Code</a>",
		data, qrURL)

	return b.Respond(ctx, msg, output, 0)
}
