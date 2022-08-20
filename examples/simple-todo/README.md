# Quickstart

Start the database via docker-compose

```
docker-compose up -d
```

build and run the server

Install swag `go install github.com/swaggo/swag/cmd/swag@latest`

```
go generate # generates swagger documentation
go build
./simple-todo
```

visit docs here: http://localhost:8080/docs/

Make some requests

```
curl -XPOST http://localhost:8080/todos -H 'Content-Type: application/json' -d '{"content": "todo1"}'

curl http://localhost:8080/todos

curl -XGET http://localhost:8080/todos/c1addf48-e20a-4dbd-932a-ce8def0aabc6 -H 'Content-Type: application/json'

curl -XDELETE http://localhost:8080/todos/dcd0c640-293a-41f7-b431-b5a4feb1f9fc -H 'Content-Type: application/json'


curl -XPUT http://localhost:8080/todos/c1addf48-e20a-4dbd-932a-ce8def0aabc6 -H 'Content-Type: application/json' -d '{"is_completed": true, "content": "todo1"}'

curl http://localhost:8080/todos?content=todo

curl http://localhost:8080/todos?content=random

curl http://localhost:8080/todos?is_completed=t

curl 'http://localhost:8080/todos?is_completed=t&content=1'
```
