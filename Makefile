SUPERSET_REPO = https://github.com/apache/incubator-superset.git

superset:
	git clone --quiet $(SUPERSET_REPO) superset

build: superset
	docker build -t smacker/superset:demo -f docker/Dockerfile .
