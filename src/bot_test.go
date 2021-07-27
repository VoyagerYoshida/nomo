package main

import (
	"os"
	"testing"
	"time"

	"github.com/slack-go/slack"
)

func TestBot(t *testing.T) {
	read_file := "/var/www/monitored/TestImage.pdf"

	var timeZoneJST = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.Local = timeZoneJST
	time.LoadLocation("Asia/Tokyo")
	const time_layout = "20060102_150405"

	_, err := os.Stat(read_file)
	if err != nil {
		t.Error("[Error]", err)
	}

	time_now := time.Now()
	format_time := time_now.Format(time_layout)
	write_file := "/var/www/out/" + format_time + ".pdf"

	err = copy_file(read_file, write_file)
	if err != nil {
		t.Error("[Error]", err)
	}
	err = remove_file(read_file)
	if err != nil {
		t.Error("[Error]", err)
	}
	err = post_file(slack.New(os.Getenv("SLACK_API_TOKEN")), format_time)
	if err != nil {
		t.Error("[Error]", err)
	}
}
