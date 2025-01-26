---
title: Understanding the application
linkTitle: Understand the application
weight: 20 #
keywords: go, golang, prometheus, grafana, containerize, monitor
description: Understand the Golang application that you will containerize and monitor with Prometheus and Grafana.
aliases:
  - /guides/go-monitoring/application/
---

## Prerequisites

* You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

## Overview

To make sure our application is working an intended, monitoring is really important. In this guide, you'll learn how to containerize a Golang application and monitor it with Prometheus and Grafana. In this guide, you'll learn how to:

## Getting the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following
command to clone the repository:

```console
$ git clone https://github.com/dockersamples/bun-docker.git
```

You should now have the following contents in your `bun-docker` directory.

```text

// TODO:

│── go-monitoring/
│ ├── Docker
│ │   ├── grafana.yml
│ │   └── prometheus.yml
│ ├── Dockerfile
│ ├── compose.yml
│ ├── go.mod
│ ├── go.sum
│ ├── README.md
│ └── main.go
```

## Understanding the application

The complete application logic is in `main.go` file. Let's understand the application in detail before jumping into containerizing it. 

To connect the pieces better, here is the complete code.

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	customRegistry = prometheus.NewRegistry()
)

var (
	HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests processed by the API",
	}, []string{"path", "status"})

	HttpRequestErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of errors processed by the API",
	}, []string{"path", "status"})
)

func main() {
	r := gin.Default()

	r.Use(RequestMetricsMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/metrics", func(ctx *gin.Context) {
		h := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
		h.ServeHTTP(ctx.Writer, ctx.Request)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})
	r.Run(":8080")
}

func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		if status <= 400 {
			HttpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
			return
		}
		HttpRequestErrorTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
	}
}

func init() {
	customRegistry.MustRegister(HttpRequestTotal)
	customRegistry.MustRegister(HttpRequestErrorTotal)
}
```

In this part of the code, we have define the packages that we are going to use in the application. We are using `gin-gonic/gin` as the web framework, `prometheus/client_golang` for monitoring and other supporting packages. Below we have a few variables, `customRegistry` we created a new registry only to register our custom metrics because by default, Prometheus registers all the metrics in the default registry. Below that we have defined two custom metrics, `HttpRequestTotal` and `HttpRequestErrorTotal`. `HttpRequestTotal` is a counter that increments every time a successful request is made to the API and `HttpRequestErrorTotal` is a counter that increments every time an error occurs in the API.

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	customRegistry = prometheus.NewRegistry()
)

var (
	HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests processed by the API",
	}, []string{"path", "status"})

	HttpRequestErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of errors processed by the API",
	}, []string{"path", "status"})
)
```

In the `main` function, we have created a new instance of the `gin` framework and registered a middleware function `RequestMetricsMiddleware` that will be called for every request made to the API. We have defined three routes, `/`, `/metrics`, and `/health`. The `/` route returns a JSON response with a message `ok`, the `/metrics` route returns the metrics that Prometheus scrapes, and the `/health` route returns a JSON response with a message `healthy` which we will use in Docker Compose to check the health of the application.

```golang
func main() {
	r := gin.Default()

	r.Use(RequestMetricsMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/metrics", func(ctx *gin.Context) {
		h := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
		h.ServeHTTP(ctx.Writer, ctx.Request)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})
	r.Run(":8080")
}
```

Now comes the middleware function `RequestMetricsMiddleware`. This function is called for every request made to the API. It increments the `HttpRequestTotal` counter (different counter for different paths and status codes) if the status code is less than or equal to 400. If the status code is greater than 400, it increments the `HttpRequestErrorTotal` counter (different counter for different paths and status codes).

At the end of the file, we have an `init` function it is really important as it registers the custom metrics in the custom registry that we created at the beginning of the file. And that will show up in the Prometheus metrics on the `/metrics` route.

```golang
func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		if status <= 400 {
			HttpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
			return
		}
		HttpRequestErrorTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
	}
}

func init() {
	customRegistry.MustRegister(HttpRequestTotal)
	customRegistry.MustRegister(HttpRequestErrorTotal)
}
```

That was the complete explanation of the application. For experts, yes, we can have better structure the app with other practices. But for the sake of this guide, and to keep it simple, we have kept it this way. Now let's run the application and see if it works as intended.

## Run the application

Inside the `bun-docker` directory, run the following command in a terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000). You will see a message `{"Status" : "OK"}` in the browser.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `bun-docker` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3000](http://localhost:3000).


In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

## Summary

In this section, you learned how you can containerize and run your Bun
application using Docker.

Related information:

 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file](/reference/dockerfile.md#dockerignore-file)
 - [Docker Compose overview](/manuals/compose/_index.md)
 - [Compose file reference](/reference/compose-file/_index.md)

## Next steps

In the next section, you'll learn how you can develop your application using
containers.
