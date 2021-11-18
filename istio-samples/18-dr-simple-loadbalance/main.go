package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	counter = 10000
	verChan = make(chan string, 100)
)

func main() {
	loop(counter)
}

func loop(n int) {

	set := make(map[string]int)

	go func() {
		for ver := range verChan {
			set[ver] = set[ver] + 1
		}
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			verChan <- do()
		}()
	}

	wg.Wait()
	close(verChan)

	for k, v := range set {
		fmt.Printf("%s => %d 次\n", k, v)
	}
}

type Data struct {
	Version string `json:"version"`
}

func do() string {
	resp, err := http.Get("http://istio.tangx.in/prod/list")
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	v := &Data{}
	_ = json.Unmarshal(data, v)
	return v.Version
}
