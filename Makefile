# ── Variables ────────────────────────────────────────────────────────

BIN_DIR    := $(CURDIR)/bin

UNAME_S    := $(shell uname -s | tr '[:upper:]' '[:lower:]')
UNAME_M    := $(shell uname -m)

ifeq ($(UNAME_M),arm64)
  ARCH := aarch64
else ifeq ($(UNAME_M),aarch64)
  ARCH := aarch64
else
  ARCH := $(UNAME_M)
endif

ifeq ($(UNAME_S),darwin)
  DPRINT_TARGET := $(ARCH)-apple-darwin
else ifeq ($(UNAME_S),linux)
  DPRINT_TARGET := $(ARCH)-unknown-linux-gnu
endif

DPRINT_VERSION     := 0.54.0
DPRINT             := $(BIN_DIR)/dprint
DPRINT_RELEASE_URL := https://github.com/dprint/dprint/releases/download/$(DPRINT_VERSION)/dprint-$(DPRINT_TARGET).zip

OBS := $(BIN_DIR)/obs

# ── Setup ────────────────────────────────────────────────────────────

.PHONY: setup clean reset-projects

setup: tools reset-projects

clean:
	rm -rf $(BIN_DIR)

reset-projects:
	@./scripts/reset-projects.sh

# ── Tools ────────────────────────────────────────────────────────────

.PHONY: tools obs

tools: $(DPRINT) obs

$(DPRINT):
	@mkdir -p $(BIN_DIR)
	@echo "Downloading dprint $(DPRINT_VERSION) for $(DPRINT_TARGET)..."
	@curl -fsSL "$(DPRINT_RELEASE_URL)" -o $(BIN_DIR)/dprint.zip
	@unzip -oq $(BIN_DIR)/dprint.zip -d $(BIN_DIR)
	@rm $(BIN_DIR)/dprint.zip
	@chmod +x $(DPRINT)
	@echo "Installed dprint -> $(DPRINT)"

obs:
	@mkdir -p $(BIN_DIR)
	cd tools/obs && go build -o $(OBS) ./cmd/obs

# ── Lint ─────────────────────────────────────────────────────────────

.PHONY: lint check

lint: $(DPRINT)
	$(DPRINT) fmt

check: $(DPRINT)
	$(DPRINT) check
