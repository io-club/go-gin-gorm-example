
build:
	go build -ldflags "-linkmode=external -extldflags '-static'" .
run:
	./fabric
send:
	scp ./fibric root@101.42.21.155:/root/app/