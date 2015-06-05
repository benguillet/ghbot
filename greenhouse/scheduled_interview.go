package greenhouse

import (
	"encoding/json"
	"fmt"
)

type Interviewer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Interview struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Start struct {
	DateTime string `json:"date_time"`
}

type ScheduledInterview struct {
	ID           int            `json:"id"`
	Start        *Start         `json:"start"`
	Status       string         `json:"status"`
	Location     string         `json:"location"`
	Interview    *Interview     `json:"interview"`
	Interviewers []*Interviewer `json:"interviewers"`
}

func (c *Client) GetScheduledInterviews(appID int) ([]*ScheduledInterview, error) {
	url := fmt.Sprintf("%s/applications/%d/scheduled_interviews", c.baseURL, appID)
	req, err := c.NewRequest("GET", url, nil)
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

	var scheduled []*ScheduledInterview
	err = decoder.Decode(&scheduled)
	if err != nil {
		return nil, err
	}

	return scheduled, nil
}
