package models

import "time"

type Checkin struct {
	ActivityCreatedAt time.Time `json:"activity_created_at"`
	Note              *string   `json:"note"`
	Location          *string   `json:"location"`
	Type              string    `json:"type"`
	ResourceUrl       *string   `json:"resource_url"`
}
