package smhb

import (
	"fmt"
	"testing"
)

func TestGetNTPTime(t *testing.T) {
	result, err := getNTPTime()

	if err != nil {
		t.Error(err)
		return
	}

	stamp := (*result).String()
	fmt.Printf(stamp, "\n")
}
