package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/slack-go/slack"
)

func remove_file(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func copy_file(read_file string, write_file string) error {
	read_data, err := os.Open(read_file)
	if err != nil {
		return err
	}
	defer read_data.Close()

	write_data, err := os.Create(write_file)
	if err != nil {
		return err
	}
	defer write_data.Close()

	_, err = io.Copy(write_data, read_data)
	if err != nil {
		return err
	}

	return nil
}

func post_file(api *slack.Client, filename string) error {
	_, _, err := api.PostMessage("#notice-whiteboard", slack.MsgOptionBlocks(
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: "野茂が投げた！",
			},
			Accessory: slack.NewAccessory(
				slack.NewImageBlockElement(filename, "NOMO:"+filename),
			),
		},
	))
	if err != nil {
		return err
	}

	return nil
}

func loop(read_file string) {
	var timeZoneJST = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.Local = timeZoneJST
	time.LoadLocation("Asia/Tokyo")
	const time_layout = "20060102_150405"

	for {
		_, err := os.Stat(read_file)
		if !os.IsNotExist(err) {
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}

	t := time.Now()
	format_time := t.Format(time_layout)
	fmt.Println("[" + format_time + "]\n>>> Find PDF file.")
	write_file := "/var/www/out/" + format_time + ".pdf"

	err := copy_file(read_file, write_file)
	if err != nil {
		fmt.Println("!!! Cannot copy the file.")
		return
	}
	fmt.Println(">>> Copy Succeed.")

	err = remove_file(read_file)
	if err != nil {
		fmt.Println("!!! Cannot remove original file.")
		return
	}
	fmt.Println(">>> Remove Succeed.")

	err = post_file(slack.New(os.Getenv("SLACK_API_TOKEN")), write_file)
	if err != nil {
		fmt.Println("!!! Cannot post the file.")
		return
	}
	fmt.Println(">>> Post Succeed.")

	time.Sleep(5000 * time.Millisecond)
}

func main() {
	fmt.Println("Start Nomo Process.")
	read_file := "/var/www/monitored/CBImage.pdf"
	for {
		loop(read_file)
		time.Sleep(5000 * time.Millisecond)
	}
}
