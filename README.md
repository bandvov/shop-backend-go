# Shop backend written in Golang

## Running

to start application add `.env` file with content:


```
CONN_STRING=<your connection string>
PORT=<port>
```


 run `docker compose up`

## Coverage

to generate test coverage run `go test -coverprofile cover.out` 
to see coverage in the browser run `go tool cover -html=cover.out`