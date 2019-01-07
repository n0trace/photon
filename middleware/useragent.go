package middleware

import (
	"math/rand"
	"time"

	"github.com/n0trace/photon"
)

var defaultUserAgentList = []string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36 OPR/38.0.2220.41",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0)",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
}

type UserAgentConfig struct {
	UserAgent []string
}

func UserAgent(useragents ...string) photon.MiddlewareFunc {
	if len(useragents) == 0 {
		useragents = defaultUserAgentList
	}
	config := UserAgentConfig{
		UserAgent: useragents,
	}
	return UserAgentWithConfig(config)
}

func UserAgentWithConfig(config UserAgentConfig) photon.MiddlewareFunc {
	rand.Seed(time.Now().Unix())
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || ctx.Downloaded() {
				return
			}
			var idx int
			if len(config.UserAgent) > 1 {
				idx = rand.Intn(len(config.UserAgent) - 1)
			}
			ctx.Request().Header.Set("User-Agent", config.UserAgent[idx])
		}
	}
}
