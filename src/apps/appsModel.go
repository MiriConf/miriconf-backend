package apps

import "go.mongodb.org/mongo-driver/bson/primitive"

type Apps struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Version     string             `json:"version" bson:"version"`
	Description string             `json:"description" bson:"description"`
}
