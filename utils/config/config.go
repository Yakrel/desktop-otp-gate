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
	CookieMinutes   int32
	SessionCookie   bool
	CookieSecure    bool
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
	var html = buildHTML(title)

	cookieName := _getEnv("SNO_COOKIE_NAME", "sno_session")

	cookieLength, err := strconv.ParseInt(_getEnv("SNO_COOKIE_LENGTH", "16"), 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_COOKIE_LENGTH\n%w", err)
	}

	cookieLifetime, err := strconv.ParseInt(_getEnv("SNO_COOKIE_LIFETIME", "14"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_COOKIE_LIFETIME\n%w", err)
	}

	cookieMinutes, err := strconv.ParseInt(_getEnv("SNO_COOKIE_LIFETIME_MINUTES", "0"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_COOKIE_LIFETIME_MINUTES\n%w", err)
	}

	sessionCookie, err := strconv.ParseBool(_getEnv("SNO_SESSION_COOKIE", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_SESSION_COOKIE\n%w", err)
	}

	cookieSecure, err := strconv.ParseBool(_getEnv("SNO_COOKIE_SECURE", "true"))
	if err != nil {
		return nil, fmt.Errorf("invalid SNO_COOKIE_SECURE\n%w", err)
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
		CookieMinutes:   int32(cookieMinutes),
		SessionCookie:   sessionCookie,
		CookieSecure:    cookieSecure,
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

func buildHTML(title string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>` + title + `</title>
<style>
:root{color-scheme:dark}
*{box-sizing:border-box}
body{min-height:100vh;margin:0;display:grid;place-items:center;background:#080b12;color:#f8fafc;font-family:Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif}
body:before{content:"";position:fixed;inset:0;background:radial-gradient(circle at 50% 0%,rgba(14,165,233,.22),transparent 32rem),linear-gradient(145deg,#070a11 0%,#111827 52%,#05070b 100%);z-index:-1}
main{width:min(92vw,390px);padding:32px;border:1px solid rgba(148,163,184,.22);background:rgba(15,23,42,.82);box-shadow:0 28px 70px rgba(0,0,0,.48);backdrop-filter:blur(18px);border-radius:18px}
.mark{width:44px;height:44px;border-radius:12px;display:grid;place-items:center;margin-bottom:22px;background:#0ea5e9;color:#00111f;font-weight:800;font-size:22px}
h1{margin:0 0 8px;font-size:24px;line-height:1.15;font-weight:720;letter-spacing:0}
p{margin:0 0 24px;color:#94a3b8;line-height:1.5;font-size:14px}
label{display:block;margin-bottom:8px;color:#cbd5e1;font-size:13px;font-weight:650}
.row{display:flex;gap:10px}
input{width:100%;height:48px;border:1px solid rgba(148,163,184,.28);background:#020617;color:#f8fafc;border-radius:10px;padding:0 14px;font-size:20px;letter-spacing:3px;text-align:center;outline:none}
input:focus{border-color:#38bdf8;box-shadow:0 0 0 4px rgba(56,189,248,.15)}
button{height:48px;border:0;border-radius:10px;padding:0 18px;background:#38bdf8;color:#00111f;font-size:14px;font-weight:750;cursor:pointer}
button:hover{background:#7dd3fc}
.hint{margin-top:16px;margin-bottom:0;color:#64748b;font-size:12px}
</style>
</head>
<body>
<main>
<div class="mark">#</div>
<h1>` + title + `</h1>
<p>Enter the 6-digit code from your authenticator app.</p>
<label for="auth">Authenticator code</label>
<div class="row">
<input id="auth" name="otp" inputmode="numeric" autocomplete="one-time-code" pattern="[0-9]*" maxlength="6" autofocus>
<button onclick="post()">Verify</button>
</div>
<p class="hint">Access is rate limited after repeated failed attempts.</p>
</main>
<script>
const auth=document.getElementById('auth');
function post(){const value=auth.value.replace(/\s+/g,'');if(value){window.location.href='?otp='+encodeURIComponent(value);}}
auth.addEventListener('input',()=>{auth.value=auth.value.replace(/\D/g,'').slice(0,6);if(auth.value.length===6)post();});
auth.addEventListener('keyup',event=>{if(event.key==='Enter')post();});
</script>
</body>
</html>`
}
