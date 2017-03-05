package provider

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor/objects"
	"github.com/DistributedSolutions/DIMWIT/database"
)

type Provider struct {
	Level2Cache database.IDatabase

	// API
	Router    *http.ServeMux
	Service   *ApiService
	Salt      string // Identify provider. Nice for tests
	apicloser io.Closer
}

func NewProvider(db database.IDatabase) (*Provider, error) {
	p := new(Provider)
	p.Level2Cache = db
	randData := make([]byte, 30)
	rand.Read(randData)
	hash := sha256.Sum256(randData[:])
	p.Salt = hex.EncodeToString(hash[:])

	p.Service = new(ApiService)
	p.Service.Provider = p
	p.Router = NewRouter(p.Service)

	return p, nil
}

func (p *Provider) Serve() {
	closer := ServeRouter(p.Router)
	p.apicloser = closer
}

func (p *Provider) Close() {
	p.apicloser.Close()
}

// Horribly inefficient with large data sets. Need to cache
// this data
type DatabaseStats struct {
	TotalChannels int `json:"totalchannels"`
	TotalContent  int `json:"totalcontent"`
}

func (p *Provider) GetStats() (*DatabaseStats, error) {
	ds := new(DatabaseStats)
	chans, err := p.GetAllChannels()
	if err != nil {
		return nil, err
	}

	ds.TotalChannels = len(chans)
	totalCon := 0
	for _, c := range chans {
		totalCon += len(c.Content.GetContents())
	}
	ds.TotalContent = totalCon
	return ds, nil
}

func (p *Provider) GetAllChannels() ([]common.Channel, error) {
	_, data, err := p.Level2Cache.GetAll(constants.CHANNEL_BUCKET)
	if err != nil {
		return nil, err
	}

	ret := make([]common.Channel, 0)
	for _, d := range data {
		cw := objects.NewChannelWrapper()
		err = cw.UnmarshalBinary(d)
		if err != nil {
			continue
		}

		ret = append(ret, cw.Channel)
	}

	return ret, nil
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

	if len(data) == 0 {
		return nil, nil
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

	if len(data) == 0 {
		return nil, nil
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
