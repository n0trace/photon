package photon

import (
	"net/http"
	"sync"
	"time"

	"github.com/n0trace/photon/common"
)

var (
	parallelNumber = 64
	nilHandler     = func(*Context) error {
		return nil
	}

	ticker    = time.NewTicker(time.Millisecond)
	limitFunc = func() <-chan interface{} {
		out := make(chan interface{})
		go func() {
			for t := range ticker.C {
				out <- t
			}
		}()
		return out
	}

	filterMap  = make(map[string]bool)
	filterFunc = func(r *http.Request) bool {
		_, ok := filterMap[r.URL.String()]
		return ok
	}
	newHTTPClientFunc = func() interface{} {
		return http.DefaultClient
	}
	newContextFunc = func() *Context {
		return &Context{
			Client: newHTTPClientFunc().(*http.Client),
			Stage:  StageStructure,
		}
	}
)

type (
	Photon struct {
		httpClientPool  sync.Pool
		contextPool     sync.Pool
		middleware      HandlerFunc
		callBackFuncMap map[ActionType][]HandlerFunc
		wait            sync.WaitGroup
		filterFunc      func(*http.Request) bool
		limitFunc       func() <-chan interface{}
		parallelChan    chan interface{}
	}

	HandlerFunc func(*Context) error

	MiddlewareFunc func(HandlerFunc) HandlerFunc

	PhotonOptionFunc func(*Photon)
)

func WithParallel(parallel int) PhotonOptionFunc {
	return func(p *Photon) {
		parallelNumber = parallel
	}
}

func WithFilter(f func(*http.Request) bool) PhotonOptionFunc {
	return func(p *Photon) {
		p.filterFunc = f
	}
}

func WithLimitFunc(f func() <-chan interface{}) PhotonOptionFunc {
	return func(p *Photon) {
		p.limitFunc = f
	}
}

func (p *Photon) SetFilterFunc(f func(*http.Request) bool) {
	p.filterFunc = f
}

func (p *Photon) Use(middlewares ...MiddlewareFunc) {
	if p.middleware == nil {
		p.middleware = nilHandler
	}
	h := p.middleware
	// Chain middleware
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	p.middleware = h
}

func (p *Photon) Wait() {
	p.wait.Wait()
}

func New(options ...PhotonOptionFunc) (p *Photon) {
	p = new(Photon)

	for _, option := range options {
		option(p)
	}

	p.httpClientPool.New = func() interface{} {
		return newHTTPClientFunc()
	}

	p.contextPool.New = func() interface{} {
		return newContextFunc()
	}

	p.callBackFuncMap = make(map[ActionType][]HandlerFunc)
	p.callBackFuncMap[OnRequest] = []HandlerFunc{}
	p.callBackFuncMap[OnResponse] = []HandlerFunc{}
	p.callBackFuncMap[OnError] = []HandlerFunc{}

	if p.filterFunc == nil {
		p.filterFunc = filterFunc
	}

	if p.limitFunc == nil {
		p.limitFunc = limitFunc
	}

	p.parallelChan = make(chan interface{}, parallelNumber)
	p.middleware = nilHandler
	return p
}

func (p *Photon) On(action ActionType, CallBackFuncFunc HandlerFunc) {
	before, ok := p.callBackFuncMap[action]
	if ok {
		p.callBackFuncMap[action] = append(before, CallBackFuncFunc)
		return
	}
	p.callBackFuncMap[action] = []HandlerFunc{CallBackFuncFunc}
}

func (p *Photon) Visit(url string, visitOptions ...VisitOptionFunc) {
	p.VisitRequest(common.Must2(http.NewRequest("GET", url, nil)).(*http.Request), visitOptions...)
}

func (p *Photon) VisitRequest(req *http.Request, visitOptions ...VisitOptionFunc) {
	done := p.filterFunc(req)
	if done {
		return
	}
	var visitOption = new(VisitOption)
	for _, option := range visitOptions {
		option(visitOption)
	}

	ctx := p.contextPool.Get().(*Context)
	ctx.Reset()
	if visitOption.PreContext != nil {
		ctx.Client = visitOption.PreContext.Client
	}

	ctx.SetMeta(visitOption.Meta)
	if visitOption.Client != nil {
		ctx.Client = visitOption.Client
	}

	ctx.SetRequest(req)

	common.Must(p.middleware(ctx))

	p.wait.Add(1)

	p.parallelChan <- true
	go p.process(ctx)
}

func (p *Photon) process(ctx *Context) {
	<-p.limitFunc()
	defer p.wait.Done()
	defer func() {
		p.httpClientPool.Put(ctx.Client)
		p.contextPool.Put(ctx)
		<-p.parallelChan

	}()
	var err error
	client := ctx.Client
	req := ctx.Request()
	var resp = new(Response)
	OnRequestCBS := p.callBackFuncMap[OnRequest]
	for _, cb := range OnRequestCBS {
		cb(ctx)
	}
	common.Must(p.middleware(ctx))
	ctx.Stage = StageDownloadBefore
	resp.Response, err = client.Do(req)
	ctx.Stage = StageDownloadAfter
	ctx.Response = resp
	ctx.SetError(err)
	common.Must(p.middleware(ctx))

	if ctx.Error() != nil {
		OnErrorCBS := p.callBackFuncMap[OnError]
		for _, cb := range OnErrorCBS {
			cb(ctx)
		}
		return
	}

	OnResponseCBS := p.callBackFuncMap[OnResponse]
	for _, cb := range OnResponseCBS {
		cb(ctx)
	}
}
