package main

import (
	"fmt"
	"time"

	"github.com/n0trace/photon/middleware"

	"github.com/PuerkitoBio/goquery"

	"github.com/n0trace/photon"
)

func main() {
	p := photon.New()
	p.Use(middleware.Limit(time.Second))

	p.Get("https://github.com/search?q=go", func(ctx photon.Context) {
		root, err := ctx.Document()
		if err != nil {
			panic(err)
		}
		root.Find(".repo-list .v-align-middle").Each(func(idx int, selection *goquery.Selection) {
			href, ok := selection.Attr("href")
			if !ok {
				return
			}
			u := "https://github.com" + href
			p.Get(u, printTitle)
		})
	})

	p.Wait()
}

func printTitle(ctx photon.Context) {
	root, _ := ctx.Document()
	fmt.Println(root.Find("title").Text())
}
