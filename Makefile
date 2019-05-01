SUPERSET_REPO = https://github.com/apache/incubator-superset.git
SUPERSET_VERSION = release--0.32
SUPERSET_DIR = superset
PATCH_SOURCE_DIR = srcd
ADD_FILES = superset/superset_config.py superset/bblfsh superset/assets/src/uast
OVERRIDE_FILES = \
	superset/assets/package.json \
	superset/assets/package-lock.json \
	superset/assets/webpack.config.js

all: superset patch

# Clone superset repository
superset:
	git clone --quiet --branch $(SUPERSET_VERSION) $(SUPERSET_REPO) $(SUPERSET_DIR)

# Overrides files in the superset repository
.PHONY: patch
patch:
	@for file in $(ADD_FILES) $(OVERRIDE_FILES); do \
		echo "patching $${file}"; \
		rm -rf "$(SUPERSET_DIR)/$${file}"; \
		cp -r "$(PATCH_SOURCE_DIR)/$${file}" "$(SUPERSET_DIR)/$${file}"; \
	done; \

# Overrides files in the superset repository using symlinks
.PHONY: patch-dev
patch-dev:
	@for file in $(ADD_FILES) $(OVERRIDE_FILES); do \
		echo "patching $${file}"; \
		rm -rf "$(SUPERSET_DIR)/$${file}"; \
		ln -s "$(PATCH_SOURCE_DIR)/$${file}" "$(SUPERSET_DIR)/$${file}"; \
	done; \

# Create docker image
.PHONY: build
build: superset patch
	docker build -t smacker/superset:demo-with-bblfsh -f docker/Dockerfile .
