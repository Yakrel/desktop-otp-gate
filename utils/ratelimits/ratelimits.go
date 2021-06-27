package ratelimits

import (
	"log"
	"simple-nginx-otp/utils/config"
	"sync"
	"time"
)

type RateLimit struct {
	Count  int8
	Expiry time.Time
}

var rateLimits = make(map[string]*RateLimit)
var rateLimitsMutex = sync.Mutex{}

func IsLimited(conf *config.Config, ip string) bool {
	rateLimitsMutex.Lock()
	defer rateLimitsMutex.Unlock()
	_prune()
	rateLimit, ok := rateLimits[ip]
	if !ok {
		rateLimit = &RateLimit{
			Count: 1,
		}
	} else {
		rateLimit.Count++
	}
	rateLimit.Expiry = time.Now().Add(time.Duration(conf.RateLimitExpiry) * time.Minute)
	rateLimits[ip] = rateLimit
	if rateLimit.Count == conf.RateLimitCount {
		log.Printf("rate limited ip `%s`", ip)
	}
	return rateLimit.Count > conf.RateLimitCount
}

func _prune() {
	for ip, rateLimit := range rateLimits {
		if time.Now().After(rateLimit.Expiry) {
			delete(rateLimits, ip)
		}
	}
}
