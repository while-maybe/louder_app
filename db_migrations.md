# go-migrate quick ref

## Install the go-migrate CLI tool

```bash
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

(might take a while to run, be patient)

## make it a part of path

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc && source ~/.bashrc
```

## Create the migrations folder at the project root folder

`mkdir -p migrations`

## Create the initial users table

`migrate create -ext .sql -dir ./migrations create_users_table`

this results in 2 files being created:

```bash
./migrations/20250607165547_create_users_table.up.sql
./migrations/20250607165547_create_users_table.down.sql
```

## Edit both files as follows

```bash
#./migrations/20250607165547_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS person (
    id BLOB(16) PRIMARY KEY,
    first_name VARCHAR(40) NOT NULL,
    last_name VARCHAR(40) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    dob DATETIME NOT NULL
);

#./migrations/20250607165547_create_users_table.down.sql
DROP TABLE person;
```

## Start the migrations at the project root

`migrate -database "sqlite3://louder.db" -path ./migrations up`

## Revert the migration with

`migrate -database "sqlite3://louder.db" -path ./migrations down`
