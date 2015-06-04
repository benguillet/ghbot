package slack

import (
	"fmt"

	"github.com/benjamin-guillet/ghbot/notifier"
)

const Green = "#0f0"

// Notifier is an implementation of the notifier.Notifier interface
// that forwards deployment core events to Slack.
type Notifier struct {
	client interface {
		Notify(*Payload) error
	}
}

func New(url string) *Notifier {
	c := newClient(nil)
	c.URL = url

	return &Notifier{
		client: c,
	}
}

func (n *Notifier) Notify(notification *notifier.Notification) error {
	fmt.Printf("%+v\n", notification)

	p := &Payload{
		Username: "Greenhouse Interviews Schedule",
		Channel:  "#testthebot",
		Text:     notification.Candidate,
		Attachments: []*Attachment{
			&Attachment{
				Text:     "test",
				Fallback: "test",
				Color:    Green,
			},
		},
	}

	n.client.Notify(p)

	return nil
}
