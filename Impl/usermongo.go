package dataimpl

import (
	data "main/Data"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoData struct {
	data.UserRepo
	client *mongo.Client
}

// func NewMongo(client *mongo.Client) (*UserData, error) {
// 	return &UserData{client: client}, nil
// }

// func (d UserData) FindByUsername(username string) (model.UserDetail, error) {
// 	// mongodb query
// 	// db.user.find({})
// 	return model.UserDetail{}, nil
// }
