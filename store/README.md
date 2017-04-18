# LURK store

The LURK store records all registrations

Registrations have the following layout:
- id
- status (new, work-in-progress, done, failed)
- creation date
- finalisation date
- expiration date	(who decides this: LURK or ACME?)
- CSR
- lifetime
- STAR crt URL

## lurk server -> store
- save new registration
  - in: csr, lifetime
  - out: record id

- fetch registration by id
  - in: id
  - out: the complete registration record

## acme client -> store
- fetch oldest pending
  - in: nothing
  - out: a complete registration record
  - side-effect: registration moves to state work-in-progress

- finalise success
  - in: id, STAR crt URL, lifetime, visibility ttl
  - out: nothing

- finalise failure
  - in: id, visibility ttl
  - out: nothing

# SQLite
```
$ sqlite
sqlite> .mode column
sqlite> .headers on
sqlite> SELECT * FROM registration;
```

## go-sqlite3 dependency
```
go get -u github.com/mattn/go-sqlite3
```
