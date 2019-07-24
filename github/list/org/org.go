package org

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"log"
	createOrg "omniactl/github/create/org"
	createTeam "omniactl/github/create/team"
	createUser "omniactl/github/create/user"
	listTeam "omniactl/github/list/team"
	githubLogin "omniactl/login/github"
)

func ListOrg(org string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: List information about a Github organisation")

	org = CheckFlag(org)
	result := PromptAction()
	switch result {
	case "Organisation stats":
		ListOrgStats(org)
	case "Organisation members":
		ListOrgMembers(org)
	case "Organisation repositories":
		ListOrgRepos(org)
	case "Organisation teams":
		ListOrgTeams(org)
	}
}

// CheckFlag checks input from flag and prompts user if necessary
func CheckFlag(org string) string {
	red := color.New(color.FgRed)
	if org != "" {
		check := createOrg.CheckIfOrgExists(org)
		switch check {
		case true:
			return org
		default:
			red.Printf("Organisation '%v' does not exist.", org)
			fmt.Println("")
			org = createTeam.PromptOrg()
			return org
		}
	} else {
		org = createTeam.PromptOrg()
		return org
	}
}

// PromptAction asks user to select which info they want about selected team
func PromptAction() string {
	prompt := promptui.Select{
		Label: "Select action",
		Items: []string{"Organisation stats", "Organisation members", "Organisation repositories", "Organisation teams"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

func ListOrgStats(org string) {
	Client := githubLogin.CreateClient()
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	greenBold := color.New(color.FgGreen, color.Bold)

	gitOrg, _, err := Client.Organizations.Get(context.Background(), org)
	if err != nil {
		log.Fatalln("Error retrieving org information from Github:\n", err)
	}

	fmt.Println("")
	whiteBold.Println("Organisation stats:")
	greenBold.Print("Name ")
	fmt.Println(gitOrg.GetLogin())
	greenBold.Print("ID ")
	fmt.Println(gitOrg.GetID())
	greenBold.Print("Private repos ")
	fmt.Println(gitOrg.GetTotalPrivateRepos())
	greenBold.Print("Public repos ")
	fmt.Println(gitOrg.GetPublicRepos())
	greenBold.Print("Time created ")
	fmt.Println(gitOrg.GetCreatedAt())
	greenBold.Print("Time last updated ")
	fmt.Println(gitOrg.GetUpdatedAt())
	fmt.Println("")
}

func ListOrgRepos(org string) {
	Client := githubLogin.CreateClient()
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	PageCount := 1
	NextPage := true
	RepoCount := 1

	fmt.Println("")
	whiteBold.Printf("'%v' repositories: ", org)
	fmt.Println("")
	for NextPage == true {
		url := fmt.Sprintf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/orgs/%v/repos?page=%v&per_page=100", org, PageCount)

		req, err := Client.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln("Error creating HTTP request:\n", err)
		}

		// Create repo struct to save data of HTTP request
		orgRepos := listTeam.Repos{}
		_, err = Client.Do(context.Background(), req, &orgRepos)
		if err != nil {
			log.Fatalln("Error retrieving repos for team", err)
		}

		PageCount++

		switch len(orgRepos) {
		case 0:
			NextPage = false
		default:
			for _, v := range orgRepos {
				fmt.Printf("%-5v Name: %-40v | ID: %-15v | Updated at: %-40v | Private: %-20v", RepoCount, v.Name, v.ID, v.UpdatedAt, v.Private)
				fmt.Println("")
				RepoCount++
			}
		}
	}
}

func ListOrgMembers(org string) {
	Client := githubLogin.CreateClient()
	whiteBold := color.New(color.FgHiWhite, color.Bold)

	members, _, err := Client.Organizations.ListMembers(context.Background(), org, nil)
	if err != nil {
		log.Fatalln("Error getting info about organisation:", err)
	}

	whiteBold.Println("Organisation members:")
	for _, v := range members {
		user, _, err := Client.Users.Get(context.Background(), v.GetLogin())
		if err != nil {
			log.Fatalln("Error getting info about organisation members:", err)
		}

		fmt.Printf("Name: %-25v | Login: %-20v | ID: %-15v | Site admin: %-10v\n", user.GetName(), v.GetLogin(), v.GetID(), v.GetSiteAdmin())
	}
}

func ListOrgTeams(org string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	teams := createUser.GetTeamsForOrg(org)
	Client := githubLogin.CreateClient()

	fmt.Println("")
	whiteBold.Println("Teams:")
	for k, v := range teams {
		teamID := v.ID
		team, _, err := Client.Teams.GetTeam(context.Background(), teamID)
		if err != nil {
			log.Fatalf("Error getting information about Github team '%v': %v", k, err)
		}
		fmt.Printf("Name: %-25v | ID: %-10v | Permission: %-10v | No of repos: %-10v | Description: %-35v \n", team.GetName(), team.GetID(), team.GetPermission(), team.GetReposCount(), team.GetDescription())
	}
	fmt.Println("")
}
