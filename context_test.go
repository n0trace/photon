package photon_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/n0trace/photon"
)

func Test_context_Request(t *testing.T) {
	p := photon.New()
	c := p.NewContext()
	req, _ := http.NewRequest("GET", "https://localhost", nil)
	c.SetRequest(req)
	if got := c.Request(); !reflect.DeepEqual(got, req) {
		t.Errorf("context.Request() = %v, want %v", got, req)
	}
}

func Test_context_StdResponse(t *testing.T) {
	p := photon.New()
	c := p.NewContext()
	resp := &http.Response{}
	c.SetStdResponse(resp)
	if got := c.StdResponse(); !reflect.DeepEqual(got, resp) {
		t.Errorf("context.StdResponse() = %v, want %v", got, resp)
	}
}

func Test_context_Client(t *testing.T) {
	p := photon.New()
	c := p.NewContext()
	client := &http.Client{}
	c.SetClient(client)
	if got := c.Client(); !reflect.DeepEqual(got, client) {
		t.Errorf("context.Client() = %v, want %v", got, client)
	}
}

func Test_context_Error(t *testing.T) {
	p := photon.New()
	c := p.NewContext()
	err := errors.New("new error")
	c.SetError(err)
	if got := c.Error(); !reflect.DeepEqual(got, err) {
		t.Errorf("context.Error() = %v, want %v", got, err)
	}
}

func Test_context_Get(t *testing.T) {
	p := photon.New()
	c := p.NewContext()
	var m = make(map[string]string)
	var s = make([]string, 8)
	var i = 8
	c.Set("m", m)
	c.Set("s", s)
	c.Set("i", i)
	if got, _ := c.Get("m"); !reflect.DeepEqual(got.(map[string]string), m) {
		t.Errorf("context.Get(map) = %v, want %v", got, m)
	}
	if got, _ := c.Get("s"); !reflect.DeepEqual(got.([]string), s) {
		t.Errorf("context.Get(slice) = %v, want %v", got, s)
	}
	if got, _ := c.Get("i"); !reflect.DeepEqual(got.(int), i) {
		t.Errorf("context.Get(int) = %v, want %v", got, i)
	}
}

func Test_context_Downloaded(t *testing.T) {
	p := photon.New()
	c := p.NewContext()

	if got := c.Downloaded(); !reflect.DeepEqual(got, false) {
		t.Errorf("context.Downloaded() = %v, want %v", got, false)
	}

	c.SetDownload(true)
	if got := c.Downloaded(); !reflect.DeepEqual(got, true) {
		t.Errorf("context.Downloaded() = %v, want %v", got, true)
	}

	c.SetDownload(false)
	if got := c.Downloaded(); !reflect.DeepEqual(got, false) {
		t.Errorf("context.Downloaded() = %v, want %v", got, false)
	}
}

func Test_context_WaitFinish(t *testing.T) {
	p := photon.New()
	n := 16

	exited := make(chan bool, n)
	ctx := p.NewContext()
	go func() {
		ctx.WaitFinish()
		exited <- true
	}()

	select {
	case <-exited:
		t.Fatal("Context finish too soon")
	default:
	}
	ctx.Finish()
	<-exited
}
