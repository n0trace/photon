package main

import (
	"fmt"

	"github.com/n0trace/photon"
)

func main() {
	p := photon.New()
	go p.Get("https://google.com", func(c photon.Context) {
		fmt.Println(c.Text())
	})
	p.Get("https://zhihu.com", func(c photon.Context) {
		fmt.Println(c.Text())
	})
}
