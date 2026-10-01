package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/platform"
)

type Bot struct {
	session      *discordgo.Session
	channelID    string
	orchestrator *platform.EventOrchestrator
}

func NewBot(token string, channelID string, orchestrator *platform.EventOrchestrator) (*Bot, error) {
	if token == "" || channelID == "" {
		return nil, fmt.Errorf("Discord token and channelID are required")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	bot := &Bot{
		session:      dg,
		channelID:    channelID,
		orchestrator: orchestrator,
	}

	dg.AddHandler(bot.handleInteraction)

	if err := dg.Open(); err != nil {
		return nil, fmt.Errorf("error opening websocket connection with Discord: %w", err)
	}

	log.Println("🤖 Approval Discord Bot started and listening for interactions successfully")
	return bot, nil
}

func (b *Bot) Close() {
	if b.session != nil {
		_ = b.session.Close()
	}
}

func (b *Bot) SendApprovalNotification(post *approval.ApprovalPost) error {
	if post == nil {
		return fmt.Errorf("post cannot be nil")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🤖 New LinkedIn Draft Post",
		Description: post.GeneratedPost.Content,
		Color:       0x0077B5, // LinkedIn Official Blue
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Draft UUID",
				Value:  post.UUID,
				Inline: true,
			},
			{
				Name:   "Archetype",
				Value:  string(post.GeneratedPost.Archetype),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  "🟡 PENDING (Awaiting your decision)",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Git-To-Feed • Click a button to process",
		},
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "🚀 Approve & Publish",
				Style:    discordgo.SuccessButton,
				CustomID: "approve_" + post.UUID,
			},
			discordgo.Button{
				Label:    "❌ Reject",
				Style:    discordgo.DangerButton,
				CustomID: "reject_" + post.UUID,
			},
		},
	}

	_, err := b.session.ChannelMessageSendComplex(b.channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{actions},
	})

	if err != nil {
		return fmt.Errorf("error sending approval message to Discord: %w", err)
	}

	return nil
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	ctx := context.Background()

	if len(customID) > 8 && customID[:8] == "approve_" {
		uuid := customID[8:]

		// Respond immediately to Discord to avoid 3-second timeout
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		pubRecord, err := b.orchestrator.ApproveAndPublish(ctx, uuid)
		if err != nil {
			b.updateMessageStatus(s, i, "🔴 PUBLICATION FAILED", fmt.Sprintf("Error: %v", err), 0xFF0000)
			return
		}

		b.updateMessageStatus(s, i, "✅ PUBLISHED ON LINKEDIN LIVE", fmt.Sprintf("Share URN: %s", pubRecord.LinkedInShareURN), 0x00FF00)
		return
	}

	if len(customID) > 7 && customID[:7] == "reject_" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		b.updateMessageStatus(s, i, "❌ DRAFT REJECTED", "The draft has been discarded by the administrator.", 0x808080)
		return
	}
}

func (b *Bot) updateMessageStatus(s *discordgo.Session, i *discordgo.InteractionCreate, statusTitle string, detail string, color int) {
	if i.Message == nil || len(i.Message.Embeds) == 0 {
		return
	}

	origEmbed := i.Message.Embeds[0]
	origEmbed.Color = color
	origEmbed.Fields[2] = &discordgo.MessageEmbedField{
		Name:   "Final Status",
		Value:  fmt.Sprintf("%s\n%s", statusTitle, detail),
		Inline: false,
	}

	embeds := []*discordgo.MessageEmbed{origEmbed}
	emptyComponents := []discordgo.MessageComponent{}

	// Disable buttons after decision
	_, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
		Embeds:     &embeds,
		Components: &emptyComponents, // No buttons
	})
}