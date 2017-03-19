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
	log "github.com/DistributedSolutions/logrus"
	"github.com/FactomProject/factom"
)

// Constructor builds the level 2 cache using factom-lite
type Constructor struct {
	Level2Cache database.IDatabase
	SqlGuy      *SqlWriter

	// Used per block
	ChannelCache map[string]objects.ChannelWrapper

	// State
	CompletedHeight uint32 // Height channels/content is updated to

	// Constructor only reads from Factom
	Reader lite.FactomLiteReader

	// Managing
	quit chan int
}

//dbType string,
func NewContructor(db database.IDatabase) (*Constructor, error) {
	c := new(Constructor)

	var err error
	c.SqlGuy, err = NewSqlWriter()
	if err != nil {
		return nil, err
	}
	c.Level2Cache = db
	c.loadStateFromDB()
	c.quit = make(chan int, 20)
	return c, nil
}

func (c *Constructor) InterruptClose() {
	err := c.Close()
	if err != nil {
		log.Println("Constructor failed to safely close: ", err.Error())
	} else {
		log.Println("Constructor closed safely")
	}
}

func (c *Constructor) Close() error {
	c.SqlGuy.Close()
	c.Kill() // Kill routine
	return c.Level2Cache.Close()
}

// SetReader takes in a FactomLiteReader, this is where the contructor will get it's data
func (c *Constructor) SetReader(r lite.FactomLiteReader) {
	c.Reader = r
}

// loadStateFromDB loads state information, such as last completed height. This mean's we don't have to parse
// through the blockchain again! Woot!
func (c *Constructor) loadStateFromDB() error {
	c.CompletedHeight = 0
	if c.Level2Cache != nil {
		data, err := c.Level2Cache.Get(constants.STATE_BUCKET, constants.STATE_COMP_HEIGHT)
		if err == nil {
			u, err := primitives.BytesToUint32(data)
			if err == nil {
				c.CompletedHeight = u
			}
		}

		// Either we encountered an error loading, or
		// we never had a database to begin with.
		if c.CompletedHeight == 0 {
			// TODO: Might want to check if data exists, if it does
			// then we have data, but the height is 0. This might cause issues
			err = c.Level2Cache.Put(constants.STATE_BUCKET, constants.STATE_COMP_HEIGHT, []byte{0x00, 0x00, 0x00, 0x00})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Constructor) ApplyHeight(height uint32) error {
	rh, err := c.Reader.GetReadyHeight()
	if err != nil {
		log.Error("Constructor failed to get ready height from Factom Lite")
		return err
	}

	if height > rh {
		return fmt.Errorf("Cannot apply height %d, the ready height is %d", height, rh)
	}

	ents, err := c.Reader.GrabAllEntriesAtHeight(height)
	if err != nil {
		log.Error("Constructor failed to grab entries from factom client")
		return err
	}

	ents = fixOrder(ents)

	// Fresh map per block. This is not very efficient, but
	// optimization can come later
	c.ChannelCache = make(map[string]objects.ChannelWrapper)
	// Make a channel map, and get a batch apply map
	for _, e := range ents {
		_, err := c.applyEntryToCache(e)
		if err != nil {
			log.Debug(err)
		}
	}

	chanList := make([]common.Channel, 0)

	// Level 2 Cache Write
	for _, channel := range c.ChannelCache {
		channel.CurrentHeight = c.CompletedHeight
		data, err := channel.MarshalBinary()
		if err != nil {
			continue
		}

		if channel.Channel.Status() >= constants.CHANNEL_READY {
			chanList = append(chanList, channel.Channel)
		}

		//fmt.Println(channel.Channel.RootChainID.String())
		err = c.Level2Cache.Put(constants.CHANNEL_BUCKET, channel.Channel.RootChainID.Bytes(), data)
		if err != nil {
			continue
		}

		for _, content := range channel.Channel.Content.GetContents() {
			data, err = content.MarshalBinary()
			if err != nil {
				continue
			}

			err = c.Level2Cache.Put(constants.CONTENT_BUCKET, content.ContentID.Bytes(), data)
			if err != nil {
				continue
			}
		}
	}

	log.SetLevel(log.InfoLevel)
	if len(chanList) > 0 {
		log.Debugf("DEBUG: Executing %d Channels", len(chanList))
		count := 0
		for _, c := range chanList {
			count += len(c.Tags.GetTags())
		}
		log.Debugf("DEBUG: chanList total tags length are: %d", count)
	}
	log.SetLevel(log.DebugLevel)

	// Write to SQL
	err = c.SqlGuy.AddChannelArr(chanList, height)
	if err != nil {
		return err
	}

	err = c.SqlGuy.FlushTempPlaylists(height)
	if err != nil {
		return err
	}

	log.Debugf("Finished applying height: %d", height)

	// Update State
	c.CompletedHeight = height
	err = c.Level2Cache.Put(constants.STATE_BUCKET, constants.STATE_COMP_HEIGHT, primitives.Uint32ToBytes(height))
	c.ChannelCache = nil
	return err
}

// applyEntryToCache will take an entry, and apply it to the channels we have in our cache. If
// true is returned, this signals to the caller, we should also save the channel to the database.
// This is where the magic happens
func (c *Constructor) applyEntryToCache(e *lite.EntryHolder) (bool, error) {
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
		end := len(entries)
		for i := 0; i < end; i++ {
			if bytes.Compare(entries[i].ExtIDs[2], first.ExtIDs[2]) == 0 {
				//entries = append(entries[:i], entries[i+1:]...)
				entries[i] = entries[len(entries)-1]
				entries = entries[:len(entries)-1]
				if i > 0 {
					i--
				}
				end = end - 1
			}
		}

		iae.AnswerChainEntriesInOther(convertFactomEntryToHolder(first), convertFactomEntriesToHolder(entries))
	}

	// The iae has everything it needs, let's see what it decided
	cw, wr := iae.ApplyEntry()
	if wr {
		// Ok, it told us to write this to the db. Let's put it in the map for a write
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
	err = c.Level2Cache.Put(constants.CHANNEL_BUCKET, ch.Channel.RootChainID.Bytes(), data)
	if err != nil {
		return err
	}
	return nil
}

// retrieveChannel will try to retrieve from our local map first. The local map is the channels
// we are currently working with before a write. This is only used internally
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
	data, err := c.Level2Cache.Get(constants.CHANNEL_BUCKET, chainID.Bytes())
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
