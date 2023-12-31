gate:
	@go build -o bin/gateway ./gateway
	@./bin/gateway

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calc:
	@go build -o bin/calc ./distance_calc
	@./bin/calc

agg:
	@go build -o bin/agg ./aggregator
	@./bin/agg

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

temp:
	@go run temp/main.go

prometheus:
	@./../prometheus/prometheus --config.file=.config/prometheus.yml

gokit:
	@go run go-kit-example/aggsvc/cmd/main.go

.PHONY: obu receiver temp