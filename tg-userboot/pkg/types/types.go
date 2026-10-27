package types

import (
	"context"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// Config holds runtime configuration loaded from environment or .env
type Config struct {
	AppID           int
	AppHash         string
	SessionFile     string
	SessionString   string
	CmdTrigger      string
	SudoTrigger     string
	OwnerID         int64
	SudoUsers       []int64
	GeminiAPIKey    string
	AIAPIKey        string
	AIBaseURL       string
	AIModel         string
	AISystemPrompt  string
	ImageAPIKey     string
	ImageBaseURL    string
	ImageModel      string
	OpenWeatherKey  string
	DefaultCity     string
	LogChatID       int64
	DevMode         bool
	CustomPackName  string
}

// Message represents an incoming or outgoing Telegram message wrapper
type Message struct {
	ID             int
	ChatID         int64
	SenderID       int64
	IsOutgoing     bool
	Text           string
	RawText        string
	Trigger        string
	Cmd            string
	Input          string
	Args           []string
	Flags          map[string]bool
	ReplyToMsgID   int
	ReplyToUserID  int64
	ReplyToText    string
	Peer           tg.InputPeerClass
	RawMessage     tg.MessageClass
}

// HandlerFunc is the signature for plugin command handlers
type HandlerFunc func(ctx context.Context, b *Bot, msg *Message) error

// Command represents a registered userbot command
type Command struct {
	Name        string
	Aliases     []string
	Category    string
	Description string
	Usage       string
	Flags       string
	OwnerOnly   bool
	SudoAllowed bool
	Handler     HandlerFunc
}

// Bot wraps the gotd client, raw API client, and dispatcher
type Bot struct {
	Client   *telegram.Client
	API      *tg.Client
	Config   *Config
	Self     *tg.User
	Commands map[string]*Command
	Start    time.Time
}
