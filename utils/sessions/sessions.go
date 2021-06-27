package sessions

import (
	"github.com/gofiber/fiber/v2"
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

func NewSession(conf *config.Config) (*Session, *fiber.Cookie, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, err := rand.GenerateRandomString(conf.CookieLength)
	if err != nil {
		return nil, nil, err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = conf.CookieName
	cookie.Value = session
	cookie.Expires = time.Now().Add(time.Hour * time.Duration(24*conf.CookieLifetime))
	if conf.CookieDomain != "" {
		cookie.Domain = conf.CookieDomain
	}

	sessions[session] = &Session{
		Redirect:   "/",
		Authorized: false,
		Expiry:     cookie.Expires,
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
