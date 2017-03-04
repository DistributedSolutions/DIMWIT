package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/provider/jsonrpc"
	"github.com/fatih/color"
)

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

	color.Green(fmt.Sprintf("%s: %s", r.Method, req.Method))
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

func (apiService *ApiService) GetChannels(hashes primitives.HashList) ([]common.Channel, error) {
	channelList := make([]common.Channel, 0)
	for _, channelHash := range hashes.GetHashes() {
		channel, err := apiService.Provider.GetChannel(channelHash.String())
		if err != nil {
			return channelList, err
		}
		channelList = append(channelList, *channel)
	}
	return channelList, nil
}

func (apiService *ApiService) GetContent(hash primitives.Hash) (*common.Content, error) {
	return apiService.Provider.GetContent(hash.String())
}

func (apiService *ApiService) GetStats() (*DatabaseStats, error) {
	return apiService.Provider.GetStats()
}

func (apiService *ApiService) GetCompleteHeight() (uint32, error) {
	return apiService.Provider.GetCompleteHeight()
}
