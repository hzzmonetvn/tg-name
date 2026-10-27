package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"tg-userboot/pkg/types"
)

// Load reads configuration from .env and environment variables
func Load() (*types.Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	appID, _ := strconv.Atoi(getEnv("API_ID", "2040"))
	appHash := getEnv("API_HASH", "b18441a1ff607e10a989891a5462e627")
	cmdTrigger := getEnv("CMD_TRIGGER", ".")
	sudoTrigger := getEnv("SUDO_TRIGGER", "!")
	sessionFile := getEnv("SESSION_FILE", "userbot.session")
	sessionString := getEnv("SESSION_STRING", "")
	ownerID, _ := strconv.ParseInt(getEnv("OWNER_ID", "0"), 10, 64)
	logChatID, _ := strconv.ParseInt(getEnv("LOG_CHAT", "0"), 10, 64)
	geminiKey := getEnv("GEMINI_API_KEY", "")
	aiKey := getEnv("AI_API_KEY", "")
	aiBaseURL := getEnv("AI_BASE_URL", "")
	aiModel := getEnv("AI_MODEL", "")
	aiSysPrompt := getEnv("AI_SYSTEM_PROMPT", "You are a helpful and concise Telegram AI assistant.")
	imgKey := getEnv("IMAGE_API_KEY", "")
	imgBaseURL := getEnv("IMAGE_BASE_URL", "")
	imgModel := getEnv("IMAGE_MODEL", "dall-e-3")

	// Backward compatibility and sensible defaults
	if aiKey == "" && geminiKey != "" {
		aiKey = geminiKey
	}
	if imgKey == "" {
		imgKey = aiKey
	}
	if imgBaseURL == "" {
		imgBaseURL = aiBaseURL
	}
	if aiModel == "" {
		if strings.Contains(aiBaseURL, "googleapis.com") || (aiBaseURL == "" && geminiKey != "") {
			aiModel = "gemini-1.5-flash"
		} else {
			aiModel = "gpt-4o-mini"
		}
	}
	if aiBaseURL == "" && geminiKey == "" {
		aiBaseURL = "https://api.openai.com/v1"
	}
	if imgBaseURL == "" {
		imgBaseURL = "https://api.openai.com/v1"
	}

	openWeatherKey := getEnv("WEATHER_API", "")
	defaultCity := getEnv("CITY", "Hanoi")
	customPack := getEnv("CUSTOM_PACK_NAME", "PlainUB")
	devMode := getEnv("DEV_MODE", "0") == "1"

	var sudoUsers []int64
	sudoStr := getEnv("SUDO_USERS", "")
	if sudoStr != "" {
		for _, s := range strings.Split(sudoStr, " ") {
			if id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
				sudoUsers = append(sudoUsers, id)
			}
		}
	}

	return &types.Config{
		AppID:          appID,
		AppHash:        appHash,
		SessionFile:    sessionFile,
		SessionString:  sessionString,
		CmdTrigger:     cmdTrigger,
		SudoTrigger:    sudoTrigger,
		OwnerID:        ownerID,
		SudoUsers:      sudoUsers,
		GeminiAPIKey:   geminiKey,
		AIAPIKey:       aiKey,
		AIBaseURL:      aiBaseURL,
		AIModel:        aiModel,
		AISystemPrompt: aiSysPrompt,
		ImageAPIKey:    imgKey,
		ImageBaseURL:   imgBaseURL,
		ImageModel:     imgModel,
		OpenWeatherKey: openWeatherKey,
		DefaultCity:    defaultCity,
		LogChatID:      logChatID,
		DevMode:        devMode,
		CustomPackName: customPack,
	}, nil
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return strings.TrimSpace(val)
}
