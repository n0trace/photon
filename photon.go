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

	ticker    = time.NewTicker(time.Microsecond)
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
	filterFunc = func(ctx *Context) bool {
		_, ok := filterMap[ctx.Request().URL.String()]
		filterMap[ctx.Request().URL.String()] = true
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
	callBackMutex   = &sync.RWMutex{}
	middlewareMutex = &sync.RWMutex{}
)

type (
	Photon struct {
		httpClientPool  sync.Pool
		contextPool     sync.Pool
		middleware      HandlerFunc
		callBackFuncMap map[ActionType][]HandlerFunc
		wait            sync.WaitGroup
		filterFunc      func(*Context) bool
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

func WithFilterFunc(f func(*Context) bool) PhotonOptionFunc {
	return func(p *Photon) {
		p.filterFunc = f
	}
}

func WithLimitFunc(f func() <-chan interface{}) PhotonOptionFunc {
	return func(p *Photon) {
		p.limitFunc = f
	}
}

func (p *Photon) SetFilterFunc(f func(*Context) bool) {
	p.filterFunc = f
}

func (p *Photon) Use(middlewares ...MiddlewareFunc) {
	if p.middleware == nil {
		middlewareMutex.Lock()
		p.middleware = nilHandler
		middlewareMutex.Unlock()
	}
	h := p.middleware
	// Chain middleware
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	middlewareMutex.Lock()
	p.middleware = h
	middlewareMutex.Unlock()
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
	callBackMutex.RLock()
	before, ok := p.callBackFuncMap[action]
	callBackMutex.RUnlock()
	if ok {
		callBackMutex.Lock()
		p.callBackFuncMap[action] = append(before, CallBackFuncFunc)
		callBackMutex.Unlock()
		return
	}
	callBackMutex.Lock()
	p.callBackFuncMap[action] = []HandlerFunc{CallBackFuncFunc}
	callBackMutex.Unlock()
}

func (p *Photon) Visit(url string, visitOptions ...VisitOptionFunc) {
	p.VisitRequest(common.Must2(http.NewRequest("GET", url, nil)).(*http.Request), visitOptions...)
}

func (p *Photon) VisitRequest(req *http.Request, visitOptions ...VisitOptionFunc) {
	var visitOption = new(VisitOption)
	for _, option := range visitOptions {
		option(visitOption)
	}

	ctx := p.contextPool.Get().(*Context)
	ctx.Reset()

	ctx.SetMeta(visitOption.Meta)
	if visitOption.Client != nil {
		ctx.Client = visitOption.Client
	}

	ctx.SetRequest(req)

	done := p.filterFunc(ctx)
	if !visitOption.DontFilter && done {
		return
	}

	middlewareMutex.RLock()
	middleware := p.middleware
	middlewareMutex.RUnlock()

	common.Must(middleware(ctx))

	p.wait.Add(1)
	p.parallelChan <- true
	go func(ctx *Context) {
		defer p.wait.Done()
		defer p.httpClientPool.Put(ctx.Client)
		defer p.contextPool.Put(ctx)
		defer func() {
			<-p.parallelChan
		}()
		<-p.limitFunc()
		p.process(ctx)
	}(ctx)
}

func (p *Photon) process(ctx *Context) {

	var err error
	client := ctx.Client
	req := ctx.Request()
	var resp = new(Response)

	callBackMutex.RLock()
	callBackMap := p.callBackFuncMap
	OnRequestCBS := callBackMap[OnRequest]
	OnResponseCBS := callBackMap[OnResponse]
	OnErrorCBS := callBackMap[OnError]
	callBackMutex.RUnlock()
	for _, cb := range OnRequestCBS {
		_ = cb(ctx)
	}
	middlewareMutex.RLock()
	middleware := p.middleware
	middlewareMutex.RUnlock()

	common.Must(middleware(ctx))
	ctx.Stage = StageDownloadBefore
	resp.Response, err = client.Do(req)
	ctx.Stage = StageDownloadAfter
	ctx.Response = resp
	ctx.SetError(err)
	common.Must(middleware(ctx))

	if ctx.Error() != nil {
		for _, cb := range OnErrorCBS {
			_ = cb(ctx)
		}
		return
	}

	for _, cb := range OnResponseCBS {
		_ = cb(ctx)
	}
}
