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
		return nil, fmt.Errorf("token y channelID de Discord son requeridos")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creando sesion de Discord: %w", err)
	}

	bot := &Bot{
		session:      dg,
		channelID:    channelID,
		orchestrator: orchestrator,
	}

	dg.AddHandler(bot.handleInteraction)

	if err := dg.Open(); err != nil {
		return nil, fmt.Errorf("error abriendo conexion websocket con Discord: %w", err)
	}

	log.Println("🤖 Discord Bot de Aprobacion iniciado y escuchando interacciones con exito")
	return bot, nil
}

func (b *Bot) Close() {
	if b.session != nil {
		_ = b.session.Close()
	}
}

func (b *Bot) SendApprovalNotification(post *approval.ApprovalPost) error {
	if post == nil {
		return fmt.Errorf("post no puede ser nil")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🤖 Nuevo Borrador de Post para LinkedIn",
		Description: post.GeneratedPost.Content,
		Color:       0x0077B5, // Azul oficial LinkedIn
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "UUID de Borrador",
				Value:  post.UUID,
				Inline: true,
			},
			{
				Name:   "Arquetipo",
				Value:  string(post.GeneratedPost.Archetype),
				Inline: true,
			},
			{
				Name:   "Estado",
				Value:  "🟡 PENDING (Esperando tu decisión)",
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Git-To-Feed • Presiona un botón para procesar",
		},
	}

	actions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "🚀 Aprobar y Publicar",
				Style:    discordgo.SuccessButton,
				CustomID: "approve_" + post.UUID,
			},
			discordgo.Button{
				Label:    "❌ Rechazar",
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
		return fmt.Errorf("error enviando mensaje de aprobacion a Discord: %w", err)
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

		// Responder inmediatamente a Discord para evitar timeout de 3 segundos
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		pubRecord, err := b.orchestrator.ApproveAndPublish(ctx, uuid)
		if err != nil {
			b.updateMessageStatus(s, i, "🔴 FALLÓ LA PUBLICACIÓN", fmt.Sprintf("Error: %v", err), 0xFF0000)
			return
		}

		b.updateMessageStatus(s, i, "✅ PUBLICADO EN LINKEDIN EN VIVO", fmt.Sprintf("Share URN: %s", pubRecord.LinkedInShareURN), 0x00FF00)
		return
	}

	if len(customID) > 7 && customID[:7] == "reject_" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})

		b.updateMessageStatus(s, i, "❌ BORRADOR RECHAZADO", "El borrador ha sido descartado por el administrador.", 0x808080)
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

	// Desactivar botones despues de la decision
	_, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
		Embeds:     &embeds,
		Components: &emptyComponents, // Sin botones
	})
}