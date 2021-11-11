package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	todoservice "todo-service/pb"
	v1 "todo-service/pkg/service/v1"
)

var (
	port = flag.Int("port", 50051, "The server port")
)
func main(){
	fmt.Println("Starting the server")
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
	client, err := connectDb("mongodb://rootuser:rootpassword@localhost:27017")
	if err != nil {
		panic("can not connect to database")
	}
	db :=client.Database("oms-integration")
	// register gRPC services.
	todoservice.RegisterTodoServiceServer(s, v1.NewTodoServiceServer(db, nil))
	go func() {
		fmt.Println("Starting http Server")
		mux:= runtime.NewServeMux()
		tls, err:= loadClientTLSCredentials()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(tls)}
		err = todoservice.RegisterTodoServiceHandlerFromEndpoint(context.Background(),mux, "localhost:50051", opts)
		if err != nil {
			fmt.Println("Can not register service to http server")
			os.Exit(0)
		}

		err = http.ListenAndServe(":8081",mux)
		if err != nil {
			fmt.Println("Can not start the http server")
		}
	}()

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

func registerHttpServer(ctx context.Context){

}

func loadClientTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("deployment/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("deployment/cert/client-cert.pem", "deployment/cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientCert},
	}

	return credentials.NewTLS(config), nil
}