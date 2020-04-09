# Hasura-pie 🥧

A project toolkit for hasura and golang standalone service (Unstable)

# Quick start

Install [Hasura-pie-cli](https://github.com/lulucas/hasura-pie-cli)

```
go get -u github.com/lulucas/hasura-pie-cli/pie
```

Initialize project

```
pie init myproject
```

# Project structure

```
app
  business # business app for hasura actions, events or normal http request.
    - main.go
    - Dockerfile
    - docker-compose.yml
    - docker-compose.prod.yml
    - .env
    - .env.prod
module
model
  - model_gen.go # models sync from postgres
```

