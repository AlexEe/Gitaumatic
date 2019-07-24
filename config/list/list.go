package list

import (
	// "errors"
	"fmt"
	"github.com/fatih/color"
	// "github.com/manifoldco/promptui"
	ini "gopkg.in/ini.v1"
	// "io"
	"log"
	// "os"
	// "regexp"
	// "strings"
)

var (
	github      string
	jira        string
	artifactory string
	confluence  string
	concourse   string
	vault       string
)

func ListConfigFile() {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Display config file")
	// whiteBold := color.New(color.FgHiWhite, color.Bold)
	greenBold := color.New(color.FgGreen, color.Bold)
	URLs := make(map[string]string)

	cfg, err := ini.Load(".omniactl")
	if err != nil {
		log.Fatalln("Error loading .omniactl config file:", err)
	}

	github = fmt.Sprint(cfg.Section("config").Key("github").String())
	URLs["Github"] = github

	jira = fmt.Sprint(cfg.Section("config").Key("jira").String())
	URLs["Jira"] = jira

	confluence = fmt.Sprint(cfg.Section("config").Key("confluence").String())
	URLs["Confluence"] = confluence

	artifactory = fmt.Sprint(cfg.Section("config").Key("artifactory").String())
	URLs["Artifactory"] = artifactory

	concourse = fmt.Sprint(cfg.Section("config").Key("concourse").String())
	URLs["Concourse"] = concourse

	vault = fmt.Sprint(cfg.Section("config").Key("vault").String())
	URLs["Vault"] = vault

	for name, url := range URLs {
		greenBold.Printf("%-15v ", name)
		fmt.Println(url)
	}
	fmt.Println("")
}
