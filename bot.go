package main

import (
	cmds "aigr20/botbot/lib/commands"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type HandlerFunc = func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Command struct {
	Name          string
	Handler       HandlerFunc
	Specification *discordgo.ApplicationCommand
}

type ComponentHandler struct {
	Name    string
	Handler HandlerFunc
}

type Bot struct {
	ID                 string
	Instance           *discordgo.Session
	Commands           []*Command
	ComponentHandlers  []*ComponentHandler
	RegisteredCommands []*discordgo.ApplicationCommand
}

func CreateBot(token []byte) (*Bot, error) {
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, err
	}
	bot.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s#%s\n", s.State.User.Username, s.State.User.Discriminator)
	})

	commands := []*Command{
		{
			Name:          "ping",
			Specification: cmds.PingCommandSpecification,
			Handler:       cmds.PingCmd,
		},
		{
			Name:          "quote",
			Specification: cmds.QuoteCommandSpecification,
			Handler:       cmds.QuoteCmd,
		},
		{
			Name:          "dota",
			Specification: cmds.DotaCommandSpecification,
			Handler:       cmds.DotaCmd,
		},
		{
			Name:          "nick",
			Specification: cmds.NickCommandSpecification,
			Handler:       cmds.NickCmd,
		},
	}

	componentHandlers := []*ComponentHandler{
		{
			Name:    "next-page",
			Handler: cmds.PingCmd,
		},
		{
			Name:    "prev-page",
			Handler: cmds.PingCmd,
		},
	}

	b := &Bot{
		Instance:           bot,
		Commands:           commands,
		ComponentHandlers:  componentHandlers,
		RegisteredCommands: make([]*discordgo.ApplicationCommand, len(commands)),
	}
	b.Instance.AddHandler(b.CommandHandler)

	return b, nil
}

func (bot *Bot) Run(guildID string) {
	err := bot.Instance.Open()
	if err != nil {
		log.Fatalf("Couldn't open discord connection: %s", err.Error())
	}

	bot.ID = bot.Instance.State.User.ID

	log.Println("Opened bot connection")
	log.Println("Adding commands...")

	for i, command := range bot.Commands {
		cmd, err := bot.Instance.ApplicationCommandCreate(bot.ID, guildID, command.Specification)
		if err != nil {
			log.Printf("Failed to create '%s' command: %s.\n", command.Name, err.Error())
		}
		log.Printf("Registered '%s'\n", command.Name)
		bot.RegisteredCommands[i] = cmd
	}
}

func (bot *Bot) CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		var command *Command
		for _, cmd := range bot.Commands {
			if cmd.Name == i.ApplicationCommandData().Name {
				command = cmd
				break
			}
		}
		if command != nil {
			log.Printf("[command] Handling %s\n", i.ApplicationCommandData().Name)
			go command.Handler(s, i)
		}
	case discordgo.InteractionMessageComponent:
		var componentHandler *ComponentHandler
		for _, handler := range bot.ComponentHandlers {
			if handler.Name == i.MessageComponentData().CustomID {
				componentHandler = handler
				break
			}
		}
		if componentHandler != nil {
			log.Printf("[interaction] Handling %s\n", i.MessageComponentData().CustomID)
			go componentHandler.Handler(s, i)
		}
	}
}

func (bot *Bot) Close() {
	bot.Instance.Close()
}

func (bot *Bot) RemoveLocalCommands(guildID string) {
	log.Println("Unregistering commands...")
	for _, cmd := range bot.RegisteredCommands {
		log.Printf("Unregistering '%s'\n", cmd.Name)
		err := bot.Instance.ApplicationCommandDelete(bot.ID, guildID, cmd.ID)
		if err != nil {
			log.Printf("Failed to unregister '%s': %s\n", cmd.Name, err.Error())
		}
	}
}

func (bot *Bot) RemoveAllCommands() {
	log.Println("Unregistering all commands...")
	for _, guild := range bot.Instance.State.Guilds {
		commands, err := bot.Instance.ApplicationCommands(bot.ID, guild.ID)
		if err != nil {
			log.Printf("Failed to get commands for '%s' (%s)\n", guild.Name, guild.ID)
			continue
		}

		for _, cmd := range commands {
			log.Printf("Unregistering '%s' in '%s'\n", cmd.Name, guild.Name)
			err = bot.Instance.ApplicationCommandDelete(cmd.ApplicationID, cmd.GuildID, cmd.ID)
			if err != nil {
				log.Printf("Failed to delete the command '%s': %s\n", cmd.Name, err)
			}
		}
	}
	log.Println("All commands unregistered")
}
