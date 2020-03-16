package user

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// Location stores user location data
type Location struct {
	Lat string `json:"lat,required"`
	Lng string `json:"lng,required"`
}

// User stores user data
type User struct {
	UUID      uuid.UUID `json:"uuid" bson:"_id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Name string `json:"name" bson:"name"`
}
