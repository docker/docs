---
title: "Run your tests using Go test"
keywords: build, go, golang, test
description: How to build and run your Go tests in a container
redirect_from:
- /get-started/golang/run-tests/
---

{% include_relative nav.html selected="4" %}

Testing is an essential part of modern software development. Yet, testing can mean a lot of things to different development teams. In the name of brevity, we'll only take a look at running isolated, high-level, functional tests.

## Test structure

Each test is meant to verify a single business requirement for our example application. The following test is an excerpt from `main_test.go` test suite in our example application.

{% raw %}
```go
func TestRespondsWithLove(t *testing.T) {

	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	resource, err := pool.Run("docker-gs-ping", "latest", []string{})
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})

	var resp *http.Response

	err = pool.Retry(func() error {
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8080/tcp"), "/"))
		if err != nil {
			t.Log("container not ready, waiting...")
			return err
		}
		return nil
	})
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	// Finally, test the business requirement!
	require.Contains(t, string(body), "<3", "does not respond with love?")
}
```
{% endraw %}

As you can see, this is a high-level test, unconcerned with the implementation details of our example application.

* the test is using [`ory/dockertest`](https://github.com/ory/dockertest) Go module;
* the test assumes that the Docker engine instance is running on the same machine where the test is being run.

The second test in `main_test.go` has almost identical structure but it tests _another_ business requirement of our application. You are welcome to have a look at all available tests in [`docker-gs-ping/main_test.go`](https://github.com/olliefr/docker-gs-ping/blob/main/main_test.go).

## Run tests locally

In order to run the tests, we must make sure that our application Docker image is up-to-date.

```shell
docker build -t docker-gs-ping:latest .
```

```
[+] Building 3.0s (13/13) FINISHED
...
```

In the above example we've omitted most of the output, only displaying the first line indicating that the build was successful.

Note, that the image is tagged with `latest` which is the same label we've chosen to use in our `main_test.go` tests. 

Now that the Docker image for our application had been built, we can run the tests that depend on it:

```shell
go test ./...
ok      github.com/olliefr/docker-gs-ping       2.564s
```

That was a bit... underwhelming? Let's ask it to print a bit more detail, just to be sure:

```shell
go test -v ./...
```

```
=== RUN   TestRespondsWithLove
    main_test.go:47: container not ready, waiting...
--- PASS: TestRespondsWithLove (5.24s)
=== RUN   TestHealthCheck
    main_test.go:83: container not ready, waiting...
--- PASS: TestHealthCheck (1.40s)
PASS
ok      github.com/olliefr/docker-gs-ping       6.670s
```

So, the tests do, indeed, pass. Note, how retrying using exponential back-off helped avoiding failing tests while the containers are being initialised. What happens in each test is that `ory/dockertest` module connects to the local Docker engine instance and instructs it to spin up a container using the image, identified by the tag `docker-gs-ping:latest`. Starting up a container may take a while, so our tests retry accessing the container until the container is ready to respond to requests.

## Next steps

In this module, we've seen an example of using Docker for isolated functional testing of an example Go application. There are many different ways to test an application and we have only considered the high-level, functional testing. This, however, feeds naturally into the next topic, where we are going to set up our tests to run in an automated pipeline.

In the next module, we’ll take a look at how to set up a CI/CD pipeline using GitHub Actions. See:

[Configure CI/CD](configure-ci-cd.md){: .button .outline-btn}

## Feedback

Help us improve this topic by providing your feedback. Let us know what you think by creating an issue in the [Docker Docs](https://github.com/docker/docker.github.io/issues/new?title=[Golang%20docs%20feedback]){:target="_blank" rel="noopener" class="_"} GitHub repository. Alternatively, [create a PR](https://github.com/docker/docker.github.io/pulls){:target="_blank" rel="noopener" class="_"} to suggest updates.

<br />
