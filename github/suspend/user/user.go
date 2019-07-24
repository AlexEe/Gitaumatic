package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"log"
	listUser "omniactl/github/list/user"
	githubLogin "omniactl/login/github"
)

var Client *github.Client

func SuspendUser(username string, reason string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Suspend user from Github")
	username = listUser.CheckUsername(username)
	githubUser := listUser.GetGithubUser(username)
	listUser.PrintUserInfo(githubUser)
	reason = CheckReason(reason)
	check := PromptSuspend(username)
	switch check {
	case "yes":
		SuspendFromGithub(username, reason)
	case "no":
		return
	}
}

func CheckReason(reason string) string {
	switch reason {
	case "":
		reason := PromptReason()
		return reason
	default:
		return reason
	}
}

func PromptReason() string {
	validate := func(input string) error {
		if input == "" {
			return errors.New("Enter a reason for the user's suspension")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Reason for suspension",
		Validate:  validate,
		Templates: templates,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func PromptSuspend(username string) string {
	prompt := promptui.Select{
		Label: "Suspend the user? ",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

type Reason struct {
	Reason string `json:"reason"`
}

// SuspendFromGithub suspends an account
func SuspendFromGithub(username string, reason string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	Client := githubLogin.CreateClient()
	body := Reason{reason}

	url := fmt.Sprintf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/users/%v/suspended", username)

	req, err := Client.NewRequest("PUT", url, body)
	if err != nil {
		log.Fatalln("Error creating HTTP request:\n", err)
	}
	_, err = Client.Do(context.Background(), req, nil)
	if err != nil {
		log.Fatalln("Error suspending user", err)
	}
	fmt.Println("")
	whiteBold.Printf("User '%v' suspended.", username)
	fmt.Println("")
}
