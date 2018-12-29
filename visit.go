package photon

import "net/http"

type (
	VisitOption struct {
		Meta       map[string]interface{}
		Client     *http.Client
		Filter     bool
		PreContext *Context
	}

	VisitOptionFunc func(*VisitOption)
)

func VisitWithMeta(meta map[string]interface{}) VisitOptionFunc {
	return func(option *VisitOption) {
		option.Meta = meta
	}
}

func VisitWithClient(client *http.Client) VisitOptionFunc {
	return func(option *VisitOption) {
		option.Client = client
	}
}

func VisitWithFiter(filter bool) VisitOptionFunc {
	return func(option *VisitOption) {
		option.Filter = filter
	}
}

func VisitWithPreContext(ctx *Context) VisitOptionFunc {
	return func(option *VisitOption) {
		option.PreContext = ctx
	}
}
