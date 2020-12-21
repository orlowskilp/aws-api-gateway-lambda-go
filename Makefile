CC=go
SRC=main.go
TARGET=hello
AR=zip

.PHONY: all
all: $(TARGET)
	zip $(TARGET).zip $(TARGET)

$(TARGET): $(SRC)
	$(CC) build -o $(TARGET)

.PHONY: clean
clean:
	rm -f $(TARGET) $(TARGET).zip