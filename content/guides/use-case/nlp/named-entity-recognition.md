---
title: Build a named entity recognition app
keywords: nlp, natural language processing, named entity recognition, python, spacy, ner
description: Learn how to build and run a named entity recognition application using Python, spaCy, and Docker.
---

## Overview

This guide walks you through building and running a named entity recognition
(NER) application. You'll build the application using Python with
spaCy, and then set up the environment and run the application using Docker.

The application processes input text to identify and print named entities, like people, organizations, or locations.

## Prerequisites

* You have installed the latest version of [Docker Desktop](../../../get-docker.md). Docker adds new features regularly and some parts of this guide may work only with the latest version of Docker Desktop.
* You have a [Git client](https://git-scm.com/downloads). The examples in this section use a command-line based Git client, but you can use any client.

## Get the sample application

1. Open a terminal, and clone the sample application's repository using the
   following command.

   ```console
   $ git clone https://github.com/harsh4870/Docker-NLP.git
   ```

2. Verify that you cloned the repository.

   You should see the following files in your `Docker-NLP` directory.

   ```text
   01_sentiment_analysis.py
   02_name_entity_recognition.py
   03_text_classification.py
   04_text_summarization.py
   05_language_translation.py
   entrypoint.sh
   requirements.txt
   Dockerfile
   README.md
   ```

## Explore the application code

The source code for the name recognition application is in the `Docker-NLP/02_name_entity_recognition.py` file. Open `02_name_entity_recognition.py` in a text or code editor to explore its contents in the following steps.

1. Import the required libraries.

   ```python
   import spacy
   ```
   
   This line imports the `spaCy` library. `spaCy` is a popular library in Python
   used for natural language processing (NLP).

2. Load the language model.

   ```python
   nlp = spacy.load("en_core_web_sm")
   ```
   
   Here, the `spacy.load` function loads a language model. The `en_core_web_sm`
   model is a small English language model. You can use this model for various
   NLP tasks, including tokenization, part-of-speech tagging, and named entity
   recognition.

3. Specify the main execution block.

   ```python
   if __name__ == "__main__":
   ```

   This Python idiom ensures that the following code block runs only if this
   script is the main program. It provides flexibility, allowing the script to
   function both as a standalone program and as an imported module.

4. Create an infinite loop for continuous input.

   ```python
      while True:
   ```

   This while loop runs indefinitely until it's explicitly broken. It lets
   the user continuously enter text for entity recognition until they decide
   to exit.

5. Get user input.

   ```python
   input_text = input("Enter the text for entity recognition (type 'exit' to end): ")
   ```

   This line prompts the user to enter text. The program will then perform entity recognition on this text.

6. Define an exit condition.

   ```python
   if input_text.lower() == 'exit':
      print("Exiting...")
      break
   ```

   If the user types something, the program converts the input to lowercase and
   compares it to `exit`. If they match, the program prints **Exiting...** and
   breaks out of the while loop, effectively ending the program.

7. Perform named entity recognition.

   ```python
   doc = nlp(input_text)

   for ent in doc.ents:
      print(f"Entity: {ent.text}, Type: {ent.label_}")
   ```

   - `doc = nlp(input_text)`: Here, the nlp model processes the user-input text. This creates a Doc object which contains various NLP attributes, including identified entities.
   - `for ent in doc.ents:`: This loop iterates over the entities found in the text.
   - `print(f"Entity: {ent.text}, Type: {ent.label_}")`: For each entity, it prints the entity text and its type (like PERSON, ORG, or GPE).


8. Create `requirements.txt`.

   The sample application already contains the `requirements.txt` file to specify the necessary packages that the application imports. Open `requirements.txt` in a code or text editor to explore its contents.

   ```text
   # 02 named_entity_recognition
   spacy==3.7.2

   ...
   ```

   Only the `spacy` package is required for the named recognition application.

## Explore the application environment

You'll use Docker to run the application in a container. Docker lets you
containerize the application, providing a consistent and isolated environment
for running it. This means the application will operate as intended within its
Docker container, regardless of the underlying system differences.

To run the application in a container, a Dockerfile is required. A Dockerfile is
a text document that contains all the commands you would call on the command
line to assemble an image. An image is a read-only template with instructions
for creating a Docker container.

The sample application already contains a `Dockerfile`. Open the `Dockerfile` in a code or text editor to explore its contents.

The following steps explain each part of the `Dockerfile`. For more details, see the [Dockerfile reference](/reference/dockerfile/).

1. Specify the base image.

   ```dockerfile
   FROM python:3.8-slim
   ```

   This command sets the foundation for the build. `python:3.8-slim` is a
   lightweight version of the Python 3.8 image, optimized for size and speed.
   Using this slim image reduces the overall size of your Docker image, leading
   to quicker downloads and less surface area for security vulnerabilities. This
   is particularly useful for a Python-based application where you might not
   need the full standard Python image.

2. Set the working directory.

   ```dockerfile
   WORKDIR /app
   ```

   `WORKDIR` sets the current working directory within the Docker image. By
   setting it to `/app`, you ensure that all subsequent commands in the
   Dockerfile (like `COPY` and `RUN`) are executed in this directory. This also
   helps in organizing your Docker image, as all application-related files are
   contained in a specific directory.

3. Copy the requirements file into the image.

   ```dockerfile
   COPY requirements.txt /app
   ```

   The `COPY` command transfers the `requirements.txt` file from
   your local machine into the Docker image. This file lists all Python
   dependencies required by the application. Copying it into the container
   lets the next command (`RUN pip install`) install these dependencies
   inside the image environment.

4. Install the Python dependencies in the image.

   ```dockerfile
   RUN pip install --no-cache-dir -r requirements.txt
   ```

   This line uses `pip`, Python's package installer, to install the packages
   listed in `requirements.txt`. The `--no-cache-dir` option disables
   the cache, which reduces the size of the Docker image by not storing the
   unnecessary cache data.

5. Run additional commands.

   ```dockerfile
   RUN python -m spacy download en_core_web_sm
   ```

   This step is specific to NLP applications that require the spaCy library. It downloads the `en_core_web_sm` model, which is a small English language model for spaCy.

6. Copy the application code into the image.

   ```dockerfile
   COPY *.py /app
   COPY entrypoint.sh /app
   ```

   These commands copy your Python scripts and the `entrypoint.sh` script into
   the image's `/app` directory. This is crucial because the container needs
   these scripts to run the application. The `entrypoint.sh` script is
   particularly important as it dictates how the application starts inside the
   container.

7. Set permissions for the `entrypoint.sh` script.

   ```dockerfile
   RUN chmod +x /app/entrypoint.sh
   ```

   This command modifies the file permissions of `entrypoint.sh`, making it
   executable. This step is necessary to ensure that the Docker container can
   run this script to start the application.

8. Set the entry point.

   ```dockerfile
   ENTRYPOINT ["/app/entrypoint.sh"]
   ```

    The `ENTRYPOINT` instruction configures the container to run `entrypoint.sh`
    as its default executable. This means that when the container starts, it
    automatically executes the script.
   
   You can explore the `entrypoint.sh` script by opening it in a code or text
   editor. As the sample contains several applications, the script lets you
   specify which application to run when the container starts.

## Run the application

To run the application using Docker:

1. Build the image.

   In a terminal, run the following command inside the directory of where the `Dockerfile` is located.

   ```console
   $ docker build -t basic-nlp .
   ```

   The following is a break down of the command:

   - `docker build`: This is the primary command used to build a Docker image
     from a Dockerfile and a context. The context is typically a set of files at
     a specified location, often the directory containing the Dockerfile.
   - `-t basic-nlp`: This is an option for tagging the image. The `-t` flag
     stands for tag. It assigns a name to the image, which in this case is
     `basic-nlp`. Tags are a convenient way to reference images later,
     especially when pushing them to a registry or running containers.
   - `.`: This is the last part of the command and specifies the build context.
     The period (`.`) denotes the current directory. Docker will look for a
     Dockerfile in this directory. The build context (the current directory, in
     this case) is sent to the Docker daemon to enable the build. It includes
     all the files and subdirectories in the specified directory.

   For more details, see the [docker build CLI reference](/reference/cli/docker/buildx/build/).

   Docker outputs several logs to your console as it builds the image. You'll
   see it download and install the dependencies. Depending on your network
   connection, this may take several minutes. Docker does have a caching
   feature, so subsequent builds can be faster. The console will
   return to the prompt when it's complete.

2. Run the image as a container.

   In a terminal, run the following command.

   ```console
   $ docker run -it basic-nlp 02_name_entity_recognition.py
   ```

   The following is a break down of the command:

   - `docker run`: This is the primary command used to run a new container from
     a Docker image.
   - `-it`: This is a combination of two options:
      - `-i` or `--interactive`: This keeps the standard input (STDIN) open even
        if not attached. It lets the container remain running in the
        foreground and be interactive.
      - `-t` or `--tty`: This allocates a pseudo-TTY, essentially simulating a
        terminal, like a command prompt or a shell. It's what lets you
        interact with the application inside the container.
   - `basic-nlp`: This specifies the name of the Docker image to use for
     creating the container. In this case, it's the image named `basic-nlp` that
     you created with the `docker build` command.
   - `02_name_entity_recognition.py`: This is the script you want to run inside
     the Docker container. It gets passed to the `entrypoint.sh` script, which
     runs it when the container starts.
     
   For more details, see the [docker run CLI reference](/reference/cli/docker/container/run/).


   > [!NOTE]
   >
   > For Windows users, you may get an error when running the container. Verify
   > that the line endings in the `entrypoint.sh` are `LF` (`\n`) and not `CRLF` (`\r\n`),
   > then rebuild the image. For more details, see [Avoid unexpected syntax errors, use Unix style line endings for files in containers](/desktop/troubleshoot/topics/#avoid-unexpected-syntax-errors-use-unix-style-line-endings-for-files-in-containers).

   You will see the following in your console after the container starts.

   ```console
   Enter the text for entity recognition (type 'exit' to end):
   ```

3. Test the application.

   Enter some information to get the named entity recognition.

   ```console
   Enter the text for entity recognition (type 'exit' to end): Apple Inc. is planning to open a new store in San Francisco. Tim Cook is the CEO of Apple.
   
   Entity: Apple Inc., Type: ORG
   Entity: San Francisco, Type: GPE
   Entity: Tim Cook, Type: PERSON
   Entity: Apple, Type: ORG
   ```

## Summary

This guide demonstrated how to build and run a named entity recognition
application. You learned how to build the application using Python with spaCy,
and then set up the environment and run the application using Docker.

Related information:

* [Docker CLI reference](/reference/cli/docker/)
* [Dockerfile reference](/reference/dockerfile/)
* [spaCy](https://spacy.io/)
* [Python documentation](https://docs.python.org/3/)

## Next steps

Explore more [natural language processing guides](./_index.md).
