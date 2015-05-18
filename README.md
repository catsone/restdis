# Restdis

Restdis is a small web service that provides an HTTP interface to Redis. It's very similar to
[Webdis](https://github.com/nicolasff/webdis), except without all the fancy features.

## Installing

```bash
go get github.com/catsone/restdis
```

## Usage

The interface is pretty much the same as Webdis. You can use GET or POST interchangeably.

```bash
curl localhost:7631/SET/foo/bar
→ {"response":"OK"}

curl localhost:7631/GET/foo
→ {"response":"bar"}
```
