package main

import (
	"fmt"
	"net/url"

	"github.com/imroc/req/v3"
)

const NTFY string = "https://ntfy.sh"

func SendNtfy(baseurl string, topic string, message string) error {
	if baseurl == "" {
		baseurl = NTFY
	}

	if topic == "" {
		return fmt.Errorf("topic is empty")
	}

	url, err := url.JoinPath(baseurl, topic)
	if err != nil {
		return err
	}

	client := req.C()
	resp, err := client.R().
		SetBody(message).
		SetHeader("Title", "Problem Solving").
		Post(url)

	if err != nil || resp.IsErrorState() {
		if err == nil {
			err = fmt.Errorf("error: %s", resp.String())
		}
		return err
	}
	return nil
}
