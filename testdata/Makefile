all: test

test:
	curl -X POST -d '404' -i http://127.0.0.1/health
	curl --fail-with-body -i 127.0.0.1/health || exit 0 && exit 1
	curl --fail-with-body -i 127.0.0.1/health || exit 0 && exit 1
	curl --fail-with-body -i http://127.0.0.1
	curl --fail-with-body -i http://127.0.0.1 || exit 0 && exit 1

.PHONY: all
