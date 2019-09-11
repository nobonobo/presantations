PACKAGE_PATH := github.com/nobonobo/presantations/tinygo/examples/sample1
PROJECT_PATH := /go/src/$(PACKAGE_PATH)
TARGET  := circuitplay-express
SOURCES := $(shell find . -type f -name '*.go')
OUTPUT  := sample.uf2
FLASH_PATH := /Volumes/CPLAYBOOT

build: $(OUTPUT)

clean:
	@-rm $(OUTPUT)

flash: $(OUTPUT)
	cp $(OUTPUT) $(FLASH_PATH)

$(OUTPUT): $(SOURCES)
	docker run -it --rm -v $(PWD):$(PROJECT_PATH) -w $(PROJECT_PATH) tinygo/tinygo tinygo build -target=$(TARGET) -o $(OUTPUT) .
