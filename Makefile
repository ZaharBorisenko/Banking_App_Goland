#Define the name of the executable
BINARY_NAME := bankingAppGo
BIN_DIR := bin
#Build target
build:
	@go build -o $(BIN_DIR)/$(BINARY_NAME)

#Run target
run: build
	@./$(BIN_DIR)/$(BINARY_NAME)

#Test target
test:
	@go test -v ./...