package engine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/testhelper"
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
	AddHelp("s", "Show current Channel/Content count")
	AddHelp("w", "Turn on api")
	AddHelp("a", "Shut off api")
	AddHelp("m[s/l]", "Add a random channel. S -> Small amount of data. L-> Large amount")
	AddHelp("F[HASH]", "Finds value by given hash from Provider")
	AddHelp("F[HASH]", "Finds value by given hash from Constructor")
	AddHelp("aC", "Prints out root chain ids of all channels")

	var last string
	var err error
	var chanList []common.Channel
	// Start loop
	for scanner.Scan() {
		err = nil

		cmd := scanner.Text()
		if cmd == "!!" {
			cmd = last
		}
		last = cmd
		chanList = nil

		switch {
		case cmd == "exit":
			os.Exit(1)
		case cmd == "h":
			fallthrough
		case cmd == "help":
			fmt.Println(HelpText[:len(HelpText)-1])
			fmt.Println("----------------------------------------------------------------------------------------------------")
		case cmd == "c":
			fmt.Printf("Constructor Completed Height: %d\n", w.Constructor.CompletedHeight)
		case cmd == "a":
			w.Provider.Close()
		case cmd == "w":
			w.Provider.Serve()
		case cmd == "ms":
			fmt.Println("Adding small channels....")
			chanList, err = testhelper.AddChannelsToClient(w.FactomClient, 1, true)
			fallthrough
		case cmd == "ml":
			if chanList == nil && err == nil {
				fmt.Println("Adding large channels....")
				chanList, err = testhelper.AddChannelsToClient(w.FactomClient, 1, false)
			}
			if err != nil {
				fmt.Printf("Error: " + err.Error())
			} else {
				fmt.Printf("------ %d Channels Created ------\n", len(chanList))
				for i, c := range chanList {
					fmt.Printf("Channel [%d]: %s\n", i, c.RootChainID.String())
				}
			}
			chanList, err = nil, nil
		case len(cmd) > 1 && cmd[:1] == "F":
			var resp string
			var con *common.Content
			c, err := w.Provider.GetChannel(cmd[1:])
			fmt.Println(c, err)
			if err == nil && c != nil {
				buf := new(bytes.Buffer)
				data, _ := c.CustomMarshalJSON()
				json.Indent(buf, data, "-", "\t")
				resp = fmt.Sprintf("Channel found with that hash\n%s\n", string(buf.Bytes()))
				goto Found
			}

			con, err = w.Provider.GetContent(cmd[1:])
			if err == nil && con != nil {
				buf := new(bytes.Buffer)
				data, _ := json.Marshal(con)
				json.Indent(buf, data, "-", "\t")
				resp = fmt.Sprintf("Content found with that hash\n%s\n", string(buf.Bytes()))
				goto Found
			}

			resp = "Nothing found by that hash\n"
		Found:
			fmt.Print(resp)
		case len(cmd) > 1 && cmd[:1] == "f":
			var resp string
			hash, err := primitives.HexToHash(cmd[1:])
			if err != nil {
				fmt.Printf("Error %s\n", err.Error())
				break
			}
			c, err := w.Constructor.RetrieveChannel(*hash)
			if err == nil && c != nil {
				buf := new(bytes.Buffer)
				data, _ := c.Channel.CustomMarshalJSON()
				json.Indent(buf, data, "-", "\t")
				resp = fmt.Sprintf("Channel found with that hash\n%s\n", string(buf.Bytes()))
			} else {
				resp = "Nothing found by that hash\n"
			}
			fmt.Print(resp)
		case cmd == "aC":
			chanList, err := w.Provider.GetAllChannels()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			fmt.Printf("---- Found %d channels ----\n", len(chanList))
			for _, c := range chanList {
				fmt.Printf("Channel: %s\n", c.RootChainID.String())
			}
		case cmd == "s":
			stats, err := w.Provider.GetStats()
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			fmt.Printf("%d Channels -- %d Content\n", stats.TotalChannels, stats.TotalContent)
		default:
			fmt.Printf("No command found\n")
		}
	}
}

func AddHelp(command string, text string) {
	HelpText += fmt.Sprintf("|   %-30s%s\n", command, text)
}
