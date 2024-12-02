### Command for running this program:

Copy file .env.example to .env and edit:

```shell
$ cp .env.example .env
```

After edit file config, build go project:

```shell
go build -v .
```

This command for run migrate database:

```shell
./shipping-service -m=migrate
```

This command for run api :

```shell
./shipping-service
```

### RUN WITH DOCKER

This command build and run container `api` shipping service

```shell
// build process
$ docker build --rm --tag user-service-api:latest -f Dockerfile .
// run
$ docker run --rm -p 8080:8080 --name user-service-api user-service-api:latest
```

```shell
// run docker with local port
$ docker run --rm --net=host -d -p 8080:8080 --name user-service-api user-service-api:latest
```

### Migrations

#### Setup

To install the migrate CLI tool using curl on Linux, you can follow these steps:

```shell
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
## install dependency client library of database
$ go get github.com/golang-migrate/migrate/v4/cmd/migrate
$ go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate # for mysql
$ go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate # for postgres

```

Next, create migration files using the following command:

```shell
$ migrate create -ext sql -dir database/migration/ -seq create_table_card
```

create migration file with timestamp name using the following command:

```shell
$  migrate create -ext sql -dir  database/migration/ -format "20060102150405" add_table_users
```

#### Run Migration Up (All Migrations)

```shell
$ migrate -path database/migration/ -database "postgres://root:root@tcp(localhost:3306)/user_db?multiStatements=true" up
```

#### Run Migration Up Verbose

```shell
$ migrate -path database/migration/ -database "postgres://root:root@tcp(localhost:3306)/user_db?multiStatements=true" -verbose up 1
```

#### Run Migration Down (all)

```shell
$ migrate -path database/migration/ -database "postgres://root:root@tcp(localhost:3306)/user_db?multiStatements=true" down
```

#### Run Migration Down

```shell
$ migrate -path database/migration/ -database "postgres://root:root@tcp(localhost:3306)/user_db?multiStatements=true" -verbose down 1
```

if you want run all migrations up or down versions, you can remove number after `up`/`down` command in below example.

#### Problem

- if migrate error, change value field `dirty` in table `schema_migrations` from 1 to 0
- source documentation https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md

Source: https://www.freecodecamp.org/news/database-migration-golang-migrate/
