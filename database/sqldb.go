package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"os"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/util"
)

type SqlDBWrapper struct {
	DB   *sql.DB
	Name string
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
	constants.SQL_PLAYLIST_TEMP,
}

var CREATE_TABLE = []string{
	constants.SQL_CHANNEL + "(" +
		constants.SQL_TABLE_CHANNEL__HASH + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") PRIMARY KEY, " +
		constants.SQL_TABLE_CHANNEL__TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TAG_MAX_LENGTH) + ") NOT NULL," +
		constants.SQL_TABLE_CHANNEL__DT + " datetime)",

	constants.SQL_CHANNEL_TAG + "(" +
		constants.SQL_TABLE_CHANNEL_TAG__ID + " INTEGER PRIMARY KEY, " +
		constants.SQL_TABLE_CHANNEL_TAG__NAME + " VARCHAR(" + fmt.Sprintf("%d", constants.TAG_MAX_LENGTH) + ") NOT NULL UNIQUE)",

	constants.SQL_CHANNEL_TAG_REL + "(" +
		constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + ") REFERENCES " + constants.SQL_CHANNEL_TAG +
		"(" + constants.SQL_TABLE_CHANNEL_TAG__ID + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"PRIMARY KEY (" + constants.SQL_TABLE_CHANNEL_TAG_REL__C_ID + "," + constants.SQL_TABLE_CHANNEL_TAG_REL__CT_ID + "))",

	constants.SQL_CONTENT + "(" +
		constants.SQL_TABLE_CONTENT__CONTENT_HASH + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") PRIMARY KEY, " +
		constants.SQL_TABLE_CONTENT__TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TITLE_MAX_LENGTH) + ") NOT NULL, " +
		constants.SQL_TABLE_CONTENT__SERIES_NAME + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT__PART_NAME + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT__CH_ID + " INTEGER NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT__CH_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_CONTENT_TAG + "(" +
		constants.SQL_TABLE_CONTENT_TAG__ID + " INTEGER PRIMARY KEY, " +
		constants.SQL_TABLE_CONTENT_TAG__NAME + " name VARCHAR(" + fmt.Sprintf("%d", constants.TAG_MAX_LENGTH) + ") NOT NULL UNIQUE," +
		constants.SQL_TABLE_CONTENT_TAG__DT + " datetime)",

	constants.SQL_CONTENT_TAG_REL + "(" +
		constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + ") REFERENCES " + constants.SQL_CONTENT +
		"(" + constants.SQL_TABLE_CONTENT__CONTENT_HASH + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + ") REFERENCES " + constants.SQL_CONTENT_TAG +
		"(" + constants.SQL_TABLE_CONTENT_TAG__ID + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"PRIMARY KEY (" + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + "," + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + "))",

	constants.SQL_PLAYLIST + "(" +
		constants.SQL_TABLE_PLAYLIST__ID + " INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_PLAYLIST__PLAYLIST_TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TITLE_MAX_LENGTH) + ") NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST__CHANNEL_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_PLAYLIST__CHANNEL_ID + ") REFERENCES " + constants.SQL_CHANNEL +
		"(" + constants.SQL_TABLE_CHANNEL__HASH + ") ON DELETE CASCADE ON UPDATE CASCADE)",

	constants.SQL_PLAYLIST_CONTENT_REL + "(" +
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID + ") REFERENCES " + constants.SQL_PLAYLIST +
		"(" + constants.SQL_TABLE_PLAYLIST__ID + ") ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (" + constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID + ") REFERENCES " + constants.SQL_CONTENT +
		"(" + constants.SQL_TABLE_CONTENT__CONTENT_HASH + ") ON DELETE CASCADE ON UPDATE CASCADE," +
		"PRIMARY KEY (" + constants.SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID + "," + constants.SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID + "))",

	constants.SQL_PLAYLIST_TEMP + "(" +
		constants.SQL_TABLE_PLAYLIST_TEMP__ID + " INTEGER PRIMARY KEY AUTOINCREMENT, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__TITLE + " VARCHAR(" + fmt.Sprintf("%d", constants.TITLE_MAX_LENGTH) + ") NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__HEIGHT + " INTEGER NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL, " +
		constants.SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID + " CHAR(" + fmt.Sprintf("%d", constants.HASH_BYTES_LENGTH*2) + ") NOT NULL)",
}

func CreateDB(dbName string, tableCreate []string) (*SqlDBWrapper, error) {
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
			return nil, fmt.Errorf("Error creating database: %s", err.Error())
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbPathName)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %s", err.Error())
	}

	_, err = db.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		return nil, fmt.Errorf("Error setting pragma: %s", err.Error())
	}

	//create all tables if they do not exist
	err = createAllTables(db, tableCreate)
	if err != nil {
		return nil, fmt.Errorf("Error creating tables: %s", err.Error())
	}

	dbStruct := SqlDBWrapper{db, dbName}
	return &dbStruct, nil
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

// func getDBPath() (string, error) {
// 	path := util.GetHomeDir() + constants.HIDDEN_DIR + constants.SQL_DB
// 	_, err := os.Stat(path)
// 	if os.IsNotExist(err) {
// 		return "", fmt.Errorf("DB Does NOT Exists with dir [%s]: %s", path, err.Error())
// 	}

// 	return path, nil
// }

// func GetDB() (*sql.DB, error) {
// 	if gDB == nil {
// 		dbPath, err := getDBPath()
// 		if err != nil {
// 			return nil, fmt.Errorf("Error retrieving path: %s", err.Error())
// 		}
// 		gDB, err = sql.Open("sqlite3", dbPath)
// 		if err != nil {
// 			return nil, fmt.Errorf("Error opening DB: %s", err.Error())
// 		}
// 	}
// 	return gDB, nil
// }

func DeleteTable(db *sql.DB, tableName string) error {
	s := "DELETE FROM " + tableName
	_, err := db.Exec(s)
	if err != nil {
		return fmt.Errorf("Error deleting from table %s: %s", tableName, err.Error())
	}
	return nil
}

func insertUpdateIntoTablePrepareString(tableName string, insertCols []string) string {
	s := "INSERT OR REPLACE INTO " + tableName + " ("
	s += CommaDelimiterArray(insertCols, false)
	s += ") values("
	s += QuestionString(len(insertCols))
	s += ");"
	return s
}

func PrepareStmtInsertUpdate(db *sql.DB, tableName string, insertCols []string, insertData []string) (*sql.Stmt, error) {
	icl := len(insertCols)
	idl := len(insertData)
	if idl != icl {
		return nil, fmt.Errorf("Error in argument lengths while inserting into %s (%d != %d)", tableName, icl, idl)
	}
	s := insertUpdateIntoTablePrepareString(tableName, insertCols)
	stmt, err := db.Prepare(s)
	if err != nil {
		return nil, fmt.Errorf("Error preparing inserting into table [%s] with query [%s]: %s", tableName, s, err.Error())
	}
	return stmt, nil
}

func ExecStmt(stmt *sql.Stmt, data []string) error {
	insertDataInterface := make([]interface{}, len(data))
	for index, value := range data {
		insertDataInterface[index] = value
	}

	_, err := stmt.Exec(insertDataInterface...)
	if err != nil {
		return fmt.Errorf("Error stmt exec: %s", err.Error())
	}
	return nil
}

func ExecStmtResult(stmt *sql.Stmt, data []string) (sql.Result, error) {
	insertDataInterface := make([]interface{}, len(data))
	for index, value := range data {
		insertDataInterface[index] = value
	}

	res, err := stmt.Exec(insertDataInterface...)
	if err != nil {
		return nil, fmt.Errorf("Error stmt exec: %s", err.Error())
	}
	return res, nil
}

// func InsertIntoTable(db *sql.DB, insertCols,) (sql.Result, error) {

// 	stmt, err := db.Prepare(s)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error preparing inserting into table [%s] with query [%s]: %s", tableName, s, err.Error())
// 	}
// 	defer stmt.Close()

// 	insertDataInterface := make([]interface{}, len(insertData))
// 	for index, value := range insertData {
// 		insertDataInterface[index] = value
// 	}

// 	res, err := stmt.Exec(insertDataInterface...)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error exec inserting into [%s] with query[%s]: %s", tableName, s, err.Error())
// 	}
// 	return res, nil
// }

func SelectSingleFromTable(db *sql.DB, colReturn string, tableName string, columnOn string, singleLookup string) (string, error) {
	s := "SELECT " + colReturn + " FROM " + tableName + " WHERE " + columnOn + " = (?) LIMIT 1"
	rows, err := db.Query(s, singleLookup)
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

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func QuestionString(length int) string {
	arr := make([]string, length, length)
	for i := range arr {
		arr[i] = "?"
	}
	return CommaDelimiterArray(arr, false)
}

func CommaDelimiterArray(arr []string, isString bool) string {
	s := ""
	l := len(arr)
	for i := 0; i < l; i++ {
		if isString {
			s += "'" + arr[i] + "'"
		} else {
			s += arr[i]
		}
		if i < l-1 && l > 1 {
			s += ","
		}
	}
	return s
}
