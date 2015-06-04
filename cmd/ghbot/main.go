package main

import (
	"log"
	"os"

	"github.com/benjamin-guillet/ghbot"
	"github.com/benjamin-guillet/ghbot/notifier/slack"
	"github.com/codegangsta/cli"
	"github.com/robfig/cron"
)

const (
	CmdRun        = "run"
	FlagSlackHook = "slack.hook"
	Name          = "GHBot"
	Usage         = "A bot for Greenhouse.io"
	Version       = "0.0.1"
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
				EnvVar: "SLACK_HOOK",
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
		log.Fatal("Missing --slack.hook or SLACK_HOOK env variable.")
	}

	slack := slack.New(hook)

	gh := ghbot.New(cron, slack)

	// Run on weekdays at 3:00 PM UTC, i.e. 8:00 AM PDT
	// gh.Run("0 0 15 * * 1-5", ghbot.GetInterviews)

	// TODO: schedule should be set via ENV var and flag
	gh.Run("1 * * * * *", gh.GetInterviews)
}
