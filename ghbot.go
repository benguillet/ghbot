package ghbot

import (
	"log"

	"github.com/benjamin-guillet/ghbot/notifier"
	"github.com/robfig/cron"
)

type GHBot struct {
	cron     *cron.Cron
	notifier notifier.Notifier
}

func New(cron *cron.Cron, notif notifier.Notifier) *GHBot {
	gh := new(GHBot)
	gh.cron = cron
	gh.notifier = notif
	return gh
}

func (gh *GHBot) Run(sched string, toRun func()) (bool, error) {
	done := make(chan bool)

	gh.cron.AddFunc(sched, toRun)
	gh.cron.Start()

	return <-done, nil
}

func (gh *GHBot) GetInterviews() {
	log.Printf("Getting the interviews...\n")

	notif := &notifier.Notification{
		Candidate: "Bob",
	}

	gh.notifier.Notify(notif)
}
