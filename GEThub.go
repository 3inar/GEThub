package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

type Config struct {
	Token string
}

func getAppToken(filename string) (string, error) {

	file, err := os.Open(filename)

	if err != nil {
		return "", errors.New("Cannot find config file!")
	}

	decoder := json.NewDecoder(file)
	config := Config{}

	err = decoder.Decode(&config)

	if err != nil {
		fmt.Println("error:", err)
	}

	return config.Token, nil
}

func main() {
	configFile := flag.String("configfile", "conf.json",
		"Configuration file. Where your app token lives")

	flag.Parse()

	token, err := getAppToken(*configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	client := github.NewClient(t.Client())
	repos, _, _ := client.Repositories.ListForks("uit-inf-3200", "Project-1",
		nil)
	messages := make([]string, 0)

	for _, fork := range repos {
		comms, _, _ := client.Repositories.ListCommits(*fork.Owner.Login,
			*fork.Name, nil)
		for _, commit := range comms {
			messages = append(messages, *commit.Commit.Message)
		}
	}

	for _, m := range messages {
		fmt.Println(m)
	}
}
