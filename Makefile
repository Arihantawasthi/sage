# Makefile for sage - process manager

BINARY_DAEMON=saged
BINARY_CLI=sagectl
INSTALL_DIR=/usr/local/bin
SERVICE_FILE=saged.service
UNIT_PATH=/etc/systemd/system
RUNTIME_DIR=/home/ubuntu/projects/sage
LOG_DIR=/var/log/sage

.PHONY: all build install clean

all: install

build:
	@echo "> Building daemon..."
	@go build -o $(BINARY_DAEMON) ./cmd/saged
	@echo "> Building CLI tool..."
	@go build -o $(BINARY_CLI) ./cmd/sagectl

install: build
	@echo "> Installing binaries to $(INSTALL_DIR)..."
	sudo cp $(BINARY_DAEMON) $(INSTALL_DIR)/$(BINARY_DAEMON) || echo "FAILED"
	sudo cp $(BINARY_CLI) $(INSTALL_DIR)/$(BINARY_CLI)

	@echo "> Creating runtime and log directories..."
	@sudo mkdir -p $(LOG_DIR)
	@sudo mkdir -p $(RUNTIME_DIR)

	@echo "> Installing systemd unit file..."
	@sudo cp ./scripts/$(SERVICE_FILE) $(UNIT_PATH)/

	@echo "> Reloading systemd daemon..."
	@sudo systemctl daemon-reexec
	@sudo systemctl daemon-reload

	@echo "> Done."

clean:
	@echo "> Cleaning binaries..."
	@rm -f $(BINARY_DAEMON) $(BINARY_CLI)

