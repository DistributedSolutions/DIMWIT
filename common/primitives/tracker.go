package primitives

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type TrackerList struct {
	length   uint32
	trackers []Tracker
}

func NewTrackerList() *TrackerList {
	tl := new(TrackerList)
	tl.trackers = make([]Tracker, 0)
	tl.length = 0

	return tl
}

func RandomTrackerList(max uint32) *TrackerList {
	tl := NewTrackerList()
	l := random.RandomUInt32Between(0, max)

	var i uint32
	for i = 0; i < l; i++ {
		tl.AddNewTracker(RandomTracker())
	}

	return tl
}

func (tl *TrackerList) AddNewTracker(tracker *Tracker) error {
	tl.length++
	tl.trackers = append(tl.trackers, *tracker)

	return nil
}

func (tl *TrackerList) AddNewTrackerURL(url string) error {
	t, err := NewTracker(url)
	if err != nil {
		return err
	}

	tl.length++
	tl.trackers = append(tl.trackers, *t)

	return nil
}

func (tl *TrackerList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := Uint32ToBytes(tl.length)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	for _, t := range tl.trackers {
		data, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (tl *TrackerList) UnmarshalBinary(data []byte) error {
	_, err := tl.UnmarshalBinaryData(data)
	return err
}

func (tl *TrackerList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	newData = data

	u, err := BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	tl.length = u
	newData = newData[4:]

	tl.trackers = make([]Tracker, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = tl.trackers[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

func (a *TrackerList) IsSameAs(b *TrackerList) bool {
	if a.length != b.length {
		return false
	}

	for i := range a.trackers {
		if !a.trackers[i].IsSameAs(&(b.trackers[i])) {
			return false
		}
	}

	return true
}

type Tracker string

func NewTracker(url string) (*Tracker, error) {
	d := new(Tracker)

	err := d.SetString(url)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (t *Tracker) SetString(url string) error {
	if len(url) > t.MaxLength() {
		return fmt.Errorf("Tracker url given is too long, length must be under %d, given length is %d",
			t.MaxLength(), len(url))
	}

	*t = Tracker(url)
	return nil
}

func (t *Tracker) String() string {
	return string(*t)
}

func (d *Tracker) MaxLength() int {
	return constants.TRACKER_URL_MAX_LENGTH
}

func (a *Tracker) IsSameAs(b *Tracker) bool {
	return a.String() == b.String()
}

func (d *Tracker) MarshalBinary() ([]byte, error) {
	return MarshalStringToBytes(d.String(), d.MaxLength())
}

func (t *Tracker) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (d *Tracker) UnmarshalBinaryData(data []byte) ([]byte, error) {
	newData := data
	str, newData, err := UnmarshalStringFromBytesData(newData, d.MaxLength())
	if err != nil {
		return data, err
	}

	err = d.SetString(str)
	if err != nil {
		return data, err
	}

	return newData, nil
}

func RandomTracker() *Tracker {
	l, _ := NewTracker("")
	l.SetString(random.RandStringOfSize(l.MaxLength()))

	return l
}
