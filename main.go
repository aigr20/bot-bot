package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
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
	bot, err := CreateBot(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %s\n", err.Error())
	}

	bot.Run(*guildID)
	defer bot.Close()

	if *completeDereg {
		bot.RemoveAllCommands()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	if *deregisterCommands || *completeDereg {
		bot.RemoveLocalCommands(*guildID)
	}

	log.Println("Shutting down gracefully...")
}
