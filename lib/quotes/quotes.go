package quotes

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Quote struct {
	Author string
	Quote  string
}

func (q *Quote) String() string {
	return fmt.Sprintf("\"%s\",\"%s\"", q.Quote, q.Author)
}

// Returns nil if an error is encountered
func (q *Quote) Embed(s *discordgo.Session) *discordgo.MessageEmbed {
	user, err := s.User(q.Author)
	if err != nil {
		return nil
	}

	return &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Quote",
				Value: q.Quote,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    user.Username,
			IconURL: user.AvatarURL(""),
		},
	}
}

func getFile(server string) (*os.File, error) {
	filename := fmt.Sprintf("quotes_%s.csv", server)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		log.Printf("Error opening %s: %s\n", filename, err.Error())
		return nil, errors.New("encountered an error when opening the quotes file")
	}

	return file, nil
}

func parseQuotes(file *os.File) ([]Quote, error) {
	quotes := make([]Quote, 0)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("encountered error reading csv: %s", err.Error())
	}

	for _, record := range records {
		quotes = append(quotes, Quote{Quote: record[0], Author: record[1]})
	}

	return quotes, nil
}

func getSpecificQuote(file *os.File, index int) (Quote, error) {
	reader := csv.NewReader(file)

	var record []string
	var err error
	for i := 0; i < index; i++ {
		record, err = reader.Read()
		if err != nil {
			log.Printf("Error reading quote record at index %v: %s\n", i, err.Error())
			return Quote{}, errors.New("encountered an error when looking for your quote (likely you provided an index that was too high)")
		}
	}

	return Quote{Quote: record[0], Author: record[1]}, nil
}

func AddQuote(content string, author *discordgo.User, server string) (int, error) {
	file, err := getFile(server)
	if err != nil {
		return -1, err
	}
	defer file.Close()
	quotes, err := parseQuotes(file)
	if err != nil {
		log.Printf("Error parsing quotes: %s\n", err.Error())
		return -1, errors.New("encountered an error when parsing existing quotes")
	}
	quotes = append(quotes, Quote{Quote: content, Author: author.ID})

	var lines string
	for _, quote := range quotes {
		lines += quote.String() + "\n"
	}

	_, err = file.WriteAt([]byte(lines), 0)
	if err != nil {
		log.Printf("Error writing to %s: %s\n", file.Name(), err.Error())
		return -1, errors.New("encountered an error when saving the quote to file")
	}

	return len(quotes), nil
}

func GetQuote(index int, server string) (Quote, error) {
	file, err := getFile(server)
	if err != nil {
		return Quote{}, err
	}
	defer file.Close()

	quote, err := getSpecificQuote(file, index)
	if err != nil {
		return Quote{}, err
	}
	return quote, nil
}
