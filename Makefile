.DEFAULT_GOAL = help
SRC_PATH = ./cmd/pod-metrics-exporter/

help:
	@echo "----------------------- HELP ------------------------"
	@echo "This application requires go1.16.5 installed.        "
	@echo "                                                     "
	@echo "To run the unit tests type make unittest             "
	@echo "To run the app type make run                         "
	@echo "To build a binary type make build                    "
	@echo "-----------------------------------------------------"

unittest:
	cd $(SRC_PATH) ; \
		go test .

build:
	cd $(SRC_PATH) ; \
	    go get -d -v ; \
		go build . 

run:
	cd $(SRC_PATH) ; \
		go run . --label-name app --label-value demo-db