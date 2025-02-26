package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tclemos/go-expert-rate-limiter/config"
	"github.com/tclemos/go-expert-rate-limiter/internal/cache/redis"
)

const (
	requestsPrefix  = "requests"
	bannedPrefix    = "banned"
	rateLimitErrMsg = "The %s %s has reached the maximum number of requests or actions allowed within a certain time frame, wait until %s"
)

func IPTokenRateLimiter(ctx context.Context, config config.WebServerConfig) func(http.Handler) http.Handler {
	var cache Cache = redis.NewRedisCache(ctx, config)
	var mutex sync.Mutex
	fn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// call redis
			mutex.Lock()
			defer mutex.Unlock()

			// get ip address
			ip := r.RemoteAddr
			ipRequests := getRequests(ctx, cache, ip)

			// get api key header
			apiKey := getApiKey(r)

			// if there is an api key
			if len(apiKey) > 0 {
				// check if api key  is banned
				banned, until := isBanned(ctx, cache, apiKey)
				if banned {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte(fmt.Sprintf(rateLimitErrMsg, "API KEY", apiKey, until.Format(time.RFC3339))))
					return
				}

				// check api key request limit
				apiKeyRequests := getRequests(ctx, cache, apiKey)
				fmt.Println("ip", ip, "ipRequests", ipRequests)
				fmt.Println("apiKey", apiKey, "apiKeyRequests", apiKeyRequests)

				if apiKeyRequests >= uint64(config.MaxRequestsPerSecondPerAPIKey) {
					banRequests(ctx, cache, apiKey, config.BanDuration)
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte(fmt.Sprintf(rateLimitErrMsg, "API KEY", apiKey, time.Now().Add(config.BanDuration).Format(time.RFC3339))))
					return
				} else {
					setRequests(ctx, cache, ip, ipRequests+1)
					setRequests(ctx, cache, apiKey, apiKeyRequests+1)
					next.ServeHTTP(w, r)
					return
				}
			}

			// check if ip is banned
			banned, until := isBanned(ctx, cache, ip)
			if banned {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(fmt.Sprintf(rateLimitErrMsg, "IP", ip, until.Format(time.RFC3339))))
				return
			}

			fmt.Println("ip", ip, "ipRequests", ipRequests)

			// check ip request limit
			if ipRequests >= uint64(config.MaxRequestsPerSecondPerIP) {
				banRequests(ctx, cache, ip, config.BanDuration)
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(fmt.Sprintf(rateLimitErrMsg, "IP", ip, time.Now().Add(config.BanDuration).Format(time.RFC3339))))
				return
			}
			setRequests(ctx, cache, ip, ipRequests+1)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return fn
}

func getApiKey(r *http.Request) string {
	apiKey := r.Header.Get("API_KEY")
	return apiKey
}

func getRequests(ctx context.Context, cache Cache, key string) uint64 {
	requests, err := cache.Get(ctx, fmt.Sprintf("%s:%s", requestsPrefix, key))
	if err != nil {
		return 0
	}

	r, err := strconv.ParseUint(requests, 10, 64)
	if err != nil {
		return 0
	}
	return r
}

func banRequests(ctx context.Context, cache Cache, key string, duration time.Duration) {
	bannedUntilStr := strconv.FormatInt(time.Now().Add(duration).Unix(), 10)
	cache.Set(ctx, fmt.Sprintf("%s:%s", bannedPrefix, key), bannedUntilStr, duration)
}

func isBanned(ctx context.Context, cache Cache, key string) (bool, time.Time) {
	bannedUntilStr, err := cache.Get(ctx, fmt.Sprintf("%s:%s", bannedPrefix, key))
	if err != nil {
		return false, time.Time{}
	}

	bannedUntilUnix, err := strconv.ParseInt(bannedUntilStr, 10, 64)
	if err != nil {
		return false, time.Time{}
	}

	bannedUntil := time.Unix(bannedUntilUnix, 0)

	return true, bannedUntil
}

func setRequests(ctx context.Context, cache Cache, key string, requests uint64) {
	if requests == 1 {
		go func() {
			time.Sleep(time.Second)
			cache.Del(ctx, fmt.Sprintf("%s:%s", requestsPrefix, key))
		}()
	}
	cache.Set(ctx, fmt.Sprintf("%s:%s", requestsPrefix, key), strconv.FormatUint(requests, 10), time.Second)
}
