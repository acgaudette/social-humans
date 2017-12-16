package smhb

import (
	"testing"
)

func TestGetNTPTime(t *testing.T) {
	result, err := getNTPTime()

	if err != nil {
		t.Error(err)
		return
	}

	stamp := (*result).String()
	t.Logf(stamp)
}
