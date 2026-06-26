package pokeapi

import (
	"io"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/ajaxx86/pokedex-cli/internal/pokecache"
)

type Client struct {
	cache *pokecache.Cache
	http http.Client
}
type Pokemon struct {
	ID int `json:"id"`
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
type locationArea struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
}
type locationAreaEncounters struct {
	AreaName string `json:"name"`
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL string `json:"url"`
		}`json:"pokemon"`
	}`json:"pokemon_encounters"`
}


func NewClient() Client {
	cl := Client{
		cache: pokecache.NewCache(pokecache.ReapCacheTime),
		http: http.Client{},
	}
	return cl
}


func (cl *Client) GetAreas(url string) ([]string, string, string, error) {
	if url == "" {
		return nil, "", "", fmt.Errorf("nul URL pointer passed")
	}

	entry, found := cl.cache.Get(url)
	if !found {
		res, err := cl.http.Get(url)
		if err != nil {
			return nil, "", "", err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil, "", "", fmt.Errorf("Network error: %v", res.StatusCode)
		}

		entry, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, "", "", err
		}
	}

	locArea := locationArea{}
	jErr := json.Unmarshal(entry, &locArea)
	if jErr != nil {
		return nil, "", "", jErr
	}
	cl.cache.Add(url, entry)

	areas := []string{}
	for _, area := range locArea.Results {
		areas = append(areas, area.Name)
	}

	return areas, locArea.Next, locArea.Previous, nil
}


func (cl *Client) GetEncounters(areaName string) ([]string, error) {
	fullURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", areaName)
	entry, found := cl.cache.Get(fullURL)
	if !found {
		res, err := cl.http.Get(fullURL)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil, fmt.Errorf("Network error: %v", res.StatusCode)
		}

		entry, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}

	areaEncounters := locationAreaEncounters{}
	jErr := json.Unmarshal(entry, &areaEncounters)
	if jErr != nil {
		return nil, jErr
	}
	cl.cache.Add(fullURL, entry)

	result := []string{}
	for _, encounter := range areaEncounters.Encounters {
		result = append(result, encounter.Pokemon.Name)
	}

	return result, nil
}


func (cl *Client) GetPokemonStats(name string) (Pokemon, error) {
	fullURL := "https://pokeapi.co/api/v2/pokemon/" + name
	entry, found := cl.cache.Get(fullURL)
	if !found {
		res, err := cl.http.Get(fullURL)
		if err != nil {
			return Pokemon{}, err
		}
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return Pokemon{}, fmt.Errorf("Network error: %v", res.StatusCode)
		}

		entry, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}
	}

	result := Pokemon{}
	jErr := json.Unmarshal(entry, &result)
	if jErr != nil {
		return Pokemon{}, jErr
	}

	return result, nil
}
