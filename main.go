package main

import (
  "fmt"
  "log"
	"os"
	"time"
	"strconv"

	"github.com/urfave/cli"
	"github.com/joho/godotenv"
)

var workableUrl string
var workableToken string
var workableLastPostedTreshold int
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
	slackToken = os.Getenv("SLACK_TOKEN")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	slackUsername = os.Getenv("SLACK_USERNAME")
	slackEmoji = os.Getenv("SLACK_EMOJI")
	slackMessageNew = os.Getenv("SLACK_MESSAGE_NEW")
	slackMessageAll = os.Getenv("SLACK_MESSAGE_ALL")
	slackMessageAllPost = os.Getenv("SLACK_MESSAGE_ALL_POST")

	if s, err := strconv.Atoi(os.Getenv("WORKABLE_LAST_POSTED_THRESHOLD")); err == nil {
		workableLastPostedTreshold = int(s)
	}
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
	for _, job := range jobs.Jobs {
		loc, _ := time.LoadLocation("UTC")
		// setup a start and end time
		createdAt := job.CreatedAt.In(loc)
		now := time.Now().In(loc)

		// get the diff
		diff := now.Sub(createdAt)
		diffMinutes := int(diff.Minutes())

		// Debug line
		fmt.Printf("%s - %v \n", job.Title, diffMinutes)

		if diffMinutes < workableLastPostedTreshold {
			fmt.Printf("Notifying users for job \"%s\" (%s) - posted on %v \n", job.Title, job.Id, job.CreatedAt)
			message := GetSingleJobSlackMessage(job)
			SendMessage(message)
			fired = true
		}
	}
	if !fired {
		fmt.Printf("No job matching notification criterias %v (%v older jobs found)\n", workableLastPostedTreshold, len(jobs.Jobs))
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
