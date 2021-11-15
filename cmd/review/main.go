package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/istio-in-action/cmd/model"
)

func main() {
	r := gin.Default()

	r.GET("/review/all", handler)
	if err := r.Run(":8089"); err != nil {
		panic(err)
	}
}

var reviews map[string]model.Review

func handler(c *gin.Context) {
	c.JSON(http.StatusOK, reviews)
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
