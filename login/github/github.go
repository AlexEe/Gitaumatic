package github

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	ini "gopkg.in/ini.v1"
)

// GetGithubTokens logs into Github with admin priviledges and retrieves Personal Access
// Check how to pass in structs
func GetGithubTokens() (string, string, string, int64, string) {
	cfg, err := ini.Load("/Users/alex/go/src/omniactl/.fake_vault")
	if err != nil {
		log.Fatalln("Failure retrieving tokens from Vault:", err)
	}
	username := cfg.Section("auth").Key("github_username").String()
	password := cfg.Section("auth").Key("github_password").String()
	token := cfg.Section("auth").Key("github_token").String()
	teamID, _ := cfg.Section("auth").Key("github_team").Int64()
	team := int64(teamID)

	cfg2, err := ini.Load("/Users/alex/go/src/omniactl/.omniactl")
	if err != nil {
		log.Fatalln("Failure retrieving URLs from .omniactl config file."+
			"\nThis file can be created or updated using the './omniactl config' command.", err)
	}
	address := fmt.Sprint(cfg2.Section("config").Key("github").String())

	return username, password, token, team, address
}

// GithubLogin logs user into github
func GithubLogin(s string) {
	greenBold := color.New(color.FgGreen, color.Bold)
	username, _, token, _, address := GetGithubTokens()
	// create authenticated Github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client, err := github.NewEnterpriseClient(address, address, tc)
	if err != nil {
		log.Fatalln("Creation of Github client failed:", err)
	}

	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Fatalln("Github login failed:", err)
	}

	if resp.StatusCode == 200 && s == "login" {
		fmt.Printf("Successfully logged into %v as user %v.\n", address, username)
	} else if resp.StatusCode == 200 && s == "check" {
		greenBold.Print("Github Login Status ")
		fmt.Println("OK")
	}
}

// CheckGithubLogin checks if Github returns data for authenticated user to verify login
func CheckGithubLogin() error {
	client := CreateClient()
	ctx := context.Background()
	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		err := errors.New("Github login failed: Could not create Github client")
		return err
	}

	_, resp, err = client.Users.Get(ctx, "")
	if err != nil {
		err := errors.New("Github login failed: Error retrieving current user information from Github")
		return err
	}

	if resp.StatusCode != 200 {
		err := errors.New("Github login failed: Error retrieving current user information from Github")
		return err
	}
	return nil
}

// CreateClient creates a client for interaction with github, authorized using token
func CreateClient() *github.Client {
	_, _, token, _, address := GetGithubTokens()
	// create authenticated Github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client, err := github.NewEnterpriseClient(address, address, tc)
	if err != nil {
		log.Fatalln("Creation of Github client failed:", err)
	}
	return client
}
