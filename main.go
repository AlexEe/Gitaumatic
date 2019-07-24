package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"omniactl/cmd"

	"github.com/fatih/color"
)

func main() {
	whiteBold := color.New(color.FgHiWhite, color.Bold)
	dat, err := ioutil.ReadFile("logo.txt")
	if err != nil {
		log.Fatalln("Error opening 'logo.txt':", err)
	}
	whiteBold.Print(string(dat))

	white := color.New(color.FgWhite, color.Bold, color.Italic)
	white.Print("Github project on-boarding and management\n")
	fmt.Println("")
	fmt.Println("")
	cmd.Execute()
}
