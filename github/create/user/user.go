package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	githubLogin "omniactl/login/github"
	"os"
	"regexp"
	"sort"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
)

var (
	GreenBold   *color.Color
	MagentaBold *color.Color
	WhiteBold   *color.Color
	Red         *color.Color
	Client      *github.Client
)

func init() {
	GreenBold = color.New(color.FgGreen, color.Bold)
	MagentaBold = color.New(color.FgMagenta, color.Bold, color.Underline)
	WhiteBold = color.New(color.FgHiWhite, color.Bold)
	Red = color.New(color.FgRed)
	Client = githubLogin.CreateClient()
}

// AddUser gets values required for adding a new Github User, either through CLI flags or User prompt
func AddUser(username string, email string, org string, role string, teams []string) {
	MagentaBold.Println("Action selected: Add new user to Github")
	CheckLogin()
	username = GetUsername(username)
	email = GetEmail(email)
	switch org {
	case "":
		check := PromptAddUserToOrg()
		if check == "yes" {
			orgs := GetOrgs(org, role, teams)
			username, _ = CreateUser(username, email)
			AddUserToOrgs(username, orgs)
		} else {
			_, _ = CreateUser(username, email)
		}
	default:
		orgs := GetOrgs(org, role, teams)
		username, _ = CreateUser(username, email)
		AddUserToOrgs(username, orgs)
	}
}

// CheckLogin checks Github token by logging into Github
func CheckLogin() {
	// githubLogin.GithubLogin("check")
	err := githubLogin.CheckGithubLogin()
	if err != nil {
		log.Fatalln("Connection to Github failed:", err)
	}
}

// GetUsername receives username through flag or prompt input
func GetUsername(username string) string {
	if username != "" {
		check := CheckUsernameFormat(username)
		if check == true {
			GreenBold.Print("Username ")
			fmt.Println(username)
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

// CheckUsernameFormat checks flag input has correct format
func CheckUsernameFormat(input string) bool {
	red := color.New(color.FgRed)

	isCorrectFormat := regexp.MustCompile("^e[0-9]{6}$").MatchString
	if isCorrectFormat(input) {
		check := CheckIfUserExists(input)
		switch check {
		case true:
			red.Println("Username already exists.")
			result := PromptAbort()
			if result == "add another user" {
				return false
			} else {
				os.Exit(1)
			}

		case false:
			return true
		}
	} else {
		red.Println("Format incorrect for State Street Lan ID. E.g. 'e123456'")
		return false
	}
	return false
}

func PromptAbort() string {
	prompt := promptui.Select{
		Label: "Exit program or add another user?",
		Items: []string{"exit", "add another user"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

// CheckIfUserExists checks if user already exists in Github
func CheckIfUserExists(username string) bool {
	Client := githubLogin.CreateClient()
	_, _, err := Client.Users.Get(context.Background(), username)
	if err != nil {
		return false
	}
	return true
}

// PromptUsername prompts user for username and checks input
func PromptUsername() string {
	validate := func(input string) error {
		red := color.New(color.FgRed)
		isCorrectFormat := regexp.MustCompile("^e[0-9]{6}$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Username must be in the format of a State Street Lan ID, e.g. 'e123456'")
		}
		check := CheckIfUserExists(input)
		switch check {
		case true:
			red.Print("Username already exists.")
			fmt.Println("")
			return errors.New("Username already exists.")
		default:
			return nil
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

// GetEmail gets email either from flag or prompt input
func GetEmail(email string) string {
	if email != "" {
		check := CheckEmailFormat(email)
		if check == true {
			GreenBold.Print("Email ")
			fmt.Println(email)
			return email
		} else {
			email = PromptEmail()
			return email
		}
	} else {
		email = PromptEmail()
		return email
	}
}

// CheckEmailFormat checks flag input for correct format
func CheckEmailFormat(input string) bool {
	isCorrectFormat := regexp.MustCompile("^[a-z0-9._%+\\-]+@statestreet.com$").MatchString
	if isCorrectFormat(input) {
		return true
	} else {
		fmt.Println("Email must be the new user's State Street email address, e.g. 'example@statestreet.com'")
		return false
	}
}

// PromptEmail prompts user for input and checks it
func PromptEmail() string {
	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^[a-z0-9._%+\\-]+@statestreet.com$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Email must be the new user's State Street email address, e.g. 'example@statestreet.com'")
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     "Email",
		Validate:  validate,
		Templates: templates,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// Org structs holds info on Github org, with id, role and teams
type Org struct {
	Name  string
	ID    int64
	Role  string
	Teams []Team
}

// Team struct holds info on Github teams with ids
type Team struct {
	Name string
	ID   int64
}

// OrgCount keeps track of number of selected orgs for new user
var OrgCount int

// GetOrgs checks flag and prompt input and returns a map of all selected orgs
func GetOrgs(flagOrg string, flagRole string, flagTeams []string) map[string]Org {
	allOrgs := GetAllOrgs()
	flagOrgs := make(map[string]Org)
	promptOrgs := make(map[string]Org)
	sliceTeams := []Team{}

	// Checks if flag has been set, if not: prompt for Orgs
	if flagOrg == "" {
		promptOrgs = SelectOrgs(flagOrgs)
	} else if flagOrg != "" {
		// Check if org entered at flag exists
		check := CheckOrgExists(flagOrg)
		if check == true {
			OrgCount++
			GreenBold.Print("Organisation ", OrgCount, " ")
			fmt.Println(flagOrg)

			// Add new org to flagOrg map, getting name and ID from allOrgs
			value, ok := allOrgs[flagOrg]
			if ok {
				flagOrgs[flagOrg] = value
			} else {
				log.Fatalln("Error saving organisations selected as flags.")
			}

			// Add role for selected org, either from flag or prompt
			if flagRole != "" {
				check = CheckRoleExists(flagRole)
				if check == true {
					GreenBold.Print("Role ")
					fmt.Println(flagRole)
					value.Role = flagRole
					flagOrgs[flagOrg] = value
				} else {
					Red.Printf("Role '%v' does not exist.\n", flagRole)
					flagRole = PromptRole()
					value.Role = flagRole
					flagOrgs[flagOrg] = value
				}
			} else {
				flagRole = PromptRole()
				value.Role = flagRole
				flagOrgs[flagOrg] = value
			}

			// Add team for selected org
			if len(flagTeams) != 0 {
				for _, v := range flagTeams {
					check = CheckTeamExists(flagOrg, v)
					if check == true {
						GreenBold.Print("Team ")
						fmt.Println(v)
						allTeams := GetTeamsForOrg(flagOrg)
						id, ok := allTeams[v]
						if ok {
							sliceTeams = append(sliceTeams, id)
						}
					}
				}
				value.Teams = sliceTeams
				flagOrgs[flagOrg] = value
			} else {
				teams := GetTeamsForOrg(flagOrg)
				sliceTeams = SelectTeams(teams)
				value.Teams = sliceTeams
				flagOrgs[flagOrg] = value
				result := PromptAnotherOrg()
				if result == "yes" {
					promptOrgs = SelectOrgs(flagOrgs)
				}
			}
		} else {
			Red.Printf("Organisation '%v' does not exist.\n", flagOrg)
			promptOrgs = SelectOrgs(flagOrgs)
		}
	}
	// Add orgs from prompt to flagOrgs
	for k, v := range promptOrgs {
		flagOrgs[k] = v
	}
	selectedOrgs := flagOrgs
	return selectedOrgs
}

// CheckOrgExists checks if org provided by flag
// exists in list of Github orgs
func CheckOrgExists(orgName string) bool {
	allOrgs := GetAllOrgs()
	_, ok := allOrgs[orgName]
	if ok {
		return true
	} else {
		return false
	}
}

// CheckRoleExists checks if role provided by flag exists
func CheckRoleExists(role string) bool {
	roles := make(map[string]string)
	roles["member"] = "member"
	roles["admin"] = "admin"

	_, ok := roles[role]
	if ok {
		return true
	} else {
		return false
	}
}

// CheckTeamExists checks if team name provided by flag exists
// in selected Github org
func CheckTeamExists(orgName string, teamName string) bool {
	allTeams := GetTeamsForOrg(orgName)
	_, ok := allTeams[teamName]
	if ok {
		return true
	} else {
		return false
	}
}

// SelectOrgs prompts user to select orgs, roles in orgs and teams for the new user
func SelectOrgs(flagOrgs map[string]Org) map[string]Org {
	red := color.New(color.FgRed)
	userOrgs := map[string]Org{}
	addOrgs := true
	allOrgs := DeleteFlagOrgs(flagOrgs)
	s := CreateOrgList(allOrgs)

	// While true: keep prompting for adding a new organisation
	for addOrgs == true {
		// User selects one organisation by name with prompt
		orgName := PromptOrgNames(s)

		// Deletes selected org from list to be prompted so user can't select it twice
		for i, v := range s {
			if v == orgName {
				s = append(s[:i], s[i+1:]...)
				sort.Strings(s)
				break
			}
		}

		// Add new org to userOrgs list, adding information about org from allOrgs map
		value, ok := allOrgs[orgName]
		if ok {
			userOrgs[orgName] = value
		} else {
			log.Fatalln("Error saving selected organisation.")
		}

		// User chooses role for selected Org
		role := PromptRole()
		value.Role = role
		userOrgs[orgName] = value

		// Get list of all teams in chosen org
		teamsForOrg := GetTeamsForOrg(orgName)

		if len(s) != 0 && len(teamsForOrg) != 0 {
			teams := SelectTeams(teamsForOrg)

			value.Teams = teams

			if ok {
				userOrgs[orgName] = value
			} else {
				log.Fatalln("Error saving selected organisation.")
			}

			result := PromptAnotherOrg()

			switch result {
			case "yes":
				continue
			case "no":
				addOrgs = false
				break
			default:
				log.Fatalln("Prompt for organisations failed.")
			}
		} else if len(s) == 0 && len(teamsForOrg) != 0 {
			teams := SelectTeams(teamsForOrg)

			value.Teams = teams

			if ok {
				userOrgs[orgName] = value
			} else {
				log.Fatalln("Error saving selected organisation.")
			}

			addOrgs = false

			red.Println("All available organisations have been selected.")
			break
		} else if len(s) == 0 && len(teamsForOrg) == 0 {
			red.Println("No teams available for this organisation.")
			red.Println("All available organisations have been selected.")
			addOrgs = false
			break
		} else if len(s) != 0 && len(teamsForOrg) == 0 {
			red.Println("No teams available for this organisation.")
			result := PromptAnotherOrg()

			switch result {
			case "yes":
				continue
			case "no":
				addOrgs = false
				break
			default:
				log.Fatalln("Prompt for organisations failed.")
			}
		} else {
			log.Fatalln("Error selecting teams for organisations.")
		}
	}
	return userOrgs
}

// DeleteFlagOrgs deletes organisation from prompt if it has already been set from flag
func DeleteFlagOrgs(flagOrgs map[string]Org) map[string]Org {
	allOrgs := GetAllOrgs()

	for k := range flagOrgs {
		_, ok := allOrgs[k]
		if ok {
			delete(allOrgs, k)
		}
	}
	return allOrgs
}

// CreateOrgList puts org names into list to be prompted
func CreateOrgList(allOrgs map[string]Org) []string {
	s := []string{}
	for k := range allOrgs {
		s = append(s, k)
		sort.Strings(s)
	}
	return s
}

// GetAllOrgs calls Github API to receive currently available orgs with their IDs
func GetAllOrgs() map[string]Org {
	allOrgs := make(map[string]Org)

	orgs, _, err := Client.Organizations.ListAll(context.Background(), nil)

	if err != nil {
		log.Fatalln("Error getting list of organisations from Github:", err)
	}

	for _, v := range orgs {
		name := v.GetLogin()
		// id := v.GetID()
		allOrgs[name] = Org{
			Name: v.GetLogin(),
			ID:   v.GetID(),
		}
	}
	return allOrgs
}

// PromptOrgNames prompts user to choose available orgs from a list
func PromptOrgNames(s []string) string {
	prompt := promptui.Select{
		Label: "Select Organisation(s)",
		Items: s,
	}
	_, org, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	OrgCount++

	GreenBold.Print("Organisation ", OrgCount, " ")
	fmt.Println(org)

	return org
}

// PromptAnotherOrg asks if user wants to add another org
func PromptAnotherOrg() string {
	prompt := promptui.Select{
		Label: "Add another organisation?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// PromptRole asks user to set the role of new user in chosen org
func PromptRole() string {
	prompt := promptui.Select{
		Label: "Role?",
		Items: []string{"member", "admin"},
	}
	_, role, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	GreenBold.Print("Role ")
	fmt.Println(role)
	return role
}

// GetTeamsForOrg returns a map of the teams in given org
func GetTeamsForOrg(org string) map[string]Team {
	teamsForOrg := make(map[string]Team)

	teams, _, err := Client.Teams.ListTeams(context.Background(), org, nil)
	if err != nil {
		log.Fatalln("Error getting list of teams from Github:", err)
	}

	for _, v := range teams {
		name := v.GetName()
		id := v.GetID()
		teamsForOrg[name] = Team{name, id}
	}
	return teamsForOrg
}

// CreateTeamList creates list with available Team names for prompt
func CreateTeamList(teamsForOrg map[string]Team) []string {
	s := []string{}
	for k := range teamsForOrg {
		s = append(s, k)
		sort.Strings(s)
	}
	return s
}

// SelectTeams prompts the user to select teams for the current org
func SelectTeams(teamsForOrg map[string]Team) []Team {
	userTeams := []Team{}
	addTeams := true
	s := CreateTeamList(teamsForOrg)

	if s == nil {
		addTeams = false
	}

	// Keep prompting for teams while true
	for addTeams == true {
		teamName := PromptTeams(s)

		value, ok := teamsForOrg[teamName]
		if ok {
			userTeams = append(userTeams, value)
		}

		// Deletes selected teams from list
		for i, v := range s {
			if v == teamName {
				s = append(s[:i], s[i+1:]...)
				break
			}
		}

		if len(s) == 0 {
			red := color.New(color.FgRed)
			red.Println("All available teams have been selected.")
			addTeams = false
			break
		}

		result := PromptAnotherTeam()

		switch result {
		case "yes":
			continue
		case "no":
			addTeams = false
			break
		default:
			log.Fatalln("Prompt for teams failed.")
		}
	}
	return userTeams
}

// PromptAnotherTeam asks user if they want to add another team to current org
func PromptAnotherTeam() string {
	prompt := promptui.Select{
		Label: "Add another team?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// PromptTeams offers user available teams to select
func PromptTeams(s []string) string {
	prompt := promptui.Select{
		Label: "Select Team(s)",
		Items: s,
	}

	_, team, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	GreenBold.Print("Team ")
	fmt.Println(team)

	return team
}

type User struct {
	Login string `json:"login"`
	Email string `json:"email"`
	ID    int64  `json:"id"`
}

func PromptAddUserToOrg() string {
	prompt := promptui.Select{
		Label: "Add user to Github org/ teams?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// CreateUser a new github user with the username and email provided
func CreateUser(username string, email string) (string, int64) {
	body := User{Login: username, Email: email}

	req, err := Client.NewRequest("POST", "https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/admin/users", body)
	if err != nil {
		log.Fatalln("Error creating new request:\n", err)
	}

	newUser := User{}
	_, err = Client.Do(context.Background(), req, &newUser)

	if err != nil {
		log.Fatalf("Error creating new user:\n%v", err)
	}

	username = newUser.Login
	userID := newUser.ID

	fmt.Println("")
	WhiteBold.Printf("User '%v' created.", username)
	fmt.Println("")

	return username, userID
}

// AddUserToOrgs invites new user to selected github orgs and teams
func AddUserToOrgs(username string, orgs map[string]Org) {
	for _, v := range orgs {
		orgName := v.Name
		role := v.Role

		membershipOptions := &github.Membership{
			Role: github.String(role),
		}

		_, resp, err := Client.Organizations.EditOrgMembership(context.Background(), username, orgName, membershipOptions)
		if err != nil {
			fmt.Println(resp.Status)
			log.Fatalf("Error creating adding '%v' to Github organisation: %v", username, err)
		}
		WhiteBold.Printf("User '%v' added to Github organisation '%v' as '%v'.", username, orgName, role)
		fmt.Println("")

		teams := v.Teams
		for _, vv := range teams {
			teamName := vv.Name
			teamID := vv.ID
			AddUserToTeams(username, teamID, teamName)
		}
	}
}

// AddUserToTeams adds user to selected teams within chosen org
func AddUserToTeams(username string, teamID int64, teamName string) {
	_, resp, err := Client.Teams.AddTeamMembership(context.Background(), teamID, username, nil)
	if err != nil {
		fmt.Println(resp.Status)
		log.Fatalf("Error adding '%v' to Github team: %v", username, err)
	}

	WhiteBold.Printf("User '%v' added to Github team '%v'.", username, teamName)
	fmt.Println("")
}
