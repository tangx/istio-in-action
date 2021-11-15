package model

type Product struct {
	Name    string
	Price   int
	Reviews interface{}
}

type Review struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"commment"`
}
