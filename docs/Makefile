IMAGE = grafana/grafana-docs-dev:latest
PORT = 3002:3002
CONTENT_PATH = /hugo/content/docs/plugins/caic-datasource/v0.1

.PHONY: pull
pull:
	docker pull $(IMAGE)

.PHONY: docs
docs: pull
	docker run -v $(PWD)/sources:$(CONTENT_PATH) -p $(PORT) --rm -it $(IMAGE) /bin/bash -c "make server-quick"

.PHONY: docs-test
docs-test: pull
	docker run -v $(PWD)/sources:$(CONTENT_PATH) -p $(PORT) --rm $(IMAGE) /bin/bash -c "make prod"
