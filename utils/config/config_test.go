package config

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	_, err := GetConfig()
	if err == nil {
		t.Error("expected error, but no error returned")
	}
	err = os.Setenv("SNO_SECRET", "test")
	if err != nil {
		t.Error(err)
	}
	_, err = GetConfig()
	if err != nil {
		t.Error(err)
	}
}
