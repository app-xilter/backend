package model

type Result struct {
	Link    string `json:"link"`
	Tag     int    `json:"tag"`
	TagName string `json:"tag_name"`
}

type Response struct {
	Results []Result `json:"results"`
}
