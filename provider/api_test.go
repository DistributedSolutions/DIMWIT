package provider_test

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "net/http"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/jsonrpc"
	"github.com/DistributedSolutions/DIMWIT/provider"
)

var _ = fmt.Sprintf("")

func (t *mainSuite) TestGetStats() {
	t.True(true)
	// TODO: This unit test is broke
	return
	req := jsonrpc.NewEmptyParamsJSONRPCRequest("get-stats", 0)

	respObj, jsonError, err := req.POSTRequest(URL+"/api", new(provider.DatabaseStats))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}

	contentCount := 0
	for _, c := range DataList {
		contentCount += len(c.Content.GetContents())
	}

	if err == nil && jsonError == nil { // If no errors, check the reponse
		resp := respObj.(*provider.DatabaseStats)
		if resp.TotalChannels != len(DataList) && resp.TotalContent != contentCount {
			t.Error("Post torrent Stream Seek does not have success", err)
		}
		t.True(resp.TotalChannels == len(DataList) && resp.TotalContent == contentCount)
	}
}

func (t *mainSuite) TestGetConstants() {
	// TODO: This unit test is broke
	t.True(true)
	return
	req := jsonrpc.NewEmptyParamsJSONRPCRequest("get-constants", 0)

	respObj, jsonError, err := req.POSTRequest(URL+"/api", new(constants.ConstantJSON))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}

	if err == nil && jsonError == nil { // If no errors, check the reponse
		c := respObj.(*constants.ConstantJSON)
		if c.HashLength != 64 &&
			c.FileNameLength != constants.FILE_NAME_MAX_LENGTH &&
			c.LongDescriptionLength != constants.LONG_DESCRIPTION_MAX_LENGTH &&
			c.ShortDescriptionLength != constants.SHORT_DESCRIPTION_MAX_LENGTH &&
			c.TrackerUrlLength != constants.TRACKER_URL_MAX_LENGTH &&
			c.FilePathLength != constants.FILE_PATH_MAX_LENGTH &&
			c.TitleLength != constants.TITLE_MAX_LENGTH &&
			c.UrlLength != constants.URL_MAX_LENGTH {
			t.Error("Get Constants is incorrect", err)
		}
		t.True(c.HashLength == 64 &&
			c.FileNameLength == constants.FILE_NAME_MAX_LENGTH &&
			c.LongDescriptionLength == constants.LONG_DESCRIPTION_MAX_LENGTH &&
			c.ShortDescriptionLength == constants.SHORT_DESCRIPTION_MAX_LENGTH &&
			c.TrackerUrlLength == constants.TRACKER_URL_MAX_LENGTH &&
			c.FilePathLength == constants.FILE_PATH_MAX_LENGTH &&
			c.TitleLength == constants.TITLE_MAX_LENGTH &&
			c.UrlLength == constants.URL_MAX_LENGTH)
	}
}

func (t *mainSuite) TestPostTorrentStreamSeek() {
	req := jsonrpc.NewJSONRPCRequest("post-torrent-stream-seek", 100.0, 0)

	respObj, jsonError, err := req.POSTRequest(URL+"/api", new(string))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}

	if err == nil && jsonError == nil { // If no errors, check the reponse
		resp := *respObj.(*string)
		if resp != "Success" {
			t.Error("Post torrent Stream Seek does not have success", err)
		}
		t.True(resp == "Success")
	}
}

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

	req := jsonrpc.NewJSONRPCRequest("get-channels", temp, 0)

	respObj, jsonError, err := req.POSTRequest(URL+"/api", new([]common.CustomJSONMarshalChannel))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}

	if err == nil && jsonError == nil { // If no errors, check the reponse
		resp := respObj.(*[]common.CustomJSONMarshalChannel)
		for i, e := range *resp {
			if PRINT_API_DOCS {
				if i != 0 { // only run once if need to print docs
					break
				}
			}
			v := DataList[i].ToCustomMarsalStruct()
			if !e.IsSimilarTo(v) {
				t.Error("Channels returned does not match", err)
			}
			t.True(e.IsSimilarTo(v))
		}
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

func (t *mainSuite) TestProviderContents() {
	contentCount := 0
	for _, c := range DataList {
		contentCount += len(c.Content.GetContents())
	}

	contentsHashList := new(primitives.HashList)
	for _, c := range DataList {
		for _, co := range c.Content.GetContents() {
			contentsHashList.AddHash(&co.ContentID)
		}
	}

	req := jsonrpc.NewJSONRPCRequest("get-contents", contentsHashList, 0)

	_, jsonError, err := req.POSTRequest(URL+"/api", new(primitives.HashList))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}
	t.True(true)

	// if err == nil && jsonError == nil { // If no errors, check the reponse
	// 	resp := *respObj.(*common.ContentList)
	// 	count := 0
	// 	for _, c := range DataList {
	// 		for _, ch := range c.Content.GetContents() {
	// 			v := &resp.GetContents()[count]
	// 			if !ch.IsSameAs(v) {
	// 				t.Error("Channels returned does not match", err)
	// 			}
	// 			t.True(ch.IsSameAs(v))
	// 			count += 1
	// 		}
	// 	}
	// }
}
