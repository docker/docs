---
title: Building the application
linkTitle: Understand the application
weight: 10 #
keywords: go, golang, prometheus, grafana, containerize, monitor
description: Learn how to create a Golang server to register metrics with Prometheus.
---

## Prerequisites

* You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

You will be creating a Golang server with some endpoints to simulate a real-world application. Then you will expose metrics from the server using Prometheus.

## Getting the sample application

Clone the sample application to use with this guide. Open a terminal, change
directory to a directory that you want to work in, and run the following
command to clone the repository:

```console
$ git clone https://github.com/dockersamples/go-prometheus-monitoring.git 
```

Once you cloned you will see the following content structure inside `go-prometheus-monitoring` directory,

```text
go-prometheus-monitoring
├── CONTRIBUTING.md
├── Docker
│   ├── grafana.yml
│   └── prometheus.yml
├── dashboard.json
├── Dockerfile
├── LICENSE
├── README.md
├── compose.yaml
├── go.mod
├── go.sum
└── main.go
```

- **main.go** - The entry point of the application.
- **go.mod and go.sum** - Go module files.
- **Dockerfile** - Dockerfile used to build the app.
- **Docker/** - Contains the Docker Compose configuration files for Grafana and Prometheus.
- **compose.yaml** - Compose file to launch everything (Golang app, Prometheus, and Grafana).
- **dashboard.json** - Grafana dashboard configuration file.
- **Dockerfile** - Dockerfile used to build the Golang app.
- **compose.yaml** - Docker Compose file to launch everything (Golang app, Prometheus, and Grafana).
- Other files are for licensing and documentation purposes.

## Understanding the application

The following is the complete logic of the application you will find in `main.go`.

```go
package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define metrics
var (
	HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests processed by the API",
	}, []string{"path", "status"})

	HttpRequestErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of errors returned by the API",
	}, []string{"path", "status"})
)

// Custom registry (without default Go metrics)
var customRegistry = prometheus.NewRegistry()

// Register metrics with custom registry
func init() {
	customRegistry.MustRegister(HttpRequestTotal, HttpRequestErrorTotal)
}

func main() {
	router := gin.Default()

	// Register /metrics before middleware
	router.GET("/metrics", PrometheusHandler())
	
	router.Use(RequestMetricsMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Up and running!",
		})
	})
	router.GET("/v1/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from /v1/users",
		})
	})

	router.Run(":8000")
}

// Custom metrics handler with custom registry
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Middleware to record incoming requests metrics
func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		if status < 400 {
			HttpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		} else {
			HttpRequestErrorTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		}
	}
}
```

In this part of the code, we have imported the required packages `gin`, `prometheus`, and `promhttp`. Then we have defined a couple of variables, `HttpRequestTotal` and `HttpRequestErrorTotal` are Prometheus counter metrics, and `customRegistry` is a custom registry that will be used to register these metrics. The name of the metric is a string that you can use to identify the metric. The help string is a string that will be shown when you query the `/metrics` endpoint to understand the metric. The reason we are using the custom registry is so avoid the default Go metrics that are registered by default by the Prometheus client. Then using the `init` function we are registering the metrics with the custom registry. 

```go
import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define metrics
var (
	HttpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_total",
		Help: "Total number of requests processed by the API",
	}, []string{"path", "status"})

	HttpRequestErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_http_request_error_total",
		Help: "Total number of errors returned by the API",
	}, []string{"path", "status"})
)

// Custom registry (without default Go metrics)
var customRegistry = prometheus.NewRegistry()

// Register metrics with custom registry
func init() {
	customRegistry.MustRegister(HttpRequestTotal, HttpRequestErrorTotal)
}
```

In the `main` function, you have created a new instance of the `gin` framework and created three routes. You can see the health endpoint that is on path `/health` that will return a JSON with `{"message": "Up and running!"}` and the `/v1/users` endpoint that will return a JSON with `{"message": "Hello from /v1/users"}`. The third route is for the `/metrics` endpoint that will return the metrics in the Prometheus format. Then you have `RequestMetricsMiddleware` middleware, it will be called for every request made to the API. It will record the incoming requests metrics like status codes and paths. Finally, you are running the gin application on port 8000.

```golang
func main() {
	router := gin.Default()

	// Register /metrics before middleware
	router.GET("/metrics", PrometheusHandler())
	
	router.Use(RequestMetricsMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Up and running!",
		})
	})
	router.GET("/v1/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from /v1/users",
		})
	})

	router.Run(":8000")
}
```

Now comes the middleware function `RequestMetricsMiddleware`. This function is called for every request made to the API. It increments the `HttpRequestTotal` counter (different counter for different paths and status codes) if the status code is less than or equal to 400. If the status code is greater than 400, it increments the `HttpRequestErrorTotal` counter (different counter for different paths and status codes). The `PrometheusHandler` function is the custom handler that will be called for the `/metrics` endpoint. It will return the metrics in the Prometheus format.

```golang
// Custom metrics handler with custom registry
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.HandlerFor(customRegistry, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Middleware to record incoming requests metrics
func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		c.Next()
		status := c.Writer.Status()
		if status < 400 {
			HttpRequestTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		} else {
			HttpRequestErrorTotal.WithLabelValues(path, strconv.Itoa(status)).Inc()
		}
	}
}
```

That's it, this was the complete gist of our application. Now it's time to run and test if our app is registering metrics correctly.

## Running the application

Make sure you are still inside `go-prometheus-monitoring` directory in the terminal, and run the following command. Install the dependencies by running `go mod tidy` command and then build and run the application by running `go run main.go` command. Then visit `http://localhost:8000/health` or `http://localhost:8000/v1/users`. You should see the output `{"message": "Up and running!"}` or `{"message": "Hello from /v1/users"}`. If you are able to see this then our app is successfully up and running. 


Now let's check our application's metrics by accessing the `/metrics` endpoint. 
Open `http://localhost:8000/metrics` in your browser. You should see similar output to the one below.

```sh
# HELP api_http_request_error_total Total number of errors returned by the API
# TYPE api_http_request_error_total counter
api_http_request_error_total{path="/",status="404"} 1
api_http_request_error_total{path="//v1/users",status="404"} 1
api_http_request_error_total{path="/favicon.ico",status="404"} 1
# HELP api_http_request_total Total number of requests processed by the API
# TYPE api_http_request_total counter
api_http_request_total{path="/health",status="200"} 2
api_http_request_total{path="/v1/users",status="200"} 1
```

In the terminal, press `ctrl` + `c` to stop the application.

> [!Note]
> If you don't want to run the application locally, and want to run it in a Docker container, skip to next page where we create a Dockerfile and containerize the application.

## Summary

In this section, we learned how to create an Golang app to register metrics with Prometheus. By implementing middleware functions, we were able to increment the counters based on the request path and status codes.

## Next steps

In the next section, you'll learn how to containerize your application.
