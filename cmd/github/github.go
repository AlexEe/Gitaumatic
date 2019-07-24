package github

import (
	createOrg "omniactl/github/create/org"
	createRepo "omniactl/github/create/repo"
	createTeam "omniactl/github/create/team"
	createUser "omniactl/github/create/user"
	deleteUser "omniactl/github/delete/user"
	listOrg "omniactl/github/list/org"
	listOrgs "omniactl/github/list/orgs"
	listTeam "omniactl/github/list/team"
	listTeams "omniactl/github/list/teams"
	listUser "omniactl/github/list/user"
	listUsers "omniactl/github/list/users"
	suspendUser "omniactl/github/suspend/user"
	updateUser "omniactl/github/update/user"

	"github.com/spf13/cobra"
)

var (
	username        string
	usernameSuspend string
	usernameDelete  string
	usernameList    string
	usernamesList   []string
	reasonSuspend   string
	usernameUpdate  string
	email           string
	org             string
	role            string
	teams           []string
	orgName         string
	orgProfile      string
	orgAdmin        string
	orgTeam         string
	team            string
	teamMaintainers []string
	teamPrivacy     string
	teamDescription string
	orgList         string
	teamList        string
	orgTeamList     string
	teamsList       string
	orgTeamsList    string
	repoName        string
	repoOrg         string
	repoTeam        string
	repoPrivacy     bool
	repoDescription string
)

// githubCmd represents the github command
var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Subcommand for interacting with Github API.",
	Long: "omniactl github' command allows for interacting with the Github API." +
		"For instance, run the subcommand 'omniactl github adduser' to add a new user to Github.",
}

// githubCmd represents the github command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Subcommand for interacting with Github API.",
	Long: "'omniactl github create' command allows for the following actions, depending on the chosen subcommand:" +
		"1, 'user': add a new user to Github." +
		"2, 'org': add a new organisation to Github." +
		"3, 'team': add a new team to Github.",
}

var userCreateCmd = &cobra.Command{
	Use:   "user",
	Short: "Add a new user to Github.",
	Long:  "'user' subcommand requires username and email address, optionally also: organisations and teams to create new Github user.",
	Run: func(cmd *cobra.Command, args []string) {
		createUser.AddUser(username, email, org, role, teams)
	},
}

var suspendCmd = &cobra.Command{
	Use:   "suspend",
	Short: "Subcommand for interacting with Github API.",
	Long:  "'suspend' requires a subcommand, e.g. 'user', to be executed.",
}

var userSuspendCmd = &cobra.Command{
	Use:   "user",
	Short: "Suspend a user from Github.",
	Long:  "'suspend user' subcommand requires the username/Lan ID of the user to be suspended as well as a reason for the suspension.",
	Run: func(cmd *cobra.Command, args []string) {
		suspendUser.SuspendUser(usernameSuspend, reasonSuspend)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Subcommand for interacting with Github API.",
	Long:  "'delete' requires a subcommand, e.g. 'user', to be executed.",
}

var userDeleteCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete a user from Github.",
	Long:  "Deleting a user will delete all their repositories, gists, applications, and personal settings. Suspending a user is often a better option.",
	Run: func(cmd *cobra.Command, args []string) {
		deleteUser.DeleteUser(usernameDelete)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Subcommand for interacting with Github API.",
	Long:  "'list' requires a subcommand, e.g. 'user', to be executed.",
}

var userListCmd = &cobra.Command{
	Use:   "user",
	Short: "List information for a user.",
	Long:  "Lists information about the user's username, login, email, organisations, teams, role etc.",
	Run: func(cmd *cobra.Command, args []string) {
		listUser.ListUser(usernameList)
	},
}

var usersListCmd = &cobra.Command{
	Use:   "users",
	Short: "Lists information about multiple users.",
	Long:  "Lists information about the user's username, login, email, organisations, teams, role etc.",
	Run: func(cmd *cobra.Command, args []string) {
		listUsers.ListUsers(usernamesList)
	},
}

var orgCreateCmd = &cobra.Command{
	Use:   "org",
	Short: "Creates a new Github organization.",
	Long:  "Creates an organisation, setting a login/username, profile_name and admin.",
	Run: func(cmd *cobra.Command, args []string) {
		createOrg.CreateOrg(orgName, orgProfile, orgAdmin)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Subcommand for interacting with Github API.",
	Long:  "'update' requires a subcommand, e.g. 'user', to be executed.",
}

var userUpdateCmd = &cobra.Command{
	Use:   "user",
	Short: "Updates an existing Github account.",
	Long:  "Allows adding an existing user to orgs and teams and change their admin status.",
	Run: func(cmd *cobra.Command, args []string) {
		updateUser.UpdateUser(usernameUpdate)
	},
}

var teamCreateCmd = &cobra.Command{
	Use:   "team",
	Short: "Creates a new Github team.",
	Long:  "Creates a new Github team within an existing organisation.",
	Run: func(cmd *cobra.Command, args []string) {
		createTeam.CreateTeam(team, orgTeam, teamDescription, teamMaintainers, teamPrivacy)
	},
}

var orgListCmd = &cobra.Command{
	Use:   "org",
	Short: "Lists information about a Github organization.",
	Long:  "Provides information on a Github organization's members, repos, admins.",
	Run: func(cmd *cobra.Command, args []string) {
		listOrg.ListOrg(orgList)
	},
}

var teamListCmd = &cobra.Command{
	Use:   "team",
	Short: "Lists information about a Github team.",
	Long:  "Provides information on a Github team's members, repos, admins etc.",
	Run: func(cmd *cobra.Command, args []string) {
		listTeam.ListTeam(teamList, orgTeamList)
	},
}

var orgsListCmd = &cobra.Command{
	Use:   "orgs",
	Short: "Lists information about all Github organizations.",
	Long:  "Provides information on a Github organization's members, id, repos etc.",
	Run: func(cmd *cobra.Command, args []string) {
		listOrgs.ListOrgs()
	},
}

var teamsListCmd = &cobra.Command{
	Use:   "teams",
	Short: "Lists information about all Github teams with corresponding orgs.",
	Long:  "Provides information on a Github team's members, repos, ID etc.",
	Run: func(cmd *cobra.Command, args []string) {
		listTeams.ListTeams(orgTeamsList)
	},
}

var createRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Creates a new Github repository",
	Long:  "Creates a new Github repository in a selected org and team with specific permissions, description, privacy etc",
	Run: func(cmd *cobra.Command, args []string) {
		createRepo.CreateRepo(repoName, repoOrg, repoTeam, repoDescription, repoPrivacy)
	},
}

func init() {
	// github create
	githubCmd.AddCommand(createCmd)
	createCmd.AddCommand(userCreateCmd)
	createCmd.AddCommand(orgCreateCmd)
	createCmd.AddCommand(teamCreateCmd)
	createCmd.AddCommand(createRepoCmd)

	// github update
	githubCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(userUpdateCmd)

	// github suspend
	githubCmd.AddCommand(suspendCmd)
	suspendCmd.AddCommand(userSuspendCmd)
	userSuspendCmd.MarkFlagRequired("username")
	userSuspendCmd.MarkFlagRequired("reason")

	// github delete
	githubCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(userDeleteCmd)
	userDeleteCmd.MarkFlagRequired("username")

	// github list
	githubCmd.AddCommand(listCmd)
	listCmd.AddCommand(userListCmd)
	listCmd.AddCommand(usersListCmd)
	listCmd.AddCommand(orgListCmd)
	listCmd.AddCommand(teamListCmd)
	listCmd.AddCommand(orgsListCmd)
	listCmd.AddCommand(teamsListCmd)
	userListCmd.MarkFlagRequired("username")
	usersListCmd.MarkFlagRequired("usernames")

	// flags for commands
	userCreateCmd.Flags().StringVarP(&username, "username", "u", "", "Github username = State Street Lan ID (required)")
	userCreateCmd.Flags().StringVarP(&email, "email", "e", "", "Github email = State Street email (required)")
	userCreateCmd.Flags().StringVarP(&org, "org", "o", "", "Github organisation the new user will be a member of (required)")
	userCreateCmd.Flags().StringVarP(&role, "role", "r", "", "Role the new user will have in the selected organisation, e.g. admin, direct_member (defaults to 'direct_member')")
	userCreateCmd.Flags().StringSliceVarP(&teams, "team", "t", []string{}, "Github teams the new user will be a member of (required)")
	userSuspendCmd.Flags().StringVarP(&usernameSuspend, "username", "u", "", "Username = State Street Lan ID of user to be suspended (required)")
	userSuspendCmd.Flags().StringVarP(&reasonSuspend, "reason", "r", "", "Reason the user is being suspended")
	userDeleteCmd.Flags().StringVarP(&usernameDelete, "username", "u", "", "Username = State Street Lan ID of user to be deleted (required)")
	userListCmd.Flags().StringVarP(&usernameList, "username", "u", "", "Username = State Street Lan ID of user to list (required)")
	usersListCmd.Flags().StringSliceVarP(&usernamesList, "usernames", "u", []string{}, "Usernames = State Street Lan ID of user to list (required)")
	orgCreateCmd.Flags().StringVarP(&orgName, "name", "n", "", "The new organization's username (required)")
	orgCreateCmd.Flags().StringVarP(&orgAdmin, "admin", "a", "", "The new organization's admin (required)")
	orgCreateCmd.Flags().StringVarP(&orgProfile, "profile", "p", "", "The new organization's display/ profile name")
	userUpdateCmd.Flags().StringVarP(&usernameUpdate, "username", "u", "", "Username = State Street Lan ID of user to list (required)")
	teamCreateCmd.Flags().StringVarP(&team, "team", "t", "", "Name of the team to be created (required)")
	teamCreateCmd.Flags().StringVarP(&orgTeam, "org", "o", "", "Existing Github organisation in which the new team will be created (required)")
	teamCreateCmd.Flags().StringVarP(&teamDescription, "description", "d", "", "Description of team to be created")
	teamCreateCmd.Flags().StringSliceVarP(&teamMaintainers, "maintainers", "m", []string{}, "Login names of organization members to add as maintainers of the team")
	teamCreateCmd.Flags().StringVarP(&teamPrivacy, "privacy", "p", "closed", "Level of privacy of the team: secret or closed")
	orgListCmd.Flags().StringVarP(&orgList, "org", "o", "", "Github organisation about which to list information (required)")
	teamListCmd.Flags().StringVarP(&orgTeamList, "org", "o", "", "Github org in which the team resides (required)")
	teamListCmd.Flags().StringVarP(&teamList, "team", "t", "", "Github team about which information is required (required)")
	teamsListCmd.Flags().StringVarP(&orgTeamsList, "org", "o", "", "Github organisation which contains teams to be listed")
	createRepoCmd.Flags().StringVarP(&repoName, "name", "n", "", "Name of new Github repository")
	createRepoCmd.Flags().StringVarP(&repoOrg, "org", "o", "", "Organisation in which new Github repository will be created")
	createRepoCmd.Flags().StringVarP(&repoTeam, "team", "t", "", "Team in which new Github repository will be created")
	createRepoCmd.Flags().BoolVarP(&repoPrivacy, "private", "p", false, "Select 'true' to create a private repo, 'false' to create a public repo")
	createRepoCmd.Flags().StringVarP(&repoDescription, "description", "d", "", "Description of new Github repository")
}

// AddSubCommands adds the sub-commands to the provided command
func AddSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(githubCmd)
}
