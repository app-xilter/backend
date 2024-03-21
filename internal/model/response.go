package model

type Result struct {
	Link string `json:"link"`
}

type Response struct {
	Results []Result `json:"results"`
}
