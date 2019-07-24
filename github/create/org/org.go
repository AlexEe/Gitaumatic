package org

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"log"
	githubLogin "omniactl/login/github"
	createUser "omniactl/github/create/user"
	"github.com/manifoldco/promptui"
	"errors"
	"regexp"
)

// Client which interacts with Github API
var Client *github.Client

type Org struct {
	Login       string `json:"login"`
	Admin       string `json:"admin"`
	ProfileName string `json:"profile_name"`
}

func CreateOrg(orgLogin string, orgProfile string, orgAdmin string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Create new Github organisation")

	switch orgLogin {
	case "":
		orgLogin = PromptNewOrgLogin()
		orgAdmin = PromptNewOrgAdmin()
		orgProfile = PromptNewOrgProfile()
		CreateGithubOrg(orgLogin, orgProfile, orgAdmin)
	default:
		check := CheckIfOrgExists(orgLogin)
		if check == true {
			fmt.Println("Organisation '%v' already exists.", orgLogin)
			orgLogin = PromptNewOrgLogin()
			orgAdmin = PromptNewOrgAdmin()
			orgProfile = PromptNewOrgProfile()
			CreateGithubOrg(orgLogin, orgProfile, orgAdmin)
		} else if check == false {
			switch orgAdmin {
			case "":
				orgAdmin = PromptNewOrgAdmin()
				orgProfile = PromptNewOrgProfile()
				CreateGithubOrg(orgLogin, orgProfile, orgAdmin)
			default:
				check := CheckIfUserExists(orgAdmin)
				if check == true {
					switch orgProfile {
					case "":
						orgProfile = PromptNewOrgProfile()
						CreateGithubOrg(orgLogin, orgProfile, orgAdmin)
					default:
						CreateGithubOrg(orgLogin, orgProfile, orgAdmin)
					}
				} else if check == false {
					fmt.Println("User 'v%' does not exist.", orgAdmin)
					orgAdmin = PromptNewOrgAdmin()
					orgProfile = PromptNewOrgProfile()
					CreateGithubOrg(orgLogin, orgProfile, orgAdmin)
				}
			}
		}
	}
}

func PromptNewOrgProfile() string {
	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     "Org display name",
		Templates: templates,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func PromptNewOrgAdmin() string {
	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^e[0-9]{6}$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Username must be in the format of a State Street Lan ID, e.g. 'e123456'")
		}
		check := CheckIfUserExists(input)
		switch check {
		case true:
			return nil
		default:
			return errors.New("User does not exist. Please enter a valid Github Login, e.g. 'e123456'")
		}
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     "Org admin",
		Validate: validate,
		Templates: templates,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}


func PromptNewOrgLogin() string {
	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^[a-z0-9._%+\\-]+$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Organization name can only contain the following characters: A-Z, a-z, 0-9, -, _")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     "Org login name",
		Validate: validate,
		Templates: templates,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func CheckIfUserExists(orgAdmin string) bool {
	Client := githubLogin.CreateClient()
	_, _, err := Client.Users.Get(context.Background(), orgAdmin)
	if err != nil {
		return false
	}
	return true
}

func CheckIfOrgExists(orgLogin string) bool {
	allOrgs := createUser.GetAllOrgs()
	for _, v := range allOrgs {
		if v.Name == orgLogin {
			return true
		} 
	}
	return false
}

func CreateGithubOrg(orgLogin string, orgProfile string, orgAdmin string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)

	Client = githubLogin.CreateClient()

	body := Org{orgLogin, orgAdmin, orgProfile}

	req, err := Client.NewRequest("POST", "https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/admin/organizations", body)
	if err != nil {
		log.Fatalln("Error creating new request:\n", err)
	}

	newOrg := Org{}
	_, err = Client.Do(context.Background(), req, &newOrg)

	if err != nil {
		log.Fatalf("Error creating new organisation:\n%v", err)
	}

	orgLogin = newOrg.Login
	fmt.Println("")
	whiteBold.Println("New Github organisation created: ")
	whiteBold.Println("Name:", orgLogin, " Admin:", orgAdmin, " Display name:", orgProfile)
}
