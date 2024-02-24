curl -d '{"id":7, "name":"Kirill"}' -H "Content-Type: application/json" -X POST http://localhost:8080/users/add
curl -d '{"id":2, "name":"Kirill"}' -H "Content-Type: application/json" -X PUT http://localhost:8080/users/update

curl http://localhost:8080/users