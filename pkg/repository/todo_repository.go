package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"todo-service/pkg/domain"
)

const todoCollection = "todos"

type TodoRepository interface {
	Add(ctx context.Context, todo *domain.Todo) (interface{}, error)
	GetAll(ctx context.Context) ([]domain.Todo, error)
}

type todoRepository struct {
	db *mongo.Database
}

func NewTodoRepository(db *mongo.Database) TodoRepository {
	return &todoRepository{db: db}
}

// Add Adding new todo return primitive.ObjectId and error if has.
func (t *todoRepository) Add(ctx context.Context, todo *domain.Todo) (interface{}, error) {
	one, err := t.db.Collection(todoCollection).InsertOne(ctx, todo)
	// raise an event to kafka
	if err != nil {
		return primitive.NilObjectID, errors.New(fmt.Sprintf("Error while inserting new todo %v", err))
	}
	return one.InsertedID, nil
}

func (t *todoRepository) GetAll(ctx context.Context) ([]domain.Todo, error) {
	cursor, err := t.db.Collection(todoCollection).Find(ctx, bson.M{})
	if err != nil{
		return nil, err
	}
	var todos []domain.Todo
	err = cursor.All(ctx, &todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}



