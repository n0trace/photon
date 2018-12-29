package photon

import "net/http"

type (
	VisitOption struct {
		Meta       map[string]interface{}
		Client     *http.Client
		DontFilter bool
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

func VisitWithDontFiter() VisitOptionFunc {
	return func(option *VisitOption) {
		option.DontFilter = true
	}
}
