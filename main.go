package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/pquerna/otp/totp"
	"log"
	"net/http"
	"simple-nginx-otp/utils/config"
	"simple-nginx-otp/utils/ratelimits"
	"simple-nginx-otp/utils/sessions"
	"simple-nginx-otp/utils/yubikey"
	"strings"
)

var lastURL = make(map[string]string)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		decode, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEVWAAAErtouO1BUAAAAAXRSTlMAQObYZgAAAC1JREFUCNdjYD7AwGPAYMDDoHiEwS2JwUWJoUWRof8jFPl/AiEFFpACoDLmAwAcwAw1QCe40wAAAABJRU5ErkJggg==")
		w.Header().Set("Content-Type", "image/png")
		w.Write(decode)
	})
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-Ip")
		if ip == "" && r.Header.Get("X-Forwarded-For") != "" {
			split := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
			ip = split[0]
		}
		if ip == "" {
			ip = r.RemoteAddr
		}
		ip = strings.Split(ip, ":")[0]

		var session *sessions.Session
		cookie, err := r.Cookie(conf.CookieName)
		if err == nil && cookie != nil {
			session = sessions.GetSession(cookie.Value)
		}

		// already authorized, send 200
		if session != nil && session.Authorized {
			w.WriteHeader(200)
			return
		}

		// auth_request coming from nginx with X-Original-URI header
		_, exist := lastURL[ip]
		if !exist {
			lastURL[ip] = "/"
		}
		redirect := r.Header.Get("X-Original-URI")
		if redirect != "" {
			if redirect != r.URL.RequestURI() {
				buffer := make([]byte, len(redirect))
				copy(buffer, redirect)
				lastURL[ip] = string(buffer)
				w.WriteHeader(401)
				return
			}
		}

		// user is redirected to SNO
		if session == nil {
			var cookie *http.Cookie
			var err error
			session, cookie, err = sessions.NewSession(conf)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			session.Redirect = lastURL[ip]
			delete(lastURL, ip)
			log.Printf("`%s` is attempting to access `%s`", ip, session.Redirect)
			w.Header().Set("Set-Cookie", cookie.String())
		}

		// check otp query param
		otp := r.URL.Query().Get("otp")
		if otp != "" {
			log.Printf("`%s` attempted authentication", ip)
			if !ratelimits.IsLimited(conf, ip) {
				if (len(otp) == 6 && conf.Secret != "" && totp.Validate(otp, conf.Secret)) || (len(otp) >= 6 && conf.YubiOTP != "" && yubikey.Validate(otp, conf.YubiOTP)) {
					session.Authorized = true
					log.Printf("`%s` successfully logged in, redirecting to `%s`", ip, session.Redirect)
					http.Redirect(w, r, session.Redirect, 302)
					return
				}
			}
		}

		// return form
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(401)
		w.Write(conf.HTML)
		return
	})

	log.Printf("listening on http://%s:%d", conf.IP, conf.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", conf.IP, conf.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
