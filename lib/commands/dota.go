package commands

import (
	"aigr20/botbot/lib/util"
	"encoding/csv"
	"errors"
	"log"
	"os"

	"strconv"

	"github.com/bwmarrin/discordgo"
)

const (
	steam64IdLimit = 76561197960265728
)

var DotaCommandSpecification = &discordgo.ApplicationCommand{
	Name:        "dota",
	Description: "Dota command",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "connect",
			Description: "The steam account id to connect to your discord account",
			Required:    false,
		},
	},
}

func DotaCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		processDotaCommand(s, i.Interaction)
	}
}

func processDotaCommand(s *discordgo.Session, i *discordgo.Interaction) {
	options := i.ApplicationCommandData().Options
	var output string
	ok := false

optionLoop:
	for _, option := range options {
		switch option.Name {
		case "connect":
			accountId, err := strconv.Atoi(option.StringValue())
			if err != nil {
				output = "Steam ID must be numeric"
				ok = true
				break optionLoop
			}

			if accountId > steam64IdLimit {
				accountId -= steam64IdLimit
			}

			err = registerAccount(i.Member.User.ID, accountId)
			if err != nil {
				output = err.Error()
				ok = true
				break optionLoop
			}

			output = "Steam account connected to discord account"
			ok = true
			break optionLoop
		}
	}

	if output != "" && ok {
		util.Reply(s, i, output)
	} else if !ok {
		util.Reply(s, i, "No recognized option provided")
	}
}

var errRegistrationFailed = errors.New("failed to register account")

func registerAccount(discordAccount string, steamAccount int) error {
	file, err := os.OpenFile("steamconnections.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Println(err)
		return errRegistrationFailed
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{discordAccount, strconv.Itoa(steamAccount)})
	if err != nil {
		log.Println(err)
		return errRegistrationFailed
	}

	return nil
}
