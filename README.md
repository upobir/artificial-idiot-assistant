# Artificial Idiot Assistant (AIA)
## Development
### Mongo DB
Use a mongo db for local testing, suggested to use docker. Use following yaml file
```yaml
version: '3.7'
services:
  db:
    image: mongo:latest
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: rootpassword
    ports:
      - 27017:27017
    volumes:
      - mongodb_data_container:/data/db
volumes:
  mongodb_data_container:
```
and run `docker compose up`


### Run
Run with 
```bash
go run ./cmd/aia
```

## Building
Build with
```bash
go build -o ./build ./cmd/aia
```
Binary will be built at `build` folder