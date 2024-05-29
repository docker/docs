---
title: Develop with containers  
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you how to develop with containers
---

{{< youtube-embed 9JnqOmJ96ds >}}

## Explanation

Now that you have Docker Desktop installed, you’re are ready to do some application development. Specifically, you will do the following:

1. Clone and start a development project
2. Make changes to the backend and frontend
3. See the changes immediately

Enough talk... time to get started.

## Try it out

## Start the project

1. To get started, either clone or [download the project as a ZIP file](https://github.com/docker/getting-started-todo-app/archive/refs/heads/main.zip) to your local machine.

    ```console
    $ git clone https://github.com/docker/getting-started-todo-app
    ```

    And after the project is cloned, navigate into the new directory created by the clone:

    ```console
    $ cd getting-started-todo-app
    ```

2. Once you have the project, start the development environment using Docker Compose.


    To start the project using the CLI, run the following command:

   ```console
   $ docker compose watch
   ```

   You will see an output that shows container images being pulled down, containers starting, and more. Don't worry if you don't understand it all at this point. But, within a moment or two, things should stabilize and finish.


3. Open your browser to [http://localhost](http://localhost) to see the application up and running. It may take a few minutes for the app to run. The app is a simple to-do application, so feel free to add an item or two, mark some as done, or even delete an item.

    ![Screenshot of the getting started to-do app after its first launch](images/develop-getting-started-app-first-launch.webp)


### What's in the environment?

Now that the environment is up and running, what's actually in it? At a high-level, there are several containers (or processes) that each serve a specific need for the application:

- React frontend - a Node container that's running the React dev server, using [Vite](https://vitejs.dev/)
- Node backend - the backend provides an API that provides the ability to retrieve, create, and delete to-do items
- MySQL database* - a database to store the list items
- phpMyAdmin - a web-based interface to interact with the database that is accessible at [http://db.localhost](http://db.localhost)
- Traefik proxy - [Traefik](https://traefik.io/traefik/) is an application proxy that routes requests to the right service. It sends all requests for `localhost/api/*` to the backend, requests for `localhost/*` to the frontend, and then requests for `db.localhost` to phpMyAdmin. This provides the ability to access all applications using port 80 (instead of different ports for each service).

With this environment, you as the developer don’t need to install or configure any services, populate a database schema, configure database credentials, or anything. You only need Docker Desktop. The rest just works.


## Make changes to the app

With this environment up and running, you’re ready to make a few changes to the application and see how Docker helps provide a fast feedback loop.

### Change the greeting

The greeting at the top of the page is populated by an API call at `/api/greeting`. Currently, it always returns "Hello world!". You’ll now modify it to return one of three randomized messages (that you'll get to choose).

1. Open the `backend/src/routes/getGreeting.js` file. This file provides the handler for the API endpoint.

2. Modify the variable at the top to an array of greetings. Feel free to use the following modifications or any customize it to your own liking.

    ```js {linenos=table,hl_lines=["1-5"],linenostart=1}
    const GREETINGS = [
        "Whalecome!",
        "All hands on deck!",
        "Charting the course ahead!",
    ];

    module.exports = async (req, res) => {
    ...
    ```

3. Now, update the endpoint to send a random greeting from this list by making the following change:

    ```js {linenos=table,hl_lines=[3],linenostart=7}
    module.exports = async (req, res) => {
        res.send({
            greeting: GREETINGS[ Math.floor( Math.random() * GREETINGS.length )],
        });
    };
    ```

4. If you haven't done so yet, save the file. If you refresh your browser, you should see a new greeting. If you keep refreshing, you should see all of the messages appear.

    ![Screenshot of the to-do app with a new greeting](images/develop-app-with-greetings.webp)


### Change the placeholder text

When you look at the app, you'll see the placeholder text is simply "New Item". You’ll now make that a little more descriptive and fun. You’ll also make a few changes to the styling of the app too.

1. Open the `client/src/components/AddNewItemForm.jsx` file. This provides the component to add a new item to the to-do list.

2. Modify the `placeholder` attribute of the `Form.Control` element to whatever you'd like to display.

    ```js {linenos=table,hl_lines=[5],linenostart=33}
    <Form.Control
        value={newItem}
        onChange={(e) => setNewItem(e.target.value)}
        type="text"
        placeholder="What do you need to do?"
        aria-label="New item"
    />
    ```

3. Save the file and go back to your browser. You should see the change already hot-reloaded into your browser. If you don't like it, feel free to tweak it until it looks just right.

![Screenshot of the to-do app with an updated placeholder in the add item text field"](images/develop-app-with-updated-placeholder.webp)


### Change the background color

Before you consider the application finalized, you need to make the colors better.

1. Open the `client/src/index.scss` file.

2. Adjust the `background-color` attribute to any color you'd like. The provided snippet is a soft blue to go along with Docker's nautical theme.

    If you're using an IDE, you can pick a color using the integrated color pickers. Otherwise, feel free to use an online [Color Picker](https://www.w3schools.com/colors/colors_picker.asp).

    ```css {linenos=table,hl_lines=2,linenostart=3}
    body {
        background-color: #99bbff;
        margin-top: 50px;
        font-family: 'Lato';
    }
    ```

    Each save should let you see the change immediately in the browser. Keep adjusting it until it's the perfect setup for you.


    ![Screenshot of the to-do app with a new placeholder and background color"](images/develop-app-with-updated-client.webp)

And with that, you're done. Congrats on updating your website. 


## Recap

Before you move on, take a moment and reflect on what happened here. Within a few moments, you were able to:

- Start a complete development project with zero installation effort. The containerized environment provided the development environment, ensuring you have everything you need. You didn't have to install Node, MySQL, or any of the other dependencies directly on your machine. All you needed was Docker Desktop and a code editor.

- Make changes and see them immediately. This was made possible because 1) the processes running in each container are watching and responding to file changes and 2) the files are shared with the containerized environment.

Docker Desktop enables all of this and so much more. Once you start thinking with containers, you can create almost any environment and easily share it with your team.

## Next steps

Now that the application has been updated, you’re ready to learn about packaging it as a container image and pushing it to a registry, specifically Docker Hub.

{{< button text="Build and push your first image" url="build-and-push-first-image" >}}

