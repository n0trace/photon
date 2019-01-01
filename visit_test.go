package photon_test

import (
	"net/http"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/n0trace/photon"
)

func TestVisitWithMeta(t *testing.T) {
	p := photon.New()
	server := newTestServer()
	meta := map[string]interface{}{"hello": "world"}
	var respMeta map[string]interface{}
	p.Visit(server.URL+"/user-agent", photon.VisitWithMeta(meta))
	p.On(photon.OnResponse, func(ctx *photon.Context) error {
		respMeta = ctx.Meta()
		return nil
	})
	p.Wait()
	if !reflect.DeepEqual(meta, respMeta) {
		t.Errorf("VisitWithMeta() = %v, want %v", respMeta, meta)
	}
}

func TestVisitWithClient(t *testing.T) {
	p := photon.New()
	client := &http.Client{
		Timeout: time.Second,
	}
	var wantClient *http.Client
	p.Visit(newTestServer().URL+"/user-agent", photon.VisitWithClient(client))
	p.On(photon.OnResponse, func(ctx *photon.Context) error {
		wantClient = ctx.Client
		return nil
	})
	p.Wait()
	if !reflect.DeepEqual(client, wantClient) {
		t.Errorf("VisitWithClient() = %v, want %v", client, wantClient)
	}
}

func TestVisitWithDontFilter(t *testing.T) {
	server := newTestServer()
	p := photon.New()
	var times int64

	p.On(photon.OnResponse, func(ctx *photon.Context) error {
		atomic.AddInt64(&times, 1)
		return nil
	})

	for i := 0; i < 100; i++ {
		p.Visit(server.URL+"/user-agent", photon.VisitWithDontFiter())
	}

	p.Wait()

	if atomic.LoadInt64(&times) != 100 {
		t.Errorf("VisitWithDontFilter() visit same url times = %v, want %v", atomic.LoadInt64(&times), 100)
	}
}

func TestVisitWithFilter(t *testing.T) {
	server := newTestServer()
	p := photon.New()
	var times1 int64

	p.On(photon.OnResponse, func(ctx *photon.Context) error {
		atomic.AddInt64(&times1, 1)
		return nil
	})

	for i := 0; i < 100; i++ {
		p.Visit(server.URL + "/user-agent")
	}

	p.Wait()

	if atomic.LoadInt64(&times1) != 1 {
		t.Errorf("VisitWithFilter() visit same url times = %v, want %v", atomic.LoadInt64(&times1), 1)
	}
}
