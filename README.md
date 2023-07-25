# toll-tracker
OBU Tracker with data aggregation and gateway back-end service.

## Stack: 
Kafka, GRPC, Protobuffers, Prometheus, Grafana, Microservices

## Transport between microservices:
- HTTP(JSON API) - native http, net servers and listeners, no fancy packages.
- TCP(gRPC & Protobuffers)

## Roadmap
1. OBUs (On-Board Units) => Generate OBUID with coordinates and sends them to data_receiver

2. Receiver/Producer (data_receiver) => Awaits for OBU Data and produces to Kafka

3. Kafka (docker installation with zookeeper)

4. Distance Calculator (distance_calc) => Polls data from Kafka and calculates the distance base on received

coordinates. Then uses the (5. Aggregator)'s Client to store it and process via gRPC or HTTP Transport layer.

5. Aggregator => CRUD operations for OBU Distance, invoicer, and data provider for the Gateway.

6. Gateway =>

## NOTE
This is a demo of how to structure, connect, and build transport between microservices.

It may not include the full business logic or data storage to aggregate distance.

All operations are handled in memory to prove the concept and make the application work properly.

Consider having a browse from my other projects to see the full implementation.

# Installation guide
## Kafka docker
```
docker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true bitnami/kafka:latest 
```
Also available via included docker-compose.yml file
```
docker-compose up      (or with daemon / background) -d
```

## Installing protobuf compiler (protoc compiler) 
For linux users or (WSL2) 
```
sudo apt install -y protobuf-compiler
```

For Mac users use Homebrew
```
brew install protobuf
```

## Installing GRPC and Protobuffer plugins for Golang.
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

2. GRPC 
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

3. Set the /go/bin directory to its local dir
```
PATH="${PATH}:${HOME}/go/bin"
```

4. Install the package dependencies
4.1 protobuffer package
```
go get google.golang.org/protobuf
```
4.2 grpc package
```
go get google.golang.org/grpc/
```