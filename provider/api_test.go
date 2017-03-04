package provider_test

import (
	"fmt"
	//"io/ioutil"
	//"net/http"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/provider/jsonrpc"
)

var _ = fmt.Sprintf("")

func (t *mainSuite) TestProviderChannel() {
	for _, c := range DataList {
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
		}
	}
}
