install:
	dep ensure -v

build:
	@echo [ building ]
	go build -o=./bin/run ./main.go
	@echo [ successfully built into "bin" folder with name "run" ]