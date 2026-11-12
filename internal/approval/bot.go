package approval

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Approver interface {
	ApproveAndPublish(ctx context.Context, uuid string) (string, error)
	DiscardDraft(ctx context.Context, uuid, reason string) error
}

type Bot struct {
	session  *discordgo.Session
	channels ChannelConfig
	approver Approver
}

func validateVariables(token string, channels ChannelConfig) error {
	if token == "" {
		return fmt.Errorf("se requiere el token de discord")
	}

	if channels.DefaultChannelID == "" {
		return fmt.Errorf("se requiere el identificador del canal por defecto")
	}

	return nil
}

func (b *Bot) Close() {
	if b.session != nil {
		_ = b.session.Close()
	}
}

func (b *Bot) IsConnected() bool {
	return b.session != nil && b.session.State != nil
}

func NewBot(token string, channels ChannelConfig, approver Approver) (*Bot, error) {
	err := validateVariables(token, channels)
	if err != nil {
		return nil, err
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error al intentar crear una sesión de Discord: %w", err)
	}

	bot := &Bot{
		session:  dg,
		channels: channels,
		approver: approver,
	}

	dg.AddHandler(bot.handleInteraction)
	if err := dg.Open(); err != nil {
		return nil, fmt.Errorf("error al abrir la conexión WebSocket con Discord: %w", err)
	}

	return bot, nil
}

func (b *Bot) resolveChannel(targetChannel string, channelType string) string {
	if targetChannel != "" {
		return targetChannel
	}
	return b.channels.DefaultChannelID
}

func (b *Bot) SendApprovalNotification(post *ApprovalPost) error {
	if post == nil {
		return fmt.Errorf("el post no puede ser nil")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Nuevo Borrador para LinkedIn",
		Description: post.Content,
		Color:       0x0077B5,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "UUID del Borrador",
				Value:  post.UUID,
				Inline: true,
			},
			{
				Name:   "Arqueotipo",
				Value:  string(post.Archetype),
				Inline: true,
			},
			{
				Name:   "Estado",
				Value:  "🟡 PENDIENTE (Esperando tu confirmación)",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Git To Feed • Pulsa un botón para continuar",
		},
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Publicar",
				Style:    discordgo.SuccessButton,
				CustomID: "approve_" + post.UUID,
			},
			discordgo.Button{
				Label:    "Descartar",
				Style:    discordgo.DangerButton,
				CustomID: "discard_" + post.UUID,
			},
		},
	}

	var targetChannel string
	if post.Archetype == "release" {
		targetChannel = b.resolveChannel(b.channels.ReleasesChannelID, "releases")
	} else {
		targetChannel = b.resolveChannel(b.channels.FeaturesChannelID, "features")
	}

	_, err := b.session.ChannelMessageSendComplex(targetChannel, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{actions},
	})

	if err != nil {
		return fmt.Errorf("error al intentar enviar el mesaje a Discord: %w", err)
	}

	return nil
}

func (b *Bot) SendSupersededNotification(post *ApprovalPost, releaseTag string) error {
	targetChannel := b.resolveChannel(b.channels.FeaturesChannelID, "features")
	embed := &discordgo.MessageEmbed{
		Title:       "⚡ Borrador sustituido por Release",
		Description: fmt.Sprintf("El borrador `%s` ha sido reemplazado por la Release `%s`. Revise el borrador nuevo vs el original antes de su divulgación.", post.UUID, releaseTag),
		Color:       0xFEE75C,
	}
	_, err := b.session.ChannelMessageSendEmbed(targetChannel, embed)
	return err
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	ctx := context.Background()

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	if len(customID) > 8 && customID[:8] == "approve_" {
		uuid := customID[8:]

		pubRecord, err := b.approver.ApproveAndPublish(ctx, uuid)
		if err != nil {
			b.updateMessageStatus(s, i, "🔴 FALLO EN PUBLICACIÓN", fmt.Sprintf("Error: %v", err), 0xFF0000)
			return
		}

		auditChannel := b.resolveChannel(b.channels.AuditChannelID, "audit")
		auditEmbed := &discordgo.MessageEmbed{
			Title:       "Disponible en LinkedIn",
			Description: fmt.Sprintf("Publicación aprobada y publicada correctamente.\n\n**Share URN:** `%s`", pubRecord),
			Color:       0x0077B5,
		}
		_, _ = s.ChannelMessageSendEmbed(auditChannel, auditEmbed)
		return
	}

	if len(customID) > 8 && customID[:8] == "discard_" {
		uuid := customID[8:]
		err := b.approver.DiscardDraft(ctx, uuid, "manual_discard")
		if err != nil {
			b.updateMessageStatus(s, i, "🔴 FALLO AL DESCARTAR", fmt.Sprintf("Error: %v", err), 0xFF0000)
			return
		}
		b.updateMessageStatus(s, i, "🔴 BORRADOR DESCARTADO", "El borrador fue descartado por el usuario.", 0x808080)
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
		Name:   "Estado Final",
		Value:  fmt.Sprintf("%s\n%s", statusTitle, detail),
		Inline: false,
	}

	embeds := []*discordgo.MessageEmbed{origEmbed}
	emptyComponents := []discordgo.MessageComponent{}

	_, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
		Embeds:     &embeds,
		Components: &emptyComponents,
	})
}
