package provider

import (
	"github.com/DistributedSolutions/DIMWIT/provider/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	// "log"
	// "net/http"
)

// func main() {
// 	r := mux.NewRouter()
// 	jsonRPC := rpc.NewServer()
// 	jsonCodec := json.NewCodec()
// 	jsonRPC.RegisterCodec(jsonCodec, "application/json")
// 	jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8")
// 	jsonRPC.RegisterService(new(api.HelloService), "")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	// r.HandleFunc("/widgets", api.GetWidgets).Methods("GET").Name("GetWidgets")
	// r.HandleFunc("/widgets", api.PostWidget).Methods("POST").Name("PostWidget")
	// r.HandleFunc("/widgets/{id}", api.GetWidget).Methods("GET").Name("GetWidget")
	jsonRPC := rpc.NewServer()
	jsonCodec := json.NewCodec()
	jsonRPC.RegisterCodec(jsonCodec, "application/json")
	jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8")
	jsonRPC.RegisterService(new(api.HelloService), "")
	r.Handle("/api", jsonRPC)

	return r
}
