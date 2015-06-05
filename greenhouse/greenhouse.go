package greenhouse

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
)

type Client struct {
	client  *http.Client
	baseUrl string
	token   string
}

func New(baseUrl string, token string) *Client {
	c := http.DefaultClient

	return &Client{
		client:  c,
		baseUrl: baseUrl,
		token:   token,
	}
}

func (c *Client) NewRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	cred := fmt.Sprintf("%s:", c.token)
	basic := fmt.Sprintf(" Basic %s", base64.StdEncoding.EncodeToString([]byte(cred)))
	req.Header.Add("Authorization", basic)
	return req, nil
}
