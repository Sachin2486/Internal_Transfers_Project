# Internal Transfers System -: 
An internal transfers application that facilitates financial transactions between accounts with HTTP endpoints for account management and balance queries.

## Dependencies:

- Go Version >= 1.24.1
- PostgreSQL >= 14

## Setup:
- Modify below config.yaml params to your DB credentials
```
db:
host: localhost
user : admin
password: admin
name: transfers
port : 5432 
```
- Run the application using `go run main.go`
