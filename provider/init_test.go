package provider_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/provider"
	"github.com/DistributedSolutions/DIMWIT/testhelper"

	"github.com/Emyrk/test2doc/test"
	"github.com/adams-sarah/prettytest"
	// "github.com/gorilla/mux"
)

var PRINT_API_DOCS bool = false
var DataList []common.Channel
var router *http.ServeMux
var server *test.Server
var URL string

type mainSuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	var err error

	fake, dataList, err := testhelper.PopulateFakeClient(true, 5)
	if err != nil {
		t.Error(err)
	}
	DataList = dataList

	con, cache, err := testhelper.PopulateLevel2Cache(fake)
	if err != nil {
		t.Error(err)
	}
	defer con.Close()

	prov, err := provider.NewProvider(cache, fake)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Validate the data before we continue
	for _, c := range DataList {
		if cc, err := prov.GetChannel(c.RootChainID.String()); err != nil {
			t.Error(err)
			t.FailNow()
		} else {
			if !cc.IsSameAs(&c) {
				t.Error("Channel was not put into cache correctly")
			}
		}
	}

	router = prov.Router // provider.NewRouter()
	//router.KeepContext = true

	test.RegisterURLVarExtractor(provider.Vars)

	server, err = test.NewServer(router)
	if err != nil {
		panic(err.Error())
	}
	defer server.Finish()
	/*prov.Serve()
	defer prov.Close()
	URL = "http://localhost:8080"*/
	URL = server.URL

	fmt.Printf("Tests running on : %s\n", URL)
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(mainSuite),
	)
}
