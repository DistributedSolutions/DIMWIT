package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestLongDesc(t *testing.T) {
	d, err := NewLongDescription("hello")
	if err != nil {
		t.Error(err)
	}

	if d.String() != "hello" {
		t.Error("String was not set")
	}

	var _ = d
}
