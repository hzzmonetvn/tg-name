package fun

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"tg-userboot/pkg/dispatcher"
	"tg-userboot/pkg/types"
)

// Register registers all fun and amusement commands
func Register(d *dispatcher.Dispatcher) {
	d.Register(&types.Command{
		Name:        "joke",
		Category:    "Fun",
		Description: "Fetch a random joke",
		Usage:       ".joke",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleJoke,
	})

	d.Register(&types.Command{
		Name:        "cowsay",
		Category:    "Fun",
		Description: "Have a cow say your message",
		Usage:       ".cowsay [text]",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleCowsay,
	})

	d.Register(&types.Command{
		Name:        "shrug",
		Category:    "Fun",
		Description: "Send shrug emoticon",
		Usage:       ".shrug",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleShrug,
	})

	d.Register(&types.Command{
		Name:        "decide",
		Category:    "Fun",
		Description: "Randomly make a decision",
		Usage:       ".decide",
		OwnerOnly:   false,
		SudoAllowed: true,
		Handler:     handleDecide,
	})
}

func handleJoke(ctx context.Context, b *types.Bot, msg *types.Message) error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://official-joke-api.appspot.com/random_joke", nil)
	resp, err := client.Do(req)
	if err != nil {
		return b.Respond(ctx, msg, "❌ Could not fetch joke.", 3*time.Second)
	}
	defer resp.Body.Close()

	var joke struct {
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&joke); err != nil {
		return b.Respond(ctx, msg, "❌ Failed to parse joke.", 3*time.Second)
	}

	output := fmt.Sprintf("😂 <b>%s</b>\n\n👉 <i>%s</i>", joke.Setup, joke.Punchline)
	return b.Respond(ctx, msg, output, 0)
}

func handleCowsay(ctx context.Context, b *types.Bot, msg *types.Message) error {
	text := strings.TrimSpace(msg.Input)
	if text == "" {
		text = "Moo!"
	}

	border := strings.Repeat("-", len(text)+2)
	cow := fmt.Sprintf(` < %s >
 %s
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`, text, border)

	return b.Respond(ctx, msg, fmt.Sprintf("<pre>%s</pre>", cow), 0)
}

func handleShrug(ctx context.Context, b *types.Bot, msg *types.Message) error {
	return b.Respond(ctx, msg, `¯\_(ツ)_/¯`, 0)
}

func handleDecide(ctx context.Context, b *types.Bot, msg *types.Message) error {
	answers := []string{"Yes, absolutely!", "No, definitely not.", "Maybe...", "Ask again later.", "Without a doubt."}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	choice := answers[r.Intn(len(answers))]
	return b.Respond(ctx, msg, fmt.Sprintf("🎲 <b>Decision:</b> <i>%s</i>", choice), 0)
}
