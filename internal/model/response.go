package model

type Result struct {
	Link string `json:"link"`
	Tag  int    `json:"tag"`
}

type Response struct {
	Results []Result `json:"results"`
}
