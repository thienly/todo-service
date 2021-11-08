package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"time"
	todoservice "todo-service/pkg/api/v1"
	v1 "todo-service/pkg/service/v1"
)

var (
	port = flag.Int("port", 50051, "The server port")
)
func main(){
	flag.Parse()
	lis, err:= net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	transportCredentials, err := loadTLS()
	if err != nil {
		log.Fatalf("Failed to load certificate %v", err)
	}
	s:= grpc.NewServer(grpc.Creds(transportCredentials))
	//client, err := connectDb("mongodb://rootuser:rootpassword@localhost:27017")
	if err != nil {
		panic("can not connect to database")
	}
	//db :=client.Database("oms-integration")
	//producer:= newProducer()
	todoservice.RegisterTodoServiceServer(s, v1.NewTodoServiceServer(nil, nil))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	
}
func loadTLS() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile("deployment/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	serverCert, err := tls.LoadX509KeyPair("deployment/cert/server-cert.pem","deployment/cert/server-key.pem")
	if err != nil {
		log.Fatalf(" can not load server tls %v", err)
		return nil, err
	}
	config:= &tls.Config{Certificates: []tls.Certificate{serverCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool}
	return credentials.NewTLS(config), nil
}

func connectDb(uri string) (*mongo.Client, error) {
	client, err:= mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	err = client.Connect(ctx)
	if err != nil {
	}
	return client, err
}
func newProducer() *kafka.Producer{
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092 todo-topic"})
	if err != nil {
		panic(fmt.Sprintf("can not connect to kafka %v", err))
	}
	return p
}
