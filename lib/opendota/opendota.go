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
var ErrNoMatches = errors.New("no matches played")

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

func GetWinrateAs(accountId int, heroId int) (WinRateResponse, error) {
	response, err := http.Get(fmt.Sprintf("%s/players/%v/wl?hero_id=%v", api_base, accountId, heroId))
	if err != nil {
		log.Println(err)
		return WinRateResponse{}, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	var result WinRateResponse
	err = decoder.Decode(&result)
	if err != nil {
		log.Println(err)
		return WinRateResponse{}, err
	}

	return result, nil
}

func GetLatestMatchAs(accountId int, heroId int) (Match, error) {
	response, err := http.Get(fmt.Sprintf("%s/players/%v/matches?hero_id=%v&limit=1&sort=start_time", api_base, accountId, heroId))
	if err != nil {
		log.Println(err)
		return Match{}, ErrInternalError
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	var result []Match
	err = decoder.Decode(&result)
	if err != nil {
		log.Println(err)
		return Match{}, ErrInternalError
	}

	if len(result) == 0 {
		return Match{}, ErrNoMatches
	}
	return result[0], nil
}
