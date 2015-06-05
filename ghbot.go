package ghbot

import (
	"fmt"
	"log"

	"github.com/benjamin-guillet/ghbot/greenhouse"
	"github.com/benjamin-guillet/ghbot/notifier"
	"github.com/robfig/cron"
)

type GHBot struct {
	cron     *cron.Cron
	notifier notifier.Notifier
	ghclient *greenhouse.Client
}

func New(cron *cron.Cron, notif notifier.Notifier, ghclient *greenhouse.Client) *GHBot {
	gh := new(GHBot)
	gh.cron = cron
	gh.notifier = notif
	gh.ghclient = ghclient
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

	app, err := gh.ghclient.GetApplications()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", app)
	// gh.ghclient.GetScheduledInterviews()

	// notif := &notifier.Notification{
	// 	Candidate: "Bob",
	// }
	// gh.notifier.Notify(notif)
}
