package api_test

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"

	"github.com/DistributedSolutions/DIMWIT/provider/api"
	"github.com/DistributedSolutions/DIMWIT/provider/jsonrpc"
)

func (t *mainSuite) TestHelloService() {
	// Make Args
	helloArgs := api.HelloArgs{"James"}

	// Pack into JsonRPC
	req := jsonrpc.NewJSONRPCRequest("HelloService.Say", helloArgs, 0)

	respObj, jsonError, err := req.POSTRequest(server.URL+"/api", new(api.HelloReply))
	if err != nil { // Go Error
		t.Error(err)
	}
	if jsonError != nil { // Error in json response
		t.Error(jsonError.Message)
	}

	if err == nil && jsonError == nil { // If no errors, check the reponse
		resp := respObj.(*api.HelloReply)
		fmt.Println(resp.Message) // The reponse!
	}
}
