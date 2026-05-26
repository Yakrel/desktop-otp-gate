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
		CookieSecure:   true,
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
	if cookie.Path != "/" {
		t.Errorf("`%s` not `/`", cookie.Path)
	}
	if !cookie.HttpOnly {
		t.Error("cookie not HttpOnly")
	}
	if !cookie.Secure {
		t.Error("cookie not Secure")
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

func TestSessionCookie(t *testing.T) {
	conf := &config.Config{
		CookieName:    "test",
		CookieLength:  2,
		CookieMinutes: 60,
		SessionCookie: true,
	}
	session, cookie, err := NewSession(conf)
	if err != nil {
		t.Error(err)
	}
	if !cookie.Expires.IsZero() {
		t.Errorf("session cookie should not set Expires: `%s`", cookie.Expires)
	}
	after := time.Now().Add(59 * time.Minute)
	if after.After(session.Expiry) {
		t.Errorf("`%s` expires too early", session.Expiry)
	}
	before := time.Now().Add(61 * time.Minute)
	if before.Before(session.Expiry) {
		t.Errorf("`%s` expires too late", session.Expiry)
	}
}
