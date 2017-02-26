package provider

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/DistributedSolutions/DIMWIT/provider/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	// Start Server and codec for JSON-RPC 2.0
	jsonRPC := rpc.NewServer()
	jsonCodec := json.NewCodec()
	jsonRPC.RegisterCodec(jsonCodec, "application/json")
	jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8")

	// Api's Available
	jsonRPC.RegisterService(new(api.HelloService), "")
	jsonRPC.RegisterService(new(api.Arith), "")

	r.Handle("/api", jsonRPC)
	return r
}

func ServeRouter(r *mux.Router) io.Closer {
	port := ":8080"
	log.Println("Serving API on localhost" + port)
	closer, err := ListenAndServeWithClose(port, r)
	if err != nil {
		panic(err)
	}
	return closer
}

func ListenAndServeWithClose(addr string, handler http.Handler) (io.Closer, error) {
	var (
		listener  net.Listener
		srvCloser io.Closer
		err       error
	)

	srv := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		err := srv.Serve(listener.(*net.TCPListener))
		if err != nil {
			// log.Println("HTTP Server Error - ", err)
		}
		log.Printf("Closing API on %s...\n", addr)
	}()

	srvCloser = listener
	return srvCloser, nil
}
