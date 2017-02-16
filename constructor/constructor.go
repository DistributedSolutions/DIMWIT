package constructor

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor/objects"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var (
	CHANNEL_BUCKET []byte = []byte("Channels")
)

// Constructor builds the level 2 cache using factom-lite
type Constructor struct {
	Level2Cache  database.IDatabase
	ChannelCache map[string]common.Channel

	// State
	CompletedHeight uint32 // Height channels/content is updated to

	// Constructor only reads from Factom
	Reader lite.FactomLiteReader
}

func NewContructor(dbType string) (*Constructor, error) {
	c := new(Constructor)
	var db database.IDatabase
	switch dbType {
	case "Bolt":
		db = database.NewBoltDB(constants.HIDDEN_DIR + constants.LVL2_CACHE)
	case "LDB":
	case "Map":
	default:
		return nil, fmt.Errorf("DBType given not valid. Found %s, expected either: Bolt, Map, LDB", dbType)
	}

	c.Level2Cache = db
	c.LoadStateFromDB()
	return c, nil
}

// SetReader takes in a FactomLiteReader, this is where the contructor will get it's data
func (c *Constructor) SetReader(r lite.FactomLiteReader) {
	c.Reader = r
}

// LoadStateFromDB loads state information, such as last completed height. This mean's we don't have to parse
// through the blockchain again! Woot!
func (c *Constructor) LoadStateFromDB() {
	// TODO: Set current height to height in DB
}

func (c *Constructor) ApplyHeight(height uint32) error {
	rh, err := c.Reader.GetReadyHeight()
	if err != nil {
		return err
	}

	if height > rh {
		return nil
	}

	ents, err := c.Reader.GrabAllEntriesAtHeight(height)
	if err != nil {
		return err
	}

	ents = fixOrder(ents)

	// Make a map

	var _ = ents
	return nil
}

// ApplyEntryToCache will take an entry, and apply it to the channels we have in our cache. If
// true is returned, this signals to the caller, we should also save the channel to the database.
func (c *Constructor) ApplyEntryToCache(e *lite.EntryHolder) bool {

	return false
}

// fixOrder puts channel instantiation before anything else
func fixOrder(ents []*lite.EntryHolder) []*lite.EntryHolder {
	pre := make([]*lite.EntryHolder, 0)
	post := make([]*lite.EntryHolder, 0)

	for _, e := range ents {
		if len(e.Entry.ExtIDs) < 2 {
			continue
		} else if bytes.Compare(e.Entry.ExtIDs[1], []byte("Channel Root Chain")) == 0 {
			pre = append(pre, e)
		} else {
			post = append(post, e)
		}
	}

	return append(pre, post...)
}

// saveChannel overwrites whatever channel is in the cache with what it is given.
// be sure not to fuck up the data in the database. As you can see,I
func (c *Constructor) saveChannel(ch objects.ChannelWrapper) error {
	data, err := ch.MarshalBinary()
	if err != nil {
		return err
	}
	err = c.Level2Cache.Put(CHANNEL_BUCKET, ch.Channel.RootChainID.Bytes(), data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Constructor) RetrieveChannel(chainID primitives.Hash) (*objects.ChannelWrapper, error) {
	data, err := c.Level2Cache.Get(CHANNEL_BUCKET, chainID.Bytes())
	if err != nil {
		return nil, err
	}

	ch := new(objects.ChannelWrapper)
	if len(data) == 0 {

	} else {
		err = ch.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
	}
	return ch, nil
}