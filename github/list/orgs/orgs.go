package orgs

import (
createTeam "omniactl/github/create/team"
createUser "omniactl/github/create/user"
// createOrg "omniactl/github/create/org"
githubLogin "omniactl/login/github"
listOrg "omniactl/github/list/org"
"github.com/manifoldco/promptui"
"log"
"os"
"fmt"
"github.com/fatih/color"
// "errors"
// "regexp"
"context"
)

// ListOrgs lists all available Github orgs
func ListOrgs() {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: List all Github organisations")

	ListAllOrgsInfo()
	fmt.Println("")
	result := PromptMoreInfo()
	switch result {
	case "yes":
		org := createTeam.PromptOrg()
		result := listOrg.PromptAction()
		switch result {
		case "Organisation stats":
			listOrg.ListOrgStats(org)
		case "Organisation members":
			listOrg.ListOrgMembers(org)
		case "Organisation repositories":
			listOrg.ListOrgRepos(org)
		case "Organisation teams":
			listOrg.ListOrgTeams(org)
		}
	default:
		os.Exit(0)
	}
}

// ListAllOrgsInfo provides an overview over all available Github orgs
func ListAllOrgsInfo() {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	Client := githubLogin.CreateClient()
	allOrgs := createUser.GetAllOrgs()

	fmt.Println("")
	whiteBold.Println("Github organisations:")
	for k, v := range allOrgs {
		org, _, err := Client.Organizations.Get(context.Background(), k)
		if err != nil {
			log.Fatalln("Error getting organisation information from Github:", err)
		}
		fmt.Printf("Name: %-25v | ID: %-15v | Private repos: %-12v | Public repos: %-15v", v.Name, v.ID, org.GetTotalPrivateRepos(), org.GetPublicRepos())
		fmt.Println("")
	}
}

// PromptMoreInfo checks if user wants more detail about a given org
func PromptMoreInfo() string {
	prompt := promptui.Select{
		Label: "Get more information on one of the organisations?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}