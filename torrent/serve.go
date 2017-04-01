package torrent

import (
	// "net"
	"net/http"
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

func (c *TorrentClient) HandleTorrentAPI(w http.ResponseWriter, r *http.Request) {
	/*
	   func (apiService *ApiService) HandleAPICalls(w http.ResponseWriter, r *http.Request) {
	   	w.Header().Set("Content-Type", "application/json")
	   	data, err := ioutil.ReadAll(r.Body)
	   	if err != nil {
	   		jErr := jsonrpc.NewInternalRPCSError("Error reading the body of the request", 0)
	   		data, _ := jErr.CustomMarshalJSON()
	   		w.Write(data)
	   		return
	   	}

	   	var extra string
	   	var errorID uint32
	   	jErr := new(jsonrpc.JSONRPCReponse)

	   	req := jsonrpc.NewEmptyJSONRPCRequest()
	   	err = json.Unmarshal(data, req)
	   	if err != nil {
	   		jErr := jsonrpc.NewParseError(err.Error(), 0)
	   		data, _ := jErr.CustomMarshalJSON()
	   		w.Write(data)
	   		return
	   	}

	   	resp := new(jsonrpc.JSONRPCReponse)
	   	resp.Id = req.ID
	   	var result json.RawMessage

	   	color.Green(fmt.Sprintf("%s: %s: %s", time.Now().Format("15:04:05"), r.Method, req.Method))
	   	switch req.Method {
	   	case "get-channel":
	   		hash := new(primitives.Hash)
	   		err = json.Unmarshal(req.Params, hash)
	   		if err != nil {
	   			extra = "Invalid request object, " + err.Error()
	   			goto InvalidRequest // Bad request data
	   		}
	   		channel, err := apiService.GetChannel(*hash)
	   		if err != nil {
	   			extra = "Channel not found"
	   			errorID = 1
	   			goto CustomError
	   		}

	   		data, err = channel.CustomMarshalJSON()
	   		if err != nil {
	   			extra = "Failed to unmarshal channel"
	   			goto InternalError
	   		}
	   		goto Success
	   	case "get-channels":
	   		hashList := new(primitives.HashList)
	   		err = json.Unmarshal(req.Params, hashList)
	   		if err != nil {
	   			color.Red("ERRROR :(: %s", err.Error())
	   			extra = "Invalid request object, " + err.Error()
	   			goto InvalidRequest // Bad request data
	   		}
	   		channels, err := apiService.GetChannels(*hashList)
	   		if err != nil {
	   			extra = "Channels not found"
	   			errorID = 2
	   			goto CustomError
	   		}

	   		data, err = channels.CustomMarshalJSON()
	   		if err != nil {
	   			extra = "Failed to unmarshal channels"
	   			goto InternalError
	   		}
	   		goto Success
	   	case "get-content":
	   		hash := new(primitives.Hash)
	   		err = json.Unmarshal(req.Params, hash)
	   		if err != nil {
	   			extra = "Invalid request object, " + err.Error()
	   			goto InvalidRequest // Bad request data
	   		}
	   		content, err := apiService.GetContent(*hash)
	   		if err != nil {
	   			extra = "Content not found"
	   			errorID = 3
	   			goto CustomError
	   		}

	   		data, err = json.Marshal(content)
	   		if err != nil {
	   			extra = "Failed to unmarshal content"
	   			goto InternalError
	   		}
	   		goto Success
	   	case "get-contents":
	   		hashList := new(primitives.HashList)
	   		err = json.Unmarshal(req.Params, hashList)
	   		if err != nil {
	   			color.Red("ERRROR :(: %s", err.Error())
	   			extra = "Invalid request object, " + err.Error()
	   			goto InvalidRequest // Bad request data
	   		}
	   		contents, err := apiService.GetContents(*hashList)
	   		if err != nil {
	   			extra = "Contents not found"
	   			errorID = 4
	   			goto CustomError
	   		}

	   		data, err = json.Marshal(contents)
	   		if err != nil {
	   			extra = "Failed to unmarshal contents"
	   			goto InternalError
	   		}
	   		goto Success
	   	case "get-stats":
	   		stats, err := apiService.GetStats()
	   		if err != nil {
	   			extra = err.Error()
	   			errorID = 5
	   			goto CustomError
	   		}

	   		data, err = json.Marshal(stats)
	   		if err != nil {
	   			extra = "Failed to unmarshal stats"
	   			goto InternalError
	   		}
	   		goto Success
	   	case "verify-channel":
	   		verifyChannel := new(VerifyChannel)
	   		err = json.Unmarshal(req.Params, verifyChannel)
	   		if err != nil {
	   			extra = "Invalid request object, " + err.Error()
	   			goto InvalidRequest // Bad request data
	   		}
	   		hash := verifyChannel.Channel.RootChainID
	   		if hash.Empty() {
	   			newHash, err := apiService.Provider.CreateChannel(&verifyChannel.Channel, verifyChannel.Paths)
	   			if err != nil {
	   				color.Red("Error verifying channel: %s", err.Error())
	   				extra = fmt.Sprintf("Error verifying new channel with error: %s", err.Error())
	   				errorID = 6
	   				goto CustomError
	   			}
	   			hash = *newHash
	   		} else {
	   			err := apiService.Provider.UpdateChannel(&verifyChannel.Channel, verifyChannel.Paths)
	   			if err != nil {
	   				color.Red("Error verifying UpdateChanne: %s", err.Error())
	   				extra = fmt.Sprintf("Error verifying update channel with error: %s", err.Error())
	   				errorID = 6
	   				goto CustomError
	   			}
	   		}
	   		data = []byte(hash.String())
	   		goto Success
	   	case "submit-channel":
	   		submitChannel := new(SubmitChannel)
	   		err = json.Unmarshal(req.Params, submitChannel)
	   		if err != nil {
	   			extra = "Invalid request object, " + err.Error()
	   			goto InvalidRequest // Bad request data
	   		}
	   		err = apiService.Provider.SubmitChannel(submitChannel.ChannelHash)
	   		if err != nil {
	   			color.Red("Error submitting channel: %s", err.Error())
	   			extra = fmt.Sprintf("Error submiting new channel with error: %s", err.Error())
	   			errorID = 7
	   			goto CustomError
	   		}
	   		color.Blue("Finished adding in new Channel")

	   		data = []byte("{}")
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
	   	jErr = jsonrpc.NewMethodNotFoundError(extra, req.ID)
	   	data, _ = jErr.CustomMarshalJSON()
	   	w.Write(data)
	   	return
	   InvalidRequest:
	   	jErr = jsonrpc.NewInvalidRequestError(extra, req.ID)
	   	data, _ = jErr.CustomMarshalJSON()
	   	w.Write(data)
	   	return
	   CustomError:
	   	jErr = jsonrpc.NewCustomError(extra, req.ID, errorID)
	   	data, _ = jErr.CustomMarshalJSON()
	   	w.Write(data)
	   	return
	   InternalError:
	   	jErr = jsonrpc.NewInternalRPCSError(extra, req.ID)
	   	data, _ = jErr.CustomMarshalJSON()
	   	w.Write(data)
	   	return
	   }
	*/
}
