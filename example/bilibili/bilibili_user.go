package main

import (
	"fmt"
	"log"
	"time"

	"github.com/n0trace/photon"
	"github.com/n0trace/photon/middleware"
)

type (
	BilibiliUser struct {
		Card *struct {
			Name string `json:"name"`
			Sex  string `json:"sex"`
		} `json:"card"`
	}
	BilibiResp struct {
		Data *BilibiliUser `json:"data"`
	}
)

var (
	ticker    = time.NewTicker(time.Microsecond)
	limitFunc = func() <-chan interface{} {
		out := make(chan interface{})
		go func() {
			for t := range ticker.C {
				out <- t
			}
		}()
		return out
	}
)

func main() {
	p := photon.New(
		photon.WithParallel(100),
		photon.WithLimitFunc(limitFunc),
	)
	p.Use(middleware.RandomUserAgent("ABACHOBot", "008"))
	p.On(photon.OnResponse, func(ctx *photon.Context) (err error) {
		var resp = new(BilibiResp)
		err = ctx.JSON(resp)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(resp.Data.Card)
		return
	})

	p.On(photon.OnRequest, func(ctx *photon.Context) (err error) {
		req := ctx.Request()
		req.Header.Set("Content-Type", "application/json")
		return nil
	})

	p.On(photon.OnError, func(ctx *photon.Context) (err error) {
		log.Println(ctx.Error())
		return nil
	})

	for i := 1; i < 10000; i++ {
		p.Visit(fmt.Sprint("https://api.bilibili.com/x/web-interface/card?mid=", i), photon.VisitWithMeta(map[string]interface{}{}))
	}

	p.Wait()
}
