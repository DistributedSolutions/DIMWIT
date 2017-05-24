package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/jsonrpc"
	"github.com/DistributedSolutions/DIMWIT/util"
	"github.com/fatih/color"
)

const (
	invalidParameters uint8 = iota
	customError       uint8 = iota
	noError           uint8 = iota
)

type ApiBase struct{}

//method that all of api methods will be tested against
func (a ApiBase) ApiBaseMethod(json json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	return nil, nil, noError
}

// type VerifyChannel struct {
// 	Channel common.Channel `json:"channel"`
// 	Paths   []string       `json:"path"`
// }

type AddContent struct {
	Channel common.Channel `json:"channel"`
	Path    string         `json:"path"`
}

type SubmitChannel struct {
	ChannelHash primitives.Hash `json:"hash"`
}

//used to provide the methods for api calls
type ApiProvider struct {
	Provider *Provider
}

type ApiService struct {
	Provider *Provider
	Api      ApiProvider
}

func Vars(r *http.Request) map[string]string {
	return make(map[string]string)
}

func marshalErr(err *jsonrpc.JSONRPCReponse) []byte {
	data, _ := err.CustomMarshalJSON()
	return data
}

// func (apiProvider ApiProvider) GetStuff() {
// 	return
// }

// func (apiProvider ApiProvider) GetStuff2() (error, error, error) {
// 	return nil, nil, nil
// }

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

	methodCamelCase := util.DashDelimiterToCamelCase(req.Method)
	color.Green(fmt.Sprintf("%s: %s: %s", time.Now().Format("15:04:05"), r.Method, methodCamelCase))

	apiProvider := ApiProvider{
		apiService.Provider,
	}

	//get the api interface object
	_, ok := reflect.TypeOf(apiProvider).MethodByName(methodCamelCase)
	//checks to see if method exists to call it
	if !ok {
		//method does not exist
		extra = fmt.Sprintf("Method not found: %s", methodCamelCase)
		color.Red(extra)
		goto MethodNotFound
	} else {
		method := reflect.ValueOf(apiProvider).MethodByName(methodCamelCase)
		//method exists
		//successResponse JSON before
		//apiError = object containing messages
		//errorType:
		//		InvalidParameters
		//		CustomError
		in := []reflect.Value{reflect.ValueOf(req.Params)}
		resultValues := method.Call(in)

		if !resultValues[1].CanInterface() || !resultValues[2].CanInterface() {
			extra = fmt.Sprintf("Unable to interface either 1: %d, or 2: %d of the return values.", resultValues[1].CanInterface(), resultValues[2].CanInterface())
			color.Red(extra)
			goto InternalError
		}
		//if the apiError is not null assume that the call was successful
		if resultValues[1].Interface().(*util.ApiError) != nil {
			//api error log error and send response
			apiError := resultValues[1].Interface().(*util.ApiError)
			errorType := resultValues[2].Interface().(uint8)
			color.Red("Error with method: %s, apiError: %s, errorType: %d", methodCamelCase, apiError.LogError.Error(), errorType)
			extra = apiError.UserError.Error()
			switch errorType {
			case invalidParameters:
				goto InvalidParameters
			case customError:
				goto CustomError
			}
		}
		if resultValues[0].CanInterface() {
			if resultValues[0].Interface() != nil {
				data, err = json.Marshal(resultValues[0].Interface())
				if err != nil {
					extra = "Failed to marshal content"
					goto InternalError
				}
			}
			goto Success
		}
		//if can not get interface from value there is an internal error...
		extra = fmt.Sprintf("Unable to interface successful response.")
		color.Red(extra)
		goto InternalError
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
	//--------------
	// NEVER USED THIS... KEEP FOR NOW UNTIL WE DONT WANT IT
	//-------------
	// InvalidRequest:
	// 	w.Write(marshalErr(jsonrpc.NewInvalidRequestError(extra, req.ID)))
	// 	return
InvalidParameters:
	w.Write(marshalErr(jsonrpc.NewInvalidParametersError(extra, req.ID)))
	return
CustomError:
	w.Write(marshalErr(jsonrpc.NewCustomError(extra, req.ID, errorID)))
	return
InternalError:
	w.Write(marshalErr(jsonrpc.NewInternalRPCSError(extra, req.ID)))
	return
}

func (apiProvider ApiProvider) PostTorrentStreamSeek(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	seconds := new(float64)
	err := json.Unmarshal([]byte(input), seconds)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall torrent-stream-stat: %s", err.Error()),
				fmt.Errorf("Error unmarshall torrent-stream-stat: %s", err.Error()),
			},
			invalidParameters
	}
	s := "Success"
	apiProvider.Provider.TorrentClientInterface.SetTorrentSeek(*seconds)
	retVal := new(interface{})
	*retVal = s
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetStats(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	stats, err := apiProvider.Provider.GetStats()
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error retrieving stats: %s", err.Error()),
				fmt.Errorf("Error retrieving stats: %s", err.Error()),
			},
			customError
	}
	retVal := new(interface{})
	*retVal = stats
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetConstants(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	constants := constants.ConstantJSONMarshal()
	retVal := new(interface{})
	*retVal = constants
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetTorrentStreamStats(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hashString := new(string)
	err := json.Unmarshal([]byte(input), hashString)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall torrent-stream-stat: %s", err.Error()),
				fmt.Errorf("Error unmarshall torrent-stream-stat: %s", err.Error()),
			},
			invalidParameters
	}
	stats, err := apiProvider.Provider.TorrentClientInterface.GetTorrentFileMetaDataChunked(*hashString)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Failed to get torrent stream stats: %s", err.Error()),
				err,
			},
			customError
	}
	retVal := new(interface{})
	*retVal = stats
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetChannel(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hash := new(primitives.Hash)
	err := json.Unmarshal(input, hash)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall get-channel: %s", err.Error()),
				fmt.Errorf("Error unmarshall get-channel: %s", err.Error()),
			},
			invalidParameters
	}
	channel, err := apiProvider.Provider.GetChannel(hash.String())
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error channel not found: %s", err.Error()),
				fmt.Errorf("Error channel not found: %s", err.Error()),
			},
			customError
	}
	if channel == nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error channel hash not found: %s", err.Error()),
				fmt.Errorf("Error channel hash not found: %s", err.Error()),
			},
			customError
	}
	retVal := new(interface{})
	*retVal = channel.ToCustomMarsalStruct()
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetChannels(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hashList := new(primitives.HashList)
	err := json.Unmarshal(input, hashList)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall get-channels: %s", err.Error()),
				fmt.Errorf("Error unmarshall get-channels: %s", err.Error()),
			},
			invalidParameters
	}
	channelList := make([]common.Channel, 0)
	for _, channelHash := range hashList.GetHashes() {
		channel, err := apiProvider.Provider.GetChannel(channelHash.String())
		if err != nil {
			return nil,
				&util.ApiError{
					fmt.Errorf("Error channels not found: %s", err.Error()),
					fmt.Errorf("Error channels not found: %s", err.Error()),
				},
				customError
		}
		if channel == nil {
			continue
			// return nil, fmt.Errorf("Content %s not found", contentHash.String())
		}
		channelList = append(channelList, *channel)
	}
	if len(channelList) == 0 {
		return nil,
			&util.ApiError{
				fmt.Errorf("No channels found by those hashes"),
				fmt.Errorf("No channels found by those hashes"),
			},
			customError
	}
	retVal := new(interface{})
	*retVal = channelList
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetContent(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hash := new(primitives.Hash)
	err := json.Unmarshal(input, hash)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall get-content: %s", err.Error()),
				fmt.Errorf("Error unmarshall get-content: %s", err.Error()),
			},
			invalidParameters
	}
	content, err := apiProvider.Provider.GetContent(hash.String())
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error content not found: %s", err.Error()),
				fmt.Errorf("Error content not found: %s", err.Error()),
			},
			customError
	}
	if content == nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error content hash not found: %s", err.Error()),
				fmt.Errorf("Error content hash not found: %s", err.Error()),
			},
			customError
	}
	retVal := new(interface{})
	*retVal = content
	return retVal, nil, noError
}

func (apiProvider ApiProvider) GetContents(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hashList := new(primitives.HashList)
	err := json.Unmarshal(input, hashList)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall get-contents: %s", err.Error()),
				fmt.Errorf("Error unmarshall get-contents: %s", err.Error()),
			},
			invalidParameters
	}

	contentList := make([]common.Content, 0)
	for _, contentHash := range hashList.GetHashes() {
		content, err := apiProvider.Provider.GetContent(contentHash.String())
		if err != nil {
			return nil,
				&util.ApiError{
					fmt.Errorf("Error contents not found: %s", err.Error()),
					fmt.Errorf("Error contents not found: %s", err.Error()),
				},
				customError
		}
		if content == nil {
			continue
		}
		contentList = append(contentList, *content)
	}
	if len(contentList) == 0 {
		return nil,
			&util.ApiError{
				fmt.Errorf("No contents found by those hashes"),
				fmt.Errorf("No contents found by those hashes"),
			},
			customError
	}
	cList := common.ContentList{
		ContentList: contentList,
	}
	retVal := new(interface{})
	*retVal = cList
	return retVal, nil, noError
}

func (apiProvider ApiProvider) VerifyChannel(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	verifyChannel := new(common.Channel)
	err := json.Unmarshal(input, verifyChannel)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall verify-channel: %s", err.Error()),
				fmt.Errorf("Error unmarshall verify-channel: %s", err.Error()),
			},
			invalidParameters
	}

	cost, apiError := apiProvider.Provider.CreationTool.VerifyChannel(verifyChannel)
	if apiError != nil {
		return nil, apiError, customError
	}
	retVal := new(interface{})
	*retVal = cost
	return retVal, nil, noError
}

func (apiProvider ApiProvider) CreateChannel(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	verifiedChannel := new(common.Channel)
	err := json.Unmarshal(input, verifiedChannel)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall create-channel: %s", err.Error()),
				fmt.Errorf("Error unmarshall create-channel: %s", err.Error()),
			},
			invalidParameters
	}

	hash, apiError := apiProvider.Provider.CreationTool.InitiateChannel(verifiedChannel)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in create-channel: initating channel: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in create-channel: initating channel: %s", apiError.UserError.Error()),
			},
			customError
	}

	verifiedChannel.RootChainID = *hash

	apiError = apiProvider.Provider.CreationTool.UpdateChannel(verifiedChannel)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in create-channel: updating channel: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in create-channel: updating channel: %s", apiError.UserError.Error()),
			},
			customError
	}

	retVal := new(interface{})
	*retVal = verifiedChannel
	return retVal, nil, noError
}

func (apiProvider ApiProvider) UpdateChannel(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	channel := new(common.Channel)
	err := json.Unmarshal(input, channel)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall update-channel: %s", err.Error()),
				fmt.Errorf("Error unmarshall update-channel: %s", err.Error()),
			},
			invalidParameters
	}

	apiError = apiProvider.Provider.CreationTool.UpdateChannel(channel)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in update-channel: updating channel: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in update-channel: updating channel: %s", apiError.UserError.Error()),
			},
			customError
	}

	retVal := new(interface{})
	*retVal = channel
	return retVal, nil, noError
}

func (apiProvider ApiProvider) AddExistingChannel(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	pk := new(primitives.PublicKey)
	err := json.Unmarshal(input, pk)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall add-existing-channel: %s", err.Error()),
				fmt.Errorf("Error unmarshall add-existing-channel: %s", err.Error()),
			},
			invalidParameters
	}

	apiError = apiProvider.Provider.CreationTool.AddExistingChannel(pk)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in add-existing-channel: deleting channel: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in add-existing-channel: deleting channel: %s", apiError.UserError.Error()),
			},
			customError
	}

	s := "Success"
	retVal := new(interface{})
	*retVal = s
	return retVal, nil, noError
}

func (apiProvider ApiProvider) DeleteChannel(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hash := new(primitives.Hash)
	err := json.Unmarshal(input, hash)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall delete-channel: %s", err.Error()),
				fmt.Errorf("Error unmarshall delete-channel: %s", err.Error()),
			},
			invalidParameters
	}

	apiError = apiProvider.Provider.CreationTool.DeleteChannel(hash)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in delete-channel: deleting channel: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in delete-channel: deleting channel: %s", apiError.UserError.Error()),
			},
			customError
	}

	s := "Success"
	retVal := new(interface{})
	*retVal = s
	return retVal, nil, noError
}

func (apiProvider ApiProvider) VerifyContent(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	verifyContent := new(common.Content)
	err := json.Unmarshal(input, verifyContent)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall verify-content: %s", err.Error()),
				fmt.Errorf("Error unmarshall verify-content: %s", err.Error()),
			},
			invalidParameters
	}

	cost, apiError := apiProvider.Provider.CreationTool.VerifyContent(verifyContent)
	if apiError != nil {
		return nil, apiError, customError
	}
	retVal := new(interface{})
	*retVal = cost
	return retVal, nil, noError
}

func (apiProvider ApiProvider) CreateContent(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	verifiedContent := new(common.Content)
	err := json.Unmarshal(input, verifiedContent)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall create-content: %s", err.Error()),
				fmt.Errorf("Error unmarshall create-content: %s", err.Error()),
			},
			invalidParameters
	}

	hash, apiError := apiProvider.Provider.CreationTool.AddContent(verifiedContent)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in create-content: creating content: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in create-content: creating content: %s", apiError.UserError.Error()),
			},
			customError
	}

	verifiedContent.ContentID = *hash

	retVal := new(interface{})
	*retVal = verifiedContent
	return retVal, nil, noError
}

func (apiProvider ApiProvider) DeleteContent(input json.RawMessage) (successResponse *interface{}, apiError *util.ApiError, errorType uint8) {
	hash := new(primitives.Hash)
	err := json.Unmarshal(input, hash)
	if err != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error unmarshall delete-content: %s", err.Error()),
				fmt.Errorf("Error unmarshall delete-content: %s", err.Error()),
			},
			invalidParameters
	}

	apiError = apiProvider.Provider.CreationTool.DeleteContent(hash)
	if apiError != nil {
		return nil,
			&util.ApiError{
				fmt.Errorf("Error in delete-channel: deleting content: %s", apiError.LogError.Error()),
				fmt.Errorf("Error in delete-channel: deleting content: %s", apiError.UserError.Error()),
			},
			customError
	}

	s := "Success"
	retVal := new(interface{})
	*retVal = s
	return retVal, nil, noError
}

func (apiService *ApiService) GetCompleteHeight() (uint32, error) {
	return apiService.Provider.GetCompleteHeight()
}
