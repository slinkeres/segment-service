package models

import "time"

type UserSegment struct {
	UserId 		int 		`json:"user_id"`
	SegmentId 	int			`json:"segment_id"`
	CreatedAt	time.Time	`json:"created_at"`
}