package engine

import (
	"bufio"
	"fmt"
	"os"
)

var _ = fmt.Sprintf("")

var HelpText string

// Control function lasts until signal hit
func Control() {
	scanner := bufio.NewScanner(os.Stdin)
	HelpText = "------------------------------------------    Commands    ------------------------------------------\n"
	AddHelp("|---[command]---|", "|---[text]---|")

	// Commands
	// Add Helps
	AddHelp("h || help", "Display help messages")

	// Start loop
	for scanner.Scan() {
		switch scanner.Text() {
		case "exit":
			os.Exit(1)
		case "h":
			fallthrough
		case "help":
			fmt.Println(HelpText[:len(HelpText)-1])
			fmt.Println("----------------------------------------------------------------------------------------------------")
		}
	}
}

func AddHelp(command string, text string) {
	HelpText += fmt.Sprintf("|   %-30s%s\n", command, text)
}
