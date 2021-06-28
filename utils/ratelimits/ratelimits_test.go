package ratelimits

import (
	"simple-nginx-otp/utils/config"
	"testing"
)

func TestIsLimited(t *testing.T) {
	conf := &config.Config{
		RateLimitCount:  2,
		RateLimitExpiry: 1,
	}
	limited := IsLimited(conf, "0.0.0.0")
	if limited {
		t.Error("limited too early")
	}
	limited = IsLimited(conf, "0.0.0.0")
	if limited {
		t.Error("limited early")
	}
	limited = IsLimited(conf, "0.0.0.0")
	if !limited {
		t.Error("limited late")
	}
}
