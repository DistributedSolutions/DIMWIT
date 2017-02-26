package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/DistributedSolutions/DIMWIT/provider/api"
)

type Tester struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func (t *mainSuite) TestHelloService() {
	helloArgs := api.HelloArgs{"James"}
	jsonTest := Tester{JsonRpc: "2.0", Method: "HelloService.Say", Params: helloArgs}
	mjt, err := json.Marshal(jsonTest)
	if err != nil {
		fmt.Printf("ERROASDOFAOSDFOASODFO:%s", err.Error())
	}
	resp, err := http.Post(server.URL+"/api", "application/json", bytes.NewBuffer(mjt))
	if err != nil {
		fmt.Printf("Error: %s %s\n", err.Error(), resp.Body)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
	t.Must(t.Nil(err))

	t.Must(t.Equal(resp.StatusCode, http.StatusOK))

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var ws api.HelloReply
	err = decoder.Decode(&ws)
	t.Must(t.Nil(err))

	// t.Equal(len(ws), len(api.AllWidgets))
	// t.Must(t.True(len(ws) > 2))
	// fmt.Printf("YOOOOOO ")
	// t.Equal(ws[0].Id, api.AllWidgets[0].Id)
	// t.Equal(ws[2].Name, api.AllWidgets[2].Name)
	// t.Equal(ws[1].Role, api.AllWidgets[1].Role)
}
