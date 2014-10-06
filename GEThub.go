package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

type Config struct {
	Token string
	Org   string
	Repo  string
}

func readConfig(filename string) (Config, error) {

	file, err := os.Open(filename)
	config := Config{}

	if err != nil {
		return config, errors.New("Cannot find config file!")
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)

	if err != nil {
		fmt.Println("error:", err)
	}

	return config, nil
}

func main() {
	configFile := flag.String("config", "conf.json",
		"Configuration file. Where your app token lives")

	flag.Parse()

	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: config.Token},
	}

	client := github.NewClient(t.Client())
	repos, _, _ := client.Repositories.ListForks(config.Org, config.Repo, nil)
	messages := make([]string, 0)

	for _, fork := range repos {
		comms, _, _ := client.Repositories.ListCommits(*fork.Owner.Login,
			*fork.Name, nil)
		for _, commit := range comms {
			m := *commit.Commit.Message
			brief := strings.Split(m, "\n")[0]
			messages = append(messages, brief)
		}
	}

	for _, m := range messages {
		fmt.Println(m)
	}
}
