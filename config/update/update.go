package update

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	ini "gopkg.in/ini.v1"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// Defaults holds the default endpoint values of the six APIs
var Defaults = make(map[string]string)

// UrlFromFlag holds urls set through flag input
var UrlFromFlag = make(map[string]string)

// UrlConfirmed holds urls after having been checked
var UrlConfirmed = make(map[string]string)

var (
	github      string
	jira        string
	artifactory string
	confluence  string
	concourse   string
	vault       string
)

// UpdateConfigFile collects input from string or prompt and updates current config file
func UpdateConfigFile(github string, jira string, confluence string, artifactory string, concourse string, vault string) {
	magentaBold := color.New(color.FgMagenta, color.Bold, color.Underline)
	magentaBold.Println("Action selected: Update config file")

	// Map containing URLs set through flag input
	UrlFromFlag = map[string]string{
		"github":      github,
		"jira":        jira,
		"confluence":  confluence,
		"concourse":   concourse,
		"artifactory": artifactory,
		"vault":       vault,
	}

	// Map containing default URLs
	Defaults = map[string]string{
		"github":      "https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3",
		"jira":        "https://jira.dev.us-east-1.aws.galleon.c.statestr.com",
		"confluence":  "https://confluence.dev.us-east-1.aws.galleon.c.statestr.com",
		"artifactory": "https://artifactory.galleon.c.statestr.com",
		"concourse":   "https://concourse.dev.tn.galleon.c.statestr.com",
		"vault":       "https://defaultvaultaddress.com",
	}

	CheckConfigFile()
	UrlConfirmed = CheckAllFlags(UrlFromFlag)
	WriteToConfigFile(UrlConfirmed)
}

// WriteToConfigFile saves updated URL endpoints to config file
func WriteToConfigFile(Urls map[string]string) {
	nf, err := os.Create(".omniactl")
	whiteBold := color.New(color.FgHiWhite, color.Bold)

	if err != nil {
		log.Fatal("Error updating config file:", err)
	}
	defer nf.Close()

	// Write new config map with updated values
	// Create new config file
	ncf := ("[config]\n" +
		"github=" + Urls["github"] + "\n" +
		"jira=" + Urls["jira"] + "\n" +
		"confluence=" + Urls["confluence"] + "\n" +
		"artifactory=" + Urls["artifactory"] + "\n" +
		"concourse=" + Urls["concourse"] + "\n" +
		"vault=" + Urls["vault"] + "\n")

	// Save values confirmed/added by user in new config file
	_, err = io.Copy(nf, strings.NewReader(ncf))
	if err != nil {
		log.Fatal("Error creating new config file", err)
	}

	whiteBold.Println("Config file updated:")
	for name, url := range Urls {
		fmt.Printf("%-15v %v", name, url)
		fmt.Println("")
	}
}

// CheckAllFlags loops over all flags and checks if they were put or not
func CheckAllFlags(UrlFromFlag map[string]string) map[string]string {
	fmt.Println("")
	for name, url := range UrlFromFlag {
		url = CheckFlag(url, name)
		UrlFromFlag[name] = url
		fmt.Println("")
	}
	return UrlFromFlag
}

// CheckFlag checks an individual flag for correct format and replaces with prompt input when necessary
func CheckFlag(url string, name string) string {
	greenBold := color.New(color.FgGreen, color.Bold)
	cfg, err := ini.Load(".omniactl")
	if err != nil {
		log.Fatalln("Error opening '.omniactl' config file:", err)
	}

	if url != "" {
		check := CheckURLFormat(url)
		switch check {
		case true:
			greenBold.Printf("New %v URL ", name)
			fmt.Println(url)
			return url
		default:
			url = PromptURL(name)
			return url
		}
	} else {
		greenBold.Printf("Current %v URL:", name)
		fmt.Println("")
		fmt.Println(cfg.Section("config").Key(strings.ToLower(name)).String())

		result := AcceptCurrent()
		switch result {
		case "yes":
			url = fmt.Sprint(cfg.Section("config").Key(strings.ToLower(name)).String())
			return url
		default:
			url = PromptURL(name)
			return url
		}
	}
}

// CheckURLFormat checks if input from flag corresponds with address
func CheckURLFormat(URL string) bool {
	isCorrectFormat := regexp.MustCompile("^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/)?[a-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$").MatchString
	if isCorrectFormat(URL) {
		return true
	}
	return false
}

// AcceptCurrent asks if user wants to accept the current value
func AcceptCurrent() string {
	prompt := promptui.Select{
		Label: "Accept current value?",
		Items: []string{"yes", "no"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// PromptURL asks user to provide a new URL and checks its format
func PromptURL(name string) string {

	validate := func(input string) error {
		isCorrectFormat := regexp.MustCompile("^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/)?[a-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$").MatchString
		if isCorrectFormat(input) {
			return nil
		} else {
			return errors.New("Please enter a valid API endpoint, e.g. https://domainname.com/api/v3")
		}
	}

	templates := &promptui.PromptTemplates{
		Success: "{{ . | green | bold }} ",
	}

	label := fmt.Sprintf("New %v URL", name)

	prompt := promptui.Prompt{
		Label:     label,
		Validate:  validate,
		Templates: templates,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

// CheckConfigFile tries to load values from .omniactl config file,
// If file cannot be opened, it sets values with default values
func CheckConfigFile() {
	cfg, err := ini.Load(".omniactl")
	// If file cannot be opened or doesn't exist, use default values
	if err != nil {
		github = Defaults["github"]
		jira = Defaults["jira"]
		confluence = Defaults["confluence"]
		artifactory = Defaults["artifactory"]
		concourse = Defaults["concourse"]
		vault = Defaults["vault"]
	} else {
		// Load values from file. If field is empty, replace with default value
		github = fmt.Sprint(cfg.Section("config").Key("github").String())
		if github == "" {
			github = Defaults["github"]
		}
		jira = fmt.Sprint(cfg.Section("config").Key("jira").String())
		if jira == "" {
			jira = Defaults["jira"]
		}
		confluence = fmt.Sprint(cfg.Section("config").Key("confluence").String())
		if confluence == "" {
			confluence = Defaults["confluence"]
		}
		artifactory = fmt.Sprint(cfg.Section("config").Key("artifactory").String())
		if artifactory == "" {
			artifactory = Defaults["artifactory"]
		}
		concourse = fmt.Sprint(cfg.Section("config").Key("concourse").String())
		if concourse == "" {
			concourse = Defaults["concourse"]
		}
		vault = fmt.Sprint(cfg.Section("config").Key("vault").String())
		if vault == "" {
			vault = Defaults["vault"]
		}
	}
}
