package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"os"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/util"
)

// type SqlObj interface {
// 	Santize() string
// 	Type() int
// }

// type SqlString string

// func (s *SqlString) SetString(st string) {
// 	*s = SqlString(st)
// }

// func (s *SqlString) Sanitize() string {
// 	return "'" + SQLSanitize(s) + "'"
// }

// func (s *SqlString) Type() int {
// 	return constants.SQL_STRING
// }

// type SqlInt int

// func (s *SqlString) SetInt(st string) {
// 	*s = SqlInt(st)
// }

// func (s *SqlString) Sanitize() string {
// 	return s
// }

// func (s *SqlString) Type() int {
// 	return constants.SQL_OTHER
// }

var _ = sqlite3.ErrNoMask

var gDB *sql.DB

var TABLE_NAMES = []string{
	constants.SQL_CHANNEL,
	constants.SQL_CHANNEL_TAG,
	constants.SQL_CHANNEL_TAG_REL,
	constants.SQL_PLAYLIST,
	constants.SQL_CONTENT,
	constants.SQL_CONTENT_TAG,
	constants.SQL_CONTENT_TAG_REL,
}

var CREATE_TABLE = []string{
	constants.SQL_CHANNEL + "(" +
		constants.SQL_TABLE_CHANNEL__HASH + " CHAR(20) PRIMARY KEY, " +
		constants.SQL_TABLE_CHANNEL__TITLE + " VARCHAR(100) NOT NULL)",
	constants.SQL_CHANNEL_TAG + "(" +
		constants.SQL_TABLE_CHANNEL_TAG__ID + " INTEGER PRIMARY KEY, " +
		constants.SQL_TABLE_CHANNEL_TAG__NAME + " VARCHAR(100) NOT NULL UNIQUE)",
	constants.SQL_CHANNEL_TAG_REL + "(" +
		constants.SQL_TABLE_CHANNEL_TAG_REL__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + " INTEGER NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + ") REFERENCES " + constants.SQL_CHANNEL_TAG +
		"(" + constants.SQL_TABLE_CHANNEL_TAG__ID + ") ON DELETE CASCADE ON UPDATE CASCADE)",
	constants.SQL_PLAYLIST + "(" +
		constants.SQL_TABLE_PLAYLIST__ID + " INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_PLAYLIST__PLAYLIST_TITLE + " VARCHAR(100) NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST__CHANNEL_ID + " INTEGER REFERENCES channel(id))",
	constants.SQL_CONTENT + "(" +
		constants.SQL_TABLE_CONTENT__CONTENT_HASH + " CHAR(20) PRIMARY KEY, " +
		constants.SQL_TABLE_CONTENT__TITLE + " VARCHAR(100) NOT NULL, " +
		constants.SQL_TABLE_CONTENT__SERIES_NAME + " VARCHAR(100) NOT NULL, " +
		constants.SQL_TABLE_CONTENT__PART_NAME + " VARCHAR(100) NOT NULL, " +
		constants.SQL_TABLE_CONTENT__CH_ID + " INTEGER NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT__CH_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE)",
	constants.SQL_CONTENT_TAG + "(" +
		constants.SQL_TABLE_CONTENT_TAG__ID + " INTEGER PRIMARY KEY UNIQUE, " +
		constants.SQL_TABLE_CONTENT_TAG__NAME + " name VARCHAR(100) NOT NULL UNIQUE)",
	constants.SQL_CONTENT_TAG_REL + "(" +
		constants.SQL_TABLE_CONTENT_TAG_REL__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + " INTEGER NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + ") REFERENCES " + constants.SQL_CONTENT +
		"(" + constants.SQL_TABLE_CONTENT__CONTENT_HASH + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + ") REFERENCES " + constants.SQL_CONTENT_TAG +
		"(" + constants.SQL_TABLE_CONTENT_TAG__ID + ") ON DELETE CASCADE ON UPDATE CASCADE)",
}

func CreateDB(dbName string, tableCreate []string) error {
	dir := util.GetHomeDir() + constants.HIDDEN_DIR
	_, err := os.Stat(dir)
	// create directory if not exists
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0666)
	}

	dbPathName := dir + dbName
	_, err = os.Stat(dbPathName)

	//create db if not exists
	if os.IsNotExist(err) {
		file, err := os.OpenFile(dbPathName, os.O_CREATE|os.O_APPEND, 0666)
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
	if os.IsNotExist(err) {
		os.MkdirAll(dir, 0666)
	}

	dbPathName := dir + dbName
	_, err = os.Stat(dbPathName)

	if !os.IsNotExist(err) {
		err = os.Remove(dbPathName)
		if err != nil {
			return fmt.Errorf("Error deleting db with path [%s]: %s", dbPathName, err.Error())
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

func getDB() (*sql.DB, error) {
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

func InsertIntoTable(db *sql.DB, tableName string, insertCols []string, insertData []string) error {
	icl := len(insertCols)
	idl := len(insertData)
	if idl != icl {
		return fmt.Errorf("Error in argument lengths while inserting into %s (%d != %d)", tableName, icl, idl)
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
		return fmt.Errorf("Error preparing inserting into table [%s] with query [%s]: %s", tableName, s, err.Error())
	}
	defer stmt.Close()

	insertDataInterface := make([]interface{}, len(insertData))
	for index, value := range insertData {
		insertDataInterface[index] = value
	}

	_, err = stmt.Exec(insertDataInterface...)
	if err != nil {
		return fmt.Errorf("Error exec inserting into %s: %s", tableName, err.Error())
	}
	return nil
}

func AddTags() error {
	gDB, err := getDB()
	if err != nil {
		return fmt.Errorf("Error adding tags: %s", err.Error())
	}

	insertCols := []string{"name"}
	for _, e := range constants.ALLOWED_TAGS {
		insertData := []string{e}
		err = InsertIntoTable(gDB, "channelTag", insertCols, insertData)
		if err != nil {
			return fmt.Errorf("Error inserting tags: %s", err.Error())
		}
	}
	return nil
}

func DeleteTags() error {
	gDB, err := getDB()
	if err != nil {
		return fmt.Errorf("Error updating tags: %s", err.Error())
	}

	return DeleteTable(gDB, "channelTag")
}

//NOT GOING TO CHECK IF IN DB ALREADY
func AddChannel(channel *common.Channel) error {
	gDB, err := getDB()
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
	err = InsertIntoTable(gDB, constants.SQL_CHANNEL, insertCols, insertData)
	if err != nil {
		return fmt.Errorf("Error adding channel: %s", err.Error())
	}

	tags := channel.Tags.GetTags()
	for _, t := range tags {
		tag, err := getTagID(gDB, t.String())
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
		err = InsertIntoTable(gDB, constants.SQL_CHANNEL_TAG_REL, insertCols, insertData)
		if err != nil {
			return fmt.Errorf("Error inserting channel tag %s: %s", tag, err.Error())
		}
	}

	return nil
}

//NOT GOING TO CHECK IF IN DB ALREADY
func UpdateChannel(channel *common.Channel) error {
	gDB, err := getDB()
	if err != nil {
		return fmt.Errorf("Error adding channel: %s", err.Error())
	}

	_, err = gDB.Exec("DELETE FROM channel WHERE channelHash = ?",
		channel.RootChainID.String())
	if err != nil {
		return fmt.Errorf("Error deleting channel: %s", err.Error())
	}
	err = AddChannel(channel)
	if err != nil {
		return fmt.Errorf("Error inserting channel after deleting: %s", err.Error())
	}
	return nil
}

func getTagID(db *sql.DB, tagName string) (string, error) {
	rows, err := gDB.Query("SELECT id FROM channelTag WHERE name = ? LIMIT 1", tagName)
	if err != nil {
		return "", fmt.Errorf("Error retrieving tag id for [%s]: %s", tagName, err.Error())
	}

	defer rows.Close()

	rows.Next()

	var result string
	rows.Scan(&result)

	return result, nil
}

func CloseDB() {
	gDB.Close()
}
