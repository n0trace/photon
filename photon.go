package photon

import (
	"io"
	"net/http"
	"sync"

	"github.com/n0trace/photon/common"
)

var (
	defaultParallelCount = 64
	nilHandler           = func(Context) {}
)

type (
	Photon struct {
		wait         sync.WaitGroup
		parallelChan chan interface{}
		middlewares  []MiddlewareFunc
	}

	HandlerFunc    func(Context)
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

//New new Photon Instance
func New() (p *Photon) {
	p = new(Photon)
	p.SetParallel(defaultParallelCount)
	return p
}

//SetParallel set parallel
func (p *Photon) SetParallel(parallel int) {
	p.parallelChan = make(chan interface{}, parallel)
}

func (p *Photon) Use(middleware ...MiddlewareFunc) {
	p.middlewares = append(p.middlewares, middleware...)
}

//Wait Wait
func (p *Photon) Wait() {
	p.wait.Wait()
}

func (p *Photon) process(ctx Context, cb HandlerFunc, middlewares ...MiddlewareFunc) {
	if cb == nil {
		cb = nilHandler
	}
	pre := applyMiddleware(cb, append(append(middlewares, downloadMiddleware), middlewares...)...)
	pre(ctx)
}

//Get Get
func (p *Photon) Get(url string, handle HandlerFunc, middlewares ...MiddlewareFunc) Context {
	req := common.Must2(http.NewRequest("GET", url, nil)).(*http.Request)
	return p.Request(req, handle, middlewares...)
}

//Post Post
func (p *Photon) Post(url string, contentType string, body io.Reader, handle HandlerFunc, middlewares ...MiddlewareFunc) Context {
	req := common.Must2(http.NewRequest("POST", url, body)).(*http.Request)
	req.Header.Set("Content-Type", contentType)
	return p.Request(req, handle, middlewares...)
}

//Request Request
func (p *Photon) Request(req *http.Request, cb HandlerFunc, middleware ...MiddlewareFunc) Context {
	ctx := p.NewContext()
	ctx.SetRequest(req)
	ctx.SetClient(p.NewStdClient())
	p.wait.Add(1)
	p.parallelChan <- true
	go func(ctx Context) {
		defer func() {
			<-p.parallelChan
			ctx.Finish()
			p.wait.Done()
		}()
		p.process(ctx, cb, append(p.middlewares, middleware...)...)
	}(ctx)
	return ctx
}

func (p *Photon) NewContext() *context {
	ctx := &context{
		stdClient:  p.NewStdClient(),
		photon:     p,
		finishChan: make(chan interface{}, 1),
		store:      new(sync.Map),
	}
	return ctx
}

func (p *Photon) NewStdClient() *http.Client {
	return &http.Client{}
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func downloadMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx Context) {
		defer next(ctx)
		if ctx.Downloaded() {
			return
		}
		client := ctx.Client()
		req := ctx.Request()
		resp, err := client.Do(req)
		ctx.SetStdResponse(resp)
		ctx.SetError(err)
		ctx.SetDownload(true)
	}
}
