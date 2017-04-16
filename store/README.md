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

- finalise
 - in: id, status, STAR crt URL (opt, if !failed), lifetime (opt, if !failed)
 - out: nothing

# SQLite
```
$ sqlite
sqlite> .mode column
sqlite> .headers on
sqlite> .read registration.sql
[sqlite> .read populate_test.sql]
```

## go-sqlite3 dependency
```
go get -u github.com/mattn/go-sqlite3
```
