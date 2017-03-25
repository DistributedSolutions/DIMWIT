package main

import (
	"github.com/DistributedSolutions/DIMWIT/engine"
	log "github.com/DistributedSolutions/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("Debugging on")

	engine.GrabFlagsAndRun()
}
