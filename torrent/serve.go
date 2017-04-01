package torrent

import (
	// "net"
	"net/http"
)

func NewTorrentRouter(client *TorrentClient) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/stream", client.HandleStream)

	return r
}

func (c *TorrentClient) HandleStream(w http.ResponseWriter, r *http.Request) {
	c.GetFile(c.selected, w, r)
}
