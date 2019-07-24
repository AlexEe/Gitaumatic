package jira

import (

	// createProject "omniactl/jira/create/project"
	// createRole "omniactl/jira/create/role"
	// createUser "omniactl/jira/create/user"
	// listProject "omniactl/jira/list/project"
	// listRole "omniactl/jira/list/role"
	// listUser "omniactl/jira/list/user"
	// suspendUser "omniactl/jira/suspend/user"
	jiraApi "omniactl/jira"

	"github.com/spf13/cobra"
)

var (
	username        string
	name            string
	email           string
	jiraProjectName string
	usernameSuspend string
	usernameList    string
	projectNameList string
	reasonSuspend   string
	roleName        string
)

// githubCmd represents the github command
var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Subcommand for interacting with JIRA API.",
	Long: "omniactl jira' command allows for interacting with the JIRA API." +
		"For instance, run the subcommand 'omniactl jira create user' to add a new user to JIRA.",
}

// githubCmd represents the github command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Subcommand for interacting with JIRA API.",
	Long: "'omniactl jira create' command allows for the following actions, depending on the chosen subcommand:" +
		"1, 'user': add a new user to JIRA." +
		"2, 'project': add a new project to JIRA." +
		"3, 'role': add a new role to JIRA.",
}

var userCreateCmd = &cobra.Command{
	Use:   "user",
	Short: "Add a new user to JIRA.",
	Long:  "'user' subcommand requires username, name, email address and project name to create new JIRA user.",
	Run: func(cmd *cobra.Command, args []string) {
		jiraApi.AddUser(username, name, email, jiraProjectName)
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "org",
	Short: "Creates a new JIRA project.",
	Long:  "Creates a project with a provided name.",
	Run: func(cmd *cobra.Command, args []string) {
		jiraApi.CreateProject(jiraProjectName)
	},
}

var roleCreateCmd = &cobra.Command{
	Use:   "role",
	Short: "Add role to JIRA.",
	Long:  "Add a role to JIRA with the provided name",
	Run: func(cmd *cobra.Command, args []string) {
		jiraApi.CreateRole(roleName)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Subcommand for interacting with JIRA API.",
	Long:  "'list' requires a subcommand, e.g. 'user', to be executed.",
}

var userListCmd = &cobra.Command{
	Use:   "user",
	Short: "List information for a user or users.",
	Long:  "Lists information about a user or users, including username, login, email, organisations, teams, role etc.",
	Run: func(cmd *cobra.Command, args []string) {
		jiraApi.ListUser(usernameList)
	},
}

var projectListCmd = &cobra.Command{
	Use:   "project",
	Short: "List projects in JIRA.",
	Long:  "Lists the projects in JIRA",
	Run: func(cmd *cobra.Command, args []string) {
		jiraApi.ListProject(projectNameList)
	},
}

var roleListCmd = &cobra.Command{
	Use:   "role",
	Short: "List roles in JIRA.",
	Long:  "List the roles in JIRA.",
	Run: func(cmd *cobra.Command, args []string) {
		jiraApi.ListRole(roleName)
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
		jiraApi.SuspendUser(usernameSuspend, reasonSuspend)
	},
}

func init() {
	// jira create
	jiraCmd.AddCommand(createCmd)
	createCmd.AddCommand(userCreateCmd)
	createCmd.AddCommand(projectCreateCmd)
	createCmd.AddCommand(roleCreateCmd)
	// jira list
	jiraCmd.AddCommand(listCmd)
	listCmd.AddCommand(userListCmd)
	listCmd.AddCommand(projectListCmd)
	listCmd.AddCommand(roleListCmd)
	// jira suspend
	jiraCmd.AddCommand(suspendCmd)
	suspendCmd.AddCommand(userSuspendCmd)
	// jira suspend user --username --reason
	//userSuspendCmd.MarkFlagRequired("username")
	//userSuspendCmd.MarkFlagRequired("reason")
	// jira list user
	//userListCmd.MarkFlagRequired("username")

	userCreateCmd.Flags().StringVarP(&username, "username", "u", "", "JIRA username is State Street Lan ID (required)")
	userCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Full name (required)")
	userCreateCmd.Flags().StringVarP(&email, "email", "e", "", "Email is State Street email (required)")
	userCreateCmd.Flags().StringVarP(&jiraProjectName, "jira-project-name", "j", "", "JIRA project name (required)")
	projectCreateCmd.Flags().StringVarP(&jiraProjectName, "jira-project-name", "j", "", "JIRA project name (required)")
	roleCreateCmd.Flags().StringVarP(&roleName, "role", "r", "", "The role name (required)")

	userSuspendCmd.Flags().StringVarP(&usernameSuspend, "username", "u", "", "Username is State Street Lan ID of user to be suspended (required)")
	userSuspendCmd.Flags().StringVarP(&reasonSuspend, "reason", "r", "", "Reason the user is being suspended")
	userListCmd.Flags().StringVarP(&usernameList, "username", "u", "", "Username is State Street Lan ID of user to list (required)")
	projectListCmd.Flags().StringVarP(&jiraProjectName, "jira-project-name", "j", "", "JIRA project name (required)")
	roleListCmd.Flags().StringVarP(&roleName, "role", "r", "", "Full or partial role name")
}

// AddSubCommands adds the sub-commands to the provided command
func AddSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(jiraCmd)
}
