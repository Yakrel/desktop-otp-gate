package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
	"log"
	"simple-nginx-otp/utils/config"
	"simple-nginx-otp/utils/ratelimits"
	"simple-nginx-otp/utils/sessions"
	yubikey2 "simple-nginx-otp/utils/yubikey"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		decode, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAABlBMVEVWAAAErtouO1BUAAAAAXRSTlMAQObYZgAAAC1JREFUCNdjYD7AwGPAYMDDoHiEwS2JwUWJoUWRof8jFPl/AiEFFpACoDLmAwAcwAw1QCe40wAAAABJRU5ErkJggg==")
		_ = c.Type("png").Send(decode)
		return nil
	})
	app.Get("*", func(c *fiber.Ctx) error {
		ip := c.IP()
		if len(c.IPs()) > 0 {
			ip = c.IPs()[0]
		}

		cookie := c.Cookies(conf.CookieName)
		session := sessions.GetSession(cookie)

		if session != nil && session.Authorized {
			c.Status(200).Type("txt", "UTF-8")
			return nil
		}

		if session == nil {
			var cookie *fiber.Cookie
			var err error
			session, cookie, err = sessions.NewSession(conf)
			if err != nil {
				return fmt.Errorf("`%s` session creation failed\n%w", ip, err)
			}
			c.Cookie(cookie)
			log.Printf("`%s` new request", ip)
		}

		otp := c.Query("otp")
		if otp != "" {
			if ratelimits.IsLimited(conf, ip) {
				c.Status(429).Type("txt", "UTF-8")
				return nil
			}

			if (len(otp) == 6 && conf.Secret != "" && totp.Validate(otp, conf.Secret)) || (len(otp) >= 6 && conf.YubiOTP != "" && yubikey2.Validate(otp, conf.YubiOTP)) {
				session.Authorized = true
				log.Printf("`%s` successfully logged in, redirecting to `%s`", ip, session.Redirect)
				_ = c.Redirect(session.Redirect)
				return nil
			}
			log.Printf("`%s` sent invalid OTP", ip)
		}

		redirect := c.Get("X-Original-URI", "")
		if redirect != "" && redirect != c.BaseURL()+c.OriginalURL() {
			buffer := make([]byte, len(redirect))
			copy(buffer, redirect)
			session.Redirect = string(buffer)
			log.Printf("`%s` is attempting to access `%s`", ip, session.Redirect)
		}

		_ = c.Status(401).Type("html", "UTF-8").Send(conf.HTML)
		return nil
	})

	log.Printf("listening on http://%s:%d", conf.IP, conf.Port)
	_ = app.Listen(fmt.Sprintf("%s:%d", conf.IP, conf.Port))
}
