package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

// GenerateRandomScore generates a random score
// between 10.0000 ~ 99.9999
func GenerateRandomScore() float32 {
	return float32(rand.Intn(899999)+100000) / 10000
}

func GetPhotoScore(gameMasterPhotoUrl string, playerPhotoUrl1 string, playerPhotoUrl2 string) (string, error) {
	// Query Build
	method := "GET"
	baseUrl := "https://image-score.yashikota.com/lpips/url"
	gameMasterQuery := fmt.Sprintf("?gamemaster_url=%s", url.QueryEscape(gameMasterPhotoUrl))
	player1Query := fmt.Sprintf("&player1_url=%s", url.QueryEscape(playerPhotoUrl1))
	player2Query := fmt.Sprintf("&player2_url=%s", url.QueryEscape(playerPhotoUrl2))
	queryParams := fmt.Sprintf("%s%s%s", gameMasterQuery, player1Query, player2Query)
	endpoint := fmt.Sprintf("%s%s", baseUrl, queryParams)

	// Request
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resJson := struct {
		Score float32 `json:"score"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f", resJson.Score), nil
}
