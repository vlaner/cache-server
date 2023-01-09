# Simple Cache Server
### Small TCP server to cache some data

# Run
```go
go build cmd/app/main.go
./main

or

go run cmd/app/main.go
```

# Commands
* SET -**KEY**- -**VALUE**- — save desired data

* GET -**KEY**- — retrieve data by key

* DEL -**KEY**-  — delete data by key

* EXPIRE -**KEY**- -**DURATION**-  — set expiration time of key in seconds

# Examples

1. SET **hello** **"hello, world!"** — maps KEY ```"hello"``` to  ```"hello, world!" ```
2. GET **hello** — will print ```"hello, world!"```
3. DEL **hello** — deletes KEY ```"hello"```
4. EXPIRE **hello** 10 — sets expiration time of ```"hello"``` to ```10``` seconds