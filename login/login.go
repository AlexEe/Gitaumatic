package login

import (
	"bufio"
	"fmt"
	"log"
	githubLogin "omniactl/login/github"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

// Logins sets CLI flags, use `./omniactl login -g` to login to Github
var Logins struct {
	Github      bool `short:"g" long:"--github" description:"Set flag to log into Github" command:"login"`
	Jira        bool `short:"j" long:"--jira" description:"Set flag to log into Jira" command:"login"`
	Confluence  bool `short:"f" long:"--confluence" description:"Set flag to log into Confluence" command:"login"`
	Artifactory bool `short:"a" long:"--artifactory" description:"Set flag to log into Artifactory" command:"login"`
	Concourse   bool `short:"c" long:"--concourse" description:"Set flag to log into Concourse" command:"login"`
	Help        bool `short:"h" long:"--help" description:"Show help for login package." command:"login"`
}

// Login gets access tokens, usernames etc from vault for each API endpoint specified by flag input
func Login() {
	ParseFlags()
	input := CheckInput()
	SelectLogin(input)
}

// ParseFlags enables parsing with flags from Options struct
func ParseFlags() {
	p := flags.NewParser(&Logins, flags.PrintErrors|flags.PassDoubleDash)
	p.Parse()
	// err, _ := p.Parse()
	// if err != nil {
	// 	log.Fatalln("Error occured while parsing flags:", err)
	// }
	if Logins.Help == true {
		PrintFlags()
		os.Exit(1)
	}
}

// CheckInput interprets flag input
func CheckInput() string {
	if Logins.Github == true {
		return "github"
	} else if Logins.Jira == true {
		return "jira"
	} else if Logins.Artifactory == true {
		return "artifactory"
	} else if Logins.Concourse == true {
		return "concourse"
	} else if Logins.Confluence == true {
		return "confluence"
	} else {
		fmt.Print("Required flag not entered.\n")
		i := PromptInput()
		return i
	}
}

// PromptInput prompts User to select which API they want to log in to
func PromptInput() string {
	inputs := map[string]bool{
		"github":      true,
		"artifactory": true,
		"concourse":   true,
		"confluence":  true,
		"jira":        true,
	}

	fmt.Print("Please specify which API you would like to log in to: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")

	// Checks if input is in accepted inputs, otherwise, prompts User for input
	_, ok := inputs[input]
	if ok {
		fmt.Printf("'%v' selected.\n", input)
		return input
	} else {
		fmt.Printf("'%v' is not a valid option. Choose among the following:\n\n", input)
		for k := range inputs {
			fmt.Print("'"+k+"'", "\n")
		}
		fmt.Println("")
		i := PromptInput()
		return i
	}
}

// SelectLogin uses flag or prompt input to initiate login
func SelectLogin(input string) {
	if input == "github" {
		githubLogin.GithubLogin("login")
	} else if input == "jira" {
		// get Jira values from Vault
		fmt.Printf("'" + input + "' login is not set up yet.\n")
	} else if input == "confluence" {
		// get Confluence values from Vault
		fmt.Printf("'" + input + "' login is not set up yet.\n")
	} else if input == "concourse" {
		// get Concourse values from Vault
		fmt.Printf("'" + input + "' login is not set up yet.\n")
	} else if input == "artifactory" {
		// get Artifactory values from Vault
		fmt.Printf("'" + input + "' login is not set up yet.\n")
	} else {
		log.Fatalln("Error selecting login.")
	}
}

// PrintFlags prints out existing flags, short and long form
func PrintFlags() {
	fmt.Print("\n'Login' command line flags:\n")
	fmt.Print("-g\t--github\tLog into Github\n" +
		"-j\t--jira\t\tLog into Jira\n" +
		"-f\t--confluence\tLog into Confluence\n" +
		"-a\t--artifactory\tLog into Artifactory\n" +
		"-c\t--concourse\tLog into Concourse\n" +
		"E.g. 'omniactl login -g' to log into Github.\n\n")
}
