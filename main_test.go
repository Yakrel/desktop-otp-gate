package main

import (
	"net/http"
	"net/http/httptest"
	"simple-nginx-otp/utils/config"
	"simple-nginx-otp/utils/sessions"
	"testing"
	"time"
)

func testConfig() *config.Config {
	return &config.Config{
		IP:              "127.0.0.1",
		Port:            7079,
		Secret:          "JBSWY3DPEHPK3PXP",
		HTML:            []byte("<html><body>OTP</body></html>"),
		ActiveHTML:      []byte("<html><body>Session Active <a href=\"/sno/logout\">Log Out</a></body></html>"),
		CookieName:      "sno_session",
		CookieLength:    16,
		CookieLifetime:  14,
		CookieMinutes:   60,
		IdleMinutes:     15,
		SessionCookie:   true,
		CookieSecure:    false,
		RateLimitCount:  3,
		RateLimitExpiry: 1,
	}
}

func TestLogoutEndpoint(t *testing.T) {
	conf := testConfig()
	router := setupRouter(conf)

	session, cookie, err := sessions.NewSession(conf)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	session.Authorized = true

	req := httptest.NewRequest("GET", "/sno/logout", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusFound {
		t.Fatalf("expected status 302, got %d", res.StatusCode)
	}

	loc := res.Header.Get("Location")
	if loc != "/sno" {
		t.Fatalf("expected redirect to /sno, got %s", loc)
	}

	// Verify session deleted
	if sessions.GetSession(cookie.Value) != nil {
		t.Fatal("expected session to be deleted after logout")
	}

	// Verify expired cookie returned
	var expiredCookieFound bool
	for _, c := range res.Cookies() {
		if c.Name == conf.CookieName && c.MaxAge < 0 {
			expiredCookieFound = true
			break
		}
	}
	if !expiredCookieFound {
		t.Fatal("expected expired cookie header in response")
	}
}

func TestAuthorizedSessionAndTouch(t *testing.T) {
	conf := testConfig()
	router := setupRouter(conf)

	session, cookie, err := sessions.NewSession(conf)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	session.Authorized = true
	oldActivity := time.Now().Add(-5 * time.Minute)
	session.LastActivity = oldActivity

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	updatedSession := sessions.GetSession(cookie.Value)
	if updatedSession == nil {
		t.Fatal("expected session to still exist")
	}
	if !updatedSession.LastActivity.After(oldActivity) {
		t.Fatal("expected LastActivity to be refreshed via TouchSession")
	}
}

func TestIdleExpiredSessionReturns401(t *testing.T) {
	conf := testConfig()
	router := setupRouter(conf)

	session, cookie, err := sessions.NewSession(conf)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	session.Authorized = true
	// Simulate idle time exceeding IdleMinutes (15 min)
	session.LastActivity = time.Now().Add(-20 * time.Minute)

	req := httptest.NewRequest("GET", "/sno", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for idle expired session, got %d", res.StatusCode)
	}

	if sessions.GetSession(cookie.Value) != nil {
		t.Fatal("expected idle expired session to be removed")
	}
}

func TestAuthorizedDirectSnoVisitShowsActiveUI(t *testing.T) {
	conf := testConfig()
	router := setupRouter(conf)

	session, cookie, err := sessions.NewSession(conf)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	session.Authorized = true

	// Direct visit to /sno
	req := httptest.NewRequest("GET", "/sno", nil)
	req.AddCookie(cookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	body := w.Body.String()
	if body != string(conf.ActiveHTML) {
		t.Fatalf("expected ActiveHTML body, got %q", body)
	}

	// Background auth_request check from nginx (X-Original-URI != requestURL)
	reqAuth := httptest.NewRequest("GET", "/sno", nil)
	reqAuth.Header.Set("X-Original-URI", "https://code.example.com/editor")
	reqAuth.AddCookie(cookie)
	wAuth := httptest.NewRecorder()

	router.ServeHTTP(wAuth, reqAuth)

	resAuth := wAuth.Result()
	defer resAuth.Body.Close()

	if resAuth.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resAuth.StatusCode)
	}
	if wAuth.Body.Len() != 0 {
		t.Fatalf("expected empty body for background auth_request, got %d bytes", wAuth.Body.Len())
	}
}
