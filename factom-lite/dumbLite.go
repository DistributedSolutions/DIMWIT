package lite

import (
	"github.com/FactomProject/factom"
)

// A very dumb factom-lite client
type DumbLite struct {
	FactomdLocation string
}

func NewDumbLite() *DumbLite {
	d := new(DumbLite)
	d.FactomdLocation = "localhost:8088"
}

func (d *DumbLite) SubmitEntry(e factom.Entry, ec factom.ECAddress) (comId string, eHash string, err error) {
	comId, err = factom.CommitEntry(&e, &ec)
	if err != nil {
		return "", "", err
	}

	ehash, err = factom.RevealEntry(e)
	return
}
func (d *DumbLite) SubmitChain(c factom.Chain, ec factom.ECAddress) (comId string, chainID string, err error) {
	comId, err = factom.CommitChain(c, ec)
	if err != nil {
		return "", "", err
	}

	chainID, err = factom.RevealChain(c)
	return
}
