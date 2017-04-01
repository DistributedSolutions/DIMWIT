package provider

import (
	"fmt"
	"github.com/DistributedSolutions/logrus"
	"io"
	"log"
	"net"
	"net/http"
)

func NewRouter(srv *ApiService) *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/api", srv.HandleAPICalls)

	return r
}

func ServeRouter(r *http.ServeMux, portNumber int) io.Closer {
	port := fmt.Sprintf(":%d", portNumber)
	log.Println("Serving API on localhost" + port)
	closer, err := ListenAndServeWithClose(port, r)
	if err != nil {
		logrus.Warn("Oh no you got an error trying to open the server :( [" + err.Error() + "]")
		return nil
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
