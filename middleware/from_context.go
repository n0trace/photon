package middleware

import "github.com/n0trace/photon"

func FromContext(lastCtx photon.Context) photon.MiddlewareFunc {
	return func(next photon.HandlerFunc) photon.HandlerFunc {
		return func(ctx photon.Context) {
			defer next(ctx)
			if ctx.Downloaded() {
				return
			}
			ctx.FromContext(lastCtx)
		}
	}
}
