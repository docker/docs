---
published: false
---

This is an initial attempt to make it easier to test the examples in the https.md
doc.

At this point, it is a manual thing, and I've been running it in boot2docker.

My process is as following:

    $ boot2docker ssh
    root@boot2docker:/# git clone https://github.com/moby/moby
    root@boot2docker:/# cd docker/docs/articles/https
    root@boot2docker:/# make cert

lots of things to see and manually answer, as openssl wants to be interactive

**NOTE:** make sure you enter the hostname (`boot2docker` in my case) when prompted for `Computer Name`)

    root@boot2docker:/# sudo make run

Start another terminal:

    $ boot2docker ssh
    root@boot2docker:/# cd docker/docs/articles/https
    root@boot2docker:/# make client

The last connects first with `--tls` and then with `--tlsverify`, both should succeed.
