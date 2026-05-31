package sessions

import (
	"net/http"
	"simple-nginx-otp/utils/config"
	"simple-nginx-otp/utils/rand"
	"sync"
	"time"
)

type Session struct {
	Redirect   string
	Expiry     time.Time
	Authorized bool
}

var sessions = make(map[string]*Session)
var sessionsMutex = sync.Mutex{}

func init() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			sessionsMutex.Lock()
			_prune()
			sessionsMutex.Unlock()
		}
	}()
}

func NewSession(conf *config.Config) (*Session, *http.Cookie, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, err := rand.GenerateRandomString(conf.CookieLength)
	if err != nil {
		return nil, nil, err
	}

	cookie := new(http.Cookie)
	cookie.Name = conf.CookieName
	cookie.Value = session
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = conf.CookieSecure
	cookie.SameSite = http.SameSiteLaxMode
	expiry := time.Now().Add(time.Hour * time.Duration(24*conf.CookieLifetime))
	if conf.CookieMinutes > 0 {
		expiry = time.Now().Add(time.Minute * time.Duration(conf.CookieMinutes))
	}
	if !conf.SessionCookie {
		cookie.Expires = expiry
	}
	if conf.CookieDomain != "" {
		cookie.Domain = conf.CookieDomain
	}

	sessions[session] = &Session{
		Redirect:   "/",
		Authorized: false,
		Expiry:     expiry,
	}

	return sessions[session], cookie, nil
}

func GetSession(cookie string) *Session {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	_prune()
	if cookie == "" {
		return nil
	}
	session, ok := sessions[cookie]
	if !ok {
		return nil
	}
	return session
}

func _prune() {
	for cookie, session := range sessions {
		if time.Now().After(session.Expiry) {
			delete(sessions, cookie)
		}
	}
}
