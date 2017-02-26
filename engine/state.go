package engine

import (
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

type WholeState struct {
	Constructor  *constructor.Constructor
	FactomClient lite.FactomLite
}
