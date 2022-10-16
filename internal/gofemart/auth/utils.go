package auth

import (
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"time"
)

func (a *Auth) junkCleaner() {
	ticker := time.NewTicker(CookieCheck * time.Second)
	for {
		<-ticker.C
		func() {
			a.Lock.Lock()
			now := time.Now()
			expiredCookies := 0
			for id, cookieInfo := range a.cookiesMeta {
				if now.After(cookieInfo.validTill) {
					logging.Debug("Cookie with id=%s is expired", id)
					delete(a.cookiesMeta, id)
					expiredCookies += 1
				}
			}
			switch {
			case len(a.cookiesMeta) == 0:
				logging.Debug("No Cookies to track")
			case expiredCookies == 0:
				logging.Debug("All tracked Cookies are valid")
			}
			a.Lock.Unlock()
		}()

	}
}
