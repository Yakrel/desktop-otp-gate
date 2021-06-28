package sessions

import (
	"simple-nginx-otp/utils/config"
	"testing"
	"time"
)

func TestSessions(t *testing.T) {
	conf := &config.Config{
		CookieName:     "test",
		CookieLength:   2,
		CookieLifetime: 1,
		CookieDomain:   "test.example.com",
	}
	_, cookie, err := NewSession(conf)
	if err != nil {
		t.Error(err)
	}
	session := GetSession(cookie.Value)

	if cookie.Name != conf.CookieName {
		t.Errorf("`%s` not `%s`", cookie.Name, conf.CookieName)
	}
	if cookie.Domain != conf.CookieDomain {
		t.Errorf("`%s` not `%s`", cookie.Domain, conf.CookieDomain)
	}

	after := time.Now().Add(time.Duration(conf.CookieLifetime*24-1) * time.Hour)
	if after.After(session.Expiry) && after.After(cookie.Expires) {
		t.Errorf("`%s` expires too early", session.Expiry)
	}
	before := time.Now().Add(time.Duration(conf.CookieLifetime*24+1) * time.Hour)
	if before.Before(session.Expiry) && before.Before(cookie.Expires) {
		t.Errorf("`%s` expires too late", session.Expiry)
	}

	if session.Redirect != "/" {
		t.Errorf("`%s` not `/`", session.Redirect)
	}
	if session.Authorized {
		t.Error("session pre-authorized!")
	}
}
