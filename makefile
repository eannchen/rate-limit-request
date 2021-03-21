run:
	go run main.go

test:
	go test -v

dkrps-prune:
	docker image prune --force

dkrps-up:
	docker-compose up -d --build

dkrps-down: dkrps-prune
	docker-compose down --remove-orphans --volumes