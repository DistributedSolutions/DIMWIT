package database

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"os"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/util"
)

var _ = sqlite3.ErrNoMask

var TABLE_NAMES = [7]string{
	"channel",
	"channelTag",
	"channelTagRel",
	"playlist",
	"content",
	"contentTag",
	"contentTagRel",
}

var CREATE_TABLE = [7]string{
	"channel(" +
		"channelHash CHAR(20) PRIMARY KEY, " +
		"tile VARCHAR(100) NOT NULL)",
	"channelTag(" +
		"id INTEGER PRIMARY KEY, " +
		"name VARCHAR(100) NOT NULL UNIQUE)",
	"channelTagRel(" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"c_id INTEGER NOT NULL, " +
		"ct_id INTEGER NOT NULL, " +
		"FOREIGN KEY (c_id) REFERENCES channel(channelHash) ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (ct_id) REFERENCES channelTag(id) ON DELETE CASCADE ON UPDATE CASCADE)",
	"playlist(" +
		"id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, " +
		"playlistTitle VARCHAR(100) NOT NULL, " +
		"channelId INTEGER REFERENCES channel(id))",
	"content(" +
		"contentHash CHAR(20) PRIMARY KEY, " +
		"tile VARCHAR(100) NOT NULL, " +
		"seriesName VARCHAR(100) NOT NULL, " +
		"partName VARCHAR(100) NOT NULL)",
	"contentTag(" +
		"id INTEGER PRIMARY KEY UNIQUE, " +
		"name VARCHAR(100) NOT NULL UNIQUE)",
	"contentTagRel(" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"c_id INTEGER NOT NULL, " +
		"ct_id INTEGER NOT NULL, " +
		"FOREIGN KEY (c_id) REFERENCES content(contentHash) ON DELETE CASCADE ON UPDATE CASCADE, " +
		"FOREIGN KEY (ct_id) REFERENCES contentTag(id) ON DELETE CASCADE ON UPDATE CASCADE)",
}

func CreateDB() error {
	dir := util.GetHomeDir() + constants.HIDDEN_DIR
	_, err := os.Stat(dir)
	// create directory if not exists
	if os.IsNotExist(err) {
		fmt.Println("SQL:CreateDB: Creating Dir [" + dir + "]")
		os.MkdirAll(dir, 0666)
	} else {
		fmt.Println("SQL:CreateDB: Dir Already Exists")
	}

	_, err = os.Stat(dir + constants.SQL_DB)

	if os.IsNotExist(err) {
		fmt.Println("SQL:CreateDB: CREATING DB [" + dir + constants.SQL_DB + "]")
		file, err := os.OpenFile(dir+constants.SQL_DB, os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		file.Close()
	} else {
		fmt.Println("SQL:CreateDB: DB Already Exists")
	}

	db, err := sql.Open("sqlite3", dir+constants.SQL_DB)
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		return err
	}

	//create all tables if they do not exist
	err = createAllTables(db)
	if err != nil {
		return err
	}

	fmt.Println("SQL:CreateDB: DB CLOSE")
	db.Close()

	return err
}

func createAllTables(db *sql.DB) error {
	fmt.Printf("SQL:CreateAllTables: Creating tables [%d]\n", len(TABLE_NAMES))

	for _, element := range CREATE_TABLE {
		//check if table exists
		s := "CREATE TABLE IF NOT EXISTS " + element + ";"
		fmt.Printf("SQL:CreateAllTables: [%s]\n", s)
		_, err := db.Exec(s)
		if err != nil {
			fmt.Printf("SQL:CreateAllTables: [%s]\n", err)
			return err
		}
	}
	fmt.Println("SQL:CreateAllTables: Finished Creating tables")
	return nil
}

func DeleteDB() error {
	dir := util.GetHomeDir() + constants.HIDDEN_DIR
	_, err := os.Stat(dir)
	// create directory if not exists
	if os.IsNotExist(err) {
		fmt.Println("SQL:DELETE: Creating Dir [" + dir + "]")
		os.MkdirAll(dir, 0666)
	} else {
		fmt.Println("SQL:DELETE: Dir Already Exists")
	}

	_, err = os.Stat(dir + constants.SQL_DB)

	if os.IsNotExist(err) {
		fmt.Println("SQL:DELETE: DB Already Deleted [" + dir + constants.SQL_DB + "]")
	} else {
		fmt.Println("SQL:DELETE: DELETING DB [" + dir + constants.SQL_DB + "]")
		err = os.Remove(dir + constants.SQL_DB)
		if err != nil {
			return err
		}
	}
	return nil
}
