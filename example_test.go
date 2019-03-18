package photon_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

func Example() {
	p := photon.New()
	p.Get(newTestServer().URL+"/users?id=hello", func(ctx photon.Context) {
		text, _ := ctx.Text()
		fmt.Println(text)
	})
	p.Wait()
	//Output:
	//hello
}

func Example_useMiddleware() {
	rootURL := newTestServer().URL
	p := photon.New()
	p.Use(middleware.Limit(200*time.Millisecond), middleware.UserAgent("diy-agent"))
	for i := 0; i != 3; i++ {
		url := fmt.Sprintf("%s/user-agent", rootURL)
		p.Get(url, func(ctx photon.Context) {
			text, _ := ctx.Text()
			fmt.Println(text)
		})
	}
	//or:
	//p.Get(url,callback,middleware...)
	p.Wait()
	//Output:
	//diy-agent
	//diy-agent
	//diy-agent
}

func Example_keepAuth() {
	p := photon.New()

	reader := strings.NewReader("username=foo&password=bar")

	lastCtx := p.Post(newTestServer().URL+"/login",
		"application/x-www-form-urlencoded", reader,
		func(ctx photon.Context) {
			text, _ := ctx.Text()
			fmt.Println(text)
		})

	p.Get(newTestServer().URL+"/must-login", func(ctx photon.Context) {
		text, _ := ctx.Text()
		fmt.Println(text)
	}, middleware.FromContext(lastCtx))

	p.Wait()

	//Output:
	//ok
	//hello foo
}
