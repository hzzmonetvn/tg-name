package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"tg-userboot/pkg/config"
	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/session"
	"tg-userboot/pkg/types"
	"tg-userboot/plugins/admin"
	"tg-userboot/plugins/ai"
	"tg-userboot/plugins/android"
	"tg-userboot/plugins/files"
	"tg-userboot/plugins/fun"
	"tg-userboot/plugins/media"
	"tg-userboot/plugins/misc"
	"tg-userboot/plugins/sudo"
	"tg-userboot/plugins/system"
	"tg-userboot/plugins/tgtools"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("🚀 TG-USERBOOT: NEXT-GEN TELEGRAM USERBOT")
	fmt.Println("Engine: gotd/td (Pure Go MTProto 2.0 Client)")
	fmt.Println("==================================================")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	if cfg.AppID == 0 || cfg.AppHash == "" {
		log.Fatalf("❌ API_ID and API_HASH must be configured in .env or environment")
	}

	bot := &types.Bot{
		Config: cfg,
		Start:  time.Now(),
	}

	disp := dispatcher.NewDispatcher(bot)

	// Register all plugin suites
	tgtools.Register(disp)
	admin.Register(disp)
	ai.Register(disp)
	misc.Register(disp)
	media.Register(disp)
	system.Register(disp)
	fun.Register(disp)
	sudo.Register(disp)
	android.Register(disp)
	files.Register(disp)

	log.Printf("📦 Loaded %d commands across 10 plugin categories", len(bot.Commands))

	// Setup gotd client
	sessionStorage := &telegram.FileSessionStorage{
		Path: cfg.SessionFile,
	}

	client := telegram.NewClient(cfg.AppID, cfg.AppHash, telegram.Options{
		SessionStorage: sessionStorage,
		UpdateHandler: telegram.UpdateHandlerFunc(func(ctx context.Context, u tg.UpdatesClass) error {
			return disp.HandleUpdates(ctx, u)
		}),
	})

	bot.Client = client
	bot.API = client.API()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err = client.Run(ctx, func(ctx context.Context) error {
		// Authenticate if needed
		status, authErr := client.Auth().Status(ctx)
		if authErr != nil || !status.Authorized {
			log.Println("🔑 Session not authenticated, starting interactive login...")
			if err := session.Authenticate(ctx, client, cfg); err != nil {
				return fmt.Errorf("authentication error: %w", err)
			}
		}

		self, err := client.Self(ctx)
		if err != nil {
			return fmt.Errorf("failed to get self info: %w", err)
		}

		bot.Self = self
		if bot.Config.OwnerID == 0 {
			bot.Config.OwnerID = self.ID
		}

		log.Println("==================================================")
		log.Printf("✅ LOGGED IN AS: %s (@%s, ID: %d)", self.FirstName, self.Username, self.ID)
		log.Printf("⚡ Command Trigger: '%s' | Sudo Trigger: '%s'", cfg.CmdTrigger, cfg.SudoTrigger)
		log.Printf("🌐 TG-Userboot is active and listening for updates...")
		log.Println("==================================================")

		<-ctx.Done()
		log.Println("🛑 Shutting down TG-Userboot gracefully...")
		return nil
	})

	if err != nil && ctx.Err() == nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
