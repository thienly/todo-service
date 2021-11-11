package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	todoservice "todo-service/pb"
)

func main(){
	tlsCredentials, _ := loadTLSCredentials()

	conn, err:= grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		fmt.Printf("can not dial to :50051 reasonn %v", err)
	}

	todoClient := todoservice.NewTodoServiceClient(conn)
	//request := &todoservice.TodoRequest{Todo: &todoservice.Todo{
	//	Id:        &todoservice.UUID{Value: "123"},
	//	Name:      "Test",
	//	Done:      false,
	//	CreatedAt: timestamppb.Now(),
	//}}
	//response, err:= todoClient.Create(context.Background(), request)
	response, err:= todoClient.Sample(context.Background(), &todoservice.Void{})
	if err != nil{
		fmt.Printf("error while calling %v", err)
	}
	fmt.Println(response)

}
func loadTLSCredentials() (credentials.TransportCredentials, error) {
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