package main

import (
	"fmt"

	"github.com/n0trace/photon"
)

func main() {
	p := photon.New()
	url := "https://github.com"
	p.Get(url, func(ctx photon.Context) {
		root, err := ctx.Document()
		if err != nil {
			panic(err)
		}
		fmt.Println(root.Find("title").Text())
	})
	p.Wait()
}
