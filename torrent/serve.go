package torrent

import (
	// "net"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DistributedSolutions/DIMWIT/jsonrpc"
	"github.com/fatih/color"
)

func NewTorrentRouter(client *TorrentClient) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/stream", client.HandleStream)
	r.HandleFunc("/torrent", client.HandleStream)

	return r
}

func (c *TorrentClient) HandleStream(w http.ResponseWriter, r *http.Request) {
	c.GetFile(c.selected, w, r)
}

func marshalErr(err *jsonrpc.JSONRPCReponse) []byte {
	data, _ := err.CustomMarshalJSON()
	return data
}

type InfohashReq struct {
	Infohash string `json:"infohash"`
}

func (c *TorrentClient) HandleTorrentAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(marshalErr(jsonrpc.NewInternalRPCSError("Error reading the body of the request", 0)))
		return
	}

	var extra string
	var errorID uint32

	req := jsonrpc.NewEmptyJSONRPCRequest()
	err = json.Unmarshal(data, req)
	if err != nil {
		w.Write(marshalErr(jsonrpc.NewParseError(err.Error(), 0)))
		return
	}

	resp := new(jsonrpc.JSONRPCReponse)
	resp.Id = req.ID
	var result json.RawMessage

	color.Green(fmt.Sprintf("%s: %s: %s", time.Now().Format("15:04:05"), r.Method, req.Method))
	switch req.Method {
	case "select":
		ih := new(InfohashReq)
		err = json.Unmarshal(req.Params, ih)
		if err != nil {
			extra = "Invalid request object, " + err.Error()
			goto InvalidRequest // Bad request data
		}

		err := c.SelectString(ih.Infohash)
		if err != nil {
			extra = fmt.Sprintf("Error selecting torrent: %s", err.Error())
			errorID = 10
			goto CustomError
		}

		data = []byte("success")
		goto Success
	case "watch-vlc":
		ih := new(InfohashReq)
		err = json.Unmarshal(req.Params, ih)
		if err != nil {
			extra = "Invalid request object, " + err.Error()
			goto InvalidRequest // Bad request data
		}

		err := c.SelectString(ih.Infohash)
		if err != nil {
			extra = fmt.Sprintf("Error selecting torrent: %s", err.Error())
			errorID = 10
			goto CustomError
		}

		c.OpenWithVLC()
		data = []byte("success")
		goto Success
	case "info":
		ih := new(InfohashReq)
		err = json.Unmarshal(req.Params, ih)
		if err != nil {
			extra = "Invalid request object, " + err.Error()
			goto InvalidRequest // Bad request data
		}

		mih, err := HexToIH(ih.Infohash)
		if err != nil {
			extra = "Invalid request object, " + err.Error()
			goto InvalidRequest // Bad request data
		}

		t, ok := c.GetTorrent(mih)
		if !ok {
			extra = fmt.Sprintf("No torrent found with %s hash", ih)
			goto InvalidRequest // Bad request data
		}

		type TorrentInfo struct {
			InfoHash  string  `json:"infohash"`
			Name      string  `json:"name"`
			Progress  float64 `json:"progress"`
			TotalSize int64   `json:totalsize`
			HaveInfo  bool    `json:"haveinfo"`
		}

		ti := new(TorrentInfo)
		ti.InfoHash = t.InfoHash().HexString()
		ti.Name = t.Name()
		ti.Progress = c.percentage(t.InfoHash())
		ti.TotalSize = t.Length()
		ti.HaveInfo = t.Info() != nil

		data, err = json.Marshal(ti)
		if err != nil {
			extra = err.Error()
			goto InternalError
		}
		goto Success
	default:
		extra = req.Method
		goto MethodNotFound
	}

	return

	// Easier to handle general here
Success:
	result = json.RawMessage(data)
	resp.Result = &result
	data, _ = resp.CustomMarshalJSON()
	w.Write(data)
	return
MethodNotFound:
	w.Write(marshalErr(jsonrpc.NewMethodNotFoundError(extra, req.ID)))
	return
InvalidRequest:
	w.Write(marshalErr(jsonrpc.NewInvalidRequestError(extra, req.ID)))
	return
CustomError:
	w.Write(marshalErr(jsonrpc.NewCustomError(extra, req.ID, errorID)))
	return
InternalError:
	w.Write(marshalErr(jsonrpc.NewInternalRPCSError(extra, req.ID)))
	return
}
