# MP API

## Initialize
`go mod init github.com/letthefireflieslive/mp-api`

## Contribute
`go mod download`

## Run
`go run main.go`

## Testing

### Unit Tests
`go test`

### E2E
Load and run [postman collection](https://www.getpostman.com/collections/de9ac6fa670ad3fc7ce3) and 
[environment](https://hammerden.postman.co/workspace/hammerden~20b843e8-ff70-4051-8b09-2a779a657145/environment/23681075-bc2f01e6-fe9f-4e62-a238-4a7e813df886)

## Run database
Create folder

`mkdir ./mongodb_data`

Run the container

`docker run -d --name mongodb -v $(pwd)/mongodb_data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:6.0.2`

## Documentation
```
swagger generate spec -o ./swagger.json
swagger serve ./swagger.json
```

