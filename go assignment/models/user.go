package models

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name" validate:"required"`
	Email     string    `bson:"email" validate:"required,email"`
	Age       int       `bson:"age" validate:"gte=0"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
