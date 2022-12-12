# MP API

## Initialize
`go mod init github.com/letthefireflieslive/mp-api`

## Contribute
`go mod download`

## Testing

### E2E
Load and run [postman collection](https://www.getpostman.com/collections/de9ac6fa670ad3fc7ce3) and 
[environment](https://hammerden.postman.co/workspace/hammerden~20b843e8-ff70-4051-8b09-2a779a657145/environment/23681075-bc2f01e6-fe9f-4e62-a238-4a7e813df886)

## Run database
Create folder

`mkdir ./mongodb_data`

## Local: Run mongodb container
`docker run -d --name mongodb -v $(pwd)/mongodb_data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:5.0.14`

## Local: Run redis container
`docker run -d -v $(pwd)/conf/redis:/usr/local/etc/redis --name redis -p 6379:6379 redis:6.0`

## Compass connection string
`mongodb://admin:password@localhost:27017`

## Set env variables
```
export MONGO_DATABASE=mealPlan
export MONGO_URI="mongodb://admin:password@0.0.0.0:27017/mealPlan?authSource=admin"
```
> Make sure that your IDE can read this values

## Optional: Import data
`mongoimport --username admin --password password --authenticationDatabase admin --db demo --collection recipes --file meal_plans.json --jsonArray`

## Run the application 
`go run main.go`

## Documentation
```
swagger generate spec -o ./swagger.json
swagger serve ./swagger.json
```

