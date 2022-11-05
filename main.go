package main

import (
	cmds "aigr20/botbot/commands"
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

	token   []byte
	guildID = flag.String("guild", "", "Guild ID for testing. If empty register commands globally")
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
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	err = bot.Open()
	if err != nil {
		log.Fatalf("Couldn't open discord connection: %s", err.Error())
	}

	defer bot.Close()

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

	log.Println("Shutting down gracefully...")
}
