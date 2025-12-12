package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/jcoughlin11/pokedexcli/internal/pokecache"
)

const baseUrl = "https://pokeapi.co/api/v2"

// When calling GET on the location-area api WITHOUT an id,
// this is the response that you get. Check with:
// curl https://pokeapi.co/api/v2/location-area/ | jq .
type LocationResponse struct {
	Count int `json:"count"`
	// The difference between using *string and string is
	// using *string you can differentiate between an empty
	// and a missing field. It will be nil if missing and
	// *p == "" if empty. If using string, both empty and
	// missing will have "" as the value, so you can't
	// differentiate between them
	Next    *string `json:"next"`
	Prev    *string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
}

type Client struct {
	cache  pokecache.Cache
	client http.Client
}

func NewClient(reapInterval time.Duration) Client {
	return Client{cache: pokecache.NewCache(reapInterval), client: http.Client{}}
}

func (c *Client) ListLocations(pageUrl *string) (LocationResponse, error) {
	// Default to the base location-area endpoint
	url := baseUrl + "/location-area"
	if pageUrl != nil {
		url = *pageUrl
	}

	// Check the cache first
	if rawData, found := c.cache.Get(&url); found {
		response := LocationResponse{}
		err := json.Unmarshal(rawData, &response)
		if err != nil {
			return LocationResponse{}, err
		}
	}

	// Make network request if not found in cache
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationResponse{}, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return LocationResponse{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return LocationResponse{}, err
	}

	locationsResp := LocationResponse{}
	err = json.Unmarshal(data, &locationsResp)
	if err != nil {
		return LocationResponse{}, err
	}

	// Add data to cache
	c.cache.Add(url, data)

	return locationsResp, nil
}
