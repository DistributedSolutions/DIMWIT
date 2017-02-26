package constructor

import (
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/constructor/objects"
)

const (
	LOOP_DELAY time.Duration = time.Duration(1 * time.Second)
)

type SqlWriter struct {
	// Incoming channels to write to sql db
	channelQueue chan objects.ChannelWrapper

	// Stop goroutine
	quit chan int
}

// Called to make SQLWriter
func NewSqlWriter() *SqlWriter {
	sw := new(SqlWriter)
	sw.quit = make(chan int, 5)
	sw.channelQueue = make(chan objects.ChannelWrapper, 1000)

	return sw
}

// Called to send a channel to the SQLWriter.
func (sw *SqlWriter) SendChannelDownQueue(c objects.ChannelWrapper) {
	// If you want to do anything to it before it hits the go-routine
	// do so here.
	// You have access to some extra variables with the Wrapper.
	// I don't think you care about them though.
	sw.channelQueue <- c
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
			channelList := make([][]common.Channel, 1)
			heightList := make([]uint32, 0)
			channelList[0] = append(channelList[0], channel.Channel)
			heightList = append(heightList, channel.CurrentHeight)
			// Do stuff
			length := len(sw.channelQueue)
			curIndex := 0
			for i := 0; i < length; i++ {
				select {
				case newChan := <-sw.channelQueue:
					if newChan.CurrentHeight == heightList[curIndex] {
						channelList[curIndex] = append(channelList[curIndex], newChan.Channel)
					} else {
						channelList = append(channelList, []common.Channel{newChan.Channel})
						heightList = append(heightList, newChan.CurrentHeight)
						curIndex++
					}
				default:
				}
			}

			// ChannelList, play with it
			// JESSE! IMPLEMENT
			// channelList is of type [i][ii]channel. Element i corrolates to heightList[i] which is the
			// height of all chanels in channelList[i].

			// So batch write loop looks like:
			// for i in channel list
			//		batchwrite channelList[i] with height = heightList[i]
			// endfor
		default:
			// Nothing really
		}

		// Don't starve other routines
		time.Sleep(LOOP_DELAY)
	}
}
