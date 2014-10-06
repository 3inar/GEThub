package main

import (
	"fmt"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

func main() {
	// NB: go to github and make yourself an access token
	token := "your token goes here"
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
