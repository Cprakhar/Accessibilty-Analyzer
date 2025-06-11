package services

import (
	"context"
	"time"
	"backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func InitUserService(db *mongo.Database) {
	userCollection = db.Collection("users")
}

func CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = user.CreatedAt
	_, err := userCollection.InsertOne(ctx, user)
	return err
}

func FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
