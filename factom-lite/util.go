package lite

import (
	"bytes"

	"github.com/FactomProject/factom"
)

func AreEntriesSame(a *factom.Entry, b *factom.Entry) bool {
	if a.ChainID != b.ChainID {
		return false
	}

	if len(a.ExtIDs) != len(b.ExtIDs) {
		return false
	}

	for i := range a.ExtIDs {
		if bytes.Compare(a.ExtIDs[i], b.ExtIDs[i]) != 0 {
			return false
		}
	}

	return true
}

func GetECAddress() *factom.ECAddress {
	//ec := factom.NewECAddress()
	ec, _ := factom.GetECAddress("Es3yhERbyDwQiswxCsb24PyeLbHUgwnq1tiTJnyGChorYdj9c8Lt")
	return ec
}
