package pokeApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/WaronLimsakul/pokedexcli/internal/pokecache"
)

type LocationArea struct {
	Name string `json:"name"`
}

type FetchingAreasResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type ExploreAreaResponse struct {
	Name               string `json:"name"`
	PokemonsEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonDetail struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	BaseExperience int `json:"base_experience"`
}

type Config struct {
	Next     string
	Previous string
}

func FetchLocationAreas(url string, cache *pokecache.Cache) (FetchingAreasResponse, error) {
	var jsonData []byte
	cachedJson, ok := cache.Get(url)
	if ok {
		jsonData = cachedJson
	} else {
		res, err := http.Get(url)
		if err != nil {
			return FetchingAreasResponse{}, fmt.Errorf("error fethcing location areas: %v", err)
		}
		defer res.Body.Close()
		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return FetchingAreasResponse{}, fmt.Errorf("error reading http response body: %v", err)
		}
		cache.Add(url, jsonData)
	}
	var locationAreasResponse FetchingAreasResponse
	if err := json.Unmarshal(jsonData, &locationAreasResponse); err != nil {
		return FetchingAreasResponse{}, fmt.Errorf("error unmarshalling bytes data: %v", err)
	}
	return locationAreasResponse, nil
}

// I think the pokemon names are a lot, so return pointer sounds better.
func ExploreAreaPokemons(location string, cache *pokecache.Cache) (*ExploreAreaResponse, error) {
	fullURL := "https://pokeapi.co/api/v2/location-area/" + location + "/"
	var jsonData []byte
	cachedJson, ok := cache.Get(fullURL)
	if ok {
		jsonData = cachedJson
	} else {
		res, err := http.Get(fullURL)
		if err != nil {
			return &ExploreAreaResponse{}, fmt.Errorf("error fethcing area information")
		}

		if res.StatusCode != http.StatusOK {
			return &ExploreAreaResponse{}, fmt.Errorf("Area not found")
		}

		defer res.Body.Close()
		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return &ExploreAreaResponse{}, fmt.Errorf("error parsing area information")
		}
	}
	cache.Add(fullURL, jsonData)

	var response ExploreAreaResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return &ExploreAreaResponse{}, fmt.Errorf("error unmarshalling area information")
	}

	return &response, nil
}

func GetPokemonDetail(pokemon string, cache *pokecache.Cache) (*PokemonDetail, error) {
	fullURL := "https://pokeapi.co/api/v2/pokemon/" + pokemon + "/"
	var jsonData []byte
	cacheJson, ok := cache.Get(fullURL)
	if ok {
		jsonData = cacheJson
	} else {
		res, err := http.Get(fullURL)
		if err != nil {
			return &PokemonDetail{}, fmt.Errorf("error fethcing pokemon information")
		}

		if res.StatusCode != http.StatusOK {
			return &PokemonDetail{}, fmt.Errorf("Pokemon not found")
		}

		defer res.Body.Close()
		jsonData, err = io.ReadAll(res.Body)
		if err != nil {
			return &PokemonDetail{}, fmt.Errorf("error parsing pokemon information")
		}
	}

	var pokemonInfo PokemonDetail
	if err := json.Unmarshal(jsonData, &pokemonInfo); err != nil {
		return &PokemonDetail{}, fmt.Errorf("error unmarshalling pokemon detail")
	}
	return &pokemonInfo, nil
}
