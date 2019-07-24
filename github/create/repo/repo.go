package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/manifoldco/promptui"
	"log"
	createOrg "omniactl/github/create/org"
	createTeam "omniactl/github/create/team"
	createUser "omniactl/github/create/user"
	listTeam "omniactl/github/list/team"
	githubLogin "omniactl/login/github"
	"os"
	"regexp"
	"strings"
)

func CreateRepo(name string, org string, team string, description string, privacy bool) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Create a new Github repository")
	var collaborators []string

	org = CheckOrgFlag(org)
	teamMap := CheckTeamFlag(team, org)
	name = CheckNameFlag(name)
	url, repoName := CreateGithubRepo(name, org, teamMap, description, privacy)
	result := PromptCollaborators()
	switch result {
	case "Add all org members":
		collaborators = AddAllOrgMembers(org)
		AddCollaborators(collaborators, url, repoName)
	case "Add all team members":
		collaborators = AddAllTeamMembers(teamMap)
		AddCollaborators(collaborators, url, repoName)
	case "Add a specific user":
		collaborators = AddSpecificUsers()
		AddCollaborators(collaborators, url, repoName)
	default:
		os.Exit(0)
	}
}

func AddAllOrgMembers(org string) []string {
	Client := githubLogin.CreateClient()
	var Collaborators []string

	members, _, err := Client.Organizations.ListMembers(context.Background(), org, nil)
	if err != nil {
		log.Fatalln("Error getting info about organisation:", err)
	}

	for _, v := range members {
		memberLogin := v.GetLogin()
		Collaborators = append(Collaborators, memberLogin)
	}
	return Collaborators
}

func AddAllTeamMembers(teamMap map[string]createUser.Team) []string {
	Client := githubLogin.CreateClient()
	var Collaborators []string
	var TeamID int64

	for _, v := range teamMap {
		TeamID = v.ID
	}

	members, _, err := Client.Teams.ListTeamMembers(context.Background(), TeamID, nil)
	if err != nil {
		log.Fatalln("Error getting info about team members:", err)
	}

	for _, v := range members {
		memberLogin := v.GetLogin()
		Collaborators = append(Collaborators, memberLogin)
	}
	return Collaborators
}

func AddCollaborators(collaborators []string, urlRepo string, repoName string) {
	Client := githubLogin.CreateClient()
	whiteBold := color.New(color.FgHiWhite, color.Bold)

	fmt.Println("")
	for _, userLogin := range collaborators {
		urlUser := fmt.Sprintf("/collaborators/%v", userLogin)
		var urlSlice []string
		urlSlice = append(urlSlice, urlRepo)
		urlSlice = append(urlSlice, urlUser)
		urlComplete := strings.Join(urlSlice, "")

		req, err := Client.NewRequest("PUT", urlComplete, nil)
		if err != nil {
			log.Fatalln("Error creating new request:\n", err)
		}

		_, err = Client.Do(context.Background(), req, nil)
		if err != nil {
			log.Fatalln("Error adding collaborators to repository: ", err)
		}

		whiteBold.Printf("User '%v' added as collaborator to repository '%v'.", userLogin, repoName)
		fmt.Println("")
	}
}

func AddSpecificUsers() []string {
	var AddCollaborator bool
	var Collaborators []string

	AddCollaborator = true

	for AddCollaborator == true {
		userLogin := PromptUsername()
		Collaborators = append(Collaborators, userLogin)

		result := PromptAnotherCollaborator()
		switch result {
		case "yes":
			continue
		default:
			AddCollaborator = false
			break
		}
	}
	return Collaborators
}

func PromptUsername() string {
	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^e[0-9]{6}$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Username must be in the format of a State Street Lan ID, e.g. 'e123456'")
		}
		check := createUser.CheckIfUserExists(input)
		switch check {
		case true:
			return nil
		default:
			red := color.New(color.FgRed)
			red.Println("Username does not exist")
			return errors.New("Username does not exist")
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

func PromptAnotherCollaborator() string {
	prompt := promptui.Select{
		Label: "Invite another collaborator?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func PromptCollaborators() string {
	prompt := promptui.Select{
		Label: "Invite collaborators?",
		Items: []string{"Add all org members", "Add all team members", "Add a specific user", "Do not add any collaborators"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func CreateGithubRepo(name string, org string, teamMap map[string]createUser.Team, description string, privacy bool) (string, string) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	var TeamID int64
	Client := githubLogin.CreateClient()

	for _, v := range teamMap {
		TeamID = v.ID
	}

	// Define new repository
	repo := &github.Repository{
		Name:             github.String(name),
		Description:      github.String(description),
		Private:          github.Bool(false),
		TeamID:           github.Int64(TeamID),
		AllowRebaseMerge: github.Bool(true),
		AllowSquashMerge: github.Bool(true),
		AllowMergeCommit: github.Bool(true),
	}

	// Create repo inside specific organization
	repo, _, err := Client.Repositories.Create(context.Background(), org, repo)
	if err != nil {
		log.Fatalln("Error creating repo:", err)
	}

	repoName := repo.GetName()

	fmt.Println("")
	whiteBold.Printf("Repository '%v' has been created.\n", repoName)
	// fmt.Println("ID: ", repo.GetID())
	// permissions, _, _ := Client.Repositories.GetPermissionLevel(context.Background(), org, repo.GetName(), "e111111")
	// fmt.Println("Permissions: ", permissions.GetPermission())
	fmt.Print("(Go to repo: ")
	fmt.Printf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/%v/%v)", org, repoName)
	fmt.Println("")
	fmt.Println("")

	url := repo.GetURL()
	return url, repoName
}

func CheckNameFlag(name string) string {
	greenBold := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed)
	if name != "" {
		isCorrectFormat := regexp.MustCompile("^[a-z0-9._%+\\-]+$").MatchString
		if isCorrectFormat(name) != true {
			red.Println("Repo name can only contain the following characters: A-Z, a-z, 0-9, -, _")
			name = PromptName()
			return name
		} else {
			greenBold.Print("Repo name ")
			fmt.Println(name)
			return name
		}
	} else {
		name = PromptName()
		return name
	}
}

func PromptName() string {
	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^[a-zA-Z0-9._%+\\-]+$").MatchString
		if isCorrectFormat(input) != true {
			return errors.New("Repo name can only contain the following characters: A-Z, a-z, 0-9, -, _")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Repo name",
		Validate:  validate,
		Templates: templates,
	}
	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func CheckOrgFlag(org string) string {
	red := color.New(color.FgRed)
	greenBold := color.New(color.FgGreen, color.Bold)

	if org != "" {
		check := createOrg.CheckIfOrgExists(org)
		switch check {
		case true:
			greenBold.Print("Organisation ")
			fmt.Println(org)
			return org
		default:
			red.Printf("Organisation '%v' does not exist.", org)
			fmt.Println("")
			org = createTeam.PromptOrg()
			return org
		}
	} else {
		check := CheckIfUserRepo()
		switch check {
		case "User":
			org = ""
			return org
		default:
			org = createTeam.PromptOrg()
			return org
		}
	}
}

func CheckTeamFlag(team string, org string) map[string]createUser.Team {
	teamMap := make(map[string]createUser.Team)
	red := color.New(color.FgRed)
	greenBold := color.New(color.FgGreen, color.Bold)

	if org == "" {
		return teamMap
	}

	if team == "" {
		check := CheckIfTeamRepo()
		switch check {
		case "Team":
			teamMap = listTeam.PromptTeam(org)
			return teamMap
		default:
			return teamMap
		}
	} else {
		teamsForOrg := createUser.GetTeamsForOrg(org)
		check := createTeam.CheckIfTeamExists(team, teamsForOrg)
		switch check {
		case true:
			greenBold.Print("Team ")
			fmt.Println(team)
			teamMap = listTeam.CreateTeamMap(team, org)
			return teamMap
		default:
			red.Printf("Team '%v' does not exist.", team)
			fmt.Println("")
			teamMap = listTeam.PromptTeam(org)
			return teamMap
		}
	}
}

func CheckIfTeamRepo() string {
	prompt := promptui.Select{
		Label: "Create repo within the organisation itself or within a team?",
		Items: []string{"Organisation", "Team"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func CheckIfUserRepo() string {
	prompt := promptui.Select{
		Label: "Create repo for authenticated user or within an organisation/team?",
		Items: []string{"Organisation/team", "User"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     string `json:"private"`
	TeamID      int64  `json:"team_id"`
	Permissions struct {
		Admin bool `json:"admin"`
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
	} `json:"permissions"`
}

type RepoResponse struct {
	AllowMergeCommit bool        `json:"allow_merge_commit"`
	AllowRebaseMerge bool        `json:"allow_rebase_merge"`
	AllowSquashMerge bool        `json:"allow_squash_merge"`
	ArchiveURL       string      `json:"archive_url"`
	Archived         bool        `json:"archived"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CloneURL         string      `json:"clone_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	CreatedAt        string      `json:"created_at"`
	DefaultBranch    string      `json:"default_branch"`
	DeploymentsURL   string      `json:"deployments_url"`
	Description      string      `json:"description"`
	Disabled         bool        `json:"disabled"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	Fork             bool        `json:"fork"`
	ForksCount       int         `json:"forks_count"`
	ForksURL         string      `json:"forks_url"`
	FullName         string      `json:"full_name"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	HasDownloads     bool        `json:"has_downloads"`
	HasIssues        bool        `json:"has_issues"`
	HasPages         bool        `json:"has_pages"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	Homepage         string      `json:"homepage"`
	HooksURL         string      `json:"hooks_url"`
	HTMLURL          string      `json:"html_url"`
	ID               int         `json:"id"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	Language         interface{} `json:"language"`
	LanguagesURL     string      `json:"languages_url"`
	MergesURL        string      `json:"merges_url"`
	MilestonesURL    string      `json:"milestones_url"`
	MirrorURL        string      `json:"mirror_url"`
	Name             string      `json:"name"`
	NetworkCount     int         `json:"network_count"`
	NodeID           string      `json:"node_id"`
	NotificationsURL string      `json:"notifications_url"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Owner            struct {
		AvatarURL         string `json:"avatar_url"`
		EventsURL         string `json:"events_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		GravatarID        string `json:"gravatar_id"`
		HTMLURL           string `json:"html_url"`
		ID                int    `json:"id"`
		Login             string `json:"login"`
		NodeID            string `json:"node_id"`
		OrganizationsURL  string `json:"organizations_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		ReposURL          string `json:"repos_url"`
		SiteAdmin         bool   `json:"site_admin"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		Type              string `json:"type"`
		URL               string `json:"url"`
	} `json:"owner"`
	Permissions struct {
		Admin bool `json:"admin"`
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
	} `json:"permissions"`
	Private          bool     `json:"private"`
	PullsURL         string   `json:"pulls_url"`
	PushedAt         string   `json:"pushed_at"`
	ReleasesURL      string   `json:"releases_url"`
	Size             int      `json:"size"`
	SSHURL           string   `json:"ssh_url"`
	StargazersCount  int      `json:"stargazers_count"`
	StargazersURL    string   `json:"stargazers_url"`
	StatusesURL      string   `json:"statuses_url"`
	SubscribersCount int      `json:"subscribers_count"`
	SubscribersURL   string   `json:"subscribers_url"`
	SubscriptionURL  string   `json:"subscription_url"`
	SvnURL           string   `json:"svn_url"`
	TagsURL          string   `json:"tags_url"`
	TeamsURL         string   `json:"teams_url"`
	Topics           []string `json:"topics"`
	TreesURL         string   `json:"trees_url"`
	UpdatedAt        string   `json:"updated_at"`
	URL              string   `json:"url"`
	WatchersCount    int      `json:"watchers_count"`
}
