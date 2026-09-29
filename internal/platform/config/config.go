package config

import "os"

type Config struct {
	Port                string
	GitHubWebhookSecret string
	LinkedInAccessToken string
	LinkedInAuthorURN   string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	githubWebhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	linkedinAccessToken := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	linkedinAuthorURN := os.Getenv("LINKEDIN_AUTHOR_URN")

	return &Config{
		Port:                port,
		GitHubWebhookSecret: githubWebhookSecret,
		LinkedInAccessToken: linkedinAccessToken,
		LinkedInAuthorURN:   linkedinAuthorURN,
	}
}
