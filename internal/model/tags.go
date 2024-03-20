package model

type SystemTag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SystemTagResponse struct {
	Tags []SystemTag `json:"tags"`
}
