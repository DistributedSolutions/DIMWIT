package provider_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/provider"
	"github.com/DistributedSolutions/DIMWIT/testhelper"
	"github.com/DistributedSolutions/DIMWIT/writeHelper"

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

func TestMethods(t *testing.T) {
	baseMethod := reflect.TypeOf(provider.ApiBase{}).Method(0)
	st := reflect.TypeOf(provider.ApiProvider{})
	for i := 0; i < st.NumMethod(); i++ {

		if baseMethod.Type.NumIn() != st.Method(i).Type.NumIn() {
			t.Errorf("Error between [%s] and method [%s], args count: %d != %d respectively.\n\t\tDo not change base method... change your method.",
				baseMethod.Name,
				st.Method(i).Name,
				baseMethod.Type.NumIn(),
				st.Method(i).Type.NumIn())
		} else {
			for bi := 1; bi < baseMethod.Type.NumIn(); bi++ {
				if baseMethod.Type.In(bi) != st.Method(i).Type.In(bi) {
					t.Errorf("Error between [%s] and method [%s] args types!\n\t\targs type[%s] and args type[%s], do not match.\n\t\tDo not change base method... change your method.",
						baseMethod.Name,
						st.Method(i).Name,
						baseMethod.Type.In(bi),
						st.Method(i).Type.In(bi))
				}
			}
		}

		if baseMethod.Type.NumOut() != st.Method(i).Type.NumOut() {
			t.Errorf("Error between [%s] and method [%s], return count: %d != %d respectively.\n\t\tDo not change base method... change your method.",
				baseMethod.Name,
				st.Method(i).Name,
				baseMethod.Type.NumOut(),
				st.Method(i).Type.NumOut())
		} else {
			for ai := 0; ai < baseMethod.Type.NumOut(); ai++ {
				if baseMethod.Type.Out(ai) != st.Method(i).Type.Out(ai) {
					t.Errorf("Error between [%s] and method [%s] return types!\n\t\treturn type[%s] and return type[%s], do not match.\n\t\tDo not change base method... change your method.",
						baseMethod.Name,
						st.Method(i).Name,
						baseMethod.Type.Out(ai),
						st.Method(i).Type.Out(ai))
				}
			}
		}
	}
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

	w, err := writeHelper.NewWriterHelper(con, fake)
	if err != nil {
		t.Error(err)
	}

	prov, err := provider.NewProvider(cache, fake, w)
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
