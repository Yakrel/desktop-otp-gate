package config

import (
	"fmt"
	"github.com/pquerna/otp/totp"
	"os"
	"strconv"
)

type Config struct {
	IP              string
	Port            int
	Secret          string
	YubiOTP         string
	HTML            []byte
	CookieName      string
	CookieLength    int8
	CookieLifetime  int16
	CookieDomain    string
	RateLimitCount  int8
	RateLimitExpiry int16
}

var config *Config

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	ip := _getEnv("SNO_LISTEN_IP", "0.0.0.0")
	port, err := strconv.Atoi(_getEnv("SNO_LISTEN_PORT", "7079"))
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_LISTEN_PORT\n%w", err)
	}

	secret := _getEnv("SNO_SECRET", "")
	yubiotp := _getEnv("SNO_YUBIOTP", "")
	if secret == "" && yubiotp == "" {
		key, _ := totp.Generate(totp.GenerateOpts{
			Issuer:      "sno",
			AccountName: "sno",
		})
		return nil, fmt.Errorf("SNO_SECRET and SNO_YUBIOTP missing, here's a random SNO_SECRET:\n%s", key.Secret())
	}
	if len(yubiotp) > 12 {
		yubiotp = yubiotp[:12]
	}

	title := _getEnv("SNO_TITLE", "Simple Nginx OTP")
	var html = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>` + title + `</title><style>body{height:100vh;display:flex;justify-content:center;align-items:center}</style></head><body> <input id="auth" type="text" autofocus/> <button onclick="post()">Submit</button> <script>let auth=document.getElementById('auth');function post(){window.location.href="?otp="+auth.value;};auth.addEventListener("keyup",function(event){if(event.keyCode===13){post();}});</script> </body></html>`

	cookieName := _getEnv("SNO_COOKIE_NAME", "sno_session")

	cookieLength, err := strconv.ParseInt(_getEnv("SNO_COOKIE_LENGTH", "16"), 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_COOKIE_LENGTH\n%w", err)
	}

	cookieLifetime, err := strconv.ParseInt(_getEnv("SNO_COOKIE_LIFETIME", "14"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_COOKIE_LIFETIME\n%w", err)
	}

	cookieDomain := _getEnv("SNO_COOKIE_DOMAIN", "")

	rateLimitCount, err := strconv.ParseInt(_getEnv("SNO_RATE_LIMIT_COUNT", "3"), 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_RATE_LIMIT_COUNT\n%w", err)
	}

	rateLimitExpiry, err := strconv.ParseInt(_getEnv("SNO_RATE_LIMIT_LIFETIME", "1"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_RATE_LIMIT_COUNT\n%w", err)
	}

	config = &Config{
		IP:              ip,
		Port:            port,
		Secret:          secret,
		YubiOTP:         yubiotp,
		HTML:            []byte(html),
		CookieName:      cookieName,
		CookieLength:    int8(cookieLength),
		CookieLifetime:  int16(cookieLifetime),
		CookieDomain:    cookieDomain,
		RateLimitCount:  int8(rateLimitCount),
		RateLimitExpiry: int16(rateLimitExpiry),
	}
	return config, nil
}

func _getEnv(env string, def string) string {
	val := os.Getenv(env)
	if val == "" {
		return def
	}
	return val
}
