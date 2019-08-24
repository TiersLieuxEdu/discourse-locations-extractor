all: build

build:
		cd cmd/forum2csv && go build
		cd cmd/forum2geojson && go build
