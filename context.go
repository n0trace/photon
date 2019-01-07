package photon

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Context interface {
	Request() *http.Request
	SetRequest(*http.Request)
	Client() *http.Client
	SetClient(*http.Client)
	Reset()
	Error() error
	SetError(error)
	SetStdResponse(*http.Response)
	StdResponse() *http.Response
	Text() (string, error)
	Bytes() ([]byte, error)
	JSON(interface{}) error
	XML(interface{}) error
	Document() (*goquery.Document, error)
	Photon() *Photon
	Set(string, interface{})
	Get(string) interface{}
	FromContext(interface{})
	Downloaded() bool
	SetDownload(bool)
}

type context struct {
	Response
	error      error
	stdRequest *http.Request
	stdClient  *http.Client
	photon     *Photon
	store      map[string]interface{}
}

func (c *context) Request() *http.Request {
	return c.stdRequest
}

func (c *context) SetStdResponse(resp *http.Response) {
	c.Response.Response = resp
}

func (c *context) SetRequest(req *http.Request) {
	c.stdRequest = req
}

func (c *context) Client() *http.Client {
	return c.stdClient
}

func (c *context) SetClient(client *http.Client) {
	c.stdClient = client
}

func (c *context) Reset() {

}

func (c *context) Error() error {
	return c.error
}

func (c *context) SetError(err error) {
	c.error = err
}

func (c *context) StdResponse() *http.Response {
	return c.Response.Response
}

func (c *context) Photon() *Photon {
	return c.photon
}

func (c *context) Set(key string, value interface{}) {
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = value
}

func (c *context) Get(key string) interface{} {
	return c.store[key]
}

func (c *context) SetDownload(d bool) {
	c.Set("downloaded", true)
}

func (c *context) Downloaded() bool {
	switch downloaded := c.Get("downloaded").(type) {
	case bool:
		return downloaded
	default:
		return false
	}
}

func (c *context) FromContext(from interface{}) {
	f := from.(*context)
	c.store = f.store
	c.store["downloaded"] = false
	c.SetClient(c.Client())

	jar := new(Jar)
	jar.SetCookies(c.Request().URL, f.StdResponse().Cookies())
	c.Client().Jar = jar
	c.Request().Header = f.Request().Header
}
