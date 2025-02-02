package pokeApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
}

type FetchingResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

func FetchLocationAreas(url string) (FetchingResponse, error) {

	res, err := http.Get(url)
	if err != nil {
		return FetchingResponse{}, fmt.Errorf("Error fethcing location areas: %v", err)
	}

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return FetchingResponse{}, fmt.Errorf("Error reading http response body: %v", err)
	}
	var locationAreasResponse FetchingResponse
	if err = json.Unmarshal(jsonData, &locationAreasResponse); err != nil {
		return FetchingResponse{}, fmt.Errorf("Error unmarshalling bytes data: %v", err)
	}
	return locationAreasResponse, nil
}
