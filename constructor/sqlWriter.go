package constructor

import (
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
)

const (
	LOOP_DELAY time.Duration = time.Duration(1 * time.Second)
)

type SqlWriter struct {
	// Incoming channels to write to sql db
	channelQueue chan common.Channel

	// Stop goroutine
	quit chan int
}

// Called to make SQLWriter
func NewSqlWriter() *SqlWriter {
	sw := new(SqlWriter)
	sw.quit = make(chan int, 5)
	sw.channelQueue = make(chan common.Channel, 1000)

	return sw
}

// Close sqlwriter
func (sw *SqlWriter) Close() {
	sw.quit <- 0
}

// Called to run SQLWriter
func (sw *SqlWriter) DrainChannelQueue() {
	for {
		// Closeing sqlwrite
		select {
		case <-sw.quit:
			// Add your close code here
			return
		default:
		}

		// Take incoming channels
		select {
		case channel := <-sw.channelQueue:
			channelList := make([]common.Channel, 0)
			channelList = append(channelList, channel)
			// Do stuff
			length := len(sw.channelQueue)
			for i := 0; i < length; i++ {
				select {
				case newChan := <-sw.channelQueue:
					channelList = append(channelList, newChan)
				default:
				}
			}

			// ChannelList, play with it
			// JESSE! IMPLEMENT
		default:
			// Nothing really
		}

		// Don't starve other routines
		time.Sleep(LOOP_DELAY)
	}
}
