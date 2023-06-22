package commands

import (
	"aigr20/botbot/lib/util"
	"log"

	"github.com/bwmarrin/discordgo"
)

var NickCommandSpecification = &discordgo.ApplicationCommand{
	Name:        "nick",
	Description: "Set the nickname for a user",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "User to set nickname for",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "nickname",
			Description: "New nickname",
			Required:    true,
		},
	},
}

func NickCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	nickname, _ := util.GetStringOption("nickname", options)
	user, err := util.GetUserOption(s, "user", options)
	if err != nil {
		util.Replyf(s, i.Interaction, "Couldn't find the user that was requested.")
		log.Printf("Error getting user: %s\n", err.Error())
		return
	}

	oldNick := ""
	guildMember, err := s.GuildMember(i.GuildID, user.ID)
	if err == nil {
		oldNick = guildMember.Nick
	}

	err = s.GuildMemberNickname(i.GuildID, user.ID, nickname)
	if err != nil {
		util.ReplyErr(s, i.Interaction, err)
		log.Printf("Couldn't change nickname: %s\n", err.Error())
		return
	}
	if oldNick == "" && nickname != "" {
		util.Replyf(s, i.Interaction, "Set the nickname of %s to %s.", user.Username, nickname)
	} else {
		util.Replyf(s, i.Interaction, "Changed the nickname of %s from %s to %s.", user.Username, oldNick, nickname)
	}
}
