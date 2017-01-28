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
