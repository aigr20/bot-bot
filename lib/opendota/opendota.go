package opendota

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const api_base = "https://api.opendota.com/api"

func GetTotals(accountId int) (TotalResponse, error) {
	response, err := http.Get(fmt.Sprintf("%s/players/%v/totals?hero_id=1", api_base, accountId))
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
