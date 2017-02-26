package api_test

import (
	"fmt"
	"testing"

	"github.com/DistributedSolutions/DIMWIT/provider"

	"github.com/adams-sarah/prettytest"
	"github.com/adams-sarah/test2doc/test"
	"github.com/gorilla/mux"
)

var router *mux.Router
var server *test.Server

type mainSuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	var err error

	router = provider.NewRouter()
	router.KeepContext = true

	test.RegisterURLVarExtractor(mux.Vars)

	server, err = test.NewServer(router)
	if err != nil {
		panic(err.Error())
	}
	defer server.Finish()

	fmt.Printf("Tests running on : %s\n", server.URL)
	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(mainSuite),
	)
}
