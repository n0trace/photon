package middleware

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/common"
)

type RandomUAConfig struct {
	Holder           Holder
	Browser          []string
	UserAgentJSONURL string
}

var (
	DefaultRandomUAConfig = RandomUAConfig{
		Holder: DefaultHolder,
	}

	browserUserAgentMap   = make(map[string][]string)
	browserUserAgentSlice = []string{}
	cacheUserAgentOnce    sync.Once
	userAgentJSONURL      = "https://user-agent.now.sh/useragent.json"
)

func RandomUAWithConfig(config RandomUAConfig) photon.MiddlewareFunc {
	if config.Holder == nil {
		config.Holder = DownloadBeforeHolder
	}
	var url = userAgentJSONURL
	if config.UserAgentJSONURL != "" {
		url = config.UserAgentJSONURL
	}
	cacheUserAgentOnce.Do(func() { common.Must(cacheUserAgent(url)) })
	rand.Seed(time.Now().Unix())
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx *photon.Context) error {
			if !config.Holder(ctx) {
				return next(ctx)
			}
			if len(config.Browser) == 0 {
				idx := rand.Intn(len(browserUserAgentSlice) - 1)
				ctx.Request().Header.Set("User-Agent", browserUserAgentSlice[idx])
			} else {
				var browserIdx, uaIdx int
				if len(config.Browser) > 1 {
					browserIdx = rand.Intn(len(config.Browser) - 1)
				}
				browser := config.Browser[browserIdx]
				uaSlice := browserUserAgentMap[browser]
				if len(uaSlice) > 1 {
					uaIdx = rand.Intn(len(uaSlice) - 1)
				}
				ctx.Request().Header.Set("User-Agent", uaSlice[uaIdx])
			}
			return next(ctx)
		}
	}
}

func RandomUserAgent(browsers ...string) photon.MiddlewareFunc {
	config := DefaultRandomUAConfig
	config.Browser = browsers
	return RandomUAWithConfig(config)
}

func cacheUserAgent(url string) (err error) {
	resp, err := http.Get(url)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.Unmarshal(bodyBytes, &browserUserAgentMap)
	if err != nil {
		return
	}
	for _, value := range browserUserAgentMap {
		browserUserAgentSlice = append(browserUserAgentSlice, value...)
	}
	for browser, ualist := range browserUserAgentMap {
		if len(ualist) == 0 {
			delete(browserUserAgentMap, browser)
		}
	}
	if len(browserUserAgentSlice) == 0 {
		panic(errors.New("user-agent nil"))
	}
	return
}
