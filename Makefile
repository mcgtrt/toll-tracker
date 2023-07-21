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

.PHONY: obu receiver