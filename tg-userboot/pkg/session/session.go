package session

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"tg-userboot/pkg/types"
)

// TerminalAuth implements gotd/td/telegram/auth.UserAuthenticator
type TerminalAuth struct {
	PhoneNumber string
}

func (a TerminalAuth) Phone(_ context.Context) (string, error) {
	if a.PhoneNumber != "" {
		return a.PhoneNumber, nil
	}
	fmt.Print("Enter Telegram phone number (with country code, e.g. +84...): ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func (a TerminalAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA Cloud Password: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func (a TerminalAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return nil
}

func (a TerminalAuth) Code(_ context.Context, sentCode *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter Telegram login OTP code: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func (a TerminalAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("signup not supported, please register your account in Telegram app first")
}

// Authenticate runs interactive authentication if session is not authenticated
func Authenticate(ctx context.Context, client *telegram.Client, cfg *types.Config) error {
	flow := auth.NewFlow(
		TerminalAuth{},
		auth.SendCodeOptions{},
	)
	return client.Auth().IfNecessary(ctx, flow)
}
