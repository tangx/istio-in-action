package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tangx/istio-in-action/cmd/model"
	"github.com/tangx/istio-in-action/version"
)

func main() {
	r := gin.Default()

	routes(r)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func listHandler(c *gin.Context) {

	reviews, err := getReivews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "获取评论失败， 内部错误",
		})
		return
	}

	product := model.Product{
		Name:    "istio in action",
		Price:   300,
		Reviews: reviews,
	}

	c.JSON(http.StatusOK, gin.H{
		"version": version.Version,
		"data":    product,
	})
}

func routes(e *gin.Engine) {
	rg := e.Group("/prod")

	rg.GET("/list", listHandler)
}

func getReivews() (map[string]model.Review, error) {

	reviews := make(map[string]model.Review)

	resp, err := http.Get("http://svc-review/review/all")
	if err != nil {
		return nil, fmt.Errorf("reqeust svc-review failed: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %v", err)
	}

	err = json.Unmarshal(data, &reviews)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal data failed: %v", err)
	}

	return reviews, nil
}
