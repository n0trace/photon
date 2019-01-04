package photon

import (
	"net/http"
	"sync"
	"time"

	"github.com/n0trace/photon/common"
)

var (
	parallelNumber = 64
	nilHandler     = func(*Context) {
		return
	}

	limitFunc = func() <-chan interface{} {
		var ticker = time.NewTicker(time.Microsecond)
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

	callBackRWMutex   = &sync.RWMutex{}
	middlewareRWMutex = &sync.RWMutex{}
)

type (
	Photon struct {
		httpClientPool sync.Pool
		contextPool    sync.Pool
		middleware     HandlerFunc
		respCallBack   HandlerFunc
		errCallBack    HandlerFunc
		wait           sync.WaitGroup
		filterFunc     func(*Context) bool
		limitFunc      func() <-chan interface{}
		limitChan      <-chan interface{}
		parallelChan   chan interface{}
	}

	HandlerFunc func(*Context)

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

func (p *Photon) ApplyFilterFunc(f func(*Context) bool) {
	p.filterFunc = f
}

func (p *Photon) ApplyLimitFunc(f func() <-chan interface{}) {
	p.limitFunc = f
	p.limitChan = p.limitFunc()
}

func (p *Photon) Use(middlewares ...MiddlewareFunc) {
	if p.middleware == nil {
		middlewareRWMutex.Lock()
		p.middleware = nilHandler
		middlewareRWMutex.Unlock()
	}
	h := p.middleware
	// Chain middleware
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	middlewareRWMutex.Lock()
	p.middleware = h
	middlewareRWMutex.Unlock()
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

	if p.filterFunc == nil {
		p.filterFunc = filterFunc
	}

	if p.limitFunc == nil {
		p.limitFunc = limitFunc
	}

	p.ApplyLimitFunc(p.limitFunc)

	p.parallelChan = make(chan interface{}, parallelNumber)
	p.middleware = nilHandler
	return p
}

func (p *Photon) OnResponse(cb HandlerFunc) {
	callBackRWMutex.Lock()
	p.respCallBack = cb
	callBackRWMutex.Unlock()
}

func (p *Photon) OnError(cb HandlerFunc) {
	callBackRWMutex.Lock()
	p.errCallBack = cb
	callBackRWMutex.Unlock()
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

	p.middleware(ctx)

	p.wait.Add(1)
	p.parallelChan <- true
	go func(ctx *Context) {
		defer p.wait.Done()
		defer p.httpClientPool.Put(ctx.Client)
		defer p.contextPool.Put(ctx)
		defer func() {
			<-p.parallelChan
		}()
		<-p.limitChan
		p.process(ctx)
	}(ctx)
}

func (p *Photon) process(ctx *Context) {

	var err error
	client := ctx.Client
	req := ctx.Request()
	var resp = new(Response)

	ctx.Stage = StageDownloadBefore
	p.middleware(ctx)

	resp.Response, err = client.Do(req)
	ctx.Stage = StageDownloadAfter
	ctx.Response = resp
	ctx.SetError(err)

	p.middleware(ctx)

	if ctx.Error() != nil {
		callBackRWMutex.RLock()
		p.errCallBack(ctx)
		callBackRWMutex.RUnlock()
		return
	}

	callBackRWMutex.RLock()
	p.respCallBack(ctx)
	callBackRWMutex.RUnlock()
}
