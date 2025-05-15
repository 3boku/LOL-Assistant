package league

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetMatch(name, tag string) (map[string]interface{}, string, error) {
	// 라이엇 계정 정보 가져오기
	url := fmt.Sprintf("https://asia.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s?api_key=%s", name, tag, os.Getenv("RIOT_GAMES_API_KEY"))
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, "", err
	}
	log.Println(result)

	puuid := result["puuid"].(string)
	// 마지막 게임 아이디 가져오기
	url = fmt.Sprintf("https://asia.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=1&api_key=%s", puuid, os.Getenv("RIOT_GAMES_API_KEY"))
	fmt.Println(url)
	resp, err = http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var mathIds []string
	if err = json.Unmarshal(respBody, &mathIds); err != nil {
		return nil, "", err
	}
	log.Println(mathIds)

	// 게임 정보 가져오기
	url = fmt.Sprintf("https://asia.api.riotgames.com/lol/match/v5/matches/%s?api_key=%s", mathIds[0], os.Getenv("RIOT_GAMES_API_KEY"))
	fmt.Println(url)
	resp, err = http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var match map[string]interface{}
	if err = json.Unmarshal(respBody, &match); err != nil {
		return nil, "", err
	}
	log.Println(match)

	return match, puuid, nil
}
