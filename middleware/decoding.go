package middleware

import (
	"compress/flate"
	"compress/gzip"
	"strings"

	"github.com/n0trace/photon"
)

type DecodingConfig struct {
}

func DecodingWithConfig(config DecodingConfig) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Error() != nil || !ctx.Downloaded() {
				return
			}
			resp := ctx.StdResponse()
			switch strings.ToLower(resp.Header.Get("Content-Encoding")) {
			case "gzip":
				gzipReader, err := gzip.NewReader(resp.Body)
				if err != nil {
					ctx.SetError(err)
					return
				}
				resp.Body = gzipReader
				return
			case "deflate":
				flateReader := flate.NewReader(resp.Body)
				resp.Body = flateReader
			default:
			}
		}
	}
}

func Decoding() photon.MiddlewareFunc {
	config := DecodingConfig{}
	return DecodingWithConfig(config)
}
