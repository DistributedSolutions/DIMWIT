package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/provider/jsonrpc"
	"github.com/fatih/color"
)

type AddChannel struct {
	Channel common.Channel `json:"channel"`
	Path    string         `json:"path"`
}

type ApiService struct {
	Provider *Provider
}

func Vars(r *http.Request) map[string]string {
	return make(map[string]string)
}

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
			errorID = 1
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
			errorID = 1
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
			errorID = 1
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
			errorID = 1
			goto CustomError
		}

		data, err = json.Marshal(stats)
		if err != nil {
			extra = "Failed to unmarshal stats"
			goto InternalError
		}
		goto Success
	case "add-channel":
		addChannel := new(AddChannel)
		err = json.Unmarshal(req.Params, addChannel)
		if err != nil {
			extra = "Invalid request object, " + err.Error()
			goto InvalidRequest // Bad request data
		}
		s := string(addChannel.Path)
		hash, err := apiService.Provider.CreateChannel(&addChannel.Channel, s)
		if err != nil {
			extra = fmt.Sprintf("Error creating new channel with error: %s", err)
			errorID = 1
			goto CustomError
		}
		color.Blue("Adding new Channel with hash: %s", hash.String())
		err = apiService.Provider.SubmitChannel(*hash)
		if err != nil {
			extra = fmt.Sprintf("Error submiting new channel with error: %s", err)
			errorID = 1
			goto CustomError
		}
		color.Blue("Finished adding in new Channel")
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

func (apiService *ApiService) GetChannel(hash primitives.Hash) (*common.Channel, error) {
	return apiService.Provider.GetChannel(hash.String())
}

func (apiService *ApiService) GetChannels(hashes primitives.HashList) (*common.ChannelList, error) {
	channelList := make([]common.Channel, 0)
	for _, channelHash := range hashes.GetHashes() {
		channel, err := apiService.Provider.GetChannel(channelHash.String())
		if err != nil {
			return nil, err
		}
		if channel == nil {
			continue
			// return nil, fmt.Errorf("Content %s not found", contentHash.String())
		}
		channelList = append(channelList, *channel)
	}

	if len(channelList) == 0 {
		return nil, fmt.Errorf("No channels found by those hashes")
	}
	chList := common.ChannelList{
		List: channelList,
	}
	return &chList, nil
}

func (apiService *ApiService) GetContent(hash primitives.Hash) (*common.Content, error) {
	return apiService.Provider.GetContent(hash.String())
}

func (apiService *ApiService) GetContents(hashes primitives.HashList) (*common.ContentList, error) {
	contentList := make([]common.Content, 0)
	for _, contentHash := range hashes.GetHashes() {
		content, err := apiService.GetContent(contentHash)
		if err != nil {
			return nil, err
		}
		if content == nil {
			continue
			// return nil, fmt.Errorf("Content %s not found", contentHash.String())
		}
		contentList = append(contentList, *content)
	}
	if len(contentList) == 0 {
		return nil, fmt.Errorf("No content found by those hashes")
	}
	cList := common.ContentList{
		ContentList: contentList,
	}
	return &cList, nil
}

func (apiService *ApiService) GetStats() (*DatabaseStats, error) {
	return apiService.Provider.GetStats()
}

func (apiService *ApiService) GetCompleteHeight() (uint32, error) {
	return apiService.Provider.GetCompleteHeight()
}
