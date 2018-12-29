package photon

import "net/http"

type Context struct {
	*Response
	error   error
	request *http.Request
	*http.Client
	meta       map[string]interface{}
	Stage      ContextStage
	PreContext *Context
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) SetRequest(req *http.Request) {
	c.request = req
}

func (c *Context) SetMeta(m map[string]interface{}) {
	c.meta = m
}

func (c *Context) Reset() {

}

func (c *Context) Error() error {
	return c.error
}

func (c *Context) SetError(err error) {
	c.error = err
}
