package team

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"log"
	createOrg "omniactl/github/create/org"
	createUser "omniactl/github/create/user"
	githubLogin "omniactl/login/github"
	"regexp"
	"sort"
)

// CreateTeam creates a new Github team based on flag or prompt input
func CreateTeam(team string, org string, teamDescription string, teamMaintainers []string, teamPrivacy string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	greenBold := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed)
	magentaBold.Println("Action selected: Create a new Github team")

	if org == "" {
		org = PromptOrg()
		team = PromptTeam(org)
		CreateGithubTeam(team, org, teamDescription, teamMaintainers, teamPrivacy)
	} else {
		check := createOrg.CheckIfOrgExists(org)
		switch check {
		case true:
			greenBold.Print("Organisation ")
			fmt.Println(org)
			if team == "" {
				team = PromptTeam(org)
				CreateGithubTeam(team, org, teamDescription, teamMaintainers, teamPrivacy)
			} else {
				allTeams := createUser.GetTeamsForOrg(org)
				check = CheckIfTeamExists(team, allTeams)
				switch check {
				case true:
					red.Printf("Team '%v' already exists.", team)
					fmt.Println("")
					team = PromptTeam(org)
					CreateGithubTeam(team, org, teamDescription, teamMaintainers, teamPrivacy)
				default:
					greenBold.Print("Team name ")
					fmt.Println(org)
					CreateGithubTeam(team, org, teamDescription, teamMaintainers, teamPrivacy)
				}
			}
		default:
			red.Printf("Organisation '%v' does not exist.", org)
			fmt.Println("")
			org = PromptOrg()
			team = PromptTeam(org)
			CreateGithubTeam(team, org, teamDescription, teamMaintainers, teamPrivacy)
		}
	}
}

// Team holds information to be passed in Http request
type Team struct {
	Name        string `json:"name"`
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
}

// CreateGithubTeam sends an HTTP Post request to create the team with the user input
func CreateGithubTeam(team string, org string, teamDescription string, teamMaintainers []string, teamPrivacy string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	Client := githubLogin.CreateClient()
	body := Team{Name: team, Description: teamDescription, Privacy: teamPrivacy}

	if teamDescription == "" {
		teamDescription = PromptDescription()
	}

	url := fmt.Sprintf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/orgs/%v/teams", org)

	req, err := Client.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalln("Error creating new request:\n", err)
	}

	newTeam := Team{}
	_, err = Client.Do(context.Background(), req, &newTeam)
	if err != nil {
		log.Fatalf("Error creating new team:\n%v", err)
	}

	team = newTeam.Name

	fmt.Println("")
	whiteBold.Printf("Team '%v' created in Github organisation '%v'.", team, org)
	fmt.Println("")
}

// PromptDescription asks if a description for the new team should be added
func PromptDescription() string {
	prompt := promptui.Select{
		Label: "Add a description for the team?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch result {
	case "yes":
		templates := &promptui.PromptTemplates{
			Success: "{{ . | green | bold }} ",
		}
		prompt := promptui.Prompt{
			Label:     "Description",
			Templates: templates,
		}
		result, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		return result
	default:
		return ""
	}
}

// CheckIfTeamExists checks if a team with same name already exists in organisation
func CheckIfTeamExists(input string, allTeams map[string]createUser.Team) bool {
	for k := range allTeams {
		if k == input {
			return true
		}
	}
	return false
}

// PromptTeam prompts user to enter a name for new team
func PromptTeam(org string) string {
	allTeams := createUser.GetTeamsForOrg(org)
	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^[a-z0-9._%+\\-]+$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Team name can only contain the following characters: A-Z, a-z, 0-9, -, _")
		}
		check := CheckIfTeamExists(input, allTeams)
		switch check {
		case false:
			return nil
		default:
			return errors.New("Team name already exists")
		}
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Team name",
		Validate:  validate,
		Templates: templates,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// PromptOrg asks which org the new team should be created in
func PromptOrg() string {
	greenBold := color.New(color.FgGreen, color.Bold)
	allOrgs := createUser.GetAllOrgs()
	var orgsSlice []string

	for k := range allOrgs {
		orgsSlice = append(orgsSlice, k)
		sort.Strings(orgsSlice)
	}

	prompt := promptui.Select{
		Label: "Select a Github organisation",
		Items: orgsSlice,
	}

	_, org, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	greenBold.Print("Organisation ")
	fmt.Println(org)
	return org
}
