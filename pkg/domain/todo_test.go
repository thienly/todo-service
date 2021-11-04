package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestMarkDone(t *testing.T) {
	id:= primitive.NewObjectID()
	todo := Todo{
		Id:    id,
		Title: "Hello",
		Done:  false,
	}
	todo.MarkDone(id)
	if todo.Done != true{
		t.Fail()
	}
}
