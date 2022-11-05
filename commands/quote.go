package commands

import (
	"aigr20/botbot/lib/quotes"
	"aigr20/botbot/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func QuoteCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var output string
	ok := false

optionLoop:
	for _, option := range options {
		switch option.Name {
		case "add":
			output = addQuote(s, i.Interaction, option.StringValue(), options)
			ok = true
			break optionLoop
		case "list":
			break optionLoop
		case "get":
			output = getQuote(s, i.Interaction, int(option.IntValue()))
			ok = true
			break optionLoop
		case "author":
			break optionLoop
		}
	}
	if output != "" {
		util.Reply(s, i.Interaction, output)
	} else if !ok {
		util.Reply(s, i.Interaction, "No recognized option provided")
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
