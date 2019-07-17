# Package configuration
PROJECT = sourced-ce
COMMANDS = cmd/sourced
PKG_OS ?= darwin linux windows

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_PATH ?= $(shell pwd)/.ci
CI_VERSION ?= v1

MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	git clone --quiet --branch $(CI_VERSION) --depth 1 $(CI_REPOSITORY) $(CI_PATH);

-include $(MAKEFILE)

GOTEST_BASE = go test -v -timeout 20m -parallel 1 -count 1 -ldflags "$(LD_FLAGS)"
GOTEST_INTEGRATION = $(GOTEST_BASE) -tags="forceposix integration"

OS := $(shell uname)

# override clean target from CI to avoid executing `go clean`
# see https://github.com/src-d/sourced-ce/pull/154
clean:
	rm -rf $(BUILD_PATH) $(BIN_PATH) $(VENDOR_PATH)

ifeq ($(OS),Darwin)
test-integration-clean:
	$(eval TMPDIR_INTEGRATION_TEST := $(PWD)/integration-test-tmp)
	$(eval GOTEST_INTEGRATION := TMPDIR=$(TMPDIR_INTEGRATION_TEST) $(GOTEST_INTEGRATION))
	rm -rf $(TMPDIR_INTEGRATION_TEST)
	mkdir $(TMPDIR_INTEGRATION_TEST)
else
test-integration-clean:
endif

test-integration-no-build: test-integration-clean
	$(GOTEST_INTEGRATION) github.com/src-d/sourced-ce/test/

test-integration: clean build test-integration-no-build
