package commands

import (
	"aigr20/botbot/lib/util"

	"github.com/bwmarrin/discordgo"
)

var PingCommandSpecification = &discordgo.ApplicationCommand{

	Name:        "ping",
	Description: "Ping command",
}

func PingCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	util.Reply(s, i.Interaction, "Pong!")
}
