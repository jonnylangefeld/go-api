## Go API

This repository contains all best practices for writing a go API that I learned over the years. It is meant to be cloned
and renamed to get a headstart for building out your api. Or just to be pulled up while working on another project to
find inspiration.

This repo uitlizes

* Command line flags using [pflag](https://github.com/spf13/pflag)
* Structured logging using [zap](https://github.com/uber-go/zap)
* HTTP framework using [chi](https://github.com/go-chi/chi)
* Database integration using [gorm](https://github.com/go-gorm/gorm)
* Database integration tests using [dockertest](https://github.com/ory/dockertest)
* API integration tests using [gomock](https://github.com/golang/mock)
* Documentation as code using [http-swagger](https://github.com/swaggo/http-swagger)

And follows the following best practices:

* Log version on startup
* Graceful exit
* Staged Dockerfile
* Custom middlewares
* Swagger docs

### Run

To run this api you need a postgres instance. You can set the connection string via the `DB_CONNECTION` environment
variable. The following command starts an empty database and calls runs the binary:
```shell script
make run
```

### Make it Your Own

Clone this repo with 
```shell script
git clone https://github.com/jonnylangefeld/go-api.git
```
And rename the module with
```shell script
go mod edit -module <your module name>
```
