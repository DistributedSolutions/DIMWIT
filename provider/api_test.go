package provider_test

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/provider"
	"github.com/DistributedSolutions/DIMWIT/provider/jsonrpc"
)

var _ = fmt.Sprintf("")

func (t *mainSuite) TestProviderChannel() {
	for i, c := range DataList {
		if PRINT_API_DOCS {
			if i != 0 { // only run once if need to print docs
				break
			}
		}
		req := jsonrpc.NewJSONRPCRequest("get-channel", c.RootChainID.String(), 0)

		respObj, jsonError, err := req.POSTRequest(URL+"/api", new(common.CustomJSONMarshalChannel))
		if err != nil { // Go Error
			t.Error(err)
		}
		if jsonError != nil { // Error in json response
			t.Error(jsonError.Message)
		}

		if err == nil && jsonError == nil { // If no errors, check the reponse
			resp := respObj.(*common.CustomJSONMarshalChannel)
			if !resp.IsSimilarTo(c.ToCustomMarsalStruct()) {
				t.Error("Channel returned does not match", err)
			}
			t.True(resp.IsSimilarTo(c.ToCustomMarsalStruct()))
		}

	}
}

func (t *mainSuite) TestProviderChannels() {
	channelsHashList := make([]primitives.Hash, len(DataList))
	for i, c := range DataList {
		channelsHashList[i] = c.RootChainID
	}

	temp := primitives.HashList{
		List: channelsHashList,
	}

	jsonChannelsHashList, _ := json.Marshal(temp)
	req := jsonrpc.NewJSONRPCRequest("get-channels", string(jsonChannelsHashList), 0)

	channelList := common.ChannelList{
		List: DataList,
	}
	respObj, jsonError, err := req.POSTRequest(URL+"/api", new(common.ChannelList))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}

	if err == nil && jsonError == nil { // If no errors, check the reponse
		resp := respObj.(*common.ChannelList)
		if !resp.IsSimilarTo(channelList) {
			t.Error("Channels returned does not match", err)
		}
		t.True(resp.IsSimilarTo(channelList))
	}
}

func (t *mainSuite) TestProviderContent() {
	done := false
	for _, c := range DataList {
		if PRINT_API_DOCS {
			if done { // only run once if need to print docs
				break
			}
		}
		if len(c.Content.ContentList) == 0 {
			continue
		}
		done = true
		con := c.Content.ContentList[0]
		req := jsonrpc.NewJSONRPCRequest("get-content", con.ContentID.String(), 0)

		respObj, jsonError, err := req.POSTRequest(URL+"/api", new(common.Content))
		if err != nil { // Go Error
			t.Error(err)
		}
		if jsonError != nil { // Error in json response
			t.Error(jsonError.Message)
		}

		if err == nil && jsonError == nil { // If no errors, check the reponse
			resp := respObj.(*common.Content)
			// fmt.Println(con, "\n\n\n", resp)
			resp.CreationTime = con.CreationTime
			if !resp.IsSameAs(&con) {
				t.Error("Content returned does not match. Error?:", err)
			}
			t.True(resp.IsSameAs(&con))
		}
	}
}

func (t *mainSuite) TestGetStats() {
	req := jsonrpc.NewJSONRPCRequest("get-stats", "", 0)
	respObj, jsonError, err := req.POSTRequest(URL+"/api", new(provider.DatabaseStats))
	t.Nil(err)
	if jsonError != nil {
		t.Error(jsonError.Message)
	}

	resp := respObj.(*provider.DatabaseStats)
	if resp == nil {
		t.Error("Why is this nil")
	}
}
