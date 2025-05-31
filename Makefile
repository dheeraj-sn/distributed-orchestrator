# Go binary build output
BIN_DIR := bin

# Commands
build-scheduler:
	go build -o $(BIN_DIR)/scheduler ./cmd/scheduler

build-worker:
	go build -o $(BIN_DIR)/worker ./cmd/worker

build-client:
	go build -o $(BIN_DIR)/client ./cmd/client

build-all: build-scheduler build-worker build-client

run-scheduler:
	$(BIN_DIR)/scheduler

run-worker:
	SCHEDULER_ADDR=localhost:50051 $(BIN_DIR)/worker

run-client:
	SCHEDULER_ADDR=localhost:50051 $(BIN_DIR)/client

clean:
	rm -rf $(BIN_DIR)

.PHONY: build-scheduler build-worker build-client build-all run-scheduler run-worker run-client clean