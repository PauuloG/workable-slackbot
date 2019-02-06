package main

import (
  "fmt"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
)

type workableJobs struct {
	Jobs []workableJob `json:"jobs"`
}

type workableJob struct {
	Id 					string 		`json:"id"`
	Title 			string 		`json:"title"`
	Department 	string 		`json:"department"`
	Shortlink		string 		`json:"shortlink"`
	CreatedAt		time.Time	`json:"created_at"`
}

func GetWorkableJobs() workableJobs {
	client := &http.Client {}

	req, err := http.NewRequest("GET", workableUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", workableToken))
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	var jobs workableJobs
	err = json.Unmarshal(body, &jobs)

	if err != nil {
		log.Fatal(err)
	}

	return jobs
}
