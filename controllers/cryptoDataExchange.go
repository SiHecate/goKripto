package controllers

import (
	"encoding/json"
	"fmt"
	model "gokripto/Model"
	"io"
	"net/http"
)

func GetExchangeRate(cryptoName string) (model.ExchangeData, error) {

	url := fmt.Sprintf("https://api.swapzone.io/v1/exchange/get-rate?from=%s&to=usdt&amount=1", cryptoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.ExchangeData{}, err
	}
	req.Header.Set("x-api-key", "zigQHUYGX")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.ExchangeData{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return model.ExchangeData{}, fmt.Errorf("HTTP Error: %d", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return model.ExchangeData{}, err
	}

	var exchangeData model.ExchangeData
	if err := json.Unmarshal(responseData, &exchangeData); err != nil {
		return model.ExchangeData{}, err
	}

	return exchangeData, nil
}
