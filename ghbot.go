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

	apps = filterApplication(apps)
	if len(apps) == 0 {
		log.Printf("No applications active\n")
		return
	}

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
	i := 0

	for range ticker.C {
		app := apps[i]
		// fmt.Printf("%d\n", app.ID)
		candidate := &greenhouse.Candidate{}
		interviews := make([]*notifier.Interview, 0)

		scheds, err := gh.ghclient.GetScheduledInterviews(app.ID) //7825355
		if err != nil {
			log.Fatal(err)
		}

		for _, sched := range scheds {
			startTime, err := time.Parse(time.RFC3339Nano, sched.Start.DateTime)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("time: %s\n", startTime)
			// today := time.Now().UTC().Truncate(24 * time.Hour)
			// tomorrow := today.Add(24 * time.Hour)
			// fmt.Printf("%s\n", today)

			// TODO: convert every date to PST/PDT
			// if startTime.After(today) && startTime.Before(tomorrow) {
			candidate, err = gh.ghclient.GetCandidate(app.CandidateID) //6908986
			if err != nil {
				log.Fatal(err)
			}

			interview := notifier.Interview{
				Name:        sched.Interview.Name,
				Start:       startTime.String(), // TODO: PST and show hour, min only
				Interviewer: sched.Interviewers[0].Name,
			}
			interviews = append(interviews, &interview)
			// }

		}

		notification := &notifier.Notification{
			Candidate:  candidate.Name(),
			Role:       app.Jobs[0].Name,
			Interviews: interviews,
		}
		fmt.Printf("%+v\n", notification)

		for _, inter := range interviews {
			fmt.Printf("%+v\n", inter)
		}
		// notif <- notification

		i += 1
		if i == len(apps) {
			done <- true
		}
	}
}

func filterApplication(apps []*greenhouse.Application) []*greenhouse.Application {
	filteredApps := make([]*greenhouse.Application, 0)

	for i, app := range apps {
		fmt.Printf("app %d: status: %s, prospect: %t\n", i, app.Status, app.Prospect)
		if app.Status == "active" && !app.Prospect {
			filteredApps = append(filteredApps, app)
		}
	}

	return filteredApps
}
