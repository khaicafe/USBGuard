# Makefile

# Output executable name
OUTPUT = secure_app.exe
APP_NAME=secure_app

# Source file
SRC = secure_app.go

# Build flags
GOOS = windows
GOARCH = amd64

# Run directly with go run
run:
	go run -tags=dev .
	
# Garble build command
# DEV BUILD (no USB lock)
# dấu chấm cuối tương đương có nhiêu file main
build-dev:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags=dev -o $(APP_NAME)_dev.exe .

# RELEASE BUILD (with USB lock, obfuscated)
build-release:
	GOOS=$(GOOS) GOARCH=$(GOARCH) garble -literals -tiny build -o $(APP_NAME).exe .

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) garble -literals -tiny build -o $(OUTPUT) $(SRC)

clean:
	rm -f $(OUTPUT)

.PHONY: build clean

# FULL HELP
help:
	@echo "Makefile commands:"
	@echo "  make run             → Run app in dev mode (skip USB check)"
	@echo "  make build-dev       → Build development version (no obfuscation)"
	@echo "  make build-release   → Build release version (obfuscated, USB locked)"
	@echo "  make clean           → Delete all .exe files"
