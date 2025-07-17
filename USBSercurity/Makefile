# Makefile

# Output executable name
OUTPUT = secure_app.exe

# Source file
SRC = secure_app.go

# Build flags
GOOS = windows
GOARCH = amd64

# Run directly with go run
run:
	go run $(SRC)
	
# Garble build command
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) garble -literals -tiny build -o $(OUTPUT) $(SRC)

clean:
	rm -f $(OUTPUT)

.PHONY: build clean
