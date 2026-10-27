package types

import (
	"context"
	"sync"
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
	ReplyToMsgID          int
	ReplyToUserID         int64
	ReplyToUserAccessHash int64
	ReplyToText           string
	Peer                  tg.InputPeerClass
	RawMessage            tg.MessageClass
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

	entityMu sync.RWMutex
	Users    map[int64]int64 // userID -> accessHash
	Channels map[int64]int64 // channelID -> accessHash
}

// CacheEntities stores access hashes for users and channels
func (b *Bot) CacheEntities(users []tg.UserClass, chats []tg.ChatClass) {
	b.entityMu.Lock()
	defer b.entityMu.Unlock()
	if b.Users == nil {
		b.Users = make(map[int64]int64)
	}
	if b.Channels == nil {
		b.Channels = make(map[int64]int64)
	}
	for _, u := range users {
		if user, ok := u.(*tg.User); ok && user.AccessHash != 0 {
			b.Users[user.ID] = user.AccessHash
		}
	}
	for _, c := range chats {
		if ch, ok := c.(*tg.Channel); ok && ch.AccessHash != 0 {
			b.Channels[ch.ID] = ch.AccessHash
		}
	}
}

// GetUserAccessHash returns the cached access hash for a user
func (b *Bot) GetUserAccessHash(userID int64) int64 {
	b.entityMu.RLock()
	defer b.entityMu.RUnlock()
	if b.Users == nil {
		return 0
	}
	return b.Users[userID]
}

// GetChannelAccessHash returns the cached access hash for a channel
func (b *Bot) GetChannelAccessHash(channelID int64) int64 {
	b.entityMu.RLock()
	defer b.entityMu.RUnlock()
	if b.Channels == nil {
		return 0
	}
	return b.Channels[channelID]
}

// GetInputUser returns an InputUser with cached AccessHash
func (b *Bot) GetInputUser(userID int64) tg.InputUserClass {
	return &tg.InputUser{
		UserID:     userID,
		AccessHash: b.GetUserAccessHash(userID),
	}
}

// GetInputPeerUser returns an InputPeerUser with cached AccessHash
func (b *Bot) GetInputPeerUser(userID int64) tg.InputPeerClass {
	return &tg.InputPeerUser{
		UserID:     userID,
		AccessHash: b.GetUserAccessHash(userID),
	}
}

// GetInputChannel returns an InputChannel with cached AccessHash
func (b *Bot) GetInputChannel(channelID int64) *tg.InputChannel {
	return &tg.InputChannel{
		ChannelID:  channelID,
		AccessHash: b.GetChannelAccessHash(channelID),
	}
}
