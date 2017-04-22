package torrent

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type TorrentClient struct {
	client *torrent.Client

	selectedReader SeekableContent
	selected       metainfo.Hash
}

func NewTorrentClientFromConfig(con *TopLevelConfig) (*TorrentClient, error) {
	os.MkdirAll(con.AConfig.DataDir, constants.DIRECTORY_PERMISSIONS)

	c := new(TorrentClient)
	cli, err := torrent.NewClient(con.AConfig)
	if err != nil {
		return nil, err
	}

	c.client = cli
	return c, err
}

func NewTorrentClient() (*TorrentClient, error) {
	return NewTorrentClientFromConfig(NewTopLevelConfig())
}

func (c *TorrentClient) Close() {
	c.client.Close()
}

func (c *TorrentClient) SelectString(ih string) error {
	infohash, err := HexToIH(ih)
	if err != nil {
		return err
	}
	return c.Select(infohash)
}

func (c *TorrentClient) Select(infohash metainfo.Hash) error {
	if c.selectedReader != nil {
		c.selectedReader.Close()
	}
	// TODO: Check if infohash exists
	c.selected = infohash
	return nil
}

func (c *TorrentClient) AddMagnet(uri string, download bool) (*torrent.Torrent, error) {
	t, err := c.client.AddMagnet(uri)
	if err != nil {
		return nil, err
	}

	downloadTorrent(t, download)
	return t, nil
}

func downloadTorrent(t *torrent.Torrent, full bool) {
	go func(ti *torrent.Torrent, all bool) {
		<-ti.GotInfo()
		if all {
			t.DownloadAll()
		}
	}(t, full)
}

func (c *TorrentClient) GetTorrent(infohash metainfo.Hash) (torrent *torrent.Torrent, ok bool) {
	return c.client.Torrent(infohash)
}

func HexToIH(ih string) (metainfo.Hash, error) {
	if len(ih) != 40 {
		return metainfo.Hash{}, fmt.Errorf("infohash must be 40 bytes, found %d", len(ih))
	}
	_, err := hex.DecodeString(ih)
	if err != nil {
		return metainfo.Hash{}, err
	}
	return metainfo.NewHashFromHex(ih), nil
}

// percentage of torrent downloaded
func (c *TorrentClient) percentage(infohash metainfo.Hash) float64 {
	t, ok := c.client.Torrent(infohash)
	if !ok || t == nil {
		return 0
	}

	info := t.Info()

	if info == nil {
		return 0
	}

	return float64(t.BytesCompleted()) / float64(info.TotalLength()) * 100
}

// ReadyForPlayback checks if the torrent is ready for playback or not.
// We wait until 5% of the torrent to start playing.
func (c *TorrentClient) ReadyForPlayback(infohash metainfo.Hash) bool {
	return c.percentage(infohash) > 5
}

// GetFile is an http handler to serve the biggest file managed by the client.
func (c *TorrentClient) GetFile(infohash metainfo.Hash, w http.ResponseWriter, r *http.Request) error {
	target := c.getLargestFile(infohash)
	if target == nil {
		return fmt.Errorf("no file found")
	}

	t, ok := c.client.Torrent(infohash)
	if !ok || t == nil {
		return fmt.Errorf("no torrent found with infohash %s", infohash.HexString())
	}

	entry, err := NewFileReader(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("error in file entry: %s", err.Error())
	}

	defer func() {
		if err := entry.Close(); err != nil {
			log.Printf("Error closing file reader: %s\n", err)
		}
	}()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+t.Info().Name+"\"")
	http.ServeContent(w, r, target.DisplayPath(), time.Now(), entry)
	return nil
}

func (c *TorrentClient) getLargestFile(infohash metainfo.Hash) *torrent.File {
	var target torrent.File
	var maxSize int64
	t, ok := c.client.Torrent(infohash)
	if !ok || t == nil {
		return nil
	}

	for _, file := range t.Files() {
		if maxSize < file.Length() {
			maxSize = file.Length()
			target = file
		}
	}

	return &target
}

func (c *TorrentClient) ClientStatus() string {
	buf := new(bytes.Buffer)
	c.client.WriteStatus(buf)
	return string(buf.Next(buf.Len()))
}

// Gets a complete list of torrent files
func (c *TorrentClient) GetTorrentFiles(ih string) ([]torrent.File, error) {
	metainfo, err := HexToIH(ih)
	if err != nil {
		return nil, err
	}
	torrent, ok := c.GetTorrent(metainfo)
	if ok {
		return torrent.Files(), nil
	}
	return nil, errors.New("Torrent Does not exist. No files are going to be returned")
}

func (c *TorrentClient) ShortStatus() string {
	tors := c.client.Torrents()
	resp := fmt.Sprintf("--- Client ---\nTotal Torrents: %d\n", len(tors))
	for i, t := range tors {
		resp += fmt.Sprintf(" --- Torrent %d\nName:%s\nInfoHash: %s\nHaveInfo: %t\nProgress: %.2f%s\n",
			i, t.Name(), t.InfoHash().HexString(), t.Info() != nil, c.percentage(t.InfoHash()), "%")
	}

	return resp
}
