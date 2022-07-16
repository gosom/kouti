default:
	echo "use one command"

docs:
	swag init -g rest/routes.go --parseDependency

generate:
	sqlc generate

dev:
	go run main.go

dummy:
	go run cmd/hydrator/main.go

#
# Some dummy requests
#

create-user:
	- curl -i -XPOST http://localhost:8080/api/v1/users -d '{"email": "g+1@gkomninos.com", "password": "what1!#ok","firstName": "Avraam", "lastName": "Papadopoulos"}'

get-users:
	 - curl -i -XGET 'http://localhost:8080/api/v1/admin/users/?pageSize=10'

get-user:
	 - curl -i -XGET http://localhost:8080/api/v1/users/2

delete-user:
	 - curl -i -XDELETE http://localhost:8080/api/v1/users/1
	
search-users:
	 - curl -i -XGET 'http://localhost:8080/api/v1/users/search?pageSize=10&searchTerm=Justenjusten%20konn'

