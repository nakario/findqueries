# findqueries

`findqueries` find SQL queries from .go files in a package and report them in json format.

It also reports call graph of every functions in the package specified.

## Usage

```bash
$ go install github.com/nakario/findqueries@latest
$ findqueries . 1> queries.json
```

`findqueries` outputs json to stdout and outputs errors and found SQL queries to stderr.

```bash
$ cat queries.json | jq .
# {
#   "name": "main",
#   "queries": [
#     {
#       "query": "SELECT id FROM users",
#       "caller": "getID",
#       "expr": "db.Query(query)",
#       "pos": "/path/to/main.go:24:10"
#     }, ...
#   ],
#   "calls": [
#     {
#       "caller": "getID",
#       "callee": "(*database/sql.DB).Query"
#     }, ...
#   ]
# }
```

## Test this package

To run `go test`, you need to run the following command.

```bash
$ go mod download
$ go test
```
