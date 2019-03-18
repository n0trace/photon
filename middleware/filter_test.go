package middleware_test

import (
	"log"
	"sync/atomic"
	"testing"

	"github.com/n0trace/photon"

	"github.com/n0trace/photon/middleware"
)

func TestFilter(t *testing.T) {
	p := photon.New()
	p.Use(middleware.Filter())
	var times int64
	for i := 0; i < 10; i++ {
		p.Get(newTestServer().URL, func(ctx photon.Context) {
			atomic.AddInt64(&times, 1)
			log.Println("lalala")
		})
	}
	if atomic.LoadInt64(&times) != 1 {
		t.Errorf("callback function execution %v times, want %v", atomic.LoadInt64(&times), 1)
	}
}
