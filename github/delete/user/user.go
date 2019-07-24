package user

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"log"
	listUser "omniactl/github/list/user"
	githubLogin "omniactl/login/github"
)

// Client which interacts with Github API
var Client *github.Client

// DeleteUser checks flags, gets user info and, once confirmed, deletes user
func DeleteUser(username string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Delete user from Github")
	username = listUser.CheckUsername(username)
	githubUser := listUser.GetGithubUser(username)
	listUser.PrintUserInfo(githubUser)
	check := PromptDelete(username)
	switch check {
	case "yes":
		DeleteFromGithub(username)
	case "no":
		return
	}
}

// PromptDelete prompts for confirmation of deletion
func PromptDelete(username string) string {
	prompt := promptui.Select{
		Label: "Delete the user? (This will remove all their repositories, gists, applications, and personal settings)",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// DeleteFromGithub deletes a user from Github, including all their repos
func DeleteFromGithub(username string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	Client := githubLogin.CreateClient()
	url := fmt.Sprintf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/admin/users/%v", username)

	req, err := Client.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalln("Error creating new request:\n", err)
	}
	_, err = Client.Do(context.Background(), req, nil)
	if err != nil {
		log.Fatalln("Error deleting user", err)
	}
	fmt.Println("")
	whiteBold.Printf("User '%v' deleted.", username)
	fmt.Println("")
}
