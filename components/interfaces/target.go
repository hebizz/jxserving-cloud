package interfaces

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Target struct {
  Id      primitive.ObjectID `bson:"_id"`
  LabelId string             `json:"id"`
  Name    string             `json:"name" binding:"required"`
  D       [] string          `json:"labels" binding:"required"` // data [12,23,55,55]
}
