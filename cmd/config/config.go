package config

import (
	"github.com/spf13/cobra"
	listConfig "omniactl/config/list"
	updateConfig "omniactl/config/update"
)

var (
	github      string
	jira        string
	confluence  string
	artifactory string
	concourse   string
	vault       string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config command requires an additional subcommand to be executed, e.g. update.",
	Long:  "Config command requires an additional subcommand to be executed, e.g. update.",
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "'config update' command allows to update config settings",
	Long: `Config command allows user to list/update config settings for
	Github, Jira, Confluence, Artifactory, Vault and Concourse.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateConfig.UpdateConfigFile(github, jira, confluence, artifactory, concourse, vault)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "'config list' command displays current config file",
	Long:  "'config list' displays URL endpoints for Github, Jira, Confluence, Concourse, Artifactory and Vault",
	Run: func(cmd *cobra.Command, args []string) {
		listConfig.ListConfigFile()
	},
}

func init() {
	configCmd.AddCommand(updateCmd)
	configCmd.AddCommand(listCmd)
	updateCmd.Flags().StringVarP(&github, "github", "g", "", "Github URL")
	updateCmd.Flags().StringVarP(&jira, "jira", "j", "", "Jira URL")
	updateCmd.Flags().StringVarP(&confluence, "confluence", "f", "", "Confluence URL")
	updateCmd.Flags().StringVarP(&concourse, "concourse", "c", "", "Concourse ULR")
	updateCmd.Flags().StringVarP(&vault, "vault", "v", "", "Vault URL")
}

// AddSubCommands adds the sub-commands to the provided command
func AddSubCommands(cmd *cobra.Command) {
	cmd.AddCommand(configCmd)
}
