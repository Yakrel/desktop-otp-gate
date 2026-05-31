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
		if rateLimit.Count > conf.RateLimitCount {
			rateLimit.Expiry = time.Now().Add(time.Duration(conf.RateLimitExpiry) * time.Minute)
			return true
		}
		rateLimit.Count++
	}
	rateLimit.Expiry = time.Now().Add(time.Duration(conf.RateLimitExpiry) * time.Minute)
	rateLimits[ip] = rateLimit
	if rateLimit.Count == conf.RateLimitCount {
		log.Printf("`%s` has been rate limited", ip)
	}
	return rateLimit.Count > conf.RateLimitCount
}

type Status struct {
	IsLimited bool
	Remaining int
	LockTime  int
}

func GetStatus(conf *config.Config, ip string) Status {
	rateLimitsMutex.Lock()
	defer rateLimitsMutex.Unlock()
	_prune()
	rateLimit, ok := rateLimits[ip]
	if !ok {
		return Status{
			IsLimited: false,
			Remaining: int(conf.RateLimitCount),
			LockTime:  0,
		}
	}

	remaining := int(conf.RateLimitCount) - int(rateLimit.Count)
	if remaining < 0 {
		remaining = 0
	}

	isLimited := rateLimit.Count > conf.RateLimitCount
	lockTime := 0
	if isLimited {
		lockTime = int(time.Until(rateLimit.Expiry).Seconds())
		if lockTime < 0 {
			lockTime = 0
		}
	}

	return Status{
		IsLimited: isLimited,
		Remaining: remaining,
		LockTime:  lockTime,
	}
}

func _prune() {
	for ip, rateLimit := range rateLimits {
		if time.Now().After(rateLimit.Expiry) {
			delete(rateLimits, ip)
		}
	}
}
