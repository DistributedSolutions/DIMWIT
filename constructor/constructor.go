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
		db = database.NewMapDB()
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
		return fmt.Errorf("Cannot apply height %d, the ready height is %d", height, rh)
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

	// TODO: Batch write
	for _, channel := range c.ChannelCache {
		data, err := channel.MarshalBinary()
		if err != nil {
			continue
		}

		err = c.Level2Cache.Put(CHANNEL_BUCKET, channel.Channel.RootChainID.Bytes(), data)
		if err != nil {
			continue
		}
	}
	c.CompletedHeight = height
	return nil
}

// ApplyEntryToCache will take an entry, and apply it to the channels we have in our cache. If
// true is returned, this signals to the caller, we should also save the channel to the database.
// This is where the magic happens
func (c *Constructor) ApplyEntryToCache(e *lite.EntryHolder) (bool, error) {
	// Instantiate the IApplyEntry
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

	// If the iae requires it to be the first entry, we need to
	// error out if it is not
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

	// Feed it all of the entries in it's chain
	ae := iae.NeedChainEntries()
	if ae {
		entries, err := c.Reader.GetAllChainEntries(*hash)
		if err != nil {
			return false, err
		}

		iae.AnswerChainEntries(convertFactomEntriesToHolder(entries))
	}

	// Do we need another chain's entries? Usually a Content Link
	linkChain, noe := iae.RequestEntriesInOtherChain()
	if noe {
		ohash, err := primitives.HexToHash(linkChain)
		if err != nil {
			return false, err
		}

		entries, err := c.Reader.GetAllChainEntries(*ohash)
		if err != nil {
			return false, err
		}

		// Discard any element with the first entry method.
		// There should only be on with that method, so anything else
		// could be spam or malicious
		first, err := c.Reader.GetFirstEntry(*ohash)
		for i := range entries {
			if bytes.Compare(entries[i].ExtIDs[2], first.ExtIDs[2]) == 0 {
				entries = append(entries[:i], entries[i+1:]...)
			}
		}

		iae.AnswerChainEntriesInOther(convertFactomEntryToHolder(first), convertFactomEntriesToHolder(entries))
	}

	// The iae has everything it needs, let's see what it decided
	cw, wr := iae.ApplyEntry()
	// fmt.Println(iae.String(), " -- ", wr)
	if wr {
		// Ok, it told us to write this to the db. Let's put it in the map for a batch write
		c.ChannelCache[cw.Channel.RootChainID.String()] = *cw
	}

	return wr, nil
}

func convertFactomEntriesToHolder(fents []*factom.Entry) []*lite.EntryHolder {
	holder := make([]*lite.EntryHolder, 0)
	for _, e := range fents {
		holder = append(holder, convertFactomEntryToHolder(e))
	}
	return holder
}

func convertFactomEntryToHolder(fent *factom.Entry) *lite.EntryHolder {
	h := new(lite.EntryHolder)
	h.Entry = fent
	h.Timestamp = 0
	h.Height = 0
	return h
}

// fixOrder puts channel instantiation before anything else, then first entries, then following entries
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
		} else if bytes.Compare(e.Entry.ExtIDs[1], []byte("Content Signing Key")) == 0 {
			middle = append(middle, e) // Root chain is in pre, we good to put this in middle
		} else {
			post = append(post, e)
		}
	}

	pre = append(pre, middle...)
	return append(pre, post...)
}

// saveChannel overwrites whatever channel is in the cache with what it is given.
// be sure not to fuck up the data in the database. As you can see, I made it private
// to keep people from writing to my cache!
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

// retrieveChannel will try to retrieve from our local map first. The local map is the channels
// we are currently working with before a batch write. This is only used internally
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

// RetrieveChannel retrieves the channel from the Level2Cache
func (c *Constructor) RetrieveChannel(chainID primitives.Hash) (*objects.ChannelWrapper, error) {
	data, err := c.Level2Cache.Get(CHANNEL_BUCKET, chainID.Bytes())
	if err != nil {
		return nil, err
	}

	ch := objects.NewChannelWrapper()
	if len(data) == 0 {

	} else {
		err = ch.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
	}
	return ch, nil
}
