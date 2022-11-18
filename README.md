# Password Manager
This is the code for my Computer Science A Level project. The password manager is written in go (golang) and serves as a webserver with an API and Client to interact with the API. I am using `net/http` to run the http server and I am using `html/template` for all the HTML. Both of these libraries are from the robust standard library. 

# Direct External Depencencies
I am using `gorm` 
- `gorm` which is an awesome database orm, this project is using sqlite3 however gorm supports many more databases. 
- `github.com/google/uuid` which allows me to use UUID's. I use uuid's as the id's for my database.
- `golang.org/x/crypto` which is only used so I can use the pbkdf2 key derivation function.

# Building
Building a go application is very simple. The server compiles into a single executable that embeds everything in src/public/* and src/frontend/html/*
```
go build src/server.go
```

### Running
Instead of building, go allows you temporarily compile and then run.
```
go run src/server.go
```

# Testing
Due to go's fantastic testing library from the standard library, you can run go tests very easily.

### Test Backend
```
go test src/backend/tests/...
```

### Test Frontend
```
go test src/frontend/tests/...
```