package discord

import (
	"testing"
)

func TestResolveChannel(t *testing.T) {
	bot := &Bot{
		channels: ChannelConfig{
			DefaultChannelID: "default-channel-456",
		},
	}

	t.Run("should fallback to default channel when target channel is empty", func(t *testing.T) {
		result := bot.resolveChannel("", "releases")
		expected := "default-channel-456"

		if result != expected {
			t.Errorf("expected channel %q, got %q", expected, result)
		}
	})

	t.Run("should return configured channel when target is present", func(t *testing.T) {
		result := bot.resolveChannel("features-channel-789", "features")
		expected := "features-channel-789"
		if result != expected {
			t.Errorf("expected channel %q, got %q", expected, result)
		}
	})
}

func TestNewBot_Validation(t *testing.T) {
	t.Run("should return error when token is empty", func(t *testing.T) {
		channels := ChannelConfig{
			DefaultChannelID: "some-channel-id",
		}

		bot, err := NewBot("", channels, nil)
		if err == nil {
			t.Errorf("expected error when token is empty, got nil")
		}
		if bot != nil {
			t.Errorf("expected nil bot when creation fails, got %v", bot)
		}
	})

	t.Run("should return error when DefaultChannelID is empty", func(t *testing.T) {
		channels := ChannelConfig{
			DefaultChannelID: "",
		}

		bot, err := NewBot("dummy-token", channels, nil)
		if err == nil {
			t.Errorf("expected error when DefaultChannelID is empty, got nil")
		}
		if bot != nil {
			t.Errorf("expected nil bot when creation fails, got %v", bot)
		}
	})
}
