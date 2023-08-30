package main

import (
	"aigr20/botbot/lib/opendota"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Hero struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	LocalizedName    string   `json:"localized_name"`
	PrimaryAttribute string   `json:"primary_attr"`
	AttackType       string   `json:"attack_type"`
	Roles            []string `json:"roles"`
}

func getHeroes() []Hero {
	response, err := http.Get("https://api.opendota.com/api/heroes")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	heroes := make([]Hero, 0)
	err = decoder.Decode(&heroes)
	if err != nil {
		log.Fatal(err)
	}

	return heroes
}

func createAliasMap(heroes []Hero) opendota.AliasMap {
	aliasMap := make(opendota.AliasMap)
	scanner := bufio.NewScanner(os.Stdin)
	for _, hero := range heroes {
		var input string
		fmt.Printf("%s:\n", hero.LocalizedName)
		if scanner.Scan() {
			input = scanner.Text()
		} else {
			return aliasMap
		}

		aliases := strings.Split(input, ",")
		for _, alias := range aliases {
			aliasMap[alias] = opendota.HeroAlias{Id: hero.Id, Name: hero.LocalizedName}
		}
	}

	return aliasMap
}

func saveAliases(aliases opendota.AliasMap) {
	file, err := os.OpenFile("heroes.json", os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(aliases)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Saved to heroes.json")
}

func main() {
	heroes := getHeroes()
	aliases := createAliasMap(heroes)
	saveAliases(aliases)
}
