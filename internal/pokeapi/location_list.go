package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (Locations, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}
	data, exists := c.cache.GetCache(url)
	if !exists {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return Locations{}, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return Locations{}, err
		}

		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return Locations{}, err
		}
		c.cache.AddCache(url, data)
	}

	locationsResp := Locations{}
	err := json.Unmarshal(data, &locationsResp)
	if err != nil {
		return Locations{}, err
	}

	return locationsResp, nil
}
