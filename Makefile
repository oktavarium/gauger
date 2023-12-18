run:
	docker-compose up --build --abort-on-container-exit
cover:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out
cleen:
	rm cover.out
