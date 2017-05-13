package engine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/testhelper"
	"github.com/fatih/color"
)

// 22min -- Louie!
// magnet:?xt=urn:btih:852ab3218be18e12d822a31440e43a718ab66f26&dn=Louie.S05E03.HDTV.x264-ASAP%5Bettv%5D&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Fzer0day.ch%3A1337&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969

// 42 Seconds -- mpg
//		magnet:?xt=urn:btih:JYKA2VEC44GNPIPI5TYUAQKBSGST4OEP
// 6.5min
//		magnet:?xt=urn:btih:TEEM5VYAOHVPP7224HCLWLADOHT6GLMX
// 2.5 hr DVDR
//		magnet:?xt=urn:btih:H4QAQQVG67LPBINAHFYPMQKKHHKEE3RV
// 2.5hr MPEG-2
//		magnet:?xt=urn:btih:eb5156f0f92001c5897064072953a9ccd5b00556&dn=The.Matrix-ASCII.mpg&tr=http%3A%2F%2Ftracker.xpear.de%2Fannounce&tr=http%3A%2F%2Ftracker.irc.su%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.torrentbox.com%3A2710%2Fannounce&tr=http%3A%2F%2Finferno.demonoid.com%3A3393%2Fannounce&tr=http%3A%2F%2Fwww.sumotracker.com%2Fannounce&tr=http%3A%2F%2Ffr33dom.h33t.com%3A3310%2Fannounce&tr=http%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=http%3A%2F%2Frlskingbt.appspot.com%2Fannounce&tr=http%3A%2F%2Fsadnem.appspot.com%2Fannounce&tr=http%3A%2F%2Ftracker.utternerd.org%3A6969%2Fannounce&tr=http%3A%2F%2Fannounce.xxx-tracker.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.ccc.de%2Fannounce&tr=http%3A%2F%2Fcpleft.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.ccc.de%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=http%3A%2F%2Ftracker.torrentbay.to%3A6969%2Fannounce&tr=http%3A%2F%2Ftracker.ex.ua%2Fannounce&tr=http%3A%2F%2Ftracker.yify-torrents.com%3A80%2Fannounce&tr=http%3A%2F%2Fbigtorrent.org%3A2710%2Fannounce&tr=http%3A%2F%2Fsline.net%3A2710%2Fannounce

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
	AddHelp("m[s/l] [-a <amount>]", "Add a random channel. S -> Small amount of data. L-> Large amount data. A -> amount of times to add new channels.")
	AddHelp("mchlf [-c] <url to file>", "Add a channelList from file that is a json. C -> current working directory of golang.")
	AddHelp("F[HASH]", "Finds value by given hash from Provider")
	AddHelp("F[HASH]", "Finds value by given hash from Constructor")
	AddHelp("aC", "Prints out root chain ids of all channels")
	AddHelp("wdb", "Wipes db channels clean cascades to almost all other tables. Keeps tags.")
	AddHelp("[MAG_LINK]", "Torrents a magnet link")
	AddHelp("stream [IH]", "Stream the infohash to port 1111")
	AddHelp("ts[l]", "Shows torrent status, 'l' for long")
	AddHelp("i", "Increment fake factom height")
	AddHelp("hard-reset-interfaceDB", "Hard reset the db... be sure you want to do this.")

	var last string
	var err error
	//var chanList []common.Channel
	var amount int
	var fileName string
	var _ = amount
	// Start loop
	for scanner.Scan() {
		err = nil

		cmd := scanner.Text()
		if cmd == "!!" {
			cmd = last
		}
		last = cmd
		//chanList = nil
		amount = 1
		fileName = ""

		if len(cmd) > 5 && cmd[:6] == "ms -a " || len(cmd) > 5 && cmd[:6] == "ml -a " {
			strArr := strings.Split(cmd, " -a ")
			cmd = strArr[0]
			amount, err = strconv.Atoi(strArr[1])
			if err != nil {
				fmt.Printf("Error coverting string [%s] to number. Setting amount to 1.", strArr[1])
				amount = 1
			}
		}

		if len(cmd) > 4 && cmd[:5] == "mchlf" {
			strArr := strings.Split(cmd, " ")
			cmd = strArr[0]
			if len(strArr) > 1 && strings.Contains(strArr[1], "-a") {
				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
				fileName = filepath.Join(dir, strArr[2])
			} else {
				fileName = strArr[1]
			}
		}

		switch {
		case cmd == "exit":
			os.Exit(1)
		case cmd == "h":
			fallthrough
		case cmd == "help":
			fmt.Println(HelpText[:len(HelpText)-1])
			fmt.Println("----------------------------------------------------------------------------------------------------")
		case cmd == "i":
			h, err := testhelper.IncrementFakeHeight(w.FactomClient)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
			fmt.Printf("Incrementing. At height %d\n", h)
		case cmd == "c":
			fmt.Printf("Constructor Completed Height: %d\n", w.Constructor.CompletedHeight)
		case cmd == "a":
			w.Provider.Close()
		case cmd == "w":
			w.Provider.Serve()
		case cmd == "ms":
			fmt.Println("Adding small channels....")
			rc := common.RandomNewSmallChannel()
			err := w.WriteHelper.InitiateChannel(rc)
			if err != nil {
				fmt.Println(err.LogError)
			} else {
				err := w.WriteHelper.UpdateChannel(rc)
				if err != nil {
					fmt.Println(err.LogError)
				}
				for _, content := range rc.Content.ContentList {
					content.RootChainID = rc.RootChainID
					err := w.WriteHelper.AddContent(&content)
					if err != nil {
						fmt.Println(err.LogError)
					}
				}
			}
			fmt.Println("Root:", rc.RootChainID.String())
			// chanList, err = testhelper.AddChannelsToClient(w.FactomClient, amount, true)
			// fallthrough
		case cmd == "ml":
			/*if chanList == nil && err == nil {
				fmt.Println("Adding large channels....")
				chanList, err = testhelper.AddChannelsToClient(w.FactomClient, amount, false)
			}
			if err != nil {
				fmt.Printf("Error: " + err.Error())
			} else {
				fmt.Printf("------ %d Channels Created ------\n", len(chanList))
				for i, c := range chanList {
					fmt.Printf("Channel [%d]: %s\n", i, c.RootChainID.String())
				}
			}
			chanList, err = nil, nil*/
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
		case cmd == "wdb":
			err := w.Constructor.SqlGuy.DeleteDBChannels()
			if err != nil {
				fmt.Println("Error deleting DB Channels: " + err.Error())
				break
			}
			fmt.Printf("DB channels deleted. Note cascading effect.\n")
		case cmd == "mchlf":
			data, err := ioutil.ReadFile(fileName)
			color.Blue("About to read from file: %s", fileName)
			if err != nil {
				color.Red("Error reading file: %s with error: %s", fileName, err.Error())
				break
			}
			chanList := new(common.ChannelList)
			err = json.Unmarshal(data, &chanList)
			if err != nil {
				color.Red("Error unmarshaling binary for chanlist: %s", err.Error())
				break
			}
			err = testhelper.AddChannelsFromFileToClient(w.FactomClient, chanList, true)
			if err != nil {
				color.Red("Error adding channels from file: %s with error: %s", fileName, err.Error())
				break
			}
			color.Blue("Finished reading from file: %s", fileName)
		case cmd == "hard-reset-interfaceDB":
			w.Constructor.SqlGuy.DeleteDB()
			sql, _ := constructor.NewSqlWriter()
			w.Constructor.SqlGuy = sql
		case cmd == "ts":
			fmt.Printf("%s", w.TorrentClient.ShortStatus())
		case cmd == "tsl":
			fmt.Printf("%s", w.TorrentClient.ClientStatus())
		case len(cmd) > 10 && cmd[:6] == "magnet":
			// Download a torrent magnet link
			link := cmd
			t, err := w.TorrentClient.AddMagnet(link, true)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("Torrent started under infohash %s\n", t.InfoHash().HexString())
			}
			var _ = t
			w.TorrentClient.Select(t.InfoHash())
		case len(cmd) > 10 && cmd[:6] == "stream":
			ih := cmd[7:]
			w.TorrentClient.SelectString(ih)
		case cmd == "watch":
			w.TorrentClient.OpenWithVLC()
			fmt.Println("Watching!")
		default:
			fmt.Printf("No command found\n")
		}
	}
}

func AddHelp(command string, text string) {
	HelpText += fmt.Sprintf("|   %-30s%s\n", command, text)
}
