# toll-tracker
Big data aggregation tool with several microservices. Stack: Kafka, Prometheus, GRCP, Protobuffer, Grafana

## Roadmap
OBUs => Receiver => Kafka Queue

# Installation guide
## Kafka docker
```
docker run --name kafka -p 9092:9092 -e ALLOW_PLAINTEXT_LISTENER=yes -e KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true bitnami/kafka:latest 
```