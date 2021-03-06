package constants

import (
	"time"
)

//
// Constants used only inside Factom
//

// Master chain information. Where to find channels
var (
	MASTER_CHAIN_STRING string = "dcf089fc510ed0f6bdece3cb106665a9a581b85d1a6a7e8db1e2ac0b62aaa16e"
	MASTER_CHAIN_BYTES  []byte = []byte{0xdc, 0xf0, 0x89, 0xfc,
		0x51, 0x0e, 0xd0, 0xf6, 0xbd, 0xec, 0xe3,
		0xcb, 0x10, 0x66, 0x65, 0xa9, 0xa5, 0x81,
		0xb8, 0x5d, 0x1a, 0x6a, 0x7e, 0x8d, 0xb1,
		0xe2, 0xac, 0x0b, 0x62, 0xaa, 0xa1, 0x6e}
)

// Rules when creating a channel
var (
	CHAIN_PREFIX              []byte = []byte{0xDC, 0xF0, 0x00}
	CHAIN_PREFIX_LENGTH_CHECK int    = 1
	// How large can an entry be in factom in bytes
	ENTRY_MAX_SIZE int = 10240
)

// For version bytes
const (
	FACTOM_VERSION byte = 0x00
)

// Window we will allow a timestamp in an entry to differ from the dblock
const (
	ENTRY_TIMESTAMP_WINDOW int64 = 24 * 60 * 60 // in seconds
)

// Constructor Settings
var (
	CHECK_FACTOM_FOR_UPDATES time.Duration = time.Duration(5 * time.Second)
)

var (
	MAX_IMAGE_SIZE  uint32 = 1024
	MAX_BANNER_SIZE uint32 = 1024
)
