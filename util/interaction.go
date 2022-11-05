package util

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Reply(s *discordgo.Session, i *discordgo.Interaction, content string) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func Replyf(s *discordgo.Session, i *discordgo.Interaction, format string, a ...any) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(format, a...),
		},
	})
}

func ReplyErr(s *discordgo.Session, i *discordgo.Interaction, err error) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Encountered error: %s", err.Error()),
		},
	})
}

func ReplyEmbed(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
