package rand

import (
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	randomString, err := GenerateRandomString(2)
	if err != nil {
		t.Error(err)
	}
	if len(randomString) != 2 {
		t.Errorf("`%s` not length of 2", randomString)
	}
}
