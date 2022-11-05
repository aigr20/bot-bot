package util

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetIntOption(name string, options OptionArray) (int64, error) {
	for _, option := range options {
		if option.Name == name {
			return option.IntValue(), nil
		}
	}
	return -1, fmt.Errorf("no option of type int with name %s found", name)
}

func GetStringOption(name string, options OptionArray) (string, error) {
	for _, option := range options {
		if option.Name == name {
			return option.StringValue(), nil
		}
	}
	return "", fmt.Errorf("no option of type string with name %s found", name)
}

func GetUserOption(session *discordgo.Session, name string, options OptionArray) (*discordgo.User, error) {
	for _, option := range options {
		if option.Name == name {
			return option.UserValue(session), nil
		}
	}
	return nil, fmt.Errorf("no option of type *user* with name **%s** found", name)
}

func GetBoolOption(name string, options OptionArray) (bool, error) {
	for _, option := range options {
		if option.Name == name {
			return option.BoolValue(), nil
		}
	}
	return false, fmt.Errorf("no option of type bool with name %s found", name)
}
