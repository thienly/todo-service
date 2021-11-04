package message

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	todoservice "todo-service/pkg/api/v1"
	"todo-service/pkg/domain"
)

const todoTopic = "todo-topic"
type TodoMessageProducer interface {
	Produce(todo *domain.Todo, converter Converter) error
}

type Converter interface {
	Convert(data interface{}) (*kafka.Message,error)
}

type messageBusProducer struct {
	producer *kafka.Producer
}

func NewMessageBus(producer *kafka.Producer) TodoMessageProducer {
	return &messageBusProducer{producer: producer}
}

func (t messageBusProducer) Produce(todo *domain.Todo, converter Converter) error{
	msg, err:=converter.Convert(todo)
	if err != nil {
		fmt.Printf("can not convert to kafka message %v", err)
		return err
	}
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	t.producer.Produce(msg, deliveryChan)
	e:= <- deliveryChan
	m:=e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		return m.TopicPartition.Error
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
	return nil
}

type todoConverter struct {
}

func NewTodoConverter() Converter {
	return &todoConverter{}
}

func (receiver todoConverter) Convert(todo interface{}) (*kafka.Message,error)  {
	d:= todo.(*domain.Todo)
	data:= &todoservice.Todo{
		Id:        &todoservice.UUID{Value: d.Id.String()},
		Name:      d.Title,
		Done:      false,
		CreatedAt: timestamppb.Now(),
	}
	byteData, err := proto.Marshal(data)
	if err != nil {
		fmt.Printf("Error while marshal %v", err)
		return nil, err
	}
	topicName:= new(string)
	*topicName = todoTopic
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     topicName,
			Partition: kafka.PartitionAny,
		},
		Key: []byte(d.Id.String()),
		Value:          byteData,
		Headers:        []kafka.Header{{Key: "Id", Value: []byte(d.Id.String())}},
	},nil
}
