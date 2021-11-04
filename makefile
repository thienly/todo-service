cert:
	cd deployment/cert; chmod +x ./gen.sh; ./gen.sh; cd ../../
server:
	go run cmd/server/main.go
client:
	go run cmd/client/main.go
consumer:
	go run cmd/consumer/main.go localhost:9092 group1 todo-topic
kafka-stop:
	cd build/kafka; docker-compose down
kafka-start:
	cd build/kafka; docker-compose down