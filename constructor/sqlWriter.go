package constructor

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/database"
)

const (
	LOOP_DELAY time.Duration = time.Duration(1 * time.Second)
)

type ISqlWriter interface {
	AddChannelArr(channels []common.Channel, height uint32) error
	FlushTempPlaylists(height uint32) error
	DeleteDBChannels() error
	DeleteDB() error
	Close()
}

type FakeSqlWriter struct{}

func (FakeSqlWriter) AddChannelArr(channels []common.Channel, height uint32) error { return nil }
func (FakeSqlWriter) FlushTempPlaylists(height uint32) error                       { return nil }
func (FakeSqlWriter) DeleteDBChannels() error                                      { return nil }
func (FakeSqlWriter) DeleteDB() error                                              { return nil }
func (FakeSqlWriter) Close()                                                       {}

type SqlWriter struct {
	// Incoming channels to write to sql db
	db *database.SqlDBWrapper
}

// Called to make SQLWriter
func NewSqlWriter() (*SqlWriter, error) {
	sw := new(SqlWriter)

	fmt.Printf("Init SqlWriter, Creating DB\n")
	db, err := database.CreateDB(constants.SQL_DB, database.CREATE_TABLE)
	if err != nil {
		return nil, fmt.Errorf("Error creating DB!! AAAHHH: %s", err)
	}
	sw.db = db

	err = sw.db.AddTags()
	if err != nil {
		return nil, fmt.Errorf("Error adding in tags: %s", err)
	}
	db.InitInterfaceDBVerification()
	fmt.Printf("Init SqlWriter, Finished Init\n")

	return sw, nil
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

func (sw *SqlWriter) DeleteDBChannels() error {
	return sw.db.DeleteDBChannels()
}

func (sw *SqlWriter) DeleteDB() error {
	sw.Close()
	err := database.DeleteDB(sw.db.Name)
	if err != nil {
		return fmt.Errorf("sqlWriter, error deleting db: %s", err.Error())
	}
	return nil
}
