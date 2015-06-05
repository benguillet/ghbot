package main

import (
	"log"
	"os"

	"github.com/benjamin-guillet/ghbot"
	"github.com/benjamin-guillet/ghbot/greenhouse"
	"github.com/benjamin-guillet/ghbot/notifier/slack"
	"github.com/codegangsta/cli"
	"github.com/robfig/cron"
)

const (
	Name    = "GHBot"
	Usage   = "A bot for Greenhouse.io"
	Version = "0.0.1"
	CmdRun  = "run"

	FlagSlackHook       = "slack.hook"
	FlagCronSchedule    = "cron.schedule"
	FlagGreenhouseToken = "gh.token"
	FlagGreenhouseURL   = "gh.url"

	EnvSlackHook       = "SLACK_HOOK"
	EnvCronSchedule    = "CRON_SCHEDULE"
	EnvGreenhouseToken = "GH_TOKEN"
	EnvGreenhouseURL   = "GH_URL"
)

var Commands = []cli.Command{
	{
		Name:      CmdRun,
		ShortName: "r",
		Usage:     "Run the GHBot",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   FlagSlackHook,
				Value:  "",
				Usage:  "The slack incoming webhook to send notitications to",
				EnvVar: EnvSlackHook,
			},
			cli.StringFlag{
				Name:   FlagCronSchedule,
				Value:  "0 0 15 * * 1-5", // Weekdays at 3:00 PM UTC, i.e. 8:00 AM PDT
				Usage:  "The cron schedule to check the daily interviews and send a summary",
				EnvVar: EnvCronSchedule,
			},
			cli.StringFlag{
				Name:   FlagGreenhouseToken,
				Value:  "",
				Usage:  "Greenhouse Harvest API token",
				EnvVar: EnvGreenhouseToken,
			},
			cli.StringFlag{
				Name:   FlagGreenhouseURL,
				Value:  "https://harvest.greenhouse.io/v1",
				Usage:  "Greenhouse Harvest API Base URL",
				EnvVar: EnvGreenhouseURL,
			},
		},
		Action: run,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = Name
	app.Usage = Usage
	app.Version = Version
	app.Commands = Commands
	app.Run(os.Args)
}

func run(c *cli.Context) {
	cron := cron.New()

	hook := c.String(FlagSlackHook)
	if hook == "" {
		log.Fatal("Missing --%s or %s env variable.", FlagSlackHook, EnvSlackHook)
	}
	slack := slack.New(hook)

	token := c.String(FlagGreenhouseToken)
	if token == "" {
		log.Fatal("Missing --%s or %s env variable.", FlagGreenhouseToken, EnvGreenhouseToken)
	}
	ghclient := greenhouse.New(c.String(FlagGreenhouseURL), token)

	ghbot := ghbot.New(cron, slack, ghclient)
	ghbot.Run(c.String(FlagCronSchedule), ghbot.GetInterviews)
}
