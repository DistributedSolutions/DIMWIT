package objects

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

// Channel wrapper contains a channel, and all factom metadata needed to verify
// factom entries
type ChannelWrapper struct {
	Channel common.Channel

	// Factom Metadata
	// Root
	RRegistered bool
	RMadeHeight uint32
	// Manage
	MRegistered bool
	MMadeHeight uint32
	// Content
	CRegistered bool
	CMadeHeight uint32
}

func NewChannelWrapper() *ChannelWrapper {
	cw := new(ChannelWrapper)
	c := common.NewChannel()
	cw.Channel = *c

	return cw
}

func RandomChannelWrapper() *ChannelWrapper {
	cw := new(ChannelWrapper)
	cw.Channel = *common.RandomNewChannel()

	cw.RRegistered = random.RandomBool()
	cw.RMadeHeight = random.RandomUInt32Between(0, 1000)

	cw.MRegistered = random.RandomBool()
	cw.MMadeHeight = random.RandomUInt32Between(0, 1000)

	cw.CRegistered = random.RandomBool()
	cw.CMadeHeight = random.RandomUInt32Between(0, 1000)

	return cw
}

func (a *ChannelWrapper) IsSameAs(b *ChannelWrapper) bool {
	if !a.Channel.IsSameAs(&b.Channel) {
		return false
	}

	if a.RRegistered != b.RRegistered {
		return false
	}

	if a.RMadeHeight != b.RMadeHeight {
		return false
	}

	if a.MRegistered != b.MRegistered {
		return false
	}

	if a.MMadeHeight != b.MMadeHeight {
		return false
	}

	if a.CRegistered != b.CRegistered {
		return false
	}

	if a.CMadeHeight != b.CMadeHeight {
		return false
	}

	return true
}

func (cw *ChannelWrapper) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := cw.Channel.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data = primitives.BoolToBytes(cw.RRegistered)
	buf.Write(data)

	data = primitives.Uint32ToBytes(cw.RMadeHeight)
	buf.Write(data)

	data = primitives.BoolToBytes(cw.MRegistered)
	buf.Write(data)

	data = primitives.Uint32ToBytes(cw.MMadeHeight)
	buf.Write(data)

	data = primitives.BoolToBytes(cw.CRegistered)
	buf.Write(data)

	data = primitives.Uint32ToBytes(cw.CMadeHeight)
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (cw *ChannelWrapper) UnmarshalBinary(data []byte) error {
	_, err := cw.UnmarshalBinaryData(data)
	return err
}

func (cw *ChannelWrapper) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[ChannelWrapper] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()
	newData = data

	newData, err = cw.Channel.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	t := primitives.ByteToBool(newData[0])
	cw.RRegistered = t
	u, err := primitives.BytesToUint32(newData[1:5])
	if err != nil {
		return data, err
	}
	cw.RMadeHeight = u
	newData = newData[5:]

	t = primitives.ByteToBool(newData[0])
	cw.MRegistered = t
	u, err = primitives.BytesToUint32(newData[1:5])
	if err != nil {
		return data, err
	}
	cw.MMadeHeight = u
	newData = newData[5:]

	t = primitives.ByteToBool(newData[0])
	cw.CRegistered = t
	u, err = primitives.BytesToUint32(newData[1:5])
	if err != nil {
		return data, err
	}
	cw.CMadeHeight = u
	newData = newData[5:]

	return
}
