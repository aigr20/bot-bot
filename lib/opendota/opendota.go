package opendota

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

const api_base = "https://api.opendota.com/api"

var ErrInternalError = errors.New("internal error")
var ErrNoHero = errors.New("no hero going by that alias")

func GetHeroFromAlias(alias string) (HeroAlias, error) {
	file, err := os.Open("heroes.json")
	if err != nil {
		log.Println(err)
		return HeroAlias{}, ErrInternalError
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var aliasMap AliasMap
	err = decoder.Decode(&aliasMap)
	if err != nil {
		log.Println(err)
		return HeroAlias{}, ErrInternalError
	}

	for aliasKey, hero := range aliasMap {
		if aliasKey == alias {
			return hero, nil
		}
	}
	return HeroAlias{}, ErrNoHero
}

func GetTotals(accountId int, heroId int) (TotalResponse, error) {
	response, err := http.Get(fmt.Sprintf("%s/players/%v/totals?hero_id=%v", api_base, accountId, heroId))
	if err != nil {
		log.Println(err)
		return TotalResponse{}, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	result := make(TotalResponse, 0)
	err = decoder.Decode(&result)
	if err != nil {
		log.Println(err)
		return TotalResponse{}, err
	}

	return result, nil
}
