package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/avast/retry-go/v4"
	"github.com/go-co-op/gocron/v2"
	"github.com/urfave/cli/v2"
)

var lastSolved int = 0

func onSix(handle string) {
	userInfo, err := retry.DoWithData(
		func() (*UserInfo, error) {
			info, err := GetUserInfo(handle)
			if err != nil {
				return nil, err
			}
			return info, nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)

	if err == nil {
		lastSolved = userInfo.SolvedCount
		log.Printf("lastSolved: %d\n", lastSolved)
	} else {
		log.Printf("error: %v\n", err)
	}
}

func onTime(handle string, baseurl string, topic string, message string) {
	userInfo, err := retry.DoWithData(
		func() (*UserInfo, error) {
			info, err := GetUserInfo(handle)
			if err != nil {
				return nil, err
			}
			return info, nil
		},
		retry.Attempts(3),
		retry.Delay(1*time.Second),
	)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if userInfo.SolvedCount <= lastSolved {
		err = SendNtfy(baseurl, topic, message)
		log.Printf("send notification. solved count: %d\n", userInfo.SolvedCount)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

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
		}
	}

	tz, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return err
	}
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(tz))
	if err != nil {
		return err
	}
	defer func() { _ = scheduler.Shutdown() }()

	_, err = scheduler.NewJob(
		gocron.DailyJob(
			0,
			gocron.NewAtTimes(
				gocron.NewAtTime(6, 0, 10)),
		),
		gocron.NewTask(onSix, handle),
	)

	if err != nil {
		return err
	}

	atTimes := []gocron.AtTime{}
	for _, t := range parsedTimes {
		atTimes = append(atTimes, gocron.NewAtTime(
			uint(t.Hour()),
			uint(t.Minute()),
			uint(t.Second())),
		)
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(
			0,
			gocron.NewAtTimes(atTimes[0], atTimes[1:]...),
		),
		gocron.NewTask(onTime, handle, baseurl, topic, message),
	)

	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Scheduler started with values: baseurl=%s, topic=%s, handle=%s, message=%s, times=%v\n", baseurl, topic, handle, message, times)
	log.Print(msg)
	_ = SendNtfy(baseurl, topic, msg)
	scheduler.Start()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	<-sigc
	log.Println("Scheduler stopped")
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
