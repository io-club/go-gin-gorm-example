include .env

test:
	echo DB_NAME: ${DB_NAME}
	echo SERVER_SSH: ${SERVER_SSH}
dev:
	go run .
build:
	go build .
build-static:
	nix-shell shell.nix --run "go build -ldflags \"-linkmode=external -extldflags '-static'\" ."
run:
	./fabric
upload:
	scp ./fibric ${SERVER_SSH}:/root/app/

remote-stop:
	ssh ${SERVER_SSH} "pkill fibric"
remote-deploy-release:
	ssh ${SERVER_SSH} "cd app && MODE=release DB_MODE=release nohup ./fibric > log/fibric.log 2>&1 &"
remote-deploy-debug:
	ssh ${SERVER_SSH} "cd app && MODE=debug DB_MODE=release nohup ./fibric > log/fibric.log 2>&1 &"
remote-clean:
	ssh ${SERVER_SSH} "cd app && rm -f fibric"
build_and_upload : build upload remote-deploy
static-build_upload_release : build-static remote-clean upload remote-stop remote-deploy-release
static-build_upload_debug : build-static remote-clean upload remote-stop remote-deploy-debug
