# Bitburst Assessment Task

Write a rest-service that listens on localhost:9090 for POST requests on /callback.

Run the go service attached to this task. It will send requests to your service at a fixed interval of 5 seconds.

The request body will look like this:

```json
{
 "object_ids": [1,2,3,4,5,6]
}
```

The amount of IDs varies with each request. Expect up to 200 IDs.

Every ID is linked to an object whose details can be fetched from the provided
service. Our service listens on localhost:9010/objects/:id and returns the
following response:

```json
{
 "id": <id>,
 "online": true|false
}
```

Note that this endpoint has an unpredictable response time between 300ms and 4s!
* Your task is to request the object information for every incoming object_id and filter the objects by their "online" status.
* Store all objects in a PostgreSQL database along with a timestamp when the object was last seen.
* Let your service delete objects in the database when they have not been received for more than 30 seconds.

Important: due to business constraints, we are not allowed to miss any callback to our service.

Write code so that all errors are properly recovered and that the endpoint is always available.

Optimize for very high throughput so that this service could work in production.

Bonus:
* some comments in the code to explain the more complicated parts are appreciated
* it is a nice bonus if you provide some way to set up the things needed for us to

Test your code.

# Tests & Behcmarks

**NOTE:** to run tests and benchmarks Docker must be installed on system, because each test and benchmark is run in isolation using [ory/dockertest]("https://github.com/ory/dockertest") package

To run the tests run **go** command:
```bash
go test -v ./...

# print coverage as well
go test -v -coverprofile cover.out ./...
go tool cover -html cover.out
```

To run benchmarks run **go** command, it can take up to 1 minute or more:
```bash
go test -cpu=1,2,4,8 -benchmem -run=^$ -bench . ./...
```

# Configuration

You can tweek configuration from command flags, configuration file(.yaml) or environmental variables. Simply run `./bitburst --help` to see all available flags, or create a file with _yaml_ extension and use [example.yaml](config/example.yaml) as example, then you can pass it to program using `--config-path` flag. If you prefer using env vars, I suggest to download and install [direnv]("https://direnv.net"), list of envs:

* **$BITBURST_SERVER_LISTEN_ADDRESS** - listen address for http server, port must be included (default: 0.0.0.0:9090)
* **$BITBURST_CLIENT_TESTER_SERVICE_ADDRESS** - listen address of tester service (default: 127.0.0.1:9010)
* **$BITBURST_DATABASE_HOST** - address host of postgres db (default: 127.0.0.1)
* **$BITBURST_DATABASE_PORT** - address port of postgres db (default: 5432)
* **$BITBURST_DATABASE_USERNAME** - username of postgres db (default: postgres)
* **$BITBURST_DATABASE_PASSWORD** - password of postgres db (default: postgres)
* **$BITBURST_DATABASE_NAME** - database name of postgres db (default: postgres)

# Build

In order to build the application, you only need Go installed. Example:
```
go build -o ./bin/bitburst ./cmd/bitburst/main.go
```


On the other hand, if you wish to use Docker for trying out the application, you can use **docker-compose.yaml** file. Example:
```
docker-compose up

or

docker compose up
```

# Notes

In a production environment metrics, tracing and error reporting are usually set for the service, but for the sake of brewity I didn't include them in application. Also, one might think that for such small task all this packages are overkill, but in reality when writing production grade services, code complexity grows, more features are added, and so this packages help to eliminate boilerplate coding and follow DRY principle, and even improve performance

For configuration management I used [spf13/viper]("https://github.com/spf13/viper") package, as it's the best solution available for Go configuration management

For postgres driver I used [jackc/pgx]("https://github.com/jackc/pgx") package, as it's better than std database/sql package and faster than jmoiron/sqlx package

Instead of using raw sql queries, I prefer to use Go type safe generator [kyleconroy/sqlc]("https://github.com/kyleconroy/sqlc"), it's much more convenient and saves you from having typos in your queries. Besided that, it works like a charm with migrations and [jackc/pgx]("https://github.com/jackc/pgx") package

For logging I used [rs/zerolog]("https://github.com/rs/zerolog"), as it's zero allocations logger, I could use uber's zap logging package as well, but zerolog seems better for me

For tests I used [stretchr/testify]("https://github.com/stretchr/testify") package, as it provides greate set of functions that removed useless boilerplates in form of `if err != nil` and also it provides convenient logging messages. I also used [gocmp/cmp]("https://github.com/google/go-cmp") in tests where difference between input and output was needed to be shown

For working with json I used [json-iterator/go]("github.com/json-iterator/go") package, as it's much more faster than stdlib json package

For testing I used [ory/dockertest]("https://github.com/ory/dockertest"), because running each database test inside Docker mock container is really convenient