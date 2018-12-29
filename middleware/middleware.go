package middleware

import (
	"github.com/n0trace/photon"
)

type (
	Holder func(*photon.Context) bool
)

// DefaultSkipper returns false which processes the middleware.
func DefaultHolder(*photon.Context) bool {
	return true
}

func DownloadAfterHolder(ctx *photon.Context) bool {
	if ctx.Stage == photon.StageDownloadAfter {
		return true
	}
	return false
}

func DownloadBeforeHolder(ctx *photon.Context) bool {
	if ctx.Stage == photon.StageDownloadBefore {
		return true
	}
	return false
}
