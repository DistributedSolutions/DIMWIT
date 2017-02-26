package provider

import (
	"io"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor/objects"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/gorilla/mux"
)

type Provider struct {
	Level2Cache database.IDatabase

	// API
	Router    *mux.Router
	apicloser io.Closer
}

func NewProvider(db database.IDatabase) (*Provider, error) {
	p := new(Provider)
	p.Level2Cache = db
	p.Router = NewRouter()

	return p, nil
}

func (p *Provider) Serve() {
	closer := ServeRouter(p.Router)
	p.apicloser = closer
}

func (p *Provider) Close() {
	p.apicloser.Close()
}

func (p *Provider) GetChannel(channelID string) (*common.Channel, error) {
	key, err := primitives.HexToHash(channelID) // hex.DecodeString(channelID)
	if err != nil {
		return nil, err
	}

	data, err := p.Level2Cache.Get(constants.CHANNEL_BUCKET, key.Bytes())
	if err != nil {
		return nil, err
	}

	cw := objects.NewChannelWrapper()
	err = cw.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return &cw.Channel, nil
}

func (p *Provider) GetContent(contentID string) (*common.Content, error) {
	key, err := primitives.HexToHash(contentID) // hex.DecodeString(channelID)
	if err != nil {
		return nil, err
	}

	data, err := p.Level2Cache.Get(constants.CONTENT_BUCKET, key.Bytes())
	if err != nil {
		return nil, err
	}

	con := common.NewContent()
	err = con.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return con, nil
}

func (p *Provider) GetCompleteHeight() (uint32, error) {
	data, err := p.Level2Cache.Get(constants.STATE_BUCKET, constants.STATE_COMP_HEIGHT)
	if err != nil {
		return 0, err
	}

	u, err := primitives.BytesToUint32(data)
	if err != nil {
		return 0, err
	}
	return u, nil
}
