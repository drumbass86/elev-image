# Testing
For testing need start the docker cointainer with PostgreSQL and tested data:
```
./tools/run_test_postgres_docker.sh
```
This script run container - **test-postgres** available on port **5433** and mount needed data in temp directory.
When the container start it execute a mounted script - `./tools/docker-entrypoint-initdb.d/init-db-data.sh` which
create db - *test* (watch `db/create_db.sql`) and fill it data (wathc `db/data.sql`).

Now we can do unit testing cmd
```
go test ./...
```
And another test.


## Test API
To test server API use a utility `curl`.
Commands for testing:
```
curl -v http://localhost:8081/v1/path
curl -v http://localhost:8081/v1/img/1
```

To test API method - *create captured image*, use the *.json file with body request, which is located in the folder `test`.
Command for testing API method - create captured image:
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data @./test/capimg_1.json \
  http://localhost:8081/v1/img/create
```
In a response body should be the created captured image.