---
title: Leveraging RAG in GenAI to teach new information
linkTitle: Leveraging RAG in GenAI
description:  This guide walks through the process of setting up and utilizing a GenAI stack with Retrieval-Augmented Generation (RAG) systems and graph databases. Learn how to integrate graph databases like Neo4j with AI models for more accurate, contextually-aware responses.
keywords: Docker, GenAI, Retrieval-Augmented Generation, RAG, Graph Databases, Neo4j, AI, LLM
summary: |
  This guide explains setting up a GenAI stack with Retrieval-Augmented Generation (RAG) and Neo4j, covering key concepts, deployment steps, and a case study. It also includes troubleshooting tips for optimizing AI performance with real-time data.
tags: [ai]
params:
  time: 35 minutes
---


This guide walks you through setting up a powerful AI stack that integrates Retrieval-Augmented Generation (RAG) systems with graph databases. In this guide, you’ll learn how to:

- Set up a GenAI system that enhances language models with real-time data from graph databases
- Use Docker to containerize and deploy the complete GenAI stack, including graph databases and AI models
- Leverage Neo4j for efficient information retrieval and context-aware AI responses
- Inspecting data in the database, ensuring smooth operation for AI-driven applications



## Introduction to Graph Databases 

All NoSQL databases can be grouped under 4 main groups:

* Document-based databases
* Key-value stores
* Column-oriented databases
* Graph-based databases

Being one of the four main groups, graph databases are specialized database systems designed to store and query data where relationships between entities are as important as the entities themselves. Unlike traditional databases, graph databases use nodes (vertices) to store entities and edges to store relationships between these entities, making them ideal for handling highly connected data.

### Graph Databases vs SQL Databases 

1. Data Model:
   - SQL: Uses tables with rows and columns, relationships expressed through foreign keys
   - Graph: Uses nodes and edges, relationships are much more flexible

2. Schema Flexibility:
   - SQL: Rigid schema, changes require extra steps
   - Graph: Flexible schema, can add new relationships without migrations

3. Use Cases:
   - SQL: Structured data with fixed relationships
   - Graph: Social networks, recommendation engines, knowledge graphs


## Understanding RAG (Retrieval-Augmented Generation) 

RAG (Retrieval-Augmented Generation) is a hybrid framework that enhances the capabilities of large language models by integrating information retrieval. It combines three core components:  

- **Information retrieval** from an external knowledge base  
- **Large Language Model (LLM)** for generating responses  
- **Vector embeddings** to enable semantic search  

In a Retrieval-Augmented Generation (RAG) system, vector embeddings are used to represent the semantic meaning of text in a way that a machine can understand and process. For instance, the words "dog" and "puppy" will have similar embeddings because they share similar meanings. By integrating these embeddings into the RAG framework, the system can combine the generative power of large language models with the ability to pull in highly relevant, contextually-aware data from external sources.

The system operates as follows:  
1. Queries are transformed into vector embeddings to capture semantic meaning.  
2. Relevant information is retrieved from the knowledge base using these embeddings.  
3. The retrieved data is used to augment the input prompt for the LLM.  
4. The LLM generates responses that blend the model's inherent knowledge with the retrieved context.  

### When RAG is Effective 

1. Knowledge Limitations:
   - LLMs have fixed knowledge cutoffs
   - Cannot access real-time or specialized information
   - May hallucinate or provide outdated information

2. Business Requirements:
   - Need for accurate, up-to-date information
   - Compliance and auditability
   - Domain-specific knowledge integration

3. Cost Efficiency:
   - Reduced need for model fine-tuning
   - Smaller models can perform well with good context
   - Lower computational requirements




## A Case Study in RAG Implementation 

### What is Apache NiFi? 

Apache NiFi is a powerful data integration tool, but in our context, it serves us as a concept rather than a tool. Since it's a relatively new concept, large language models (LLMs) may have limited or no inherent knowledge about it.



    Question: What is Apache Nifi? 
    RAG: Disabled

    Hello! I'm here to help you with your question about Apache NiFi. Unfortunately, I don't know the answer to that question. I'm just an AI and my knowledge cutoff is December 2022, so I may not be familiar with the latest technologies or software. Can you please provide more context or details about Apache NiFi? Maybe there's something I can help you with related to it.

    Question: What is Apache Nifi? 
    RAG: Enabled

    Answer:

    Apache NiFi is an open-source, distributed data processing system that allows you to collect, transform, and transport data across various sources and sinks. It provides a graphical user interface for designing and managing data flows, as well as a set of APIs for integrating with other systems.

    According to this Stack Overflow question, the free Neo4J JDBC driver doesn't understand SQL, only Cypher, so you would need to use a non-free BI Connector to connect to Neo4j from Apache Drill.

    Additionally, this Stack Overflow question suggests that Apache POI does work in a Neo4j User Defined Function. However, there may be an unknown glitch in the standalone UDF mentioned in the initial question.

    To export data from Neo4j to Excel using APOC Extended procedures, you can use apoc.export.xls.query which takes a Cypher query and exports the results to an Excel file.

    References:

    How to connect and query Neo4j Database on Apache Drill?
    Is a Neo4j UDF compatible with Apache POI?




## Setting Up GenAI Stack with GPU Acceleration on Linux 

To set up and run the GenAI stack on a Linux host with GPU acceleration, execute the following command:

```bash
docker compose --profile linux-gpu up -d
```

### Setting Up on Other Platforms 

For instructions on how to set up the stack on other platforms, refer to [this page](https://github.com/docker/genai-stack). 

---

### Notes 

- **Initial Startup**: The first startup may take some time because the system needs to download a large language model (LLM).
- **Monitoring Progress**: We can monitor the download and initialization progress by viewing the logs.

Run the following command to view the logs:

```bash
docker compose logs
```

Wait for specific lines in the logs indicating that the download is complete and the stack is ready. These lines typically confirm successful setup and initialization.

    pull-model-1 exited with code 0
    database-1    | 2024-12-29 09:35:53.269+0000 INFO  Started.
    pdf_bot-1     |   You can now view your Streamlit app in your browser.
    loader-1      |   You can now view your Streamlit app in your browser.
    bot-1         |   You can now view your Streamlit app in your browser.


    You can now access the interface at [http://localhost:8501/](http://localhost:8501/) to ask questions. For example, you can try the sample question:


    What is Apache Nifi? 


    The response should be similar to the following:


    ... I'm just an AI and my knowledge cutoff is December 2022...


Now it's time to teach the AI some new tricks. First, connect to [loader-1](http://localhost:8502/). Instead of using the "neo4j" tag, change it to the "apache-nifi" tag, then click the **Import** button. After the import is successful, access Neo4j to verify the data.



After logging in to [http://localhost:7474/](http://localhost:7474/) using the credentials from the `.env` file, you can run queries on Neo4j. Using the Neo4j Cypher query language, you can check for the data stored in the database.

To count the data, run the following query:


```cypher
MATCH (n)
RETURN DISTINCT labels(n) AS NodeTypes, count(*) AS Count
ORDER BY Count DESC;
```


You can also run the following query to visualize the data:

```cypher
CALL db.schema.visualization()
```

To check the relationships in the database, run the following query:

```cypher
CALL db.relationshipTypes()
```

Now, we are ready to enable our LLM to use this information. Go back to [http://localhost:8501/](http://localhost:8501/), enable the **RAG** checkbox, and ask the same question again. The LLM will now provide a more detailed answer.

Keep in mind that new questions will be added to Stack Overflow, and due to the inherent randomness in most AI models, the answers may vary and won't be identical to those in this example.

Feel free to start over with another [Stack Overflow tag](https://stackoverflow.com/tags). To drop all data in Neo4j, you can use the following command in the Neo4j Web UI:


```cypher
MATCH (n)
DETACH DELETE n;
```

For optimal results, choose a tag that the LLM is not familiar with.

