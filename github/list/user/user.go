package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"log"
	createUser "omniactl/github/create/user"
	githubLogin "omniactl/login/github"
)

var (
	GreenBold *color.Color
	WhiteBold *color.Color
	Red       *color.Color
	Client    *github.Client
)

func init() {
	GreenBold = color.New(color.FgGreen, color.Bold)
	WhiteBold = color.New(color.FgHiWhite, color.Bold)
	Red = color.New(color.FgRed)
	Client = githubLogin.CreateClient()
}

// ListUser gets information about user from Github
func ListUser(username string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: List user information")
	username = CheckUsername(username)
	githubUser := GetGithubUser(username)
	PrintUserInfo(githubUser)
}

// CheckUsername checks flag input and if none was set prompts user
func CheckUsername(username string) string {
	if username != "" {
		check := createUser.CheckIfUserExists(username)
		if check == true {
			return username
		} else {
			username = PromptUsername()
			return username
		}
	} else {
		username = PromptUsername()
		return username
	}
}

// PromptUsername asks to enter user to be listed, and checks they exist
func PromptUsername() string {
	validate := func(input string) error {
		if input == "" {
			return errors.New("Enter a user to be listed")
		}
		check := createUser.CheckIfUserExists(input)
		switch check {
		case true:
			return nil
		default:
			return errors.New("User does not exist")
		}
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     "User",
		Validate:  validate,
		Templates: templates,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// GetGithubUser returns a github user when entering the username
func GetGithubUser(username string) *github.User {
	user, _, err := Client.Users.Get(context.Background(), username)
	if err != nil {
		log.Fatalln("Error getting stats about user:", err)
	}
	return user
}

// GetOrgsForUser lists all orgs which the user is a member of
func GetOrgsForUser(username string) map[string]createUser.Org {
	allOrgs := createUser.GetAllOrgs()
	UserOrgs := make(map[string]createUser.Org)

	for k, _ := range allOrgs {
		isMember, _, err := Client.Organizations.IsMember(context.Background(), k, username)
		if err != nil {
			log.Fatalln("Error listing orgs for user:", err)
		}
		switch isMember {
		case true:
			value, _ := allOrgs[k]
			UserOrgs[k] = value
			teams := GetTeamsForUser(username, k)
			value.Teams = teams
			UserOrgs[k] = value
			continue
		case false:
			continue
		}
	}
	return UserOrgs
}

// GetTeamsForUser lists all teams which user is a member of in specific org
func GetTeamsForUser(username string, org string) []createUser.Team {
	allTeams := createUser.GetTeamsForOrg(org)
	var userTeams []createUser.Team

	for _, v := range allTeams {
		isMember, _, err := Client.Teams.IsTeamMember(context.Background(), v.ID, username)
		if err != nil {
			log.Fatalln("Error listing teams for user:", err)
		}
		switch isMember {
		case true:
			userTeams = append(userTeams, v)
			continue
		default:
			continue
		}
	}
	return userTeams
}

// PrintUserInfo gets user object from Github and gets their id, email, orgs etc
func PrintUserInfo(user *github.User) {
	var OrgCount int

	GreenBold.Print("Login ")
	fmt.Println(user.GetLogin())
	GreenBold.Print("ID ")
	fmt.Println(user.GetID())
	GreenBold.Print("Username ")
	fmt.Println(user.GetName())
	GreenBold.Print("Email ")
	fmt.Println(user.GetEmail())
	GreenBold.Print("Is admin ")
	fmt.Println(user.GetSiteAdmin())
	GreenBold.Print("Time last updated ")
	fmt.Println(user.GetUpdatedAt())
	GreenBold.Print("Time created ")
	fmt.Println(user.GetCreatedAt())
	// IsOrgMember(user.GetLogin())
	userOrgs := GetOrgsForUser(user.GetLogin())
	for _, v := range userOrgs {
		OrgCount++
		GreenBold.Print("Organisation ", OrgCount, " ")
		fmt.Println(v.Name, " ")
		for _, vv := range v.Teams {
			GreenBold.Print("Team ")
			fmt.Println(vv.Name, " ")
		}
	}
}
