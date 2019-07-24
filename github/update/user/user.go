package user

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"log"
	createUser "omniactl/github/create/user"
	listUser "omniactl/github/list/user"
	githubLogin "omniactl/login/github"
)

// Client which interacts with Github API
var Client *github.Client

// UpdateUser gets info about user and allows to make changes to their status, membership
func UpdateUser(username string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Update Github user")
	username = listUser.CheckUsername(username)
	githubUser := listUser.GetGithubUser(username)
	listUser.PrintUserInfo(githubUser)
	// EmptyMap is a necessary placeholder for AddUserToOrgs function,
	// which in original function removes certain orgs from selection,
	// this is not desired in the update user package
	emptyMap := make(map[string]createUser.Org)

	result := PromptUpdate()
	switch result {
	case "yes":
		result = PromptAction()
		switch result {
		case "Add to Github organizations/ teams":
			newOrgs := createUser.SelectOrgs(emptyMap)
			createUser.AddUserToOrgs(username, newOrgs)
		case "Make site admin":
			MakeSiteAdmin(username)
		case "Remove from Github organization":
			RemoveOrgMember(username)
		}
	case "no":
		fmt.Println("")
		return
	}
}

// MakeSiteAdmin promotes user to site admin
func MakeSiteAdmin(username string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	Client := githubLogin.CreateClient()

	url := fmt.Sprintf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/users/%v/site_admin", username)

	req, err := Client.NewRequest("PUT", url, nil)
	if err != nil {
		log.Fatalln("Error creating HTTP request:\n", err)
	}
	resp, err := Client.Do(context.Background(), req, nil)
	if err != nil {
		log.Fatalln("Error making user site admin", err)
	}

	if resp.StatusCode == 204 {
		fmt.Println("")
		whiteBold.Printf("User '%v' promoted to site administrator.", username)
		fmt.Println("")
	} else {
		log.Fatalln("User was not promoted to site admin:", resp.Status)
	}
}

// PromptOrg asks to choose which org the user should be removed from
func PromptOrg(userOrgs map[string]createUser.Org) string {
	var orgsSlice []string

	for k := range userOrgs {
		orgsSlice = append(orgsSlice, k)
	}

	prompt := promptui.Select{
		Label: "Select the organisation from which the user should be removed",
		Items: orgsSlice,
	}
	_, org, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	fmt.Println(org)
	return org
}

// RemoveOrgMember removes a user from a selected org
func RemoveOrgMember(username string) {
	Client = githubLogin.CreateClient()
	userOrgs := GetOrgsForUser(username)
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	org := PromptOrg(userOrgs)

	resp, err := Client.Organizations.RemoveMember(context.Background(), org, username)
	if err != nil {
		log.Fatalln("Error removing member for organisation:", err)
	}
	if resp.StatusCode == 204 {
		fmt.Println("")
		whiteBold.Printf("User '%v' removed from Github organisation '%v'.", username, org)
		fmt.Println("")
	} else {
		log.Fatalln("User was not removed:", resp.Status)
	}
}

// PromptUpdate asks if a user should be updated
func PromptUpdate() string {
	prompt := promptui.Select{
		Label: "Update this user?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// PromptAction asks in which way user should be updated
func PromptAction() string {
	prompt := promptui.Select{
		Label: "Select action",
		Items: []string{"Add to Github organizations/ teams", "Remove from Github organization", "Make site admin"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// GetOrgsForUser gets the user's current orgs
func GetOrgsForUser(username string) map[string]createUser.Org {
	Client = githubLogin.CreateClient()
	allOrgs := createUser.GetAllOrgs()
	userOrgs := make(map[string]createUser.Org)

	for org := range allOrgs {
		isMember, _, err := Client.Organizations.IsMember(context.Background(), org, username)
		if err != nil {
			log.Fatalln("Error getting orgs for user:", err)
		}
		switch isMember {
		case true:
			userOrgs[org] = allOrgs[org]
			continue
		case false:
			continue
		}
	}
	return userOrgs
}
