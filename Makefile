run: build
	./bin/app -data-path ./data/data -http :3000

build: deps gen
	@go build -o ./bin/app

docker-run: docker-build
	docker run --name=hlc --rm -v `pwd`/data/data:/tmp/hlc/data -p 3000:80 hlc

docker-build: build-binary
	docker build -f ./docker/Dockerfile.run -t hlc . >/dev/null

build-binary: build-builder
	docker run --rm -v `pwd`:/go/src/hlc -w /go/src/hlc hlc_builder make build-static

build-static: deps gen
	go build -ldflags "-linkmode external -extldflags -static" -o ./bin/app .

gen:
	easyjson ./import.go ./models/user.go ./models/location.go ./models/visit.go

deps:
	@go get github.com/mailru/easyjson/...
	@dep ensure

build-builder:
	docker build -f ./docker/Dockerfile.build -t hlc_builder . >/dev/null

release: docker-build
	docker tag hlc stor.highloadcup.ru/travels/chinook_buildero
	docker push stor.highloadcup.ru/travels/chinook_buildero

release-volegov: docker-build
	docker tag hlc stor.highloadcup.ru/travels/utopian_dolphin
	docker push stor.highloadcup.ru/travels/utopian_dolphin

data-unpack-train: clean
	@if [ ! -d "./data" ]; then \
		echo "unpacking ./data_train.zip to ./data"; \
		unzip -q ./data_train.zip -d ./data; \
	fi

data-unpack-full: clean
	@if [ ! -d "./data" ]; then \
		echo "unpacking ./data_full.zip to ./data"; \
		unzip -q ./data_full.zip -d ./data; \
	fi

TESTER_BIN=highloadcup_tester

$(TESTER_BIN):
	@go get github.com/AterCattus/highloadcup_tester

tester: tester-1 tester-2 tester-3

tester-1: $(TESTER_BIN)
	@$(TESTER_BIN) -addr http://127.0.0.1:3000 -hlcupdocs ./data -test -phase 1

tester-2: $(TESTER_BIN)
	@$(TESTER_BIN) -addr http://127.0.0.1:3000 -hlcupdocs ./data -test -phase 2

tester-3: $(TESTER_BIN)
	@$(TESTER_BIN) -addr http://127.0.0.1:3000 -hlcupdocs ./data -test -phase 3

clean:
	rm -rf ./bin ./data