package main

import (
	cmds "aigr20/botbot/lib/commands"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name          string
	Handler       func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Specification *discordgo.ApplicationCommand
}

var (
	commands = []*Command{
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

	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"next-page": cmds.QuoteCmd,
		"prev-page": cmds.QuoteCmd,
	}

	token              []byte
	guildID            = flag.String("guild", "", "Guild ID for testing. If empty register commands globally")
	deregisterCommands = flag.Bool("clear", false, "Set to true to remove commands on shutdown.")
	completeDereg      = flag.Bool("clear-full", false, "Remove all commands from all servers")
)

func init() {
	log.Printf("Launching with discordgo version %v", discordgo.VERSION)
	var err error
	token, err = os.ReadFile("token")
	if err != nil {
		log.Fatalf("failed to open token file: %s\n", err.Error())
	}
	flag.Parse()
}

func main() {
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatalf("Couldn't create bot: %s", err.Error())
	}

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %s#%s", s.State.User.Username, s.State.User.Discriminator)
	})
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			var command *Command
			for _, cmd := range commands {
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
			if handler, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				log.Printf("[interaction] Handling %s\n", i.MessageComponentData().CustomID)
				go handler(s, i)
			}
		}
	})

	err = bot.Open()
	if err != nil {
		log.Fatalf("Couldn't open discord connection: %s", err.Error())
	}

	defer bot.Close()

	if *completeDereg {
		log.Println("Clearing commands")
		for _, guild := range bot.State.Guilds {
			commands, err := bot.ApplicationCommands(bot.State.User.ID, guild.ID)
			if err != nil {
				log.Printf("Failed to get commands for %s (%s): %s\n", guild.Name, guild.ID, err.Error())
				continue
			}
			for _, command := range commands {
				err = bot.ApplicationCommandDelete(bot.State.User.ID, guild.ID, command.ID)
				if err != nil {
					log.Printf("Failed to delete '%s' command in %s (%s): %s", command.Name, guild.Name, guild.ID, err.Error())
				}
			}
		}
		log.Println("All commands have been cleared")
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, *guildID, v.Specification)
		if err != nil {
			log.Panicf("Cannot create '%s' command: %s", v.Name, err.Error())
		}
		registeredCommands[i] = cmd
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	if *deregisterCommands || *completeDereg {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			log.Printf("Deregistering %s\n", v.Name)
			err := bot.ApplicationCommandDelete(bot.State.User.ID, *guildID, v.ID)
			if err != nil {
				log.Fatalf("failed to delete %s command: %s\n", v.Name, err.Error())
			}
		}
	}

	log.Println("Shutting down gracefully...")
}
