package provider

import (
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
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

type ChannelsReturn struct {
}

var TEST_CONST = "TEST_CONST"

type EmptyRequest struct {
}

type ApiService struct {
	Provider Provider
}

func (apiService *ApiService) GetChannels(r *http.Request, hashList *primitives.HashList, reply *common.ChannelList) error {
	for i, channelHash := range hashList.List {
		channel, err := apiService.Provider.GetChannel(channelHash.String())
		if err != nil {
			return err
		}
		reply.List[i] = *channel
	}
	return nil
}

func (apiService *ApiService) GetContent(r *http.Request, hash *primitives.Hash, reply *common.Content) error {
	content, err := apiService.Provider.GetContent(hash.String())
	if err != nil {
		return err
	}
	reply = content
	return nil
}

func (apiService *ApiService) GetCompleteHeight(r *http.Request, empty *EmptyRequest, reply *uint32) error {
	height, err := apiService.Provider.GetCompleteHeight()
	if err != nil {
		return err
	}
	reply = &height
	return nil
}
