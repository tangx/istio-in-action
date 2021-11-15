package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tangx/istio-in-action/cmd/model"
)

func main() {
	r := gin.Default()

	r.GET("/review/all", reviewHandler)
	r.GET("/review/delay", delayHanlder)

	if err := r.Run(":8089"); err != nil {
		panic(err)
	}
}

var reviews map[string]model.Review

func reviewHandler(c *gin.Context) {
	c.JSON(http.StatusOK, reviews)
}

func delayHanlder(c *gin.Context) {
	delay := c.DefaultQuery("delay", "3")
	n := func() int {
		n, err := strconv.Atoi(delay)
		if err != nil {
			return 3
		}
		return n
	}()

	// 休眠
	time.Sleep(time.Duration(n) * time.Second)

	// 返回结果
	reviewHandler(c)
}

func init() {
	reviews = map[string]model.Review{
		"1": {
			ID:      "1",
			Name:    "zhangsan",
			Comment: "istio 功能很强大， 就是配置太麻烦",
		},
		"2": {
			ID:      "1",
			Name:    "wangwu",
			Comment: "《istio in action》 真是一本了不起的书",
		},
	}
}
