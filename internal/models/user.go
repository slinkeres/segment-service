package models

import "time"

type User struct{
	Id 			int  		`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
}