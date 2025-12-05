package models

type City struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}
