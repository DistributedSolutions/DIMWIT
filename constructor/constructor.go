package constructor

import (
	"bytes"
	"fmt"

	"github.com/FactomProject/factom"
	//"github.com/DistributedSolutions/DIMWIT/common"
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
	Level2Cache database.IDatabase
	// Used per block
	ChannelCache map[string]objects.ChannelWrapper

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

	// Fresh map per block. This is not very efficient, but
	// optimization can come later
	c.ChannelCache = make(map[string]objects.ChannelWrapper)
	// Make a channel map, and get a batch apply map
	for _, e := range ents {
		c.ApplyEntryToCache(e)
	}

	var _ = ents
	return nil
}

// ApplyEntryToCache will take an entry, and apply it to the channels we have in our cache. If
// true is returned, this signals to the caller, we should also save the channel to the database.
func (c *Constructor) ApplyEntryToCache(e *lite.EntryHolder) (bool, error) {
	iae, err := objects.ParseFactomEntry(e)
	if err != nil {
		return false, err
	}

	// Almost all will request.
	chain, req := iae.RequestChannel()
	if req {
		cw, err := c.retrieveChannel(chain)
		if err != nil {
			return false, err
		}
		err = iae.AnswerChannelRequest(cw)
		if err != nil {
			return false, err
		}
	}

	hash, err := primitives.HexToHash(e.Entry.ChainID)
	if err != nil {
		return false, err
	}

	f := iae.NeedIsFirstEntry()
	if f {
		ent, err := c.Reader.GetFirstEntry(*hash)
		if err != nil {
			return false, err
		}

		same := lite.AreEntriesSame(ent, e.Entry)
		if !same {
			return false, fmt.Errorf("Entry needs to be the first entry, but it is not")
		}
	}

	ae := iae.NeedChainEntries()
	if ae {
		entries, err := c.Reader.GetAllChainEntries(*hash)
		if err != nil {
			return false, err
		}

		iae.AnswerChainEntries(converFactomEntriesToHolder(entries))
	}

	cw, wr := iae.ApplyEntry()
	if wr {
		c.ChannelCache[cw.Channel.RootChainID.String()] = *cw
	}

	return wr, nil
}

func converFactomEntriesToHolder(fents []*factom.Entry) []*lite.EntryHolder {
	holder := make([]*lite.EntryHolder, 0)
	for _, e := range fents {
		h := new(lite.EntryHolder)
		h.Entry = e
		h.Timestamp = 0
		h.Height = 0
		holder = append(holder, h)
	}
	return holder
}

// fixOrder puts channel instantiation before anything else
func fixOrder(ents []*lite.EntryHolder) []*lite.EntryHolder {
	pre := make([]*lite.EntryHolder, 0)    // Root chain
	middle := make([]*lite.EntryHolder, 0) // Other chains
	post := make([]*lite.EntryHolder, 0)   // Entries

	for _, e := range ents {
		if len(e.Entry.ExtIDs) < 2 {
			continue
		} else if bytes.Compare(e.Entry.ExtIDs[1], []byte("Channel Root Chain")) == 0 {
			pre = append(pre, e)
		} else if bytes.Compare(e.Entry.ExtIDs[1], []byte("Channel Management Chain")) == 0 ||
			bytes.Compare(e.Entry.ExtIDs[1], []byte("Channel Content Chain")) == 0 {
			middle = append(middle, e)
		} else {
			post = append(post, e)
		}
	}

	pre = append(pre, middle...)
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

func (c *Constructor) retrieveChannel(chainID string) (*objects.ChannelWrapper, error) {
	if cw, ok := c.ChannelCache[chainID]; ok {
		return &cw, nil
	}
	cid, err := primitives.HexToHash(chainID)
	if err != nil {
		return nil, err
	}
	return c.RetrieveChannel(*cid)
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
