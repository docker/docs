---
description: Testing contributions
keywords: test, source, contributions, Docker
title: Testing contributions
---

Testing is about software quality, performance, reliability, or product
usability. We develop and test Docker software before we release but we are
human. So, we make mistakes, we get forgetful, or we just don't have enough
time to do everything.

Choose to contribute testing if you want to improve Docker software and
processes. Testing is a good choice for contributors that have experience
in software testing, usability testing, or who are otherwise great at spotting
problems.

## What can you contribute to testing?

* Write a blog about <a href="http://www.appneta.com/blog/automated-testing-with-docker/" target="_blank">how your company uses Docker as its test infrastructure</a>.
* Take <a href="http://ows.io/tj/w88v3siv" target="_blank">an online usability test</a> or create a usability test about Docker.
* Test one of<a href="https://github.com/docker-library/official-images/issues"> Docker's official images</a>.
* Test the Docker documentation.


## Test the Docker documentation

Testing documentation is relatively easy:

1.  Find a page in [Docker's documentation](/) that contains a procedure or example you want to test.

    You should choose something that _should work_ on your machine. For example,
    [creating a base image](/engine/userguide/eng-image/baseimages/){: target="_blank" class="_" }
    is something anyone with Docker can do, while [changing volume directories in Kitematic](https://kitematic.com/docs/managing-volumes/){: target="_blank" class="_" }
    requires Kitematic installed on a Mac.

2.  Try and follow the procedure or recreate the example.

    Look for:

    * Are the steps clearly laid out and identifiable?
    * Are the steps in the right order?
    * Did you get the results the procedure or example said you would?

3.  If you couldn't complete the procedure or example,
    [file an issue](https://github.com/moby/moby/issues/new){: target="_blank" class="_" }{: target="_blank" class="_" }.

## Test code in Docker

If you are interested in writing or fixing test code for the Docker project, learn  about  <a href="/opensource/project/test-and-docs/" target="_blank">our test infrastructure</a>.

View <a href="https://goo.gl/JWuQPJ" target="_blank"> our open test issues</a> in Docker for something to work on. Or, create one of your own.
