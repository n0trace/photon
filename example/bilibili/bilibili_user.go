package main

import (
	"fmt"
	"log"
	"time"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

type (
	BilibiResp struct {
		Data struct {
			Card *struct {
				Name string `json:"name"`
				Sex  string `json:"sex"`
			} `json:"card"`
		} `json:"data"`
	}
)

func main() {
	p := photon.New()
	p.Use(
		middleware.Limit(time.Second),
		middleware.UserAgent(),
		middleware.Header(map[string]string{"Accept-Encoding": "gzip"}),
	)

	for i := 1; i < 10000; i++ {
		go p.Get(fmt.Sprint("https://api.bilibili.com/x/web-interface/card?mid=", i), func(ctx photon.Context) {
			var resp = new(BilibiResp)
			err := ctx.JSON(resp)
			if err != nil {
				log.Println("bilibili error", err)
				return
			}
			fmt.Println(resp.Data.Card)
		})
	}
	p.Wait()
}
