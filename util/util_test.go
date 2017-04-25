package util_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/util"
)

var _ = fmt.Sprintf("")

func TestHomeDor(t *testing.T) {
	GetHomeDir()
}

func TestDashDelimiterToCamelCase(t *testing.T) {
	var retVal string
	given := []string{"apples-are-delicious", "a-b", "", " ", " yes-sir1", "yes-sir2 ", " yes-size3 "}
	expected := []string{"ApplesAreDelicious", "AB", "", "", "YesSir1", "YesSir2", "YesSize3"}
	for i, val := range given {
		retVal = DashDelimiterToCamelCase(val)
		if retVal != expected[i] {
			t.Errorf("Recieved %s, Expected %s\n", retVal, expected[i])
		}
	}
}
