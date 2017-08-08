---
title: "Play with Docker"
keywords: get started
description: Play with Docker tests
---
{% include pwd.html %}

## Play with Docker Test 1

Here is a play with Docker terminal.
<div id="term1" class="term1" style="height: 300px; width: 400px; position: fixed; right: 20px; top: 100px; z-index:100;"></div>
<div id="term2" class="term2" style="height: 300px; width: 400px; position: fixed; right: 20px; top: 425px; z-index:100;"></div>


## Play with Docker Test 2

```.term1
docker container run hello-world
```
```.term2
docker container run hello-world
```
```.term3
docker container run hello-world
```

<script>
        pwd.newSession([{selector: '.term1'}, {selector: '.term2'}, {selector: '.term3'}], {ImageName: ''});
// If you can get a local running version of PWD, use 
// pwd.newSession([{selector: '.term1'}, {selector: '.term2'}, {selector: '.term3'}], {ImageName: ''}, baseUrl: 'http://localhost');
// Directions for running locally: github.com/play-with-docker/play-with-docker

</script>