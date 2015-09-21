Chocolat
========

## What is Chocolat?

Chocolat is a Keen.io-compatible API server for Collecting Event Data, and Analyzing. Most of data analysis services don't allow to use event data directly. Bring event data of your service.

### Why Keen.io?

[Keen.io](https://keen.io/) is awesome, and developer-friendly service. SDK is well-documented, and [a lot of open sources](https://github.com/keen) are provided. By building keen.io-compatible API server, we can use plenty of resources. If you need robust support (stability, scalability, etc), then use Keen.io.

## Demo

[Demo page](http://chocolat-demo.angdev.io/explorer)

![demo](http://i.imgur.com/rD4VyBo.png)

This demo includes [keen-js](https://github.com/keen/keen-js), [Data Explorer](https://github.com/keen/explorer). Try to use keen.js client, and explore metrics below.

## Getting Started

### Prerequisites

* Go 1.4+
* MongoDB 2.2+ (See Status)
* Sqlite3 (development)
* Mysql, Postgres (production)

### Quick Run

Before running chocolat, make sure that mongo is running (localhost). Check out `config/dbconf.yml`, `config/repoconf.yml`.

Clone this repository:

    git clone https://github.com/angdev/chocolat.git

Install go package dependencies:

    go get .

or if you use godep,

    godep restore

Build and Run:

    go build && ./chocolat -r

If you want to install:

    go install

but chocolat should be executed at chocolat directory, because of configuration files.

## Configuration

Database (Mysql, Postgres, ...) configuration is at `config/dbdonf.yml`. Repository (MongoDB, for event collections) configuration is at `config/repoconf.yml`. You can set config values directly, or using environment variable. (YAML files are preprocessed by [golang template package](http://golang.org/pkg/text/template/))

### Environment Variables

`CHOCOLAT_PORT=3000`, `CHOCOLAT_ENV=development`, `REPO_URL`, `DATABASE_DRIVER`, `DATABASE_URL` are configurable. [web.environment @ docker-compose.yml](https://github.com/angdev/chocolat/blob/develop/docker-compose.yml#L9) will be a help to you.

## Simple CLI Usage

### Create a project

    ./chocolat -c

### List projects

    ./chocolat -l

### Show API Keys

    ./chocolat -p {{.projectID}}

### Run Server

    ./chocolat -r

### Print chocolat help (All commands usage)

    ./chocolat -h

## Play with Keen Open Source

### Keen.io SDK

Create a new project:

    ./chocolat -c

Then chocolat prints `PROJECT_ID`, `MASTER_KEY`, `READ_KEY`, `WRITE_KEY`.

And you can use chocolat with keen.io sdk like below (Javascript client example):

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

If want to use with rails, then just use [Ruby Client](https://github.com/keenlabs/keen-gem).

### [Data Explorer](https://github.com/keen/explorer)

Follow instructions: [Explorer Quick Setup](https://github.com/keen/explorer#tldr)

and then just add chocolat project id, api keys!

## Docker (Experimental)

Build chocolat image:

    docker build -t chocolat .

You can test chocolat by using docker compose:

    docker-compose up

## Status

Chocolat is in active development.

Currently, it just works at local development environment. It uses MongoDB as event data store to implement aggregation features fast, MongoDB may be changed to another database. So far, I've concentrated on API specification, so chocolat just likes skeleton. From now, Data collecting and Aggregation feature should aim high performance, and scalability.

## Roadmap

Chocolat will be multi-components service. Even if chocolat is just api server now, but a lot of components (event aggregation service, message queue, ...) will be built.

### Backlog

* Project API
* Multi-Analysis API
* Funnel API
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
