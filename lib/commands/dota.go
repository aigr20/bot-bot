package commands

import (
	"aigr20/botbot/lib/opendota"
	"aigr20/botbot/lib/util"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
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
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "stats",
			Description: "The user who's stats you want to look up",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "hero",
			Description: "Hero name or abbreviation",
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
		case "stats":
			discordAccount := option.UserValue(s).ID
			steamAccount, err := getAccountByDiscordId(discordAccount)
			if err != nil && err != errAccountNotFound {
				output = "internal error"
				ok = true
				break optionLoop
			} else if err != nil && err == errAccountNotFound {
				output = err.Error()
				ok = true
				break optionLoop
			}
			output = strconv.Itoa(steamAccount)
			ok = true
			hero, err := util.GetStringOption("hero", options)
			if err != nil {
				output = "Please provide a hero name"
				ok = true
				break optionLoop
			}
			getTotals(steamAccount, hero)
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
var errAccountNotFound = errors.New("no account found")

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

func getAccountByDiscordId(discordId string) (int, error) {
	file, err := os.Open("steamconnections.csv")
	if err != nil {
		return 0, errAccountNotFound
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.ReuseRecord = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if record[0] == discordId {
			return strconv.Atoi(record[1])
		}
	}
	return 0, errAccountNotFound
}

func getTotals(account int, hero string) {
	totals, err := opendota.GetTotals(account)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(totals)
}
