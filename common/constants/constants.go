package constants

// Primitive type constants
// Byte length can be max length + 1 for strings
const (
	INFOHASH_LENGTH              int = 20
	HASH_LENGTH                  int = 32
	LONG_DESCRIPTION_MAX_LENGTH  int = 100
	SHORT_DESCRIPTION_MAX_LENGTH int = 20
	FILE_NAME_MAX_LENGTH         int = 50
	MD5_CHECKSUM_LENGTH          int = 16
	TRACKER_URL_MAX_LENGTH       int = 100
	TITLE_MAX_LENGTH             int = 30
	URL_MAX_LENGTH               int = 20
	TAG_MAX_LENGTH               int = 15
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

// Sql types
const (
	SQL_STRING byte = iota
	SQL_OTHER  byte = iota
)

// Sql Table Names
const (
	SQL_CHANNEL         = "channel"
	SQL_CHANNEL_TAG     = "channelTag"
	SQL_CHANNEL_TAG_REL = "channelTagRel"
	SQL_PLAYLIST        = "playlist"
	SQL_CONTENT         = "content"
	SQL_CONTENT_TAG     = "contentTag"
	SQL_CONTENT_TAG_REL = "contentTagRel"
)

// Sql Table Cols
const (
	SQL_TABLE_CHANNEL__HASH  = "channelHash"
	SQL_TABLE_CHANNEL__TITLE = "tile"

	SQL_TABLE_CHANNEL_TAG__ID   = "id"
	SQL_TABLE_CHANNEL_TAG__NAME = "name"

	SQL_TABLE_CHANNEL_TAG_REL__ID    = "id"
	SQL_TABLE_CHANNEL_TAG_REL__C_ID  = "c_id"
	SQL_TABLE_CHANNEL_TAG_REL__CT_ID = "ct_id"

	SQL_TABLE_PLAYLIST__ID             = "id"
	SQL_TABLE_PLAYLIST__PLAYLIST_TITLE = "playlistTitle"
	SQL_TABLE_PLAYLIST__CHANNEL_ID     = "channelId"

	SQL_TABLE_CONTENT__CONTENT_HASH = "contentHash"
	SQL_TABLE_CONTENT__TITLE        = "tile"
	SQL_TABLE_CONTENT__SERIES_NAME  = "seriesName"
	SQL_TABLE_CONTENT__PART_NAME    = "partName"
	SQL_TABLE_CONTENT__CH_ID        = "ch_id"

	SQL_TABLE_CONTENT_TAG__ID   = "id"
	SQL_TABLE_CONTENT_TAG__NAME = "name"

	SQL_TABLE_CONTENT_TAG_REL__ID    = "id"
	SQL_TABLE_CONTENT_TAG_REL__C_ID  = "c_id"
	SQL_TABLE_CONTENT_TAG_REL__CT_ID = "ct_id"
)

// For version bytes
const (
	FACTOM_VERSION byte = 0x00
)

// Constant Tags
var ALLOWED_TAGS = []string{"DIMWIT", "CLIT", "FRUIT", "Jesse", "Steve", "Go", "Node", "PEEEEEENNNNNIIIIIISSSSSS"}
