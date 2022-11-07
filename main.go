package main

import (
	cmds "aigr20/botbot/lib/commands"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Ping command",
		},
		{
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
			},
		},
		{
			Name:        "dota",
			Description: "Dota command",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Pong!",
				},
			})
		},
		"quote": cmds.QuoteCmd,
		"dota":  cmds.DotaCmd,
	}

	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"next-page": cmds.QuoteCmd,
	}

	token              []byte
	guildID            = flag.String("guild", "", "Guild ID for testing. If empty register commands globally")
	deregisterCommands = flag.Bool("clear", false, "Set to true to remove commands on shutdown.")
	completeDereg      = flag.Bool("clear-full", false, "Remove all commands from all servers")
)

func init() {
	var err error
	token, err = ioutil.ReadFile("token")
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
			if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				go handler(s, i)
			}

		case discordgo.InteractionMessageComponent:
			if handler, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
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
			fmt.Println(guild.Name)
			commands, err := bot.ApplicationCommands(bot.State.User.ID, guild.ID)
			if err != nil {
				log.Printf("Failed to get commands for %s (%s): %s\n", guild.Name, guild.ID, err.Error())
				continue
			}
			for _, command := range commands {
				fmt.Println(command.Name)
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
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, *guildID, v)
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
