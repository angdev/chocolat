Chocolat
========

## What is Chocolat?

Chocolat is a Keen.io-compatible API server for Collecting Event Data, and Analyzing. Most of data analysis services don't allow to use event data directly. Bring event data of your service.

### Why Keen.io?

[Keen.io](https://keen.io/) is awesome, and developer-friendly service. SDK is well-documented, and [a lot of open sources](https://github.com/keen) are provided. By building keen.io-compatible API server, we can use plenty of resources. If you need robust support (stability, scalability, etc), then use Keen.io.

### Status

Chocolat is in active development.

Currently, it just works at local development environment. It uses MongoDB as event data store, and aggregation framework. So far, I've concentrated on API specification, so chocolat just like skeleton. From now, Data collecting and Aggregation feature should aim high performance, and scalability.


## Getting Started

### Prerequisites

* Go 1.4+
* MongoDB 2.2+
* Sqlite3 (development)
* Mysql, Postgres (production)

### How to run

Clone this repository:

    git clone https://github.com/angdev/chocolat.git

Install go package dependencies:

    go get . && go get bitbucket.org/liamstask/goose/cmd/goose

Initialize a database:

    goose up

Build and Run:

    go build && ./chocolat -r

Before you run chocolat, check out `config/dbconf.yml`, `config/repoconf.yml`.

### Using with Keen.io SDK

Create a new project:

    ./chocolat -c

Then chocolat prints `PROJECT_ID`, `MASTER_KEY`, `READ_KEY`, `WRITE_KEY`.

And you can use chocolat with keen.io sdk like below (JS SDK example):

```javascript
<script src="//cdn.jsdelivr.net/keen.js/3.2.5/keen.min.js" type="text/javascript"></script>
<script type="text/javascript">
  var client = new Keen({
    projectId: "{{PROJECT_ID}}",
    writeKey: "{{WRITE_KEY}}",
    readKey: "{{READ_KEY}}",
    protocol: "http",
    host: "localhost:3000",
  });
</script>
```

## Roadmap

Chocolat will be multi-components service. Even if chocolat is just api server now, but a lot of components (event aggregation service, message queue, ...) will be built.

### Backlog

* Testing (Unit test, Integration test)
* Mass event collecting
* Fast event aggregation
* Worker (e.g. Daily Report)
* Dockerize

## License

MIT License. See LICENSE for details.

## Credits

* [Keen.io](https://keen.io/)
* Keen.io open sources contributors