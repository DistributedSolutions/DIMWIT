package constructor

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/constructor/objects"
	"github.com/DistributedSolutions/DIMWIT/database"
)

const (
	LOOP_DELAY time.Duration = time.Duration(1 * time.Second)
)

type SqlWriter struct {
	// Incoming channels to write to sql db
	channelQueue chan objects.ChannelWrapper
	db           *database.SqlDBWrapper
}

// Called to make SQLWriter
func NewSqlWriter() *SqlWriter {
	sw := new(SqlWriter)
	sw.channelQueue = make(chan objects.ChannelWrapper, 1000)

	fmt.Printf("Init SqlWriter, Creating DB\n")
	db, err := database.CreateDB(constants.SQL_DB, database.CREATE_TABLE)
	if err != nil {
		fmt.Printf("Error creating DB!! AAAHHH: %s", err)
	}
	err = sw.db.AddTags()
	if err != nil {
		fmt.Printf("Error adding in tags: %s", err)
	}
	fmt.Printf("Init SqlWriter, Finished Init\n")

	sw.db = db

	return sw
}

// Called to send a channel to the SQLWriter.
func (sw *SqlWriter) SendChannelDownQueue(c objects.ChannelWrapper) {
	// If you want to do anything to it before it hits the go-routine
	// do so here.
	// You have access to some extra variables with the Wrapper.
	// I don't think you care about them though.
	sw.channelQueue <- c
}

// Close sqlwriter
func (sw *SqlWriter) Close() {
	database.CloseDB(sw.db.DB)
}

func (sw *SqlWriter) AddChannelArr(channels []common.Channel, height uint32) error {
	return sw.db.AddChannelArr(channels, height)
}

func (sw *SqlWriter) FlushTempPlaylists(height uint32) error {
	return sw.db.FlushPlaylistTempTable(height)
}
