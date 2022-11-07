package commands

import (
	"aigr20/botbot/lib/quotes"
	"aigr20/botbot/lib/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func QuoteCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		processCommand(s, i.Interaction)
	case discordgo.InteractionMessageComponent:
		processComponent(s, i.Interaction)
	}
}

func processCommand(s *discordgo.Session, i *discordgo.Interaction) {
	options := i.ApplicationCommandData().Options
	var output string
	ok := false

optionLoop:
	for _, option := range options {
		switch option.Name {
		case "add":
			output = addQuote(s, i, option.StringValue(), options)
			ok = true
			break optionLoop
		case "list":
			output = listQuotes(s, i)
			ok = true
			break optionLoop
		case "get":
			output = getQuote(s, i, int(option.IntValue()))
			ok = true
			break optionLoop
		case "author":
			break optionLoop
		}
	}
	if output != "" {
		util.Reply(s, i, output)
	} else if !ok {
		util.Reply(s, i, "No recognized option provided")
	}
}

func processComponent(s *discordgo.Session, i *discordgo.Interaction) {
	switch i.MessageComponentData().CustomID {
	case "next-page":
		quotes, _ := quotes.ListQuotes(i.User.ID)
		messageString := ""
		for _, quote := range quotes {
			messageString += quote.Quote + "\n"
		}
		s.ChannelMessageEdit(i.ChannelID, i.Message.ID, messageString)
		util.Acknowledge(s, i)
	}
}

func addQuote(s *discordgo.Session, i *discordgo.Interaction, content string, options util.OptionArray) string {
	author, err := util.GetUserOption(s, "author", options)
	if err != nil {
		author = i.Member.User
	}

	index, err := quotes.AddQuote(content, author, i.GuildID)
	if err != nil {
		util.ReplyErr(s, i, err)
		return ""
	}

	return fmt.Sprintf("Added quote at index %v", index)
}

func getQuote(s *discordgo.Session, i *discordgo.Interaction, index int) string {
	quote, err := quotes.GetQuote(index, i.GuildID)
	if err != nil {
		util.ReplyErr(s, i, err)
		return ""
	}

	author, err := s.User(quote.Author)
	if err != nil {
		log.Printf("Error creating user from retrieved id: %s\n", err.Error())
		util.Reply(s, i, "encountered an error when getting the author of the quote")
	}

	embed := quote.Embed(s)
	if embed != nil {
		util.ReplyEmbed(s, i, embed)
		return ""
	}
	return fmt.Sprintf("\"%s\" - %s", quote.Quote, author.Username)
}

func listQuotes(s *discordgo.Session, i *discordgo.Interaction) string {
	quotes.QuoteListOffsets[i.Member.User.ID] = &quotes.ListOffsetTracker{
		Server: i.GuildID,
		Offset: 0,
	}
	quotes, _ := quotes.ListQuotes(i.Member.User.ID)
	messageString := ""
	for _, n := range quotes {
		messageString += n.Quote + "\n"
	}

	channel, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("Error creating user channel for id (%s): %s\n", i.Member.User.ID, err.Error())
		util.Reply(s, i, "encountered an error when attempting to DM you")
	}
	_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: messageString,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Next page",
						CustomID: "next-page",
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending message to user with id (%s): %s\n", i.Member.User.ID, err.Error())
		util.Reply(s, i, "encountered an error when attempting to DM you")
	}
	util.Acknowledge(s, i)
	return ""
}
