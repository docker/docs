---
description: Learn how to use Claude Code with Docker MCP Toolkit to generate production-ready Docker Compose files from natural language using the Docker Hub MCP server.
keywords: mcp, model context protocol, docker, docker desktop, claude code, docker hub, compose, automation
title: Generate Docker Compose Files with Claude Code and Docker MCP Toolkit
summary: |
  This guide shows how to wire Claude Code to the Docker MCP Toolkit so it can search Docker Hub images and generate complete Docker Compose stacks from natural language.
  You’ll enable the Docker Hub MCP server, connect Claude Code, verify MCP access, and create a Node.js + PostgreSQL stack with a conversational prompt.
tags: [ai]
languages: [bash, json, yaml, javascript]
aliases:
  - /guides/use-case/claude-compose-mcp/
params:
  time: 15 minutes
---

This guide introduces how to use Claude Code together with Docker MCP Toolkit so Claude can search Docker Hub in real time and generate a complete `docker-compose.yml` from natural language.

Instead of manually writing YAML or looking for image tags, you describe your stack once — Claude uses the Model Context Protocol (MCP) to query Docker Hub and build a production-ready Compose file.

In this guide, you’ll learn how to:

- Enable Docker MCP Toolkit in Docker Desktop  
- Add the Docker Hub MCP server  
- Connect Claude Code to the MCP Gateway (GUI or CLI)  
- Verify MCP connectivity inside Claude  
- Ask Claude to generate and save a Compose file for a Node.js + PostgreSQL app  
- Deploy it instantly with `docker compose up`  

---

## Use Claude Code and Docker MCP Toolkit to generate a Docker Compose file from natural language


- **Setup**: Enable MCP Toolkit → Add Docker Hub MCP server → Connect Claude Code  
- **Use Claude**: Describe your stack in plain English  
- **Automate**: Claude queries Docker Hub via MCP and builds a complete `docker-compose.yml`  
- **Deploy**: Run `docker compose up` → Node.js + PostgreSQL live on `localhost:3000`  
- **Benefit**: Zero YAML authoring. Zero image searching. Describe once → Claude builds it.

**Estimated time**: ~15 minutes

---

## 1. What you are building

The goal is simple: use Claude Code together with the Docker MCP Toolkit to search Docker Hub images and generate a complete Docker Compose file for a Node.js and PostgreSQL setup. 

The Model Context Protocol (MCP) bridges Claude Code and Docker Desktop, giving Claude real-time access to Docker's tools. Instead of context-switching between Docker, terminal commands, and YAML editors, you describe your requirements once and Claude handles the infrastructure details.

**Why this matters:** This pattern scales to complex multi-service setups, database migrations, networking, security policies — all through conversational prompts.

---

## 2. Prerequisites

Make sure you have:

- Docker Desktop installed
- Docker Desktop updated with MCP Toolkit support
- Claude Code installed

---

## 3. Enable Docker MCP Toolkit in Docker Desktop

1. Open Docker Desktop  
1. Go to **Settings**  
1. Open **Beta features**  
1. Turn on **Enable Docker MCP Toolkit**  
1. Select **Apply and restart**

![Enable MCP](./images/enable_mcp.png "Enable MCP")


You should now see **MCP Toolkit** in the sidebar.

---

## 4. Install the Docker Hub MCP server

1. Open **Docker Desktop**  
1. Select **MCP Toolkit**  
1. Go to the **Catalog** tab  
1. Search for **Docker Hub**  
1. Select the **Docker Hub MCP server**  
1. Select **+ Add**

![Docker Hub](./images/catalog_docker_hub.png "Docker Hub")

Public images work without credentials. For private repositories, you can add your Docker Hub username and token later.

![Docker Hub Secrets](./images/dockerhub_secrets.png "Docker Hub Secrets")


---

## 5. Connect Claude Code to Docker MCP Toolkit

You can connect from Docker Desktop or using the CLI.

### Option A. Connect with Docker Desktop

1. Open **MCP Toolkit**  
1. Go to the **Clients** tab  
1. Locate **Claude Code**  
1. Select **Connect**

![Docker Connection](./images/docker-connect-claude.png)

### Option B. Connect using the CLI

```bash
claude mcp add MCP_DOCKER -s user -- docker mcp gateway run
```

Verify the connection:

```bash
claude mcp list
```

You should see:
![MCP List](./images/claude-mcp-list.png "MCP List")

```bash
MCP_DOCKER  docker mcp gateway run  user
```

---

## 6. Verify MCP servers inside Claude Code

1. Navigate to your project folder:

```bash
cd /path/to/project
```

1. Start Claude Code:

```bash
claude code
```

![Claude Home](./images/claude-home.png)

1. In the input box, type:

```text
/mcp
```

You should now see:

- The MCP gateway (for example `MCP_DOCKER`)
- Tools provided by the Docker Hub MCP server

![mcp-docker](./images/mcp-servers.png)

If not, restart Claude Code or check Docker Desktop to confirm the connection.

---

## 7. Ask Claude Code to design your Docker Compose stack

Paste this message into Claude Code:

```text
Using the Docker Hub MCP server:

Search Docker Hub for an official Node.js image and a PostgreSQL image.
Choose stable, commonly used tags such as the Node LTS version and a recent major Postgres version.

Generate a Docker Compose file (`docker-compose.yml`) with:
- app: running on port 3000
- db: running on port 5432 using a named volume

Include:
- Environment variables for Postgres
- A shared bridge network
- Healthchecks where appropriate
```

Claude will search images through MCP and generate a Compose file for you.

---

## 8. Save the generated Docker Compose file

Tell Claude:

```text
Save the final Docker Compose file (docker-compose.yml) into the current project directory.
```

Verify it:

```bash
cat docker-compose.yml
```

You should see something like this:

```yaml
services:
  app:
    image: node:<tag>
    working_dir: /usr/src/app
    volumes:
      - .:/usr/src/app
    ports:
      - "3000:3000"
    depends_on:
      - db
    networks:
      - app-net

  db:
    image: postgres:<tag>
    environment:
      POSTGRES_USER: example
      POSTGRES_PASSWORD: example
      POSTGRES_DB: appdb
    volumes:
      - db-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - app-net

volumes:
  db-data:

networks:
  app-net:
    driver: bridge
```

---

## 9. Add a basic Node.js app

Inside your project folder, create a folder named `app`:

```bash
mkdir app
cd app
npm init -y
npm install express
```

Create `index.js`:

```js
const express = require("express");
const app = express();

app.get("/", (req, res) => {
  res.send("Node.js, Docker, and MCP Toolkit are working together!");
});

app.listen(3000, () => {
  console.log("Server running on port 3000");
});
```

Add a start script to `package.json`:

```json
"scripts": {
  "start": "node index.js"
}
```

---

## 10. Run the Docker Compose stack

From your project root:

```bash
docker compose up
```

Docker will:

- Pull the Node and Postgres images selected through Docker Hub MCP  
- Create networks and volumes  
- Start the containers  

Open your browser:

```text
http://localhost:3000
```
![Local Host](./images/Localhost.png)

Your Node.js app should now be running.

---

## Conclusion

By combining Claude Code with the Docker MCP Toolkit, Docker Desktop, and the Docker Hub MCP server, you can describe your stack in natural language and let MCP handle the details. This removes context switching and replaces it with a smooth, guided workflow powered by model context protocol integrations.

---

### Next steps

- Explore the 220+ MCP servers available in the [Docker MCP catalog](https://hub.docker.com/mcp)
- Connect Claude Code to your databases, internal APIs, and team tools  
- Share your MCP setup with your team so everyone works consistently  

The future of development is not about switching between tools. It is about tools working together in a simple, safe, and predictable way. The Docker MCP Toolkit brings that future into your everyday workflow.



## Learn more

- **[Explore the MCP Catalog](https://hub.docker.com/mcp):** Discover containerized, security-hardened MCP servers  
- **[Get started with MCP Toolkit in Docker Desktop](https://hub.docker.com/open-desktop?url=https://open.docker.com/dashboard/mcp):** Requires version 4.48 or newer to launch automatically  
- **[Read the MCP Horror Stories series](https://www.docker.com/blog/mcp-horror-stories-the-supply-chain-attack/):** Learn about common MCP security pitfalls and how to avoid them  
