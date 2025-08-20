OS := $(shell uname)
RM := rm -rf
BUILD_DIR := ./build
SUBMODULE_DIR := ./core/llama.cpp

all: run

.PHONY: run clean

run:
ifeq ($(OS), Linux)
	@echo "Running on Linux"
	@./scripts/build.sh
else ifeq ($(OS), Darwin)
	@echo "Running on macOS"
	@./scripts/build.sh
else
	@echo "Running on Windows or Unknown OS"
	@powershell -ExecutionPolicy Bypass -File ./scripts/build_windows.ps1
endif

clean:
	@echo "Cleaning $(BUILD_DIR) directory"
ifeq ($(OS), Linux)
	$(RM) $(BUILD_DIR)
	cd $(SUBMODULE_DIR) && git checkout . && git clean -fdx
else ifeq ($(OS), Darwin)
	$(RM) $(BUILD_DIR)
	cd $(SUBMODULE_DIR) && git checkout . && git clean -fdx
else
	@powershell -Command "if (Test-Path '$(BUILD_DIR)') { Remove-Item -Recurse -Force '$(BUILD_DIR)' }"
	cd $(SUBMODULE_DIR) && git checkout . && git clean -fdx
endif