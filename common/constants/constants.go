package constants

import (
	"os"
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
	FILE_PATH_MAX_LENGTH         int = 400
	TITLE_MAX_LENGTH             int = 100
	URL_MAX_LENGTH               int = 20
	TAG_MAX_LENGTH               int = 100
)

// Common type constants
const (
	MAX_CONTENT_TAGS int = 4
	MAX_CHANNEL_TAGS int = 4
)

type ConstantJSON struct {
	HashLength             int `json:"hashlength"`
	FileNameLength         int `json:"filenamelength"`
	LongDescriptionLength  int `json:"shortdesclength"`
	ShortDescriptionLength int `json:"longdesclength"`
	TrackerUrlLength       int `json:"trackerurllength"`
	FilePathLength         int `json:"filepathlength"`
	TitleLength            int `json:"titlelength"`
	UrlLength              int `json:"urllength"`
	MaxChannelTags         int `json:"maxchanneltags"`
	MaxContentTags         int `json:"maxcontenttags"`
}

func ConstantJSONMarshal() *ConstantJSON {
	c := new(ConstantJSON)
	c.HashLength = 64
	c.FileNameLength = FILE_NAME_MAX_LENGTH
	c.LongDescriptionLength = LONG_DESCRIPTION_MAX_LENGTH
	c.ShortDescriptionLength = SHORT_DESCRIPTION_MAX_LENGTH
	c.TrackerUrlLength = TRACKER_URL_MAX_LENGTH
	c.FilePathLength = FILE_PATH_MAX_LENGTH
	c.TitleLength = TITLE_MAX_LENGTH
	c.UrlLength = URL_MAX_LENGTH
	c.MaxChannelTags = MAX_CHANNEL_TAGS
	c.MaxContentTags = MAX_CONTENT_TAGS
	return c
}

// Channel Status Values
const (
	CHANNEL_NOT_READY int = iota // Missing critcal items to be in blockchain
	CHANNEL_READY                // Can be put in blockchain
	CHANNEL_FULL                 // Has all elements
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
	HIDDEN_DIR  = ".DistroSols/"
	SQL_DB      = "sql.db"
	LVL2_CACHE  = "level_two.cache"
	TORRENT_DIR = "torrent/"
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
	SQL_TABLE_CHANNEL__TITLE = "title"
	SQL_TABLE_CHANNEL__DT    = "dt"

	//CHANNEL TAG
	SQL_TABLE_CHANNEL_TAG__ID   = "id"
	SQL_TABLE_CHANNEL_TAG__NAME = "name"

	//CHANNEL + TAG + REL
	SQL_TABLE_CHANNEL_TAG_REL__C_ID  = "c_id"
	SQL_TABLE_CHANNEL_TAG_REL__CT_ID = "ct_id"

	//CONTENT
	SQL_TABLE_CONTENT__CONTENT_HASH = "contentHash"
	SQL_TABLE_CONTENT__TITLE        = "tile"
	SQL_TABLE_CONTENT__SERIES_NAME  = "seriesName"
	SQL_TABLE_CONTENT__PART_NAME    = "partName"
	SQL_TABLE_CONTENT__CH_ID        = "ch_id"
	SQL_TABLE_CONTENT__DT           = "dt"

	//CONTENT TAG
	SQL_TABLE_CONTENT_TAG__ID   = "id"
	SQL_TABLE_CONTENT_TAG__NAME = "name"

	//CONTENT + TAG + REL
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

// Constant Tags
var ALLOWED_TAGS = []string{"action", "drama", "comedy", "sports", "classic", "crime", "horror"}
