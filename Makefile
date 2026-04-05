SHELL := /bin/bash # For Completion Now

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean

# Paths and names
SRC_FOLDER=./cmd/hermInvestCli
SRC_PATH=$(SRC_FOLDER)
BIN_NAME=hermInvestCli
DB=./internal/app/database/dev-database.db
CREATE_DB_CMD=./cmd/internal/createDBSchema/createDBSchema.go
SEED_DATA_CMD=./cmd/internal/seedSampleData/seedSampleData.go

# Completion
WORKING_DIR=$(shell pwd)
GENERATE_COMPLETION=$(WORKING_DIR)/$(BIN_NAME) completion bash --no-descriptions
BASH_COMPLETION_DIR=/usr/share/bash-completion/completions

# Detect the OS
ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
    detected_OS := Windows
else
    detected_OS := $(shell uname)  # same as "uname -s"  is Linux
endif

# Windows-specific variables and commands
ifeq ($(detected_OS),Windows)
    BIN_EXT=.exe
    RM=del /f
else
    BIN_EXT=
    RM=rm -f
endif

all: exec

$(BIN_NAME)$(BIN_EXT): $(SRC_PATH)
	$(MAKE) build

build:
	$(GOBUILD) -o $(BIN_NAME)$(BIN_EXT) $(SRC_PATH)

exec: $(BIN_NAME)$(BIN_EXT)
	./$(BIN_NAME)$(BIN_EXT)

clean:
	$(GOCLEAN)
	$(RM) $(BIN_NAME)$(BIN_EXT)

run:
	$(GORUN) $(SRC_PATH).go

example: $(DB) $(BIN_NAME)$(BIN_EXT)
	./$(BIN_NAME)$(BIN_EXT)
	./$(BIN_NAME)$(BIN_EXT) stock query --all

$(DB):
	@echo "DB $(DB) does not exist"
	$(GORUN) $(CREATE_DB_CMD)
	$(GORUN) $(SEED_DATA_CMD)

# sudo required
completion:
	$(GENERATE_COMPLETION) > $(BASH_COMPLETION_DIR)/$(BIN_NAME)
	source $(BASH_COMPLETION_DIR)/$(BIN_NAME)