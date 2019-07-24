package teams

import (
	createUser "omniactl/github/create/user"
	createOrg "omniactl/github/create/org"
	// githubLogin "omniactl/login/github"
	createTeam "omniactl/github/create/team"
	// "github.com/manifoldco/promptui"
	// "log"
	"fmt"
	"github.com/fatih/color"
	// "errors"
	// "regexp"
	// "context"
	// "sort"
)

// ListTeams lists all teams for all orgs or all teams for a specific org when flag is provided
func ListTeams(org string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: List all Github teams for all or selected organisations")

	CheckFlag(org)
}

func CheckFlag(org string) {
	if org == "" {
		ListAllTeams()
	} else {
		check := createOrg.CheckIfOrgExists(org)
		switch check {
		case true:
			ListOrgTeams(org)
		default:
			org := createTeam.PromptOrg()
			ListOrgTeams(org)
		}
	}
}

func ListAllTeams() {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	allOrgs := createUser.GetAllOrgs()

	for k, _ := range allOrgs {
		teams := createUser.GetTeamsForOrg(k)
		fmt.Println("")
		whiteBold.Print("Organisation:")
		fmt.Println("\t" + k)
		whiteBold.Print("Teams:")
		if len(teams) != 0 {
			for kk := range teams {
				fmt.Printf("\t\t%-30v", kk)
				fmt.Println("")
			}
		} else {
			fmt.Print("\n")
		}
	}
}


func ListOrgTeams(org string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	teams := createUser.GetTeamsForOrg(org)

	fmt.Println("")
	whiteBold.Print("Organisation:")
	fmt.Println("\t" + org)
	whiteBold.Print("Teams:")
	if len(teams) != 0 {
		for kk := range teams {
			fmt.Printf("\t\t%-30v", kk)
			fmt.Println("")
		}
	} else {
		fmt.Print("\n")
	}
	fmt.Println("")
}