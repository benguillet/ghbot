package greenhouse

import (
	"encoding/json"
	"fmt"
)

type Application struct {
	ID int `json:"id"`
}

func (c *Client) GetApplications() (*[]Application, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/applications", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	// TODO: Need pagination support
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var a []Application
	err2 := decoder.Decode(&a)
	if err2 != nil {
		fmt.Printf("%+v\n", err2)
	}

	return &a, nil
}