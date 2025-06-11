package models

type User struct {
	ID           string `bson:"_id,omitempty" json:"_id"`
	Email        string `bson:"email" json:"email"`
	PasswordHash string `bson:"passwordHash" json:"-"`
	Name         string `bson:"name" json:"name"`
	CreatedAt    int64  `bson:"createdAt" json:"createdAt"`
	UpdatedAt    int64  `bson:"updatedAt" json:"updatedAt"`
}
