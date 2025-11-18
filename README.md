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

 Assumptions:
1. All accounts use the same currency (no currency conversion needed).
2. No authentication or authorization is implemented.
3. Account IDs are provided by the client during creation.
4. Decimal precision supports up to 5 decimal places (format: `decimal(10,5)`).
5. Initial balance must be non-negative.
6. Transfer amount must be greater than zero.
7. Self-transfers (same source and destination) are not allowed.
8. Both source and destination accounts must exist before a transfer.
9. Concurrent transactions are handled using database-level locking.
