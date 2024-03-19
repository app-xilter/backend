package model

type Tweet struct {
	Link string `json:"link" validate:"required"`
	Text string `json:"text" validate:"required"`
}

type Tweets struct {
	Tweets []Tweet `json:"tweets" validate:"required,dive"`
}
