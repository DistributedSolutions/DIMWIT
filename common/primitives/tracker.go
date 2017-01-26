package primitives

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

type TrackerList struct {
	Length   int
	Trackers []Tracker
}

func NewTrackerList() *TrackerList {
	tl := new(TrackerList)
	tl.Trackers = make([]Tracker, 0)
	tl.Length = 0

	return tl
}

func (tl *TrackerList) AddNewTracker(url string) error {
	t, err := NewTracker(url)
	if err != nil {
		return err
	}

	tl.Length++
	tl.Trackers = append(tl.Trackers, *t)

	return nil
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
	if len(url) > constants.TRACKER_URL_MAX_LENGTH {
		return fmt.Errorf("Tracker url given is too long, length must be under %d, given length is %d", constants.TRACKER_URL_MAX_LENGTH, len(url))
	}

	*t = Tracker(url)
	return nil
}

func (t *Tracker) String() string {
	return string(*t)
}
