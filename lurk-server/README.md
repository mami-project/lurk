# Building
```
$ go get
$ go build
```

# Testing
```
$ brew install httpie
```

## Create registration request
```
$ http --verbose POST :8080/star/registration csr=abcd lifetime:=1234
```

## List pending registrations
```
$ http --verbose GET :8080/star/registrations
```

## Poll registration by id
```
$ http --verbose GET :8080/star/registration/1
```
