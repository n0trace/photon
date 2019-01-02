package main

import (
	"encoding/json"
	"fmt"

	"github.com/n0trace/photon"

	"github.com/PuerkitoBio/goquery"
)

var (
	userAgentMap = make(map[string][]string)
)

func main() {
	p := photon.New()
	p.OnResponse(func(ctx *photon.Context) {
		document, err := ctx.Document()
		if err != nil {
			panic(err)
		}
		document.Find("h3").Each(func(index int, selectionH3 *goquery.Selection) {
			userAgentMap[selectionH3.Text()] = []string{}
			selectionH3.NextUntil("h3").Find("ul li a").Each(func(idx int, selectionA *goquery.Selection) {
				userAgentMap[selectionH3.Text()] = append(userAgentMap[selectionH3.Text()], selectionA.Text())
			})
		})
		bs, err := json.Marshal(userAgentMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bs))
	})

	p.OnError(func(ctx *photon.Context) {
		panic(ctx.Error())
	})

	p.Visit("http://useragentstring.com/pages/useragentstring.php?name=All")

	p.Wait()
}
