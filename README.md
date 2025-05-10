# Go Gin Boilerplate

> A starter project with Golang, Gin and postgresql

Golang Gin boilerplate with postgresql resource. Supports multiple configuration environments.

![](header.jpg)

## Boilerplate structure

```
.
├── Makefile
├── Procfile
├── README.md
├── build
├── config
│   ├── config.go
│   ├── dev.yaml
│   ├── prod.yaml
│   └── stage.yaml
├── controllers
│   └── health.go
├── data
│   └── dummy.json
├── db
│   └── db.go
├── logger
│   └── logger.go
├── header.jpg
├── main.go
├── middlewares
│   └── auth.go
├── schema
│   └── db.go
├── util
│   └── http.go
└── server
    ├── router.go
    └── server.go
```

## Installation

```sh
make deps
```

## Usage example

`curl http://localhost:3000/health`
