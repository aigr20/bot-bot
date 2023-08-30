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
	var embed *discordgo.MessageEmbed

optionLoop:
	for _, option := range options {
		switch option.Name {
		case "connect":
			accountId, err := strconv.Atoi(option.StringValue())
			if err != nil {
				output = "Steam ID must be numeric"
				break optionLoop
			}

			if accountId > steam64IdLimit {
				accountId -= steam64IdLimit
			}

			err = registerAccount(i.Member.User.ID, accountId)
			if err != nil {
				output = err.Error()
				break optionLoop
			}

			output = "Steam account connected to discord account"
			break optionLoop
		case "stats":
			discordAccount := option.UserValue(s).ID
			steamAccount, err := getAccountByDiscordId(discordAccount)
			if err != nil && err != errAccountNotFound {
				output = "internal error"
				break optionLoop
			} else if err != nil && err == errAccountNotFound {
				output = err.Error()
				break optionLoop
			}

			hero, err := util.GetStringOption("hero", options)
			if err != nil {
				output = "Please provide a hero name"
				break optionLoop
			}
			embed = getTotals(steamAccount, hero)
			break optionLoop
		default:
			output = "No primary action provided"
		}
	}
	if embed != nil {
		util.ReplyEmbed(s, i, embed)
	} else {
		util.Reply(s, i, output)
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

func getTotals(account int, heroAlias string) *discordgo.MessageEmbed {
	hero, err := opendota.GetHeroFromAlias(heroAlias)
	if err != nil {
		return nil
	}

	totals, err := opendota.GetTotals(account, hero.Id)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(totals)
	embed := &discordgo.MessageEmbed{Title: fmt.Sprintf("Totals as %s", hero.Name), Fields: totalFields(totals)}
	return embed
}

func addFieldAvg(fields []*discordgo.MessageEmbedField, totals []opendota.TotalField, field string, title string) []*discordgo.MessageEmbedField {
	if f, ok := opendota.GetTotalField(field, totals); ok {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   title,
			Value:  util.FloatString(f.Sum/float64(f.Matches), 2),
			Inline: true,
		})
	}
	return fields
}

var emptyField = discordgo.MessageEmbedField{Name: "-", Value: "-", Inline: true}

func totalFields(totals []opendota.TotalField) []*discordgo.MessageEmbedField {
	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = addFieldAvg(fields, totals, "kills", "Average Kills")
	fields = addFieldAvg(fields, totals, "deaths", "Average Deaths")
	fields = addFieldAvg(fields, totals, "assists", "Average Assists")
	fields = addFieldAvg(fields, totals, "last_hits", "Average Last Hits")
	fields = addFieldAvg(fields, totals, "denies", "Average Denies")
	fields = append(fields, &emptyField)
	fields = addFieldAvg(fields, totals, "hero_damage", "Average Hero Damage")
	fields = addFieldAvg(fields, totals, "tower_damage", "Average Tower Damage")
	fields = addFieldAvg(fields, totals, "hero_healing", "Average Hero Healing")
	return fields
}
