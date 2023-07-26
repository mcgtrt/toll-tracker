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

5. 1. Prometheus implementation will only be present as the middleware for aggregator interface methods

6. Gateway => User facing API, communicating with Aggregator API

## NOTE
This is a demo of how to structure, connect, and build transport between microservices.

It may not include the full business logic, all endpoints or data storage (like NoSQL).

All operations are handled in memory to prove the concept and make the application work properly.

For business logic, check out Hotel API, Crypto, or other projects.

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
4. 1. Protobuffer package
```
go get google.golang.org/protobuf
```

4. 2. gRPC Package
```
go get google.golang.org/grpc/
```

## Prometheus
1. Install Prometheus (Docker)
```
docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus
```

2. Install prometheus client
```
go get github.com/prometheus/client_golang/prometheus
```

## Grafana
```
docker run -d -p 3000:3000 --name=grafana grafana/grafana-enterprise
```