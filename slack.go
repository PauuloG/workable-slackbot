package main

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"net/http"
	"log"
	"encoding/json"
)

var attachmentsColors = []string{
	"#6735EE",
	"#00BCD6",
	"#2646E3",
	"#1BCE40",
}

var attachmentsEmojis = []string{
	"🔍",
	"✨",
	"🚀",
	"☄️",
}

type SlackMessage struct {
	// Username      string `json:"username"`
	// IconEmoji     string `json:"icon_emoji"`
	Text					string `json:"text"`
	Channel 			string `json:"channel"`
	Token					string `json:"token"`
	Attachments   []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Fallback		string `json:"fallback"`
	Color				string `json:"color"`
	Fields			[]SlackAttachmentField `json:"fields"`
}

type SlackAttachmentField struct {
	Title				string `json:"title"`
	Value				string `json:"value"`
}

func GetDefaultSlackMessage() SlackMessage {
	var slackMessage SlackMessage
	// slackMessage.Username = slackUsername
	// slackMessage.IconEmoji = slackEmoji
	slackMessage.Channel = slackChannel

	return slackMessage
}

func GetSingleJobSlackMessage(job workableJob) SlackMessage {

	messageString := fmt.Sprintf(slackMessageNew, job.Title, job.Department, job.Shortlink)

	slackMessage := GetDefaultSlackMessage()
	slackMessage.Text = messageString

	return slackMessage
}

func GetAllJobsSlackMessage(jobs workableJobs) SlackMessage {
	messageString := fmt.Sprintf(slackMessageAll)

	message := GetDefaultSlackMessage()
	message.Text = messageString

	var attachments []SlackAttachment

	for i, job := range jobs.Jobs {
		modulo := i % 4
		color := attachmentsColors[modulo]
		emoji := attachmentsEmojis[modulo]

		var slackAttachmentField SlackAttachmentField
		slackAttachmentField.Title = fmt.Sprintf("%s", job.Title)
		slackAttachmentField.Value = fmt.Sprintf("%s %s \n 🔗 %s", emoji, job.Department, job.Shortlink)

		var slackAttachment SlackAttachment
		slackAttachment.Fallback = fmt.Sprintf("%s - %s - %s", job.Title, job.Department, job.Shortlink)
		slackAttachment.Color = color
		slackAttachment.Fields = [] SlackAttachmentField {slackAttachmentField}

		attachments = append(attachments, slackAttachment)
	}

	message.Attachments = attachments

	return message
}

func GetAllJobsSlackPostMessage () SlackMessage {

	slackMessage := GetDefaultSlackMessage()
	slackMessage.Text = slackMessageAllPost

	return slackMessage
}

func SendMessage(message SlackMessage) []byte {

	message.Token = slackToken
	jsonBody, err := json.Marshal(message)

	if err != nil {
		log.Fatal(err)
	}

	slackUrl := fmt.Sprintf("%s/%s", "https://slack.com/api/", "chat.postMessage")

	client := &http.Client {}
	req, err := http.NewRequest("POST", slackUrl, bytes.NewBuffer(jsonBody))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", slackToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}
