[![Go Report Card](https://goreportcard.com/badge/github.com/kecci/goscription)](https://goreportcard.com/report/github.com/kecci/goscription)
# Goscription

Goscription is a sample of template RESTful API Project

## Features

Some features & libraries used on this template:
1. REST API (**labstack/echo**)
2. Dependency Injection (**uber-go/fx**)
3. NSQ Messaging (**segmentio/nsq-go**)
4. Custom CLI (**spf13/cobra**)
5. Custom Config File (**spf13/viper**)
6. SQL Generator (**squirrel**)
7. Migrations (**go-migrate**)
8. Swagger API Docs Generator (**swaggo/swag**)
9. Mock Generator (**vektra/mockery**)
10. Custom Logger (**sirupsen/logrus**)
11. Dockerize an Application (**docker**)
12. Circuit Breaker (**hystrix-go/hystrix** && **eapache/go-resiliency**)
13. Message Streams: Pub/Sub, Kafka, RabbitMQ, etc (**ThreeDotsLabs/watermill**) [Soon]

## Installation

* Clone the Repos
```bash
$ git clone https://github.com/kecci/goscription.git
```

## Command
This command is provided in Makefile. You can see all command with make:
```bash
$ make

 Choose a command run in goscription:

  mockery-prepare    install mockery before generate
  mockery-generate   generate all mock
  test-docker        test docker integration
  test-unit          run all unit test
  go-build           build to compile the project & swagger docs
  go-run             run project
  docker-build       dockerize the project
  docker-up          run docker compose up
  docker-down        run docker compose down
  swagger-init       initialize swagger to folder ./docs
  swagger-validate   validate swagger.yaml in folder ./docs
  migrate-prepare    prepare migrate the schema with mysql
  migrate-up         run migration up to latest version
  migrate-down       run migration down to oldest version
```

## Project Structure

```bash
.
├── app
│   ├── cmd
│   └── main.go
├── internal
│   ├── controller
│   ├── database
│   │   └── mysql
│   │       └── migrations
│   ├── middleware
│   ├── outbound
│   └── service
├── mocks
├── models
└── util
```

## Dependency Injection
I found this library is very useful and you no need to generate anything. Just code. Very modular and clear layer.

### uber-go/fx
A dependency injection based application framework for Go.

```go
func main() {
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		fx.Provide( 
			NewTimeOutContext, 
			NewDbConn, 
		),
		repository.Module, // Provide All Repository Module
		service.Module, // Provide All Service Module
		outbound.Module, // Provide All Outbound Module
		server.Module, // Provide Runner/Stopper Server
		fx.Invoke( // Invoke List of Controller
      controller.InitFooController, 
      controller.InitBarController,
		),
	)
}
```
For more information about uber-go/fx: https://github.com/uber-go/fx

## Config Properties
### spf13/viper
Viper is a complete configuration solution for Go applications. viper can reading from: 
* JSON, 
* TOML, 
* YAML, 
* HCL, 
* INI,
* envfile, and 
* Java properties config files

Reading Config Files:
```go
viper.SetConfigName("config") // name of config file (without extension)
viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
viper.AddConfigPath(".")               // optionally look for config in the working directory
err := viper.ReadInConfig() // Find and read the config file
if err != nil { // Handle errors reading the config file
	panic(fmt.Errorf("Fatal error config file: %s \n", err))
}
```

Use Value from Config:
```go
viper.Get("name") // this would be "steve"
```

You can do a lot more with viper, see more the documetation: https://github.com/spf13/viper

### TOML
We are using toml in this sample project, for example:
```toml
title="Configuration File for Goscription"
debug=true
contextTimeout="5"
[server]
  address= ":9090"
[database]
  host="mysql"
  port="3306"
  user="root"
  pass="root"
  name="article"
```

## Swagger

### swag UI
Because our project setup the path of swagger in: `/swagger/*`. You can access swagger UI in here: http://localhost:9090/swagger/index.html

### swaggo/swag
Swag converts Go annotations to Swagger Documentation 2.0. We've created a variety of plugins for popular Go web frameworks. This allows you to quickly integrate with an existing Go project (using Swagger UI).

Supported Web Frameworks:
* gin
* echo
* buffalo
* net/http

### swag annotation
Swag has handled your swagger docs. So you no longer need to write `swagger.yml` or `swagger.json`. What you need to do is just write annotations. This is an example:

```go
// @title Blueprint Swagger API
// @version 1.0
// @description Swagger API for Golang Project Blueprint.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email martin7.heinz@gmail.com

// @license.name MIT
// @license.url https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

// @BasePath /api/v1
func main() {
    ...
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    ...
}
```

You can see more about swag annotation in here: https://github.com/swaggo/swag.

## Benchmark
Using https://github.com/wg/wrk

Script:
```sh
wrk -t8 -c256 -d30s http://localhost:9090/health
```

Output:
```sh
Running 30s test @ http://localhost:9090/health
  8 threads and 256 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    13.17ms   11.67ms 153.74ms   71.04%
    Req/Sec     2.82k   660.71     6.40k    69.88%
  675770 requests in 30.08s, 150.16MB read
  Socket errors: connect 0, read 121, write 0, timeout 0
Requests/sec:  22464.24
Transfer/sec:      4.99MB
```

## Sources
This template is inspired & modified from https://github.com/golangid/menekel
