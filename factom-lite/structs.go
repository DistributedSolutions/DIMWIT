package lite

import (
	"github.com/FactomProject/factom"
)

type DBlock struct {
	Dbentries []struct {
		Chainid string `json:"chainid"`
		Keymr   string `json:"keymr"`
	} `json:"dbentries"`
	Dbhash string `json:"dbhash"`
	Header struct {
		Blockcount   int    `json:"blockcount"`
		Bodymr       string `json:"bodymr"`
		Chainid      string `json:"chainid"`
		Dbheight     int    `json:"dbheight"`
		Networkid    int    `json:"networkid"`
		Prevfullhash string `json:"prevfullhash"`
		Prevkeymr    string `json:"prevkeymr"`
		Timestamp    int    `json:"timestamp"`
		Version      int    `json:"version"`
	} `json:"header"`
	Keymr string `json:"keymr"`
}

type EntryHolder struct {
	Entry     *factom.Entry
	Timestamp int64
	Height    uint32
}
