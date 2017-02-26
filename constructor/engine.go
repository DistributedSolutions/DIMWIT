package constructor

import (
	// "log"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

// The engine controls the constructor object. It manages the constructor state and go routines
// that read from the factom reader and update the appropriate databases

// StartConstructor has the constructor continuously check the next blocks for more information
func (c *Constructor) StartConstructor() {
	InitEnginePrometheus()
	for {
		select {
		case <-c.quit:
			return
		default:
			constructorEngineHeight.Set(float64(c.CompletedHeight))
			err := c.ApplyHeight(c.CompletedHeight + 1)
			if err != nil {
				// log.Println("[ConstructorError] ", err.Error())
				time.Sleep(constants.CHECK_FACTOM_FOR_UPDATES)
			}
		}
	}
}
