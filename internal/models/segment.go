package models

import "time"

type Segment struct {
	Id 			int  		`json:"id"`
	Name 		string		`json:"name"`
	CreatedAt 	time.Time 	`json:"created_at"`
}