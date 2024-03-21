package model

type Tweet struct {
	Link string `json:"link" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type Tag struct {
	Id   int    `json:"id" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type Request struct {
	Tweets []Tweet `json:"tweets" validate:"required,omitempty,min=1,dive"`
	Tags   []Tag   `json:"tags" validate:"required,omitempty,min=1,dive"`
}
