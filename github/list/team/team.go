package team

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"log"
	createOrg "omniactl/github/create/org"
	createTeam "omniactl/github/create/team"
	createUser "omniactl/github/create/user"
	githubLogin "omniactl/login/github"
	"time"
	// "github.com/fatih/color"
	// "errors"
	// "regexp"
	"context"
	"sort"
)

// ListTeam receives flag input and shows structure of package
func ListTeam(team string, org string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: List information about a Github team")

	githubTeam := CheckFlag(team, org)
	ListGithubTeam(githubTeam)
}

// CheckFlag checks if input was put via flags, checks input or prompts user
func CheckFlag(team string, org string) map[string]createUser.Team {
	greenBold := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed)
	githubTeam := make(map[string]createUser.Team)
	if org != "" {
		check := createOrg.CheckIfOrgExists(org)
		switch check {
		case true:
			greenBold.Print("Organisation ")
			fmt.Println(org)
			if team == "" {
				githubTeam = PromptTeam(org)
				return githubTeam
			} else {
				teamsForOrg := createUser.GetTeamsForOrg(org)
				check := createTeam.CheckIfTeamExists(team, teamsForOrg)
				switch check {
				case true:
					greenBold.Print("Team ")
					fmt.Println(team)
					githubTeam = CreateTeamMap(team, org)
					return githubTeam
				default:
					red.Printf("Team '%v' does not exist.", team)
					fmt.Println("")
					githubTeam = PromptTeam(org)
					return githubTeam
				}
			}
		default:
			red.Printf("Organisation '%v' does not exist.", org)
			fmt.Println("")
			org = createTeam.PromptOrg()
			githubTeam = PromptTeam(org)
			return githubTeam
		}
	} else {
		org = createTeam.PromptOrg()
		githubTeam = PromptTeam(org)
		return githubTeam
	}
}

// ListGithubTeam prints out information on chosen github team
func ListGithubTeam(githubTeam map[string]createUser.Team) {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	greenBold := color.New(color.FgGreen, color.Bold)
	Client := githubLogin.CreateClient()

	for k, v := range githubTeam {
		teamID := v.ID
		team, _, err := Client.Teams.GetTeam(context.Background(), teamID)
		if err != nil {
			log.Fatalf("Error getting information about Github team '%v': %v", k, err)
		}

		result := PromptAction()

		switch result {
		case "Team stats":
			fmt.Println("")
			whiteBold.Println("Team stats:")
			greenBold.Print("Name ")
			fmt.Println(team.GetName())
			greenBold.Print("ID ")
			fmt.Println(team.GetID())
			greenBold.Print("Description ")
			fmt.Println(team.GetDescription())
			greenBold.Print("Organisation ")
			fmt.Println(team.GetOrganization().GetLogin())
			greenBold.Print("Permission ")
			fmt.Println(team.GetPermission())
			greenBold.Print("Privacy ")
			fmt.Println(team.GetPrivacy())
			greenBold.Print("No of repos ")
			fmt.Println(team.GetReposCount())
		case "Team repositories":
			fmt.Println("")
			whiteBold.Println("Team repositories:")
			GetRepos(teamID)
		case "Team members":
			fmt.Println("")
			whiteBold.Println("Team members:")
			members, _, err := Client.Teams.ListTeamMembers(context.Background(), teamID, nil)
			if err != nil {
				log.Fatalln("Error getting info about team members:", err)
			}

			for _, v := range members {
				user, _, err := Client.Users.Get(context.Background(), v.GetLogin())
				if err != nil {
					log.Fatalln("Error getting info about organisation members:", err)
				}

				fmt.Printf("Name: %-25v | Login: %-20v | ID: %-15v | Site admin: %-10v\n", user.GetName(), v.GetLogin(), v.GetID(), v.GetSiteAdmin())
			}
		default:
			log.Fatalln("Error selecting action.")
		}
	}
}

// PromptAction asks user to select which info they want about selected team
func PromptAction() string {
	prompt := promptui.Select{
		Label: "Select action",
		Items: []string{"Team stats", "Team members", "Team repositories"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}

// Repos struct created to contain info on team repos received from HTTP request
type Repos []struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	Private          bool        `json:"private"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ArchiveURL       string      `json:"archive_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	ForksURL         string      `json:"forks_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	LanguagesURL     string      `json:"languages_url"`
	MergesURL        string      `json:"merges_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	PullsURL         string      `json:"pulls_url"`
	ReleasesURL      string      `json:"releases_url"`
	SSHURL           string      `json:"ssh_url"`
	StargazersURL    string      `json:"stargazers_url"`
	StatusesURL      string      `json:"statuses_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	TagsURL          string      `json:"tags_url"`
	TeamsURL         string      `json:"teams_url"`
	TreesURL         string      `json:"trees_url"`
	CloneURL         string      `json:"clone_url"`
	MirrorURL        string      `json:"mirror_url"`
	HooksURL         string      `json:"hooks_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Language         interface{} `json:"language"`
	ForksCount       int         `json:"forks_count"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Size             int         `json:"size"`
	DefaultBranch    string      `json:"default_branch"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Topics           []string    `json:"topics"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	HasDownloads     bool        `json:"has_downloads"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	PushedAt         time.Time   `json:"pushed_at"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Permissions      struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	SubscribersCount int `json:"subscribers_count"`
	NetworkCount     int `json:"network_count"`
	License          struct {
		Key    string `json:"key"`
		Name   string `json:"name"`
		SpdxID string `json:"spdx_id"`
		URL    string `json:"url"`
		NodeID string `json:"node_id"`
	} `json:"license"`
}

// GetRepos sends an HTTP request to get names and stats of team repos
func GetRepos(teamID int64) {
	Client := githubLogin.CreateClient()
	PageCount := 1
	NextPage := true
	RepoCount := 1

	for NextPage == true {
		url := fmt.Sprintf("https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/teams/%v/repos?page=%v&per_page=100", teamID, PageCount)

		req, err := Client.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln("Error creating HTTP request:\n", err)
		}

		// Create repo struct to save data of HTTP request
		teamRepos := Repos{}
		_, err = Client.Do(context.Background(), req, &teamRepos)
		if err != nil {
			log.Fatalln("Error retrieving repos for team", err)
		}

		PageCount++

		switch len(teamRepos) {
		case 0:
			NextPage = false
		default:
			for _, v := range teamRepos {
				fmt.Printf("%-5v Name: %-40v | ID: %-15v | Updated at: %-40v | Private: %-20v", RepoCount, v.Name, v.ID, v.UpdatedAt, v.Private)
				fmt.Println("")
				RepoCount++
			}
		}
	}
}

// CreateTeamMap takes team name string and returns a map containing the team
func CreateTeamMap(team string, org string) map[string]createUser.Team {
	teamsForOrg := createUser.GetTeamsForOrg(org)
	githubTeam := make(map[string]createUser.Team)

	// Checks if team name is a key in available teams,
	// if so fills new team map with corresponding value
	value, ok := teamsForOrg[team]
	if ok {
		githubTeam[team] = value
	}
	return githubTeam
}

// PromptTeam prompts user to select team for which to provide info
func PromptTeam(org string) map[string]createUser.Team {
	// Creates list with available Team names for prompt
	githubTeam := make(map[string]createUser.Team)
	teamsForOrg := createUser.GetTeamsForOrg(org)

	s := []string{}
	for k := range teamsForOrg {
		s = append(s, k)
		sort.Strings(s)
	}

	prompt := promptui.Select{
		Label: "Select Team",
		Items: s,
	}

	_, team, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	githubTeam = CreateTeamMap(team, org)

	greenBold := color.New(color.FgGreen, color.Bold)
	greenBold.Print("Team ")
	fmt.Println(team)

	return githubTeam
}
