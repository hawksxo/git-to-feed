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

	return &Config{
		Port:                port,
		GitHubWebhookSecret: githubWebhookSecret,
		LinkedInAccessToken: linkedinAccessToken,
		LinkedInAuthorURN:   linkedinAuthorURN,
		DatabaseURL:         databaseURL,
		DiscordBotToken:     discordBotToken,
		DiscordChannelID:    discordChannelID,
	}
}
