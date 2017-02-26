package database

import (
	"fmt"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"os"
	"strings"
	"testing"
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
	//add in tags to db that was just created
	err = AddTags(testDB.DB)
	if err != nil {
		t.Error(err)
	}
	//deletes the tags from the newly created db
	err := DeleteTags(testDB.DB)
	if err != nil {
		t.Error(err)
	}

	//check if tables were created right
	var count int
	s := "SELECT COUNT(" + constants.SQL_TABLE_CONTENT_TAG__ID + ") FROM " + constants.SQL_CONTENT_TAG
	err = testDB.DB.QueryRow(s).Scan(&count)
	if err != nil {
		t.Error(fmt.Errorf("Error deleting content tags [%s]: [%s]\n", s, err))
	} else if count != 0 {
		t.Error(fmt.Errorf("All content tags not deleted with query [%s] and count [%d]\n", s, count))
	}

	s = "SELECT COUNT(" + constants.SQL_TABLE_CHANNEL_TAG__ID + ") FROM " + constants.SQL_CHANNEL_TAG
	err = testDB.DB.QueryRow(s).Scan(&count)
	if err != nil {
		t.Error(fmt.Errorf("Error deleting channel tags [%s]: [%s]\n", s, err))
	} else if count != 0 {
		t.Error(fmt.Errorf("All channel tags not deleted with query [%s] and count [%d]\n", s, count))
	}

	err = AddTags(testDB.DB)
	if err != nil {
		t.Error(err)
	}

	//check that the tags were created for both tables
	q := "SELECT " + constants.SQL_TABLE_CONTENT_TAG__ID + " FROM " + constants.SQL_CONTENT_TAG +
		" WHERE " + constants.SQL_TABLE_CONTENT_TAG__NAME + " NOT IN (" + CommaDelimiterArray(constants.ALLOWED_TAGS, true) + ")"
	rows, err := testDB.DB.Query(q)
	if err != nil {
		t.Error(fmt.Errorf("Error retrieving tables with query[%s]: %s\n", q, err))
	}
	for rows.Next() {
		t.Error(fmt.Errorf("Error extra content rows\n"))
	}
	rows.Close()

	q = "SELECT " + constants.SQL_TABLE_CHANNEL_TAG__ID + " FROM " + constants.SQL_CHANNEL_TAG + " " +
		" WHERE " + constants.SQL_TABLE_CHANNEL_TAG__NAME + " NOT IN (" + CommaDelimiterArray(constants.ALLOWED_TAGS, true) + ")"
	rows, err = testDB.DB.Query(q)
	if err != nil {
		t.Error(fmt.Errorf("Error retrieving tables[%s]: [%s]\n", q, err))
	}
	for rows.Next() {
		t.Error(fmt.Errorf("Error extra content rows\n"))
	}
	rows.Close()
}

func TestAddChannel(t *testing.T) {

	//////////ADD 3 RANDOM CHANNELS/////////////////
	c := common.RandomNewChannel()
	channels := make([]*common.Channel, 1, 1)
	channels[0] = c
	err := AddChannelArr(testDB.DB, channels, -1) //HEIGHT PLACEHOLDER IN SECOND VALUE CURRENT IS TEMPORARY
	if err != nil {
		t.Error(err)
	}

	//check if channel tags were added
	tags := c.Tags.GetTagsAsStringArr()
	q := "SELECT " + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID +
		" FROM " + constants.SQL_CHANNEL_TAG_REL +
		" WHERE " + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + " NOT IN(" +
		testSubQueryGetTags(constants.SQL_CHANNEL_TAG, constants.SQL_TABLE_CHANNEL_TAG__ID, constants.SQL_TABLE_CHANNEL_TAG__NAME, tags) +
		") AND " + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + " = '" + c.ContentChainID.String() + "'"
	rows, err := testDB.DB.Query(q)
	if err != nil {
		t.Error(fmt.Errorf("Error retrieving tables with query[%s]: %s\n", q, err.Error()))
	}
	for rows.Next() {
		t.Error(fmt.Errorf("Error extra channel tag rows with query[%s]\n", q))
	}
	rows.Close()

	//check if content was added
	contents := c.Content.GetContents()
	for _, content := range contents {
		var count int
		//check if content was added
		q = "SELECT COUNT(" + constants.SQL_TABLE_CONTENT__CONTENT_HASH + ") FROM " + constants.SQL_CONTENT +
			" WHERE " + constants.SQL_TABLE_CONTENT__CONTENT_HASH + " = '" + content.ContentID.String() + "'"
		err := testDB.DB.QueryRow(q).Scan(&count)
		if err != nil {
			t.Error(fmt.Errorf("Error retrieving tables with query[%s]: [%s]\n", q, err))
		}
		if count != 1 {
			t.Error(fmt.Errorf("Error content not found with query [%s] and count [%d]\n", q, count))
		}

		//check if content tags were set
		tags = content.Tags.GetTagsAsStringArr()
		q = "SELECT " + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + " FROM " + constants.SQL_CONTENT_TAG_REL +
			" WHERE " + constants.SQL_TABLE_CONTENT_TAG_REL__CT_ID + " NOT IN(" +
			testSubQueryGetTags(constants.SQL_CONTENT_TAG, constants.SQL_TABLE_CONTENT_TAG__ID,
				constants.SQL_TABLE_CONTENT_TAG__NAME, tags) +
			") AND " + constants.SQL_TABLE_CONTENT_TAG_REL__C_ID + " = '" + content.ContentID.String() + "'"
		rows, err := testDB.DB.Query(q)
		if err != nil {
			t.Error(fmt.Errorf("Error retrieving tables with query[%s]: [%s]\n", q, err))
		}
		for rows.Next() {
			t.Error(fmt.Errorf("Error extra content tag rows with query[%s]\n", q))
		}
		rows.Close()
	}

	//NO need to check if the playlist was inserted, would have thrown an error if it was not inserted correctly
}

func TestFlushPlaylistTemp(t *testing.T, height int) {
	err := FlushPlaylistTempTable(testDB.DB, height) //HEIGHT PLACEHOLDER IN SECONd VALUE CURRENT IS TEMPORARY
	if err != nil {
		t.Error(err)
	}
	//*************WILL HAVE TO CHANGE IN FUTURE WHEN HEIGHT MANEGEMENT IS CHANGED******************
}

func TestCloseDB(t *testing.T) {
	CloseDB(testDB.DB)
}

func testSubQueryGetTags(tableName string, colOn string, colWhere string, tags []string) string {
	q := "SELECT " + colOn +
		" FROM " + tableName +
		" WHERE " + colWhere + " IN (" + CommaDelimiterArray(tags, true) + ")"
	return q
}
