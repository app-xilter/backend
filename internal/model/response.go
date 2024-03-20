package model

type Result struct {
	Link     string `json:"link"`
	Category int    `json:"category"`
}

type Response struct {
	Results []Result `json:"results"`
}
