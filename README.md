# Hasura-pie 🥧

A project toolkit for hasura and golang (Unstable)

# Quick start

Install CLI

```
go get github.com/lulucas/hasura-pie-cli
```

# Project structure

```
apps
modules
models
```

# CLI

## Command

Generate module

```
pie g m account
```

Sync action from hasura

## Config

Create config.yml in cli work directory

```
postgres:
  host: 127.0.0.1
  user: postgres
  pass: postgres
hasura:
  endpoint: http://localhost:8080
  admin_key: 123
```

# Related project

* https://github.com/ekhabarov/go-pg-generator
