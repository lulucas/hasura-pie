# Hasura-pie 🥧

A project toolkit for hasura and golang

# Quick start

```
github.com/lulucas/hasura-pie
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
