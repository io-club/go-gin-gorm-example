
build:
	go build -ldflags "-linkmode=external -extldflags '-static'" .
run:
	./fabric
dev:
	go run .
upload:
	scp ./fibric root@101.42.21.155:/root/app/
build_and_upload : build upload