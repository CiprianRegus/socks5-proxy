build:
	go build -o proxy.exe proxy.go

run:
	go build -o proxy.exe proxy.go  && ./proxy.exe