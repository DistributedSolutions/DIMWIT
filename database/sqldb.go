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

var _ = sqlite3.ErrNoMask

var gDB *sql.DB

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
		"partName VARCHAR(100) NOT NULL, " +
		"ch_id INTEGER NOT NULL, " +
		"FOREIGN KEY (ch_id) REFERENCES channel(channelHash) ON DELETE CASCADE ON UPDATE CASCADE)",
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

func getDBPath() (string, error) {
	path := util.GetHomeDir() + constants.HIDDEN_DIR + constants.SQL_DB
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("SQL:GetDB: DB DOES NOT Exists")
		return "", fmt.Errorf("DB Does NOT Exists: %s", err.Error())
	} else {
		fmt.Println("SQL:GetDB: DB Path exists :)")
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

func AddTags() error {
	gDB, err := getDB()
	if err != nil {
		return fmt.Errorf("Error adding tags: %s", err.Error())
	}

	s := ""
	for _, t := range constants.ALLOWED_TAGS {
		s += "INSERT INTO channelTag(name) VALUES('" + t + "');"
	}

	fmt.Println("SQL:AddTags: Tags[" + s + "]")
	_, err = gDB.Exec(s)
	if err != nil {
		return fmt.Errorf("Error inserting tags: %s", err.Error())
	}
	return nil
}

func DeleteTags() error {
	gDB, err := getDB()
	if err != nil {
		return fmt.Errorf("Error updating tags: %s", err.Error())
	}

	_, err = gDB.Exec("DELETE FROM channelTag")
	if err != nil {
		return fmt.Errorf("Error deleting tags: %s", err.Error())
	}
	return nil
}

//NOT GOING TO CHECK IF IN DB ALREADY
func AddChannel(channel *common.Channel) error {
	gDB, err := getDB()
	if err != nil {
		return fmt.Errorf("Error adding channel: %s", err.Error())
	}

	_, err = gDB.Exec("INSERT INTO channel(channelHash,tile) VALUES(?,?)",
		channel.RootChainID.String(),
		channel.ChannelTitle.String())
	if err != nil {
		return fmt.Errorf("Error inserting channel: %s", err.Error())
	}

	tags := channel.Tags.GetTags()
	fmt.Printf("SQL: Tags for channel len [%d]\n", len(tags))
	for _, t := range tags {
		tag, err := getTagID(gDB, t.String())
		if err != nil {
			return fmt.Errorf("Error retrieving tag: %s", err.Error())
		}
		s := "INSERT INTO channelTagRel(c_id,ct_id) VALUES(" +
			channel.RootChainID.String() +
			"," +
			tag +
			")"
		_, err = gDB.Exec(s)
		if err != nil {
			return fmt.Errorf("Error inserting tag for channel: %s", err.Error())
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
		return "", fmt.Errorf("Error deleting channel: %s", err.Error())
	}

	defer rows.Close()

	rows.Next()

	var result string
	rows.Scan(&result)

	return result, nil
}
