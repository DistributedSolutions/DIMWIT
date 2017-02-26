package api

import (
	"log"
	"net/http"
)

type HelloArgs struct {
	Who string
}

type HelloReply struct {
	Message string
}

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
	log.Printf(args.Who)
	reply.Message = "Hello, " + args.Who + "!"
	log.Printf(reply.Message)
	return nil
}
