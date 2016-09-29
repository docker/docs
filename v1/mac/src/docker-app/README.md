# Docker.app

Docker.app is the native Mac port of the Docker engine, but it also comes with a nice set of tools to get the best out of Docker, on OS X. 

## Goal

Docker.app is mainly an integration project. It assembles a bunch of different pieces:

- A status item: a whale icon with a menu, in the menu bar. 
- A graphical interface (currently implemented with [Electron](http://electron.atom.io))
- A preference pane
- A Go agent that will later embed a complete runtime for the Docker platform, but for V1.10 it would be great if it could handle: logs (+[bugsnag](https://bugsnag.com)), and a simple API to get information about the hyperkit VM (status, CPU usage...etc)
- A virtual networking program (vmnetd), required by the hyperkit VM for networking.
- The hyperkit VM, host for the Docker engine.
- A few binaries that have to be installed (docker-cli, docker-compose, docker-machine...)


The purpose of **docker-app2** is to bundle everything in one main application. The hyperkit VM may need to run in its own process, maybe the GUI as well for some time, because it may not be easy to build the Electorn application as a simple window manager library... But we shouldn't have more than 3 running processes:

- com.docker.app
- com.docker.ui (will be merged in com.docker.app at some point)
- com.docker.engine

Docker.app should be an agent application (no dock icon), but it can dynamically transform its process to foreground when showing the GUI (dock icon appears).

The preference pane and binaries are installed on first launch and can all be updated by the main application itself. BUT there's no dedicated installer that has to kill itself when done, the main application handles this.

Except for the GUI, all interactions with the user (alerts, loading bars, prompts...etc) should be handled by the same (small) window that manages an interaction queue. Interaction requests can be added anywhere in the queue, adding them in front displays them right away.

The Go agent is embeded as a library it provides a set of functions, maybe also an API (over ssh). (not sure for v1.10)

The preference pane should just save preferences, it's not communicating with anything. Docker.app checks preferences when launched, but also watches for changes when running.

Docker.app may not handle going from agent to foreground mode in v1.10 if we can't embed the GUI. If the GUI has to remain a separate application, then it will be foreground, and the Docker.app will always run as an agent. (working on something with @jeffdm & @FrenchBen)

The best would be to also include the hyperkit VM has a library, and start it from within the application, but it's ok to keep a dedicated process for it. As long as it's tight to the application, none of these should be running without the other.

The app should be added by default to user's login items. But the user should be able to quit Docker.app, and everything should be running then after that (agent, vm...). Ideally, when typing a docker command, or when launching the GUI directly (if it's a seperate app), Docker.app should be opened (if needed) and display a "Docker starting..." message.

## Language

Docker.app should be written 100% in **Swift**. Everything that is not Swift should be a .a library, or embeded binary.

## Build

It should ALWAYS be possible to build the **Docker.xcworkspace** project, as is, from XCode.

The build scripts should basically build all libraries and binaries, to update them, before building Docker.app itself. It should be much easier to sign the whole package this way.

