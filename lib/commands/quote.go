package commands

import (
	"aigr20/botbot/lib/quotes"
	"aigr20/botbot/lib/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var QuoteCommandSpecification = &discordgo.ApplicationCommand{
	Name:        "quote",
	Description: "Quote command",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "add",
			Description: "The quote to add",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "author",
			Description: "The author of the quote",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "get",
			Description: "Get quote at index",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "list",
			Description: "Send a list of all quotes in your DMs",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "delete",
			Description: "Delete the quote at index",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "edit",
			Description: "Index of quote to edit",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "content",
			Description: "New quote content",
			Required:    false,
		},
	},
}

func QuoteCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		processQuoteCommand(s, i.Interaction)
	case discordgo.InteractionMessageComponent:
		processQuoteComponent(s, i.Interaction)
	}
}

func processQuoteCommand(s *discordgo.Session, i *discordgo.Interaction) {
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
		case "delete":
			output = deleteQuote(s, i, int(option.IntValue()))
			ok = true
			break optionLoop
		case "edit":
			newContent, errContent := util.GetStringOption("content", options)
			newAuthor, errAuthor := util.GetUserOption(s, "author", options)
			if errContent != nil && errAuthor != nil {
				output = "You must provide at least one of the 'content' or 'author' options when editing a quote."
				break optionLoop
			}
			output = editQuote(i, int(option.IntValue()), newContent, newAuthor)
			ok = true
			break optionLoop
		}
	}
	if output != "" {
		util.Reply(s, i, output)
	} else if !ok {
		util.Reply(s, i, "No recognized option provided")
	}
}

func processQuoteComponent(s *discordgo.Session, i *discordgo.Interaction) {
	var offsetMod int
	switch i.MessageComponentData().CustomID {
	case "next-page":
		offsetMod = 10
	case "prev-page":
		offsetMod = -10
	}

	updateList(s, i, offsetMod)
	util.Acknowledge(s, i)
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
	quotes.ListTrackerMap[i.Member.User.ID] = &quotes.ListTracker{
		Server: i.GuildID,
		Offset: 0,
	}
	quoteList, _ := quotes.ListQuotes(i.Member.User.ID, 10)
	messageString := ""
	for _, n := range quoteList {
		messageString += n.Quote + "\n"
	}

	offset := quotes.ListTrackerMap[i.Member.User.ID]
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
						Label:    "Prev. page",
						CustomID: "prev-page",
						Disabled: true,
					},
					discordgo.Button{
						Label:    "Next page",
						CustomID: "next-page",
						Disabled: offset.Offset >= offset.TotalQuotes,
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Error sending message to user with id (%s): %s\n", i.Member.User.ID, err.Error())
		util.Reply(s, i, "encountered an error when attempting to DM you")
	}

	return "Sent list in DMs"
}

func updateList(s *discordgo.Session, i *discordgo.Interaction, offsetMod int) {
	offset := quotes.ListTrackerMap[i.User.ID]
	quotes, _ := quotes.ListQuotes(i.User.ID, offsetMod)
	messageString := ""
	for _, quote := range quotes {
		messageString += quote.Quote + "\n"
	}
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: i.ChannelID,
		ID:      i.Message.ID,
		Content: &messageString,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Prev. page",
						CustomID: "prev-page",
						Disabled: offset.Offset-10 == 0,
					},
					discordgo.Button{
						Label:    "Next page",
						CustomID: "next-page",
						Disabled: offset.Offset >= offset.TotalQuotes,
					},
				},
			},
		},
	})
}

func deleteQuote(s *discordgo.Session, i *discordgo.Interaction, index int) string {
	if index < 1 {
		return "The index must be larger than 0."
	}

	err := quotes.DeleteQuote(index-1, i.GuildID)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Quote at index %v has been deleted.", index)
}

func editQuote(i *discordgo.Interaction, index int, newContent string, newAuthor *discordgo.User) string {
	if index < 1 {
		return "The index must be larger than 0."
	}

	if newContent != "" {
		err := quotes.EditQuote(index-1, newContent, i.GuildID)
		if err != nil {
			return err.Error()
		}
	}
	if newAuthor != nil {
		err := quotes.ChangeAuthor(index-1, newAuthor, i.GuildID)
		if err != nil {
			return err.Error()
		}
	}
	return fmt.Sprintf("Quote at index %v has been updated.", index)
}
