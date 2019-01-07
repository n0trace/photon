package photon

import (
	"net/http"
	"sync"
	"time"

	"github.com/n0trace/photon/common"
)

var (
	parallelNumber = 64

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
		}
	}

	callBackRWMutex = &sync.RWMutex{}
)

type (
	Photon struct {
		httpClientPool sync.Pool
		contextPool    sync.Pool

		wait         sync.WaitGroup
		filterFunc   func(*Context) bool
		limitFunc    func() <-chan interface{}
		limitChan    <-chan interface{}
		parallelChan chan interface{}
		callback
	}
	callback struct {
		requestCallbacks []HandlerFunc
		respCallbacks    []HandlerFunc
		errCallbacks     []HandlerFunc
		scrapedCallbacks []HandlerFunc
	}

	HandlerFunc      func(*Context)
	MiddlewareFunc   func(*Photon)
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
	for i := 0; i < len(middlewares); i++ {
		middlewares[i](p)
	}
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
	return p
}

func (p *Photon) OnRequest(cb HandlerFunc) {
	callBackRWMutex.Lock()
	p.requestCallbacks = append(p.requestCallbacks, cb)
	callBackRWMutex.Unlock()
}

func (p *Photon) OnResponse(cb HandlerFunc) {
	callBackRWMutex.Lock()
	p.respCallbacks = append(p.respCallbacks, cb)
	callBackRWMutex.Unlock()
}

func (p *Photon) OnError(cb HandlerFunc) {
	callBackRWMutex.Lock()
	p.errCallbacks = append(p.errCallbacks, cb)
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

	p.execHandleFunc(ctx, p.requestCallbacks)
	resp.Response, err = client.Do(req)
	ctx.Response = resp
	ctx.SetError(err)

	if ctx.Error() != nil {
		callBackRWMutex.RLock()
		p.execHandleFunc(ctx, p.errCallbacks)
		callBackRWMutex.RUnlock()
		return
	}

	callBackRWMutex.RLock()
	p.execHandleFunc(ctx, p.respCallbacks)
	callBackRWMutex.RUnlock()
}

func (p *Photon) execHandleFunc(ctx *Context, callbacks []HandlerFunc) {
	for i := 0; i < len(callbacks); i++ {
		callbacks[i](ctx)
	}
}
