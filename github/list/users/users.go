package users

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"log"
	createUser "omniactl/github/create/user"
	listUser "omniactl/github/list/user"
)

// Client which interacts with Github API
var Client *github.Client

// ListUsers lists information about multiple users
func ListUsers(usernames []string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: List information for multiple users")
	var username string
	switch len(usernames) {
	case 0:
		username := createUser.GetUsername(username)
		githubUser := listUser.GetGithubUser(username)
		listUser.PrintUserInfo(githubUser)
		fmt.Println("")
		PromptAnotherUser()
	default:
		for _, username := range usernames {
			username = createUser.GetUsername(username)
			githubUser := listUser.GetGithubUser(username)
			listUser.PrintUserInfo(githubUser)
			fmt.Println("")
		}
	}
}

// PromptAnotherUser asks if info about another user should be fetched
func PromptAnotherUser() {
	var username string
	prompt := promptui.Select{
		Label: "List another user?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if result == "yes" {
		username := createUser.GetUsername(username)
		githubUser := listUser.GetGithubUser(username)
		listUser.PrintUserInfo(githubUser)
		fmt.Println("")
		PromptAnotherUser()
	} else if result == "no" {
		return
	}
}
