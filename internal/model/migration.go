package model

import "time"

type Tags struct {
	Id   int    `gorm:"primaryKey index autoIncrement"`
	Name string `gorm:"unique"`
}

type Tweets struct {
	Id        int    `gorm:"primaryKey"`
	Link      string `gorm:"unique index"`
	TagId     int    `gorm:"foreignKey:TagsId"`
	Content   string
	IPAddress string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
