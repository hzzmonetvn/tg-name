package system

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all host system inspection commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "neofetch",
		Aliases:     []string{"sysinfo", "sys"},
		Category:    "System",
		Description: "Display host system information and hardware specs",
		Usage:       ".neofetch",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleNeofetch,
	})
}

func handleNeofetch(ctx context.Context, b *types.Bot, msg *types.Message) error {
	_ = b.Respond(ctx, msg, "🖥️ <i>Inspecting host system...</i>", 0)

	var uname string
	if out, err := exec.Command("uname", "-sr").Output(); err == nil {
		uname = strings.TrimSpace(string(out))
	} else {
		uname = runtime.GOOS
	}

	var uptime string
	if out, err := exec.Command("uptime", "-p").Output(); err == nil {
		uptime = strings.TrimSpace(string(out))
	} else {
		uptime = "N/A"
	}

	var memInfo string
	if out, err := exec.Command("free", "-h").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			memInfo = strings.Join(strings.Fields(lines[1]), " | ")
		}
	}
	if memInfo == "" {
		memInfo = "N/A"
	}

	output := fmt.Sprintf(`🖥️ <b>Host System Specs</b>
━━━━━━━━━━━━━━━━━
🐧 <b>Kernel:</b> <code>%s</code>
⏱️ <b>System Uptime:</b> <code>%s</code>
⚙️ <b>Architecture:</b> <code>%s</code>
🧠 <b>Logical CPUs:</b> <code>%d</code>
💾 <b>Memory (Mem):</b> <code>%s</code>
🐹 <b>Go Runtime:</b> <code>%s</code>
━━━━━━━━━━━━━━━━━`,
		uname,
		uptime,
		runtime.GOARCH,
		runtime.NumCPU(),
		memInfo,
		runtime.Version(),
	)

	return b.Respond(ctx, msg, output, 0)
}
