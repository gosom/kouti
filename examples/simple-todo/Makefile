default:
	echo "use one command"

run:
	go run main.go

swaggen:
	swag init -g rest/routes.go --parseDependency

enter-docker-db:
	docker exec -it simple-todo_db_1 psql -U postgres todo
