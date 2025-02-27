---
title: Java quickstart
linkTitle: Java
description: How to install Testcontainers for Java and run your first container
keywords: testcontainers, testcontainers java, testcontainers oss, testcontainers oss java, testcontainers java quickstart,
    junit, junit4, junit5, junit 4, junit 5
toc_max: 3
weight: 20
---

- [JUnit 4 Quickstart](#junit-4-quickstart)
- [JUnit 5 Quickstart](#junit-5-quickstart)
- [Spock Quickstart](#spock-quickstart)

## JUnit 4 Quickstart

It's easy to add Testcontainers to your project - let's walk through a quick example to see how.

Let's imagine we have a simple program that has a dependency on Redis, and we want to add some tests for it.
In our imaginary program, there is a `RedisBackedCache` class which stores data in Redis.
 
You can see an example test that could have been written for it (without using Testcontainers):

```java
public class RedisBackedCacheIntTestStep0 {

    private RedisBackedCache underTest;

    @Before
    public void setUp() {
        // Assume that we have Redis running locally?
        underTest = new RedisBackedCache("localhost", 6379);
    }

    @Test
    public void testSimplePutAndGet() {
        underTest.put("test", "example");

        String retrieved = underTest.get("test");
        assertThat(retrieved).isEqualTo("example");
    }
}
```

Notice that the existing test has a problem - it's relying on a local installation of Redis, which is a red flag for test reliability.
This may work if we were sure that every developer and CI machine had Redis installed, but would fail otherwise.
We might also have problems if we attempted to run tests in parallel, such as state bleeding between tests, or port clashes.

Let's start from here, and see how to improve the test with Testcontainers:  

### Step 1: Add Testcontainers as a test-scoped dependency

First, add Testcontainers as a dependency as follows:

{{< tabs group="lang" >}}
{{< tab name="Gradle" >}}

```groovy
testImplementation "org.testcontainers:testcontainers:{{latest_version}}"
```

{{< /tab >}}
{{< tab name="Maven" >}}

```xml
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>testcontainers</artifactId>
    <version>{{latest_version}}</version>
    <scope>test</scope>
</dependency>
```

{{< /tab >}}
{{< /tabs >}}

### Step 2: Get Testcontainers to run a Redis container during our tests

Simply add the following to the body of our test class:

```java
@Rule
public GenericContainer redis = new GenericContainer(DockerImageName.parse("redis:6-alpine"))
    .withExposedPorts(6379);
```

The `@Rule` annotation tells JUnit to notify this field about various events in the test lifecycle.
In this case, our rule object is a Testcontainers `GenericContainer`, configured to use a specific Redis image from Docker Hub, and configured to expose a port.

If we run our test as-is, then regardless of the actual test outcome, we'll see logs showing us that Testcontainers:

* was activated before our test method ran
* discovered and quickly tested our local Docker setup
* pulled the image if necessary
* started a new container and waited for it to be ready
* shut down and deleted the container after the test

### Step 3: Make sure our code can talk to the container

Before Testcontainers, we might have hardcoded an address like `localhost:6379` into our tests.

Testcontainers uses *randomized ports* for each container it starts, but makes it easy to obtain the actual port at runtime.
We can do this in our test `setUp` method, to set up our component under test:

```java
String address = redis.getHost();
Integer port = redis.getFirstMappedPort();

// Now we have an address and port for Redis, no matter where it is running
underTest = new RedisBackedCache(address, port);
```

> [!TIP]
>
> Notice that we also ask Testcontainers for the container's actual address with `redis.getHost();`, 
> rather than hard-coding `localhost`. `localhost` may work in some environments but not others - for example it may
> not work on your current or future CI environment. As such, **avoid hard-coding** the address, and use 
> `getHost()` instead.

### Step 4: Run the tests!

That's it!

Let's look at our complete test class to see how little we had to add to get up and running with Testcontainers:

```java
public class RedisBackedCacheIntTest {

    private RedisBackedCache underTest;

    // rule {
    @Rule
    public GenericContainer redis = new GenericContainer(DockerImageName.parse("redis:6-alpine"))
        .withExposedPorts(6379);

    // }

    @Before
    public void setUp() {
        String address = redis.getHost();
        Integer port = redis.getFirstMappedPort();

        // Now we have an address and port for Redis, no matter where it is running
        underTest = new RedisBackedCache(address, port);
    }

    @Test
    public void testSimplePutAndGet() {
        underTest.put("test", "example");

        String retrieved = underTest.get("test");
        assertThat(retrieved).isEqualTo("example");
    }
}
```

## JUnit 5 Quickstart

It's easy to add Testcontainers to your project - let's walk through a quick example to see how.

Let's imagine we have a simple program that has a dependency on Redis, and we want to add some tests for it.
In our imaginary program, there is a `RedisBackedCache` class which stores data in Redis.
 
You can see an example test that could have been written for it (without using Testcontainers):

```java
public class RedisBackedCacheIntTestStep0 {

    private RedisBackedCache underTest;

    @BeforeEach
    public void setUp() {
        // Assume that we have Redis running locally?
        underTest = new RedisBackedCache("localhost", 6379);
    }

    @Test
    public void testSimplePutAndGet() {
        underTest.put("test", "example");

        String retrieved = underTest.get("test");
        assertThat(retrieved).isEqualTo("example");
    }
}
```

Notice that the existing test has a problem - it's relying on a local installation of Redis, which is a red flag for test reliability.
This may work if we were sure that every developer and CI machine had Redis installed, but would fail otherwise.
We might also have problems if we attempted to run tests in parallel, such as state bleeding between tests, or port clashes.

Let's start from here, and see how to improve the test with Testcontainers:  

### Step 1: Add Testcontainers as a test-scoped dependency

First, add Testcontainers as a dependency as follows:

{{< tabs group="lang" >}}
{{< tab name="Gradle" >}}

```groovy
testImplementation "org.junit.jupiter:junit-jupiter:5.8.1"
testImplementation "org.testcontainers:testcontainers:{{latest_version}}"
testImplementation "org.testcontainers:junit-jupiter:{{latest_version}}"
```

{{< /tab >}}
{{< tab name="Maven" >}}

```xml
<dependency>
    <groupId>org.junit.jupiter</groupId>
    <artifactId>junit-jupiter</artifactId>
    <version>5.8.1</version>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>testcontainers</artifactId>
    <version>{{latest_version}}</version>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>junit-jupiter</artifactId>
    <version>{{latest_version}}</version>
    <scope>test</scope>
</dependency>
```

{{< /tab >}}
{{< /tabs >}}

### Step 2: Get Testcontainers to run a Redis container during our tests

First, you'll need to annotate the test class with `@Testcontainers`. Furthermore, add the following to the body of our test class:

```java
@Container
public GenericContainer redis = new GenericContainer(DockerImageName.parse("redis:6-alpine"))
    .withExposedPorts(6379);
```

The `@Container` annotation tells JUnit to notify this field about various events in the test lifecycle.
In this case, our rule object is a Testcontainers `GenericContainer`, configured to use a specific Redis image from Docker Hub, and configured to expose a port.

If we run our test as-is, then regardless of the actual test outcome, we'll see logs showing us that Testcontainers:

* was activated before our test method ran
* discovered and quickly tested our local Docker setup
* pulled the image if necessary
* started a new container and waited for it to be ready
* shut down and deleted the container after the test

### Step 3: Make sure our code can talk to the container

Before Testcontainers, we might have hardcoded an address like `localhost:6379` into our tests.

Testcontainers uses *randomized ports* for each container it starts, but makes it easy to obtain the actual port at runtime.
We can do this in our test `setUp` method, to set up our component under test:

```java
String address = redis.getHost();
Integer port = redis.getFirstMappedPort();

// Now we have an address and port for Redis, no matter where it is running
underTest = new RedisBackedCache(address, port);
```

> [!TIP]
>
> Notice that we also ask Testcontainers for the container's actual address with `redis.getHost();`, 
> rather than hard-coding `localhost`. `localhost` may work in some environments but not others - for example it may
> not work on your current or future CI environment. As such, **avoid hard-coding** the address, and use 
> `getHost()` instead.

### Step 4: Additional attributes

Additional attributes are available for the `@Testcontainers` annotation.
Those attributes can be helpful when:

* Tests should be skipped instead of failing because Docker is unavailable in the
current environment. Set `disabledWithoutDocker` to `true`.
* Enable parallel container initialization instead of sequential (by default). Set `parallel` to `true`.

### Step 5: Run the tests!

That's it!

Let's look at our complete test class to see how little we had to add to get up and running with Testcontainers:

```java
@Testcontainers
public class RedisBackedCacheIntTest {

    private RedisBackedCache underTest;

    // container {
    @Container
    public GenericContainer redis = new GenericContainer(DockerImageName.parse("redis:6-alpine"))
        .withExposedPorts(6379);

    // }

    @BeforeEach
    public void setUp() {
        String address = redis.getHost();
        Integer port = redis.getFirstMappedPort();

        // Now we have an address and port for Redis, no matter where it is running
        underTest = new RedisBackedCache(address, port);
    }

    @Test
    public void testSimplePutAndGet() {
        underTest.put("test", "example");

        String retrieved = underTest.get("test");
        assertThat(retrieved).isEqualTo("example");
    }
}
```

## Spock Quickstart

It's easy to add Testcontainers to your project - let's walk through a quick example to see how.

Let's imagine we have a simple program that has a dependency on Redis, and we want to add some tests for it.
In our imaginary program, there is a `RedisBackedCache` class which stores data in Redis.
 
You can see an example test that could have been written for it (without using Testcontainers):

```groovy
class RedisBackedCacheIntTestStep0 extends Specification {
    private RedisBackedCache underTest

    void setup() {
        // Assume that we have Redis running locally?
        underTest = new RedisBackedCache("localhost", 6379)
    }

    void testSimplePutAndGet() {
        setup:
        underTest.put("test", "example")

        when:
        String retrieved = underTest.get("test")

        then:
        retrieved == "example"
    }
}
```

Notice that the existing test has a problem - it's relying on a local installation of Redis, which is a red flag for test reliability.
This may work if we were sure that every developer and CI machine had Redis installed, but would fail otherwise.
We might also have problems if we attempted to run tests in parallel, such as state bleeding between tests, or port clashes.

Let's start from here, and see how to improve the test with Testcontainers:  

### Step 1: Add Testcontainers as a test-scoped dependency

First, add Testcontainers as a dependency as follows:

{{< tabs group="lang" >}}
{{< tab name="Gradle" >}}

```groovy
testImplementation "org.testcontainers:spock:{{latest_version}}"
```

{{< /tab >}}
{{< tab name="Maven" >}}

```xml
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>spock</artifactId>
    <version>{{latest_version}}</version>
    <scope>test</scope>
</dependency>
```

{{< /tab >}}
{{< /tabs >}}

### Step 2: Get Testcontainers to run a Redis container during our tests

Annotate the Spock specification class with the Testcontainers extension:

```groovy
@org.testcontainers.spock.Testcontainers
class RedisBackedCacheIntTest extends Specification {
```

And add the following field to the body of our test class:

```groovy
GenericContainer redis = new GenericContainer<>("redis:6-alpine")
    .withExposedPorts(6379)
```

This tells Spock to start a Testcontainers `GenericContainer`, configured to use a specific Redis image from Docker Hub, and configured to expose a port.

If we run our test as-is, then regardless of the actual test outcome, we'll see logs showing us that Testcontainers:

* was activated before our test method ran
* discovered and quickly tested our local Docker setup
* pulled the image if necessary
* started a new container and waited for it to be ready
* shut down and deleted the container after the test

### Step 3: Make sure our code can talk to the container

Before Testcontainers, we might have hardcoded an address like `localhost:6379` into our tests.

Testcontainers uses *randomized ports* for each container it starts, but makes it easy to obtain the actual port at runtime.
We can do this in our test `setup` method, to set up our component under test:

```groovy
String address = redis.host
Integer port = redis.firstMappedPort

// Now we have an address and port for Redis, no matter where it is running
underTest = new RedisBackedCache(address, port)
```

> [!TIP]
>
> Notice that we also ask Testcontainers for the container's actual address with `redis.containerIpAddress`, 
> rather than hard-coding `localhost`. `localhost` may work in some environments but not others - for example it may
> not work on your current or future CI environment. As such, **avoid hard-coding** the address, and use 
> `containerIpAddress` instead.

### Step 4: Run the tests!

That's it!

Let's look at our complete test class to see how little we had to add to get up and running with Testcontainers:

```groovy
@org.testcontainers.spock.Testcontainers
class RedisBackedCacheIntTest extends Specification {

    private RedisBackedCache underTest

    GenericContainer redis = new GenericContainer<>("redis:6-alpine")
        .withExposedPorts(6379)

    void setup() {
        String address = redis.host
        Integer port = redis.firstMappedPort

        // Now we have an address and port for Redis, no matter where it is running
        underTest = new RedisBackedCache(address, port)
    }

    void testSimplePutAndGet() {
        setup:
        underTest.put("test", "example")

        when:
        String retrieved = underTest.get("test")

        then:
        retrieved == "example"
    }
}
```
