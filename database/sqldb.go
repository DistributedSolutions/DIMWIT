package database

import (
	"database/sql"
	"fmt"
	// _ "github.com/mattn/go-sqlite3"
	"os"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/util"
)

func CreateDB() error {
	dir := util.GetHomeDir() + constants.HIDDEN_DIR
	println(dir)
	d, err := os.Stat(dir)
	println(d)
	println(err)
	// create directory if not exists
	if os.IsNotExist(err) {
		println("CREATING DIR")
		os.MkdirAll(dir, 0666)
	} else {
		println("DIR ALREADY EXISTS")
	}

	d, err = os.Stat(dir + constants.SQL_DB)

	if os.IsNotExist(err) {
		println("CREATING DB")
		file, err := os.OpenFile(dir+constants.SQL_DB, os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		file.Close()
	} else {
		println("DB ALREADY EXISTS")
	}

	fmt.Println(dir)
	db, err := sql.Open("sqlite3", dir+constants.SQL_DB)
	if err != nil {
		return err
	}

	s := "CREATE TABLE test (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"value INTEGER NOT NULL);"
	res, err := db.Exec(s)
	println(s)
	fmt.Println(res)
	if err != nil {
		return err
	}

	s = "INSERT INTO Foo (value) VALUES (100);"
	res, err = db.Exec(s)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	fmt.Println(id)
	if err != nil {
		return err
	}

	db.Close()

	return err
}
