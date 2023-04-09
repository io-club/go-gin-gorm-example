include .env

test:
	echo DB_NAME: ${DB_NAME}
	echo SERVER_SSH: ${SERVER_SSH}
dev:
	go run .
build:
	go build -ldflags "-linkmode=external -extldflags '-static'" .
run:
	./fabric
upload:
	scp ./fibric ${SERVER_SSH}:/root/app/

remote-stop:
	ssh ${SERVER_SSH} "pkill fibric"
remote-deploy:
	ssh ${SERVER_SSH} "cd app && nohup ./fibric > log/fibric.log 2>&1 &"
remote-clean:
	ssh ${SERVER_SSH} "cd app && rm -f fibric"
build_and_upload : build upload remote-deploy
