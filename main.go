package main

import (
  "fmt"
  "log"
	"os"
	"io/ioutil"
	"regexp"

	"github.com/urfave/cli"
	"github.com/joho/godotenv"
)

var workableUrl string
var workableToken string
var workableLastSentId string
var slackToken string
var slackChannel string
var slackUsername string
var slackEmoji string
var slackMessageNew string
var slackMessageAll string
var slackMessageAllPost string

func init() {
	err := godotenv.Load("/go/bin/.env")
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  workableUrl = os.Getenv("WORKABLE_URL")
	workableToken = os.Getenv("WORKABLE_TOKEN")
	workableLastSentId = os.Getenv("WORKABLE_LAST_SENT_ID")
	slackToken = os.Getenv("SLACK_TOKEN")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	slackUsername = os.Getenv("SLACK_USERNAME")
	slackEmoji = os.Getenv("SLACK_EMOJI")
	slackMessageNew = os.Getenv("SLACK_MESSAGE_NEW")
	slackMessageAll = os.Getenv("SLACK_MESSAGE_ALL")
	slackMessageAllPost = os.Getenv("SLACK_MESSAGE_ALL_POST")
}

func main() {
  app := cli.NewApp()
  app.Name = "Workable-slack"
  app.Usage = "Notifies about Workable jobs"

	app.Commands = []cli.Command{
    {
      Name:    "new",
      Aliases: []string{"n"},
      Usage:   "notifies for the last added workable job",
      Action:  func(c *cli.Context) error {
				NotifyNewJob()
				return nil
      },
    },
    {
      Name:    "all",
      Aliases: []string{"a"},
      Usage:   "notifies for all jobs",
      Action:  func(c *cli.Context) error {
				NotifyAllJobs()
				return nil
      },
    },
  }

	err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
	}
}


func NotifyNewJob() {
	fmt.Println("Fetching new job offers from Workable")
	jobs := GetWorkableJobs()
	fired := false
	matched := false

	for _, job := range jobs.Jobs {

		// Debug line
		fmt.Printf("%s - %s \n", job.Title, job.Id)

		if (matched) {
			dotEnvContent, err := ioutil.ReadFile("/go/bin/.env")
			if err != nil {
				panic(err)
			}

			regex := regexp.MustCompile(`WORKABLE_LAST_SENT_ID=(.*)`)
			dotEnvContentString := string(dotEnvContent[:])
			dotEnvNewContentString := regex.ReplaceAllString(dotEnvContentString, fmt.Sprintf("WORKABLE_LAST_SENT_ID=%s", job.Id))

			fmt.Printf("Writing job %s id %s to .env \n", job.Title, job.Id)

			err = ioutil.WriteFile("/go/bin/.env", []byte(dotEnvNewContentString), 0)
			if err != nil {
				panic(err)
			}

			fired = true
			message := GetSingleJobSlackMessage(job)
			SendMessage(message)
	  }

		if (job.Id == workableLastSentId) {
			matched = true
		}
  }
	if !fired {
		fmt.Printf("No new job found (%v older jobs found)\n", len(jobs.Jobs))
	}
}

func NotifyAllJobs() {
	fmt.Println("Fetching all job offers from Workable")
	jobs := GetWorkableJobs()
	fmt.Printf("Notifying users for all %v jobs \n", len(jobs.Jobs))

	message := GetAllJobsSlackMessage(jobs)
	SendMessage(message)
	postMessage := GetAllJobsSlackPostMessage()
	SendMessage(postMessage)
}
