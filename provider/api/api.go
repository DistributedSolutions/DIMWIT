package api

import (
	"errors"
	"log"
	"net/http"
)

type HelloArgs struct {
	Who string `json: "who"`
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

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(r *http.Request, args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(r *http.Request, args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
