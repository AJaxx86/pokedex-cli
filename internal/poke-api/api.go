package pokeapi

import (
	"io"
	"net/http"
	"encoding/json"
	"fmt"
)

const URL = "https://pokeapi.co/api/v2/"
const defaultLimit = 20

type locationArea struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
}


func GetAreas(url string) ([]string, string, string, error) {
	if url == "" {
		return nil, "", "", fmt.Errorf("nul URL pointer passed")
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, "", "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return nil, "", "", fmt.Errorf("Network error: %v", res.StatusCode)
	}
	if err != nil {
		return nil, "", "", err
	}

	locArea := locationArea{}
	jErr := json.Unmarshal(body, &locArea)
	if jErr != nil {
		return nil, "", "", jErr
	}

	areas := []string{}
	for _, area := range locArea.Results {
		areas = append(areas, area.Name)
	}
	
	return areas, locArea.Next, locArea.Previous, nil
}
