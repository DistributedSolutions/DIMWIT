package database_test

import (
	"fmt"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"os"
	"strings"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/database"
)

var _ = fmt.Sprintf("")

var testDB *DB

func TestCreateDB(t *testing.T) {
	err := DeleteDB(constants.SQL_DB)
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(constants.SQL_DB)
	if !os.IsNotExist(err) {
		t.Error(fmt.Errorf("Error DB not deleted"))
	}

	testDB, err = CreateDB(constants.SQL_DB, CREATE_TABLE)
	if err != nil {
		t.Error(err)
	}

	if len(TABLE_NAMES) != len(CREATE_TABLE) {
		t.Error(fmt.Errorf("TableNames [%d] and CreateTable [%d] not same length\n", len(TABLE_NAMES), len(CREATE_TABLE)))
	}

	//check if tables were created right
	q := "SELECT tbl_name, sql FROM sqlite_master WHERE type='table' AND tbl_name NOT LIKE 'sqlite_sequence'"
	rows, err := testDB.DB.Query(q)
	if err != nil {
		t.Error(fmt.Errorf("Error retrieving tables\n"))
	}
	defer rows.Close()

	var name string
	var sql string
	row := 0
	for rows.Next() {
		err := rows.Scan(&name, &sql)
		if err != nil {
			t.Error(fmt.Errorf("Error reading names of tables and sql values\n"))
		}
		// fmt.Printf("Name [%s] sql [%s] row[%d] len[%d]]\n", name, sql, row, len(TABLE_NAMES))
		if row >= len(TABLE_NAMES) {
			fmt.Errorf("Error rows pass the amount of tables created extra name[%s] extra sql[%s] row count[%d]",
				name,
				sql,
				row)
		} else if !strings.EqualFold(TABLE_NAMES[row], name) || !strings.EqualFold(CREATE_TABLE[row], sql) {
			fmt.Errorf("Error tables were not created correct [%s] vs [%s] and [%s] vs [%s]\n",
				name,
				TABLE_NAMES[row],
				sql,
				CREATE_TABLE[row])
		}
		row++
	}
	if err := rows.Err(); err != nil {
		t.Error(fmt.Printf("ERROR when retrieving rows from playlistTemp\n"))
	}
}

func TestAddTags(t *testing.T) {
	err := DeleteTags(testDB.DB)
	if err != nil {
		t.Error(err)
	}

	err = AddTags(testDB.DB)
	if err != nil {
		t.Error(err)
	}
}

func TestAddChannel(t *testing.T) {
	c := common.RandomNewChannel()
	err := AddChannel(testDB.DB, c, -1) //HEIGHT PLACEHOLDER IN SECONd VALUE CURRENT IS TEMPORARY
	if err != nil {
		t.Error(err)
	}
}

func TestFlushPlaylistTemp(t *testing.T) {
	err := FlushPlaylistTempTable(testDB.DB, -1) //HEIGHT PLACEHOLDER IN SECONd VALUE CURRENT IS TEMPORARY
	if err != nil {
		t.Error(err)
	}
}

func TestCloseDB(t *testing.T) {
	CloseDB(testDB.DB)
}
