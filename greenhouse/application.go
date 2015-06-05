package greenhouse

import (
	"encoding/json"
	"fmt"
)

type Job struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Application struct {
	ID          int    `json:"id"`
	CandidateID int    `json:"candidate_id"`
	Prospect    bool   `json:"prospect"`
	Status      string `json:"status"`
	Jobs        []*Job `json:"jobs"`
}

func (c *Client) GetApplications() ([]*Application, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/applications?last_activity_after=2015-06-01T00:00:00Z&per_page=500", c.baseURL), nil)
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

	var a []*Application
	err = decoder.Decode(&a)
	if err != nil {
		return nil, err
	}

	return a, nil
}
