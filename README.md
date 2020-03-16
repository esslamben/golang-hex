# golang-hex
A brief example of a golang project using a hexagonal architecture design.

## Installation
You can run the app using normal ```go run``` command or ```go build```. It uses 2 environment variables but both have a default
 
 - MONGURI: mongodb://mongo:27017
 - PORT: 5000
 
For the above to work you will need a mongdb instance running on the shown port.
 
## Run
I've took the liberty of creating a docker-compose file that builds the app into a binary and runs a mongo instance that it has access to.

To start the server just run ``docker-compose up --build``

The API will be available on ```http://localhost:5000/user```

## Available endpoints
- POST http://localhost:5000/user   which takes some json with name and returns user
- GET http://localhost:5000/user/{uuid} which returns a user for the given uuid
- PUT http://localhost:5000/user/{uuid} which uses the given json to update the user
- DELETE http://localhost:5000/user/{uuid} which deletes a user for the given uuid