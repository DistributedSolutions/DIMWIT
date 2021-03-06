package constructor

import (
	"log"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/fatih/color"
)

// The engine controls the constructor object. It manages the constructor state and go routines
// that read from the factom reader and update the appropriate databases

var _ = log.Prefix()
var _ = time.Second
var _ = constants.MAX_CHANNEL_TAGS

// StartConstructor has the constructor continuously check the next blocks for more information
func (c *Constructor) StartConstructor() {
	InitEnginePrometheus()
	for {
		select {
		case <-c.quit:
			return
		default:
			h, err := c.GetReadyHeight()
			if err != nil {
				log.Println("Error getting ready height: %s", err.Error())
			}
			if (c.CompletedHeight + 1) <= h {
				//process
				// constructorEngineHeight.Set(float64(c.CompletedHeight))
				err := c.ApplyHeight(c.CompletedHeight + 1)
				if err != nil {
					// log.Println(err.Error())
					// time.Sleep(10 * time.Millisecond)
					color.Red("Error Applying Height: %s", err.Error())
				} else {
					// Height X was applied
				}
			} else {
				//do not process sleep instead
				time.Sleep(constants.CHECK_FACTOM_FOR_UPDATES)
			}
		}
	}
	// If I die, so does the SQLGuy. He should be dead at this point, but gatta be sure
	c.quit <- 0
}

func (c *Constructor) Kill() {
	c.quit <- 0
}
