package database

import (
	"fmt"
	"github.com/mattn/go-sqlite3"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	// "github.com/DistributedSolutions/logrus"
)

type PlayListTempStoreEntry struct {
	Title     string
	ChannelId string
	ContentId string
	Id        string
}

func (sqlDB *SqlDBWrapper) DeleteTags() error {
	db := sqlDB.DB
	err := DeleteTable(db, constants.SQL_CHANNEL_TAG)
	if err != nil {
		return nil
	}
	return DeleteTable(db, constants.SQL_CONTENT_TAG)
}

func (sqlDB *SqlDBWrapper) AddChannelArr(channels []common.Channel, height uint32) error {
	err := sqlDB.addChannels(channels)
	if err != nil {
		return err
	}
	err = sqlDB.addChannelsTags(channels)
	if err != nil {
		return err
	}
	err = sqlDB.addChannelsContents(channels)
	if err != nil {
		return err
	}
	err = sqlDB.addChannelsContentsTags(channels)
	if err != nil {
		return err
	}
	err = sqlDB.addChannelsPlaylistsTemps(channels, height)
	if err != nil {
		return err
	}
	return nil
}

func (sqlDB *SqlDBWrapper) addChannels(channels []common.Channel) error {
	db := sqlDB.DB
	insertCols := []string{
		constants.SQL_TABLE_CHANNEL__HASH,
		constants.SQL_TABLE_CHANNEL__TITLE,
	}
	insertData := make([]string, 2, 2)

	stmt, err := PrepareStmtInsertUpdate(db, constants.SQL_CHANNEL, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channels: %s", err.Error())
	}
	defer (*stmt).Close()
	for _, c := range channels {
		insertData[0] = c.RootChainID.String()
		insertData[1] = c.ChannelTitle.String()
		err = ExecStmt(stmt, insertData)
		if err != nil {
			return fmt.Errorf("Error insert/update channel: %s", err.Error())
		}
	}
	return nil
}

func (sqlDB *SqlDBWrapper) addChannelsTags(channels []common.Channel) error {
	db := sqlDB.DB
	insertCols := []string{
		constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID,
		constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID,
	}

	insertData := make([]string, 2, 2)

	stmt, err := PrepareStmtInsertUpdate(db, constants.SQL_CHANNEL_TAG_REL, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channel tags: %s", err.Error())
	}
	defer (*stmt).Close()
	for _, c := range channels {
		tags := c.Tags.GetTags()
		fmt.Printf("Tags found amount [%d]\n", len(tags))
		for _, t := range tags {
			tag, err := SelectSingleFromTable(db,
				constants.SQL_TABLE_CHANNEL_TAG__ID,
				constants.SQL_CHANNEL_TAG,
				constants.SQL_TABLE_CHANNEL_TAG__NAME,
				t.String())
			if err != nil {
				fmt.Printf("Error retrieving channel tag id for tag name [%s]: %s\n", t.String(), err.Error())
				continue
			}
			insertData[0] = c.RootChainID.String()
			insertData[1] = tag
			fmt.Printf("Inserting tag [%s] for hash [%s]\n", insertData[1], insertData[0])
			err = ExecStmt(stmt, insertData)
			if err != nil {
				return fmt.Errorf("Error insert/update channel tag [%s] with length [%s]: %s", tag, t, err.Error())
			}
		}
	}
	return nil
}

func (sqlDB *SqlDBWrapper) addChannelsContents(channels []common.Channel) error {
	db := sqlDB.DB
	insertCols := []string{
		constants.SQL_TABLE_CONTENT__CONTENT_HASH,
		constants.SQL_TABLE_CONTENT__TITLE,
		constants.SQL_TABLE_CONTENT__SERIES_NAME,
		constants.SQL_TABLE_CONTENT__PART_NAME,
		constants.SQL_TABLE_CONTENT__CH_ID,
	}

	insertData := make([]string, 5, 5)

	stmt, err := PrepareStmtInsertUpdate(db, constants.SQL_CONTENT, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channel contents: %s", err.Error())
	}
	defer (*stmt).Close()
	for _, channel := range channels {
		//Add Content for Channel
		contents := channel.Content.GetContents()
		for _, c := range contents {
			s, _ := primitives.BytesToUint32(append([]byte{0x00, 0x00, 0x00}, c.Series))
			p, _ := primitives.BytesToUint32(append([]byte{0x00, 0x00}, c.Part[:]...))

			insertData[0] = c.ContentID.String()
			insertData[1] = c.ShortDescription.String()
			insertData[2] = fmt.Sprintf("%d", s)
			insertData[3] = fmt.Sprintf("%d", p)
			insertData[4] = channel.RootChainID.String()

			err = ExecStmt(stmt, insertData)
			if err != nil {
				return fmt.Errorf("Error inserting content with hash[%s]: %s", c.ContentID.String(), err.Error())
			}
		}
	}
	return nil
}

func (sqlDB *SqlDBWrapper) addChannelsContentsTags(channels []common.Channel) error {
	db := sqlDB.DB
	insertCols := []string{
		constants.SQL_TABLE_CONTENT_TAG_REL__C_ID,
		constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID,
	}

	insertData := make([]string, 2, 2)

	stmt, err := PrepareStmtInsertUpdate(db, constants.SQL_CONTENT_TAG_REL, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channel content tags: %s", err.Error())
	}
	defer (*stmt).Close()
	for _, channel := range channels {
		//Add Content for Channel
		contents := channel.Content.GetContents()
		for _, c := range contents {
			tags := c.Tags.GetTags()
			for _, t := range tags {
				tag, err := SelectSingleFromTable(db,
					constants.SQL_TABLE_CONTENT_TAG__ID,
					constants.SQL_CONTENT_TAG,
					constants.SQL_TABLE_CONTENT_TAG__NAME,
					t.String())
				if err != nil {
					fmt.Printf("Error retrieving content tag id for tag name [%s]: %s\n", t.String(), err.Error())
					continue
				}
				insertData[0] = c.ContentID.String()
				insertData[1] = tag

				err = ExecStmt(stmt, insertData)
				if err != nil {
					if err == sqlite3.ErrConstraintUnique {
						fmt.Printf("WARNING attempted to insert duplicate tag for content: %s\n", err.Error())
					} else {
						return fmt.Errorf("Error inserting channel tag [%s] with length [%s]: %s", tag, t, err.Error())
					}
				}
			}
		}
	}
	return nil
}

func (sqlDB *SqlDBWrapper) addChannelsPlaylistsTemps(channels []common.Channel, height uint32) error {
	db := sqlDB.DB
	insertCols := []string{
		constants.SQL_TABLE_PLAYLIST_TEMP__TITLE,
		constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT,
		constants.SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID,
		constants.SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID,
	}

	insertData := make([]string, 4, 4)

	stmt, err := PrepareStmtInsertUpdate(db, constants.SQL_PLAYLIST_TEMP, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channel playlists: %s", err.Error())
	}
	defer (*stmt).Close()
	for _, channel := range channels {
		//Add Content for Channel
		playlists := channel.Playlist.GetPlaylists()
		for _, p := range playlists {
			//go through the content hash list for each playlist
			for _, ph := range p.Playlist.GetHashes() {

				insertData[0] = p.Title.String()
				insertData[1] = fmt.Sprintf("%d", height)
				insertData[2] = channel.RootChainID.String()
				insertData[3] = ph.String()

				err = ExecStmt(stmt, insertData)
				if err != nil {
					return fmt.Errorf("Error inserting playlist with hash[%s]: %s", ph.String(), err.Error())
				}
			}
		}
	}
	return nil
}

func (sqlDB *SqlDBWrapper) FlushPlaylistTempTable(currentHeight uint32) error {
	db := sqlDB.DB
	rowQuery := "SELECT COUNT(" + constants.SQL_TABLE_PLAYLIST_TEMP__ID + ") " +
		" FROM " + constants.SQL_PLAYLIST_TEMP +
		" WHERE " + constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT + " = ?"
	var rowCount int
	err := db.QueryRow(rowQuery, currentHeight).Scan(&rowCount)
	if err != nil {
		return fmt.Errorf("Error retrieving row count for flush playlist [%s]: %s", rowQuery, err.Error())
	}

	s := "SELECT " + constants.SQL_TABLE_PLAYLIST_TEMP__TITLE + ", " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID + ", " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID + ", " +
		constants.SQL_TABLE_PLAYLIST_TEMP__ID +
		" FROM " + constants.SQL_PLAYLIST_TEMP +
		" WHERE " + constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT + " = ?"
	rows, err := db.Query(s, currentHeight)
	if err != nil {
		return fmt.Errorf("Error select all from playlistTemp with query [%s]: %s", s, err.Error())
	}

	nRows := 0
	var title string
	var channelId string
	var contentId string
	var id string

	tableEntries := make([]PlayListTempStoreEntry, rowCount)

	//for each row attempt to insert into the playlist table
	for rows.Next() {
		err := rows.Scan(&title, &channelId, &contentId, &id)
		if err != nil {
			return fmt.Errorf("Error reading from playlistTemp: %s", err.Error())
		}
		tE := PlayListTempStoreEntry{title, channelId, contentId, id}
		tableEntries[nRows] = tE
		nRows++
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("ERROR when retrieving rows from playlistTemp\n")
	}
	// } else {
	// 	fmt.Printf("TempPlaylist rows went through [%d]\n", nRows)
	// }
	rows.Close()

	insertColsPlayList := []string{
		constants.SQL_TABLE_PLAYLIST__PLAYLIST_TITLE,
		constants.SQL_TABLE_PLAYLIST__CHANNEL_ID,
	}

	insertColsPlayListRel := []string{
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID,
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID,
	}

	insertData := make([]string, 2, 2)

	stmtPlayList, err := PrepareStmtInsertUpdate(db, constants.SQL_PLAYLIST, insertColsPlayList, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channel playlists: %s", err.Error())
	}
	stmtPlayListRel, err := PrepareStmtInsertUpdate(db, constants.SQL_PLAYLIST_CONTENT_REL, insertColsPlayListRel, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add channel playlists: %s", err.Error())
	}
	defer (*stmtPlayList).Close()
	defer (*stmtPlayListRel).Close()

	for i := 0; i < nRows; i++ {
		var retValChannel int
		var retValContent int
		s := "SELECT 1 FROM " + constants.SQL_CHANNEL + " WHERE " + constants.SQL_TABLE_CHANNEL__HASH + " = ?"
		err := db.QueryRow(s, tableEntries[i].ChannelId).Scan(&retValChannel)
		if err != nil {
			fmt.Errorf("Error retrieving rows channel: %s", err.Error())
			continue
		}
		s = "SELECT 1 FROM " + constants.SQL_CONTENT + " WHERE " + constants.SQL_TABLE_CONTENT__CH_ID + " = ?"
		err = db.QueryRow(s, tableEntries[i].ContentId).Scan(&retValContent)
		if err != nil {
			fmt.Errorf("Error retrieving rows content: %s", err.Error())
			continue
		}
		if retValChannel == 1 && retValContent == 1 {
			//Insert into playlist table
			insertData[0] = tableEntries[i].Title
			insertData[1] = tableEntries[i].ChannelId
			res, err := ExecStmtResult(stmtPlayList, insertData)
			if err != nil {
				return fmt.Errorf("Error adding playlist from playlistTemp: %s\n", err.Error())
			}
			//retrieve id
			id, err := res.LastInsertId()
			if err != nil {
				fmt.Printf("WARNING retrieving returned id from temp playlist with title[%s] and channelId[%s] and contentID[%s]\n", title, channelId, contentId)
			}

			//Insert into playlistRel table
			insertData[0] = fmt.Sprintf("%d", id)
			insertData[1] = tableEntries[i].ContentId
			err = ExecStmt(stmtPlayListRel, insertData)
			if err != nil {
				return fmt.Errorf("Error adding playlistRel from playlistTemp: %s\n", err.Error())
			}
		}
	}

	deleteQuery := "DELETE FROM " + constants.SQL_PLAYLIST_TEMP + " WHERE " + constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT + " <= ?"
	_, err = db.Exec(deleteQuery, currentHeight)
	if err != nil {
		fmt.Printf("ERROR!! CRUCIAL problems delete query deleting playlsit temp index's with query [%s]: Error [%s]\n", deleteQuery, err.Error())
	}

	return nil
}

func (sqlDB *SqlDBWrapper) AddTags() error {
	db := sqlDB.DB
	insertColsChannel := []string{constants.SQL_TABLE_CHANNEL_TAG__NAME}
	insertColsContent := []string{constants.SQL_TABLE_CONTENT_TAG__NAME}

	insertData := make([]string, 1, 1)

	stmtChannel, err := PrepareStmtInsertUpdate(db, constants.SQL_CHANNEL_TAG, insertColsChannel, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add table tags channel: %s", err.Error())
	}
	stmtContent, err := PrepareStmtInsertUpdate(db, constants.SQL_CONTENT_TAG, insertColsContent, insertData)
	if err != nil {
		return fmt.Errorf("Error preparing to add table tags content: %s", err.Error())
	}
	for _, e := range constants.ALLOWED_TAGS {
		_, err := stmtChannel.Exec(e)
		if err != nil {
			return fmt.Errorf("Error inserting table tags channel: %s", err.Error())
		}
		_, err = stmtContent.Exec(e)
		if err != nil {
			return fmt.Errorf("Error inserting table tags content: %s", err.Error())
		}
	}
	return nil
}

func (sqlDB *SqlDBWrapper) DeleteDBChannels() error {
	db := sqlDB.DB
	s := "DELETE FROM " + constants.SQL_CHANNEL
	_, err := db.Exec(s)
	if err != nil {
		return fmt.Errorf("Error deleting items from db with query[%s]: %s", s, err.Error())
	}
	return nil
}
