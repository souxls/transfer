# Transfer

`Transfer` 主要用于临时上传和下载，支持分享给他人。此项目是`Transfer`项目后台 API.

## Dependent Tools

```bash
go get -u github.com/cosmtrek/air
go get -u github.com/google/wire/cmd/wire
go get -u github.com/swaggo/swag/cmd/swag
```

- [air](https://github.com/cosmtrek/air) -- Live reload for Go apps
- [wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go
- [swag](https://github.com/swaggo/swag) -- Automatically generate RESTful API documentation with Swagger 2.0 for Go.

## Dependent Library

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [GORM](https://gorm.io/) -- The fantastic ORM library for Golang
- [Casbin](https://casbin.org/) -- An authorization library that supports access control models like ACL, RBAC, ABAC in Golang
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## Build

```bash
cd transfer
go build
```

## Usage

```bash
transfer -h
transfer file to minio API

Usage:
  transfer [flags]

Flags:
      --config string   config file (default is $HOME/.transfer.yaml)
  -h, --help            help for transfer
  -t, --toggle          Help message for toggle
```
