package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/jsonrpc"
	"github.com/DistributedSolutions/DIMWIT/torrent"
	"github.com/fatih/color"
)

type VerifyChannel struct {
	Channel common.Channel `json:"channel"`
	Paths   []string       `json:"path"`
}

type AddContent struct {
	Channel common.Channel `json:"channel"`
	Path    string         `json:"path"`
}

type SubmitChannel struct {
	ChannelHash primitives.Hash `json:"hash"`
}

type ApiService struct {
	Provider *Provider
}

func Vars(r *http.Request) map[string]string {
	return make(map[string]string)
}

func marshalErr(err *jsonrpc.JSONRPCReponse) []byte {
	data, _ := err.CustomMarshalJSON()
	return data
}

func (apiService *ApiService) HandleAPICalls(w http.ResponseWriter, r *http.Request) {
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
	case "get-constants":
		tempData, err := constants.ConstantJSONMarshal()
		if err != nil {
			extra = "Failed to unmarshal channel"
			goto InternalError
		}
		data = *tempData
		goto Success
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
		if channel == nil {
			extra = "Channel hash not found"
			errorID = 8
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
	case "torrent-stream-stat":
		var hashString string
		err := json.Unmarshal(req.Params, &hashString)
		if err != nil {
			color.Red("Error unmarshall torrent-stream-stat: %s", err.Error())
			extra = fmt.Sprintf("Error unmarshall torrent-stream-stat: %s", err.Error())
			goto InvalidRequest // Bad Request data
		}
		stats, err := apiService.GetTorrentStreamStats(hashString)
		if err != nil {
			extra = err.Error()
			errorID = 9
			goto CustomError
		}

		data, err = json.Marshal(stats)
		if err != nil {
			extra = fmt.Sprintf("Failed to unmarshal torrent file stats: %s", err.Error())
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
		var verifiedChannel *common.Channel
		hash := verifyChannel.Channel.RootChainID
		if hash.Empty() {
			verifiedChannel, err = apiService.Provider.CreateChannel(&verifyChannel.Channel, verifyChannel.Paths)
			if err != nil {
				color.Red("Error verifying channel: %s", err.Error())
				extra = fmt.Sprintf("Error verifying new channel with error: %s", err.Error())
				errorID = 6
				goto CustomError
			}
		} else {
			verifiedChannel, err = apiService.Provider.UpdateChannel(&verifyChannel.Channel, verifyChannel.Paths)
			if err != nil {
				color.Red("Error verifying UpdateChanne: %s", err.Error())
				extra = fmt.Sprintf("Error verifying update channel with error: %s", err.Error())
				errorID = 7
				goto CustomError
			}
		}

		data, err = (*verifiedChannel).CustomMarshalJSON()
		if err != nil {
			fmt.Println(err.Error())
			extra = "Failed to unmarshal verified channel"
			goto InternalError
		}
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

func (apiService *ApiService) GetTorrentStreamStats(torrentHash string) (*torrent.JSONFiles, error) {
	return apiService.Provider.TorrentClientInterface.GetTorrentFileMetaData(torrentHash)
}

func (apiService *ApiService) GetCompleteHeight() (uint32, error) {
	return apiService.Provider.GetCompleteHeight()
}

func (apiService *ApiService) VerifyChannel(ch *common.Channel) (int, error) {
	cost, apiError := (*apiService.Provider.CreationTool.ApiTool).VerifyChannel(ch)
	if apiError.LogError != nil {
		color.Red("VerifyChannel Error: %s", apiError.LogError.Error())
		return cost, apiError.UserError
	}
	return cost, nil
}

func (apiService *ApiService) InitiateChannel(ch *common.Channel) error {
	_, _, apiError := (*apiService.Provider.CreationTool.ApiTool).InitiateChannel(ch)
	if apiError.LogError != nil {
		color.Red("InitiateChannel Error: %s", apiError.LogError.Error())
		return apiError.UserError
	}
	return nil
}

func (apiService *ApiService) UpdateChannel(ch *common.Channel) error {
	_, _, apiError := (*apiService.Provider.CreationTool.ApiTool).UpdateChannel(ch)
	if apiError.LogError != nil {
		color.Red("UpdateChannel Error: %s", apiError.LogError.Error())
		return apiError.UserError
	}
	return nil
}

func (apiService *ApiService) DeleteChannel(h *primitives.Hash) error {
	apiError := (*apiService.Provider.CreationTool.ApiTool).DeleteChannel(h)
	if apiError.LogError != nil {
		color.Red("DeleteChannel Error: %s", apiError.LogError.Error())
		return apiError.UserError
	}
	return nil
}

func (apiService *ApiService) VerifyContent(c *common.Content) (int, error) {
	cost, apiError := (*apiService.Provider.CreationTool.ApiTool).VerifyContent(c)
	if apiError.LogError != nil {
		color.Red("VerifyContent Error: %s", apiError.LogError.Error())
		return cost, apiError.UserError
	}
	return cost, nil
}

func (apiService *ApiService) AddContent(c *common.Content) error {
	_, _, apiError := (*apiService.Provider.CreationTool.ApiTool).AddContent(c, &c.ContentID)
	if apiError.LogError != nil {
		color.Red("UpdateChannel Error: %s", apiError.LogError.Error())
		return apiError.UserError
	}
	return nil
}

func (apiService *ApiService) DeleteContent(h *primitives.Hash) error {
	apiError := (*apiService.Provider.CreationTool.ApiTool).DeleteContent(h)
	if apiError.LogError != nil {
		color.Red("DeleteContent Error: %s", apiError.LogError.Error())
		return apiError.UserError
	}
	return nil
}
