package constants

import (
	"os"
)

var (
	MASTER_CHAIN_STRING string = "dcf089fc510ed0f6bdece3cb106665a9a581b85d1a6a7e8db1e2ac0b62aaa16e"
	MASTER_CHAIN_BYTES  []byte = []byte{0xdc, 0xf0, 0x89, 0xfc,
		0x51, 0x0e, 0xd0, 0xf6, 0xbd, 0xec, 0xe3,
		0xcb, 0x10, 0x66, 0x65, 0xa9, 0xa5, 0x81,
		0xb8, 0x5d, 0x1a, 0x6a, 0x7e, 0x8d, 0xb1,
		0xe2, 0xac, 0x0b, 0x62, 0xaa, 0xa1, 0x6e}
)

// Primitive type constants
// Byte length can be max length + 1 for strings
const (
	// Bytes
	INFOHASH_BYTES_LENGTH     int = 20 // String is 40 in length
	HASH_BYTES_LENGTH         int = 32 // String is 64 in length
	MD5_CHECKSUM_BYTES_LENGTH int = 16 // String is 32 in length

	// Strings
	LONG_DESCRIPTION_MAX_LENGTH  int = 100
	SHORT_DESCRIPTION_MAX_LENGTH int = 20
	FILE_NAME_MAX_LENGTH         int = 50
	TRACKER_URL_MAX_LENGTH       int = 100
	TITLE_MAX_LENGTH             int = 100
	URL_MAX_LENGTH               int = 20
	TAG_MAX_LENGTH               int = 100
)

var (
	CHAIN_PREFIX              []byte = []byte{0xDC, 0xF0, 0x00}
	CHAIN_PREFIX_LENGTH_CHECK int    = 2
	// How large can an entry be in factom in bytes
	ENTRY_MAX_SIZE int = 10240
)

// Common type constants
const (
	MAX_CONTENT_TAGS int = 4
	MAX_CHANNEL_TAGS int = 4
)

// Content Types
const (
	CONTENT_TYPE_VIDEO byte = iota
)

// Image Types
const (
	IMAGE_JPEG byte = iota
)

// Hidden File Directory
const (
	HIDDEN_DIR = ".DistroSols/"
	SQL_DB     = "sql.db"
)

const (
	DIRECTORY_PERMISSIONS os.FileMode = 0777
	FILE_PERMISSIONS      os.FileMode = 0777
)

// Sql Extras
const (
	SQL_HEIGHT_FLUSH_VALUE = 10
)

// Sql Table Names
const (
	SQL_CHANNEL              = "channel"
	SQL_CHANNEL_TAG          = "channelTag"
	SQL_CHANNEL_TAG_REL      = "channelTagRel"
	SQL_CONTENT              = "content"
	SQL_CONTENT_TAG          = "contentTag"
	SQL_CONTENT_TAG_REL      = "contentTagRel"
	SQL_PLAYLIST             = "playlist"
	SQL_PLAYLIST_CONTENT_REL = "playlistContentRel"
	SQL_PLAYLIST_TEMP        = "playlistTemp"
)

// Sql Table Cols
const (
	//CHANNEL
	SQL_TABLE_CHANNEL__HASH  = "channelHash"
	SQL_TABLE_CHANNEL__TITLE = "tile"

	//CHANNEL TAG
	SQL_TABLE_CHANNEL_TAG__ID   = "id"
	SQL_TABLE_CHANNEL_TAG__NAME = "name"

	//CHANNEL + TAG + REL
	SQL_TABLE_CHANNEL_TAG_REL__ID    = "id"
	SQL_TABLE_CHANNEL_TAG_REL__C_ID  = "c_id"
	SQL_TABLE_CHANNEL_TAG_REL__CT_ID = "ct_id"

	//CONTENT
	SQL_TABLE_CONTENT__CONTENT_HASH = "contentHash"
	SQL_TABLE_CONTENT__TITLE        = "tile"
	SQL_TABLE_CONTENT__SERIES_NAME  = "seriesName"
	SQL_TABLE_CONTENT__PART_NAME    = "partName"
	SQL_TABLE_CONTENT__CH_ID        = "ch_id"

	//CONTENT TAG
	SQL_TABLE_CONTENT_TAG__ID   = "id"
	SQL_TABLE_CONTENT_TAG__NAME = "name"

	//CONTENT + TAG + REL
	SQL_TABLE_CONTENT_TAG_REL__ID    = "id"
	SQL_TABLE_CONTENT_TAG_REL__C_ID  = "c_id"
	SQL_TABLE_CONTENT_TAG_REL__CT_ID = "ct_id"

	//PLAYLIST
	SQL_TABLE_PLAYLIST__ID             = "id"
	SQL_TABLE_PLAYLIST__PLAYLIST_TITLE = "title"
	SQL_TABLE_PLAYLIST__CHANNEL_ID     = "channelId"

	//PLAYLIST + CONTENT + REL
	SQL_TABLE_PLAYLIST_CONTENT_REL__ID    = "id"
	SQL_TABLE_PLAYLIST_CONTENT_REL__P_ID  = "p_id"
	SQL_TABLE_PLAYLIST_CONTENT_REL__CT_ID = "ct_id"

	//PLAYLIST TEMP
	SQL_TABLE_PLAYLIST_TEMP__ID         = "id"
	SQL_TABLE_PLAYLIST_TEMP__TITLE      = "title"
	SQL_TABLE_PLAYLIST_TEMP__HEIGHT     = "height"
	SQL_TABLE_PLAYLIST_TEMP__CHANNEL_ID = "channelId"
	SQL_TABLE_PLAYLIST_TEMP__CONTENT_ID = "contentId"
)

// For version bytes
const (
	FACTOM_VERSION byte = 0x00
)

// Constant Tags
var ALLOWED_TAGS = []string{"DIMWIT", "CLIT", "FRUIT", "Jesse", "Steve", "Go", "Node", "PEEEEEENNNNNIIIIIISSSSSS"}
