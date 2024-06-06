package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err == nil {
		return t, nil
	}
	t2, err2 := time.Parse(time.TimeOnly, s)
	if err2 == nil {
		return t2, nil
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s (available formats: TT:MM, TT:MM:SS)", s)
}

func appAction(c *cli.Context) error {
	baseurl := c.String("baseurl")
	topic := c.String("topic")
	handle := c.String("handle")
	message := c.String("message")
	times := c.StringSlice("times")

	parsedTimes := []time.Time{}
	for _, s := range times {
		t, err := parseTime(s)
		if err == nil {
			parsedTimes = append(parsedTimes, t)
		} else {
			fmt.Printf("error: %v %s\n", err, s)
		}
	}

	fmt.Printf("baseurl: %s\n", baseurl)
	fmt.Printf("topic: %s\n", topic)
	fmt.Printf("handle: %s\n", handle)
	fmt.Printf("message: %s\n", message)
	fmt.Printf("times: %v\n", times)

	for _, pt := range parsedTimes {
		fmt.Printf("parsed time: %v %d %d %d\n", pt, pt.Hour(), pt.Minute(), pt.Second())
	}

	return nil
}

func createApp() *cli.App {
	app := &cli.App{
		Name:   "dops",
		Usage:  "send do problem solving notification",
		Action: appAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "baseurl",
				Aliases:     []string{"b"},
				Usage:       "base `URL` for ntfy",
				Value:       "https://ntfy.sh",
				EnvVars:     []string{"DOPS_BASEURL"},
				DefaultText: "https://ntfy.sh",
			},
			&cli.StringFlag{
				Name:     "topic",
				Aliases:  []string{"t"},
				Usage:    "`TOPIC` for the notification",
				EnvVars:  []string{"DOPS_TOPIC"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "handle",
				Aliases:  []string{"H"},
				Usage:    "solved.ac handle",
				EnvVars:  []string{"DOPS_HANDLE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m"},
				Usage:   "message for the notification",
				Value:   "문제풀어!",
				EnvVars: []string{"DOPS_MESSAGE"},
			},
			&cli.StringSliceFlag{
				Name:        "times",
				Aliases:     []string{"T"},
				Usage:       "times for the notification (TT:MM, TT:MM:SS)",
				Value:       cli.NewStringSlice("09:00", "21:00"),
				DefaultText: "[09:00 21:00]",
				Action: func(c *cli.Context, args []string) error {
					for _, s := range args {
						if _, err := parseTime(s); err != nil {
							return err
						}
					}
					return nil
				},
			},
		}}

	return app
}

func main() {
	app := createApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
