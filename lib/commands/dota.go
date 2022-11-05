package commands

import (
	"aigr20/botbot/lib/util"

	"github.com/bwmarrin/discordgo"
)

func DotaCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	util.Reply(s, i.Interaction, "To be implemented")
}
