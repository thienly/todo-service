package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	Id primitive.ObjectID `bson:"_id"`
	Title string `bson:"title"`
	Done bool `bson:"done"`
}

func (t *Todo) MarkDone(id primitive.ObjectID) {
	t.Done = true
}
