package greenhouse

import (
	"encoding/json"
	"fmt"
)

type Candidate struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *Candidate) Name() string {
	return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
}

func (c *Client) GetCandidate(candidateID int) (*Candidate, error) {
	url := fmt.Sprintf("%s/candidates/%d", c.baseURL, candidateID)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var candidate *Candidate
	err = decoder.Decode(&candidate)
	if err != nil {
		return nil, err
	}

	return candidate, nil
}
