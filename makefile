run:
	go run main.go

test:
	go test

dkrps-prune:
	docker image prune --force

dkrps-up:
	docker-compose up -d --build

dkrps-down: dkrps-prune
	docker-compose down --remove-orphans --volumes