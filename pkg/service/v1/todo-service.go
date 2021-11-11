package v1

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	todoservice "todo-service/pb"
	"todo-service/pkg/domain"
	message "todo-service/pkg/kafka"
	"todo-service/pkg/repository"
)

type todoServiceServer struct {
	todoRepo repository.TodoRepository
	todoMessage message.TodoMessageProducer
	todoservice.UnimplementedTodoServiceServer
}

func (t *todoServiceServer) Create(ctx context.Context, request *todoservice.TodoRequest) (*todoservice.TodoResponse, error) {
	return &todoservice.TodoResponse{Id: primitive.NewObjectID().String()}, nil
	todo:= request.Todo
	domainTodo := &domain.Todo{
		Id:    primitive.NewObjectID(),
		Title: todo.Name,
		Done:  false,
	}
	add, err := t.todoRepo.Add(ctx, domainTodo)
	if err != nil {
		log.Fatalf("Can not insert into database %v", err)
		return nil, err
	}
	converter:= message.NewTodoConverter()
	err = t.todoMessage.Produce(domainTodo, converter)
	if err != nil {
		log.Fatalf("Can not send to kafka %v", err)
	}
	idStr := add.(primitive.ObjectID).String()
	return &todoservice.TodoResponse{Id: idStr},nil
}

func (t *todoServiceServer)	 GetAll(ctx context.Context, empty *todoservice.Void) (*todoservice.TodoList, error) {
	all, err := t.todoRepo.GetAll(ctx)
	if err != nil{
		return nil, err
	}
	var result todoservice.TodoList
	for _, v := range all {
		result.Data = append(result.Data, &todoservice.Todo{
			Id:        &todoservice.UUID{Value: v.Id.String()},
			Name:      v.Title,
			Done:      v.Done,
			CreatedAt: nil,
		})
	}
	return &result, nil
}
func (t *todoServiceServer)  Sample(ctx context.Context, emp *todoservice.Void) (*todoservice.TodoResponse, error){
	return &todoservice.TodoResponse{Id: primitive.NewObjectID().String()}, nil
}

func (t todoServiceServer) mustEmbedUnimplementedTodoServiceServer() {
	panic("implement me")
}

func NewTodoServiceServer(db *mongo.Database, producer *kafka.Producer) todoservice.TodoServiceServer {
	return &todoServiceServer{
		todoRepo:                       repository.NewTodoRepository(db),
		todoMessage:                    message.NewMessageBus(producer),
		UnimplementedTodoServiceServer: todoservice.UnimplementedTodoServiceServer{},
	}
}


