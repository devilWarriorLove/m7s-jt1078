package main

import (
	"context"
	"fmt"
	"github.com/cuteLittleDevil/go-jt808/service"
	_ "github.com/cuteLittleDevil/m7s-jt1078/v5"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"m7s.live/v5"
	_ "m7s.live/v5/plugin/flv"
	_ "m7s.live/v5/plugin/mp4"
	_ "m7s.live/v5/plugin/preview"
)

func init() {
	go func() {
		_ = m7s.Run(context.Background(), "./config.yaml")
	}()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	{
		goJt808 := service.New(
			service.WithHostPorts("0.0.0.0:11001"),
			service.WithCustomTerminalEventer(func() service.TerminalEventer {
				return &LogTerminal{}
			}),
		)
		go goJt808.Run()
		r.Use(func(c *gin.Context) {
			c.Set("jt808", goJt808)
		})
	}

	group := r.Group("/api/v1/jt808/")
	{
		group.POST("/9101", p9101)
		group.POST("/9102", p9102)
	}
	onEvent := r.Group("/api/v1/jt808/event/")
	{
		onEvent.POST("/join-audio", onEventJoinAudio)
		onEvent.POST("/real-time-join", onEventRealTimeJoin)
		onEvent.POST("/real-time-leave", onEventRealTimeLeave)
	}
	r.Static("/", "./static")
	fmt.Println("服务已启动 默认首页:", "http://124.221.30.46:11000")
	_ = r.Run(":11000")
}
