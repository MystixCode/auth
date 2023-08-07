# auth

auth is a go api

It's a container with the api and a container with a postgresql database.


<!-- <img src="https:///.png" width="100%" height="100%"> -->


## run compose in dev mode

> in dev mode the air package is used for hot reload

```bash
docker-compose -f docker-compose.dev.yml up --build
```

## run with compose

```bash
docker-compose up -d --build
```
show logs:
```bash
docker-compose logs -f auth
docker-compose logs -f postgres
```

## run with docker

```bash
docker build -t auth .
docker run -p 8080:8080 auth
```

## migrate db

```bash
# in config.yaml change database.host to localhost

# in docker-compose.yml add port mapping to postgres
# ports:
    #   - 5432:5432

# Build and start Database
docker-compose up --build -d postgres

# Run migrate cli command
go run auth migrate
```
