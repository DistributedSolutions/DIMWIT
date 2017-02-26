package engine

import (
	"bufio"
	"fmt"
	"os"
)

var _ = fmt.Sprintf("")

var HelpText string

// Control function lasts until signal hit
func Control(w *WholeState) {
	scanner := bufio.NewScanner(os.Stdin)
	HelpText = "------------------------------------------    Commands    ------------------------------------------\n"
	AddHelp("|---[command]---|", "|---[text]---|")

	// Commands
	// Add Helps
	AddHelp("h || help", "Display help messages")
	AddHelp("c", "Display Constructor completed height")
	AddHelp("a", "Turn on api")
	AddHelp("a", "Shut off api")

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
		case "c":
			fmt.Printf("Constructor Completed Height: %d\n", w.Constructor.CompletedHeight)
		case "a":
			w.Provider.Close()
		case "w":
			w.Provider.Serve()
		}
	}
}

func AddHelp(command string, text string) {
	HelpText += fmt.Sprintf("|   %-30s%s\n", command, text)
}
