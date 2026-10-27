package android

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers Android enthusiast commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "magisk",
		Aliases:     []string{"root"},
		Category:    "Android",
		Description: "Get official Magisk & KernelSU releases and documentation links",
		Usage:       ".magisk",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleMagisk,
	})

	d.Register(&types.Command{
		Name:        "specs",
		Category:    "Android",
		Description: "Search GSM Arena device specs",
		Usage:       ".specs [device name]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleSpecs,
	})
}

func handleMagisk(ctx context.Context, b *types.Bot, msg *types.Message) error {
	output := `📱 <b>Android Root & Custom ROM Resources</b>
━━━━━━━━━━━━━━━━━
🔥 <b>Magisk:</b> <a href="https://github.com/topjohnwu/Magisk/releases">GitHub Releases</a>
⚡ <b>KernelSU:</b> <a href="https://github.com/tiann/KernelSU/releases">GitHub Releases</a>
🛡️ <b>APatch:</b> <a href="https://github.com/bmax121/APatch/releases">GitHub Releases</a>
📦 <b>XDA Forums:</b> <a href="https://xdaforums.com">Visit XDA</a>
━━━━━━━━━━━━━━━━━`
	return b.Respond(ctx, msg, output, 0)
}

func handleSpecs(ctx context.Context, b *types.Bot, msg *types.Message) error {
	device := strings.TrimSpace(msg.Input)
	if device == "" {
		return b.Respond(ctx, msg, "<b>Usage:</b> <code>.specs [device name]</code>", 3*time.Second)
	}

	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s+gsmarena", strings.ReplaceAll(device, " ", "+"))
	output := fmt.Sprintf("📱 <b>Device Specifications:</b> <a href=\"%s\">Search %s on GSMArena</a>", searchURL, device)
	return b.Respond(ctx, msg, output, 0)
}
