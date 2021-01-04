CC=go
SRC=main.go pkg/dynamodb/dynamodb.go
TARGET=hello
AR=zip

FILES_TO_MOCK=pkg/dynamodb/dynamodb.go
MOCK_FILES=$(FILES_TO_MOCK:%.go=%_mock.go)
FILES_TO_TEST=$(FILES_TO_MOCK)
TEST_FILES=$(FILES_TO_TEST:%.go=%_test.go)
INTEGRATION_TEST_FILES=$(FILES_TO_TEST:%.go=%_integration_test.go)

DDB_SETUP_SCRIPT=dynamodb/create-table.sh
DDB_POPULATE_SCRIPT=dynamodb/populate-table.sh
DDB_CLEANUP_SCRIPT=dynamodb/delete-table.sh

# Artifact building recipes
.PHONY: all
all: $(TARGET)
	zip $(TARGET).zip $(TARGET)

$(TARGET): $(SRC)
	$(CC) get -d -v ./...
	$(CC) build -o $(TARGET)

# Mock file generating recipes
.PHONY: mock
mock: $(MOCK_FILES)

%_mock.go: %.go
	mockgen -source=$< -package=dynamodb -destination=$@

# Runs only unit tests
.PHONY: test
test: $(FILES_TO_TEST) $(TEST_FILES)
	$(CC) get -d -t -v ./...
	$(CC) test $^

# Runs only integration tests (on local DynamoDB instance)
.PHONY: integration_test
integration_test: $(FILES_TO_TEST) $(INTEGRATION_TEST_FILES) $(DDB_SETUP_SCRIPT) $(DDB_POPULATE_SCRIPT) $(DDB_CLEANUP_SCRIPT)
	@echo "##!!! If this fails, remove local DynamoDB table and try again !!!##"
	$(CC) get -d -t -v ./...
	./$(DDB_SETUP_SCRIPT)
	./$(DDB_POPULATE_SCRIPT)
	$(CC) test $(FILES_TO_TEST) $(INTEGRATION_TEST_FILES)
	./$(DDB_CLEANUP_SCRIPT)

.PHONY: clean
clean:
	rm -f $(TARGET) $(TARGET).zip