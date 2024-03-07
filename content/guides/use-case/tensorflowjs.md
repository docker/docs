---
description: Learn how to deploy pre-trained models in a TensorFlow.js web applications to perform face detection.
keywords: tensorflow.js, machine learning, ml, mediapipe, blazeface, face detection
title: Face detection with TensorFlow.js
---

This guide introduces the seamless integration of TensorFlow.js with Docker to
perform face detection. In this guide, you'll explore:

- How to run a containerized TensorFlow.js application with Docker
- The code for a face detection application using TensorFlow.js
- The Dockerfile to build a TensorFlow.js web application

> **Acknowledgment**
>
> Docker would like to thank [Harsh Manvar](https://github.com/harsh4870) for
> his contribution to this guide.

## Prerequisites

* You have installed the latest version of
  [Docker Desktop](../../../get-docker.md).
* You have a [Git client](https://git-scm.com/downloads). The examples in this
  guide use a command-line based Git client, but you can use any client.

## What is TensorFlow.js?

[TensorFlow.js](https://www.tensorflow.org/js) is an open-source JavaScript
library for machine learning that enables you to train and deploy ML models in
the browser or on Node.js. It supports creating new models from scratch or using
pre-trained models, facilitating a wide range of ML applications directly in web
environments. TensorFlow.js offers efficient computation, making sophisticated
ML tasks accessible to web developers without deep ML expertise.

## Why Use TensorFlow.js and Docker together?

- Environment consistency and simplified deployment: Docker packages
  TensorFlow.js applications and their dependencies into containers, ensuring
  consistent runs across all environments and simplifying deployment.
- Efficient development and easy scaling: Docker enhances development efficiency
  with features like hot reloading and facilitates easy scaling of -
  TensorFlow.js applications using orchestration tools like Kubernetes.
- Isolation and enhanced security: Docker isolates TensorFlow.js applications in
  secure environments, minimizing conflicts and security vulnerabilities while
  running applications with limited permissions.


## Get and run the sample application

In a terminal, clone the sample application using the following command.

```console
$ git clone https://github.com/harsh4870/TensorJS-Face-Detection
```

After cloning the application, you'll notice the application has a `Dockerfile`.
This Dockerfile lets you build and run the application locally with nothing more
than Docker.

Before you run the application as a container, you must build it into an image.
Run the following command inside the `TensorJS-Face-Detection` directory to
build an image named `face-detection-tensorjs`.

```console
$ docker build -t face-detection-tensorjs .
```

The command builds the application into an image. Depending on your network
connection, it can take several minutes to download the necessary components the
first time you run the command.

To run the image as a container, run the following command in a terminal.

```console
$ docker run -p 80:80 face-detection-tensorjs
```

The command runs the container and maps port 80 in the container to port 80 on
your system.

Once the application is running, open a web browser and access the application
at [http://localhost:80](http://localhost:80). You may need to grant access to
your webcam for the application.

In the web application, you can change the backend to use one of the following:
- WASM
- WebGL
- CPU

To stop the application, press `ctrl`+`c` in the terminal.

## About the application

The sample application performs real-time face detection using
[MediaPipe](https://developers.google.com/mediapipe/), a comprehensive framework
for building multimodal machine learning pipelines. It's specifically using the
BlazeFace model, a lightweight model for detecting faces in images.

In the context of TensorFlow.js or similar web-based machine learning
frameworks, the WASM, WebGL, and CPU backends can be used to
execute operations. Each of these backends utilizes different resources and
technologies available in modern browsers and has its strengths and limitations.
The following sections are a brief breakdown of the different backends.

### WASM

WebAssembly (WASM) is a low-level, assembly-like language with a compact binary
format that runs at near-native speed in web browsers. It allows code written in
languages like C/C++ to be compiled into a binary that can be executed in the
browser.

It's a good choice when high performance is required, and either the WebGL
backend is not supported or you want consistent performance across all devices
without relying on the GPU.

### WebGL

WebGL is a browser API that allows for GPU-accelerated usage of physics and
image processing and effects as part of the web page canvas.

WebGL is well-suited for operations that are parallelizable and can
significantly benefit from GPU acceleration, such as matrix multiplications and
convolutions commonly found in deep learning models.

### CPU

The CPU backend uses pure JavaScript execution, utilizing the device's central
processing unit (CPU). This backend is the most universally compatible and
serves as a fallback when neither WebGL nor WASM backends are available or
suitable.

## Explore the application's code

Explore the purpose of each file and their contents in the following sections.

### The index.html file

The `index.html` file serves as the frontend for the web application that
utilizes TensorFlow.js for real-time face detection from the webcam video feed.
It incorporates several technologies and libraries to facilitate machine
learning directly in the browser. It uses several TensorFlow.js libraries,
including:

- tfjs-core and tfjs-converter for core TensorFlow.js functionality and model
  conversion.
- tfjs-backend-webgl, tfjs-backend-cpu, and the tf-backend-wasm script
  for different computational backend options that TensorFlow.js can use for
  processing. These backends allow the application to perform machine learning
  tasks efficiently by leveraging the user's hardware capabilities.
- The BlazeFace library, a TensorFlow model for face detection.

It also uses the following additional libraries:

- dat.GUI for creating a graphical interface to interact with the application's
  settings in real-time, such as switching between TensorFlow.js backends.
- Stats.min.js for displaying performance metrics (like FPS) to monitor the
  application's efficiency during operation.

{{< accordion title="index.html" >}}

```html
<style>
  body {
    margin: 25px;
  }

  .true {
    color: green;
  }

  .false {
    color: red;
  }

  #main {
    position: relative;
    margin: 50px 0;
  }

  canvas {
    position: absolute;
    top: 0;
    left: 0;
  }

  #description {
    margin-top: 20px;
    width: 600px;
  }

  #description-title {
    font-weight: bold;
    font-size: 18px;
  }
</style>

<body>
  <div id="main">
    <video id="video" playsinline style="
      -webkit-transform: scaleX(-1);
      transform: scaleX(-1);
      width: auto;
      height: auto;
      ">
    </video>
    <canvas id="output"></canvas>
    <video id="video" playsinline style="
      -webkit-transform: scaleX(-1);
      transform: scaleX(-1);
      visibility: hidden;
      width: auto;
      height: auto;
      ">
    </video>
  </div>
</body>
<script src="https://unpkg.com/@tensorflow/tfjs-core@2.1.0/dist/tf-core.js"></script>
<script src="https://unpkg.com/@tensorflow/tfjs-converter@2.1.0/dist/tf-converter.js"></script>

<script src="https://unpkg.com/@tensorflow/tfjs-backend-webgl@2.1.0/dist/tf-backend-webgl.js"></script>
<script src="https://unpkg.com/@tensorflow/tfjs-backend-cpu@2.1.0/dist/tf-backend-cpu.js"></script>
<script src="./tf-backend-wasm.js"></script>

<script src="https://unpkg.com/@tensorflow-models/blazeface@0.0.5/dist/blazeface.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/dat-gui/0.7.6/dat.gui.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/stats.js/r16/Stats.min.js"></script>
<script src="./index.js"></script>
```

{{< /accordion >}}

### The index.js file

The `index.js` file conducts the facial detection logic. It demonstrates several
advanced concepts in web development and machine learning integration. Here's a
breakdown of some of its key components and functionalities:

- Stats.js: The script starts by creating a Stats instance to monitor and
  display the frame rate (FPS) of the application in real time. This is helpful
  for performance analysis, especially when testing the impact of different
  TensorFlow.js backends on the application's speed.
- TensorFlow.js: The application allows users to switch between different
  computation backends (wasm, webgl, and cpu) for TensorFlow.js through a
  graphical interface provided by dat.GUI. Changing the backend can affect
  performance and compatibility depending on the device and browser. The
  addFlagLabels function dynamically checks and displays whether SIMD (Single
  Instruction, Multiple Data) and multithreading are supported, which are
  relevant for performance optimization in the wasm backend.
- setupCamera function: Initializes the user's webcam using the MediaDevices Web
  API. It configures the video stream to not include audio and to use the
  front-facing camera (facingMode: 'user'). Once the video metadata is loaded,
  it resolves a promise with the video element, which is then used for face
  detection.
- BlazeFace: The core of this application is the renderPrediction function,
  which performs real-time face detection using the BlazeFace model, a
  lightweight model for detecting faces in images. The function calls
  model.estimateFaces on each animation frame to detect faces from the video
  feed. For each detected face, it draws a red rectangle around the face and
  blue dots for facial landmarks on a canvas overlaying the video.


{{< accordion title="index.js" >}}

```javascript
const stats = new Stats();
stats.showPanel(0);
document.body.prepend(stats.domElement);

let model, ctx, videoWidth, videoHeight, video, canvas;

const state = {
  backend: 'wasm'
};

const gui = new dat.GUI();
gui.add(state, 'backend', ['wasm', 'webgl', 'cpu']).onChange(async backend => {
  await tf.setBackend(backend);
  addFlagLables();
});

async function addFlagLables() {
  if(!document.querySelector("#simd_supported")) {
    const simdSupportLabel = document.createElement("div");
    simdSupportLabel.id = "simd_supported";
    simdSupportLabel.style = "font-weight: bold";
    const simdSupported = await tf.env().getAsync('WASM_HAS_SIMD_SUPPORT');
    simdSupportLabel.innerHTML = `SIMD supported: <span class=${simdSupported}>${simdSupported}<span>`;
    document.querySelector("#description").appendChild(simdSupportLabel);
  }

  if(!document.querySelector("#threads_supported")) {
    const threadSupportLabel = document.createElement("div");
    threadSupportLabel.id = "threads_supported";
    threadSupportLabel.style = "font-weight: bold";
    const threadsSupported = await tf.env().getAsync('WASM_HAS_MULTITHREAD_SUPPORT');
    threadSupportLabel.innerHTML = `Threads supported: <span class=${threadsSupported}>${threadsSupported}</span>`;
    document.querySelector("#description").appendChild(threadSupportLabel);
  }
}

async function setupCamera() {
  video = document.getElementById('video');

  const stream = await navigator.mediaDevices.getUserMedia({
    'audio': false,
    'video': { facingMode: 'user' },
  });
  video.srcObject = stream;

  return new Promise((resolve) => {
    video.onloadedmetadata = () => {
      resolve(video);
    };
  });
}

const renderPrediction = async () => {
  stats.begin();

  const returnTensors = false;
  const flipHorizontal = true;
  const annotateBoxes = true;
  const predictions = await model.estimateFaces(
    video, returnTensors, flipHorizontal, annotateBoxes);

  if (predictions.length > 0) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    for (let i = 0; i < predictions.length; i++) {
      if (returnTensors) {
        predictions[i].topLeft = predictions[i].topLeft.arraySync();
        predictions[i].bottomRight = predictions[i].bottomRight.arraySync();
        if (annotateBoxes) {
          predictions[i].landmarks = predictions[i].landmarks.arraySync();
        }
      }

      const start = predictions[i].topLeft;
      const end = predictions[i].bottomRight;
      const size = [end[0] - start[0], end[1] - start[1]];
      ctx.fillStyle = "rgba(255, 0, 0, 0.5)";
      ctx.fillRect(start[0], start[1], size[0], size[1]);

      if (annotateBoxes) {
        const landmarks = predictions[i].landmarks;

        ctx.fillStyle = "blue";
        for (let j = 0; j < landmarks.length; j++) {
          const x = landmarks[j][0];
          const y = landmarks[j][1];
          ctx.fillRect(x, y, 5, 5);
        }
      }
    }
  }

  stats.end();

  requestAnimationFrame(renderPrediction);
};

const setupPage = async () => {
  await tf.setBackend(state.backend);
  addFlagLables();
  await setupCamera();
  video.play();

  videoWidth = video.videoWidth;
  videoHeight = video.videoHeight;
  video.width = videoWidth;
  video.height = videoHeight;

  canvas = document.getElementById('output');
  canvas.width = videoWidth;
  canvas.height = videoHeight;
  ctx = canvas.getContext('2d');
  ctx.fillStyle = "rgba(255, 0, 0, 0.5)";

  model = await blazeface.load();

  renderPrediction();
};

setupPage();
```

{{< /accordion >}}

### The tf-backend-wasm.js file

The `tf-backend-wasm.js` file is part of the
[TensorFlow.js library](https://github.com/tensorflow/tfjs/tree/master/tfjs-backend-wasm).
It contains initialization logic for the TensorFlow.js WASM backend, some
utilities for interacting with the WASM binaries, and functions to set custom
paths for the WASM binaries.

### The tfjs-backend-wasm-simd.wasm file

The `tfjs-backend-wasm-simd.wasm` file is part of the
[TensorFlow.js library](https://github.com/tensorflow/tfjs/tree/master/tfjs-backend-wasm).
It's a WASM binary that's used for the WebAssembly
backend, specifically optimized to utilize SIMD (Single Instruction, Multiple
Data) instructions.

## The Dockerfile

In a Docker-based project, the Dockerfile serves as the foundational
asset for building your application's environment.

A Dockerfile is a text file that instructs Docker how to create an image of your
application's environment. An image contains everything you want and
need when running application, such as files, packages, and tools.

The following is the Dockerfile for this project.

```dockerfile
FROM nginx:stable-alpine3.17-slim
WORKDIR /usr/share/nginx/html
COPY . .
```

This Dockerfile defines an image that serves static content using Nginx from an
Alpine Linux base image.

## Summary

In this guide, you explored how to leverage TensorFlow.js with Docker to run a
pre-trained model for face detection in a web application.

Related information:

- [TensorFlow.js website](https://www.tensorflow.org/js)
- [MediaPipe website](https://developers.google.com/mediapipe/)
- [Dockerfile reference](/reference/dockerfile/)
- [Docker CLI reference](/reference/cli/docker/)
- [Docker Blog: Accelerating Machine Learning with TensorFlow.js](https://www.docker.com/blog/accelerating-machine-learning-with-tensorflow-js-using-pretrained-models-and-docker/)