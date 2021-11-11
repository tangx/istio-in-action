package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/istio-in-action/version"
)

func main() {
	r := gin.Default()

	routes(r)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

type Product struct {
	Name    string
	Price   int
	Reviews interface{}
}

func listHandler(c *gin.Context) {

	product := Product{
		Name:  "istio in action",
		Price: 300,
	}
	c.JSON(http.StatusOK, gin.H{
		"version": version.Version,
		"data":    product,
	})
}

func routes(e *gin.Engine) {
	rg := e.Group("/prods")

	rg.GET("/list", listHandler)
}
