package ghbot

import (
	"fmt"
	"log"
	"time"

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
	apps, err := gh.ghclient.GetApplications()
	if err != nil {
		log.Fatal(err)
	}

	// for _, app := range apps {
	// 	fmt.Printf("%+v\n", app.ID)
	// }

	notification := make(chan *notifier.Notification)
	ticker := time.NewTicker(time.Second * 1)
	done := make(chan bool)
	go getSchedules(gh, apps, ticker, done, notification)

	for {
		select {
		case notif := <-notification:
			log.Printf("Sending schedule!\n")
			fmt.Printf("%+v\n", notif)
			// gh.notifier.Notify(notif)
		case <-done:
			log.Printf("Done retrieving schedules\n")
			ticker.Stop()
			return
		}
	}
}

func getSchedules(gh *GHBot, apps []*greenhouse.Application, ticker *time.Ticker, done chan<- bool, notif chan<- *notifier.Notification) {
	// i := 0

	for range ticker.C {
		// app := apps[i]
		// fmt.Printf("%d\n", app.ID)

		scheds, err := gh.ghclient.GetScheduledInterviews(7825355)
		if err != nil {
			log.Fatal(err)
		}

		for _, sched := range scheds {
			t, err := time.Parse(time.RFC3339Nano, sched.Start.DateTime)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("time: %s", t)
			notif <- &notifier.Notification{
				Candidate: "Laura",
			}
		}

		// i += 1
		// if i == len(apps) {
		done <- true
		// }
	}
}
