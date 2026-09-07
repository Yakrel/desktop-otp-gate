package sessions

import (
	"net/http"
	"simple-nginx-otp/utils/config"
	"simple-nginx-otp/utils/rand"
	"sync"
	"time"
)

type Session struct {
	Redirect     string
	Expiry       time.Time
	LastActivity time.Time
	IdleTimeout  time.Duration
	Authorized   bool
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
	var expiry time.Time
	if conf.CookieMinutes > 0 {
		expiry = time.Now().Add(time.Minute * time.Duration(conf.CookieMinutes))
	} else if conf.CookieLifetime > 0 {
		expiry = time.Now().Add(time.Hour * time.Duration(24*conf.CookieLifetime))
	} else {
		expiry = time.Now().Add(time.Hour * 24 * 14)
	}
	if !conf.SessionCookie {
		cookie.Expires = expiry
	}
	if conf.CookieDomain != "" {
		cookie.Domain = conf.CookieDomain
	}

	var idleTimeout time.Duration
	if conf.IdleMinutes > 0 {
		idleTimeout = time.Minute * time.Duration(conf.IdleMinutes)
	}

	sessions[session] = &Session{
		Redirect:     "/",
		Authorized:   false,
		Expiry:       expiry,
		LastActivity: time.Now(),
		IdleTimeout:  idleTimeout,
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
	if session.IdleTimeout > 0 && time.Since(session.LastActivity) > session.IdleTimeout {
		delete(sessions, cookie)
		return nil
	}
	return session
}

func TouchSession(cookie string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	if session, ok := sessions[cookie]; ok {
		session.LastActivity = time.Now()
	}
}

func DeleteSession(cookie string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	delete(sessions, cookie)
}

func _prune() {
	now := time.Now()
	for cookie, session := range sessions {
		if now.After(session.Expiry) || (session.IdleTimeout > 0 && now.Sub(session.LastActivity) > session.IdleTimeout) {
			delete(sessions, cookie)
		}
	}
}
