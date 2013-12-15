# Cortado

A URL shortener written in Go and base on the [Traffic Web Framework](https://github.com/pilu/traffic).

<http://♡♡♡.ws>

```bash
$ curl -XPOST -i http://localhost:7000/ -d "url=http://gravityblast.com"
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Cortado-Version: 0.1.0
Location: http://localhost:7000/1

{"long_url":"http://gravityblast.com","short_url":"http://localhost:7000/1","shorty":"1"}
```

```bash
$ curl -i http://localhost:7000/1
HTTP/1.1 301 Moved Permanently
Content-Type: application/json; charset=utf-8
Cortado-Version: 0.1.0
Location: http://gravityblast.com

```

```bash
curl -i http://localhost:7000/1+
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Cortado-Version: 0.1.0

{"clicks":"30","long_url":"http://gravityblast.com","short_url":"http://localhost:7000/1","shorty":"1"}
```

## Author

* [Andrea Franz](http://gravityblast.com)

