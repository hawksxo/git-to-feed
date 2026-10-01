package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	GitHubWebhookSecret string
	LinkedInAccessToken string
	LinkedInAuthorURN   string
	DatabaseURL         string
	DiscordBotToken     string
	DiscordChannelID    string
	GeminiAPIKey        string
	GeminiModelName     string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	githubWebhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	linkedinAccessToken := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	linkedinAuthorURN := os.Getenv("LINKEDIN_AUTHOR_URN")
	databaseURL := os.Getenv("DATABASE_URL")
	discordBotToken := os.Getenv("DISCORD_BOT_TOKEN")
	discordChannelID := os.Getenv("DISCORD_CHANNEL_ID")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	geminiModelName := os.Getenv("GEMINI_MODEL_NAME")
	if geminiModelName == "" {
		geminiModelName = "gemini-3.5-flash-lite"
	}

	return &Config{
		Port:                port,
		GitHubWebhookSecret: githubWebhookSecret,
		LinkedInAccessToken: linkedinAccessToken,
		LinkedInAuthorURN:   linkedinAuthorURN,
		DatabaseURL:         databaseURL,
		DiscordBotToken:     discordBotToken,
		DiscordChannelID:    discordChannelID,
		GeminiAPIKey:        geminiAPIKey,
		GeminiModelName:     geminiModelName,
	}
}
