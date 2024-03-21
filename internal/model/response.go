package model

type Result struct {
	Link     string `json:"link"`
	Filtered bool   `json:"filtered"`
}

type Response struct {
	Results []Result `json:"results"`
}
