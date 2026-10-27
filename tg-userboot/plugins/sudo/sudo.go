package sudo

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers sudo management and help commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "help",
		Aliases:     []string{"cmds"},
		Category:    "Sudo",
		Description: "List all loaded userbot commands and categories",
		Usage:       ".help [optional command_name]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleHelp,
	})

	d.Register(&types.Command{
		Name:        "addsudo",
		Category:    "Sudo",
		Description: "Add user ID to active sudo list",
		Usage:       ".addsudo [user_id or reply]",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleAddSudo,
	})

	d.Register(&types.Command{
		Name:        "delsudo",
		Category:    "Sudo",
		Description: "Remove user ID from active sudo list",
		Usage:       ".delsudo [user_id or reply]",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleDelSudo,
	})

	d.Register(&types.Command{
		Name:        "sudolist",
		Category:    "Sudo",
		Description: "View list of authorized sudo users",
		Usage:       ".sudolist",
		OwnerOnly:   true,
		SudoAllowed: false,
		Handler:     handleSudoList,
	})
}

func handleHelp(ctx context.Context, b *types.Bot, msg *types.Message) error {
	arg := strings.ToLower(strings.TrimSpace(msg.Input))

	if arg != "" {
		cmd, exists := b.Commands[arg]
		if !exists {
			return b.Respond(ctx, msg, fmt.Sprintf("❌ Command <code>%s</code> not found.", arg), 3*time.Second)
		}

		aliases := "None"
		if len(cmd.Aliases) > 0 {
			aliases = strings.Join(cmd.Aliases, ", ")
		}

		text := fmt.Sprintf(`📖 <b>Command:</b> <code>%s</code>
━━━━━━━━━━━━━━━━━
📂 <b>Category:</b> %s
📝 <b>Description:</b> %s
💡 <b>Usage:</b> <code>%s</code>
🔄 <b>Aliases:</b> <code>%s</code>
👑 <b>Owner Only:</b> <code>%t</code>
━━━━━━━━━━━━━━━━━`,
			cmd.Name,
			cmd.Category,
			cmd.Description,
			cmd.Usage,
			aliases,
			cmd.OwnerOnly,
		)
		return b.Respond(ctx, msg, text, 0)
	}

	// Group commands by category
	categories := make(map[string][]string)
	for name, cmd := range b.Commands {
		cat := cmd.Category
		if cat == "" {
			cat = "General"
		}
		categories[cat] = append(categories[cat], name)
	}

	var catNames []string
	for cat := range categories {
		catNames = append(catNames, cat)
		sort.Strings(categories[cat])
	}
	sort.Strings(catNames)

	var sb strings.Builder
	sb.WriteString("📚 <b>TG-Userboot Commands Directory</b>\n━━━━━━━━━━━━━━━━━\n")
	for _, cat := range catNames {
		sb.WriteString(fmt.Sprintf("📂 <b>%s:</b>\n<code>%s</code>\n\n", cat, strings.Join(categories[cat], ", ")))
	}
	sb.WriteString(fmt.Sprintf("💡 <i>Type <code>%shelp [cmd]</code> for command syntax and details.</i>", b.Config.CmdTrigger))

	return b.Respond(ctx, msg, sb.String(), 0)
}

func handleAddSudo(ctx context.Context, b *types.Bot, msg *types.Message) error {
	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, "❌ Reply to a user or specify user ID to add to sudo.", 5*time.Second)
	}

	for _, id := range b.Config.SudoUsers {
		if id == targetID {
			return b.Respond(ctx, msg, fmt.Sprintf("⚠️ User <code>%d</code> is already a sudoer.", targetID), 3*time.Second)
		}
	}

	b.Config.SudoUsers = append(b.Config.SudoUsers, targetID)
	return b.Respond(ctx, msg, fmt.Sprintf("✅ <b>Added Sudo User:</b> <code>%d</code>", targetID), 5*time.Second)
}

func handleDelSudo(ctx context.Context, b *types.Bot, msg *types.Message) error {
	targetID, _, err := b.ExtractUserAndReason(msg)
	if err != nil {
		return b.Respond(ctx, msg, "❌ Reply to a user or specify user ID to remove from sudo.", 5*time.Second)
	}

	var updated []int64
	found := false
	for _, id := range b.Config.SudoUsers {
		if id == targetID {
			found = true
			continue
		}
		updated = append(updated, id)
	}

	if !found {
		return b.Respond(ctx, msg, fmt.Sprintf("⚠️ User <code>%d</code> was not in the sudoers list.", targetID), 3*time.Second)
	}

	b.Config.SudoUsers = updated
	return b.Respond(ctx, msg, fmt.Sprintf("🗑️ <b>Removed Sudo User:</b> <code>%d</code>", targetID), 5*time.Second)
}

func handleSudoList(ctx context.Context, b *types.Bot, msg *types.Message) error {
	if len(b.Config.SudoUsers) == 0 {
		return b.Respond(ctx, msg, "ℹ️ No sudo users registered.", 5*time.Second)
	}

	var lines []string
	for idx, id := range b.Config.SudoUsers {
		lines = append(lines, fmt.Sprintf("%d. <code>%d</code>", idx+1, id))
	}

	output := fmt.Sprintf("👑 <b>Registered Sudo Users:</b>\n━━━━━━━━━━━━━━━━━\n%s", strings.Join(lines, "\n"))
	return b.Respond(ctx, msg, output, 0)
}
