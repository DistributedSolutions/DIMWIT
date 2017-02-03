package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"os"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/util"
)

type PlayListTempStoreEntry struct {
	Title     string
	ChannelId string
	ContentId string
	Id        string
}

var _ = sqlite3.ErrNoMask

var gDB *sql.DB

var TABLE_NAMES = []string{
	constants.SQL_CHANNEL,
	constants.SQL_CHANNEL_TAG,
	constants.SQL_CHANNEL_TAG_REL,
	constants.SQL_CONTENT,
	constants.SQL_CONTENT_TAG,
	constants.SQL_CONTENT_TAG_REL,
	constants.SQL_PLAYLIST,
	constants.SQL_PLAYLIST_CONTENT_REL,
}

var CREATE_TABLE = []string{
	constants.SQL_CHANNEL + "(" +
		constants.SQL_TABLE_CHANNEL__HASH + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") PRIMARY KEY, " +
		constants.SQL_TABLE_CHANNEL__TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TAG_MAX_LENGTH) + ") NOT NULL)",

	constants.SQL_CHANNEL_TAG + "(" +
		constants.SQL_TABLE_CHANNEL_TAG__ID + " INTEGER PRIMARY KEY, " +
		constants.SQL_TABLE_CHANNEL_TAG__NAME + " VARCHAR(" + fmt.Sprintf("%d", constants.TAG_MAX_LENGTH) + ") NOT NULL UNIQUE)",

	constants.SQL_CHANNEL_TAG_REL + "(" +
		constants.SQL_TABLE_CHANNEL_TAG_REL__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + ") REFERENCES " + constants.SQL_CHANNEL_TAG +
		"(" + constants.SQL_TABLE_CHANNEL_TAG__ID + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_CONTENT + "(" +
		constants.SQL_TABLE_CONTENT__CONTENT_HASH + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") PRIMARY KEY, " +
		constants.SQL_TABLE_CONTENT__TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TITLE_MAX_LENGTH) + ") NOT NULL, " +
		constants.SQL_TABLE_CONTENT__SERIES_NAME + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT__PART_NAME + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT__CH_ID + " INTEGER NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT__CH_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_CONTENT_TAG + "(" +
		constants.SQL_TABLE_CONTENT_TAG__ID + " INTEGER PRIMARY KEY UNIQUE, " +
		constants.SQL_TABLE_CONTENT_TAG__NAME + " name VARCHAR(" + fmt.Sprintf("%d", constants.TAG_MAX_LENGTH) + ") NOT NULL UNIQUE)",

	constants.SQL_CONTENT_TAG_REL + "(" +
		constants.SQL_TABLE_CONTENT_TAG_REL__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + ") REFERENCES " + constants.SQL_CONTENT +
		"(" + constants.SQL_TABLE_CONTENT__CONTENT_HASH + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + ") REFERENCES " + constants.SQL_CONTENT_TAG +
		"(" + constants.SQL_TABLE_CONTENT_TAG__ID + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_PLAYLIST + "(" +
		constants.SQL_TABLE_PLAYLIST__ID + " INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_PLAYLIST__PLAYLIST_TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TITLE_MAX_LENGTH) + ") NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST__CHANNEL_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_PLAYLIST__CHANNEL_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_PLAYLIST_CONTENT_REL + "(" +
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID + ") REFERENCES " + constants.SQL_PLAYLIST +
		"(" + constants.SQL_TABLE_PLAYLIST__ID + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID + ") REFERENCES " + constants.SQL_CONTENT +
		"(" + constants.SQL_TABLE_CONTENT__CONTENT_HASH + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_PLAYLIST_TEMP + "(" +
		constants.SQL_TABLE_PLAYLIST_TEMP__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TITLE_MAX_LENGTH) + ") NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL)",
}

func CreateDB(dbName string, tableCreate []string) error {
	dir := util.GetHomeDir() + constants.HIDDEN_DIR
	_, err := os.Stat(dir)
	// create directory if not exists
	if os.IsNotExist(err) {
		os.MkdirAll(dir, constants.DIRECTORY_PERMISSIONS)
	}

	dbPathName := dir + dbName
	_, err = os.Stat(dbPathName)

	//create db if not exists
	if os.IsNotExist(err) {
		file, err := os.OpenFile(dbPathName, os.O_CREATE|os.O_RDWR, constants.FILE_PERMISSIONS)
		if err != nil {
			return fmt.Errorf("Error creating database: %s", err.Error())
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbPathName)
	if err != nil {
		return fmt.Errorf("Error opening database: %s", err.Error())
	}

	gDB = db

	_, err = db.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		return fmt.Errorf("Error setting pragma: %s", err.Error())
	}

	//create all tables if they do not exist
	err = createAllTables(db, tableCreate)
	if err != nil {
		return fmt.Errorf("Error creating tables: %s", err.Error())
	}

	return err
}

func createAllTables(db *sql.DB, tableCreate []string) error {
	for _, element := range tableCreate {
		//check if table exists
		s := "CREATE TABLE IF NOT EXISTS " + element + ";"
		_, err := db.Exec(s)
		if err != nil {
			return fmt.Errorf("Error creating table with query [%s]: %s", s, err.Error())
		}
	}
	return nil
}

func DeleteDB(dbName string) error {
	dir := util.GetHomeDir() + constants.HIDDEN_DIR
	_, err := os.Stat(dir)
	// create directory if not exists
	if !os.IsNotExist(err) {
		dbPathName := dir + dbName
		_, err = os.Stat(dbPathName)

		if !os.IsNotExist(err) {
			err = os.Remove(dbPathName)
			if err != nil {
				return fmt.Errorf("Error deleting db with path [%s]: %s", dbPathName, err.Error())
			}
		}
	}
	return nil
}

func getDBPath() (string, error) {
	path := util.GetHomeDir() + constants.HIDDEN_DIR + constants.SQL_DB
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("DB Does NOT Exists with dir [%s]: %s", path, err.Error())
	}

	return path, nil
}

func GetDB() (*sql.DB, error) {
	if gDB == nil {
		dbPath, err := getDBPath()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving path: %s", err.Error())
		}
		gDB, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, fmt.Errorf("Error opening DB: %s", err.Error())
		}
	}
	return gDB, nil
}

func DeleteTable(db *sql.DB, tableName string) error {
	s := "DELETE FROM " + tableName
	_, err := gDB.Exec(s)
	if err != nil {
		return fmt.Errorf("Error deleting from table %s: %s", tableName, err.Error())
	}
	return nil
}

func InsertIntoTable(db *sql.DB, tableName string, insertCols []string, insertData []string) (sql.Result, error) {
	icl := len(insertCols)
	idl := len(insertData)
	if idl != icl {
		return nil, fmt.Errorf("Error in argument lengths while inserting into %s (%d != %d)", tableName, icl, idl)
	}

	s := "INSERT INTO " + tableName + " ("
	for i, e := range insertCols {
		s += e
		if i < icl-1 && icl > 1 {
			s += ","
		}
	}
	s += ") values("

	for i := 0; i < idl; i++ {
		s += "?"
		if i < idl-1 && idl > 1 {
			s += ","
		}
	}
	s += ")"

	stmt, err := db.Prepare(s)
	if err != nil {
		return nil, fmt.Errorf("Error preparing inserting into table [%s] with query [%s]: %s", tableName, s, err.Error())
	}
	defer stmt.Close()

	insertDataInterface := make([]interface{}, len(insertData))
	for index, value := range insertData {
		insertDataInterface[index] = value
	}

	res, err := stmt.Exec(insertDataInterface...)
	if err != nil {
		return nil, fmt.Errorf("Error exec inserting into [%s] with query[%s]: %s", tableName, s, err.Error())
	}
	return res, nil
}

func AddTags() error {
	gDB, err := GetDB()
	if err != nil {
		return fmt.Errorf("Error adding tags: %s", err.Error())
	}

	insertColsChannel := []string{constants.SQL_TABLE_CHANNEL_TAG__NAME}
	insertColsContent := []string{constants.SQL_TABLE_CONTENT_TAG__NAME}
	for _, e := range constants.ALLOWED_TAGS {
		insertData := []string{e}
		_, err = InsertIntoTable(gDB, constants.SQL_CHANNEL_TAG, insertColsChannel, insertData)
		if err != nil {
			return fmt.Errorf("Error inserting tags: %s", err.Error())
		}
		_, err = InsertIntoTable(gDB, constants.SQL_CONTENT_TAG, insertColsContent, insertData)
		if err != nil {
			return fmt.Errorf("Error inserting tags: %s", err.Error())
		}
	}
	return nil
}

func DeleteTags() error {
	gDB, err := GetDB()
	if err != nil {
		return fmt.Errorf("Error updating tags: %s", err.Error())
	}

	return DeleteTable(gDB, "channelTag")
}

//NOT GOING TO CHECK IF IN DB ALREADY
func AddChannel(channel *common.Channel, height int) error {
	gDB, err := GetDB()
	if err != nil {
		return fmt.Errorf("Error adding channel: %s", err.Error())
	}

	insertCols := []string{
		constants.SQL_TABLE_CHANNEL__HASH,
		constants.SQL_TABLE_CHANNEL__TITLE,
	}
	insertData := []string{
		channel.RootChainID.String(),
		channel.ChannelTitle.String(),
	}
	_, err = InsertIntoTable(gDB, constants.SQL_CHANNEL, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error adding channel: %s", err.Error())
	}

	//Add channel tags for channel
	tags := channel.Tags.GetTags()
	for _, t := range tags {
		tag, err := SelectSingleFromTable(gDB,
			constants.SQL_TABLE_CHANNEL_TAG__ID,
			constants.SQL_CHANNEL_TAG,
			constants.SQL_TABLE_CHANNEL_TAG__NAME,
			t.String())
		if err != nil {
			return fmt.Errorf("Error retrieving tag id: %s", err.Error())
		}
		insertCols = []string{
			constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID,
			constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID,
		}
		insertData = []string{
			channel.RootChainID.String(),
			tag,
		}
		_, err = InsertIntoTable(gDB, constants.SQL_CHANNEL_TAG_REL, insertCols, insertData)
		if err != nil {
			return fmt.Errorf("Error inserting channel tag [%s] with length [%s]: %s", tag, t, err.Error())
		}
	}

	//Add Content for Channel
	contents := channel.Content.GetContents()
	for _, c := range contents {
		if err != nil {
			return fmt.Errorf("Error retrieving tag id: %s", err.Error())
		}
		insertCols = []string{
			constants.SQL_TABLE_CONTENT__CONTENT_HASH,
			constants.SQL_TABLE_CONTENT__TITLE,
			constants.SQL_TABLE_CONTENT__SERIES_NAME,
			constants.SQL_TABLE_CONTENT__PART_NAME,
			constants.SQL_TABLE_CONTENT__CH_ID,
		}
		s, err := primitives.BytesToUint32(append([]byte{0x00, 0x00, 0x00}, c.Series))
		if err != nil {
			return fmt.Errorf("FUCK YOUR BYTES s: %s", err.Error())
		}
		p, err := primitives.BytesToUint32(append([]byte{0x00, 0x00}, c.Part[:]...))
		if err != nil {
			return fmt.Errorf("FUCK YOUR BYTES p: %s", err.Error())
		}

		insertData = []string{
			c.ContentID.String(),
			c.ShortDescription.String(),
			fmt.Sprintf("%d", s),
			fmt.Sprintf("%d", p),
			channel.RootChainID.String(),
		}
		_, err = InsertIntoTable(gDB, constants.SQL_CONTENT, insertCols, insertData)
		if err != nil {
			return fmt.Errorf("Error inserting content with hash[%s]: %s", c.ContentID.String(), err.Error())
		}

		//Add Content Tags for content
		tags := c.Tags.GetTags()
		for _, t := range tags {
			tag, err := SelectSingleFromTable(gDB,
				constants.SQL_TABLE_CONTENT_TAG__ID,
				constants.SQL_CONTENT_TAG,
				constants.SQL_TABLE_CONTENT_TAG__NAME,
				t.String())
			if err != nil {
				return fmt.Errorf("Error retrieving tag id: %s", err.Error())
			}
			insertCols = []string{
				constants.SQL_TABLE_CONTENT_TAG_REL__C_ID,
				constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID,
			}
			insertData = []string{
				c.ContentID.String(),
				tag,
			}
			_, err = InsertIntoTable(gDB, constants.SQL_CONTENT_TAG_REL, insertCols, insertData)
			if err != nil {
				return fmt.Errorf("Error inserting channel tag [%s] with length [%s]: %s", tag, t, err.Error())
			}
		}
	}

	//Add PlaylistsTemp for Channel
	playlists := channel.Playlist.GetPlaylists()
	for _, p := range playlists {
		//go through the content hash list for each playlist
		for _, ph := range p.Playlist.GetHashes() {

			insertCols = []string{
				constants.SQL_TABLE_PLAYLIST_TEMP__TITLE,
				constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT,
				constants.SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID,
				constants.SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID,
			}
			insertData = []string{
				p.Title.String(),
				fmt.Sprintf("%d", height),
				channel.RootChainID.String(),
				ph.String(),
			}
			_, err := InsertIntoTable(gDB, constants.SQL_PLAYLIST_TEMP, insertCols, insertData)
			if err != nil {
				return fmt.Errorf("Error inserting playlistTemp title [%s]: %s", p.Title.String(), err.Error())
			}
		}
	}

	return nil
}

func FlushPlaylistTempTable(db *sql.DB, currentHeight int) error {
	rowQuery := "SELECT COUNT(" + constants.SQL_TABLE_PLAYLIST_TEMP__ID + ") " +
		" FROM " + constants.SQL_PLAYLIST_TEMP
	var rowCount int
	err := db.QueryRow(rowQuery).Scan(&rowCount)
	if err != nil {
		return fmt.Errorf("Error retrieving row count for flush playlist [%s]: %s", rowQuery, err.Error())
	}

	s := "SELECT " + constants.SQL_TABLE_PLAYLIST_TEMP__TITLE + ", " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID + ", " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID + ", " +
		constants.SQL_TABLE_PLAYLIST_TEMP__ID +
		" FROM " + constants.SQL_PLAYLIST_TEMP
	rows, err := db.Query(s)
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
		// tE.Title = title
		// tE.ChannelId = channelId
		// tE.ContentId = contentId
		// tE.Id = id
		tableEntries[nRows] = tE
		nRows++
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("ERROR when retrieving rows from playlistTemp\n")
	} else {
		fmt.Printf("TempPlaylist rows went through [%d]\n", nRows)
	}
	rows.Close()
	deleteQuery := "DELETE FROM " + constants.SQL_PLAYLIST_TEMP + " WHERE " + constants.SQL_TABLE_PLAYLIST_TEMP__ID + " IN("

	for i := 0; i < nRows; i++ {
		//Insert into playlist table
		insertCols := []string{
			constants.SQL_TABLE_PLAYLIST__PLAYLIST_TITLE,
			constants.SQL_TABLE_PLAYLIST__CHANNEL_ID,
		}
		insertData := []string{
			tableEntries[i].Title,
			tableEntries[i].ChannelId,
		}
		res, err := InsertIntoTable(gDB, constants.SQL_PLAYLIST, insertCols, insertData)
		if err != nil {
			return fmt.Errorf("Error adding channel: %s\n", err.Error())
		}

		//retrieve id
		id, err := res.LastInsertId()
		if err != nil {
			fmt.Printf("WARNING retrieving returned id from temp playlist with title[%s] and channelId[%s] and contentID[%s]\n", title, channelId, contentId)
		}

		//Insert into playlistRel table
		insertCols = []string{
			constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID,
			constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID,
		}
		insertData = []string{
			fmt.Sprintf("%d", id),
			tableEntries[i].ContentId,
		}
		_, err = InsertIntoTable(gDB, constants.SQL_PLAYLIST_CONTENT_REL, insertCols, insertData)
		if err != nil {
			fmt.Printf("WARNING 'MOST LIKELY FOREIGN KEY CONSTRAINT FAIL' inserting into playlist rel table with title[%s] and channelId[%s] and contentID[%s] error message is [%s]\n", title, channelId, contentId, err.Error())
		}

		deleteQuery += fmt.Sprintf("%d", id) + ","
	}

	deleteQuery = deleteQuery[0:(len(deleteQuery)-1)] + ")"
	_, err = db.Exec(deleteQuery)
	if err != nil {
		fmt.Printf("ERROR!! CRUCIAL problems deleting index's with query [%s]: Error [%s]\n", deleteQuery, err.Error())
	}

	return nil
}

//NOT GOING TO CHECK IF IN DB ALREADY
func UpdateChannel(db *sql.DB, channel *common.Channel) error {
	db, err := GetDB()
	if err != nil {
		return fmt.Errorf("Error adding channel: %s", err.Error())
	}

	_, err = db.Exec("DELETE FROM channel WHERE channelHash = ?",
		channel.RootChainID.String())
	if err != nil {
		return fmt.Errorf("Error deleting channel: %s", err.Error())
	}
	err = AddChannel(channel, -1)
	if err != nil {
		return fmt.Errorf("Error inserting channel after deleting: %s", err.Error())
	}
	return nil
}

func SelectSingleFromTable(db *sql.DB, colReturn string, tableName string, columnOn string, singleLookup string) (string, error) {
	s := "SELECT " + colReturn + " FROM " + tableName + " WHERE " + columnOn + " = (?) LIMIT 1"
	rows, err := gDB.Query(s, singleLookup)
	if err != nil {
		return "", fmt.Errorf("Error SINGLE select single from table with query [%s]: %s", s, err.Error())
	}

	rows.Next()

	var result string
	rows.Scan(&result)

	// fmt.Printf("----SINGLE select returned empty with query [%s] with result [%s] [%s]\n", s, result, singleLookup)
	rows.Close()
	return result, nil
}

func CloseDB() {
	if gDB != nil {
		gDB.Close()
	}
}
